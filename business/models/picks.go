package models

import "github.com/google/uuid"

type Pick struct {
	ID      uuid.UUID
	Word    string
	Letters []string
}
