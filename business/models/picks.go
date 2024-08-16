package models

import "github.com/google/uuid"

type Pick struct {
	ID      uuid.UUID
	Letters []string
}

type Score struct {
	Word  string
	Score int
	Exist bool
}
