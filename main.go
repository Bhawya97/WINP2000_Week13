package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

type TimeResponse struct {
	CurrentTime string `json:"current_time"`
	Error       string `json:"error,omitempty"`
}

func main() {
	// Set up logging to a file
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Load .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database credentials
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Connect to MySQL database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	log.Println("Database connected successfully")

	// Start the HTTP server
	http.HandleFunc("/current-time", currentTimeHandler)
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func currentTimeHandler(w http.ResponseWriter, r *http.Request) {
	// Get Toronto timezone
	loc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		writeErrorResponse(w, "Failed to load Toronto timezone", err)
		log.Printf("Failed to load Toronto timezone: %v", err)
		return
	}

	// Get current time
	currentTime := time.Now().In(loc)

	// Log the time to the database
	err = logTimeToDatabase(currentTime)
	if err != nil {
		writeErrorResponse(w, "Failed to log time to database", err)
		log.Printf("Failed to log time to database: %v", err)
		return
	}

	// Return the current time as JSON
	response := TimeResponse{
		CurrentTime: currentTime.Format("2006-01-02 15:04:05"),
	}
	writeJSONResponse(w, response)
	log.Printf("Current time logged: %s", currentTime.Format("2006-01-02 15:04:05"))
}

func logTimeToDatabase(timestamp time.Time) error {
	_, err := db.Exec("INSERT INTO time_log (timestamp) VALUES (?)", timestamp)
	if err != nil {
		log.Printf("Error logging time to database: %v", err)
	}
	return err
}

func writeErrorResponse(w http.ResponseWriter, message string, err error) {
	log.Printf("%s: %v", message, err)
	w.WriteHeader(http.StatusInternalServerError)
	writeJSONResponse(w, TimeResponse{Error: message})
}

func writeJSONResponse(w http.ResponseWriter, response TimeResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
