package cert

import (
	"crypto/rand"
	"math/big"
)

// GenerateRandomBigInt generates a random big.int (decimal) needed for certificate serial numbers.
func GenerateRandomBigInt() (*big.Int, error) {
	var n *big.Int
	var err error

	max := *big.NewInt(99999999999)

	n, err = rand.Int(rand.Reader, &max)

	if err != nil {
		return nil, err
	}

	return n, nil
}
