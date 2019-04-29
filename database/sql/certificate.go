package sql

import (
	"github.com/mvmaasakkers/certificates/cert"
)

func (sqldb *sqlDB) GetCertificateRepository() cert.CertificateRepository {
	return &CertificateRepository{sqldb}
}

type CertificateRepository struct {
	sqldb *sqlDB
}

func (repo *CertificateRepository) List() ([]*cert.Certificate, error) {

	return nil, nil
}

func (repo *CertificateRepository) GetByUUID(uuid string) (*cert.Certificate, error) {
	return nil, nil
}

func (repo *CertificateRepository) GetBySerialNumber(serialNumber string) (*cert.Certificate, error) {
	crt := &cert.Certificate{}
	if err := repo.sqldb.conn.Where("serial_number = ?", serialNumber).First(crt).Error; err != nil {
		return nil, GetError(err)
	}
	return crt, nil
}

func (repo *CertificateRepository) Create(certificate *cert.Certificate) error {

	if err := repo.sqldb.conn.Create(certificate).Error; err != nil {
		return GetError(err)
	}
	return nil
}

func (repo *CertificateRepository) Update(certificate *cert.Certificate) error {
	return nil
}


func (repo *CertificateRepository) Delete(certificate *cert.Certificate) error {
	return nil
}

func (repo *CertificateRepository) DeleteBySerialNumber(serialNumber string) error {
	return nil
}

type Certificate struct {
	GormModel
	cert.Certificate
}
