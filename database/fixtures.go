package database

import (
	"log"
)

// InsertFixtures inserts data needed for tests
func InsertFixtures(db DB) {
	cs := db.GetCertificateRepository()
	for _, item := range fixtureCertificates {
		if err := cs.Create(item); err != nil {
			log.Printf("certificate write error %+v\n", err)
		}
	}
}

// ClearFixtures cleans up data needed for tests
func ClearFixtures(db DB) {
	cs := db.GetCertificateRepository()
	for _, item := range fixtureCertificates {
		if err := cs.DeleteByNameSerialNumber(item.NameSerialNumber); err != nil {
			log.Printf("%s: certificate delete error %+v\n", item.SerialNumber, err)
		}
	}
}

var fixtureCertificates = []*Certificate{
	{
		NameSerialNumber: "testserial",
		CommonName:       "test.id",
	},
}
