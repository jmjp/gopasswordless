package main

import (
	"database/sql"
	delivery "hyperzoop/internal/adapters/delivery/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	postgres, _ := sql.Open("postgres", os.Getenv("database_url"))
	if err := postgres.Ping(); err != nil {
		zap.L().Error("failed to connect pg database")
		panic(err)
	}
	defer postgres.Close()
	delivery.NewHTTPServer(3000, postgres).Start()
}
