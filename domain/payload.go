package domain

// Payload is the data structure for the websocket
type Payload struct {
	ProjectID string  `json:"project_id"`
	UserID    string  `json:"user_id"`
	Content   Content `json:"content"`
}

// Content is the data structure for the websocket
type Content struct {
	// status
	// json camel case
	UserID      string `json:"user_id"`
	ContentID   string `json:"content_id"`
	ContentType string `json:"content_type"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Text        string `json:"text"`
	ImageURL    string `json:"image_url"` // could be empty
}

// Project is the data structure for the project
type Project struct {
	Status    string    `json:"status"` // will be enum later
	ProjectID string    `json:"project_id"`
	Contents  []Content `json:"contents"`
}
