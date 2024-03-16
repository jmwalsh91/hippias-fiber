package models

import "time"

// Author represents an author of books.
// swagger:model Author
type Author struct {
	// The unique identifier for the author
	// example: 1
	// required: true
	ID int `json:"id"`

	// The name of the author
	// example: Jean Baudrillard (the illest fr fr)
	// required: true
	Name string `json:"name"`

	// The nationality of the author
	// example: (we are all) American
	Nationality string `json:"nationality"`

	// A short description of the author
	// example: John Doe is a renowned American author known for his compelling novels.
	// required: true
	Description string `json:"description"`

	// The time when the author record was created
	// example: 2020-01-01T00:00:00Z
	// required: true
	CreatedAt time.Time `json:"createdAt"`

	// The last time the author record was updated
	//
}
