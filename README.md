# Certificates helper

This is an opinionated helper for generating tls certificates.
It outputs only in PEM format but this enables you easily generate certificate
chains for MA TLS.

## Usage

### Generate a CA set

`certificate generate-ca`

```
NAME:
   Certificates certificate generate-ca - generate ca

USAGE:
   Certificates certificate generate-ca [command options] [arguments...]

OPTIONS:
   --ca value             Filename to write the ca cert to (default: "ca.crt")
   --ca-key value         Filename to write the ca key to (default: "ca.key")
   --cn value             Common name attached to the ca cert
   --org value            Organisation
   --country value        Country
   --province value       Province
   --locality value       Locality
   --postalcode value     PostalCode
   --streetaddress value  StreetAddress
```

### Generate a certificate

This needs a pregenerated CA certificate and key (see "Generate a CA set")

`certificate generate`

```
NAME:
   Certificates certificate generate - generate certificate

USAGE:
   Certificates certificate generate [command options] [arguments...]

OPTIONS:
   --stdout                   Send pem to stdout instead of to file
   --ca value                 CA Certificate file (default: "ca.crt")
   --ca-key value             CA Key file (default: "ca.key")
   --ca-db-type value         CA DB type (sql) (default: "sql")
   --ca-db-sql-dialect value  SQL Dialect (default: "sqlite3")
   --ca-db-sql-cs value       SQL Connection String (default: "sql.db")
   --crt value                Filename to write certificate to (default: "certificate.crt")
   --key value                Filename to write key to (default: "certificate.key")
   --cn value                 Common name attached to the cert
   --org value                Organisation
   --country value            Country
   --province value           Province
   --locality value           Locality
   --postalcode value         PostalCode
   --streetaddress value      StreetAddress
   --serialnumber value       SerialNumber
   --subject-alt-name value   Subject Alt Name
```

The CA database can be one of the following flavours of sql: sqlite3, mysql, postgresql or mssql. 

## Development setup

No external dependencies are needed. Just:

- `dep ensure`
- `go run main.go`

And this will output:

```bash
NAME:
   Certificates - An opinionated TLS certificate generator.

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   v0.0.1-alpha2

DESCRIPTION:
   An opinionated TLS certificate generator.

COMMANDS:
     certificate, cert  certificate commands
     help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```