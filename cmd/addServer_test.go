package cmd

import (
	"log"
	"os"
	"testing"
)

func TestAddServerPrompt(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}
}
