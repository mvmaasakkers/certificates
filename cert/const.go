package cert

import "errors"

var (
	// ErrorInvalidCommonName is given if an empty or invalid common name is given
	ErrorInvalidCommonName = errors.New("invalid common name")
	// ErrorInvalidSubjectAltName is given if an empty or invalid subject alt  name is given
	ErrorInvalidSubjectAltName = errors.New("invalid subject alt name")
)
