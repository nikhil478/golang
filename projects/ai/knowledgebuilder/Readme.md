# Idea of knowldge builder is to create knowledge base from attachments of jira ticket

## Initial Flow

                                        JIRA Ticket
                                            |
                                        Download Attachments
                                            |
                                        Attachment Registry
                                            |
                                        File Classifier
                                            |
                                        Log Analyzer
                                            |
                                        Image Analyzer
                                            |
                                        Fact Extractor
                                            |
                                        Knowledge Builder
                                            |
                                        Knowledge Package (JSON)


Project Structure:

cmd/
  jira-analyzer/

internal/
  jira/
    client.go

  attachments/
    downloader.go
    registry.go

  classifier/
    classifier.go

  analyzers/
    logs/
      analyzer.go
      parser.go

    images/
      analyzer.go

  facts/
    extractor.go

  knowledge/
    builder.go
    graph.go

  storage/
    sqlite.go

  models/
    attachment.go
    fact.go
    knowledge.go

docker run --env-file .env kbbuilder