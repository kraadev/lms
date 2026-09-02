package meetings

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"

	"lms/internal/chat"
	"lms/internal/middleware"
	"lms/internal/models"
	"lms/internal/utils"
)

type Handler struct {
	service          *Service
	accessController *middleware.AccessController
	hub              *chat.Hub
	activeMu         sync.RWMutex
	activePeers      map[int64]map[int64]map[string]interface{}
}

func NewHandler(service *Service, accessController *middleware.AccessController, hub *chat.Hub) *Handler {
	return &Handler{
		service:          service,
		accessController: accessController,
		hub:              hub,
		activePeers:      make(map[int64]map[int64]map[string]interface{}),
	}
}

func (h *Handler) ListByClass(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	classID, _ := strconv.ParseInt(chi.URLParam(r, "classId"), 10, 64)

	hasAccess, err := h.accessController.CheckClassAccess(user.ID, user.Role, classID)
	if err != nil || !hasAccess {
		utils.Forbidden(w, "You do not have access to this class's meetings.")
		return
	}

	meetings, err := h.service.ListByClass(classID)
	if err != nil {
		utils.InternalServerError(w, "Failed to load meetings")
		return
	}

	utils.JSON(w, http.StatusOK, meetings)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	classID, _ := strconv.ParseInt(chi.URLParam(r, "classId"), 10, 64)

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, classID)
	if err != nil || !canManage {
		utils.Forbidden(w, "Only the assigned teacher or an admin can start class meetings.")
		return
	}

	var req CreateMeetingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid JSON payload")
		return
	}

	req.ClassID = classID
	req.TeacherID = user.ID

	meeting, err := h.service.CreateMeeting(req)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, meeting)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	meeting, err := h.service.GetByID(id)
	if err != nil || meeting == nil {
		utils.NotFound(w, "Meeting not found")
		return
	}

	hasAccess, err := h.accessController.CheckClassAccess(user.ID, user.Role, meeting.ClassID)
	if err != nil || !hasAccess {
		utils.Forbidden(w, "Access denied to this meeting")
		return
	}

	utils.JSON(w, http.StatusOK, meeting)
}

func (h *Handler) Join(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	meeting, err := h.service.GetByID(id)
	if err != nil || meeting == nil {
		utils.NotFound(w, "Meeting not found")
		return
	}

	hasAccess, err := h.accessController.CheckClassAccess(user.ID, user.Role, meeting.ClassID)
	if err != nil || !hasAccess {
		utils.Forbidden(w, "You are not enrolled in the class for this meeting.")
		return
	}

	res, err := h.service.JoinMeeting(id, user)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	// Register peer in in-memory active list
	peerInfo := map[string]interface{}{
		"id":             user.ID,
		"name":           user.Name,
		"role":           user.Role,
		"isHost":         user.Role == models.RoleTeacher || user.Role == models.RoleAdmin || meeting.TeacherID == user.ID,
		"isAudioEnabled": true,
		"isVideoEnabled": false,
	}

	h.activeMu.Lock()
	if _, ok := h.activePeers[id]; !ok {
		h.activePeers[id] = make(map[int64]map[string]interface{})
	}
	h.activePeers[id][user.ID] = peerInfo
	h.activeMu.Unlock()

	// Broadcast peer presence to WebSocket room
	if h.hub != nil {
		roomKey := fmt.Sprintf("meeting:%d", id)
		joinedBytes, _ := json.Marshal(chat.OutgoingEvent{
			Type:    "meeting.peer_joined",
			Payload: peerInfo,
		})
		h.hub.BroadcastToRoom(roomKey, joinedBytes)
	}

	utils.JSON(w, http.StatusOK, res)
}

func (h *Handler) GetParticipants(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	h.activeMu.RLock()
	peersMap, ok := h.activePeers[id]
	var list []map[string]interface{}
	if ok {
		for _, p := range peersMap {
			list = append(list, p)
		}
	} else {
		list = []map[string]interface{}{}
	}
	h.activeMu.RUnlock()

	utils.JSON(w, http.StatusOK, list)
}

func (h *Handler) Leave(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	h.activeMu.Lock()
	if peersMap, ok := h.activePeers[id]; ok {
		delete(peersMap, user.ID)
		if len(peersMap) == 0 {
			delete(h.activePeers, id)
		}
	}
	h.activeMu.Unlock()

	// Broadcast peer left
	if h.hub != nil {
		roomKey := fmt.Sprintf("meeting:%d", id)
		leaveBytes, _ := json.Marshal(chat.OutgoingEvent{
			Type: "meeting.peer_left",
			Payload: map[string]interface{}{
				"user_id": user.ID,
			},
		})
		h.hub.BroadcastToRoom(roomKey, leaveBytes)
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Left meeting"})
}

func (h *Handler) End(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	meeting, err := h.service.GetByID(id)
	if err != nil || meeting == nil {
		utils.NotFound(w, "Meeting not found")
		return
	}

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, meeting.ClassID)
	if err != nil || !canManage {
		if meeting.TeacherID != user.ID && user.Role != models.RoleTeacher && user.Role != models.RoleAdmin {
			utils.Forbidden(w, "Only the assigned teacher, host, or an admin can end this meeting.")
			return
		}
	}

	if err := h.service.EndMeeting(id); err != nil {
		utils.InternalServerError(w, "Failed to end meeting")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Meeting ended successfully"})
}
