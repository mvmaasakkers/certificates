package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/mvmaasakkers/certificates/database"
	"math/big"
	"time"
)

type CertRequest struct {
	Organization    string
	Country         string
	Province        string
	Locality        string
	StreetAddress   string
	PostalCode      string
	CommonName      string
	SerialNumber    string
	SubjectAltNames []string

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

	for _, n := range req.SubjectAltNames {
		if n == "" {
			return ErrorInvalidSubjectAltName
		}
	}

	return nil
}

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

	if req.SerialNumber != "" {
		name.SerialNumber = req.SerialNumber
	}

	return name
}

func (req *CertRequest) GenerateCertificate(db database.DB, caCrt []byte, caKey []byte) ([]byte, []byte, error) {
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

	// Store in CA DB
	dbCert := database.NewCertificate()
	dbCert.Status = "valid"
	dbCert.ExpirationDate = req.NotAfter
	dbCert.RevocationDate = nil
	dbCert.SerialNumber = req.SerialNumber
	dbCert.CommonName = req.CommonName

	certificateRepository := db.GetCertificateRepository()
	if err := certificateRepository.Create(dbCert); err != nil {
		return nil, nil, err
	}

	return pemOut.Bytes(), pemKeyOut.Bytes(), nil
}
