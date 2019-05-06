package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/mvmaasakkers/certificates/database"
)

type sqlDB struct {
	dialect          string
	connectionString string
	conn             *gorm.DB
}

// NewDB bootstraps a new File DB instance
func NewDB(dialect, connectionString string) database.DB {
	sqldb := &sqlDB{
		dialect:          dialect,
		connectionString: connectionString,
	}
	return sqldb
}
