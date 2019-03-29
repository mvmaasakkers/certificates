package sql

import "github.com/mvmaasakkers/certificates/database"

func (sqldb *sqlDB) Close() error {
	if sqldb.conn == nil {
		return database.ErrorNilConnection
	}
	return sqldb.conn.Close()
}
