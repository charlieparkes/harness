package model

// AutomationPolicy is the approval mode for each workflow phase.
type AutomationPolicy struct {
	Worktree AutomationPolicyType
	Plan     AutomationPolicyType
	Execute  AutomationPolicyType
	Review   AutomationPolicyType
	Apply    AutomationPolicyType
}
