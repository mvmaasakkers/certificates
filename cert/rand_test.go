package cert

import (
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
