package models

type Like struct {
	TargetID int    `json:"target_id" db:"target_id"`
	UserID   int    `json:"user_id" db:"user_id"`
	Type     string `json:"type" db:"type"`
}
