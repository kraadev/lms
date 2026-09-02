package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"lms/internal/middleware"
	"lms/internal/models"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 65536
)

type Client struct {
	hub              *Hub
	conn             *websocket.Conn
	send             chan []byte
	user             *models.User
	roomKey          string
	currentClassID   int64
	repo             *Repository
	accessController *middleware.AccessController
}

type IncomingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type JoinPayload struct {
	ClassID int64 `json:"class_id"`
}

type SendPayload struct {
	ClassID int64  `json:"class_id"`
	Message string `json:"message"`
}

type OutgoingEvent struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WS] Read error for user %s: %v", c.user.Email, err)
			}
			break
		}

		var incoming IncomingMessage
		if err := json.Unmarshal(messageData, &incoming); err != nil {
			c.sendError("INVALID_PAYLOAD", "Could not parse WebSocket message")
			continue
		}

		switch incoming.Type {
		case "class.join":
			var p JoinPayload
			if err := json.Unmarshal(incoming.Payload, &p); err != nil || p.ClassID <= 0 {
				c.sendError("BAD_REQUEST", "Valid class_id is required")
				continue
			}

			// Validate class access server-side
			hasAccess, err := c.accessController.CheckClassAccess(c.user.ID, c.user.Role, p.ClassID)
			if err != nil || !hasAccess {
				c.sendError("FORBIDDEN", "You are not authorized to join this class chat room.")
				continue
			}

			roomKey := fmt.Sprintf("class:%d", p.ClassID)
			c.currentClassID = p.ClassID
			c.hub.JoinRoom(c, roomKey)

			c.sendEvent("class.joined", map[string]interface{}{
				"class_id": p.ClassID,
				"room":     roomKey,
				"status":   "connected",
			})

		case "chat.send":
			var p SendPayload
			if err := json.Unmarshal(incoming.Payload, &p); err != nil {
				c.sendError("BAD_REQUEST", "Invalid chat payload")
				continue
			}

			p.Message = strings.TrimSpace(p.Message)
			if p.Message == "" {
				c.sendError("BAD_REQUEST", "Message content cannot be empty")
				continue
			}
			if len(p.Message) > 2000 {
				c.sendError("BAD_REQUEST", "Message exceeds maximum length of 2000 characters")
				continue
			}

			targetClassID := p.ClassID
			if targetClassID <= 0 {
				targetClassID = c.currentClassID
			}

			if targetClassID <= 0 {
				c.sendError("BAD_REQUEST", "You must join a class room before sending messages")
				continue
			}

			// Validate class access server-side
			hasAccess, err := c.accessController.CheckClassAccess(c.user.ID, c.user.Role, targetClassID)
			if err != nil || !hasAccess {
				c.sendError("FORBIDDEN", "You are not a member of this class")
				continue
			}

			// Persist message to PostgreSQL
			savedMsg, err := c.repo.SaveMessage(targetClassID, c.user.ID, p.Message)
			if err != nil {
				log.Printf("[WS ERROR] Failed to save chat message: %v", err)
				c.sendError("INTERNAL_ERROR", "Failed to save chat message")
				continue
			}

			// Broadcast to class room
			roomKey := fmt.Sprintf("class:%d", targetClassID)
			broadcastEvent := OutgoingEvent{
				Type: "chat.message",
				Payload: map[string]interface{}{
					"id":       savedMsg.ID,
					"class_id": savedMsg.ClassID,
					"user": map[string]interface{}{
						"id":         savedMsg.UserID,
						"name":       savedMsg.UserName,
						"role":       savedMsg.UserRole,
						"avatar_url": savedMsg.UserAvatar,
					},
					"message":    savedMsg.Message,
					"created_at": savedMsg.CreatedAt,
				},
			}

			payloadBytes, _ := json.Marshal(broadcastEvent)
			c.hub.BroadcastToRoom(roomKey, payloadBytes)

		case "meeting.join":
			var p struct {
				MeetingID int64 `json:"meeting_id"`
				IsAudio   bool  `json:"is_audio"`
				IsVideo   bool  `json:"is_video"`
			}
			if err := json.Unmarshal(incoming.Payload, &p); err != nil || p.MeetingID <= 0 {
				c.sendError("BAD_REQUEST", "Valid meeting_id is required")
				continue
			}

			roomKey := fmt.Sprintf("meeting:%d", p.MeetingID)
			c.hub.JoinRoom(c, roomKey)

			// Get all existing peers in this room
			peers := c.hub.GetRoomPeers(roomKey)

			// Send current peers list to the newly joined user
			c.sendEvent("meeting.peers", map[string]interface{}{
				"meeting_id": p.MeetingID,
				"peers":      peers,
			})

			// Broadcast to other peers that a new user joined
			joinedBytes, _ := json.Marshal(OutgoingEvent{
				Type: "meeting.peer_joined",
				Payload: map[string]interface{}{
					"id":             c.user.ID,
					"name":           c.user.Name,
					"role":           c.user.Role,
					"isAudioEnabled": p.IsAudio,
					"isVideoEnabled": p.IsVideo,
					"isHost":         c.user.Role == "teacher" || c.user.Role == "admin",
				},
			})
			c.hub.BroadcastToRoom(roomKey, joinedBytes)

		case "meeting.media":
			var p struct {
				MeetingID int64 `json:"meeting_id"`
				IsAudio   bool  `json:"is_audio"`
				IsVideo   bool  `json:"is_video"`
			}
			if err := json.Unmarshal(incoming.Payload, &p); err == nil && p.MeetingID > 0 {
				roomKey := fmt.Sprintf("meeting:%d", p.MeetingID)
				mediaBytes, _ := json.Marshal(OutgoingEvent{
					Type: "meeting.peer_media",
					Payload: map[string]interface{}{
						"user_id":        c.user.ID,
						"isAudioEnabled": p.IsAudio,
						"isVideoEnabled": p.IsVideo,
					},
				})
				c.hub.BroadcastToRoom(roomKey, mediaBytes)
			}

		case "meeting.leave":
			var p struct {
				MeetingID int64 `json:"meeting_id"`
			}
			if err := json.Unmarshal(incoming.Payload, &p); err == nil && p.MeetingID > 0 {
				roomKey := fmt.Sprintf("meeting:%d", p.MeetingID)
				leaveBytes, _ := json.Marshal(OutgoingEvent{
					Type: "meeting.peer_left",
					Payload: map[string]interface{}{
						"user_id": c.user.ID,
					},
				})
				c.hub.BroadcastToRoom(roomKey, leaveBytes)
			}

		default:
			c.sendError("UNKNOWN_EVENT", fmt.Sprintf("Unknown event type: %s", incoming.Type))
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Drain queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write([]byte{'\n'})
				_, _ = w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) sendEvent(eventType string, payload interface{}) {
	event := OutgoingEvent{
		Type:    eventType,
		Payload: payload,
	}
	data, err := json.Marshal(event)
	if err == nil {
		select {
		case c.send <- data:
		default:
		}
	}
}

func (c *Client) sendError(code, message string) {
	c.sendEvent("error", map[string]string{
		"code":    code,
		"message": message,
	})
}
