package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/nikhil478/knowledgebuilder/internal/analyzers/images"
	"github.com/nikhil478/knowledgebuilder/internal/analyzers/logs"
	"github.com/nikhil478/knowledgebuilder/internal/classifier"
	"github.com/nikhil478/knowledgebuilder/internal/client/jira"
	"github.com/nikhil478/knowledgebuilder/internal/config"
	"github.com/nikhil478/knowledgebuilder/internal/export"
	"github.com/nikhil478/knowledgebuilder/internal/facts"
	"github.com/nikhil478/knowledgebuilder/internal/knowledge"
	"github.com/nikhil478/knowledgebuilder/internal/models"
	"github.com/nikhil478/knowledgebuilder/internal/storage"
)

func main() {

	_ = godotenv.Load()

	if len(os.Args) < 2 {
		log.Fatal("usage: knowledgebuilder <ticket-id>")
	}

	ticketID := os.Args[1]

	cfg := config.Load()

	db, err := storage.NewSQLite("jira.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	workspaceDir := filepath.Join("workspace", ticketID)
	attachmentsDir := filepath.Join(workspaceDir, "attachments")
	knowledgeDir := filepath.Join(workspaceDir, "knowledge")

	if err := os.MkdirAll(attachmentsDir, 0755); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(knowledgeDir, 0755); err != nil {
		log.Fatal(err)
	}

	jiraClient := jira.New(
		cfg.JiraURL,
		cfg.JiraEmail,
		cfg.JiraToken,
	)

	fmt.Printf("\nLoading JIRA issue %s\n", ticketID)

	issue, err := jiraClient.GetIssue(ticketID)
	if err != nil {
		log.Fatalf("failed to load jira issue: %v", err)
	}

	fmt.Printf("Loaded issue: %s\n", issue.Key)

	// ======================================================
	// ALWAYS SYNC ATTACHMENTS (no folder-length shortcut)
	// ======================================================

	fmt.Println("Syncing Jira attachments...")

	err = jiraClient.DownloadAttachments(issue, attachmentsDir)
	if err != nil {
		log.Fatalf("attachment sync failed: %v", err)
	}

	files, err := os.ReadDir(attachmentsDir)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nFound %d attachment(s)\n\n", len(files))

	// ======================================================
	// ANALYSIS PIPELINE
	// ======================================================

	var allFacts []models.Fact

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		path := filepath.Join(attachmentsDir, file.Name())

		attachmentType, err := classifier.Classify(path)
		if err != nil {
			log.Printf("classification failed: %s : %v", path, err)
			continue
		}

		fmt.Printf("Processing %s (%s)\n", file.Name(), attachmentType)

		switch attachmentType {

		case classifier.LogAttachment:

			logFacts, err := logs.ParseLog(path)
			if err != nil {
				log.Printf("log parsing failed: %v", err)
				continue
			}

			allFacts = append(allFacts, logFacts...)

		case classifier.ImageAttachment:

			imageFacts, err := images.Analyze(path)
			if err != nil {
				log.Printf("image analysis failed: %v", err)
				continue
			}

			allFacts = append(allFacts, imageFacts...)

		default:

			fmt.Printf("Skipping unsupported file: %s\n", file.Name())
		}
	}

	fmt.Printf("\nRaw facts extracted: %d\n", len(allFacts))

	correlatedFacts := facts.Correlate(allFacts)

	fmt.Printf("Facts after correlation: %d\n", len(correlatedFacts))

	rootCause := logs.DetectRootCause(correlatedFacts)

	if rootCause != nil {
		fmt.Printf("Root Cause Candidate: %s\n", rootCause.Value)
	}

	kp := knowledge.Build(ticketID, correlatedFacts)

	outputFile := filepath.Join(knowledgeDir, "knowledge.json")

	err = export.SaveKnowledge(outputFile, kp)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("======================================")
	fmt.Printf("Ticket: %s\n", ticketID)
	fmt.Printf("Systems: %d\n", len(kp.Systems))
	fmt.Printf("Errors: %d\n", len(kp.Errors))
	fmt.Printf("Facts: %d\n", len(kp.Facts))

	if rootCause != nil {
		fmt.Printf("Root Cause: %s\n", rootCause.Value)
	}

	fmt.Printf("Knowledge Package: %s\n", outputFile)
	fmt.Println("======================================")
}
