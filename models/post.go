package models

import "time"

// PostData for rendering create_post.html
type PostData struct {
	Title           string
	Content         string
	Category        []string // Keep for backward compatibility
	CategoryStr     string   // Single category string for form display
	BookTitle       string
	Author          string // Book author's name
	UserID          string // User creating the post (kept for backward compatibility)
	Tags            []string
	Error           string
	IsAuthenticated bool
	CurrentUser     *User
}

// Post represents a post in the system, matching the posts table in the database
// Fields: id, title, content, likes, files
// Add json tags for potential API use
type Post struct {
	ID             int       `json:"id" db:"id"`
	Title          string    `json:"title" db:"title"`
	Content        string    `json:"content" db:"content"`
	ContentExcerpt string    `json:"content_excerpt"`
	UserID         int       `json:"user_id" db:"user_id"`
	Username       string    `json:"username"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	LikesRaw       string    `json:"-" db:"likes"`
	Likes          []Like    `json:"likes"`
	Category       string    `json:"category" db:"category"`
	CategoryName   string    `json:"category_name"`
	Tags           []string  `json:"tags" db:"tags"`
	Files          string    `json:"files"`
	ViewsCount     int       `json:"views_count"`
	LikesCount     int       `json:"likes_count"`
	CommentsCount  int       `json:"comments_count"`
}

// PostPageData represents all data needed for the view_post.html
type PostPageData struct {
	Post     *Post
	Comments []Comment
	User     *User
}

// PostsPageData represents all data needed for the posts.html template
type PostsPageData struct {
	Posts           []Post          `json:"posts"`
	IsAuthenticated bool            `json:"is_authenticated"`
	Filter          string          `json:"filter"`
	CategoryFilter  string          `json:"category_filter"`
	CurrentUser     *User           `json:"current_user"`
	Pagination      *PaginationData `json:"pagination"`
	TotalCount      int             `json:"total_count"`
}

// PaginationData represents pagination information
type PaginationData struct {
	CurrentPage  int  `json:"current_page"`
	TotalPages   int  `json:"total_pages"`
	HasPrevious  bool `json:"has_previous"`
	HasNext      bool `json:"has_next"`
	PreviousPage int  `json:"previous_page"`
	NextPage     int  `json:"next_page"`
}

// PostDetailData for rendering post_detail.html
type PostDetailData struct {
	Post            *Post                 `json:"post"`
	Comments        []CommentWithUsername `json:"comments"`
	IsAuthenticated bool                  `json:"is_authenticated"`
	CurrentUser     *User                 `json:"current_user"`
	HasLiked        bool                  `json:"has_liked"`
	RelatedPosts    []Post                `json:"related_posts"`
	PopularTags     []string              `json:"popular_tags"`
}

// CommentWithUsername represents a comment with username for template rendering
type CommentWithUsername struct {
	Comment
	Username string `json:"username"`
}

// Comment represents a comment in the system, matching the comments table in the database
// Fields: id, postid, content, likes, files
// Add json tags for potential API use
type Comment struct {
	ID      int       `json:"id"`
	Content string    `json:"content"`
	PostID  int       `json:"postid"`
	UserID  int       `json:"user_id" db:"user_id"`
	Time    time.Time `json:"time" db:"time"`
	Likes   []Like    `json:"likes"`
	Tags    []string  `json:"tags"`
	Files   []string  `json:"files"`
}
