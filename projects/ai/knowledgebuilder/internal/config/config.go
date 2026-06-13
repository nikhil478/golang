package config

import "os"

type Config struct {
	JiraURL   string
	JiraToken string
	GeminiKey string
}

func Load() Config {

	return Config{
		JiraURL:   os.Getenv("JIRA_URL"),
		JiraToken: os.Getenv("JIRA_TOKEN"),
		GeminiKey: os.Getenv("GEMINI_API_KEY"),
	}
}
