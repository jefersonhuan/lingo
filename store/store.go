package store

import (
	"github.com/sdomino/scribble"
	"os"
	"path/filepath"
)

var Conn *scribble.Driver

func InitDatabase() {
	var err error

	dir, err := configDir()
	if err != nil {
		panic(err)
	}

	Conn, err = scribble.New(dir, nil)
	if err != nil {
		panic(err)
	}
}

func configDir() (dir string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	path := filepath.Join(home, ".lingo")

	err = os.Mkdir(path, 0755)
	if os.IsExist(err) {
		return path, nil
	}

	return
}
