package database

import (
	"log"
)

func InsertFixtures(db DB) {
	cs := db.GetCertificateRepository()
	for _, item := range fixtureCertificates {
		if err := cs.Create(item); err != nil {
			log.Printf("certificate write error %+v\n", err)
		}
	}

}

func ClearFixtures(db DB) {
	cs := db.GetCertificateRepository()
	for _, item := range fixtureCertificates {
		if err := cs.DeleteByNameSerialNumber(item.NameSerialNumber); err != nil {
			log.Printf("%s: certificate delete error %+v\n", item.SerialNumber, err)
		}
	}
}

func GetFixtureCertificates() []*Certificate {
	return fixtureCertificates
}

var fixtureCertificates = []*Certificate{
	{
		NameSerialNumber: "testserial",
		CommonName:   "test.id",
	},
}
