package database

import (
	"testing"
)

func TestServer_FromURI(t *testing.T) {
	type args struct {
		username string
		password string
		host     string
		port     int
	}

	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			name: "without username and password",
			args: args{
				username: "",
				password: "",
				host:     "localhost",
				port:     27017,
			},
			expected: "mongodb://localhost:27017",
		},
		{
			name: "full URI",
			args: args{
				username: "root",
				password: "root",
				host:     "localhost",
				port:     9000,
			},
			expected: "mongodb://root:root@localhost:9000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{}

			server.BuildURI(tt.args.username, tt.args.password, tt.args.host, tt.args.port)

			if tt.expected != server.URI {
				t.Errorf("Expected %s, got %s", tt.expected, server.URI)
			}
		})
	}
}
