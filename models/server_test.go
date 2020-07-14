package models

import (
	"testing"
)

func TestServer_ToURIWithoutUserAndPass(t *testing.T) {
	server := Server{
		Host: "localhost",
		Port: 27017,
	}

	got := server.ToURI()

	if got != "mongodb://localhost:27017" {
		t.Errorf("server.ToURI() = %s; want mongodb://localhost:27017", got)
	}
}

func TestServer_ToURIWithUserAndPass(t *testing.T) {
	server := Server{
		Host:     "localhost",
		Port:     27017,
		User:     "root",
		Password: "root",
	}

	got := server.ToURI()
	expected := "mongodb://root:root@localhost:27017"

	if got != expected {
		t.Errorf("server.ToURI() = %s; want %s", got, expected)
	}
}
