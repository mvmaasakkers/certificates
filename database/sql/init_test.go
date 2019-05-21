package sql

import (
	"github.com/mvmaasakkers/certificates/database"
	"github.com/mvmaasakkers/certificates/database/test"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var (
	testDB database.DB
)

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	filename := filepath.Join(os.TempDir(), "certificates_sql_test.db")
	testDB = NewDB("sqlite3", filename)
	if err := testDB.Open(); err != nil {
		log.Println(err)
		return 1
	}
	defer testDB.Close()
	defer os.Remove(filename) // clean up

	if err := testDB.Provision(); err != nil {
		return 1
	}

	test.InsertFixtures(testDB)

	defer database.ClearFixtures(testDB)

	return m.Run()
}
