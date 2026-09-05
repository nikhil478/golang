package storage

type table map[string][]byte

type Storage interface {
	Set(key string, data []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	SaveSnapshot() error
	LoadSnapshot() (table, error)
}
