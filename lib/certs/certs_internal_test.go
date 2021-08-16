package certs

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGetCertificateSubject tests the generation of the default subject name
func TestGetCertificateSubject(t *testing.T) {
	name := getCertificateSubject()
	assert.IsType(t, pkix.Name{}, name, "expected a valid PKI Name object")
	assert.Equal(t, "Unravelling Technologies GmbH", name.Organization[0],
		"Expected organization name to match")
	assert.Equal(t, "DE", name.Country[0],
		"Expected country name to match")
}

// TestGeneratePrivateKey tests the generation of the private keys
func TestGeneratePrivateKey(t *testing.T) {
	privateKey, privateKeyPEM, err := generatePrivateKey()

	assert.Nil(t, err, "No errors expected during creation of private key")
	assert.IsType(t, rsa.PrivateKey{}, *privateKey, "Expected private key to be a valid rsa key")

	block, _ := pem.Decode(privateKeyPEM.Bytes())
	assert.Equal(t, "RSA PRIVATE KEY", block.Type, "Expected a private key in the PEM")
	// does this really bring something?
	assert.Equal(t, x509.MarshalPKCS1PrivateKey(privateKey), block.Bytes, "Expected same content in the bytes")
}

// TestGetCATemplate tests the generation of the default CA template
func TestGetCATemplate(t *testing.T) {
	caTemplate := getCATemplate()
	assert.IsType(t, x509.Certificate{}, *caTemplate)
}

// TestGetServerCertificateTemplate tests the generation of the server certificate template
func TestGetServerCertificateTemplate(t *testing.T) {
	serverTemplate := getServerCertificateTemplate()
	assert.IsType(t, x509.Certificate{}, *serverTemplate)
}

// TestGenerateCA Tests the generation of the CA certificate
func TestGenerateCA(t *testing.T) {
	caTemplate, caPEM, caPrivateKey, err := generateCA()
	assert.Nil(t, err, "Expected no errors to occur during creation of the CA")
	assert.IsType(t, x509.Certificate{}, *caTemplate)
	assert.IsType(t, rsa.PrivateKey{}, *caPrivateKey, "Expected private key to be a valid rsa key")
	block, _ := pem.Decode(caPEM.Bytes())
	assert.Equal(t, "CERTIFICATE", block.Type, "Expected a private key in the PEM")
}

// TestGenerateServerCertificates Tests the generation of the server certificate
func TestGenerateServerCertificates(t *testing.T) {
	// TODO: figure out how to appropriately check the validity of the generated certificates

}