package file

import (
	"github.com/mvmaasakkers/certificates/database"
)

// GetCertificateRepository returns a bootstrapped certificate repository
func (db *db) GetCertificateRepository() database.CertificateRepository {
	return &CertificateRepository{db}
}

// CertificateRepository implements the CertificateRepository interface from database package
type CertificateRepository struct {
	db *db
}

// GetByNameSerialNumber gets a certificate by NameSerialNumber
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

// Create creates a certificate
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

// DeleteByNameSerialNumber deletes a certificate by NameSerialNumber
func (repo *CertificateRepository) DeleteByNameSerialNumber(nameSerialNumber string) error {
	return nil
}

// Certificate is the implementation for the Certificate struct in the database package
type Certificate struct {
	database.Certificate
}
