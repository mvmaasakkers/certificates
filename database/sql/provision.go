package sql

import (
	// Gorm expects implementation packages
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func (sqldb *sqlDB) Provision() error {

	return sqldb.conn.AutoMigrate(&Certificate{}).Error
}
