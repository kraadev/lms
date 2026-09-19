package chat

import (
	"testing"

	"lms/internal/models"
)

func TestHubRoomManagement(t *testing.T) {
	hub := NewHub()

	user1 := &models.User{ID: 1, Name: "User One", Email: "one@lms.local"}
	user2 := &models.User{ID: 2, Name: "User Two", Email: "two@lms.local"}

	client1 := &Client{user: user1, send: make(chan []byte, 10)}
	client2 := &Client{user: user2, send: make(chan []byte, 10)}

	roomKey := "class_10"

	// 1. Join
	hub.JoinRoom(client1, roomKey)
	hub.JoinRoom(client2, roomKey)

	if count := hub.ClientCount(roomKey); count != 2 {
		t.Errorf("expected 2 clients in room %s, got %d", roomKey, count)
	}

	// 2. Unregister
	hub.UnregisterClient(client1)
	if count := hub.ClientCount(roomKey); count != 1 {
		t.Errorf("expected 1 client in room %s after unregister, got %d", roomKey, count)
	}

	// 3. Dead connection cleanup
	deadClient := &Client{user: &models.User{ID: 99}, send: nil}
	hub.JoinRoom(deadClient, "room_dead")
	cleaned := hub.CleanupDeadConnections()
	if cleaned < 1 {
		t.Errorf("expected at least 1 dead connection cleaned, got %d", cleaned)
	}
}
