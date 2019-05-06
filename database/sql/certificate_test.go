package sql

import (
	"github.com/mvmaasakkers/certificates/database/test"
	"testing"
)

func TestCertificateCertificate(t *testing.T) {
	test.TestCertificateCertificate(t, testDB.GetCertificateRepository())
}

func TestCertificateCreateCertificate(t *testing.T) {
	test.TestCertificateCreateCertificate(t, testDB.GetCertificateRepository())
}
