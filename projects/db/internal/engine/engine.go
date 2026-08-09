package engine

import (
	"github.com/nikhil478/golang/projects/db/internal/common"
	"github.com/nikhil478/golang/projects/db/internal/storage"
	"go.uber.org/zap"
)

type Engine interface {
	Set(key string, data []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}

type engineSVC struct {
	storage storage.Storage
	logger  common.Logger
}

func NewEngine(storage storage.Storage, logger common.Logger) Engine {
	return &engineSVC{
		storage: storage,
		logger:  logger,
	}
}

func (svc *engineSVC) Set(key string, data []byte) error {
	svc.logger.Debug("setting key in engine", zap.String("key", key))

	if err := svc.storage.Set(key, data); err != nil {
		svc.logger.Error("failed to set key in engine", zap.String("key", key), zap.Error(err))
		return err
	}

	svc.logger.Debug("key set successfully in engine", zap.String("key", key))
	return nil
}

func (svc *engineSVC) Get(key string) ([]byte, error) {
	svc.logger.Debug("getting key from engine", zap.String("key", key))
	data, err := svc.storage.Get(key)
	if err != nil {
		svc.logger.Error("failed to get key from engine", zap.String("key", key), zap.Error(err))
		return nil, err
	}
	svc.logger.Debug("key retrieved successfully from engine", zap.String("key", key))
	return data, nil
}

func (svc *engineSVC) Delete(key string) error {
	svc.logger.Debug("deleting key from engine", zap.String("key", key))
	if err := svc.storage.Delete(key); err != nil {
		svc.logger.Error("failed to delete key from engine", zap.String("key", key), zap.Error(err))
		return err
	}
	svc.logger.Debug("key deleted successfully from engine", zap.String("key", key))
	return nil
}
