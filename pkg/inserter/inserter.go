package inserter

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"
	// postgres side effects

	_ "github.com/lib/pq"
)

type Inserter struct {
	db  *sql.DB
	Log *slog.Logger

	Stutter time.Duration
	ConnStr string

	wChan chan int
	rChan chan int
}

func (in *Inserter) Go(ctx context.Context, wCount, rCount int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	for i := 0; i < rCount; i++ {
		g.Go(func() error { return in.read(ctx) })
	}

	for i := 0; i < wCount; i++ {
		g.Go(func() error { return in.write(ctx) })
	}

	g.Go(func() error { return count(ctx, in.wChan, "write", in.Log) })
	g.Go(func() error { return count(ctx, in.rChan, "read", in.Log) })

	err := g.Wait()
	return err
}

func count(ctx context.Context, c <-chan int, name string, log *slog.Logger) error {
	for {
		select {
		case i := <-c:
			if i%250 == 0 {
				log.Info(name, "count", i)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (in *Inserter) read(ctx context.Context) error {
	i := 0
	for {
		_, err := in.db.Exec("SELECT id, name, age  FROM example_table LIMIT 1;")
		if err != nil {
			in.Log.Error("failed to read", "error", err)
			return err
		}
		select {
		case <-ctx.Done():
			return nil
		case in.rChan <- i:
			i = i + 1
		}
	}
}

func (in *Inserter) write(ctx context.Context) error {
	i := 0
	for {
		_, err := in.db.Exec("INSERT INTO example_table(name, age) VALUES($1, $2);", "foo", i)
		if err != nil {
			in.Log.Error("failed to write", "error", err)
		}
		select {
		case <-ctx.Done():
			return nil
		case in.wChan <- i:
			i = i + 1
		}
	}
}

func (in *Inserter) RebuildPool() error {
	in.db.Close()

	return in.SpawnPool()
}

func (in *Inserter) SpawnPool() error {
	db, err := sql.Open("postgres", in.ConnStr)
	in.Log.Info(in.ConnStr)

	if err != nil {
		in.Log.Error("failed to create a postgres connection", "error", err)
		return err
	}
	in.wChan = make(chan int, 5)
	in.rChan = make(chan int, 5)

	in.db = db
	return nil
}

func (in *Inserter) Close() {
	in.db.Close()
}
