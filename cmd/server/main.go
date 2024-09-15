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

	"github.com/buzzryan/zenbu/internal/config"
	"github.com/buzzryan/zenbu/internal/httputil"
	"github.com/buzzryan/zenbu/internal/logutil"
	"github.com/buzzryan/zenbu/internal/rdbutil"
	userctrl "github.com/buzzryan/zenbu/internal/user/controller"
	userinfra "github.com/buzzryan/zenbu/internal/user/infra"
)

func main() {
	logutil.InitDefaultLogger()

	cfg := config.LoadConfigFromEnv()
	rdb := rdbutil.MustConnectMySQL(cfg.MySQLConfig)
	slog.Info("mysql connected")

	mux := http.NewServeMux()

	userRepo := userinfra.NewUserRepo(rdb)
	tokenManager := userinfra.NewJWSTokenManager(cfg.JWSSigningKey)
	userctrl.Init(&userctrl.InitOpts{
		Mux:          mux,
		UserRepo:     userRepo,
		TokenManager: tokenManager,
	})

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
