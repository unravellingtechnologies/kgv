// Package cli has the functions to process the command line interface options
package cli

import (
	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"
	"github.com/unravellingtechnologies/kgv/lib/fs"
	"os"
)

type CLI struct {
	Start struct {
		Port string `short:"p" env:"PORT" help:"Port where to listen for requests" default:"8443"`
		// Don't want to define defaults for cert and key. If they are user defined, then should exit on error
		// if files are not found. If using the defaults, then generates a default certificate
		TlsCert string `help:"TLS Certificate for the HTTPS server"`
		TlsKey  string `help:"TLS Key for the HTTPS server"`
	} `cmd default:"withargs" help:"Starts the server."`
}

const defaultPath = "/etc/kgv/certs"

var options CLI

// Parse function parses the CLI options (heavy lifting done by kong)
func Parse() (*kong.Context, *CLI) {
	ctx := kong.Parse(&options, kong.Configuration(kong.JSON, "/etc/kgv/config.json", "~/.kgv/config.json"))
	return ctx, &options
}

// ParseOpts processes the parsed the options on the cli into usable values in our application
func ParseOpts(cli *CLI) (port string, tlsCert string, tlsKey string) {
	if cli.Start.TlsCert != "" || cli.Start.TlsKey != "" {
		if !fs.Exists(cli.Start.TlsCert) || !fs.Exists(cli.Start.TlsKey) {
			log.Error("configured certificates not found")
			os.Exit(1)
		}

		// User provided certificates found, will be used
		tlsCert = cli.Start.TlsCert
		tlsKey = cli.Start.TlsKey
	}

	return cli.Start.Port, tlsCert, tlsKey
}
