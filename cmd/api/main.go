package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/3HunnaWeight/file-service/internal/config"
	apphttp "github.com/3HunnaWeight/file-service/internal/http"
	"github.com/3HunnaWeight/file-service/internal/repository/postgres"
	"github.com/3HunnaWeight/file-service/internal/service"
	s3storage "github.com/3HunnaWeight/file-service/internal/storage/s3"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg := config.MustLoad()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pool, err := postgres.New(ctx, cfg.PostgresDSN)
	if err != nil {
		slog.Error("postgres connect failed", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := postgres.Ping(ctx, pool); err != nil {
		slog.Error("postgres ping failed", "err", err)
		os.Exit(1)
	}
	slog.Info("postgres connected")

	storage, err := s3storage.New(
		ctx,
		cfg.S3Endpoint,
		cfg.S3AccessKey,
		cfg.S3SecretKey,
		cfg.S3Region,
		cfg.S3Bucket,
	)
	if err != nil {
		slog.Error("s3 init failed", "err", err)
		os.Exit(1)
	}

	if err := storage.EnsureBucket(ctx); err != nil {
		slog.Error("s3 bucket check failed", "err", err)
		os.Exit(1)
	}
	slog.Info("s3 storage ready", "bucket", cfg.S3Bucket)

	repo := postgres.NewFileRepository(pool)
	fileService := service.NewFileService(repo, storage, cfg.S3Bucket)
	fileHandler := apphttp.NewFileHandler(fileService)
	router := apphttp.NewRouter(fileHandler)

	addr := ":" + cfg.ServerPort
	srv := &apphttp.Server{}
	go srv.Start(addr, router)

	<-ctx.Done()
	slog.Info("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	srv.Shutdown(shutdownCtx)
}
