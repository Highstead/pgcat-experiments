package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/highstead/pgcat-experiments/pkg/server"
	cli "github.com/urfave/cli/v2"
)

const (
	addrFlag    = "addr"
	tlsCertFlag = "tls-cert"
	tlsKeyFlag  = "tls-key"
	connString  = "conn-String"
)

var app = &cli.App{
	Name:    "inserter",
	Version: "0.1.0",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  addrFlag,
			Usage: "The host to bind on (host:port)",
			Value: ":8887",
		},
		&cli.StringFlag{
			Name:  tlsCertFlag,
			Usage: "The PEM encoded TLS cert file",
		},
		&cli.StringFlag{
			Name:  tlsKeyFlag,
			Usage: "The PEM encoded TLS key file",
		},
		&cli.StringFlag{
			Name:  connString,
			Usage: "Connection string",
			Value: "postgres://postgres:postgres@localhost:6432/my_database?sslmode=disable&binary_parameters=yes",
		},
	},
	Action: func(ctx *cli.Context) error {
		logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))

		opts := []server.Option{
			server.WithLogger(logger),
			server.WithContext(ctx.Context),
			server.WithConnString(ctx.String(connString)),
		}

		if ctx.String(tlsKeyFlag) != "" {
			opts = append(opts, server.WithTLS(ctx.String(tlsKeyFlag), ctx.String(tlsCertFlag)))
		}

		_, done := server.ListenAndServe(ctx.String(addrFlag), opts...)
		<-done
		return nil
	},
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "There was an error: %v", err)
	}
}
