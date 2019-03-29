package sql

import (
	"github.com/mvmaasakkers/certificates/database"
)

func (sqldb *sqlDB) GetCertificateRepository() database.CertificateRepository {
	return &CertificateRepository{sqldb}
}

type CertificateRepository struct {
	sqldb *sqlDB
}

func (repo *CertificateRepository) List() ([]*database.Certificate, error) {

	return nil, nil
}

func (repo *CertificateRepository) GetByUUID(uuid string) (*database.Certificate, error) {
	return nil, nil
}

func (repo *CertificateRepository) GetBySerialNumber(serialNumber string) (*database.Certificate, error) {
	return nil, nil
}

func (repo *CertificateRepository) Create(certificate *database.Certificate) error {
	if err := repo.sqldb.conn.Create(certificate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *CertificateRepository) Update(certificate *database.Certificate) error {
	return nil
}

func (repo *CertificateRepository) Delete(certificate *database.Certificate) error {
	return nil
}

type Certificate struct {
	GormModel
	database.Certificate
}
