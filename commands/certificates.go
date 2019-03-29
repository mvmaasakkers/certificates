package commands

import (
	"fmt"
	"github.com/mvmaasakkers/certificates/cert"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
)

var CertificateCommand = cli.Command{
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
					Name:  "serialnumber",
					Value: "",
					Usage: "SerialNumber",
				},
				cli.StringSliceFlag{
					Name: "subject-alt-name",
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
				cr.SerialNumber = c.String("serialnumber")
				cr.SubjectAltNames = c.StringSlice("subject-alt-name")

				caCrt, err := ioutil.ReadFile(c.String("ca"))
				if err != nil {
					return err
				}
				caKey, err := ioutil.ReadFile(c.String("ca-key"))
				if err != nil {
					return err
				}

				crt, key, err := cr.GenerateCertificate(caCrt, caKey)
				if err != nil {
					return err
				}

				if c.Bool("stdout") {
					fmt.Println(string(key))
					fmt.Println(string(crt))
					return nil
				}

				if err := ioutil.WriteFile(c.String("crt"), crt, 0600); err != nil {
					return err
				}

				if err := ioutil.WriteFile(c.String("key"), key, 0600); err != nil {
					return err
				}

				return nil
			},
		},
	},
}
