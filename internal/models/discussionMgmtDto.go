package models

type DiscussionMgmtDto struct {
	Discussion   Discussion             `json:"discussion"`
	Participants []CourseParticipantDto `json:"participants"`
	Readings     []ReadingDto           `json:"readings"`
	Attendance   []DiscussionAttendance `json:"attendance"`
}
