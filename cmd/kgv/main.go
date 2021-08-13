package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/unravellingtechnologies/kgv/pkg/certs"
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

func HelloWorld(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Hello world!\n"))
}

func main() {
	ctx, cliOpts := cli.Parse()
	port, tlsCert, tlsKey = cli.ParseOpts(cliOpts)
	switch ctx.Command() {
	case "start":
	default:
		panic(ctx.Command())
	}

	// cert, _ := tls.LoadX509KeyPair(tlsCert, tlsKey)

	serverTLSConf, _, err := certs.TLSSetup(tlsCert, tlsKey)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", HelloWorld)

	s := &http.Server{
		Addr:      ":" + port,
		Handler:   mux,
		TLSConfig: serverTLSConf,
	}

	log.Info("Listening for requests on port ", port)
	log.Fatal(s.ListenAndServeTLS("", ""))
}
