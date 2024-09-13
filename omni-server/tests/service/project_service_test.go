package integration

import (
	"context"
	"database/sql"
	"testing"

	"github.com/khaossystems/omni-server/internal/pkg/krest"
	"github.com/khaossystems/omni-server/internal/pkg/krest_orm"
	"github.com/khaossystems/omni-server/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

func TestCreateProject(t *testing.T) {
	// Create in-memory database.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Create the project service.
	projectRepository := krest_orm.NewGenericPostgresRepository[models.Project](db)
	projectService := krest_orm.NewGenericService(projectRepository)

	// Create a project.
	project := models.Project{
		Name: "Test Project",
	}
	_, err = projectService.Create(context.Background(), project)
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	// Make sure the project was created.
	projects, err := projectService.List(context.Background(), krest.CollectionQuery{})
	if err != nil {
		t.Fatalf("failed to list projects: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}
	if projects[0].Name != project.Name {
		t.Fatalf("expected project name %s, got %s", project.Name, projects[0].Name)
	}
}
