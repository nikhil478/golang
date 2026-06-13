package models

type Fact struct {
	ID string

	// What kind of fact this is:
	// error | service | database | queue | dependency | event
	Type string

	// Human readable value
	Value string

	// Graph support (CRITICAL for images)
	From string
	To   string

	// Confidence from LLM or parser
	Confidence float64

	// Source tracking
	SourceFile string
	SourceLine int

	// Optional timestamp (VERY useful for logs)
	Timestamp string

	// Where this came from
	// log | image | jira | llm
	SourceType string
}
