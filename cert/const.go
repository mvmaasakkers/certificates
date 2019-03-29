package cert

import "errors"

var (
	ErrorInvalidCommonName = errors.New("invalid common name")
	ErrorInvalidSubjectAltName = errors.New("invalid subject alt name")
)
