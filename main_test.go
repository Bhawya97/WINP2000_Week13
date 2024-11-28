package main

import (
	"bytes"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var testDb *sql.DB

// Initialize the test database (SQLite in-memory for testing)
func setupTestDB() error {
	var err error
	// Use SQLite in-memory database for testing (no CGO required)
	testDb, err = sql.Open("mysql", "root:password@tcp(localhost:3306)/testdb")
	if err != nil {
		return err
	}

	// Create the time_log table for testing
	_, err = testDb.Exec("CREATE TABLE IF NOT EXISTS time_log (id INT PRIMARY KEY AUTO_INCREMENT, timestamp DATETIME)")
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

func TestAllTimesHandler(t *testing.T) {
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

	// Insert some test times into the database
	times := []string{
		"2024-11-28 10:00:00",
		"2024-11-28 11:00:00",
	}
	for _, timeStr := range times {
		_, err := testDb.Exec("INSERT INTO time_log (timestamp) VALUES (?)", timeStr)
		if err != nil {
			t.Fatalf("Error inserting test data: %v", err)
		}
	}

	// Simulate an HTTP request to the /all-times endpoint
	req, err := http.NewRequest("GET", "/all-times", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(allTimesHandler)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check that the status code is OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body to ensure it contains the correct times
	expectedResponse := `{"times":["2024-11-28 10:00:00","2024-11-28 11:00:00"]}`
	assert.JSONEq(t, expectedResponse, rr.Body.String())
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

	// Simulate a failed time zone conversion (invalid database query)
	err = logTimeToDatabase(time.Now())
	if err == nil {
		t.Fatalf("Expected error logging to database")
	}

	// Check that the log contains an error message
	expectedErrorMessage := "Error logging time to database"
	if !bytes.Contains(buf.Bytes(), []byte(expectedErrorMessage)) {
		t.Errorf("Expected log to contain %q, but got %q", expectedErrorMessage, buf.String())
	}
}
