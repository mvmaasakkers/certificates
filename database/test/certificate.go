package test

import (
	"github.com/mvmaasakkers/certificates/database"
	"testing"
)

var certificateTests = []struct {
	ID          string
	Error       error
	Certificate *database.Certificate
}{
	{
		ID:          "test.id",
		Error:       nil,
		Certificate: &database.Certificate{CommonName: "test.id", NameSerialNumber: "testserial"},
	},
	{
		ID:          "testnotfound",
		Error:       database.ErrorObjectNotFound,
		Certificate: &database.Certificate{},
	},
}

// TestCertificateCertificate tests
func TestCertificateCertificate(t *testing.T, certificateRepository database.CertificateRepository) {
	for _, test := range certificateTests {
		_, err := certificateRepository.GetByNameSerialNumber(test.Certificate.NameSerialNumber)
		if err != test.Error {
			t.Errorf("%s: expected error %+v, got error %+v", test.ID, test.Error, err)
			t.Fail()
		}
	}
}

var createCertificateTests = []struct {
	Error       error
	DeleteError error
	Certificate *database.Certificate
}{
	{
		Certificate: &database.Certificate{CommonName: "test.id", NameSerialNumber: "testserial"},
		Error:       database.ErrorDuplicateObject,
	},
	{
		Certificate: &database.Certificate{CommonName: "testid_2", NameSerialNumber: "two"},
		Error:       nil,
	},
}

// TestCertificateCreateCertificate tests
func TestCertificateCreateCertificate(t *testing.T, certificateRepository database.CertificateRepository) {
	for _, test := range createCertificateTests {
		err := certificateRepository.Create(test.Certificate)
		if err != test.Error {
			t.Errorf("%s: expected error %+v, got error %+v", test.Certificate.SerialNumber, test.Error, err)
			t.Fail()
		}

		if err := certificateRepository.DeleteByNameSerialNumber(test.Certificate.NameSerialNumber); err != test.DeleteError {
			t.Errorf("%s: expected delete error %+v, got error %+v", test.Certificate.SerialNumber, test.Error, err)
			t.Fail()
		}
	}
}
