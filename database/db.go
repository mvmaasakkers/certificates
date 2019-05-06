package database

import (
	"github.com/google/uuid"
	"math/big"
	"time"
)

// DB is the interface for database implementations
type DB interface {
	Open() error
	Close() error
	Provision() error

	GetCertificateRepository() CertificateRepository
}

// CertificateRepository is the interface for certificate repository implementations
type CertificateRepository interface {
	GetByNameSerialNumber(nameSerialNumber string) (*Certificate, error)
	Create(certificate *Certificate) error
	DeleteByNameSerialNumber(nameSerialNumber string) error
}

// Certificate is the struct for the database certificate object
type Certificate struct {
	Meta

	Status           string
	ExpirationDate   time.Time
	RevocationDate   *time.Time
	NameSerialNumber string   `gorm:"unique"`
	SerialNumber     *big.Int `gorm:"unique"`
	CommonName       string
}

// NewCertificate creates a new Certificate object with a generated ID and settings the CreatedAt and UpdatedAt to now
func NewCertificate() *Certificate {
	uid, _ := uuid.NewRandom()

	c := &Certificate{}
	c.UUID = uid.String()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	return c
}

// Meta is an extra set of recurring data used in database objects
type Meta struct {
	UUID      string `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
