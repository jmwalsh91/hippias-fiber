package models

type CourseMgmtDto struct {
	Course       Course                 `json:"course"`
	Discussions  []DiscussionDto        `json:"discussions"`
	Participants []CourseParticipantDto `json:"participants"`
}
