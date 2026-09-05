package engine

import (
	"fmt"
	"time"

	"github.com/nikhil478/golang/projects/db/internal/common"
	"github.com/nikhil478/golang/projects/db/internal/storage"
	"github.com/nikhil478/golang/projects/db/internal/wal"
	"go.uber.org/zap"
)

type Engine interface {
	Set(key string, data []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}

type engineSVC struct {
	storage storage.Storage
	wal     wal.WAL
	logger  common.Logger
}

func OpenEngine(
	storage storage.Storage,
	wal wal.WAL,
	logger common.Logger,
) (Engine, error) {
	svc := &engineSVC{
		storage: storage,
		wal:     wal,
		logger:  logger,
	}

	if err := svc.recover(); err != nil {
		logger.Error(
			"failed to recover database from WAL",
			zap.Error(err),
		)

		return nil, err
	}

	logger.Info("database recovery completed")

	return svc, nil
}

func (svc *engineSVC) Set(key string, data []byte) error {
	svc.logger.Debug(
		"setting key in engine",
		zap.String("key", key),
	)

	entry := wal.Entry{
		Timestamp: time.Now(),
		Operation: wal.OPERATION_SET,
		Key:       key,
		Data:      data,
	}

	if err := svc.wal.Append(entry); err != nil {
		svc.logger.Error(
			"failed to append set operation to WAL",
			zap.String("key", key),
			zap.Error(err),
		)

		return err
	}

	if err := svc.storage.Set(key, data); err != nil {
		svc.logger.Error(
			"failed to set key in storage",
			zap.String("key", key),
			zap.Error(err),
		)

		return err
	}

	svc.logger.Debug(
		"key set successfully in engine",
		zap.String("key", key),
	)

	return nil
}

func (svc *engineSVC) Get(key string) ([]byte, error) {
	svc.logger.Debug(
		"getting key from engine",
		zap.String("key", key),
	)

	data, err := svc.storage.Get(key)
	if err != nil {
		if err == common.ErrKeyNotFound {
			svc.logger.Debug(
				"key not found in engine",
				zap.String("key", key),
			)

			return nil, err
		}

		svc.logger.Error(
			"failed to get key from engine",
			zap.String("key", key),
			zap.Error(err),
		)

		return nil, err
	}

	svc.logger.Debug(
		"key retrieved successfully from engine",
		zap.String("key", key),
	)

	return data, nil
}

func (svc *engineSVC) Delete(key string) error {
	svc.logger.Debug(
		"deleting key from engine",
		zap.String("key", key),
	)

	entry := wal.Entry{
		Timestamp: time.Now(),
		Operation: wal.OPERATION_DELETE,
		Key:       key,
	}

	if err := svc.wal.Append(entry); err != nil {
		svc.logger.Error(
			"failed to append delete operation to WAL",
			zap.String("key", key),
			zap.Error(err),
		)

		return err
	}

	if err := svc.storage.Delete(key); err != nil {
		svc.logger.Error(
			"failed to delete key in storage",
			zap.String("key", key),
			zap.Error(err),
		)

		return err
	}

	svc.logger.Debug(
		"key deleted successfully in engine",
		zap.String("key", key),
	)

	return nil
}

func (svc *engineSVC) recover() error {
	entries, err := svc.wal.Replay()
	if err != nil {
		return err
	}

	svc.logger.Info(
		"replaying WAL entries",
		zap.Int("entries", len(entries)),
	)

	for _, entry := range entries {
		switch entry.Operation {
		case wal.OPERATION_SET:
			if err := svc.storage.Set(entry.Key, entry.Data); err != nil {
				return err
			}

		case wal.OPERATION_DELETE:
			if err := svc.storage.Delete(entry.Key); err != nil {
				return err
			}

		default:
			return fmt.Errorf(
				"unknown WAL operation: %s",
				entry.Operation,
			)
		}
	}

	return nil
}
