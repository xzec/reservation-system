package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func Run(ctx context.Context, w io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logger := createLogger(w)
	slog.SetDefault(logger)

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("unable to connect to the database: %v\n", err)
	}
	defer func() {
		slog.Info("closing database pool connection ...")
		pool.Close()
	}()

	server := NewServer(pool)
	httpServer := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: server,
	}

	go func() {
		slog.Info(fmt.Sprintf("listening and serving on %s", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		slog.Info(fmt.Sprintf("shutdown initiated, reason: %v", ctx.Err()))
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		slog.Info("gracefully shutting down the http server ...")
		if err = httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down the http server: %s\n", err)
		}
	}()
	wg.Wait()

	return nil
}
