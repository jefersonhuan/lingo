package cmd

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io"
	"os"
	"path/filepath"
)

const configFile = ".mongo-transfer"

func openFile(flags int) (file *os.File, err error) {
	home, _ := homedir.Dir()

	path := filepath.Join(home, configFile)
	fullpath := filepath.Join(path, "servers.xml")

	file, err = os.OpenFile(fullpath, flags, 0755)
	if os.IsNotExist(err) && flags != os.O_RDONLY {
		if err = os.Mkdir(path, 0755); err != nil {
			return
		}

		return os.OpenFile(fullpath, flags, 0755)
	}

	return
}

func readServers(target *[]Server) (err error) {
	var file *os.File

	file, err = openFile(os.O_RDONLY)
	if err != nil {
		return
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	var bytes []byte

	for {
		bytes, err = reader.ReadBytes('\n')

		if err == io.EOF {
			return nil
		} else if err != nil {
			return
		}

		if err = xml.Unmarshal(bytes, &target); err != nil {
			return
		}
	}
}

func saveServer(v Server) (err error) {
	var file *os.File

	file, err = openFile(os.O_CREATE | os.O_WRONLY | os.O_APPEND)
	if err != nil {
		return
	}

	defer file.Close()

	body, err := xml.MarshalIndent(v, "", "   ")
	if err != nil {
		return fmt.Errorf("couldn't save config file: %v", err)
	}

	if _, err := file.Write(body); err != nil {
		return err
	}

	_, _ = file.WriteString("\n")

	_ = file.Sync()

	return nil
}
