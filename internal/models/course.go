package models

import "time"

type Course struct {
	ID            int    `json:"id"`
	FacilitatorID int    `json:"facilitator_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	PhotoUrl      string `json:"photo_url"`
}

type CourseDetails struct {
	Course      Course       `json:"course"`
	Facilitator Facilitator  `json:"facilitator"`
	Books       []Book       `json:"books"`
	Schedules   []CourseWeek `json:"schedules"`
}

type CourseWeek struct {
	Week     int       `json:"week"`
	Meetings []Meeting `json:"meetings"`
}

type Meeting struct {
	Day        string    `json:"day"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	LocationID int       `json:"locationId"`
}
