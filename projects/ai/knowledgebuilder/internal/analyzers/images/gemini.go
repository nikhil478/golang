package images

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

func analyzeWithGemini(path string) (string, error) {

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY not configured")
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return "", err
	}

	imageBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{Text: ArchitecturePrompt},
				{
					InlineData: &genai.Blob{
						MIMEType: detectMimeType(path),
						Data:     imageBytes,
					},
				},
			},
		},
	}

	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		contents,
		nil,
	)

	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates returned")
	}

	var out string
	for _, p := range resp.Candidates[0].Content.Parts {
		if p.Text != "" {
			out += p.Text
		}
	}

	return out, nil
}