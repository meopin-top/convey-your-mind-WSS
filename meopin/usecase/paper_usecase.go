package usecase

import (
	"meopin-top-wss/domain"
	"meopin-top-wss/domain/singleton"

	"github.com/gofiber/contrib/websocket"
)

type paperUsecase struct {
	paperRepo domain.PaperRepository
	broaker   singleton.Broaker
}

// NewPaperUsecase is the constructor for paperUsecase
func NewPaperUsecase(paperRepo domain.PaperRepository, broaker singleton.Broaker) domain.PaperUsecase {
	return &paperUsecase{
		paperRepo: paperRepo,
		broaker:   broaker,
	}
}

func (p *paperUsecase) Subscribe(paperID string, conn *websocket.Conn) error {
	p.broaker.Add(paperID, conn)
	return nil
}

func (p *paperUsecase) ReceiveMessage(payload domain.Payload) error {
	db := p.paperRepo
	// Get paper from redis
	paper, err := db.Get(payload.PaperID)
	if err != nil {
		return err
	}

	// Add message to paper
	err = db.GetAndAdd(payload.PaperID, paper)
	if err != nil {
		return err
	}

	return nil
}

func (p *paperUsecase) GetData(payload domain.Payload) (string, error) {
	db := p.paperRepo
	// Get paper from redis
	paper, err := db.Get(payload.PaperID)
	if err != nil {
		return "", err
	}

	return paper, nil
}

func (p *paperUsecase) BroadcastMessage(paperID string, msg string) error {

	p.broaker.Broadcast(paperID, []byte(msg))
	return nil
}

func (p *paperUsecase) Remove(paperID string, conn *websocket.Conn) error {
	p.broaker.Remove(paperID, conn)
	return nil
}
