package main

import (
	"encoding/json"
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

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:channel_id", websocket.New(func(c *websocket.Conn) {
		channelID := c.Params("channel_id")
		subscribe(channelID)

		defer func() {
			unsubscribe(channelID)
		}()

		// Get local value
		allowed := c.Locals("allowed").(bool)

		log.Println(channelID, allowed)

		go func() {
			// when broadcast is called get message from channel and send to client
			for {
				msg := <-channel

				jsonMsg, err := json.Marshal(msg)
				if err != nil {
					log.Println("marshal:", err)
					break
				}

				// send message to client
				err = c.WriteMessage(websocket.TextMessage, jsonMsg)
				if err != nil {
					log.Println("write:", err)
					break
				}

				log.Printf("sent: %s\n", jsonMsg)
			}
		}()

		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			var payload Payload
			err = json.Unmarshal(msg, &payload)
			if err != nil {
				log.Println("unmarshal:", err)
				break
			}

			// update channel with new message
			updateChannel(channelID, payload)

			// Send message to all subscribers
			broadcast(channelID, payload)

			log.Printf("recv: %s\n", msg)
			err = c.WriteMessage(mt, msg)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}

	}))

	log.Fatal(app.Listen(":3000"))
	// Access ws://localhost:3000/ws/{channel_id}
}

func unsubscribe(channelID string) {
	panic("unimplemented")
}

func broadcast(channelID string, payload Payload) {
	panic("unimplemented")
}

func updateChannel(channelID string, payload Payload) {
	panic("unimplemented")
}

func subscribe(channelID string) {
	panic("unimplemented")
}
