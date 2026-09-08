package agentmodel

import (
	"time"

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

func (p Plan) Apply(prev model.Plan, now time.Time) model.Plan {
	rev := prev.Revision + 1

	changed := prev.Title.Set(rev, p.Title)
	changed = prev.Description.Set(rev, p.Description) || changed
	changed = prev.Risks.Set(rev, p.Risks) || changed
	changed = prev.DefinitionOfDone.Set(rev, p.DefinitionOfDone) || changed

	prevSteps, _ := prev.Steps.Value()
	prevByID := make(map[string]model.Step, len(prevSteps))
	for _, s := range prevSteps {
		prevByID[s.ID] = s
	}
	steps := transform.Slice(p.Steps, func(s Step) model.Step {
		return s.Apply(prevByID[s.ID], rev)
	})
	changed = prev.Steps.Set(rev, steps) || changed

	if changed {
		prev.Revision = rev
		prev.UpdatedAt = &now
	}
	return prev
}
