package main

import (
	"os"
	"testing"
)

func Test_run(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
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
