package subscribe

import "github.com/gofiber/contrib/websocket"

// Subscribe is the data structure for the websocket
type Subscribe struct {
	ChannelMap map[string][]*websocket.Conn
}

// singleton instance
var instance *Subscribe

// GetInstance returns the singleton instance
func GetInstance() *Subscribe {
	if instance == nil {
		instance = &Subscribe{
			ChannelMap: make(map[string][]*websocket.Conn),
		}
	}
	return instance
}

// Add adds a websocket connection to the channel
func (s *Subscribe) Add(channelID string, conn *websocket.Conn) {
	s.ChannelMap[channelID] = append(s.ChannelMap[channelID], conn)
}

// Remove removes a websocket connection from the channel
func (s *Subscribe) Remove(channelID string, conn *websocket.Conn) {
	conns := s.ChannelMap[channelID]
	for i, c := range conns {
		if c == conn {
			s.ChannelMap[channelID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
}

// Broadcast sends a message to all connections in the channel
func (s *Subscribe) Broadcast(channelID string, message []byte) {
	conns := s.ChannelMap[channelID]
	for _, conn := range conns {
		if error := conn.WriteMessage(websocket.TextMessage, message); error != nil {
			s.Remove(channelID, conn)
		}
	}
}
