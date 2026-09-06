package model

//go:generate go tool go-enum

// TargetPurpose is whether a target is read-only or writable.
// ENUM(Unspecified, Read, ReadWrite)
type TargetPurpose string
