package agentmodel

import (
	"github.com/charlieparkes/go-transform"
	"github.com/charlieparkes/harness/internal/domain/model"
)

// Step is an agent-defined step in a list of steps necessary to complete a task.
//
// Definition of a step should include answers to, "what", "where", "why", and "how".
//
// A step should contain enough information, independent of the overall plan, for a
// person or agent to execute it without additional context.
type Step struct {
	// Agent-sourced identifier for this question
	// For example, S1, S2, S3, etc
	ID string `json:"id"`

	// Short sentence describing step
	Title string `json:"title"`

	// Description of goals, reasoning, and purpose
	// May be as short as one sentence, or as long as multiple paragraphs.
	Summary string `json:"summary"`

	// List of proposed changes or additions to any files: business logic, api definitions, documentation, etc.
	ProposedChanges []ProposedChange `json:"proposed_changes"`

	// List of proposed changes or additions to any tests or test files.
	// A step with proposed changes MUST propose test updates or additions, or it will be rejected.
	ProposedTests []ProposedTest `json:"proposed_tests"`

	// List of verification steps to take, including commands to run
	Verifications []string `json:"verifications"`

	// List of identified risks for this particular step
	Risks []string `json:"risks"`

	// List of step IDs which must be completed before this step may run
	Dependencies []string `json:"dependencies"`
}

// NewStep returns the agent step document for step.
// dependencies may include edges for other steps; only edges for step.ID are kept.
func NewStep(step model.Step) Step {
	return Step{
		ID:              step.ID,
		Title:           step.Title,
		Summary:         step.Summary,
		ProposedChanges: transform.Slice(step.ProposedChanges, NewProposedChange),
		ProposedTests:   transform.Slice(step.ProposedTests, NewProposedTest),
		Verifications:   step.Verifications,
		Risks:           step.Risks,
		Dependencies:    step.Dependencies,
	}
}

func (s Step) model() model.Step {
	return model.Step{
		ID:              s.ID,
		Title:           s.Title,
		Summary:         s.Summary,
		ProposedChanges: transform.Slice(s.ProposedChanges, ProposedChange.model),
		ProposedTests:   transform.Slice(s.ProposedTests, ProposedTest.model),
		Verifications:   s.Verifications,
		Risks:           s.Risks,
		Dependencies:    s.Dependencies,
	}
}
