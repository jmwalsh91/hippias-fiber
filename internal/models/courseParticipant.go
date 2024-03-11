package models

import "time"

type CourseParticipant struct {
	ID        int       `json:"id"`
	CourseID  int       `json:"courseId"`
	UserID    int       `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
