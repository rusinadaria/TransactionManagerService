package main

import (
	"net/http"
	// "fmt"
	"log/slog"
	"os"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
    "Transaction-manager/internal"
	"log"
)

func main() {
	logger := configLogger()

    p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Ошибка создания продюсера: %s", err)
	}
	defer p.Close()

    r := internal.ConfigRouter(p, logger)

    err = http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Error("failed start server")
		panic(err)
	}
}

func configLogger() *slog.Logger {
	var logger *slog.Logger

	f, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
        slog.Error("Unable to open a file for writing")
    }

	opts := &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }

	logger = slog.New(slog.NewJSONHandler(f, opts))
	return logger
}


