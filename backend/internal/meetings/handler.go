package meetings

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"lms/internal/middleware"
	"lms/internal/utils"
)

type Handler struct {
	service          *Service
	accessController *middleware.AccessController
}

func NewHandler(service *Service, accessController *middleware.AccessController) *Handler {
	return &Handler{
		service:          service,
		accessController: accessController,
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

	utils.JSON(w, http.StatusOK, res)
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
		utils.Forbidden(w, "Only the assigned teacher or an admin can end this meeting.")
		return
	}

	if err := h.service.EndMeeting(id); err != nil {
		utils.InternalServerError(w, "Failed to end meeting")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Meeting ended successfully"})
}
