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

type CertRequest struct {
	Organization  string
	Country       string
	Province      string
	Locality      string
	StreetAddress string
	PostalCode    string
	CommonName    string
	SerialNumber  string

	NotBefore time.Time
	NotAfter  time.Time
}

func NewCertRequest() *CertRequest {
	return &CertRequest{
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 0, 1),
	}
}

func (req *CertRequest) Validate() error {
	if req.CommonName == "" {
		return ErrorInvalidCommonName
	}

	return nil
}

func (req *CertRequest) GenerateCertificate(caCrt []byte, caKey []byte) ([]byte, []byte, error) {
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

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{req.Organization},
			Country:       []string{req.Country},
			Province:      []string{req.Province},
			Locality:      []string{req.Locality},
			StreetAddress: []string{req.StreetAddress},
			PostalCode:    []string{req.PostalCode},
			CommonName:    req.CommonName,
			SerialNumber:  req.SerialNumber,
		},
		NotBefore:   req.NotBefore,
		NotAfter:    req.NotAfter,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
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
