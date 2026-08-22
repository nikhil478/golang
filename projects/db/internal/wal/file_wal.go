package wal

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/nikhil478/golang/projects/db/internal/common"
	"go.uber.org/zap"
)

type fileWAL struct {
	filePath string
	file     *os.File
	logger   common.Logger
}

func NewFileWAL(filePath string, logger common.Logger) (WAL, error) {
	logger.Info(
		"creating file WAL",
		zap.String("path", filePath),
	)

	file, err := os.OpenFile(
		filePath,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		logger.Error(
			"failed to open WAL file",
			zap.String("path", filePath),
			zap.Error(err),
		)

		return nil, err
	}

	return &fileWAL{
		filePath: filePath,
		file:     file,
		logger:   logger,
	}, nil
}

func (w *fileWAL) Append(entry Entry) error {

	encoded, err := json.Marshal(entry)
	if err != nil {
		w.logger.Error(
			"failed to marshal WAL entry",
			zap.Error(err),
		)

		return err
	}

	encoded = append(encoded, '\n')

	if _, err := w.file.Write(encoded); err != nil {
		w.logger.Error(
			"failed to write WAL entry",
			zap.String("key", entry.Key),
			zap.Error(err),
		)

		return err
	}

	if err := w.file.Sync(); err != nil {
		w.logger.Error(
			"failed to sync WAL",
			zap.String("key", entry.Key),
			zap.Error(err),
		)

		return err
	}

	w.logger.Debug(
		"WAL entry appended",
		zap.String("operation", string(entry.Operation)),
		zap.String("key", entry.Key),
	)

	return nil
}

func (w *fileWAL) Replay() ([]Entry, error) {
	w.logger.Debug(
		"replaying WAL",
		zap.String("path", w.filePath),
	)

	file, err := os.Open(w.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			w.logger.Debug(
				"WAL file does not exist, nothing to replay",
				zap.String("path", w.filePath),
			)

			return []Entry{}, nil
		}

		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	entries := make([]Entry, 0)

	for scanner.Scan() {
		var entry Entry

		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			w.logger.Error(
				"failed to decode WAL entry",
				zap.Error(err),
			)

			return nil, err
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		w.logger.Error(
			"failed while reading WAL",
			zap.Error(err),
		)

		return nil, err
	}

	w.logger.Info(
		"WAL replay completed",
		zap.String("path", w.filePath),
		zap.Int("entries", len(entries)),
	)

	return entries, nil
}
