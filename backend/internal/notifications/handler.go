package notifications

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"lms/internal/middleware"
	"lms/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	notifs, err := h.service.ListUserNotifications(user.ID)
	if err != nil {
		utils.InternalServerError(w, "Failed to load notifications")
		return
	}

	unreadCount := 0
	for _, n := range notifs {
		if !n.IsRead {
			unreadCount++
		}
	}

	resp := map[string]interface{}{
		"notifications": notifs,
		"unread_count":  unreadCount,
	}

	utils.JSON(w, http.StatusOK, resp)
}

func (h *Handler) MarkRead(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err := h.service.MarkAsRead(id, user.ID); err != nil {
		utils.InternalServerError(w, "Failed to mark notification as read")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Notification marked as read"})
}

func (h *Handler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	if err := h.service.MarkAllRead(user.ID); err != nil {
		utils.InternalServerError(w, "Failed to mark notifications as read")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "All notifications marked as read"})
}
