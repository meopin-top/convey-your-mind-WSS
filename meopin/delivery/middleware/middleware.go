package middleware

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Middleware struct {
}

func (m *Middleware) CheckWebsocketUpgrade(ctx *fiber.Ctx) error {
	log.Println("CheckWebsocketUpgrade")
	if websocket.IsWebSocketUpgrade(ctx) {
		log.Println("CheckWebsocketUpgrade success")
		ctx.Locals("allowed", true)
		return ctx.Next()
	}
	log.Println("CheckWebsocketUpgrade failed")
	return fiber.ErrUpgradeRequired
}

func InitMiddleware() *Middleware {
	return &Middleware{}
}
