package handler

import (
	"log"
	"net/http"
	"time"

	"messenger/internal/domain"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Be more restrictive in production
	},
}

type WebSocketHandler struct {
	service   domain.MessageService
	clients   map[*websocket.Conn]bool
	broadcast chan domain.Message
}

func NewWebSocketHandler(service domain.MessageService) *WebSocketHandler {
	return &WebSocketHandler{
		service:   service,
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan domain.Message),
	}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
		return
	}
	defer conn.Close()

	h.clients[conn] = true

	// Send recent messages to new client
	messages, err := h.service.GetRecentMessages(50)
	if err != nil {
		log.Printf("Error getting recent messages: %v", err)
	} else {
		for _, msg := range messages {
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("Error sending message: %v", err)
				break
			}
		}
	}

	for {
		var msg domain.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			delete(h.clients, conn)
			break
		}

		msg.Timestamp = time.Now()
		if err := h.service.SendMessage(msg); err != nil {
			log.Printf("Error saving message: %v", err)
		}

		h.broadcast <- msg
	}
}

func (h *WebSocketHandler) StartBroadcasting() {
	for {
		msg := <-h.broadcast
		for client := range h.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error broadcasting message: %v", err)
				client.Close()
				delete(h.clients, client)
			}
		}
	}
}
