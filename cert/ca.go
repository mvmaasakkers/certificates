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

func (req *CARequest) GenerateCA() ([]byte, []byte, error) {
	if err := req.Validate(); err != nil {
		return nil, nil, err
	}

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization:  []string{req.Organization},
			Country:       []string{req.Country},
			Province:      []string{req.Province},
			Locality:      []string{req.Locality},
			StreetAddress: []string{req.StreetAddress},
			PostalCode:    []string{req.PostalCode},
			CommonName:    req.CommonName,
		},
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
