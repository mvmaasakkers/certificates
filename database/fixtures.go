package database

import (
	"github.com/mvmaasakkers/certificates/cert"
	"log"
)

func InsertFixtures(db cert.DB) {
	cs := db.GetCertificateRepository()
	for _, item := range fixtureCertificates {
		if err := cs.Create(item); err != nil {
			log.Printf("certificate write error %+v\n", err)
		}
	}

}

func ClearFixtures(db cert.DB) {
	cs := db.GetCertificateRepository()
	for _, item := range fixtureCertificates {
		if err := cs.DeleteBySerialNumber(item.SerialNumber); err != nil {
			log.Printf("%s: certificate delete error %+v\n", item.SerialNumber, err)
		}
	}
}

func GetFixtureCertificates() []*cert.Certificate {
	return fixtureCertificates
}

var fixtureCertificates = []*cert.Certificate{
	{
		SerialNumber: "testserial",
		CommonName:   "test.id",
	},
}
