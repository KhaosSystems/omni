package krest_orm_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/khaossystems/omni-server/internal/pkg/krest"
	"github.com/khaossystems/omni-server/internal/pkg/krest_orm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db      *sql.DB
	service *krest_orm.GenericService[TestType]
)

type TestType struct {
	UUID uuid.UUID `json:"uuid" krest_orm:"pk"`
	Name string    `json:"name"`
}

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		fmt.Printf("failed to open database: %v\n", err)
		os.Exit(1)
	}

	repository := krest_orm.NewGenericPostgresRepository[TestType](db)
	service = krest_orm.NewGenericService(repository)

	// Run the tests
	code := m.Run()
	// Cleanup
	db.Close()
	os.Exit(code)
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	testData := TestType{
		UUID: id,
		Name: "TestName",
	}

	createdData, err := service.Create(ctx, testData)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if createdData.UUID != testData.UUID || createdData.Name != testData.Name {
		t.Errorf("Create returned incorrect data: got %+v, want %+v", createdData, testData)
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	testData := TestType{
		UUID: id,
		Name: "TestName",
	}

	_, err := service.Create(ctx, testData)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	gotData, err := service.Get(ctx, id, krest.ResourceQuery{})
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if gotData.UUID != testData.UUID || gotData.Name != testData.Name {
		t.Errorf("Get returned incorrect data: got %+v, want %+v", gotData, testData)
	}
}

func TestList(t *testing.T) {
	ctx := context.Background()
	id1 := uuid.New()
	id2 := uuid.New()

	testData1 := TestType{
		UUID: id1,
		Name: "TestName1",
	}

	testData2 := TestType{
		UUID: id2,
		Name: "TestName2",
	}

	_, err := service.Create(ctx, testData1)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	_, err = service.Create(ctx, testData2)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	gotData, err := service.List(ctx, krest.CollectionQuery{})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(gotData) != 2 {
		t.Errorf("List returned incorrect number of items: got %d, want %d", len(gotData), 2)
	}

	// Check if the items are in the result
	var found1, found2 bool
	for _, item := range gotData {
		if item.UUID == testData1.UUID {
			found1 = true
		}
		if item.UUID == testData2.UUID {
			found2 = true
		}
	}

	if !found1 || !found2 {
		t.Errorf("List did not return all items: got %+v, want %+v and %+v", gotData, testData1, testData2)
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	testData := TestType{
		UUID: id,
		Name: "TestName",
	}

	_, err := service.Create(ctx, testData)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	updatedData := TestType{
		UUID: id,
		Name: "UpdatedName",
	}

	result, err := service.Update(ctx, id, updatedData)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if result.Name != updatedData.Name {
		t.Errorf("Update returned incorrect data: got %+v, want %+v", result, updatedData)
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	testData := TestType{
		UUID: id,
		Name: "TestName",
	}

	_, err := service.Create(ctx, testData)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = service.Delete(ctx, id)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = service.Get(ctx, id, krest.ResourceQuery{})
	if err == nil {
		t.Fatalf("Expected error when getting deleted item, got nil")
	}
}
