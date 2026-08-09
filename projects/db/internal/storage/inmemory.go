package storage

import (
	"encoding/json"
	"os"

	"github.com/nikhil478/golang/projects/db/internal/common"
	"go.uber.org/zap"
)

func NewInMemoryStorage(relativePath string, logger common.Logger) Storage {
	logger.Info(
		"creating in-memory storage",
		zap.String("path", relativePath),
	)

	table := make(table)

	inMemoryStorage := &inMemoryStorage{
		table:        table,
		relativePath: relativePath,
		logger:       logger,
	}

	_, err := os.Stat(relativePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Debug(
				"snapshot file does not exist, starting with empty table",
				zap.String("path", relativePath),
			)
		} else {
			logger.Warn(
				"failed to check snapshot file",
				zap.String("path", relativePath),
				zap.Error(err),
			)
		}
	}

	table, err = inMemoryStorage.loadSnapshot()
	if err != nil {
		logger.Warn(
			"failed to load snapshot, starting with current table",
			zap.String("path", relativePath),
			zap.Error(err),
		)
	} else {
		inMemoryStorage.table = table

		logger.Info(
			"snapshot loaded successfully",
			zap.String("path", relativePath),
			zap.Int("entries", len(table)),
		)
	}

	return inMemoryStorage
}

type inMemoryStorage struct {
	table        table
	relativePath string
	logger       common.Logger
}

func (db *inMemoryStorage) Set(key string, data []byte) error {
	db.logger.Debug(
		"setting key in storage",
		zap.String("key", key),
		zap.Int("data_size", len(data)),
	)

	db.table[key] = data

	if err := db.saveSnapshot(); err != nil {
		db.logger.Error(
			"failed to save snapshot after setting key",
			zap.String("key", key),
			zap.Error(err),
		)

		return err
	}

	db.logger.Debug(
		"key set successfully in storage",
		zap.String("key", key),
	)

	return nil
}

func (db *inMemoryStorage) Get(key string) ([]byte, error) {
	db.logger.Debug(
		"getting key from storage",
		zap.String("key", key),
	)

	val, ok := db.table[key]
	if !ok {
		db.logger.Debug(
			"key not found in storage",
			zap.String("key", key),
		)

		return nil, common.ErrKeyNotFound
	}

	db.logger.Debug(
		"key retrieved successfully from storage",
		zap.String("key", key),
	)

	return val, nil
}

func (db *inMemoryStorage) Delete(key string) error {
	db.logger.Debug(
		"deleting key from storage",
		zap.String("key", key),
	)

	delete(db.table, key)

	if err := db.saveSnapshot(); err != nil {
		db.logger.Error(
			"failed to save snapshot after deleting key",
			zap.String("key", key),
			zap.Error(err),
		)

		return err
	}

	db.logger.Debug(
		"key deleted successfully from storage",
		zap.String("key", key),
	)

	return nil
}

func (db *inMemoryStorage) saveSnapshot() error {
	db.logger.Debug(
		"saving storage snapshot",
		zap.String("path", db.relativePath),
		zap.Int("entries", len(db.table)),
	)

	file, err := os.Create(db.relativePath)
	if err != nil {
		db.logger.Error(
			"failed to create snapshot file",
			zap.String("path", db.relativePath),
			zap.Error(err),
		)

		return err
	}
	defer file.Close()

	data, err := json.Marshal(db.table)
	if err != nil {
		db.logger.Error(
			"failed to marshal storage table",
			zap.Error(err),
		)

		return err
	}

	if _, err := file.Write(data); err != nil {
		db.logger.Error(
			"failed to write snapshot",
			zap.String("path", db.relativePath),
			zap.Error(err),
		)

		return err
	}

	db.logger.Debug(
		"storage snapshot saved successfully",
		zap.String("path", db.relativePath),
		zap.Int("size", len(data)),
	)

	return nil
}

// This function is useful for recovering the database state after a crash.
func (db *inMemoryStorage) loadSnapshot() (table, error) {
	db.logger.Debug(
		"loading storage snapshot",
		zap.String("path", db.relativePath),
	)

	var table table

	data, err := os.ReadFile(db.relativePath)
	if err != nil {
		db.logger.Error(
			"failed to read snapshot",
			zap.String("path", db.relativePath),
			zap.Error(err),
		)

		return table, err
	}

	if err := json.Unmarshal(data, &table); err != nil {
		db.logger.Error(
			"failed to unmarshal snapshot",
			zap.String("path", db.relativePath),
			zap.Error(err),
		)

		return table, err
	}

	db.logger.Debug(
		"storage snapshot loaded",
		zap.String("path", db.relativePath),
		zap.Int("entries", len(table)),
	)

	return table, nil
}