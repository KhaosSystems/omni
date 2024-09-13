package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/khaossystems/omni-server/internal/pkg/krest"
	"github.com/khaossystems/omni-server/internal/pkg/krest_orm"
	"github.com/khaossystems/omni-server/pkg/models"
	"github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func connectToDatabase() (*sql.DB, error) {
	// Get database connection information from environment variables.
	dbHost := os.Getenv("OMNI_DB_HOST")
	dbPort := os.Getenv("OMNI_DB_PORT")
	dbUser := os.Getenv("OMNI_DB_USER")
	dbPass := os.Getenv("OMNI_DB_PASSWORD")
	dbName := os.Getenv("OMNI_DB_DATABASE")

	// Validate required environment variables.
	if dbHost == "" || dbPort == "" || dbUser == "" || dbPass == "" || dbName == "" {
		log.Fatal("Missing required environment variables")
	}

	// Connect to the Postgres server.
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL server: %v", err)
	}
	defer db.Close()

	// Attempt to create the database if it does not exist.
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		// Check if the error is due to the database already existing.
		if pgError, ok := err.(*pq.Error); ok && pgError.Code == "42P04" {
			// 42P04 is the code for "database already exists"
			log.Println("Database already exists, skipping creation.")
		} else {
			log.Fatalf("Failed to create database: %v", err)
		}
	} else {
		log.Println("Database created successfully.")
	}

	// Now connect to the specific database.
	db.Close() // Close the previous connection
	connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Test the connection to the specific database.
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	return db, nil
}

func connectToSQLiteDatabase() (*sql.DB, error) {
	// Get database connection information from environment variables.
	dbPath := os.Getenv("OMNI_DATA_PATH")

	// Validate required environment variables.
	if dbPath == "" {
		log.Fatal("Missing required environment variables")
	}

	// Connect to the SQLite database.
	db, err := sql.Open("sqlite3", dbPath+"/omni.db")
	if err != nil {
		log.Fatalf("Failed to connect to SQLite database: %v", err)
	}

	// Test the connection to the SQLite database.
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the SQLite database: %v", err)
	}

	return db, nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func createRouter(db *sql.DB) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	userRepository := krest_orm.NewGenericPostgresRepository[models.User](db)
	userService := krest_orm.NewGenericService(userRepository)
	userHandler := krest.NewHandler(userService)

	taskRepository := krest_orm.NewGenericPostgresRepository[models.Task](db)
	taskService := krest_orm.NewGenericService(taskRepository)
	taskHandler := krest.NewHandler(taskService)

	projectRepository := krest_orm.NewGenericPostgresRepository[models.Project](db)
	projectService := krest_orm.NewGenericService(projectRepository)
	projectHandler := krest.NewHandler(projectService)

	router.Route("/v1", func(v2 chi.Router) {
		// Users
		v2.Get("/users/{uuid}", userHandler.Get)
		v2.Get("/users", userHandler.List)
		v2.Post("/users", userHandler.Create)
		v2.Patch("/users/{uuid}", userHandler.Update)
		v2.Delete("/users/{uuid}", userHandler.Delete)

		// Tasks
		v2.Get("/tasks/{uuid}", taskHandler.Get)
		v2.Get("/tasks", taskHandler.List)
		v2.Post("/tasks", taskHandler.Create)
		v2.Patch("/tasks/{uuid}", taskHandler.Update)
		v2.Delete("/tasks/{uuid}", taskHandler.Delete)

		// Projects
		v2.Get("/projects/{uuid}", projectHandler.Get)
		v2.Get("/projects", projectHandler.List)
		v2.Post("/projects", projectHandler.Create)
		v2.Patch("/projects/{uuid}", projectHandler.Update)
		v2.Delete("/projects/{uuid}", projectHandler.Delete)
	})

	return router
}

func main() {
	// Load .env file.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	} else {
		log.Println("Successfully loaded .env file")
	}

	// Connect to the database.
	//db, err := connectToDatabase()
	db, err := connectToSQLiteDatabase()
	if err != nil {
		log.Panicf("Error connecting to database: %v", err)
	} else {
		log.Println("Successfully connected to database")
	}
	defer db.Close()

	// Create api router.
	router := createRouter(db)

	// Get port from environment variable.
	port := os.Getenv("OMNI_API_PORT")
	if port == "" {
		port = "30090"
	}

	// Start the server.
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Panicf("Error starting server: %v", err)
	} else {
		log.Println("Successfully started on port " + port)
	}
}
