package chat

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	"lms/internal/config"
	"lms/internal/middleware"
	"lms/internal/utils"
)

type Handler struct {
	hub              *Hub
	repo             *Repository
	accessController *middleware.AccessController
	cfg              *config.Config
	upgrader         websocket.Upgrader
}

func NewHandler(hub *Hub, repo *Repository, accessController *middleware.AccessController, cfg *config.Config) *Handler {
	return &Handler{
		hub:              hub,
		repo:             repo,
		accessController: accessController,
		cfg:              cfg,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow localhost / local frontend origins
				return true
			},
		},
	}
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil {
		utils.Unauthorized(w, "Authentication required for WebSocket connection")
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS ERROR] Upgrade failed: %v", err)
		return
	}

	client := &Client{
		hub:              h.hub,
		conn:             conn,
		send:             make(chan []byte, 256),
		user:             user,
		repo:             h.repo,
		accessController: h.accessController,
	}

	h.hub.register <- client

	// Start reader & writer routines
	go client.writePump()
	go client.readPump()
}

func (h *Handler) ListMessages(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	classID, _ := strconv.ParseInt(chi.URLParam(r, "classId"), 10, 64)

	hasAccess, err := h.accessController.CheckClassAccess(user.ID, user.Role, classID)
	if err != nil || !hasAccess {
		utils.Forbidden(w, "You are not enrolled in this class.")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	messages, err := h.repo.ListMessages(classID, limit, offset)
	if err != nil {
		utils.InternalServerError(w, "Failed to load chat history")
		return
	}

	utils.JSON(w, http.StatusOK, messages)
}
