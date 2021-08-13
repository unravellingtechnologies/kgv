package fs

import (
	log "github.com/sirupsen/logrus"
	"os"
)

// Exists checks if a given file exists
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			log.Error(err)
			return false
		}
	}
	return true
}
