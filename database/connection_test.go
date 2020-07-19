package database

import (
	"errors"
	"mongo-transfer/utils"
	"testing"
)

func TestServer_Connect(t *testing.T) {
	type fields struct {
		ID  string
		URI string
	}

	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		errMessage error
	}{
		{
			name: "connects existing database and loads struct fields",
			fields: fields{
				URI: "mongodb://localhost:27017",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{
				URI: tt.fields.URI,
			}

			err := server.Connect()

			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && server.Client == nil {
				t.Errorf("Connect() expected server.client to be initialized")
			}
		})
	}
}

func TestServer_Ping(t *testing.T) {
	type fields struct {
		URI string
	}
	tests := []struct {
		name       string
		fields     fields
		startMongo bool
		wantErr    bool
		errMessage error
	}{
		{
			name:       "fulfills ping for existing database",
			fields:     fields{URI: "mongodb://localhost:27017"},
			startMongo: true,
		},
		{
			name: "error with non-existing database",
			fields: fields{
				URI: "mongodb://some-random-domain:7000",
			},
			wantErr:    true,
			errMessage: errors.New("context deadline exceeded"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := Server{URI: tt.fields.URI}

			if tt.startMongo {
				if err := utils.StartMongo(); err != nil {
					t.Skip(err)
				}

				defer utils.StopMongo()
			}

			if err := server.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
