package main

import (
	"encoding/json"
	"log"
	"meopin-top-wss/domain/singleton"
	"meopin-top-wss/meopin/repository/redis"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// Payload is the data structure for the websocket
type Payload struct {
	PaperID     string `json:"paper_id"`
	UserID      string `json:"user_id"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

func main() {
	app := fiber.New()

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	app.Use("/ws", func(ctx *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(ctx) {
			ctx.Locals("allowed", true)
			return ctx.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:paper_id", websocket.New(func(conn *websocket.Conn) {
		paperID := conn.Params("channel_id")
		subscriber := singleton.GetBroakerInstance()
		// add connection to channel
		subscriber.Add(paperID, conn)

		defer func() {
			// remove connection from channel
			subscriber.Remove(paperID, conn)
			// close websocket connection
			conn.Close()
		}()

		// Get local value: 필요 없어지면 삭제 예정
		allowed := conn.Locals("allowed").(bool)
		log.Println(paperID, allowed)

		for {
			// wait for new message
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("read message failed:", err)
				break
			}

			log.Printf("recv: %s\n", msg)

			var payload Payload
			err = json.Unmarshal(msg, &payload)
			if err != nil {
				log.Println("unmarshal failed:", err)
				break
			}

			// push message to redis
			db := redis.GetInstance()
			err = db.GetAndAdd(paperID, string(msg))
			if err != nil {
				log.Println("publish failed:", err)
				break
			}

			strMsg, err := db.Get(paperID)
			if err != nil {
				log.Println("get failed:", err)
				break
			}

			msg = []byte(strMsg)
			go subscriber.Broadcast(paperID, msg)

			log.Printf("send: %s\n", msg)
		}

	}))

	log.Fatal(app.Listen(":3000"))
	// Access ws://localhost:3000/ws/{channel_id}
}
