package utils

import (
	"strings"
	"testing"
)

// func TestPasswordHashing(t *testing.T) {
// 	password := "testpassword123"

// 	// Test password hashing
// 	hash, err := handlers.HashPassword(password)
// 	if err != nil {
// 		t.Fatalf("Failed to hash password: %v", err)
// 	}

// 	if hash == password {
// 		t.Fatal("Password hash should not be the same as the original password")
// 	}

// 	// Test password checking
// 	if !utils.CheckPassword(password, hash) {
// 		t.Fatal("Password check should return true for correct password")
// 	}

// 	// Test wrong password
// 	if CheckPassword("wrongpassword", hash) {
// 		t.Fatal("Password check should return false for wrong password")
// 	}
// }

func TestEmailValidation(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name@domain.co.uk",
		"user+tag@example.org",
	}

	invalidEmails := []string{
		"invalid-email",
		"@example.com",
		"user@",
		"",
	}

	for _, email := range validEmails {
		if !isValidEmail(email) {
			t.Errorf("Email %s should be valid", email)
		}
	}

	for _, email := range invalidEmails {
		if isValidEmail(email) {
			t.Errorf("Email %s should be invalid", email)
		}
	}
}

func TestUsernameValidation(t *testing.T) {
	validUsernames := []string{
		"user123",
		"user_name",
		"User123",
		"abc",
		"username123",
	}

	invalidUsernames := []string{
		"user-name",
		"user.name",
		"user name",
		"ab",
		"",
		"user@name",
	}

	for _, username := range validUsernames {
		if !isValidUsername(username) {
			t.Errorf("Username %s should be valid", username)
		}
	}

	for _, username := range invalidUsernames {
		if isValidUsername(username) {
			t.Errorf("Username %s should be invalid", username)
		}
	}
}

// Helper functions for testing (copied from handlers)
func isValidEmail(email string) bool {
	// Simple email validation - check for @ symbol and basic structure
	return strings.Contains(email, "@") &&
		strings.Contains(email, ".") &&
		len(email) >= 3 &&
		email[0] != '@' &&
		email[len(email)-1] != '@' &&
		email[0] != '.' &&
		email[len(email)-1] != '.'
}

func isValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	for _, char := range username {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}
	return true
}
