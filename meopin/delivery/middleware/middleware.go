package middleware

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Middleware struct {
}

func (m *Middleware) CheckWebsocketUpgrade(ctx *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		ctx.Locals("allowed", true)
		return ctx.Next()
	}
	return fiber.ErrUpgradeRequired
}

func InitMiddleware() *Middleware {
	return &Middleware{}
}
