package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mvmaasakkers/certificates/cert"
	"github.com/mvmaasakkers/certificates/database"
	"github.com/mvmaasakkers/certificates/database/sql"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"log"
	"math/big"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Certificates"
	app.Usage = "An opinionated TLS certificate generator."
	app.Version = "v0.0.1-alpha3"
	app.Description = "An opinionated TLS certificate generator."
	app.Commands = []cli.Command{
		certificateCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

var certificateCommand = cli.Command{
	Name:    "certificate",
	Aliases: []string{"cert"},
	Usage:   "certificate commands",
	Subcommands: []cli.Command{
		{
			Name:    "generate-ca",
			Aliases: []string{"gen-ca"},
			Usage:   "generate ca",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ca",
					Value: "ca.crt",
					Usage: "Filename to write the ca cert to",
				},
				cli.StringFlag{
					Name:  "ca-key",
					Value: "ca.key",
					Usage: "Filename to write the ca key to",
				},
				cli.StringFlag{
					Name:  "cn",
					Value: "",
					Usage: "Common name attached to the ca cert",
				},
				cli.StringFlag{
					Name:  "org",
					Value: "",
					Usage: "Organisation",
				},
				cli.StringFlag{
					Name:  "country",
					Value: "",
					Usage: "Country",
				},
				cli.StringFlag{
					Name:  "province",
					Value: "",
					Usage: "Province",
				},
				cli.StringFlag{
					Name:  "locality",
					Value: "",
					Usage: "Locality",
				},
				cli.StringFlag{
					Name:  "postalcode",
					Value: "",
					Usage: "PostalCode",
				},
				cli.StringFlag{
					Name:  "streetaddress",
					Value: "",
					Usage: "StreetAddress",
				},
			},
			Action: func(c *cli.Context) error {

				ca := cert.NewCARequest()
				ca.CommonName = c.String("cn")
				ca.Organization = c.String("org")
				ca.Country = c.String("country")
				ca.Province = c.String("province")
				ca.Locality = c.String("locality")
				ca.PostalCode = c.String("postalcode")
				ca.StreetAddress = c.String("streetaddress")

				caCrt, caKey, err := ca.GenerateCA()
				if err != nil {
					return err
				}

				if err := ioutil.WriteFile(c.String("ca"), caCrt, 0600); err != nil {
					return err
				}

				if err := ioutil.WriteFile(c.String("ca-key"), caKey, 0600); err != nil {
					return err
				}

				return nil
			},
		},
		{
			Name:    "generate",
			Aliases: []string{"gen"},
			Usage:   "generate certificate",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "stdout",
					Usage: "Send pem to stdout instead of to file",
				},
				cli.StringFlag{
					Name:  "ca",
					Value: "ca.crt",
					Usage: "CA Certificate file",
				},
				cli.StringFlag{
					Name:  "ca-key",
					Value: "ca.key",
					Usage: "CA Key file",
				},
				cli.StringFlag{
					Name:  "ca-db-type",
					Value: "sql",
					Usage: "CA DB type (sql)",
				},
				cli.StringFlag{
					Name:  "ca-db-sql-dialect",
					Value: "sqlite3",
					Usage: "SQL Dialect",
				},
				cli.StringFlag{
					Name:  "ca-db-sql-cs",
					Value: "sql.db",
					Usage: "SQL Connection String",
				},
				cli.StringFlag{
					Name:  "crt",
					Value: "certificate.crt",
					Usage: "Filename to write certificate to",
				},
				cli.StringFlag{
					Name:  "key",
					Value: "certificate.key",
					Usage: "Filename to write key to",
				},
				cli.StringFlag{
					Name:  "cn",
					Value: "",
					Usage: "Common name attached to the cert",
				},
				cli.StringFlag{
					Name:  "org",
					Value: "",
					Usage: "Organisation",
				},
				cli.StringFlag{
					Name:  "country",
					Value: "",
					Usage: "Country",
				},
				cli.StringFlag{
					Name:  "province",
					Value: "",
					Usage: "Province",
				},
				cli.StringFlag{
					Name:  "locality",
					Value: "",
					Usage: "Locality",
				},
				cli.StringFlag{
					Name:  "postalcode",
					Value: "",
					Usage: "PostalCode",
				},
				cli.StringFlag{
					Name:  "streetaddress",
					Value: "",
					Usage: "StreetAddress",
				},
				cli.StringFlag{
					Name:  "name-serialnumber",
					Value: "",
					Usage: "Name SerialNumber",
				},
				cli.Int64Flag{
					Name:  "serialnumber",
					Value: 0,
					Usage: "SerialNumber",
				},
				cli.StringSliceFlag{
					Name:  "subject-alt-name",
					Usage: "Subject Alt Name",
				},
			},
			Action: func(c *cli.Context) error {
				cr := cert.NewCertRequest()
				cr.CommonName = c.String("cn")
				cr.Organization = c.String("org")
				cr.Country = c.String("country")
				cr.Province = c.String("province")
				cr.Locality = c.String("locality")
				cr.PostalCode = c.String("postalcode")
				cr.StreetAddress = c.String("streetaddress")
				cr.NameSerialNumber = c.String("name-serialnumber")

				cr.SubjectAltNames = c.StringSlice("subject-alt-name")

				if cr.NameSerialNumber == "" {
					// Generating serial number
					sn, err := uuid.NewRandom()
					if err != nil {
						return err
					}
					cr.NameSerialNumber = sn.String()
				}

				if c.Int64("serialnumber") != 0 {
					cr.SerialNumber = big.NewInt(c.Int64("serialnumber"))
				} else {
					cr.SerialNumber, _ = cert.GenerateRandomBigInt()
				}

				caCrt, err := ioutil.ReadFile(c.String("ca"))
				if err != nil {
					return err
				}
				caKey, err := ioutil.ReadFile(c.String("ca-key"))
				if err != nil {
					return err
				}

				// DB
				var db database.DB
				switch c.String("ca-db-type") {
				case "sql":
					db = sql.NewDB(c.String("ca-db-sql-dialect"), c.String("ca-db-sql-cs"))
				}

				if db == nil {
					return database.ErrorNilConnection
				}

				if err := db.Open(); err != nil {
					return err
				}
				defer db.Close()

				if err := db.Provision(); err != nil {
					return err
				}

				crt, key, err := cr.GenerateCertificate(caCrt, caKey)
				if err != nil {
					return err
				}

				// Store in CA DB
				dbCert := database.NewCertificate()
				dbCert.Status = "valid"
				dbCert.ExpirationDate = cr.NotAfter
				dbCert.RevocationDate = nil
				dbCert.SerialNumber = cr.SerialNumber
				dbCert.NameSerialNumber = cr.NameSerialNumber
				dbCert.CommonName = cr.CommonName

				if err := db.GetCertificateRepository().Create(dbCert); err != nil {
					return err
				}

				if c.Bool("stdout") {
					fmt.Println(string(key))
					fmt.Println(string(crt))
					return nil
				}

				fmt.Printf("Generated certificate with serial number %s\n", cr.SerialNumber)
				fmt.Printf("Writing certificate to %s\n", c.String("crt"))
				if err := ioutil.WriteFile(c.String("crt"), crt, 0600); err != nil {
					return err
				}
				fmt.Printf("Writing key to %s\n", c.String("key"))
				if err := ioutil.WriteFile(c.String("key"), key, 0600); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
