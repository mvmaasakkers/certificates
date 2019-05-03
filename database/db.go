package database

import (
	"github.com/google/uuid"
	"math/big"
	"time"
)

type DB interface {
	Open() error
	Close() error
	Provision() error

	GetCertificateRepository() CertificateRepository
}

type CertificateRepository interface {
	List() ([]*Certificate, error)
	GetByUUID(uuid string) (*Certificate, error)
	GetByNameSerialNumber(nameSerialNumber string) (*Certificate, error)
	Create(certificate *Certificate) error
	Update(certificate *Certificate) error
	DeleteByNameSerialNumber(nameSerialNumber string) error
}

type Certificate struct {
	Meta

	Status           string
	ExpirationDate   time.Time
	RevocationDate   *time.Time
	NameSerialNumber string   `gorm:"unique"`
	SerialNumber     *big.Int `gorm:"unique"`
	CommonName       string
}

func NewCertificate() *Certificate {
	uid, _ := uuid.NewRandom()

	c := &Certificate{}
	c.UUID = uid.String()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	return c
}

type Meta struct {
	UUID      string `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
