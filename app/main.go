package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"meopin-top-wss/meopin/delivery"
	"meopin-top-wss/meopin/repository/redis"
	"meopin-top-wss/meopin/usecase"

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

	// middleware := middleware.InitMiddleware()

	// app.Use("/ws", middleware.CheckWebsocketUpgrade)

	paperRepo := redis.GetInstance()
	if err := paperRepo.Client.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	paperUsecase := usecase.NewPaperUsecase(paperRepo)
	delivery.NewWsHandler(app, paperUsecase)

	data, _ := json.MarshalIndent(app.Stack(), "", "  ")
	fmt.Println(string(data))
	log.Fatal(app.Listen(":3000"))
	// Access ws://localhost:3000/ws/{channel_id}
}
