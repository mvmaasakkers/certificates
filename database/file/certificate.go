package file

import (
	"github.com/mvmaasakkers/certificates/database"
)

func (db *db) GetCertificateRepository() database.CertificateRepository {
	return &CertificateRepository{db}
}

type CertificateRepository struct {
	db *db
}

func (repo *CertificateRepository) List() ([]*database.Certificate, error) {

	return nil, nil
}

func (repo *CertificateRepository) GetByUUID(uuid string) (*database.Certificate, error) {
	return nil, nil
}

func (repo *CertificateRepository) GetByNameSerialNumber(nameSerialNumber string) (*database.Certificate, error) {
	repo.db.stateLock.Lock()
	defer repo.db.stateLock.Unlock()

	for _, c := range repo.db.state.Certificates {
		if c.NameSerialNumber == nameSerialNumber {
			return c, nil
		}
	}

	return nil, database.ErrorObjectNotFound
}

func (repo *CertificateRepository) Create(certificate *database.Certificate) error {
	repo.db.stateLock.Lock()
	defer repo.db.stateLock.Unlock()

	if repo.db.state.Certificates == nil {
		repo.db.state.Certificates = make(map[string]*database.Certificate, 0)
	}
	
	if _, ok := repo.db.state.Certificates[certificate.CommonName]; ok {
		return database.ErrorDuplicateObject
	}

	repo.db.state.Certificates[certificate.CommonName] = certificate
	return repo.db.writeState()
}

func (repo *CertificateRepository) Update(certificate *database.Certificate) error {
	return nil
}

func (repo *CertificateRepository) Delete(certificate *database.Certificate) error {
	return nil
}

func (repo *CertificateRepository) DeleteByNameSerialNumber(nameSerialNumber string) error {
	return nil
}

type Certificate struct {
	database.Certificate
}
