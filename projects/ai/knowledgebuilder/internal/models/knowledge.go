package models

type TimelineEvent struct {
	Timestamp string `json:"timestamp"`
	Event     string `json:"event"`
}

type KnowledgePackage struct {
	TicketID string `json:"ticket_id"`

	Summary string `json:"summary"`

	Systems      []string        `json:"systems"`
	Errors       []string        `json:"errors"`
	Dependencies []string        `json:"dependencies"`
	Timeline     []TimelineEvent `json:"timeline"`

	Facts []Fact `json:"facts"`
}
