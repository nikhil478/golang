package vector

import (
	"context"

	"github.com/nikhil478/golang/projects/tiny-memory/internal/chunk"
)

type VectorStore interface {
	Upsert(ctx context.Context, chunks []chunk.Chunk) error
	Search(ctx context.Context, vector []float32, limit int) ([]chunk.Chunk, error)
}
