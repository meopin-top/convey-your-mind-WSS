package domain

// Payload is the data structure for the websocket
type Payload struct {
	PaperID     string `json:"paper_id"`
	UserID      string `json:"user_id"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}
