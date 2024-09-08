package omni

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type OmniService struct {
	db *sql.DB
}

type Project struct {
	UUID  uuid.UUID `json:"uuid"`
	Title string    `json:"title"`
}

type Task struct {
	/*
	* The unique identifier of the task.
	 */
	UUID        uuid.UUID `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`

	/*
	* The owning project.
	 */
	ProjectUUID uuid.UUID `json:"project_uuid"`

	/*
	* Assignees of the task.
	 */
	Assignees []uuid.UUID `json:"assignees"`

	/*
	* Reporter of the task.
	 */
	Reporters []uuid.UUID `json:"reporter"`
}

func NewOmniService(db *sql.DB) (*OmniService, error) {
	// Test the connection to the database.
	err := db.Ping()
	if err != nil {
		return nil, err
	}

	// Create the tasks table if it does not exist.
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		uuid UUID PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT
	)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create tasks table: %v", err)
	}

	// Create the projects table if it does not exist.
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS projects (
		uuid UUID PRIMARY KEY,
		title TEXT NOT NULL
	)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create projects table: %v", err)
	}

	return &OmniService{db: db}, nil
}

// GetTask returns a task from the database.
func (s *OmniService) GetTask(uuid uuid.UUID) (Task, error) {
	var task Task
	err := s.db.QueryRow("SELECT * FROM tasks WHERE uuid = $1", uuid).Scan(&task.UUID, &task.Title, &task.Description)
	if err != nil {
		return task, err
	}

	return task, nil
}

// CreateTask creates a new task in the database.
func (s *OmniService) CreateTask(title, description string) (uuid.UUID, error) {
	// Generate a new UUID for the task.
	uuid, err := uuid.NewRandom()
	if err != nil {
		return uuid, err
	}

	// Insert the task into the database.
	_, err = s.db.Exec("INSERT INTO tasks (uuid, title, description) VALUES ($1, $2, $3)", uuid, title, description)
	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

// CreateProject creates a new project in the database.
func (s *OmniService) CreateProject(title string) (uuid.UUID, error) {
	// Generate a new UUID for the project.
	uuid, err := uuid.NewRandom()
	if err != nil {
		return uuid, err
	}

	// Insert the project into the database.
	_, err = s.db.Exec("INSERT INTO projects (uuid, title) VALUES ($1, $2)", uuid, title)
	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

type ListProjectsOptions struct {
	Limit  int
	Offset int
}

type ListProjectsResponse struct {
	Projects []Project `json:"projects"`
	Count    int       `json:"count"`
	Total    int       `json:"total"`
}

// ListProjects returns a list of projects from the database.
func (s *OmniService) ListProjects(options ListProjectsOptions) (ListProjectsResponse, error) {
	query := `
	SELECT *,
		COUNT(*) OVER() AS total
	FROM projects
`
	args := []interface{}{}
	argIdx := 1

	// Validate the limit and offset values.
	if options.Limit < 0 {
		return ListProjectsResponse{}, fmt.Errorf("invalid limit value: %d", options.Limit)
	}
	if options.Offset < 0 {
		return ListProjectsResponse{}, fmt.Errorf("invalid offset value: %d", options.Offset)
	}

	// Add the limit and offset to the query.
	if options.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, options.Limit)
		argIdx++
	}

	if options.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, options.Offset)
	}

	// Print the query for debugging purposes.
	fmt.Println(query, args)

	// Query the database for all projects.
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return ListProjectsResponse{}, err
	}
	defer rows.Close()

	// Iterate over the rows and create a slice of projects.
	projects := []Project{}
	var total int
	for rows.Next() {
		var project Project
		err := rows.Scan(&project.UUID, &project.Title, &total)
		if err != nil {
			return ListProjectsResponse{}, err
		}
		projects = append(projects, project)
	}

	return ListProjectsResponse{
		Projects: projects,
		Count:    len(projects),
		Total:    total,
	}, nil
}

type ListTasksOptions struct {
	Limit  int
	Offset int
}

type ListTasksResponse struct {
	Tasks []Task `json:"tasks"`
	Count int    `json:"count"`
	Total int    `json:"total"`
}

// ListTasks returns a list of task from the database.
func (s *OmniService) ListTasks(options ListTasksOptions) (ListTasksResponse, error) {
	query := `
		SELECT *,
			COUNT(*) OVER() AS total
		FROM tasks
	`
	args := []interface{}{}
	argIdx := 1

	// Validate the limit and offset values.
	if options.Limit < 0 {
		return ListTasksResponse{}, fmt.Errorf("invalid limit value: %d", options.Limit)
	}
	if options.Offset < 0 {
		return ListTasksResponse{}, fmt.Errorf("invalid offset value: %d", options.Offset)
	}

	// Add the limit and offset to the query.
	if options.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, options.Limit)
		argIdx++
	}

	if options.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, options.Offset)
	}

	// Print the query for debugging purposes.
	fmt.Println(query, args)

	// Query the database for all tasks.
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return ListTasksResponse{}, err
	}
	defer rows.Close()

	// Iterate over the rows and create a slice of tasks.
	tasks := []Task{}
	var total int
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.UUID, &task.Title, &task.Description, &total)
		if err != nil {
			return ListTasksResponse{}, err
		}
		tasks = append(tasks, task)
	}

	return ListTasksResponse{
		Tasks: tasks,
		Count: len(tasks),
		Total: total,
	}, nil
}
