package main

import (
	"github.com/google/uuid"
	"github.com/nikhil478/golang/projects/db/internal/common"
	"github.com/nikhil478/golang/projects/db/internal/engine"
	"github.com/nikhil478/golang/projects/db/internal/storage"
	"go.uber.org/zap"
)

func main() {
	logger := common.NewLogger()

	requestID := uuid.NewString()

	requestLogger := logger.With(
		zap.String("request_id", requestID),
	)

	store := storage.NewInMemoryStorage(
		"greet.json",
		requestLogger,
	)

	engine := engine.NewEngine(
		store,
		requestLogger,
	)

	logger.Info("starting database operations")

	if err := engine.Set("greet", []byte("HELLO NIKHIL !")); err != nil {
		logger.Error(
			"failed to set value in engine",
			zap.String("key", "greet"),
			zap.Error(err),
		)

		return
	}

	logger.Info(
		"value set successfully",
		zap.String("key", "greet"),
	)

	val, err := engine.Get("greet")
	if err != nil {
		logger.Error(
			"failed to get value from engine",
			zap.String("key", "greet"),
			zap.Error(err),
		)

		return
	}

	logger.Info(
		"value retrieved successfully",
		zap.String("key", "greet"),
		zap.ByteString("value", val),
	)

	if err := engine.Delete("greet"); err != nil {
		logger.Error(
			"failed to delete value from engine",
			zap.String("key", "greet"),
			zap.Error(err),
		)

		return
	}

	logger.Info(
		"value deleted successfully",
		zap.String("key", "greet"),
	)
}
