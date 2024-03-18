package models

import "time"

type Discussion struct {
	ID          int       `json:"id"`
	CourseID    int       `json:"course_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DateTime    time.Time `json:"date_time"`
}
