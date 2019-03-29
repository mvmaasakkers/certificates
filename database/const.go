package database

import "errors"

var (
	ErrorNilConnection = errors.New("connection is nil")
	ErrorObjectNotFound = errors.New("object not found")
)
