package main

import (
	"os"
	"testing"
)

func Test_run(t *testing.T) {
	// Clean up files first
	os.Remove("ca.crt")
	os.Remove("ca.key")
	os.Remove("file.db")


	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty-args",
			args: args{
				args: nil,
			},
			wantErr: false,
		},
		{
			name: "invalid-ca",
			args: args{
				args: []string{"cert", "gen-ca"},
			},
			wantErr: true,
		},
		{
			name: "invalid-crt",
			args: args{
				args: []string{"cert", "gen"},
			},
			wantErr: true,
		},
		{
			name: "valid-ca",
			args: args{
				args: []string{"cert", "gen-ca", "--cn=common.test.name"},
			},
			wantErr: false,
		},
		{
			name: "valid-ca-stdout",
			args: args{
				args: []string{"cert", "gen-ca", "--cn=common.test.name", "--stdout"},
			},
			wantErr: false,
		},
		{
			name: "valid-crt",
			args: args{
				args: []string{"cert", "gen", "--cn=common.test.name"},
			},
			wantErr: false,
		},
		{
			name: "valid-crt",
			args: args{
				args: []string{"cert", "gen", "--cn=common.test.name.stdout", "--stdout"},
			},
			wantErr: false,
		},
		{
			name: "valid-crt-duplicate",
			args: args{
				args: []string{"cert", "gen", "--cn=common.test.name", "--stdout"},
			},
			wantErr: true,
		},
		{
			name: "valid-crt-expiration-dates",
			args: args{
				args: []string{"cert", "gen", "--cn=common.test.name.two", "--stdout", "--notbefore=2019-01-01", "--notafter=2019-10-20"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := os.Args[0:1]
			args = append(args, tt.args.args...)
			if err := run(args); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
