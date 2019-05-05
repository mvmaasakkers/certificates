package file

import (
	"fmt"
	"github.com/mvmaasakkers/certificates/database"
	"os"
	"testing"
)

var (
	testDB database.DB
)

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	testDB = NewDB("file.db")

	if err := testDB.Open(); err != nil {
		fmt.Println(err)
		return 1
	}
	if err := testDB.Provision(); err != nil {
		return 1
	}

	database.InsertFixtures(testDB)

	defer database.ClearFixtures(testDB)

	return m.Run()
}
