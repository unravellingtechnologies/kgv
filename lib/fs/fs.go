package fs

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func WriteToFile(filename string, data *bytes.Buffer) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Error("error creating file %s", filename, err)
		os.Exit(1)
	}

	defer file.Close()

	_, err = io.WriteString(file, data.String())
	if err != nil {
		return err
	}
	return file.Sync()
}