package file

import (
	"encoding/json"
	"github.com/mvmaasakkers/certificates/database"
	"io/ioutil"
	"sync"
	"time"
)

type db struct {
	filename  string
	stateLock *sync.Mutex
	fileLock  *sync.Mutex
	state     *state
}

type state struct {
	LastSync     time.Time
	Certificates map[string]*database.Certificate
}

// NewDB bootstraps a new File DB instance
func NewDB(filename string) database.DB {
	return &db{
		filename: filename,
		state:    &state{},
	}
}

func (db *db) readState() error {
	db.fileLock.Lock()
	defer db.fileLock.Unlock()

	d, err := ioutil.ReadFile(db.filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(d, db.state)
}

func (db *db) writeState() error {
	db.fileLock.Lock()
	defer db.fileLock.Unlock()

	db.state.LastSync = time.Now()

	d, err := json.Marshal(db.state)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(db.filename, d, 0644)
}
