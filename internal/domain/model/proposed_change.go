package model

import "github.com/charlieparkes/harness/internal/lineage"

type ProposedChange struct {
	Description string                          `json:"description"`
	Reason      string                          `json:"reason"`
	Files       lineage.RevisionedField[[]File] `json:"files"`
}

type ProposedTest struct {
	ProposedChange

	TestCases lineage.RevisionedField[[]string] `json:"test_cases"`
}
