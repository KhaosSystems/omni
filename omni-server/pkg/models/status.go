package models

import "github.com/google/uuid"

/*
* Status represents the status of a task.
 */
type Status struct {
	UUID        uuid.UUID `json:"uuid" krest_orm:"pk"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
