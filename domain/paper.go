package domain

import "github.com/gofiber/contrib/websocket"

// PaperRepository is the interface for the paper repository
type PaperRepository interface {
	// Get returns the paper with the given id
	Get(id string) (string, error)
	GetLock(id string) (string, error)
	IncrLock(id string) error
	DecrLock(id string) error
	// GetAndAdd returns the paper with the given id and adds the value to the key
	Add(id string, value string) error
	Set(id string, value string) error
}

// PaperUsecase is the interface for the paper usecase
type PaperUsecase interface {
	PushData(payload Payload) error
	BroadcastMessage(paperID string, msg string) error
	Subscribe(paperID string, conn *websocket.Conn) error
	Remove(paperID string, conn *websocket.Conn) error
	GetData(payload Payload) (string, error)
	GetProject(projectID string) (string, error)

	// Create Dummy Project
	CreateDummyProject() error

	// BroadCast(paperID string, userID string, message string, messageType string) error
}
