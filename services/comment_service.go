package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"literary-lions/models"
	"literary-lions/utils"
	"log"
	"time"
)

// CreateComment allows a registered user to comment on a post
func CreateComment(content string, userID, postID int, tags []string) (*models.Comment, error) {
	if utils.CheckExistDB("posts", "id", postID) {
		if userID == 0 {
			return nil, errors.New("user not authenticated")
		} else {
			// Check for duplicate comment
			if CheckDuplicateComment(content, userID, postID) {
				return nil, errors.New("duplicate comment: you have already posted this exact comment on this post")
			}

			comment := models.Comment{
				Content: content,
				PostID:  postID,
				UserID:  userID,
				Time:    time.Now(),
				Tags:    tags,
			}
			InsertOrUpdateComment(comment)
			return &comment, nil
		}
	} else {
		return nil, errors.New("post does not exist")
	}
}

func InsertCommentDB(c models.Comment) error {
	_, err := utils.Getdata().Exec(`INSERT INTO comments (content, user_id, post_id, time, tags, files) VALUES (?, ?, ?, ?, ?, ?)`,
		c.Content, c.UserID, c.PostID, c.Time, utils.ToJson(c.Tags), utils.ToJson(c.Files))
	if err != nil {
		log.Printf("Error inserting comment: %v", err)
		return err
	}
	return nil
}

// UpdateCommentDB updates an existing comment
func UpdateCommentDB(c models.Comment) error {
	_, err := utils.Getdata().Exec(`UPDATE comments SET content = ?, tags = ?, files = ? WHERE id = ?`,
		c.Content, utils.ToJson(c.Tags), utils.ToJson(c.Files), c.ID)
	if err != nil {
		log.Printf("Error updating comment: %v", err)
		return err
	}
	return nil
}

// CheckDuplicateComment checks if a comment with same content from same user on same post exists
func CheckDuplicateComment(content string, userID, postID int) bool {
	comments := GetAllComments()
	for _, comment := range comments {
		if comment.Content == content && comment.UserID == userID && comment.PostID == postID {
			return true
		}
	}
	return false
}

// InsertOrUpdateComment checks if comment exists and updates it, otherwise inserts new comment
func InsertOrUpdateComment(c models.Comment) error {
	if c.ID != 0 && utils.CheckExistDB("comments", "id", c.ID) {
		return UpdateCommentDB(c)
	} else {
		return InsertCommentDB(c)
	}
}

func GetAllComments() (output []models.Comment) {
	var c models.Comment
	var tags, timeStr, likes, files string
	row, err := utils.Getdata().Query(`SELECT * FROM comments`)
	if err != nil {
		log.Printf("Error querying comments: %s", err)
	}
	for row.Next() {
		err := row.Scan(&c.ID, &c.Content, &c.PostID, &c.UserID, &timeStr, &likes, &tags, &files)
		if err != nil {
			log.Printf("Error querying comments: %s", err)
		} else {
			// Parse time string
			c.Time = utils.ParseTimeString(timeStr)
			// Parse likes from JSON string to slice
			if likes != "none" && likes != "" {
				json.Unmarshal([]byte(likes), &c.Likes)
			}
			json.Unmarshal([]byte(tags), &c.Tags)
			json.Unmarshal([]byte(files), &c.Files)
			output = append(output, c)
		}
	}
	return output
}

func GetCommentByID(id int) (*models.Comment, error) {
	var c models.Comment
	var tags, likes, files, timeStr string
	err := utils.Getdata().QueryRow(`SELECT * FROM comments WHERE id = ?`, id).Scan(&c.ID, &c.Content, &c.PostID, &c.UserID, &timeStr, &likes, &tags, &files)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Printf("Error querying comments with ID %d: %v", id, err)
		return nil, err
	}
	// Parse time string
	c.Time = utils.ParseTimeString(timeStr)
	json.Unmarshal([]byte(likes), &c.Likes)
	json.Unmarshal([]byte(tags), &c.Tags)
	json.Unmarshal([]byte(files), &c.Files)
	return &c, nil
}

func GetCommentsByPostID(postID int) (output []models.Comment, err error) {
	for _, v := range GetAllComments() {
		if v.PostID == postID {
			output = append(output, v)
		}
	}
	return output, nil
}

func CountComments(postID int) (count int) {
	for _, v := range GetAllComments() {
		if v.PostID == postID {
			count++
		}
	}
	return count
}

func GetCommentsWithUsernamesByPostID(postID int) ([]models.CommentWithUsername, error) {
	comments, err := GetCommentsByPostID(postID)
	if err != nil {
		return nil, err
	}

	var commentsWithUsernames []models.CommentWithUsername
	for _, comment := range comments {
		usernameInterface, err := utils.FetchWithID("users", "username", comment.UserID)
		if err != nil {
			log.Printf("Error fetching username for user ID %d: %v", comment.UserID, err)
			continue
		}

		username := ""
		if usernameInterface != nil {
			username = usernameInterface.(string)
		}

		commentWithUsername := models.CommentWithUsername{
			Comment:  comment,
			Username: username,
		}
		commentsWithUsernames = append(commentsWithUsernames, commentWithUsername)
	}

	return commentsWithUsernames, nil
}
