package server

import (
	"context"
	"log/slog"
	"os"
	"syscall"
	"time"
)

type options struct {
	trap      chan os.Signal
	log       *slog.Logger
	done      chan bool
	sigs      []os.Signal
	sdTimeout time.Duration
	ctx       context.Context

	connString string
	addr       string

	tlsKey  string
	tlsCert string
}

// String I put this here to not accidently leak the connString
func (o *options) String() {

}

// Option describes a setup option for the RPC server.
type Option interface {
	apply(*options)
}

type optFunc func(*options)

func (f optFunc) apply(o *options) { f(o) }

func WithLogger(log *slog.Logger) Option {
	return optFunc(func(o *options) { o.log = log })
}

func WithConnString(cString string) Option {
	return optFunc(func(o *options) { o.connString = cString })
}

func WithContext(ctx context.Context) Option {
	return optFunc(func(o *options) { o.ctx = ctx })
}

// WithTLS specifies the cert and key used to secure communication with the server.
func WithTLS(key, cert string) Option {
	return optFunc(func(o *options) {
		o.tlsKey = key
		o.tlsCert = cert
	})
}

func makeOptions(svrOptions []Option) *options {
	opts := &options{
		done:      make(chan bool),
		log:       defaultLogger(),
		sigs:      []os.Signal{syscall.SIGINT, syscall.SIGTERM},
		trap:      make(chan os.Signal),
		sdTimeout: 5 * time.Second,

		connString: "postgres://postgres:postgres@localhost:6432/my_database?sslmode=disable",
	}

	for _, opt := range svrOptions {
		opt.apply(opts)
	}

	return opts
}

func defaultLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{}))
}
