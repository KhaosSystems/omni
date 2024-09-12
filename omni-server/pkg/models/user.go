package models

import "github.com/google/uuid"

/*
* User represents a user of the system.
 */
type User struct {
	UUID     uuid.UUID `json:"uuid" krest_orm:"pk"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}
