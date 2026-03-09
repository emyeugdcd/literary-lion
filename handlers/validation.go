package handlers

import (
	"literary-lions/services"
	"literary-lions/utils"
	"strings"
)

// ValidateRegistrationInput performs validation for registration form inputs and returns a slice of error messages.
func ValidateRegistrationInput(email, username, password, confirmPassword string) []string {
	var errs []string

	// Trim whitespace from inputs
	email = strings.TrimSpace(email)
	username = strings.TrimSpace(username)

	// Email validation
	if email == "" {
		errs = append(errs, "Email is required")
	} else if !services.IsValidEmail(email) {
		errs = append(errs, "Please enter a valid email address")
	} else if utils.CheckExistDB("users", "email", email) {
		errs = append(errs, "Email already exists")
	}

	// Username validation
	if username == "" {
		errs = append(errs, "Username is required")
	} else if len(username) < 3 || len(username) > 20 {
		errs = append(errs, "Username must be between 3 and 20 characters")
	} else if utils.CheckExistDB("users", "username", username) {
		errs = append(errs, "Username already exists")
	}

	// else if !IsValidUsername(username) {
	// 	errs = append(errs, "Username can only contain letters, numbers, and underscores")
	// }

	// Password validation
	if password == "" {
		errs = append(errs, "Password is required")
	} else if len(password) < 8 {
		errs = append(errs, "Password must be at least 8 characters long")
	}

	// Confirm password validation
	if password != confirmPassword {
		errs = append(errs, "Passwords do not match")
	}

	return errs
}

// ValidateLoginInput performs validation for login form inputs and returns error messages if any.
func ValidateLoginInput(email, password string) []string {
	var errs []string
	// Trim whitespace from email
	email = strings.TrimSpace(email)
	if email == "" || password == "" {
		errs = append(errs, "Email and password are required")
	}
	return errs
}
