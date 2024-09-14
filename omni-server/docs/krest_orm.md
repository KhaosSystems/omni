Krest ORM is a **_highly opinionated_** ORM-like framework for Go, designed to work with the Krest RESTful specification, and accompanying library Krest.

This library is not a replacement to learning SQL, but a tool that can help experienced developers with the burden that is baisc boilerplate SQL operations and architecture. It is designed to be highly customizable, and the automatic schema generation -by design- wil not cover all use casese.

The library is currently being used in production at Khaos Group companies, with a Postgres backend. But should be easily extendable to other SQL databases- and possibly NoSQL databases.

## Features
 - Automatic schema generation from structs.
 - Postgres support.
 - Generic service pattern and repository pattern implementation for basic CRUD operations. (extendable)
 - Helpers for easy unit tests.

## Usage
Customizable though the 'krest_orm' tag.
 - pk: Primary key.
 - fk: Foreign key.
 - ignore: Ignore the field in automatic schema generation.
 - custom: Custom SQL for the field- if the automatic schema generation is not cutting it (which is wont- this is not a replacement to learning SQL.).
```go
/*
* Task represents a task in the system.
 */
type Task struct {
	UUID        uuid.UUID `json:"uuid" krest_orm:"pk"`
	Summary     string    `json:"summary" krest:"expandable"`
	Description string    `json:"description" krest:"expandable"`
	Type        *TaskType `json:"type" krest:"expandable" krest_orm:"custom:'task_type UUID REFERENCES task_types(uuid)'"`
	Status      *Status   `json:"status" krest:"expandable" krest_orm:"ignore"`
	Project     *Project  `json:"project" krest:"expandable" krest_orm:"ignore"`
}
```
