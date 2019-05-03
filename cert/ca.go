package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

type CARequest struct {
	Organization  string
	Country       string
	Province      string
	Locality      string
	StreetAddress string
	PostalCode    string
	CommonName    string

	SerialNumber     *big.Int
	NameSerialNumber string

	NotBefore time.Time
	NotAfter  time.Time
}

// NewCARequest will create a new CARequest struct and set the NotBefore to now and the NotAfter to one day from now
func NewCARequest() *CARequest {
	return &CARequest{
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 0, 1),
	}
}

// Validate will check the validity of the CARequest object
func (req *CARequest) Validate() error {
	if req.CommonName == "" {
		return ErrorInvalidCommonName
	}

	return nil
}

// GetPKIXName extracts the CARequest object into a PKIX Name format for usage in constructing the certificate
func (req *CARequest) GetPKIXName() pkix.Name {
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

// GenerateCA will generate a CA certificate pair and will return certificate, key and a possible error
func (req *CARequest) GenerateCA() ([]byte, []byte, error) {
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
