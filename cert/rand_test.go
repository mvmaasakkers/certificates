package cert

import (
	"math/big"
	"reflect"
	"testing"
)

func TestGenerateRandomBigInt(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "one",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateRandomBigInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRandomBigInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRandomResults(t *testing.T) {
	checks := 10
	results := []big.Int{}

	for i := 0; i < checks; i++ {
		r, _ := GenerateRandomBigInt()
		if contains(results, *r) {
			t.Error("GenerateRandomBigInt() must produce unique values")
			return
		}
		results = append(results, *r)
	}
}

func contains(s []big.Int, e big.Int) bool {
	for _, a := range s {
		if reflect.DeepEqual(a, e) {
			return true
		}
	}
	return false
}
