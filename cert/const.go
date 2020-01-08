package cert

import "errors"

var (
	// ErrorInvalidCommonName is given if an empty or invalid common name is given
	ErrorInvalidCommonName = errors.New("invalid common name")
	// ErrorInvalidSubjectAltName is given if an empty or invalid subject alt  name is given
	ErrorInvalidSubjectAltName = errors.New("invalid subject alt name")
	// ErrorInvalidBitSize is given if an invalid bitsize is given
	ErrorInvalidBitSize = errors.New("invalid bit size")
)
