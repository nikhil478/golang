package main

import (
	"github.com/google/uuid"
	"github.com/nikhil478/golang/projects/db/internal/common"
	"github.com/nikhil478/golang/projects/db/internal/engine"
	"github.com/nikhil478/golang/projects/db/internal/storage"
	"github.com/nikhil478/golang/projects/db/internal/wal"
	"go.uber.org/zap"
)

const (
	snapshotPath = "greet.json"
	walPath      = "greet.wal"
)

func main() {
	logger := common.NewLogger()

	requestID := uuid.NewString()

	requestLogger := logger.With(
		zap.String("request_id", requestID),
	)

	// ------------------------------------------------------------
	// FIRST STARTUP
	// ------------------------------------------------------------

	logger.Info("opening database - first startup")

	store := storage.NewInMemoryStorage(
		snapshotPath,
		requestLogger,
	)

	walStore, err := wal.NewFileWAL(
		walPath,
		requestLogger,
	)
	if err != nil {
		logger.Error(
			"failed to create WAL",
			zap.Error(err),
		)

		return
	}

	db, err := engine.OpenEngine(
		store,
		walStore,
		requestLogger,
	)
	if err != nil {
		logger.Error(
			"failed to open database engine",
			zap.Error(err),
		)

		return
	}

	logger.Info("database opened successfully")

	// ------------------------------------------------------------
	// WRITE
	// ------------------------------------------------------------

	logger.Info("setting value")

	if err := db.Set(
		"greet",
		[]byte("HELLO NIKHIL !"),
	); err != nil {
		logger.Error(
			"failed to set value",
			zap.String("key", "greet"),
			zap.Error(err),
		)

		return
	}

	logger.Info(
		"value set successfully",
		zap.String("key", "greet"),
	)

	// ------------------------------------------------------------
	// SIMULATE RESTART
	// ------------------------------------------------------------

	logger.Info("simulating database restart")

	// Create completely new storage and WAL instances
	// using the same snapshot and WAL files.

	store = storage.NewInMemoryStorage(
		snapshotPath,
		requestLogger,
	)

	walStore, err = wal.NewFileWAL(
		walPath,
		requestLogger,
	)
	if err != nil {
		logger.Error(
			"failed to reopen WAL",
			zap.Error(err),
		)

		return
	}

	db, err = engine.OpenEngine(
		store,
		walStore,
		requestLogger,
	)
	if err != nil {
		logger.Error(
			"failed to recover database",
			zap.Error(err),
		)

		return
	}

	logger.Info("database recovered successfully")

	// ------------------------------------------------------------
	// VERIFY RECOVERY
	// ------------------------------------------------------------

	val, err := db.Get("greet")
	if err != nil {
		logger.Error(
			"failed to get recovered value",
			zap.String("key", "greet"),
			zap.Error(err),
		)

		return
	}

	logger.Info(
		"recovered value successfully",
		zap.String("key", "greet"),
		zap.ByteString("value", val),
	)

	// ------------------------------------------------------------
	// DELETE
	// ------------------------------------------------------------

	logger.Info("deleting value")

	if err := db.Delete("greet"); err != nil {
		logger.Error(
			"failed to delete value",
			zap.String("key", "greet"),
			zap.Error(err),
		)

		return
	}

	logger.Info(
		"value deleted successfully",
		zap.String("key", "greet"),
	)

	// ------------------------------------------------------------
	// SIMULATE SECOND RESTART
	// ------------------------------------------------------------

	logger.Info("simulating second database restart")

	store = storage.NewInMemoryStorage(
		snapshotPath,
		requestLogger,
	)

	walStore, err = wal.NewFileWAL(
		walPath,
		requestLogger,
	)
	if err != nil {
		logger.Error(
			"failed to reopen WAL",
			zap.Error(err),
		)

		return
	}

	db, err = engine.OpenEngine(
		store,
		walStore,
		requestLogger,
	)
	if err != nil {
		logger.Error(
			"failed to recover database after delete",
			zap.Error(err),
		)

		return
	}

	// ------------------------------------------------------------
	// VERIFY DELETE RECOVERY
	// ------------------------------------------------------------

	_, err = db.Get("greet")

	if err == common.ErrKeyNotFound {
		logger.Info(
			"deletion successfully recovered from WAL",
			zap.String("key", "greet"),
		)

		return
	}

	if err != nil {
		logger.Error(
			"unexpected error while checking deleted key",
			zap.String("key", "greet"),
			zap.Error(err),
		)

		return
	}

	logger.Error(
		"key still exists after recovery",
		zap.String("key", "greet"),
	)
}