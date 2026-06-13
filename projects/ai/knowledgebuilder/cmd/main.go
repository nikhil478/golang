package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nikhil478/knowledgebuilder/internal/classifier"
	"github.com/nikhil478/knowledgebuilder/internal/knowledge"
	"github.com/nikhil478/knowledgebuilder/internal/models"
	"github.com/nikhil478/knowledgebuilder/internal/storage"
)

func main() {

	godotenv.Load()

	if len(os.Args) < 2 {
		log.Fatal("usage: jira-analyzer <ticket-key>")
	}

	ticketID := os.Args[1]

	db, err := storage.NewSQLite("jira.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	files, err := os.ReadDir("./data/" + ticketID)
	if err != nil {
		log.Fatal(err)
	}

	var allFacts []models.Fact

	for _, f := range files {
		path := "./data/" + ticketID + "/" + f.Name()

		t, err := classifier.Classify(path)
		if err != nil {
			continue
		}

		fmt.Printf("processing %s (%s)\n", path, t)
	}

	kp := knowledge.Build(ticketID, allFacts)

	fmt.Printf("knowledge package built %+v\n", kp)
}
