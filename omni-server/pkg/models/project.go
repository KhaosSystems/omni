package models

import "github.com/google/uuid"

/*
* Project represents a project in the system.
 */
type Project struct {
	UUID uuid.UUID `json:"uuid" krest_orm:"pk"`
	Name string    `json:"name"`
	//Tasks []*Task   `json:"tasks" krest:"expandable" krest_orm:"fk:Project"`

	/*
	* The project key is a short, human-readable identifier for the project.
	* The key is NOT garanteed to be unique.
	 */
	Key string `json:"key"`
}
