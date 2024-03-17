package models

type Reading struct {
	ID               int    `json:"id"`
	DiscussionID     int    `json:"discussion_id"`
	Type             string `json:"type"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	URL              string `json:"url"`
	BookID           int    `json:"book_id"`
	VideoURL         string `json:"video_url"`
	DiscussionPrompt string `json:"discussion_prompt"`
}
