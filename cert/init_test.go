package cert

import (
	"fmt"
	"os"
	"testing"
)

type testData struct {
	Request *CertRequest
	Crt []byte
	Key []byte
}

var testCA *testData


func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {

	testCA = &testData{
		Request: &CertRequest{
			CommonName: "test.local",
		},
	}
	var err error
	testCA.Crt, testCA.Key, err = GenerateCA(testCA.Request)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	return m.Run()
}
