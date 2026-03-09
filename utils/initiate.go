package utils

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func Getdata() *sql.DB {
	data, err := sql.Open("sqlite3", "./database/data.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil
	}
	return data
}

func InitiateDB() error {
	var err error
	fmt.Println("Initiating database")
	defer Getdata().Close()

	categories := []string{
		"book review",
		"book recommendation",
		"author discussion",
		"reading challenge",
		"book club",
		"writing tips",
		"genre discussion",
		"general",
	}
	// Create tables for users
	_, err = Getdata().Exec(
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER NOT NULL PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
			);`,
	)

	if err != nil {
		log.Printf("Error creating users table: %s", err)
		return err
	}

	// Add bio column if not exists
	_, err = Getdata().Exec(`ALTER TABLE users ADD COLUMN bio TEXT;`)
	if err != nil && !strings.Contains(err.Error(), "duplicate column name") {
		log.Printf("Error adding bio column: %v", err)
	}

	// Add created_at column if not exists (without default)
	_, err = Getdata().Exec(`ALTER TABLE users ADD COLUMN created_at DATETIME;`)
	if err != nil && !strings.Contains(err.Error(), "duplicate column name") {
		log.Printf("Error adding created_at column: %v", err)
	} else {
		// backfill created_at for existing rows
		_, err = Getdata().Exec(`UPDATE users SET created_at = CURRENT_TIMESTAMP WHERE created_at IS NULL;`)
		if err != nil {
			log.Printf("Error backfilling created_at: %v", err)
		}
	}

	// Create tables for posts
	_, err = Getdata().Exec(
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER NOT NULL PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			time DATETIME DEFAULT CURRENT_TIMESTAMP,
			likes TEXT,
			categories TEXT NOT NULL,
			tags TEXT DEFAULT 'none',
			files TEXT DEFAULT 'none'
		);`,
	)
	if err != nil {
		log.Printf("Error creating posts table: %s", err)
		return err
	}

	// Create tables for comments
	_, err = Getdata().Exec(
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER NOT NULL PRIMARY KEY,
			content TEXT NOT NULL,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			time STRING NOT NULL,
			likes TEXT,
			tags TEXT DEFAULT 'none',
			files TEXT DEFAULT 'none'
		);`,
	)
	if err != nil {
		log.Printf(": %s", err)
		return err
	}

	// Create tables for categories
	_, err = Getdata().Exec(
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER NOT NULL PRIMARY KEY,
			title TEXT NOT NULL
		);`,
	)
	if err != nil {
		log.Printf("Error creating categories table: %s", err)
		return err
	}
	for i := 0; i < len(categories); i++ {
		_, err = Getdata().Exec(`INSERT INTO categories (title) VALUES (?)`, categories[i])
	}
	if err != nil {
		log.Printf("Error filling categories table: %s", err)
		return err
	}

	// Create tables for books, right now not being used but kept for future use
	_, err = Getdata().Exec(
		`CREATE TABLE IF NOT EXISTS books (
			isbn INTEGER NOT NULL,
			title TEXT NOT NULL,
			subtitle TEXT NOT NULL,
			genre TEXT NOT NULL,
			author TEXT NOT NULL,
			published STRING NOT NULL,
			publisher STRING NOT NULL,
			pages INTEGER NOT NULL,
			description STRING NOT NULL
		);`,
	)
	if err != nil {
		log.Printf("Error creating books table: %s", err)
		return err
	}
	for i := 0; i < len(categories); i++ {
		_, err = Getdata().Exec(`INSERT INTO categories (title) VALUES (?)`, categories[i])
	}
	if err != nil {
		log.Printf("Error filling books table: %s", err)
		return err
	}
	fmt.Println("Database successfully initiated")
	return nil
}

func FetchWithID(table, field string, id int) (any, error) {
	var output any
	err := Getdata().QueryRow(fmt.Sprintf(`SELECT %s FROM %s WHERE id = ?`, field, table), id).Scan(&output)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error fetching %s from %s with ID %d: %v", field, table, id, err)
		return nil, err
	}
	return output, nil
}

func GetLatestID(table string) int {
	var lastID int
	err := Getdata().QueryRow(fmt.Sprintf(`SELECT MAX(id) FROM %s`, table)).Scan(&lastID)
	if err != nil {
		return 0
	}
	return lastID
}
