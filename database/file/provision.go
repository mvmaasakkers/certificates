package file

func (db *db) Provision() error {
	return db.writeState()
}