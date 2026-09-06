package agentmodel

import (
	"github.com/charlieparkes/harness/internal/domain/model"
)

type ProposedChange struct {
	// One to three sentence description of the proposed change
	Description string `json:"description"`

	// One to three sentence reason for the proposed change
	Reason string `json:"reason"`

	// List of files affected by this change
	Files []model.File `json:"files"`
}

func NewProposedChange(c model.ProposedChange) ProposedChange {
	files, _ := c.Files.Value()
	return ProposedChange{
		Description: c.Description,
		Reason:      c.Reason,
		Files:       files,
	}
}

type ProposedTest struct {
	ProposedChange

	TestCases []string `json:"test_cases"`
}

func NewProposedTest(c model.ProposedTest) ProposedTest {
	testCases, _ := c.TestCases.Value()
	return ProposedTest{
		ProposedChange: NewProposedChange(c.ProposedChange),
		TestCases:      testCases,
	}
}
