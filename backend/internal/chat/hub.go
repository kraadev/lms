package chat

import (
	"log"
	"sync"
)

type Hub struct {
	mu         sync.RWMutex
	rooms      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Just tracking, room subscription happens explicitly in JoinRoom
			_ = client
		case client := <-h.unregister:
			h.UnregisterClient(client)
		}
	}
}

func (h *Hub) JoinRoom(client *Client, roomKey string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// If client was previously in another room, remove
	if client.roomKey != "" && client.roomKey != roomKey {
		if room, ok := h.rooms[client.roomKey]; ok {
			delete(room, client)
			if len(room) == 0 {
				delete(h.rooms, client.roomKey)
			}
		}
	}

	client.roomKey = roomKey
	if _, ok := h.rooms[roomKey]; !ok {
		h.rooms[roomKey] = make(map[*Client]bool)
	}
	h.rooms[roomKey][client] = true
	log.Printf("[WS HUB] User %s (%s) joined room %s", client.user.Name, client.user.Email, roomKey)
}

func (h *Hub) UnregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client.roomKey != "" {
		if room, ok := h.rooms[client.roomKey]; ok {
			delete(room, client)
			if len(room) == 0 {
				delete(h.rooms, client.roomKey)
			}
		}
	}
	close(client.send)
	log.Printf("[WS HUB] Client disconnected: %s", client.user.Email)
}

func (h *Hub) BroadcastToRoom(roomKey string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, ok := h.rooms[roomKey]
	if !ok {
		return
	}

	for client := range room {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(room, client)
		}
	}
}
