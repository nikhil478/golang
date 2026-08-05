package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nikhil478/golang/projects/db/projects/db/internal/common"
)

func NewInMemoryStorage(relativePath string) Storage {
	table := make(table)
	inMemoryStorage := inMemoryStorage{
		table:        table,
		relativePath: relativePath,
	}
	_, err := os.Stat(relativePath)
	if err == nil {
		table, err = inMemoryStorage.loadSnapshot()
		if err == nil {
			inMemoryStorage.table = table
		}
	}
	return &inMemoryStorage
}

type inMemoryStorage struct {
	table table

	relativePath string
}

func (db *inMemoryStorage) Set(key string, data []byte) error {
	db.table[key] = data
	db.saveSnapShot()
	return nil
}

func (db *inMemoryStorage) Get(key string) ([]byte, error) {
	fmt.Printf("db state before getting value from db : %v \n", db.table)
	val, ok := db.table[key]
	if !ok {
		return nil, common.ErrKeyNotFound
	}
	return val, nil
}

func (db *inMemoryStorage) Delete(key string) error {
	delete(db.table, key)
	return nil
}

func (db *inMemoryStorage) saveSnapShot() error {
	file, err := os.Create(db.relativePath)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(db.table)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		fmt.Printf("error while writing table to json file %v", data)
	}

	return nil
}

// . this function will be useful when crash happened
func (db *inMemoryStorage) loadSnapshot() (table, error) {
	var table table
	data, err := os.ReadFile(db.relativePath)
	if err != nil {
		return table, err
	}
	if err := json.Unmarshal(data, &table); err != nil {
		return table, err
	}
	return table, nil
}
