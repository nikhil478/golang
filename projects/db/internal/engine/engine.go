package engine

import "github.com/nikhil478/golang/projects/db/projects/db/internal/storage"

type Engine interface {
	Set(key string, data []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}

type engineSVC struct {
	storage storage.Storage
}

func NewEngine() Engine {
	return &engineSVC{
		storage: storage.NewInMemoryStorage(),
	}
}

func (svc *engineSVC) Set(key string, data []byte) error {
	return svc.storage.Set(key, data)
}
func (svc *engineSVC) Get(key string) ([]byte, error) {
	return svc.storage.Get(key)
}
func (svc *engineSVC) Delete(key string) error {
	return svc.storage.Delete(key)
}
