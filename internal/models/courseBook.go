package models

type CourseBook struct {
	ID        int    `json:"id"`
	CourseID  int    `json:"course_id"`
	BookID    int    `json:"book_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
