package knowledge

import (
	"github.com/nikhil478/knowledgebuilder/internal/models"
)

func Build(ticketID string, facts []models.Fact) models.KnowledgePackage {

	systems := map[string]bool{}
	errors := map[string]bool{}

	var systemsList []string
	var errorsList []string

	for _, fact := range facts {

		switch fact.Type {

		case "service":

			if !systems[fact.Value] {
				systems[fact.Value] = true
				systemsList = append(
					systemsList,
					fact.Value,
				)
			}

		case "error":

			if !errors[fact.Value] {
				errors[fact.Value] = true
				errorsList = append(
					errorsList,
					fact.Value,
				)
			}
		}
	}

	return models.KnowledgePackage{
		TicketID: ticketID,
		Systems:  systemsList,
		Errors:   errorsList,
		Facts:    facts,
	}
}
