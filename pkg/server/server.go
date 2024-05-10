package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/highstead/pgcat-experiments/pkg/inserter"
)

var shuttingDown bool

// ListenAndServe starts the gRPC server that serves API requests.
func ListenAndServe(addr string, opts ...Option) (chan<- os.Signal, <-chan bool) {
	svrOpts := makeOptions(opts)

	// Start HTTP server
	go func(opts *options) {
		opts.log.Info("Starting HTTP server...")

		http.HandleFunc("/ping", ping)
		http.HandleFunc("/healthz", healthz)

		go func() {
			if err := http.ListenAndServe(opts.addr, nil); err != nil {
				opts.log.Error("HTTP server error %v\n", err)
				return
			}
		}()

		ctx, cancel := context.WithCancel(opts.ctx)
		ins := &inserter.Inserter{
			ConnStr: opts.connString,
			Stutter: time.Second * 1,
			Log:     opts.log,
		}
		opts.log.Info("Spawning pool")
		err := ins.SpawnPool()
		if err != nil {
			close(opts.done)
			opts.log.Error("Unable to establish connection pool", "error", err)
			return
		}
		opts.log.Info("doing stuff")
		doStuff(ctx, ins)

		for {
			select {
			case sig := <-opts.trap:
				if sig == syscall.SIGTERM {
					opts.log.Info("Recieved SIGHTERM, Shutting down")
					gracefulShutdown(opts.sdTimeout, cancel)
					close(opts.done)
				}

				if sig == syscall.SIGHUP {
					opts.log.Info("Recieved SIGHUP, reloading config")

					// Kill go routines by cancelling context
					cancel()

					// Spawn a new context and rebuild the pool and spawn new routines
					ctx, cancel = context.WithCancel(opts.ctx)
					ins.RebuildPool()
					doStuff(ctx, ins)
				}
			}
		}
	}(svrOpts)

	return svrOpts.trap, svrOpts.done
}

func ping(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second) // Simulate some work
	fmt.Fprintf(w, "pong\n")
}

func healthz(w http.ResponseWriter, r *http.Request) {
	if !shuttingDown {
		w.WriteHeader(http.StatusOK)
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

func gracefulShutdown(dur time.Duration, cancel context.CancelFunc) {
	shuttingDown = true
	cancel()

	<-time.After(dur)
}

func rebuildPool() {

}

func doStuff(ctx context.Context, ins *inserter.Inserter) {
	ins.Go(ctx, 2, 2)
}
