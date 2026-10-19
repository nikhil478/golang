package vector

import (
	"context"
	"fmt"

	"github.com/TIANLI0/chromem-go"
	"github.com/nikhil478/golang/projects/tiny-memory/internal/chunk"
)

const collectionName = "tiny-memory"

type ChromemStore struct {
	db         *chromem.DB
	collection *chromem.Collection
}

func NewChromemStore(path string) (*ChromemStore, error) {
	db, err := chromem.NewPersistentDB(path, false)
	if err != nil {
		return nil, fmt.Errorf("create chromem db: %w", err)
	}

	collection, err := db.GetOrCreateCollection(
		collectionName,
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create collection: %w", err)
	}

	return &ChromemStore{
		db:         db,
		collection: collection,
	}, nil
}

func (s *ChromemStore) Upsert(
	ctx context.Context,
	chunks []chunk.Chunk,
) error {
	if len(chunks) == 0 {
		return nil
	}

	ids := make([]string, 0, len(chunks))
	embeddings := make([][]float32, 0, len(chunks))
	metadatas := make([]map[string]string, 0, len(chunks))
	contents := make([]string, 0, len(chunks))

	for _, c := range chunks {
		if len(c.Vector) == 0 {
			return fmt.Errorf(
				"chunk %q has empty embedding",
				c.ID,
			)
		}

		ids = append(ids, c.ID)

		embeddings = append(
			embeddings,
			c.Vector,
		)

		metadatas = append(
			metadatas,
			map[string]string{
				"file_name": c.FileName,
			},
		)

		contents = append(
			contents,
			c.Text,
		)
	}

	err := s.collection.Add(
		ctx,
		ids,
		embeddings,
		metadatas,
		contents,
	)
	if err != nil {
		return fmt.Errorf("add chunks: %w", err)
	}

	return nil
}

func (s *ChromemStore) Search(
	ctx context.Context,
	vector []float32,
	limit int,
) ([]chunk.Chunk, error) {
	if len(vector) == 0 {
		return nil, fmt.Errorf("search vector is empty")
	}

	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than zero")
	}

	results, err := s.collection.QueryEmbedding(
		ctx,
		vector,
		limit,
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("search vectors: %w", err)
	}

	chunks := make([]chunk.Chunk, 0, len(results))

	for _, result := range results {
		chunks = append(chunks, chunk.Chunk{
			ID:       result.ID,
			Text:     result.Content,
			FileName: result.Metadata["file_name"],
			Vector:   result.Embedding,
		})
	}

	return chunks, nil
}
