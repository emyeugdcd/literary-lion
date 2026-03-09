package handlers

import (
	"fmt"
	"literary-lions/models"
	"literary-lions/services"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

var tmpl *template.Template

func init() {
	var err error
	// Create function map for template functions
	funcMap := template.FuncMap{
		"add":  func(a, b int) int { return a + b },
		"sub":  func(a, b int) int { return a - b },
		"join": func(slice []string, sep string) string { return strings.Join(slice, sep) },
	}

	// Parse templates from both directories
	tmpl, err = template.New("").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		fmt.Println("Error parsing main templates:", err)
		return
	}

	// Parse templates from elements directory
	tmpl, err = tmpl.ParseGlob("templates/elements/*.html")
	if err != nil {
		fmt.Println("Error parsing element templates:", err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	isAuth := IsAuthenticated(r)
	fmt.Println("IsAuthenticated:", isAuth)
	var currentUser *models.User
	if isAuth {
		var err error
		currentUser, err = GetCurrentUser(r)
		if err != nil {
			// If there's an error getting the user, clear the session
			DestroySession(w, r)
			isAuth = false
			currentUser = nil
		}
	}

	data := map[string]interface{}{
		"IsAuthenticated": isAuth,
		"CurrentUser":     currentUser,
	}

	err := tmpl.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		fmt.Println("Error executing template:", err)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		email := strings.TrimSpace(r.FormValue("email"))
		username := strings.TrimSpace(r.FormValue("username"))
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		// Perform validation using centralized function
		errs := ValidateRegistrationInput(email, username, password, confirmPassword)
		if len(errs) > 0 {
			data := map[string]interface{}{
				"Error":    strings.Join(errs, ", "),
				"Email":    email,
				"Username": username,
			}
			tmpl.ExecuteTemplate(w, "register.html", data)
			return
		}

		// Hash the password
		hashedPassword, err := HashPassword(password)
		if err != nil {
			http.Error(w, "Error processing registration", http.StatusInternalServerError)
			return
		}

		// Insert user into database
		_, err = services.CreateUser(email, username, hashedPassword)
		if err != nil {
			log.Printf("Error creating account: %v", err)
			data := map[string]interface{}{
				"Error":    "Error creating account. Please try again.",
				"Email":    email,
				"Username": username,
			}
			tmpl.ExecuteTemplate(w, "register.html", data)
			return
		}

		// Redirect to login with success message
		http.Redirect(w, r, "/login?success=Account created successfully! Please log in.", http.StatusSeeOther)
		return
	}

	// Show registration form with empty default values
	data := map[string]interface{}{
		"Email":    "",
		"Username": "",
	}
	err := tmpl.ExecuteTemplate(w, "register.html", data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		fmt.Println("Error executing template:", err)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		// Perform validation using centralized function
		errs := ValidateLoginInput(email, password)
		if len(errs) > 0 {
			data := map[string]interface{}{
				"Error": strings.Join(errs, ", "),
				"Email": email,
			}
			tmpl.ExecuteTemplate(w, "login.html", data)
			return
		}

		// Get user from database
		user := services.GetUserByEmailDB(email)
		if user == nil {
			data := map[string]interface{}{
				"Error": "Invalid email or password",
				"Email": email,
			}
			tmpl.ExecuteTemplate(w, "login.html", data)
			return
		}

		// Check password
		if !CheckPassword(password, user.Password) {
			data := map[string]interface{}{
				"Error": "Invalid email or password",
				"Email": email,
			}
			tmpl.ExecuteTemplate(w, "login.html", data)
			return
		}

		// Create session
		err := CreateSession(w, r, *user)
		if err != nil {
			http.Error(w, "Error creating session", http.StatusInternalServerError)
			return
		}

		// Redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Check for success message from registration
	success := r.URL.Query().Get("success")
	data := map[string]interface{}{
		"Success": success,
		"Email":   "",
	}

	err := tmpl.ExecuteTemplate(w, "login.html", data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		fmt.Println("Error executing template:", err)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := DestroySession(w, r)
	if err != nil {
		http.Error(w, "Error logging out", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	isAuth := IsAuthenticated(r)
	if !isAuth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get current user
	currentUser, err := GetCurrentUser(r)
	if err != nil {
		log.Printf("Error getting current user: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Extract form values
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))
		categoryStr := strings.TrimSpace(r.FormValue("category"))
		category := []string{categoryStr} // Convert to slice for service function
		bookTitle := strings.TrimSpace(r.FormValue("book_title"))
		author := strings.TrimSpace(r.FormValue("author"))
		tags := strings.Split(strings.TrimSpace(r.FormValue("tags")), " ")

		// Validate required fields
		var validationErrors []string
		if title == "" {
			validationErrors = append(validationErrors, "Title is required")
		}
		if len(title) > 200 {
			validationErrors = append(validationErrors, "Title must be less than 200 characters")
		}
		if content == "" {
			validationErrors = append(validationErrors, "Content is required")
		}
		if len(content) < 50 {
			validationErrors = append(validationErrors, "Content must be at least 50 characters")
		}
		if categoryStr == "" {
			validationErrors = append(validationErrors, "Category is required")
		}

		// If there are validation errors, show the form again with errors
		if len(validationErrors) > 0 {
			errorMessage := strings.Join(validationErrors, "; ")
			data := models.PostData{
				Title:           title,
				Content:         content,
				CategoryStr:     categoryStr, // Use string for template
				BookTitle:       bookTitle,
				Author:          author,
				Tags:            tags,
				Error:           errorMessage,
				IsAuthenticated: isAuth,
				CurrentUser:     currentUser,
			}
			err := tmpl.ExecuteTemplate(w, "create_post.html", data)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				return
			}
			return
		}

		// Prepare content fields
		var fullContent string
		if bookTitle != "" || author != "" {
			fullContent = content
			if bookTitle != "" {
				fullContent += "\n\nBook: " + bookTitle
			}
			if author != "" {
				fullContent += "\nAuthor: " + author
			}
		} else {
			fullContent = content
		}

		// Insert post into database
		_, err = services.CreatePost(title, fullContent, currentUser.ID, category, tags)
		if err != nil {
			log.Printf("Error inserting post: %v", err)
			data := models.PostData{
				Title:           title,
				Content:         content,
				CategoryStr:     categoryStr, // Use string for template
				BookTitle:       bookTitle,
				Author:          author,
				Tags:            tags,
				Error:           "Failed to create post. Please try again.",
				IsAuthenticated: isAuth,
				CurrentUser:     currentUser,
			}
			err := tmpl.ExecuteTemplate(w, "create_post.html", data)
			if err != nil {
				log.Printf("Error executing template: %v", err)
				return
			}
			return
		}

		// Redirect to home page or post
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Handle GET request - show the form
	data := models.PostData{
		IsAuthenticated: isAuth,
		CurrentUser:     currentUser,
	}
	err = tmpl.ExecuteTemplate(w, "create_post.html", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	// Get filter and category from query parameters
	filter := r.URL.Query().Get("filter")
	if filter == "" {
		filter = "all"
	}

	category := r.URL.Query().Get("category")

	// Get page parameter
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Set posts per page
	postsPerPage := 10

	// Check if user is authenticated
	isAuthenticated := IsAuthenticated(r)
	var userID int
	var currentUser *models.User
	if isAuthenticated {
		user, err := GetCurrentUser(r)
		if err == nil {
			userID = user.ID
			currentUser = user
		}
	}

	// Get posts from database with pagination
	var posts []models.Post
	var totalCount int
	var err error

	switch filter {
	case "my":
		if isAuthenticated {
			posts, totalCount, err = services.FilterPostsByUserWithPagination(userID, page, postsPerPage)
		} else {
			posts, totalCount, err = services.GetPostsWithPagination(page, postsPerPage)
		}
	case "liked":
		if isAuthenticated {
			posts, totalCount, err = services.FilterPostsLikedByUserWithPagination(userID, page, postsPerPage)
		} else {
			posts, totalCount, err = services.GetPostsWithPagination(page, postsPerPage)
		}
	case "hot", "recent":
		// For now, treat hot and recent the same as all posts
		// TODO: Implement proper sorting by likes/views for "hot" and by creation date for "recent"
		if category != "" {
			posts, totalCount, err = services.FilterPostsByCategoryWithPagination(category, page, postsPerPage)
		} else {
			posts, totalCount, err = services.GetPostsWithPagination(page, postsPerPage)
		}
	default:
		if category != "" {
			posts, totalCount, err = services.FilterPostsByCategoryWithPagination(category, page, postsPerPage)
		} else {
			posts, totalCount, err = services.GetPostsWithPagination(page, postsPerPage)
		}
	}

	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		http.Error(w, "Error loading posts", http.StatusInternalServerError)
		return
	}

	// Calculate pagination data
	totalPages := (totalCount + postsPerPage - 1) / postsPerPage
	var pagination *models.PaginationData
	if totalPages > 1 {
		pagination = &models.PaginationData{
			CurrentPage:  page,
			TotalPages:   totalPages,
			HasPrevious:  page > 1,
			HasNext:      page < totalPages,
			PreviousPage: page - 1,
			NextPage:     page + 1,
		}
	}

	// Add TotalCount to PostsPageData
	data := models.PostsPageData{
		Posts:           posts,
		IsAuthenticated: isAuthenticated,
		Filter:          filter,
		CategoryFilter:  category,
		CurrentUser:     currentUser,
		Pagination:      pagination,
		TotalCount:      totalCount,
	}

	// Execute template
	err = tmpl.ExecuteTemplate(w, "posts.html", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	// Extract post ID from URL path
	pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathSegments) < 2 {
		http.Error(w, "Invalid post URL", http.StatusBadRequest)
		return
	}

	postIDStr := pathSegments[1] // Get and then convert the ID from the URL /post/{id}
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Check if user is authenticated
	isAuthenticated := IsAuthenticated(r)
	var userID int
	var currentUser *models.User
	if isAuthenticated {
		user, err := GetCurrentUser(r)
		if err == nil {
			userID = user.ID
			currentUser = user
		}
	}

	// Get the post from database
	post, err := services.GetPostByID(postID)
	if err != nil {
		log.Printf("Error fetching post: %v", err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	// print category for debugging
	// Increment view count
	if !isAuthenticated || userID != post.UserID {
		err = services.IncrementPostViews(postID)
		if err != nil {
			log.Printf("Error incrementing post views: %v", err)
		}
	}

	// Get comments for this post
	comments, err := services.GetCommentsWithUsernamesByPostID(postID)
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		comments = []models.CommentWithUsername{}
	}

	// Check if current user has liked this post
	var hasLiked bool
	if isAuthenticated {
		hasLiked, err = services.HasUserLikedPost(userID, postID)
		if err != nil {
			log.Printf("Error checking user like status: %v", err)
			hasLiked = false
		}
	}

	// Get related posts (posts from same category, excluding current post)
	var relatedPosts []models.Post
	if post.Category != "" && strings.TrimSpace(post.Category) != "" {
		var err error
		relatedPosts, err = services.FilterPostsByCategory(post.Category)
		if err != nil {
			log.Printf("Error fetching related posts: %v", err)
			relatedPosts = []models.Post{}
		}
	} else {
		log.Printf("Post has no category, skipping related posts fetch")
		relatedPosts = []models.Post{}
	}
	// Filter out the current post from related posts
	var filteredRelatedPosts []models.Post
	for _, relatedPost := range relatedPosts {
		if relatedPost.ID != postID {
			filteredRelatedPosts = append(filteredRelatedPosts, relatedPost)
		}
	}
	// Limit to 3 related posts
	if len(filteredRelatedPosts) > 3 {
		filteredRelatedPosts = filteredRelatedPosts[:3]
	}

	// Get popular tags (for now, just some sample tags)
	popularTags := []string{"fiction", "classic", "review", "discussion", "recommendation"}

	// Prepare template data
	data := models.PostDetailData{
		Post:            post,
		Comments:        comments,
		IsAuthenticated: isAuthenticated,
		CurrentUser:     currentUser,
		HasLiked:        hasLiked,
		RelatedPosts:    filteredRelatedPosts,
		PopularTags:     popularTags,
	}

	// Execute template
	err = tmpl.ExecuteTemplate(w, "post_detail.html", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func NewsletterHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated (optional, but maybe new user can also read newsletter? I don't know)
	isAuth := IsAuthenticated(r)
	var currentUser *models.User
	if isAuth {
		var err error
		currentUser, err = GetCurrentUser(r)
		if err != nil {
			DestroySession(w, r)
			isAuth = false
			currentUser = nil
		}
	}

	data := map[string]interface{}{
		"IsAuthenticated": isAuth,
		"CurrentUser":     currentUser,
	}

	err := tmpl.ExecuteTemplate(w, "newsletter.html", data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		fmt.Println("Error executing template:", err)
	}
}

func TermsHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated (optional, same as above, we can also remove if needed)
	isAuth := IsAuthenticated(r)
	var currentUser *models.User
	if isAuth {
		var err error
		currentUser, err = GetCurrentUser(r)
		if err != nil {
			DestroySession(w, r)
			isAuth = false
			currentUser = nil
		}
	}

	data := map[string]interface{}{
		"IsAuthenticated": isAuth,
		"CurrentUser":     currentUser,
		"LastUpdated":     "January 2025",
		"ContactEmail":    "support@literarylions.com",
	}

	err := tmpl.ExecuteTemplate(w, "terms.html", data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		fmt.Println("Error executing template:", err)
	}
}

// LikePostHandler handles liking/unliking a post
func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is authenticated
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get current user
	currentUser, err := GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	// Extract post ID from URL path
	pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathSegments) < 3 {
		http.Error(w, "Invalid post URL", http.StatusBadRequest)
		return
	}

	postIDStr := pathSegments[2] // post/like/{postID} -> index 2
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Check if user has already liked this post
	hasLiked, err := services.HasUserLikedPost(currentUser.ID, postID)
	if err != nil {
		log.Printf("Error checking user like status: %v", err)
		http.Error(w, "Error processing like", http.StatusInternalServerError)
		return
	}

	// Toggle like status
	if hasLiked {
		// Unlike the post - remove the like
		err = services.LikeOrDislike(postID, currentUser.ID, "posts", "")
	} else {
		// Like the post
		err = services.LikeOrDislike(postID, currentUser.ID, "posts", "+")
	}

	if err != nil {
		log.Printf("Error toggling like: %v", err)
		http.Error(w, "Error processing like", http.StatusInternalServerError)
		return
	}

	// Redirect back to the post detail page
	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}

// AddCommentHandler handles adding a comment to a post
func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is authenticated
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get current user
	currentUser, err := GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	// Extract post ID from URL path
	pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathSegments) < 3 {
		http.Error(w, "Invalid post URL", http.StatusBadRequest)
		return
	}

	postIDStr := pathSegments[2] // post/comment/{postID} -> index 2
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		http.Error(w, "Comment content cannot be empty", http.StatusBadRequest)
		return
	}

	// Create the comment
	_, err = services.CreateComment(content, currentUser.ID, postID, []string{})
	if err != nil {
		log.Printf("Error creating comment: %v", err)
		http.Error(w, "Error creating comment", http.StatusInternalServerError)
		return
	}

	// Redirect back to the post detail page
	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}

// LikeCommentHandler handles liking/unliking a comment
func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is authenticated
	if !IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get current user
	currentUser, err := GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	// Extract comment ID from URL path
	pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathSegments) < 3 {
		http.Error(w, "Invalid comment URL", http.StatusBadRequest)
		return
	}

	commentIDStr := pathSegments[2] // comment/like/{commentID} -> index 2
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Get the comment to find the post ID for redirect
	comment, err := services.GetCommentByID(commentID)
	if err != nil {
		log.Printf("Error getting comment: %v", err)
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	// Check if user has already liked this comment
	likes := services.GetLikesByID(commentID, "comments")
	hasLiked := false
	for _, like := range likes {
		if like.UserID == currentUser.ID && like.Type == "+" {
			hasLiked = true
			break
		}
	}

	// Toggle like status
	if hasLiked {
		// Unlike the comment - remove the like
		err = services.LikeOrDislike(commentID, currentUser.ID, "comments", "")
	} else {
		// Like the comment
		err = services.LikeOrDislike(commentID, currentUser.ID, "comments", "+")
	}

	if err != nil {
		log.Printf("Error toggling comment like: %v", err)
		http.Error(w, "Error processing like", http.StatusInternalServerError)
		return
	}

	// Redirect back to the post detail page
	http.Redirect(w, r, fmt.Sprintf("/post/%d", comment.PostID), http.StatusSeeOther)
}

