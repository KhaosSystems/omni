package db

import "github.com/khaossystems/omni-server/pkg/models"

func (q *Queries) GetTask(uuid string) (models.Task, error) {
	var task models.Task
	err := q.db.QueryRow("SELECT * FROM tasks WHERE uuid = $1", uuid).Scan(&task.UUID, &task.Title, &task.Description, &task.Project.UUID)
	if err != nil {
		return task, err
	}

	return task, nil
}
