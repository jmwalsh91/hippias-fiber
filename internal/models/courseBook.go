package models

import "time"

type CourseBook struct {
	ID        int       `json:"id"`
	CourseID  int       `json:"courseId"`
	BookID    int       `json:"bookId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
