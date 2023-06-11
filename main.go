package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/j1mmyson/wss/db"
)

type Payload struct {
	ChannelID   string `json:"channel_id"`
	UserID      string `json:"user_id"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

func main() {
	app := fiber.New()
	// client := redis.NewClient(&redis.Options{
	// 	Addr: "localhost:6379",
	// })

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

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		var (
			// mt  int
			msg []byte
			err error
		)

		redisConn, err := db.NewDatabase("localhost:6379")
		if err != nil {
			panic(err)
		}
		defer redisConn.Client.Close()

		sub := redisConn.Subscribe(c.Params("id"))
		defer sub.Close()

		for {
			if _, msg, err = c.ReadMessage(); err != nil {
				log.Println("read: ", err)
				break
			}
			var payload Payload
			if err := json.Unmarshal(msg, &payload); err != nil {
				panic(err)
			}
			log.Println(payload)

			// if pmsg, err := sub.ReceiveMessage(context.Background()); err != nil {
			// 	panic(err)
			// } else {
			// 	if err := c.WriteMessage(mt, pmsg); err != nil {
			// 		log.Println("write: ", err)
			// 		break
			// 	}
			// }
		}
	}))

	log.Fatal(app.Listen(":3000"))
	// Access ws://localhost:3000/ws/123?v=1.0
}

// func Broadcast(channel string, message string) {
// 	// client.Publish(channel, message)
// }
