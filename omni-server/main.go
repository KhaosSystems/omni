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
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/khaossystems/omni-server/api"
	"github.com/khaossystems/omni-server/internal/omni"
	"github.com/lib/pq"
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

func createRouter(service *omni.OmniService) *chi.Mux {
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

	router.Route("/v1", func(v1 chi.Router) {
		// Get task.
		v1.Get("/tasks/{uuid}", func(w http.ResponseWriter, r *http.Request) {
			// Get the UUID from the URL parameters.
			uuidStr := chi.URLParam(r, "uuid")

			// Parse the UUID from the URL parameters.
			uuid, err := uuid.Parse(uuidStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Get the resource query parameters from the request.
			_, err = api.ParseResourceQueryParams(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Get the task from the service.
			task, err := service.GetTask(uuid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			api.WriteResourceResponse(w, http.StatusOK, task, nil)
		})
		// List tasks.
		v1.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {
			// Parse the meta query parameters from the request.
			metaQueryParams, err := api.ParseMetaQueryParams(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Parse the collection query parameters from the request.
			collectionQueryParams, err := api.ParseCollectionQueryParams(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Get the list of tasks from the service.
			res, err := service.ListTasks(omni.ListTasksOptions{
				Limit:  collectionQueryParams.Limit,
				Offset: collectionQueryParams.Offset,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			api.WriteCollectionResponse(w, http.StatusOK, res.Tasks, res.Count, res.Total, collectionQueryParams, metaQueryParams)
		})
		// Create task.
		v1.Post("/tasks", func(w http.ResponseWriter, r *http.Request) {
			var task omni.Task

			// Check the Content-Type to determine if the request is JSON or form submission
			contentType := r.Header.Get("Content-Type")

			if contentType == "application/json" {
				// Handle JSON request
				err := json.NewDecoder(r.Body).Decode(&task)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else if contentType == "application/x-www-form-urlencoded" {
				// Handle form submission
				err := r.ParseForm()
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				task.Title = r.FormValue("title")
				task.Description = r.FormValue("description")
			} else {
				http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
				return
			}

			// Create the task using the service
			uuid, err := service.CreateTask(task.Title, task.Description)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Respond with the UUID of the created task, inside a JSON object
			object := map[string]interface{}{
				"@links": map[string]string{
					"self": fmt.Sprintf("/v1/tasks/%s", uuid),
				},
				"uuid": uuid,
			}

			respondWithJSON(w, http.StatusCreated, object)
		})
		// List projects.
		v1.Get("/projects", func(w http.ResponseWriter, r *http.Request) {
			// Parse the meta query parameters from the request.
			metaQueryParams, err := api.ParseMetaQueryParams(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Parse the collection query parameters from the request.
			collectionQueryParams, err := api.ParseCollectionQueryParams(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Get the list of projects from the service.
			res, err := service.ListProjects(omni.ListProjectsOptions{
				Limit:  collectionQueryParams.Limit,
				Offset: collectionQueryParams.Offset,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			api.WriteCollectionResponse(w, http.StatusOK, res.Projects, res.Count, res.Total, collectionQueryParams, metaQueryParams)
		})
		// Create project.
		v1.Post("/projects", func(w http.ResponseWriter, r *http.Request) {
			var project omni.Project

			// Check the Content-Type to determine if the request is JSON or form submission
			contentType := r.Header.Get("Content-Type")

			if contentType == "application/json" {
				// Handle JSON request
				err := json.NewDecoder(r.Body).Decode(&project)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else if contentType == "application/x-www-form-urlencoded" {
				// Handle form submission
				err := r.ParseForm()
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				project.Title = r.FormValue("title")
			} else {
				http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
				return
			}

			// Create the project using the service
			uuid, err := service.CreateProject(project.Title)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Respond with the UUID of the created project, inside a JSON object
			object := map[string]interface{}{
				"@links": map[string]string{
					"self": fmt.Sprintf("/v1/projects/%s", uuid),
				},
				"uuid": uuid,
			}

			respondWithJSON(w, http.StatusCreated, object)
		})
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
	db, err := connectToDatabase()
	if err != nil {
		log.Panicf("Error connecting to database: %v", err)
	} else {
		log.Println("Successfully connected to database")
	}
	defer db.Close()

	// Create the OmniService.
	omniService, err := omni.NewOmniService(db)
	if err != nil {
		log.Panicf("Error creating OmniService: %v", err)
	} else {
		log.Println("Successfully created service")
	}

	// Create api router.
	router := createRouter(omniService)

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
