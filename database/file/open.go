package file

import (
	"os"
	"sync"
)

func (db *db) Open() error {
	db.fileLock = &sync.Mutex{}
	db.stateLock = &sync.Mutex{}

	if _, err := os.Stat(db.filename); os.IsNotExist(err) {
		if err := db.writeState(); err != nil {
			return err
		}
	}

	if err := db.readState(); err != nil {
		return err
	}

	return nil
}
