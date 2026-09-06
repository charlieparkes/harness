package model

//go:generate go tool go-enum

// AgentActivityType is the kind of agent activity event.
// ENUM(Unspecified, Start, ToolCall, Read, Write, Error, Complete)
type AgentActivityType string
