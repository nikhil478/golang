package images

import (
	"encoding/json"
	"strings"

	"github.com/nikhil478/knowledgebuilder/internal/models"
)

func Analyze(path string) ([]models.Fact, error) {

	resp, err := analyzeWithGemini(path)
	if err != nil {
		return nil, err
	}

	resp = cleanJSON(resp)

	var result Analysis

	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		return nil, err
	}

	return ToFacts(result, path), nil
}

func cleanJSON(input string) string {

	// remove markdown fences
	input = strings.TrimSpace(input)

	input = strings.ReplaceAll(input, "```json", "")
	input = strings.ReplaceAll(input, "```", "")

	return strings.TrimSpace(input)
}
