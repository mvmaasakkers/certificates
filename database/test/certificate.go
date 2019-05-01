package test

import (
	"github.com/mvmaasakkers/certificates/cert"
	"github.com/mvmaasakkers/certificates/database"
	"testing"
)

var certificateTests = []struct {
	Id          string
	Error       error
	Certificate *cert.Certificate
}{
	{
		Id:          "test.id",
		Error:       nil,
		Certificate: &cert.Certificate{CommonName: "test.id", SerialNumber: "testserial"},
	},
	{
		Id:          "testnotfound",
		Error:       database.ErrorObjectNotFound,
		Certificate: &cert.Certificate{},
	},
}

func TestCertificate_Certificate(t *testing.T, certificateRepository cert.CertificateRepository) {
	for _, test := range certificateTests {
		_, err := certificateRepository.GetBySerialNumber(test.Certificate.SerialNumber)
		if err != test.Error {
			t.Errorf("%s: expected error %+v, got error %+v", test.Id, test.Error, err)
			t.Fail()
		}
	}
}

var createCertificateTests = []struct {
	Error       error
	DeleteError error
	Certificate *cert.Certificate
}{
	{
		Certificate: &cert.Certificate{CommonName: "test.id", SerialNumber: "testserial"},
		Error:       database.ErrorDuplicateObject,
	},
	{
		Certificate: &cert.Certificate{CommonName: "testid_2", SerialNumber: "two"},
		Error:       nil,
	},
}

func TestCertificate_CreateCertificate(t *testing.T, certificateRepository cert.CertificateRepository) {
	for _, test := range createCertificateTests {
		err := certificateRepository.Create(test.Certificate)
		if err != test.Error {
			t.Errorf("%s: expected error %+v, got error %+v", test.Certificate.SerialNumber, test.Error, err)
			t.Fail()
		}

		if err := certificateRepository.DeleteBySerialNumber(test.Certificate.SerialNumber); err != test.DeleteError {
			t.Errorf("%s: expected delete error %+v, got error %+v", test.Certificate.SerialNumber, test.Error, err)
			t.Fail()
		}
	}
}
