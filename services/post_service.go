package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"literary-lions/models"
	"literary-lions/utils"
	"log"
	"strings"
	"time"
)

// parseTags safely parses tags from JSON string with fallback to comma-separated format
func parseTags(tagsString string) []string {
	var tags []string

	// Try to parse as JSON first
	err := json.Unmarshal([]byte(tagsString), &tags)
	if err != nil {
		// Fallback: try to parse as comma-separated string
		if strings.Contains(tagsString, ",") {
			tagSlice := strings.Split(tagsString, ",")
			for _, tag := range tagSlice {
				cleaned := strings.TrimSpace(tag)
				if cleaned != "" {
					tags = append(tags, cleaned)
				}
			}
		} else if tagsString != "" && tagsString != "none" {
			cleaned := strings.TrimSpace(tagsString)
			if cleaned != "" {
				tags = []string{cleaned}
			}
		}
	}

	return tags
}

// CheckDuplicatePost checks if a post with same title from same user exists
func CheckDuplicatePost(title string, userID int) bool {
	posts := GetAllPosts()
	for _, post := range posts {
		if post.Title == title && post.UserID == userID {
			return true
		}
	}
	return false
}

// CreatePost allows a registered user to create a post with categories
func CreatePost(title, content string, userID int, categories, tags []string) (*models.Post, error) {
	if userID == 0 {
		return nil, errors.New("user not authenticated")
	}

	// Check for duplicate post
	if CheckDuplicatePost(title, userID) {
		return nil, errors.New("duplicate post: you have already created a post with this title")
	}

	// Handle category - use the first category as the main category for DB storage
	var category string
	if len(categories) > 0 {
		category = categories[0]
	}

	post := models.Post{
		Title:     title,
		Content:   content,
		UserID:    userID,
		Category:  category, // Set the singular category for DB storage
		Tags:      tags,
		CreatedAt: time.Now(),
	}
	InsertOrUpdatePost(post)
	return &post, nil
}

