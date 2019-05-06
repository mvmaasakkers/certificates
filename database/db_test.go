package database

import (
	"testing"
)

func TestNewCertificate(t *testing.T) {
	tests := []struct {
		name string
		want *Certificate
	}{
		{
			name: "valid",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCertificate(); got == nil || got.UUID == "" {
				t.Errorf("NewCertificate() = %v, want %v", got, tt.want)
			}
		})
	}
}
