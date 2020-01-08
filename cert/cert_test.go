package cert

import (
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"reflect"
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
		SubjectAltNames  []string
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
		{
			name: "invalid_subject_alt_name",
			fields: fields{
				Organization:     "",
				Country:          "",
				Province:         "",
				Locality:         "",
				StreetAddress:    "",
				PostalCode:       "",
				CommonName:       "valid.common.name",
				SubjectAltNames:  []string{""},
				NameSerialNumber: "",
				NotBefore:        time.Time{},
				NotAfter:         time.Time{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{
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
				SubjectAltNames:  tt.fields.SubjectAltNames,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Request.Validate() error = %v, wantErr %v", err, tt.wantErr)
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
		BitSize          int
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
				Organization:     "",
				Country:          "",
				Province:         "",
				Locality:         "",
				StreetAddress:    "",
				PostalCode:       "",
				CommonName:       "",
				NameSerialNumber: "",
				SubjectAltNames:  []string{},
				NotBefore:        time.Time{},
				NotAfter:         time.Time{},
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
				Organization:     "",
				Country:          "",
				Province:         "",
				Locality:         "",
				StreetAddress:    "",
				PostalCode:       "",
				CommonName:       "cn",
				NameSerialNumber: "",
				SubjectAltNames:  []string{""},
				NotBefore:        time.Time{},
				NotAfter:         time.Time{},
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
			name: "valid_certificate",
			fields: fields{
				Organization:     "",
				Country:          "",
				Province:         "",
				Locality:         "",
				StreetAddress:    "",
				PostalCode:       "",
				CommonName:       "valid.test.local",
				NameSerialNumber: "",
				SubjectAltNames:  []string{"valid.subject.alt.name"},
				NotBefore:        time.Time{},
				NotAfter:         time.Time{},
			},
			args: args{
				caCrt: testCA.Crt,
				caKey: testCA.Key,
			},
			want:    nil,
			want1:   nil,
			wantErr: false,
		},
		{
			name: "valid_certificate_2048",
			fields: fields{
				Organization:     "",
				Country:          "",
				Province:         "",
				Locality:         "",
				StreetAddress:    "",
				PostalCode:       "",
				CommonName:       "valid.test.local",
				NameSerialNumber: "",
				SubjectAltNames:  []string{"valid.subject.alt.name"},
				NotBefore:        time.Time{},
				NotAfter:         time.Time{},
				BitSize:          2048,
			},
			args: args{
				caCrt: testCA.Crt,
				caKey: testCA.Key,
			},
			want:    nil,
			want1:   nil,
			wantErr: false,
		},
		{
			name: "valid_certificate_invalid_bitsize",
			fields: fields{
				Organization:     "",
				Country:          "",
				Province:         "",
				Locality:         "",
				StreetAddress:    "",
				PostalCode:       "",
				CommonName:       "valid.test.local",
				NameSerialNumber: "",
				SubjectAltNames:  []string{"valid.subject.alt.name"},
				NotBefore:        time.Time{},
				NotAfter:         time.Time{},
				BitSize:          1234,
			},
			args: args{
				caCrt: testCA.Crt,
				caKey: testCA.Key,
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
		{
			name: "invalid_ca_key",
			fields: fields{
				Organization:     "",
				Country:          "",
				Province:         "",
				Locality:         "",
				StreetAddress:    "",
				PostalCode:       "",
				CommonName:       "valid.test.local",
				NameSerialNumber: "",
				SubjectAltNames:  []string{},
				NotBefore:        time.Time{},
				NotAfter:         time.Time{},
			},
			args: args{
				caCrt: testCA.Crt,
				caKey: []byte("invalid_key"),
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{
				Organization:     tt.fields.Organization,
				Country:          tt.fields.Country,
				Province:         tt.fields.Province,
				Locality:         tt.fields.Locality,
				StreetAddress:    tt.fields.StreetAddress,
				PostalCode:       tt.fields.PostalCode,
				CommonName:       tt.fields.CommonName,
				NameSerialNumber: tt.fields.NameSerialNumber,
				SubjectAltNames:  tt.fields.SubjectAltNames,
				NotBefore:        tt.fields.NotBefore,
				NotAfter:         tt.fields.NotAfter,
				BitSize:          tt.fields.BitSize,
			}
			got, got1, err := GenerateCertificate(req, tt.args.caCrt, tt.args.caKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request.GenerateCertificate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if len(got) == 0 {
				t.Errorf("Request.GenerateCertificate() got = %v, want %v", got, tt.want)
			}
			if len(got1) == 0 {
				t.Errorf("Request.GenerateCertificate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

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
			req := &Request{
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
			got, got1, err := GenerateCA(req)
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
			req := &Request{
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

func TestNewRequest(t *testing.T) {
	now := time.Now()
	if got := NewRequest(); (got.NotBefore.Before(now) || got.NotBefore.Equal(now)) && got.NotAfter.After(now) {
		fmt.Println(got.NotBefore, got.NotAfter, time.Now())
		t.Error("NewRequest() = empty times")
	}
}

func TestRequest_GetPKIXName(t *testing.T) {
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
	tests := []struct {
		name   string
		fields fields
		want   pkix.Name
	}{
		{
			name: "first",
			fields: fields{
				Organization:     "",
				Country:          "",
				Province:         "",
				Locality:         "",
				StreetAddress:    "",
				PostalCode:       "",
				CommonName:       "",
				SerialNumber:     &big.Int{},
				NameSerialNumber: "",
				SubjectAltNames:  nil,
				NotBefore:        time.Time{},
				NotAfter:         time.Time{},
			},
			want: pkix.Name{
				Country:            nil,
				Organization:       nil,
				OrganizationalUnit: nil,
				Locality:           nil,
				Province:           nil,
				StreetAddress:      nil,
				PostalCode:         nil,
				SerialNumber:       "",
				CommonName:         "",
				Names:              nil,
				ExtraNames:         nil,
			},
		},
		{
			name: "filled",
			fields: fields{
				Organization:     "Org",
				Country:          "NLD",
				Province:         "Zuid-Holland",
				Locality:         "Rotterdam",
				StreetAddress:    "Street",
				PostalCode:       "1234AB",
				CommonName:       "common.name",
				SerialNumber:     &big.Int{},
				NameSerialNumber: "abc1234cba",
				SubjectAltNames:  nil,
				NotBefore:        time.Time{},
				NotAfter:         time.Time{},
			},
			want: pkix.Name{
				Country:            []string{"NLD"},
				Organization:       []string{"Org"},
				OrganizationalUnit: nil,
				Locality:           []string{"Rotterdam"},
				Province:           []string{"Zuid-Holland"},
				StreetAddress:      []string{"Street"},
				PostalCode:         []string{"1234AB"},
				SerialNumber:       "abc1234cba",
				CommonName:         "common.name",
				Names:              nil,
				ExtraNames:         nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &Request{
				Organization:     tt.fields.Organization,
				Country:          tt.fields.Country,
				Province:         tt.fields.Province,
				Locality:         tt.fields.Locality,
				StreetAddress:    tt.fields.StreetAddress,
				PostalCode:       tt.fields.PostalCode,
				CommonName:       tt.fields.CommonName,
				SerialNumber:     tt.fields.SerialNumber,
				NameSerialNumber: tt.fields.NameSerialNumber,
				SubjectAltNames:  tt.fields.SubjectAltNames,
				NotBefore:        tt.fields.NotBefore,
				NotAfter:         tt.fields.NotAfter,
			}
			if got := req.GetPKIXName(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request.GetPKIXName() = %v, want %v", got, tt.want)
			}
		})
	}
}
