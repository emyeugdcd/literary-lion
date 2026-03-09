package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"literary-lions/models"
	"literary-lions/services"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	// Store for session management
	store = sessions.NewCookieStore([]byte("your-secret-key-change-this-in-production"))

	// Session configuration
	sessionName     = "literary-lions-session"
	sessionDuration = 24 * time.Hour // 24 hours
)

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a password with its hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSessionToken creates a random session token
func GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// CreateSession creates a new session for a user
func CreateSession(w http.ResponseWriter, r *http.Request, user models.User) error {
	session, err := store.Get(r, sessionName)
	if err != nil {
		return err
	}

	// Set session values
	session.Values["user_id"] = user.ID
	session.Values["email"] = user.Email
	session.Values["username"] = user.Username
	session.Values["authenticated"] = true

	// Set session options
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(sessionDuration.Seconds()),
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	return session.Save(r, w)
}

// GetSession retrieves the current session
func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, sessionName)
}

// IsAuthenticated checks if the user is authenticated
func IsAuthenticated(r *http.Request) bool {
	session, err := GetSession(r)
	if err != nil {
		return false
	}

	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}

// GetCurrentUser retrieves the current user from the session
func GetCurrentUser(r *http.Request) (*models.User, error) {
	session, err := GetSession(r)
	if err != nil {
		return nil, err
	}

	if !IsAuthenticated(r) {
		return nil, errors.New("user not authenticated")
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		return nil, errors.New("invalid user ID in session")
	}
	userData, err := services.GetUserByID(userID)

	if userData == nil {
		// User doesn't exist in database, session is invalid
		return nil, errors.New("user not found in database")
	}

	return &models.User{
		ID:        userID,
		Email:     userData.Email,
		Username:  userData.Username,
		Password:  "", // Don't return password
		Bio:       userData.Bio,
		CreatedAt: userData.CreatedAt,
	}, nil
}

// DestroySession removes the user session
func DestroySession(w http.ResponseWriter, r *http.Request) error {
	session, err := GetSession(r)
	if err != nil {
		return err
	}

	// Clear session values
	session.Values["authenticated"] = false
	session.Values["user_id"] = nil
	session.Values["email"] = nil
	session.Values["username"] = nil

	// Set session to expire immediately
	session.Options.MaxAge = -1

	return session.Save(r, w)
}

// AuthMiddleware is middleware to protect routes that require authentication
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
