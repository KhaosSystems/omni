package models

import "github.com/google/uuid"

/*
* Task represents a task in the system.
 */
type Task struct {
	UUID        uuid.UUID `json:"uuid" krest_orm:"pk"`
	Summary     string    `json:"summary" krest:"expandable"`
	Description string    `json:"description" krest:"expandable"`
	Type        *TaskType `json:"type" krest:"expandable" krest_orm:"ignore,fk"`
	Status      *Status   `json:"status" krest:"expandable" krest_orm:"ignore"`
	Project     *Project  `json:"project" krest:"expandable" krest_orm:"ignore"`
}
