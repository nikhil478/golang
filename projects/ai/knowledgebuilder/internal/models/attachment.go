package models

type Attachment struct {
	ID          string
	TicketID    string
	FileName    string
	ContentType string
	LocalPath   string
	Size        int64
}
