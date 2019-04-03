package sql

import (
	"github.com/mvmaasakkers/certificates/database/test"
	"testing"
)

func TestCertificate_Certificate(t *testing.T) {
	test.TestCertificate_Certificate(t, testDB)
}

func TestCertificate_CreateCertificate(t *testing.T) {
	test.TestCertificate_CreateCertificate(t, testDB)
}