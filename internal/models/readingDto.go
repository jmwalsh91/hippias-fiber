package models

type ReadingDto struct {
	Reading Reading         `json:"reading"`
	Ratings []ReadingRating `json:"ratings"`
}
