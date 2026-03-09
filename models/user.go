package models

import (
	"database/sql"
	"time"
)

// User represents a user in the system, matching the users table in the database
// Fields: id, email, username, password
// Add json tags for potential API use
type User struct {
	ID        int          `json:"id"`
	Email     string       `json:"email"`
	Username  string       `json:"username"`
	Password  string       `json:"password"`
	Bio       string       `json:"bio"`
	CreatedAt sql.NullTime `json:"created_at"`
}

// GetCreatedAtFormatted returns the formatted created_at time or a default value if null
func (u *User) GetCreatedAtFormatted(layout string) string {
	if u.CreatedAt.Valid {
		return u.CreatedAt.Time.Format(layout)
	}
	return "Not set"
}

// GetCreatedAt returns the time.Time value or current time if null
func (u *User) GetCreatedAt() time.Time {
	if u.CreatedAt.Valid {
		return u.CreatedAt.Time
	}
	return time.Now() // or return a default time
}

// Profile Page Data struct for profile.html
type ProfileData struct {
	Username        string
	Bio             string
	JoinDate        string
	IsOwnProfile    bool
	IsAuthenticated bool
}
