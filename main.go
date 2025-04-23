package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/hiteshchouhan29/fraudblocker/config"
	"github.com/hiteshchouhan29/fraudblocker/handlers"
	_ "github.com/lib/pq"
)

var logger *log.Logger

func initLogger() *log.Logger {
	return log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	// Initialize logger
	logger := initLogger()
	logger.Println("Starting Phone Number Blocking System...")

	// Load configuration
	config, err := config.LoadConfig("config.yaml")
	if err != nil {
		logger.Println("Failed to load configuration file, using default settings")

	}

	// Connect to database
	db, err := sql.Open("postgres", config.GetDSN())
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		logger.Fatalf("Failed to ping database: %v", err)
	}
	logger.Printf("Connected to database: %s", config.Database.DBName)

	// Run migrations
	err = handlers.RunMigrations(db, "db/migrations/init.sql")
	if err != nil {
		logger.Fatalf("Migration failed: %v", err)
	}
	logger.Println("Database migrations applied successfully")

	// Parse templates
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		logger.Fatalf("Failed to parse templates: %v", err)
	}
	logger.Println("Templates parsed successfully")

	// Set up routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(w, r, db, templates, logger)
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddNumberHandler(w, r, db, logger)
	})

	http.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
		handlers.RemoveNumberHandler(w, r, db, logger)
	})

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		handlers.CheckNumberHandler(w, r, db, logger)
	})

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the server
	serverAddr := fmt.Sprintf(":%d", config.Server.Port)
	logger.Printf("Server starting on http://localhost%s", serverAddr)
	logger.Fatal(http.ListenAndServe(serverAddr, nil))
}
