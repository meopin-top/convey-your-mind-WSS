package main

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Payload struct {
	ChannelID   string `json:"channel_id"`
	UserID      string `json:"user_id"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	log.Fatal(app.Listen(":3000"))
	// Access ws://localhost:3000/ws/123?v=1.0
}
