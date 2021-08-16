package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	certs "github.com/unravellingtechnologies/kgv/lib/certs"
	"github.com/unravellingtechnologies/kgv/pkg/cli"
	"github.com/unravellingtechnologies/kgv/pkg/webhook"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	serverTLSConf, _, err := certs.TLSSetup(tlsCert, tlsKey)
	if err != nil {
		panic(err)
	}

	mux := webhook.SetupListeners()

	s := &http.Server{
		Addr:      ":" + port,
		Handler:   mux,
		TLSConfig: serverTLSConf,
	}

	go func() {
		if err := s.ListenAndServeTLS("", ""); err != nil {
			log.Errorf("Failed to listen and serve: %v", err)
		}
	}()

	log.Info("Listening for requests on port ", port)

	// listen shutdown signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Infof("Shutdown gracefully...")
	if err := s.Shutdown(context.Background()); err != nil {
		log.Error(err)
	}
}
