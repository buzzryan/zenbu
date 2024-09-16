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

	awscfg "github.com/aws/aws-sdk-go-v2/config"

	"github.com/buzzryan/zenbu/internal/config"
	"github.com/buzzryan/zenbu/internal/httputil"
	"github.com/buzzryan/zenbu/internal/logutil"
	"github.com/buzzryan/zenbu/internal/nosqlutil"
	"github.com/buzzryan/zenbu/internal/storageutil"
	userctrl "github.com/buzzryan/zenbu/internal/user/controller"
	userinfra "github.com/buzzryan/zenbu/internal/user/infra"
)

func main() {
	logutil.InitDefaultLogger()

	cfg := config.LoadConfigFromEnv()

	awsCfg, err := awscfg.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Panicf("failed to load AWS config: %v", err)
	}

	ddb := nosqlutil.ConnectDDB(awsCfg, cfg.DynamoConfig)
	slog.Info("dynamoDB connected")

	storage := storageutil.NewS3Storage(awsCfg, cfg.S3Config)

	mux := http.NewServeMux()

	userRepo := userinfra.NewDynamoUserRepo(ddb, cfg.TableName)
	tokenManager := userinfra.NewJWSTokenManager(cfg.JWSSigningKey)
	userctrl.Init(&userctrl.InitOpts{
		Mux:          mux,
		UserRepo:     userRepo,
		TokenManager: tokenManager,
		Storage:      storage,
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
