package sql

import (
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func (sqldb *sqlDB) Provision() error {

	if err := sqldb.conn.AutoMigrate(&Certificate{}).Error; err != nil {
		return err
	}

	return nil
}
