package logs

import (
	"strings"

	"github.com/nikhil478/knowledgebuilder/internal/models"
)

func DetectRootCause(facts []models.Fact) *models.Fact {

	for i := len(facts) - 1; i >= 0; i-- {

		f := facts[i]

		if f.Type != "error" {
			continue
		}

		if strings.Contains(f.Value, "Caused by") {
			return &f
		}

		if strings.Contains(f.Value, "Exception") {
			return &f
		}

		if strings.Contains(f.Value, "timeout") {
			return &f
		}
	}

	return nil
}
