package agentmodel

import (
	"github.com/charlieparkes/go-transform"
	"github.com/charlieparkes/harness/internal/domain/model"
)

// Plan is the expected json document structure of an agent-defined task plan.
type Plan struct {
	// Short sentence describing step
	Title string `json:"title"`

	// Description of goals, reasoning, and purpose
	// May be as short as one sentence, or as long as multiple paragraphs.
	Description string `json:"description"`

	// A list of steps which should be executed to complete the plan
	Steps []Step `json:"steps"`

	// Definition of done for the overall plan
	DefinitionOfDone []string `json:"definition_of_done"`

	// List of identified risks for the overall plan
	Risks []string `json:"risks"`
}

func NewPlan(p model.Plan) Plan {
	title, _ := p.Title.Value()
	description, _ := p.Description.Value()
	steps, _ := p.Steps.Value()
	definitionOfDone, _ := p.DefinitionOfDone.Value()
	risks, _ := p.Risks.Value()
	return Plan{
		Title:            title,
		Description:      description,
		Steps:            transform.Slice(steps, NewStep),
		DefinitionOfDone: definitionOfDone,
		Risks:            risks,
	}
}
