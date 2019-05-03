# Certificates helper 

[![Build Status](https://travis-ci.com/mvmaasakkers/certificates.svg?branch=master)](https://travis-ci.com/mvmaasakkers/certificates) 
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/mvmaasakkers/certificates?status.svg)](https://godoc.org/github.com/mvmaasakkers/certificates)

This is an opinionated helper for generating tls certificates.
It outputs only in PEM format but this enables you easily generate certificate
chains for MA TLS.

## cert package

The cert package can be used directly in your application without the need of the command 
line interface, underlying database layer or external dependencies. This way certificate 
generation can be easily embedded. 

Documentation can be found [here](https://godoc.org/github.com/mvmaasakkers/certificates/cert).

## Usage

### Generate a CA set

You can generate a CA set by using the generate-ca subcommand like the following example:

`certificates cert gen-ca --cn=*.test.domain --stdout`

This will output the key and certificate directly to stdout like this (parts are omitted for readability):

```
-----BEGIN RSA PRIVATE KEY-----
MIIJJwIBAAKCAgEA0txN/brNlBcGrU8mAxL8V19pS1dWEVVTF82LDahI7FMsPPkM
sg5iBCLwYJhnVRPucUmcGC1NyljCy/yW0Cbwl5aNWozAfEkiUpWsukn/ZcMuXvac
qsPRK0Xswbr305NDRnlphoeutyzXAhW2P4FQGCwSfx/Mlaezphc7AreLKg==
-----END RSA PRIVATE KEY-----

-----BEGIN CERTIFICATE-----
MIIE3zCCAsegAwIBAgIFANHEYb4wDQYJKoZIhvcNAQELBQAwDzENMAsGA1UEAxME
P9g8SpNaf6jNS0ULG8+DJ7dwdHes7IWA0BtjDkur4Ya+ey/FwowgMeEnc/h10Adc
az7b
-----END CERTIFICATE-----

```

By default the certificates are written to files `ca.key` and `ca.crt`.

### Generate a certificate

This needs a pregenerated CA certificate and key (see "Generate a CA set").

To generate a signed certificate pair you can use the following example:

`certificates cert gen --cn=local.test.domain --stdout`

This will output the key and certificate directly to stdout like this (parts are omitted for readability):

```
-----BEGIN RSA PRIVATE KEY-----
MIIJFAIBAAKCAf0Z7/5ZYgOo4gHfAPAPN0vKWEVJ5D97wvnYUq00DcaRPCZZopXl
XUcctgAb3kw27ohTm31KnVEnN8ibeUg2fz+LO/xYVvhD2BMkoe1gk/2JAogPUi1l
jWjI7fuKGwlyHimeYnUx1ADRlShBgHGr
-----END RSA PRIVATE KEY-----

-----BEGIN CERTIFICATE-----
MIIE/TCCAuWgAwIBAgIFFPmGQ70wDQYJKoZIhvcNAQELBQAwDzENMAsGA1UEAxME
V964wCgh6TgfUtt9RabcM3MWtAR18N0vedYg46jhxDa1b+/brQWLuxXDsKIVHrRP
M6ZzVSUF1PH+Ok2Fm7EP26Yax3RkoPrgmlLqL/1fRJaJ
-----END CERTIFICATE-----

```

By default a sqlite (ca.db) database is created to keep track of unique certificate serialnumbers. The CA database can be one of the following flavours of sql: sqlite3, mysql, postgresql or mssql. 

## Development setup

This module uses [Go modules](https://github.com/golang/go/wiki/Modules) for dependency management.
To run: 

- `go run main.go`

And this will output:

```bash
NAME:
   Certificates - An opinionated TLS certificate generator.

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   v0.0.1-alpha4

DESCRIPTION:
   An opinionated TLS certificate generator.

COMMANDS:
     certificate, cert  certificate commands
     help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```