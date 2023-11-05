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
	c.Get("/dummy", handler.CreateDummyProject)
}

func (m *wsHandler) Ping(c *fiber.Ctx) error {
	return c.SendString("pong")
}

func (m *wsHandler) CreateDummyProject(c *fiber.Ctx) error {
	m.paperUsecase.CreateDummyProject()

	return nil
}

func (m *wsHandler) WebsocketConnection(conn *websocket.Conn) {
	paperID := conn.Params("paper_id")
	log.Printf("new websocket connection: %s\n", paperID)
	m.paperUsecase.Subscribe(paperID, conn)

	// remove connection from channel when connection is closed
	defer func() {
		m.paperUsecase.Remove(paperID, conn)
		conn.Close()
	}()

	// Get local value, 필요 없어지면 삭제 예정
	allowed := conn.Locals("allowed").(bool)
	log.Println(paperID, allowed)

	// 최초 연결 시 프로젝트 데이터 수신
	project, err := m.paperUsecase.GetProject(paperID)
	if err != nil {
		log.Println("get project failed:", err)
		return
	}

	// send project data to client
	err = conn.WriteMessage(websocket.TextMessage, []byte(project))

	var payload domain.Payload
	for {
		// wait for new message
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read message failed:", err)
			break
		}
		log.Println("receive message:", string(msg))

		if err = json.Unmarshal(msg, &payload); err != nil {
			log.Println("unmarshal failed:", err)
			break
		}

		go m.paperUsecase.PushData(payload) // fire and forget
		m.paperUsecase.BroadcastMessage(paperID, string(msg))
		log.Printf("broadcast message: %s\n", msg)
	}
	return
}
