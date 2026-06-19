package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	ServerPort  string
	PostgresDSN string
	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
	S3Region    string
}

func MustLoad() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		ServerPort:  os.Getenv("SERVER_PORT"),
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
		S3Endpoint:  os.Getenv("S3_ENDPOINT"),
		S3AccessKey: os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey: os.Getenv("S3_SECRET_KEY"),
		S3Bucket:    os.Getenv("S3_BUCKET"),
		S3Region:    os.Getenv("S3_REGION"),
	}

	if cfg.ServerPort == "" {
		panic("SERVER_PORT IS MISSING")
	}
	if cfg.PostgresDSN == "" {
		panic("POSTGRES_DSN IS MISSING")
	}

	return cfg
}
