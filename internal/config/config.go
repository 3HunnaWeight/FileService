package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  string
	PostgresDSN string
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3Region    string
	MaxUploadMB int
}

func MustLoad() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
		S3Endpoint:  os.Getenv("S3_ENDPOINT"),
		S3AccessKey: os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey: os.Getenv("S3_SECRET_KEY"),
		S3Bucket:    os.Getenv("S3_BUCKET"),
		S3Region:    os.Getenv("S3_REGION"),
		MaxUploadMB: 32,
	}

	var missing []string
	if cfg.PostgresDSN == "" {
		missing = append(missing, "POSTGRES_DSN")
	}
	if cfg.S3Endpoint == "" {
		missing = append(missing, "S3_ENDPOINT")
	}
	if cfg.S3AccessKey == "" {
		missing = append(missing, "S3_ACCESS_KEY")
	}
	if cfg.S3SecretKey == "" {
		missing = append(missing, "S3_SECRET_KEY")
	}
	if cfg.S3Bucket == "" {
		missing = append(missing, "S3_BUCKET")
	}
	if cfg.S3Region == "" {
		missing = append(missing, "S3_REGION")
	}
	if len(missing) > 0 {
		panic(fmt.Sprintf("missing required env vars: %v", missing))
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
