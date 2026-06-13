package images

import (
	"mime"
	"path/filepath"
)

func detectMimeType(path string) string {

	ext := filepath.Ext(path)

	m := mime.TypeByExtension(ext)
	if m == "" {
		return "image/png"
	}

	return m
}