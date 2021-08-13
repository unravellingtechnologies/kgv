package certs

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	log "github.com/sirupsen/logrus"
	"github.com/unravellingtechnologies/kgv/lib/fs"
	"math/big"
	"net"
	"os"
	"time"
)

func GeneratePrivateKey() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Error("error while generating private key", err)
		os.Exit(1)
	}

	return key
}

func GenerateCertificate(template, parent *x509.Certificate, publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) (*x509.Certificate, []byte) {
	certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, publicKey, privateKey)
	if err != nil {
		log.Error("failed to create certificate :%s", err.Error())
		os.Exit(1)
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		log.Error("failed to parse certificate: %s", err.Error())
		os.Exit(1)
	}

	pemBlock := pem.Block{Type: "CERTIFICATE", Bytes: certBytes}
	certPEM := pem.EncodeToMemory(&pemBlock)

	return cert, certPEM
}

func CreateCaTemplate() *x509.Certificate {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2021),
		Subject: pkix.Name{
			Organization:  []string{"Unravelling Technologies GmbH"},
			Country:       []string{"DE"},
			Province:      []string{"Bayern"},
			Locality:      []string{"Hengersberg"},
			StreetAddress: []string{"Bayerische Wald"},
			PostalCode:    []string{"94491"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	return ca
}

func GenerateCa() (*x509.Certificate, *bytes.Buffer){
	caTemplate := CreateCaTemplate()
	caPrivateKey := GeneratePrivateKey()
	caBytes, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		log.Error("failed creating the caTemplate certificate", err)
		os.Exit(1)
	}

	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivateKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivateKey),
	})

	return caTemplate, caPrivateKeyPEM
}

func CreateCert() *x509.Certificate {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"Unravelling Technologies GmbH"},
			Country:       []string{"DE"},
			Province:      []string{"Bayern"},
			Locality:      []string{"Hengersberg"},
			StreetAddress: []string{"Bayerische Wald"},
			PostalCode:    []string{"94491"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	return cert
}

func GenerateCerts(path string) {
	ca, caPrivateKey := GenerateCa()

	cert := CreateCert()
	certPrivateKey := GeneratePrivateKey()

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		log.Error("failed creating the certificate", err)
		os.Exit(1)
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivateKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivateKey),
	})

	err = fs.WriteToFile(path + "/tls.crt", certPEM)
	err = fs.WriteToFile(path + "tls.key", certPrivateKeyPEM)
}