package cert

import (
	"math/big"
	"testing"
	"time"
)

func TestCertRequest_Validate(t *testing.T) {
	type fields struct {
		Organization     string
		Country          string
		Province         string
		Locality         string
		StreetAddress    string
		PostalCode       string
		CommonName       string
		SerialNumber     *big.Int
		NameSerialNumber string
		NotBefore        time.Time
		NotAfter         time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "invalid_common_name",
			fields: fields{
				Organization:     "",
				Country:          "",
				Province:         "",
				Locality:         "",
				StreetAddress:    "",
				PostalCode:       "",
				CommonName:       "",
				SerialNumber:     nil,
				NameSerialNumber: "",
				NotBefore:        time.Time{},
				NotAfter:         time.Time{},
			},
			wantErr: true,
		},
		{
			name: "valid_common_name",
			fields: fields{
				Organization:     "",
				Country:          "",
				Province:         "",
				Locality:         "",
				StreetAddress:    "",
				PostalCode:       "",
				CommonName:       "valid.common.name",
				NameSerialNumber: "",
				NotBefore:        time.Time{},
				NotAfter:         time.Time{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &CertRequest{
				Organization:     tt.fields.Organization,
				Country:          tt.fields.Country,
				Province:         tt.fields.Province,
				Locality:         tt.fields.Locality,
				StreetAddress:    tt.fields.StreetAddress,
				PostalCode:       tt.fields.PostalCode,
				CommonName:       tt.fields.CommonName,
				NameSerialNumber: tt.fields.NameSerialNumber,
				NotBefore:        tt.fields.NotBefore,
				NotAfter:         tt.fields.NotAfter,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("CertRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCertRequest_GenerateCertificate(t *testing.T) {
	type fields struct {
		Organization     string
		Country          string
		Province         string
		Locality         string
		StreetAddress    string
		PostalCode       string
		CommonName       string
		SerialNumber     *big.Int
		NameSerialNumber string
		SubjectAltNames  []string
		NotBefore        time.Time
		NotAfter         time.Time
	}
	type args struct {
		caCrt []byte
		caKey []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		want1   []byte
		wantErr bool
	}{
		{
			name: "invalid_common_name",
			fields: fields{
				Organization:    "",
				Country:         "",
				Province:        "",
				Locality:        "",
				StreetAddress:   "",
				PostalCode:      "",
				CommonName:      "",
				NameSerialNumber:    "",
				SubjectAltNames: []string{},
				NotBefore:       time.Time{},
				NotAfter:        time.Time{},
			},
			args: args{
				caCrt: nil,
				caKey: nil,
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
		{
			name: "invalid_subject_alt_name",
			fields: fields{
				Organization:    "",
				Country:         "",
				Province:        "",
				Locality:        "",
				StreetAddress:   "",
				PostalCode:      "",
				CommonName:      "cn",
				NameSerialNumber:    "",
				SubjectAltNames: []string{""},
				NotBefore:       time.Time{},
				NotAfter:        time.Time{},
			},
			args: args{
				caCrt: nil,
				caKey: nil,
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &CertRequest{
				Organization:  tt.fields.Organization,
				Country:       tt.fields.Country,
				Province:      tt.fields.Province,
				Locality:      tt.fields.Locality,
				StreetAddress: tt.fields.StreetAddress,
				PostalCode:    tt.fields.PostalCode,
				CommonName:    tt.fields.CommonName,
				NameSerialNumber:  tt.fields.NameSerialNumber,
				NotBefore:     tt.fields.NotBefore,
				NotAfter:      tt.fields.NotAfter,
			}
			got, got1, err := req.GenerateCertificate(tt.args.caCrt, tt.args.caKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("CertRequest.GenerateCertificate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if len(got) == 0 {
				t.Errorf("CertRequest.GenerateCertificate() got = %v, want %v", got, tt.want)
			}
			if len(got1) == 0 {
				t.Errorf("CertRequest.GenerateCertificate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
