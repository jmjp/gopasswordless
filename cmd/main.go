package main

import (
	"context"
	"database/sql"
	delivery "hyperzoop/internal/adapters/delivery/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	godotenv.Load()
	logger := zap.Must(zap.NewProduction())
	if os.Getenv("env") == "dev" {
		file, err := os.OpenFile(os.Getenv("log_file"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			zap.L().Error("failed to open log file", zap.Error(err))
			// handle error
		}
		defer file.Close()
		logger = zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(file),
			zapcore.InfoLevel,
		))
	}
	zap.ReplaceGlobals(logger)
	postgres, _ := sql.Open("postgres", getEnvDB())
	if err := postgres.Ping(); err != nil {
		zap.L().Error("failed to connect pg database")
		panic(err)
	}
	opt, err := redis.ParseURL(os.Getenv("redis_url"))
	if err != nil {
		panic(err)
	}
	redis := redis.NewClient(opt)
	if err := redis.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	defer postgres.Close()
	defer redis.Close()
	delivery.NewHTTPServer(3000, postgres, redis).Start()
}

// !TODO: move to .env
func getEnvDB() string {
	switch os.Getenv("env") {
	case "prod":
		return os.Getenv("prod_db")
	case "stagging":
		return os.Getenv("stagging_db")
	default:
		return os.Getenv("dev_db")
	}
}
