package export

import (
	"encoding/json"
	"os"

	"github.com/nikhil478/knowledgebuilder/internal/models"
)

func SaveKnowledge(
	path string,
	k models.KnowledgePackage,
) error {

	f, err := os.Create(path)

	if err != nil {
		return err
	}

	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	return enc.Encode(k)
}
