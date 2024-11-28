package main

import (
	"bytes"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/glebarez/sqlite" // SQLite driver without CGO dependency
	"github.com/stretchr/testify/assert"
)

var testDb *sql.DB

// Initialize the test database (SQLite in-memory)
func setupTestDB() error {
	var err error
	// Use SQLite in-memory database for testing (without CGO)
	testDb, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		return err
	}

	// Create the time_log table for testing
	_, err = testDb.Exec("CREATE TABLE time_log (id INTEGER PRIMARY KEY, timestamp DATETIME)")
	return err
}

func tearDownTestDB() {
	if testDb != nil {
		testDb.Close()
	}
}

func TestCurrentTimeHandler(t *testing.T) {
	// Redirect logs to a buffer to capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// Setup the in-memory SQLite database
	err := setupTestDB()
	if err != nil {
		t.Fatalf("Error setting up test database: %v", err)
	}
	defer tearDownTestDB()

	// Set the global db variable to use the testDb
	db = testDb

	// Simulate an HTTP request to the /current-time endpoint
	req, err := http.NewRequest("GET", "/current-time", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(currentTimeHandler)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check that the status code is OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check if the log contains the correct log message for the current time
	expectedLogMessage := "Current time logged:"
	if !bytes.Contains(buf.Bytes(), []byte(expectedLogMessage)) {
		t.Errorf("Expected log to contain %q, but got %q", expectedLogMessage, buf.String())
	}
}

func TestLogError(t *testing.T) {
	// Redirect logs to a buffer to capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// Setup the in-memory SQLite database
	err := setupTestDB()
	if err != nil {
		t.Fatalf("Error setting up test database: %v", err)
	}
	defer tearDownTestDB()

	// Set the global db variable to use the testDb
	db = testDb

	// Force an error in logTimeToDatabase by making the query invalid
	// This is done by simulating a failure (wrong SQL query)
	db.Exec("DROP TABLE IF EXISTS time_log") // Force an error by dropping the table

	// Call the logTimeToDatabase to simulate an error
	err = logTimeToDatabase(time.Now()) // This should trigger the error
	if err == nil {
		t.Fatalf("Expected error logging to database, but got nil")
	}

	// Check if the log contains the expected error message
	expectedErrorMessage := "Error logging time to database"
	if !bytes.Contains(buf.Bytes(), []byte(expectedErrorMessage)) {
		t.Errorf("Expected log to contain %q, but got %q", expectedErrorMessage, buf.String())
	}
}
