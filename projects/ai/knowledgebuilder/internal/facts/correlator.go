package facts

import (
	"strings"

	"github.com/nikhil478/knowledgebuilder/internal/models"
)

func Correlate(input []models.Fact) []models.Fact {

	seen := map[string]bool{}

	var output []models.Fact

	for _, fact := range input {

		key := strings.ToLower(
			fact.Type + ":" + fact.Value,
		)

		if seen[key] {
			continue
		}

		seen[key] = true

		output = append(output, fact)
	}

	return output
}
