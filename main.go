package main

import (
	"database/sql"
	"fmt"
	"literary-lions/handlers"

	// "literary-lions/models"
	// "literary-lions/services"
	"literary-lions/utils"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ServerTesting()

	db, err := sql.Open("sqlite3", "./database/data.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Database opened successfully")

}

func ServerTesting() {
	//utils.InitiateDB()
	// Check if database is running
	if utils.IsDBRunning() {
		fmt.Println("Database is running properly!")
	} else {
		fmt.Println("Database is not running or has issues.")
	}

	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	// Handle routes
	// Homepage and authentication routes
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)
	mux.HandleFunc("/newsletter", handlers.NewsletterHandler)
	mux.HandleFunc("/terms", handlers.TermsHandler)

	// Post routes
	mux.HandleFunc("/posts", handlers.PostsHandler)
	mux.HandleFunc("/createpost", handlers.CreatePostHandler)
	mux.HandleFunc("/post/", handlers.ViewPostHandler)

	// Like and comment routes
	mux.HandleFunc("/post/like/", handlers.LikePostHandler)
	mux.HandleFunc("/post/comment/", handlers.AddCommentHandler)
	mux.HandleFunc("/comment/like/", handlers.LikeCommentHandler)

	// Search route
	mux.HandleFunc("/search", handlers.SearchHandler)

	//User profile routes
	mux.HandleFunc("/profile/", handlers.ProfileHandler)
	mux.HandleFunc("/profile/edit", handlers.EditProfileHandler)

	// Search suggestions API but using JavaScript so not needed here but kept for reference and future study
	//mux.HandleFunc("/api/suggestions", handlers.SearchSuggestionsHandler)

	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
