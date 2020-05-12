package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

// Request is the struct needed to generate a CA or Certificate pair
type Request struct {
	Organization  string
	Country       string
	Province      string
	Locality      string
	StreetAddress string
	PostalCode    string
	CommonName    string

	SerialNumber     *big.Int
	NameSerialNumber string

	SubjectAltNames []string

	NotBefore time.Time
	NotAfter  time.Time

	BitSize int
}

const defaultBitSize = 4096

var bitSizeOptions = []int{2048, 4096}

func isValidBitSizeOption(option int) bool {
	for _, v := range bitSizeOptions {
		if v == option {
			return true
		}
	}

	return false
}

// NewRequest will create a new Request struct and set the NotBefore to now and the NotAfter to one day from now
func NewRequest() *Request {
	return &Request{
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 0, 1),
		BitSize:   defaultBitSize,
	}
}

// Validate will check the validity of the Request object
//
// The checks are:
// - A Common Name is mandatory
// - If a list of SubjectAltNames is given, none of them can be empty
func (req *Request) Validate() error {
	if req.CommonName == "" {
		return ErrorInvalidCommonName
	}

	if req.BitSize == 0 {
		req.BitSize = defaultBitSize
	}

	if !isValidBitSizeOption(req.BitSize) {
		return ErrorInvalidBitSize
	}

	for _, n := range req.SubjectAltNames {
		if n == "" {
			return ErrorInvalidSubjectAltName
		}
	}

	return nil
}

// GetPKIXName extracts the Request object into a PKIX Name format for usage in constructing the certificate
// The NameSerialNumber is used as pkix.Name.SerialNumber here (if given).
func (req *Request) GetPKIXName() pkix.Name {
	name := pkix.Name{}

	if req.Organization != "" {
		name.Organization = []string{req.Organization}
	}

	if req.Country != "" {
		name.Country = []string{req.Country}
	}

	if req.Province != "" {
		name.Province = []string{req.Province}
	}

	if req.Locality != "" {
		name.Locality = []string{req.Locality}
	}

	if req.StreetAddress != "" {
		name.StreetAddress = []string{req.StreetAddress}
	}

	if req.PostalCode != "" {
		name.PostalCode = []string{req.PostalCode}
	}

	if req.Locality != "" {
		name.Locality = []string{req.Locality}
	}

	if req.CommonName != "" {
		name.CommonName = req.CommonName
	}

	if req.NameSerialNumber != "" {
		name.SerialNumber = req.NameSerialNumber
	}

	return name
}

// ReadCSR reads csr into a x509.CertificateRequest and converts it into a Request
func ReadCSR(csrFile []byte) (*Request, error) {
	block, _ := pem.Decode(csrFile)
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return nil, fmt.Errorf("failed to decode PEM block containing certificate request")
	}

	csr, err :=  x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, err
	}

	request := NewRequest()
	if len(csr.Subject.Organization) > 0 {
		request.Organization = csr.Subject.Organization[0]
	}
	if len(csr.Subject.Country) > 0 {
		request.Country = csr.Subject.Country[0]
	}
	if len(csr.Subject.Province) > 0 {
		request.Province = csr.Subject.Province[0]
	}
	if len(csr.Subject.Locality) > 0 {
		request.Locality = csr.Subject.Locality[0]
	}
	if len(csr.Subject.StreetAddress) > 0 {
		request.StreetAddress = csr.Subject.StreetAddress[0]
	}
	if len(csr.Subject.PostalCode) > 0 {
		request.PostalCode = csr.Subject.PostalCode[0]
	}
	request.CommonName = csr.Subject.CommonName
	request.NameSerialNumber = csr.Subject.SerialNumber
	request.SubjectAltNames = csr.DNSNames

	if request.SerialNumber == nil {
		randInt, err := GenerateRandomBigInt()
		if err != nil {
			return nil, err
		}

		request.SerialNumber = randInt
	}

	return request, nil
}

// GenerateCertificate will generate a signed certificate pair and will return certificate, key and a possible error
// The Generated key will be in RSA format and has a bit size of 4096 and output of the Certificate
// and Key will be returned in PEM format as bytes.
//
// The certificate will be signed by the given CA Certificate pair (caCrt and caKey). Validity of the CA Certificate
// pair is checked.
func GenerateCertificate(req *Request, caCrt []byte, caKey []byte) ([]byte, []byte, error) {
	if err := req.Validate(); err != nil {
		return nil, nil, err
	}

	catls, err := tls.X509KeyPair(caCrt, caKey)
	if err != nil {
		return nil, nil, err
	}

	ca, err := x509.ParseCertificate(catls.Certificate[0])
	if err != nil {
		return nil, nil, err
	}

	if req.SerialNumber == nil {
		randInt, err := GenerateRandomBigInt()
		if err != nil {
			return nil, nil, err
		}

		req.SerialNumber = randInt
	}

	cert := &x509.Certificate{
		SerialNumber: req.SerialNumber,
		Subject:      req.GetPKIXName(),
		NotBefore:    req.NotBefore,
		NotAfter:     req.NotAfter,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	if len(req.SubjectAltNames) > 0 {
		cert.DNSNames = req.SubjectAltNames
	}

	priv, err := rsa.GenerateKey(rand.Reader, req.BitSize)
	if err != nil {
		return nil, nil, err
	}
	pub := &priv.PublicKey

	certB, err := x509.CreateCertificate(rand.Reader, cert, ca, pub, catls.PrivateKey)
	if err != nil {
		return nil, nil, err
	}

	pemOut := bytes.NewBuffer([]byte{})
	if err := pem.Encode(pemOut, &pem.Block{Type: "CERTIFICATE", Bytes: certB}); err != nil {
		return nil, nil, err
	}

	pemKeyOut := bytes.NewBuffer([]byte{})
	if err := pem.Encode(pemKeyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return nil, nil, err
	}

	return pemOut.Bytes(), pemKeyOut.Bytes(), nil
}

// GenerateCA will generate a CA certificate pair and will return certificate, key and a possible error
// The Generated key will be in RSA format and has a bit size of 4096 and output of the Certificate
// and Key will be returned in PEM format as bytes.
func GenerateCA(req *Request) ([]byte, []byte, error) {
	if err := req.Validate(); err != nil {
		return nil, nil, err
	}

	if req.SerialNumber == nil {
		randInt, err := GenerateRandomBigInt()
		if err != nil {
			return nil, nil, err
		}

		req.SerialNumber = randInt
	}

	ca := &x509.Certificate{
		SerialNumber:          req.SerialNumber,
		Subject:               req.GetPKIXName(),
		NotBefore:             req.NotBefore,
		NotAfter:              req.NotAfter,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}
	pub := &priv.PublicKey
	caB, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		return nil, nil, err
	}

	crtB := bytes.NewBuffer([]byte{})
	if err := pem.Encode(crtB, &pem.Block{Type: "CERTIFICATE", Bytes: caB}); err != nil {
		return nil, nil, err
	}

	keyB := bytes.NewBuffer([]byte{})
	if err := pem.Encode(keyB, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return nil, nil, err
	}

	return crtB.Bytes(), keyB.Bytes(), nil
}
