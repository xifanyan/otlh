package main

import (
	"os"

	otlh "github.com/xifanyan/otlh/pkg"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:    "otlh",
		Version: "0.5.4-beta",
		Usage:   "Command Line Interface to access Opentext LegalHold service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "domain",
				Aliases: []string{"x"},
				Usage:   "domain name for Opentext legahold service",
				EnvVars: []string{"LHN_DOMAIN"},
				Value:   otlh.DEFAULT_DOMAIN,
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "port",
				Value:   otlh.DEFAULT_PORT,
			},
			&cli.StringFlag{
				Name:    "proxy",
				Aliases: []string{"y"},
				Usage:   "http proxy",
				EnvVars: []string{"LHN_HTTPPROXY"},
				Value:   "",
			},
			&cli.StringFlag{
				Name:    "tenant",
				Aliases: []string{"t"},
				Usage:   "tenant name",
				EnvVars: []string{"LHN_TENANT"},
				Value:   "",
			},
			&cli.StringFlag{
				Name:    "authToken",
				Aliases: []string{"a"},
				Usage:   "token to access legalhold web service",
				EnvVars: []string{"LHN_AUTHTOKEN"},
				Value:   "",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "LHN json config file",
				EnvVars: []string{"LHN_CONFIG"},
				Value:   "",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Debug Mode",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "trace",
				Aliases: []string{"z"},
				Usage:   "Trace Mode",
				Value:   false,
			},
		},
		Commands: Commands,
		Before: func(c *cli.Context) error {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			if c.Bool("debug") {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			}
			if c.Bool("trace") {
				zerolog.SetGlobalLevel(zerolog.TraceLevel)
			}
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error().Msgf("error: %s", err)
		os.Exit(1)
	}
}