// SearchHandler handles search requests
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("search")
	if query == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get page parameter for pagination
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Get search results with pagination
	results, totalCount, err := services.SearchPostsWithPagination(query, page, 10)
	if err != nil {
		log.Printf("Error searching posts: %v", err)
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}

	// Calculate pagination info
	totalPages := (totalCount + 9) / 10 // Ceiling division
	hasNext := page < totalPages
	hasPrev := page > 1

	// Check if user is authenticated
	isAuthenticated := IsAuthenticated(r)
	var currentUser *models.User
	if isAuthenticated {
		currentUser, err = GetCurrentUser(r)
		if err != nil {
			http.Error(w, "Error getting user", http.StatusInternalServerError)
			return
		}
	}

	data := struct {
		Query           string
		Results         []models.Post
		ResultsCount    int
		TotalCount      int
		CurrentPage     int
		TotalPages      int
		HasNext         bool
		HasPrev         bool
		IsAuthenticated bool
		CurrentUser     *models.User
	}{
		Query:           query,
		Results:         results,
		ResultsCount:    len(results),
		TotalCount:      totalCount,
		CurrentPage:     page,
		TotalPages:      totalPages,
		HasNext:         hasNext,
		HasPrev:         hasPrev,
		IsAuthenticated: isAuthenticated,
		CurrentUser:     currentUser,
	}

	// Create template with helper functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}

	tmpl, err := template.New("search.html").Funcs(funcMap).ParseFiles("templates/search.html", "templates/elements/header.html", "templates/elements/footer.html")
	if err != nil {
		log.Printf("Error parsing search template: %v", err)
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing search template: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

// SearchSuggestionsHandler provides autocomplete suggestions,
// but requires Javascript to work so for now just for testing purposes and sake of studying
// func SearchSuggestionsHandler(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query().Get("q")
// 	if len(query) < 2 {
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode([]string{})
// 		return
// 	}

// 	suggestions, err := services.GetSearchSuggestions(query)
// 	if err != nil {
// 		log.Printf("Error getting search suggestions: %v", err)
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode([]string{})
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(suggestions)
// }

// Handlers for user profile
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract username from URL path
	username := strings.TrimPrefix(r.URL.Path, "/profile/")
	if username == "" {
		http.Error(w, "Username required", http.StatusBadRequest)
		return
	}

	// Get user data from database
	userData := services.GetUserByUsernameDB(username)
	if userData == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if user is authenticated
	isAuth := IsAuthenticated(r)
	var currentUser *models.User
	var isOwnProfile bool

	if isAuth {
		var err error
		currentUser, err = GetCurrentUser(r)
		if err != nil {
			// If there's an error getting the user, clear the session
			DestroySession(w, r)
			isAuth = false
			currentUser = nil
		} else {
			isOwnProfile = currentUser.Username == username
		}
	}

	// Prepare template data
	data := map[string]interface{}{
		"Username":        userData.Username,
		"Bio":             userData.Bio,
		"JoinDate":        userData.GetCreatedAtFormatted("January 2, 2006"),
		"IsOwnProfile":    isOwnProfile,
		"IsAuthenticated": isAuth,
		"CurrentUser":     currentUser,
	}

	// Render template using your existing tmpl variable
	err := tmpl.ExecuteTemplate(w, "profile.html", data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		fmt.Println("Error executing template:", err)
	}
}

func EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	isAuth := IsAuthenticated(r)
	if !isAuth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	currentUser, err := GetCurrentUser(r)
	if err != nil {
		DestroySession(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		// Show edit form
		data := map[string]interface{}{
			"Username":        currentUser.Username,
			"Bio":             currentUser.Bio,
			"JoinDate":        currentUser.GetCreatedAtFormatted("January 2, 2006"),
			"IsOwnProfile":    true,
			"IsAuthenticated": true,
			"CurrentUser":     currentUser,
		}

		err := tmpl.ExecuteTemplate(w, "edit_profile.html", data)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			fmt.Println("Error executing template:", err)
		}
	} else if r.Method == "POST" {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		bio := strings.TrimSpace(r.FormValue("bio"))

		// Validate bio length (optional, can also change in the future if 500 is too short? maybe?)
		if len(bio) > 500 {
			data := map[string]interface{}{
				"Error":           "Bio too long (max 500 characters)",
				"Username":        currentUser.Username,
				"Bio":             bio,
				"JoinDate":        currentUser.GetCreatedAtFormatted("January 2, 2006"),
				"IsOwnProfile":    true,
				"IsAuthenticated": true,
				"CurrentUser":     currentUser,
			}
			tmpl.ExecuteTemplate(w, "edit_profile.html", data)
			return
		}

		// Update database
		err := services.UpdateUserBio(currentUser.ID, bio)
		if err != nil {
			data := map[string]interface{}{
				"Error":           "Error updating profile. Please try again.",
				"Username":        currentUser.Username,
				"Bio":             bio,
				"JoinDate":        currentUser.GetCreatedAtFormatted("January 2, 2006"),
				"IsOwnProfile":    true,
				"IsAuthenticated": true,
				"CurrentUser":     currentUser,
			}
			tmpl.ExecuteTemplate(w, "edit_profile.html", data)
			return
		}

		// Redirect back to profile
		http.Redirect(w, r, "/profile/"+currentUser.Username, http.StatusSeeOther)
	}
}
