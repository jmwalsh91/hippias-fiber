package models

type DiscussionDto struct {
	Discussion
	Readings   []Reading              `json:"readings"`
	Ratings    []ReadingRating        `json:"ratings"`
	Attendance []DiscussionAttendance `json:"attendance"`
}
