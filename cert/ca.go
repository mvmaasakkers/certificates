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

	NotBefore time.Time
	NotAfter  time.Time
}

func NewCARequest() *CARequest {
	return &CARequest{
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 0, 1),
	}
}

func (req *CARequest) Validate() error {
	if req.CommonName == "" {
		return ErrorInvalidCommonName
	}

	return nil
}

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

	return name
}

func (req *CARequest) GenerateCA() ([]byte, []byte, error) {
	if err := req.Validate(); err != nil {
		return nil, nil, err
	}

	ca := &x509.Certificate{
		SerialNumber:          big.NewInt(1653),
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
