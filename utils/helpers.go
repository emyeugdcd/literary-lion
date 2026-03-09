package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// ParseTimeString parses time string with multiple formats
func ParseTimeString(timeStr string) time.Time {
	// Try different time formats
	formats := []string{
		time.RFC3339Nano,                      // 2006-01-02T15:04:05.999999999Z07:00 (RFC3339 with nanoseconds)
		time.RFC3339,                          // 2006-01-02T15:04:05Z07:00 (RFC3339)
		"2006-01-02T15:04:05.999999999-07:00", // ISO format with nanoseconds and timezone
		"2006-01-02T15:04:05.999999999+07:00", // ISO format with nanoseconds and timezone (positive)
		"2006-01-02 15:04:05.999999999-07:00", // With timezone and nanoseconds
		"2006-01-02 15:04:05.999999999+07:00", // With timezone and nanoseconds (positive)
		"2006-01-02 15:04:05-07:00",           // With timezone
		"2006-01-02 15:04:05+07:00",           // With timezone (positive)
		"2006-01-02 15:04:05",                 // Standard format
		"2006-01-02T15:04:05Z",                // ISO format
		"2006-01-02T15:04:05.999999999Z",      // ISO format with nanoseconds
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t
		}
	}

	// If all parsing fails, return current time
	log.Printf("Could not parse time string: %s", timeStr)
	return time.Now()
}

// IsDatabaseRunning checks if the database is accessible and functioning properly
func IsDBRunning() bool {
	db, err := sql.Open("sqlite3", "./database/data.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return false
	}
	defer db.Close()

	// Ping the database to verify connection
	if err := db.Ping(); err != nil {
		log.Printf("Error pinging database: %v", err)
		return false
	}

	// Try to query a table to further verify functionality
	_, err = db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='users'")
	if err != nil {
		log.Printf("Error querying database: %v", err)
		return false
	}

	return true
}

func CheckExistDB(table, field string, value any) bool {
	var output any
	err := Getdata().QueryRow(fmt.Sprintf(`SELECT %s FROM %s WHERE %s = ?`, field, table, field), value).Scan(&output)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Printf("Error checking existence in %s: %v", table, err)
		return false
	}
	return true
}

func Delete() {

}

func ToJson(input any) string {
	output, err := json.Marshal(input)
	if err != nil {
		log.Printf("Error marshaling: %s", err)
	}
	return string(output)
}
