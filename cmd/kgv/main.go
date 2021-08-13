package main

import (
	"crypto/tls"
	log "github.com/sirupsen/logrus"
	"github.com/unravellingtechnologies/kgv/pkg/cli"
	"net/http"
)

var (
	port, tlsCert, tlsKey string
)

func init() {
	// Initialize logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}

func main() {
	ctx, cliOpts := cli.Parse()
	port, tlsCert, tlsKey = cli.ParseOpts(cliOpts)
	switch ctx.Command() {
		case "start":
		default:
			panic(ctx.Command())
	}

	cert, _ := tls.LoadX509KeyPair(tlsCert, tlsKey)

	s := &http.Server{
		Addr: ":" + port,
		Handler: nil,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	log.Info("Listening for requests on port %s", port)
	log.Fatal(s.ListenAndServeTLS("", ""))
}