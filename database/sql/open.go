package sql

import (
	"github.com/jinzhu/gorm"
	// Gorm expects implementation packages
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Open opens the connection to the database
func (sqldb *sqlDB) Open() error {
	db, err := gorm.Open(sqldb.dialect, sqldb.connectionString)
	if err != nil {
		return err
	}
	sqldb.conn = db

	return nil
}
