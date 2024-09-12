package models

import "github.com/google/uuid"

/*
* Task represents a task in the system.
 */
type Task struct {
	UUID        uuid.UUID `json:"uuid" krest_orm:"pk"`
	Summary     string    `json:"summary,omitempty" krest:"expandable"`
	Description string    `json:"description,omitempty" krest:"expandable"`
	Type        *TaskType `json:"type,omitempty" krest:"expandable" krest_orm:"ignore"`
	Status      *Status   `json:"status,omitempty" krest:"expandable" krest_orm:"ignore"`
	Project     *Project  `json:"project,omitempty" krest:"expandable" krest_orm:"ignore"`
}
