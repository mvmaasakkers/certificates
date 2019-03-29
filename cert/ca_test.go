package cert

import (
	"testing"
	"time"
)

func TestCARequest_GenerateCA(t *testing.T) {
	type fields struct {
		Organization  string
		Country       string
		Province      string
		Locality      string
		StreetAddress string
		PostalCode    string
		CommonName    string
		NotBefore     time.Time
		NotAfter      time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		want1   []byte
		wantErr bool
	}{
		{
			name: "check_for_invalid_input",
			fields: fields{
				Organization:  "",
				Country:       "",
				Province:      "",
				Locality:      "",
				StreetAddress: "",
				PostalCode:    "",
				CommonName:    "",
				NotBefore:     time.Time{},
				NotAfter:      time.Time{},
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &CARequest{
				Organization:  tt.fields.Organization,
				Country:       tt.fields.Country,
				Province:      tt.fields.Province,
				Locality:      tt.fields.Locality,
				StreetAddress: tt.fields.StreetAddress,
				PostalCode:    tt.fields.PostalCode,
				CommonName:    tt.fields.CommonName,
				NotBefore:     tt.fields.NotBefore,
				NotAfter:      tt.fields.NotAfter,
			}
			got, got1, err := req.GenerateCA()
			if (err != nil) != tt.wantErr {
				t.Errorf("CARequest.GenerateCA() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if len(got) == 0 {
				t.Errorf("CARequest.GenerateCA() got = %v, want %v", got, tt.want)
			}
			if len(got1) == 0 {
				t.Errorf("CARequest.GenerateCA() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCARequest_Validate(t *testing.T) {
	type fields struct {
		Organization  string
		Country       string
		Province      string
		Locality      string
		StreetAddress string
		PostalCode    string
		CommonName    string
		NotBefore     time.Time
		NotAfter      time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "invalid_common_name",
			fields: fields{
				Organization:  "",
				Country:       "",
				Province:      "",
				Locality:      "",
				StreetAddress: "",
				PostalCode:    "",
				CommonName:    "",
				NotBefore:     time.Time{},
				NotAfter:      time.Time{},
			},
			wantErr: true,
		},
		{
			name: "valid_common_name",
			fields: fields{
				Organization:  "",
				Country:       "",
				Province:      "",
				Locality:      "",
				StreetAddress: "",
				PostalCode:    "",
				CommonName:    "valid.common.name",
				NotBefore:     time.Time{},
				NotAfter:      time.Time{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &CARequest{
				Organization:  tt.fields.Organization,
				Country:       tt.fields.Country,
				Province:      tt.fields.Province,
				Locality:      tt.fields.Locality,
				StreetAddress: tt.fields.StreetAddress,
				PostalCode:    tt.fields.PostalCode,
				CommonName:    tt.fields.CommonName,
				NotBefore:     tt.fields.NotBefore,
				NotAfter:      tt.fields.NotAfter,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("CARequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
