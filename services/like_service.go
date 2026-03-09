package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"literary-lions/models"
	"literary-lions/utils"
	"log"
)

// LikeOrDislike allows a user to like/dislike a post or comment
func LikeOrDislike(targetID, userID int, table, likeType string) error {
	if userID == 0 {
		return errors.New("user not authenticated")
	}
	likes := GetLikesByID(targetID, table)
	// Remove previous like/dislike by this user for this target
	for i, l := range likes {
		if l.UserID == userID {
			likes = append(likes[:i], likes[i+1:]...)
			break
		}
	}

	// Only add a new like if likeType is not empty (empty means unlike/remove)
	if likeType != "" {
		like := models.Like{
			TargetID: targetID,
			UserID:   userID,
			Type:     likeType,
		}
		likes = append(likes, like)
	}

	InsertLikeDB(targetID, table, likes)
	return nil
}

func InsertLikeDB(targetID int, table string, likesSlice []models.Like) error {
	_, err := utils.Getdata().Exec(fmt.Sprintf(`UPDATE %s SET likes = ? WHERE id = ?`, table), utils.ToJson(likesSlice), targetID)
	if err != nil {
		log.Printf("Error inserting likes: %v", err)
		return err
	}
	return nil
}

func GetAllLikes() (postLikes, commentLikes [][]models.Like) {
	var jsonLikes string
	var targetID int
	row, err := utils.Getdata().Query(`SELECT likes, id FROM posts`)
	if err != nil {
		log.Printf("Error querying posts: %s", err)
	}
	for row.Next() {
		err := row.Scan(&jsonLikes, &targetID)
		if err != nil {
			log.Printf("Error querying posts: %s", err)
		}
		postLikes = append(postLikes, LikeConvert(jsonLikes))
	}
	row, err = utils.Getdata().Query(`SELECT likes, id FROM comments`)
	if err != nil {
		log.Printf("Error querying posts: %s", err)
	}
	for row.Next() {
		err := row.Scan(&jsonLikes, &targetID)
		if err != nil {
			log.Printf("Error querying posts: %s", err)
		}
		commentLikes = append(commentLikes, LikeConvert(jsonLikes))
	}
	return postLikes, commentLikes
}

func GetLikesByID(id int, table string) (tableLikes []models.Like) {
	var jsonLikes string
	var targetID int
	err := utils.Getdata().QueryRow(fmt.Sprintf(`SELECT likes, id FROM %s WHERE id = ?`, table), id).Scan(&jsonLikes, &targetID)
	if err != nil {
		log.Printf("Error querying %s with ID %d: %s", table, id, err)
	}
	return LikeConvert(jsonLikes)
}

func LikeConvert(jsonLikes string) (likes []models.Like) {
	json.Unmarshal([]byte(jsonLikes), &likes)
	return likes
}

func CountLikes(id int, table string) (likeCount, dislikeCount int) {
	for _, v := range GetLikesByID(id, table) {
		switch v.Type {
		case "+":
			likeCount++
		case "-":
			dislikeCount++
		}
	}
	return likeCount, dislikeCount
}
