package images

import "github.com/nikhil478/knowledgebuilder/internal/models"

func ToFacts(a Analysis, source string) []models.Fact {

	var facts []models.Fact

	// Nodes → facts
	for _, n := range a.Nodes {
		facts = append(facts, models.Fact{
			Type:  "node",
			Value: n.ID,
		})
	}

	// Edges → relationship facts (CRITICAL PART)
	if a.Edges != nil {
		for _, e := range a.Edges {
			facts = append(facts, models.Fact{
				Type:       "dependency",
				Value:      e.From + " -> " + e.To,
				From:       e.From,
				To:         e.To,
				SourceType: "image",
			})
		}
	}

	return facts
}
