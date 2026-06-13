package config

import "os"

type Config struct {
	JiraURL   string
	JiraEmail string
	JiraToken string
	GeminiKey string
}

func Load() Config {

	return Config{
		JiraURL:   os.Getenv("JIRA_URL"),
		JiraEmail: os.Getenv("JIRA_EMAIL"),
		JiraToken: os.Getenv("JIRA_TOKEN"),
		GeminiKey: os.Getenv("GEMINI_API_KEY"),
	}
}
