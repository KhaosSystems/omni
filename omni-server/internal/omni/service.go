package omni

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type OmniService struct {
	db *sql.DB
}

type Task struct {
	UUID        uuid.UUID `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
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
