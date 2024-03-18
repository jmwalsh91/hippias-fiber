package models

type DiscussionAttendance struct {
	ID           int  `json:"id"`
	DiscussionID int  `json:"discussion_id"`
	UserID       int  `json:"user_id"`
	Attended     bool `json:"attended"`
}
