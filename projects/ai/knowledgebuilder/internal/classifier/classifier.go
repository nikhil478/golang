package classifier

import (
	"path/filepath"
	"strings"
)

type AttachmentType string

const (
	LogAttachment   AttachmentType = "log"
	ImageAttachment AttachmentType = "image"
	Unknown         AttachmentType = "unknown"
)

func Classify(path string) (AttachmentType, error) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".log", ".txt":
		return LogAttachment, nil

	case ".png", ".jpg", ".jpeg":
		return ImageAttachment, nil

	default:
		return Unknown, nil
	}
}
