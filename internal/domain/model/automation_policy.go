package model

// AutomationPolicy is the approval mode for each workflow phase.
type AutomationPolicy struct {
	Worktree AutomationPolicyType `json:"worktree"`
	Plan     AutomationPolicyType `json:"plan"`
	Execute  AutomationPolicyType `json:"execute"`
	Review   AutomationPolicyType `json:"review"`
	Apply    AutomationPolicyType `json:"apply"`
}
