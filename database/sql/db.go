package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/mvmaasakkers/certificates/cert"
)

type sqlDB struct {
	dialect string
	connectionString string
	conn *gorm.DB
}

func NewDB(dialect, connectionString string) cert.DB {
	sqldb := &sqlDB{
		dialect: dialect,
		connectionString: connectionString,
	}
	return sqldb
}
