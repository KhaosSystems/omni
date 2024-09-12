package models

import "github.com/google/uuid"

/*
* TaskType defines how a family of tasks should be handled.
 */
type TaskType struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
