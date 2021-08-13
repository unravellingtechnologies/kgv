// Package fs Tests for the fs package
package fs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestExists tests the happy path of the Exists function
func TestExists(t *testing.T) {
	assert.Truef(t, Exists("../../testdata/lib/fs/existing"),
		"Expected file to exist")

}

// TestsDoesNotExist test the case where the file does not exist
func TestDoesNotExist(t *testing.T) {
	assert.False(t, Exists("testdata/lib/fs/not-existing"),
		"Expected file not to be found")
}
