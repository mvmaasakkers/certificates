package sql

import (
	"github.com/jinzhu/gorm"
	"github.com/mvmaasakkers/certificates/database"
)

type sqlDB struct {
	dialect string
	connectionString string
	conn *gorm.DB
}

func NewDB(dialect, connectionString string) database.DB {
	sqldb := &sqlDB{
		dialect: dialect,
		connectionString: connectionString,
	}
	return sqldb
}
