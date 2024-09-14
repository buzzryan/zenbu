package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/buzzryan/zenbu/internal/httputil"
	"github.com/buzzryan/zenbu/internal/logutil"
)

func main() {
	logutil.InitDefaultLogger()

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: httputil.WithGlobalMiddlewares(mux),
	}

	// gracefully shut down
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan

		slog.Info("caught shutting down signal", slog.Any("signal", sig.String()))
		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownRelease()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("HTTP shutdown error: %v", err)
		}
	}()

	slog.Info("http: server start")
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server error: %v", err)
	}
	slog.Info("http: server down")
}
