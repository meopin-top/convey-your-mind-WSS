package delivery

import (
	"encoding/json"
	"log"
	"meopin-top-wss/domain"
	"meopin-top-wss/meopin/delivery/middleware"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type wsHandler struct {
	paperUsecase domain.PaperUsecase
}

// NewWsHandler is the constructor for wsHandler
func NewWsHandler(c *fiber.App, paperUsecase domain.PaperUsecase) {
	handler := &wsHandler{
		paperUsecase: paperUsecase,
	}

	_middleware := middleware.InitMiddleware()

	c.Get("/ping", handler.Ping)
	c.Use("/ws", _middleware.CheckWebsocketUpgrade)
	c.Get("/ws/:paper_id", websocket.New(handler.WebsocketConnection))
}

func (m *wsHandler) Ping(c *fiber.Ctx) error {
	return c.SendString("pong")
}

func (m *wsHandler) WebsocketConnection(conn *websocket.Conn) {
	paperID := conn.Params("channel_id")
	m.paperUsecase.Subscribe(paperID, conn)

	defer func() {
		// remove connection from channel
		m.paperUsecase.Remove(paperID, conn)
		// close websocket connection
		conn.Close()
	}()

	// Get local value, 필요 없어지면 삭제 예정
	allowed := conn.Locals("allowed").(bool)
	log.Println(paperID, allowed)

	var payload domain.Payload
	for {
		// wait for new message
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read message failed:", err)
			break
		}

		log.Printf("recv: %s\n", msg)

		err = json.Unmarshal(msg, &payload)
		if err != nil {
			log.Println("unmarshal failed:", err)
			break
		}

		// push message to redis
		m.paperUsecase.ReceiveMessage(payload)

		if err != nil {
			log.Println("publish failed:", err)
			break
		}

		strMsg, err := m.paperUsecase.GetData(payload)

		if err != nil {
			log.Println("get data failed:", err)
			break
		}
		m.paperUsecase.BroadcastMessage(paperID, strMsg)

		log.Printf("send: %s\n", msg)
	}
	return
}
