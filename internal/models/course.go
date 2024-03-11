package models

import "time"

type Course struct {
	ID            int       `json:"id"`
	FacilitatorID int       `json:"facilitatorId"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
