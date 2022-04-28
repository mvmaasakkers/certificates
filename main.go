package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mvmaasakkers/certificates/cert"
	"github.com/mvmaasakkers/certificates/database"
	"github.com/mvmaasakkers/certificates/database/file"
	"github.com/mvmaasakkers/certificates/database/sql"
	"github.com/tkuchiki/parsetime"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)

func main() {
	if err := run(os.Args); err != nil {
		os.Exit(1)
	}
}

func run(args []string) error {
	if args == nil {
		args = os.Args[0:1]
	}

	app := cli.NewApp()
	app.Name = "certificates"
	app.Usage = "An opinionated TLS certificate generator."
	app.Version = "v0.7.1"
	app.Description = "An opinionated TLS certificate generator."
	app.Commands = []cli.Command{
		certificateCommand,
	}

	return app.Run(args)
}

var certificateCommand = cli.Command{
	Name:    "certificate",
	Aliases: []string{"cert"},
	Usage:   "certificate commands",
	Subcommands: []cli.Command{
		{
			Name:        "generate-ca",
			Aliases:     []string{"gen-ca"},
			Usage:       "Generate a CA pair",
			Description: `To generate a CA pair you need to supply at least a valid name.`,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "stdout",
					Usage: "Send pem to stdout instead of to file",
				},
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
				cli.StringFlag{
					Name:  "notbefore",
					Value: time.Now().String(),
					Usage: "NotBefore sets the NotBefore timestamp of the certificate request",
				},
				cli.StringFlag{
					Name:  "notafter",
					Value: time.Now().AddDate(0, 0, 30).String(),
					Usage: "NotAfter sets the NotAfter timestamp of the certificate request. The default is 30 days from now.",
				},
				cli.StringFlag{
					Name:  "timezone",
					Value: "UTC",
					Usage: "Timezone to use. Default is set to UTC.",
				},
			},
			Action: func(c *cli.Context) error {

				ca := cert.NewRequest()
				ca.CommonName = c.String("cn")
				ca.Organization = c.String("org")
				ca.Country = c.String("country")
				ca.Province = c.String("province")
				ca.Locality = c.String("locality")
				ca.PostalCode = c.String("postalcode")
				ca.StreetAddress = c.String("streetaddress")

				p, err := parsetime.NewParseTime(c.String("timezone"))
				if err != nil {
					fmt.Printf("Error parsing --timezone: %s\n", err.Error())
					return err
				}

				if c.String("notbefore") != "" {
					notBefore, err := p.Parse(c.String("notbefore"))
					if err != nil {
						fmt.Printf("Error parsing --notbefore: %s\n", err.Error())
						return err
					}
					ca.NotBefore = notBefore
				}

				if c.String("notafter") != "" {
					notAfter, err := p.Parse(c.String("notafter"))
					if err != nil {
						fmt.Printf("Error parsing --notafter: %s\n", err.Error())
						return err
					}
					ca.NotAfter = notAfter
				}

				caCrt, caKey, err := cert.GenerateCA(ca)
				if err != nil {
					fmt.Printf("Error generating CA: %s\n", err.Error())
					return err
				}

				if c.Bool("stdout") {
					fmt.Println(string(caKey))
					fmt.Println(string(caCrt))
					return nil
				}

				if err := ioutil.WriteFile(c.String("ca"), caCrt, 0600); err != nil {
					fmt.Printf("Error writing CA: %s\n", err.Error())
					return err
				}

				if err := ioutil.WriteFile(c.String("ca-key"), caKey, 0600); err != nil {
					fmt.Printf("Error writing CA: %s\n", err.Error())
					return err
				}

				return nil
			},
		},
		{
			Name:    "generate",
			Aliases: []string{"gen"},
			Usage:   "Generate a signed certificate pair",
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
					Name:  "ca-DB-type",
					Value: "file",
					Usage: "CA DB type (file)",
				},

				cli.StringFlag{
					Name:  "ca-DB-file",
					Value: "file.DB",
					Usage: "File DB filename",
				},
				cli.StringFlag{
					Name:  "ca-DB-sql-dialect",
					Value: "mysql",
					Usage: "SQL Dialect",
				},
				cli.StringFlag{
					Name:  "ca-DB-sql-cs",
					Value: "user:pass@tcp(localhost:3306)/test?charset=utf8&parseTime=True",
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
					Name:  "csr",
					Value: "",
					Usage: "Give the filepath to an existing CSR if you want to sign using a pre-existing CSR",
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
				cli.StringFlag{
					Name:  "notbefore",
					Value: time.Now().String(),
					Usage: "NotBefore sets the NotBefore timestamp of the certificate request",
				},
				cli.StringFlag{
					Name:  "notafter",
					Value: time.Now().AddDate(0, 0, 30).String(),
					Usage: "NotAfter sets the NotAfter timestamp of the certificate request. The default is 30 days from now.",
				},
				cli.StringFlag{
					Name:  "timezone",
					Value: "UTC",
					Usage: "Timezone",
				},
				cli.IntFlag{
					Name:  "bitsize",
					Value: 4096,
					Usage: "Encryption key bitsize",
				},
			},
			Action: func(c *cli.Context) error {

				var cr *cert.Request
				if c.String("csr") != "" {
					csrFile, err := ioutil.ReadFile(c.String("csr"))
					if err != nil {
						fmt.Printf("Error reading CSR: %s\n", err.Error())
						return err
					}

					cr, err = cert.ReadCSR(csrFile)
					if err != nil {
						fmt.Printf("Error reading CSR: %s\n", err.Error())
						return err
					}
				} else {
					cr = cert.NewRequest()
					cr.CommonName = c.String("cn")
					cr.Organization = c.String("org")
					cr.Country = c.String("country")
					cr.Province = c.String("province")
					cr.Locality = c.String("locality")
					cr.PostalCode = c.String("postalcode")
					cr.StreetAddress = c.String("streetaddress")
					cr.NameSerialNumber = c.String("name-serialnumber")
					cr.BitSize = c.Int("bitsize")

					cr.SubjectAltNames = c.StringSlice("subject-alt-name")

					if cr.NameSerialNumber == "" {
						// Generating serial number
						sn, err := uuid.NewRandom()
						if err != nil {
							fmt.Printf("Error generating serial number: %s\n", err.Error())
							return err
						}
						cr.NameSerialNumber = sn.String()
					}
				}

				if c.Int64("serialnumber") != 0 {
					cr.SerialNumber = big.NewInt(c.Int64("serialnumber"))
				} else {
					cr.SerialNumber, _ = cert.GenerateRandomBigInt()
				}

				p, err := parsetime.NewParseTime(c.String("timezone"))
				if err != nil {
					fmt.Printf("Error parsing --timezone: %s\n", err.Error())
					return err
				}

				if c.String("notbefore") != "" {
					notBefore, err := p.Parse(c.String("notbefore"))
					if err != nil {
						fmt.Printf("Error parsing --notbefore: %s\n", err.Error())
						return err
					}
					cr.NotBefore = notBefore
				}

				if c.String("notafter") != "" {
					notAfter, err := p.Parse(c.String("notafter"))
					if err != nil {
						fmt.Printf("Error parsing --notafter: %s\n", err.Error())
						return err
					}
					cr.NotAfter = notAfter
				}

				caCrt, err := ioutil.ReadFile(c.String("ca"))
				if err != nil {
					fmt.Printf("Error reading CA certificate: %s\n", err.Error())
					return err
				}
				caKey, err := ioutil.ReadFile(c.String("ca-key"))
				if err != nil {
					fmt.Printf("Error reading CA key: %s\n", err.Error())
					return err
				}

				// DB
				var DB database.DB
				switch c.String("ca-DB-type") {
				case "sql":
					DB = sql.NewDB(c.String("ca-DB-sql-dialect"), c.String("ca-DB-sql-cs"))
				case "file":
					DB = file.NewDB(c.String("ca-DB-file"))
				}

				if DB == nil {
					return database.ErrorNilConnection
				}

				if err := DB.Open(); err != nil {
					fmt.Printf("Error opening DB: %s\n", err.Error())
					return err
				}
				defer DB.Close()

				if err := DB.Provision(); err != nil {
					fmt.Printf("Error provisioning DB: %s\n", err.Error())
					return err
				}

				crt, key, err := cert.GenerateCertificate(cr, caCrt, caKey)
				if err != nil {
					fmt.Printf("Error generating certificate: %s\n", err.Error())
					return err
				}

				// Store in CA DB
				DBCert := database.NewCertificate()
				DBCert.Status = "valid"
				DBCert.ExpirationDate = cr.NotAfter
				DBCert.RevocationDate = nil
				DBCert.SerialNumber = cr.SerialNumber
				DBCert.NameSerialNumber = cr.NameSerialNumber
				DBCert.CommonName = cr.CommonName

				if err := DB.GetCertificateRepository().Create(DBCert); err != nil {
					fmt.Printf("Error saving certificate to DB: %s\n", err.Error())
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
					fmt.Printf("Error writing certificate to file: %s\n", err.Error())
					return err
				}
				fmt.Printf("Writing key to %s\n", c.String("key"))
				if err := ioutil.WriteFile(c.String("key"), key, 0600); err != nil {
					fmt.Printf("Error writing certificate key to file: %s\n", err.Error())
					return err
				}
				return nil
			},
		},
	},
}
