package database

type DB interface {
	Open() error
	Close() error
	Provision() error

	GetCertificateRepository() CertificateRepository
}
