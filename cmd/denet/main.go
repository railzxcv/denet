package main

import (
	"context"
	"denet-app/internal/config"
	"denet-app/internal/storage"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"denet-app/internal/routes"
	"errors"

	"denet-app/internal/token"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// check env
	var env string
	flag.StringVar(&env, "env", "local", "Path to config file")
	flag.Parse()

	// config
	cfg := config.Load(env)

	// logger
	log := setupLogger(env)

	// token manager
	TokenMn, err := token.NewTokenManagerRSA(cfg.PublicKeyPath)
	if err != nil {
		log.Error("error created a new token manager", "error", err)
		os.Exit(1)
	}

	// storage
	storage, err := storage.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Error("failed to create storage", "error", err)
		return
	}

	// routes
	router := routes.NewRouter(log, storage, TokenMn)
	
	// server & gracefull shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server", "error", err)
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", "error", err)
		return
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
