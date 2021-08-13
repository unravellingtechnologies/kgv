package cli

import (
	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"
	"github.com/unravellingtechnologies/kgv/lib/fs"
	"github.com/unravellingtechnologies/kgv/pkg/certs"
	"os"
)

type CLI struct {
	Start struct {
		Port string `short:"p" env:"PORT" help:"Port where to listen for requests" default:"8443"`
		// Don't want to define defaults for cert and key. If they are user defined, then should exit on error
		// if files are not found. If using the defaults, then generates a default certificate
		TlsCert string `help:"TLS Certificate for the HTTPS server"`
		TlsKey string `help:"TLS Key for the HTTPS server"`
	} `cmd default:"withargs" help:"Starts the server."`
}

const defaultPath = "/etc/kgv/certs"

var options CLI

func Parse() (*kong.Context, *CLI) {
	ctx := kong.Parse(&options, kong.Configuration(kong.JSON, "/etc/kgv/config.json", "~/.kgv/config.json"))
	return ctx, &options
}

func ParseOpts(cli *CLI) (port string, tlsCert string, tlsKey string) {
	if cli.Start.TlsCert != "" || cli.Start.TlsKey != "" {
		if ! fs.Exists(cli.Start.TlsCert) || ! fs.Exists(cli.Start.TlsKey) {
			log.Error("configured certificates not found")
			os.Exit(1)
		}

		// User provided certificates found, will be used
		tlsCert = cli.Start.TlsCert
		tlsKey = cli.Start.TlsKey
	} else {
		// User didn't provide certificate paths, will set defaults and generate
		tlsCert = defaultPath + "/tls.crt"
		tlsKey = defaultPath + "/etc/kgv/certs/tls.key"

		certs.GenerateCerts(defaultPath)
	}

	return cli.Start.Port, tlsCert, tlsKey
}