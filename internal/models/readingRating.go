package models

type ReadingRating struct {
	ID        int `json:"id"`
	ReadingID int `json:"reading_id"`
	UserID    int `json:"user_id"`
	Rating    int `json:"rating"`
}
