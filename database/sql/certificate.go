package sql

import (
	"github.com/mvmaasakkers/certificates/database"
)

// GetCertificateRepository returns a bootstrapped certificate repository
func (sqldb *sqlDB) GetCertificateRepository() database.CertificateRepository {
	return &CertificateRepository{sqldb}
}

// CertificateRepository implements the CertificateRepository interface from database package
type CertificateRepository struct {
	sqldb *sqlDB
}

// GetByNameSerialNumber gets a certificate by NameSerialNumber
func (repo *CertificateRepository) GetByNameSerialNumber(nameSerialNumber string) (*database.Certificate, error) {
	crt := &database.Certificate{}
	if err := repo.sqldb.conn.Where("name_serial_number = ?", nameSerialNumber).First(crt).Error; err != nil {
		return nil, GetError(err)
	}
	return crt, nil
}

// Create creates a certificate
func (repo *CertificateRepository) Create(certificate *database.Certificate) error {

	if err := repo.sqldb.conn.Create(certificate).Error; err != nil {
		return GetError(err)
	}
	return nil
}

// DeleteByNameSerialNumber deletes a certificate by NameSerialNumber
func (repo *CertificateRepository) DeleteByNameSerialNumber(nameSerialNumber string) error {
	return nil
}

// Certificate is the implementation for the Certificate struct in the database package
type Certificate struct {
	GormModel
	database.Certificate
}
