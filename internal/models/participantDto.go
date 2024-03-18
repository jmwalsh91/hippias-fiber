package models

type CourseParticipantDto struct {
	CourseParticipant
	User User `json:"user"`
}
