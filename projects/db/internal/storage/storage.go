package storage

import "github.com/nikhil478/golang/projects/db/projects/db/internal/common"

type table map[string][]byte

type Storage interface {
	Set(key string, data []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}

func NewInMemoryStorage() Storage {
	return &storage{
		table: make(table),
	}
}

type storage struct {
	table table
}

func (db *storage) Set(key string, data []byte) error {
	db.table[key] = data
	return nil
}

func (db *storage) Get(key string) ([]byte, error) {
	val, ok := db.table[key]
	if !ok {
		return nil, common.ErrKeyNotFound
	}
	return val, nil
}

func (db *storage) Delete(key string) error {
	delete(db.table, key)
	return nil
}
