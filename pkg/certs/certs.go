// Package certs provides functions to manage certificates and keys in the context of the kgv project
package certs

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	log "github.com/sirupsen/logrus"
	"github.com/unravellingtechnologies/kgv/lib/fs"
	"math/big"
	"net"
	"time"
)

// Generates the default subject for the certificates
func getCertificateSubject() pkix.Name {
	return pkix.Name{
		Organization:  []string{"Unravelling Technologies GmbH"},
		Country:       []string{"DE"},
		Province:      []string{"Bayern"},
		Locality:      []string{"Hengersberg"},
		StreetAddress: []string{"Bayerische Wald"},
		PostalCode:    []string{"94491"},
	}
}

// Generates a RSA private key
func generatePrivateKey() (*rsa.PrivateKey, *bytes.Buffer, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Error("error while generating private key", err)
		return nil, nil, err
	}

	privateKeyPEM := new(bytes.Buffer)
	err = pem.Encode(privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		log.Error("failed encoding private key into PEM", err)
		return nil, nil, err
	}

	return privateKey, privateKeyPEM, nil
}

// Generates the CA Template
func getCATemplate() *x509.Certificate {
	return &x509.Certificate{
		SerialNumber:          big.NewInt(2019),
		Subject:               getCertificateSubject(),
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}

// Gets the server certificate template
func getServerCertificateTemplate() *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject:      getCertificateSubject(),
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
}

// Generates a CA certificate
func generateCACertificate(caTemplate *x509.Certificate, caPrivateKey *rsa.PrivateKey) (*bytes.Buffer, error) {
	caBytes, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		log.Error("failed to generate certificate", err)
		return nil, err
	}

	// pem encode
	caPEM := new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		log.Error("failed to encode CA certificate into PEM", err)
		return nil, err
	}

	return caPEM, nil
}

// Generates the CA certificate
func generateCA() (*x509.Certificate, *bytes.Buffer, *rsa.PrivateKey, error) {
	// set up our CA certificate
	caTemplate := getCATemplate()

	// create our private and public key
	caPrivateKey, _, err := generatePrivateKey()
	if err != nil {
		return nil, nil, nil, err
	}

	// create the CA
	caPEM, err := generateCACertificate(caTemplate, caPrivateKey)
	if err != nil {
		return nil, nil, nil, err
	}

	return caTemplate, caPEM, caPrivateKey, nil
}

// Generates the server certificates
func generateServerCertificates() (*tls.Certificate, *bytes.Buffer, error) {
	caTemplate, caPEM, caPrivateKey, err := generateCA()
	if err != nil {
		log.Error("failed to generate CA certificate", err)
		return nil, nil, err
	}

	// set up our server certificate
	cert := getServerCertificateTemplate()

	privateKey, privateKeyPEM, err := generatePrivateKey()
	if err != nil {
		return nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caTemplate, &privateKey.PublicKey, caPrivateKey)
	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		log.Error("Error while generating the server certificates")
		return nil, nil, err
	}

	serverCertificates, err := tls.X509KeyPair(certPEM.Bytes(), privateKeyPEM.Bytes())
	if err != nil {
		return nil, nil, err
	}

	return &serverCertificates, caPEM, nil
}

// Loads certificates from disk
func loadCertificates(certificate string, privateKey string) (*tls.Certificate, error) {
	serverCertificates, err := tls.LoadX509KeyPair(certificate, privateKey)
	if err != nil {
		log.Error("error loading provided certificates", err)
		return nil, err
	}

	return &serverCertificates, nil
}

// TLSSetup produces the TLS configuration for usage in the HTTPS endpoint
func TLSSetup(certFile string, keyFile string) (serverTLSConf *tls.Config, clientTLSConf *tls.Config, err error) {
	var serverCertificates *tls.Certificate
	var caPEM *bytes.Buffer = nil

	// first use case, there are provided certificates to use
	if certFile != "" && keyFile != "" && fs.Exists(certFile) && fs.Exists(keyFile) {
		serverCertificates, err = loadCertificates(certFile, keyFile)
		if err != nil {
			return nil, nil, err
		}
	} else {
		// second use case, generate our own self-signed certificate
		serverCertificates, caPEM, err = generateServerCertificates()
		if err != nil {
			return nil, nil, err
		}
	}

	serverTLSConf = &tls.Config{
		Certificates: []tls.Certificate{*serverCertificates},
	}

	if caPEM != nil {
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(caPEM.Bytes())
		clientTLSConf = &tls.Config{
			RootCAs: certPool,
		}
	}

	return
}
