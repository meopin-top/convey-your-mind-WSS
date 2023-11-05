package usecase

import (
	"encoding/json"
	"log"
	"meopin-top-wss/domain"
	"meopin-top-wss/domain/singleton"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type paperUsecase struct {
	paperRepo domain.PaperRepository
	broaker   singleton.Broaker
}

var (
	waitLockDelay = 100 * time.Millisecond
)

// NewPaperUsecase is the constructor for paperUsecase
func NewPaperUsecase(paperRepo domain.PaperRepository) domain.PaperUsecase {
	return &paperUsecase{
		paperRepo: paperRepo,
		broaker:   *singleton.GetBroakerInstance(),
	}
}

func (p *paperUsecase) Subscribe(paperID string, conn *websocket.Conn) error {
	p.broaker.Add(paperID, conn)
	return nil
}

func (p *paperUsecase) CreateDummyProject() error {
	db := p.paperRepo
	dummy := domain.Project{
		Status:    "active",
		ProjectID: "abc",
		Contents: []domain.Content{
			{
				UserID:      "byungwook",
				ContentID:   "123",
				ContentType: "text",
				X:           100,
				Y:           100,
				Width:       200,
				Height:      300,
				Text:        "hello world",
				ImageURL:    "",
			},
		},
	}
	dummyString, _ := json.Marshal(dummy)
	if err := db.Set(dummy.ProjectID, string(dummyString)); err != nil {
		log.Println("create dummy failed:", err)
		return err
	}

	return nil
}

func (p *paperUsecase) PushData(payload domain.Payload) error {
	// https://stackoverflow.com/questions/46536234/how-to-parse-json-in-golang-without-unmarshaling-twice

	db := p.paperRepo
	var lock string
	var err error

	defer func() {
		if err := db.DecrLock(payload.ProjectID); err != nil {
			log.Println("unlock failed:", err)
			return
		}
	}()

	// Redis Lock
	for {
		if lock, err = db.GetLock(payload.ProjectID); err != nil {
			log.Println("get lock failed:", err)
			return err
		}
		if lock == "0" || lock == "" {
			if err = db.IncrLock(payload.ProjectID); err != nil {
				log.Println("lock failed:", err)
				return err
			}
			break
		}
		time.Sleep(waitLockDelay)
	}

	// Get paper from redis
	paper, err := db.Get(payload.ProjectID)
	if err != nil {
		return err
	}
	log.Println("data loaded: " + paper)

	// Processing paper data
	// jsonMap := map[string]json.RawMessage{}
	jsonMap := map[string]string{}

	if err = json.Unmarshal([]byte(paper), &jsonMap); err != nil {
		log.Println("json unmarshal failed:", err)
		return err
	}

	log.Println("jsonMap: ", jsonMap)

	contents := jsonMap["contents"]

	contentsMap := map[string]string{}

	if err = json.Unmarshal([]byte(contents), &contentsMap); err != nil {
		return err
	}
	v, _ := json.Marshal(payload.Content)
	contentsMap[payload.Content.ContentID] = string(v)
	v, _ = json.Marshal(contentsMap)
	jsonMap["contents"] = string(v)

	jsonMapString, err := json.Marshal(jsonMap)
	if err != nil {
		log.Println("json marshal failed:", err)
		return err
	}

	log.Println("data saved: " + string(jsonMapString))
	// Set paper to redis
	if err = db.Set(payload.ProjectID, string(jsonMapString)); err != nil {
		log.Println("set failed:", err)
		return err
	}
	// err = db.Set(payload.ProjectID, string(jsonMapString))

	// Redis Unlock

	return nil
}

func (p *paperUsecase) GetProject(projectID string) (string, error) {
	db := p.paperRepo
	// Get paper from redis
	paper, err := db.Get(projectID)
	if err != nil {
		return "", err
	}

	return paper, nil
}

func (p *paperUsecase) GetData(payload domain.Payload) (string, error) {
	db := p.paperRepo
	// Get paper from redis
	paper, err := db.Get(payload.ProjectID)
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
