package services

import (
	"database/sql"
	"errors"
	"literary-lions/models"
	"literary-lions/utils"
	"log"
	"net/mail"
)

func CreateUser(email, username, password string) (*models.User, error) {
	if !IsValidEmail(email) || !IsValidUsername(username) || !IsvalidPassword(password) {
		return nil, errors.New("user parameters not valid")
	}
	post := models.User{
		Email:    email,
		Username: username,
		Password: password, // Hash here for ease
	}
	InsertUserDB(post)
	return &post, nil
}

func InsertUserDB(u models.User) error {
	_, err := utils.Getdata().Exec(`INSERT INTO users (email, username, password) VALUES (?, ?, ?)`, u.Email, u.Username, u.Password)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return err
	}
	return nil
}

func GetAllUsers() (output []models.User) {
	var u models.User
	rows, err := utils.Getdata().Query(`SELECT * FROM users`)
	if err != nil {
		log.Printf("Error querying users: %s", err)
		return output
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.Bio, &u.CreatedAt)
		if err != nil {
			log.Printf("Error scanning user: %s", err)
			continue
		}
		output = append(output, u)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating users: %s", err)
	}

	return output
}

func GetUserByID(id int) (*models.User, error) {
	var u models.User
	err := utils.Getdata().QueryRow(`SELECT * FROM users WHERE id = ?`, id).Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.Bio, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Printf("Error querying users with ID %d: %v", id, err)
		return nil, err
	}
	return &u, nil
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	// Check if email already exists in database
	var count int
	err = utils.Getdata().QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		log.Printf("Error checking email existence: %v", err)
		return false
	}

	return count == 0 // Valid if email doesn't exist
}

func IsValidUsername(username string) bool {
	// Check if username already exists in database
	var count int
	err := utils.Getdata().QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		log.Printf("Error checking username existence: %v", err)
		return false
	}

	return count == 0 // Valid if username doesn't exist
}

func IsvalidPassword(password string) bool {
	return len(password) >= 8
}

func GetUserByEmailDB(email string) *models.User {
	var user models.User
	err := utils.Getdata().QueryRow(`SELECT id, email, username, password FROM users WHERE email = ?`, email).Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // User not found
		}
		log.Printf("Error querying user by email %s: %v", email, err)
		return nil
	}
	return &user
}

func GetUserByUsernameDB(username string) *models.User {
	var user models.User
	err := utils.Getdata().QueryRow(
		"SELECT id, email, username, bio, created_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Email, &user.Username, &user.Bio, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Printf("Error getting user by username: %v", err)
		return nil
	}
	return &user
}

func UpdateUserBio(userID int, bio string) error {
	_, err := utils.Getdata().Exec(
		"UPDATE users SET bio = ? WHERE id = ?",
		bio, userID,
	)
	if err != nil {
		log.Printf("Error updating user bio: %v", err)
		return err
	}
	return nil
}
