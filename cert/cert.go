package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

// CertRequest is the struct needed to generate a CA or Certificate pair
type CertRequest struct {
	Organization    string
	Country         string
	Province        string
	Locality        string
	StreetAddress   string
	PostalCode      string
	CommonName      string

	SerialNumber    *big.Int
	NameSerialNumber string

	SubjectAltNames []string

	NotBefore time.Time
	NotAfter  time.Time
}

// NewCertRequest will create a new CertRequest struct and set the NotBefore to now and the NotAfter to one day from now
func NewCertRequest() *CertRequest {
	return &CertRequest{
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 0, 1),
	}
}

// Validate will check the validity of the CertRequest object
//
// The checks are:
// - A Common Name is mandatory
// - If a list of SubjectAltNames is given, none of them can be empty
func (req *CertRequest) Validate() error {
	if req.CommonName == "" {
		return ErrorInvalidCommonName
	}

	for _, n := range req.SubjectAltNames {
		if n == "" {
			return ErrorInvalidSubjectAltName
		}
	}

	return nil
}

// GetPKIXName extracts the CertRequest object into a PKIX Name format for usage in constructing the certificate
// The NameSerialNumber is used as pkix.Name.SerialNumber here (if given).
func (req *CertRequest) GetPKIXName() pkix.Name {
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

// GenerateCertificate will generate a signed certificate pair and will return certificate, key and a possible error
// The Generated key will be in RSA format and has a bit size of 4096 and output of the Certificate
// and Key will be returned in PEM format as bytes.
//
// The certificate will be signed by the given CA Certificate pair (caCrt and caKey). Validity of the CA Certificate
// pair is checked.
func GenerateCertificate(req *CertRequest, caCrt []byte, caKey []byte) ([]byte, []byte, error) {
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

	priv, err := rsa.GenerateKey(rand.Reader, 4069)
	if err != nil {
		return nil, nil, err
	}
	pub := &priv.PublicKey

	cert_b, err := x509.CreateCertificate(rand.Reader, cert, ca, pub, catls.PrivateKey)

	pemOut := bytes.NewBuffer([]byte{})
	if err := pem.Encode(pemOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert_b}); err != nil {
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
func GenerateCA(req *CertRequest) ([]byte, []byte, error) {
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
	ca_b, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		return nil, nil, err
	}

	crtB := bytes.NewBuffer([]byte{})
	if err := pem.Encode(crtB, &pem.Block{Type: "CERTIFICATE", Bytes: ca_b}); err != nil {
		return nil, nil, err
	}

	keyB := bytes.NewBuffer([]byte{})
	if err := pem.Encode(keyB, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return nil, nil, err
	}

	return crtB.Bytes(), keyB.Bytes(), nil
}
