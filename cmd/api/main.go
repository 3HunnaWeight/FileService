package main

import (
	"context"
	"log"

	"github.com/3HunnaWeight/file-service/internal/config"
	apphttp "github.com/3HunnaWeight/file-service/internal/http"
	"github.com/3HunnaWeight/file-service/internal/repository/postgres"
	"github.com/3HunnaWeight/file-service/internal/service"
	s3storage "github.com/3HunnaWeight/file-service/internal/storage/s3"
)

func main() {
	cfg := config.MustLoad()

	// DB
	db, err := postgres.New(cfg.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewFileRepository(db)

	// S3 storage
	storage, err := s3storage.New(
		context.Background(),
		cfg.S3Endpoint,
		cfg.S3AccessKey,
		cfg.S3SecretKey,
		cfg.S3Region,
		cfg.S3Bucket,
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = storage.EnsureBucket(context.Background())

	// service layer
	fileService := service.NewFileService(repo, storage, cfg.S3Bucket)

	// http layer
	fileHandler := apphttp.NewFileHandler(fileService)
	router := apphttp.NewRouter(fileHandler)
	log.Println("server started :" + cfg.ServerPort)

	err = apphttp.StartServer(":"+cfg.ServerPort, router)
	if err != nil {
		log.Fatal(err)
	}
}
