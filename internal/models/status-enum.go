package models

// StatusEnum defines the possible status values
type StatusEnum string

const (
	Active   StatusEnum = "A"
	Inactive StatusEnum = "I"
	Deleted  StatusEnum = "D"
)
