package database

import "errors"

var (
	// ErrorNilConnection is used when a connection is not available but needed
	ErrorNilConnection = errors.New("connection is nil")
	// ErrorObjectNotFound is used when an object is not found in DB based on given filter
	ErrorObjectNotFound = errors.New("object not found")
	// ErrorDuplicateObject is used when a Serial Number (which is the primary id) is already in the db
	ErrorDuplicateObject = errors.New("object is a duplicate")
)