func InsertPostDB(p models.Post) error {
	_, err := utils.Getdata().Exec(`INSERT INTO posts (title, content, user_id, time, category, tags, files) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		p.Title, p.Content, p.UserID, p.CreatedAt, p.Category, utils.ToJson(p.Tags), p.Files)
	if err != nil {
		log.Printf("Error inserting post: %v", err)
		return err
	}
	return nil
}

// UpdatePostDB updates an existing post
func UpdatePostDB(p models.Post) error {
	_, err := utils.Getdata().Exec(`UPDATE posts SET title = ?, content = ?, category = ?, tags = ?, files = ? WHERE id = ?`,
		p.Title, p.Content, p.Category, utils.ToJson(p.Tags), p.Files, p.ID)
	if err != nil {
		log.Printf("Error updating post: %v", err)
		return err
	}
	return nil
}

// InsertOrUpdatePost checks if post exists and updates it, otherwise inserts new post
func InsertOrUpdatePost(p models.Post) error {
	if p.ID != 0 && utils.CheckExistDB("posts", "id", p.ID) {
		return UpdatePostDB(p)
	} else {
		return InsertPostDB(p)
	}
}

// GetAllPosts retrieves all posts from the database
func GetAllPosts() (output []models.Post) {
	var p models.Post
	var tags, likes, files, timeStr string
	row, err := utils.Getdata().Query(`SELECT p.id, p.title, p.content, p.user_id, p.time, p.likes, p.category, p.tags, p.files, p.views_count, u.username FROM posts p LEFT JOIN users u ON p.user_id = u.id ORDER BY p.time DESC`)
	if err != nil {
		log.Printf("Error querying posts: %s", err)
	}
	for row.Next() {
		err := row.Scan(&p.ID, &p.Title, &p.Content, &p.UserID, &timeStr, &likes, &p.Category, &tags, &files, &p.ViewsCount, &p.Username)
		if err != nil {
			log.Printf("Error querying posts: %s", err)
		} else {
			// Parse time
			p.CreatedAt = utils.ParseTimeString(timeStr)

			// Parse tags with error handling
			p.Tags = parseTags(tags)

			p.LikesRaw = likes
			// Parse likes count
			if likes != "none" && likes != "" {
				var likesSlice []models.Like
				json.Unmarshal([]byte(likes), &likesSlice)
				p.LikesCount = len(likesSlice)
			} else {
				p.LikesCount = 0
			}
			// Set comments count to 0 for now (can be calculated later if needed)
			p.CommentsCount = CountComments(p.ID)
			p.ContentExcerpt = ContentExcerpt(p.Content, 200) // Generate content excerpt
			p.Files = files
			output = append(output, p)
		}
	}
	return output
}

// GetPostsWithPagination returns paginated posts
func GetPostsWithPagination(page, limit int) ([]models.Post, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // Default limit
	}

	offset := (page - 1) * limit

	var posts []models.Post
	var p models.Post
	var categories, tags, likes string

	// Get total count first
	var totalCount int
	err := utils.Getdata().QueryRow("SELECT COUNT(*) FROM posts").Scan(&totalCount)
	if err != nil {
		log.Printf("Error counting posts: %v", err)
		return nil, 0, err
	}

	// Get paginated posts
	query := `SELECT p.*, u.username FROM posts p LEFT JOIN users u ON p.user_id = u.id ORDER BY p.time DESC LIMIT ? OFFSET ?`
	rows, err := utils.Getdata().Query(query, limit, offset)
	if err != nil {
		log.Printf("Error querying paginated posts: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.UserID, &p.CreatedAt, &likes, &categories, &tags, &p.Files, &p.ViewsCount, &p.Username)
		if err != nil {
			log.Printf("Error scanning post: %v", err)
			continue
		}

		p.Tags = parseTags(tags)
		p.LikesRaw = likes

		// Parse likes count
		if likes != "none" && likes != "" {
			var likesSlice []models.Like
			json.Unmarshal([]byte(likes), &likesSlice)
			p.LikesCount = len(likesSlice)
		} else {
			p.LikesCount = 0
		}

		p.CommentsCount = CountComments(p.ID)
		p.ContentExcerpt = ContentExcerpt(p.Content, 200)
		posts = append(posts, p)
	}

	return posts, totalCount, nil
}

// FilterPostsByUser returns posts created by a user
func FilterPostsByUser(userID int) ([]models.Post, error) {
	var result []models.Post
	checkExist := utils.CheckExistDB("users", "id", userID)
	// print
	log.Printf("Checking if user with ID %d exists: %v", userID, checkExist)
	if checkExist {
		for _, v := range GetAllPosts() {
			if v.UserID == userID {
				result = append(result, v)
			}
		}
	} else {
		return nil, errors.New("user does not exist")
	}
	return result, nil
}

// FilterPostsByUserWithPagination returns paginated posts created by a user
func FilterPostsByUserWithPagination(userID, page, limit int) ([]models.Post, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	if !utils.CheckExistDB("users", "id", userID) {
		return nil, 0, errors.New("user does not exist")
	}

	var posts []models.Post
	var p models.Post
	var categories, tags, likes string

	// Get total count first
	var totalCount int
	err := utils.Getdata().QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&totalCount)
	if err != nil {
		log.Printf("Error counting user posts: %v", err)
		return nil, 0, err
	}

	// Get paginated posts
	query := `SELECT * FROM posts WHERE user_id = ? ORDER BY time DESC LIMIT ? OFFSET ?`
	rows, err := utils.Getdata().Query(query, userID, limit, offset)
	if err != nil {
		log.Printf("Error querying paginated user posts: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.UserID, &p.CreatedAt, &likes, &categories, &tags, &p.Files)
		if err != nil {
			log.Printf("Error scanning post: %v", err)
			continue
		}

		json.Unmarshal([]byte(tags), &p.Tags)
		p.LikesRaw = likes

		// Parse likes count
		if likes != "none" && likes != "" {
			var likesSlice []models.Like
			json.Unmarshal([]byte(likes), &likesSlice)
			p.LikesCount = len(likesSlice)
		} else {
			p.LikesCount = 0
		}

		p.CommentsCount = CountComments(p.ID)
		p.ContentExcerpt = ContentExcerpt(p.Content, 200)
		posts = append(posts, p)
	}

	return posts, totalCount, nil
}

// FilterPostsByCategory returns posts filtered by category
func FilterPostsByCategory(category string) ([]models.Post, error) {
	var result []models.Post

	// Check if category is empty or just whitespace
	if strings.TrimSpace(category) == "" {
		return nil, errors.New("category cannot be empty")
	}

	if utils.CheckExistDB("categories", "title", category) {
		for _, v := range GetAllPosts() {
			if v.Category == category {
				result = append(result, v)
			}
		}
	} else {
		return nil, errors.New("category does not exist")
	}
	return result, nil
}

// FilterPostsByCategoryWithPagination returns paginated posts filtered by category
func FilterPostsByCategoryWithPagination(category string, page, limit int) ([]models.Post, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Check if category is empty or just whitespace
	if strings.TrimSpace(category) == "" {
		return nil, 0, errors.New("category cannot be empty")
	}

	if !utils.CheckExistDB("categories", "title", category) {
		return nil, 0, errors.New("category does not exist")
	}

	// For now, we'll use the existing logic but add pagination
	// This is not the most efficient way, but maintains compatibility
	allPosts := GetAllPosts()
	var filteredPosts []models.Post

	for _, post := range allPosts {
		if post.Category == category {
			filteredPosts = append(filteredPosts, post)
		}
	}

	totalCount := len(filteredPosts)
	offset := (page - 1) * limit

	// Apply pagination to filtered results
	var paginatedPosts []models.Post
	for i := offset; i < offset+limit && i < len(filteredPosts); i++ {
		paginatedPosts = append(paginatedPosts, filteredPosts[i])
	}

	return paginatedPosts, totalCount, nil
}

func GetPostByID(id int) (*models.Post, error) {
	var p models.Post
	var tags, likes, files, timeStr string

	// Try to get the display name for the category (with spaces) or fall back to the slug
	query := `SELECT p.id, p.title, p.content, p.user_id, p.time, p.likes, p.category, p.tags, p.files, p.views_count, u.username, 
		COALESCE(
			(SELECT title FROM categories WHERE title = REPLACE(p.category, '-', ' ') LIMIT 1),
			(SELECT title FROM categories WHERE title = p.category LIMIT 1),
			p.category
		) as category_name 
		FROM posts p 
		LEFT JOIN users u ON p.user_id = u.id 
		WHERE p.id = ?`

	err := utils.Getdata().QueryRow(query, id).Scan(&p.ID, &p.Title, &p.Content, &p.UserID, &timeStr, &likes, &p.Category, &tags, &files, &p.ViewsCount, &p.Username, &p.CategoryName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Printf("Error querying posts with ID %d: %v", id, err)
		return nil, err
	}
	// Parse time
	p.CreatedAt = utils.ParseTimeString(timeStr)
	p.LikesRaw = likes
	// Parse likes count
	if likes != "none" && likes != "" {
		var likesSlice []models.Like
		json.Unmarshal([]byte(likes), &likesSlice)
		p.LikesCount = len(likesSlice)
	} else {
		p.LikesCount = 0
	}
	// Set comments count to 0 for now
	p.CommentsCount = CountComments(p.ID)
	p.Tags = parseTags(tags)
	json.Unmarshal([]byte(files), &p.Files)

	return &p, nil
}

// FilterPostsLikedByUser returns posts liked by a user
func FilterPostsLikedByUser(userID int) ([]models.Post, error) {
	var result []models.Post
	checkExist := utils.CheckExistDB("users", "id", userID)
	// print
	log.Printf("Checking1 if user with ID %d exists: %v", userID, checkExist)
	if checkExist {
		for _, v := range GetAllPosts() {
			for i := range v.Likes {
				if v.Likes[i].UserID == userID {
					result = append(result, v)
				}
			}
		}
	} else {
		return nil, errors.New("user does not exist")
	}
	return result, nil
}

// FilterPostsLikedByUserWithPagination returns paginated posts liked by a user
func FilterPostsLikedByUserWithPagination(userID, page, limit int) ([]models.Post, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	if !utils.CheckExistDB("users", "id", userID) {
		return nil, 0, errors.New("user does not exist")
	}

	// For now, we'll use the existing logic but add pagination
	// This could be optimized with a proper SQL query joining likes table
	allLikedPosts, err := FilterPostsLikedByUser(userID)
	if err != nil {
		return nil, 0, err
	}

	totalCount := len(allLikedPosts)
	offset := (page - 1) * limit

	// Apply pagination to filtered results
	var paginatedPosts []models.Post
	for i := offset; i < offset+limit && i < len(allLikedPosts); i++ {
		paginatedPosts = append(paginatedPosts, allLikedPosts[i])
	}

	return paginatedPosts, totalCount, nil
}

// IncrementPostViews increments the view count for a post
func IncrementPostViews(postID int) error {
	db, err := sql.Open("sqlite3", "./database/data.db")
	if err != nil {
		return err
	}
	defer db.Close()

	query := `UPDATE posts SET views_count = COALESCE(views_count, 0) + 1 WHERE id = ?`
	_, err = db.Exec(query, postID)
	return err
}

// HasUserLikedPost checks if a user has liked the current post
func HasUserLikedPost(userID, postID int) (bool, error) {
	likes := GetLikesByID(postID, "posts")
	for _, like := range likes {
		if like.UserID == userID && like.Type == "+" {
			return true, nil
		}
	}
	return false, nil
}

// ContentExcerpt returns a truncated version of the post content for previews
func ContentExcerpt(content string, maxLength int) string {
	if len(content) <= maxLength {
		return content
	}

	// Find the last space before maxLength to avoid cutting words
	truncated := content[:maxLength]
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > 0 && lastSpace > maxLength-20 { // Don't go too far back
		truncated = content[:lastSpace]
	}

	return truncated + "..."
}

// SearchPosts searches for posts based on title, content, tags, and categories
func SearchPosts(query string) ([]models.Post, error) {
	if query == "" {
		return nil, errors.New("search query cannot be empty")
	}

	var result []models.Post
	searchTerm := strings.ToLower(strings.TrimSpace(query))

	for _, post := range GetAllPosts() {
		// Search in title
		if strings.Contains(strings.ToLower(post.Title), searchTerm) {
			result = append(result, post)
			continue
		}

		// Search in content
		if strings.Contains(strings.ToLower(post.Content), searchTerm) {
			result = append(result, post)
			continue
		}

		// Search in tags
		for _, tag := range post.Tags {
			if strings.Contains(strings.ToLower(tag), searchTerm) {
				result = append(result, post)
				goto nextPost // Found match, skip to next post
			}
		}

		// Search in category
		if strings.Contains(strings.ToLower(post.Category), searchTerm) {
			result = append(result, post)
			continue
		}

		// Search in username
		if strings.Contains(strings.ToLower(post.Username), searchTerm) {
			result = append(result, post)
			continue
		}

	nextPost:
	}

	return result, nil
}

// SearchPostsWithPagination returns paginated search results
func SearchPostsWithPagination(query string, page, limit int) ([]models.Post, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Get all search results first
	allResults, err := SearchPosts(query)
	if err != nil {
		return nil, 0, err
	}

	totalCount := len(allResults)
	offset := (page - 1) * limit

	// Apply pagination to search results
	var paginatedPosts []models.Post
	for i := offset; i < offset+limit && i < len(allResults); i++ {
		paginatedPosts = append(paginatedPosts, allResults[i])
	}

	return paginatedPosts, totalCount, nil
}

// GetSearchSuggestions returns search suggestions based on partial query
func GetSearchSuggestions(query string) ([]string, error) {
	if len(query) < 2 {
		return nil, errors.New("query too short")
	}

	var suggestions []string
	searchTerm := strings.ToLower(strings.TrimSpace(query))
	suggestionSet := make(map[string]bool) // To avoid duplicates

	for _, post := range GetAllPosts() {
		// Check titles
		if strings.Contains(strings.ToLower(post.Title), searchTerm) {
			if !suggestionSet[post.Title] {
				suggestions = append(suggestions, post.Title)
				suggestionSet[post.Title] = true
			}
		}

		// Check tags
		for _, tag := range post.Tags {
			if strings.Contains(strings.ToLower(tag), searchTerm) {
				if !suggestionSet[tag] {
					suggestions = append(suggestions, tag)
					suggestionSet[tag] = true
				}
			}
		}

		// Check categories

		if strings.Contains(strings.ToLower(post.Category), searchTerm) {
			if !suggestionSet[post.Category] {
				suggestions = append(suggestions, post.Category)
				suggestionSet[post.Category] = true
			}
		}

		// Limit suggestions to prevent too many results
		if len(suggestions) >= 10 {
			break
		}
	}

	return suggestions, nil
}
