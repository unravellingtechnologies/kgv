package certs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestTLSSetupWithCertGeneration Tests the ssl setup using generated certificates
func TestTLSSetupWithCertGeneration(t *testing.T) {
	// TODO: fix this test to be more appropriate
	_, _,  err := TLSSetup("", "")

	assert.Nil(t, err, "No errors expected during setup of the configuration using certificate generation")
}

// TestTLSSetupWithFiles Tests the ssl setup using supplied certificates
func TestTLSSetupWithFiles(t *testing.T) {
	// TODO: fix this test to be more appropriate
	_, _,  err := TLSSetup("../../testdata/pkg/certs/tls.crt", "../../testdata/pkg/certs/tls.key")

	assert.Nil(t, err, "No errors expected during setup of the configuration using certificate generation")
}