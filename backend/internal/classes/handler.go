package classes

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

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	classes, err := h.service.ListUserClasses(user.ID, user.Role)
	if err != nil {
		utils.InternalServerError(w, "Failed to list classes")
		return
	}

	utils.JSON(w, http.StatusOK, classes)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	if user == nil {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	var req CreateClassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid JSON payload")
		return
	}

	class, err := h.service.CreateClass(req, user.ID, user.Role)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, class)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	classID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	hasAccess, err := h.accessController.CheckClassAccess(user.ID, user.Role, classID)
	if err != nil || !hasAccess {
		utils.Forbidden(w, "You are not enrolled in this class.")
		return
	}

	class, err := h.service.GetClassByID(classID)
	if err != nil {
		utils.InternalServerError(w, "Failed to load class")
		return
	}
	if class == nil {
		utils.NotFound(w, "Class not found")
		return
	}

	utils.JSON(w, http.StatusOK, class)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	classID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, classID)
	if err != nil || !canManage {
		utils.Forbidden(w, "Only the assigned teacher or an admin can update this class.")
		return
	}

	existing, err := h.service.GetClassByID(classID)
	if err != nil || existing == nil {
		utils.NotFound(w, "Class not found")
		return
	}

	var req struct {
		Name         *string `json:"name"`
		Description  *string `json:"description"`
		AcademicYear *string `json:"academic_year"`
		Status       *string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid JSON payload")
		return
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.AcademicYear != nil {
		existing.AcademicYear = *req.AcademicYear
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}

	if err := h.service.UpdateClass(existing); err != nil {
		utils.InternalServerError(w, "Failed to update class")
		return
	}

	utils.JSON(w, http.StatusOK, existing)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	classID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, classID)
	if err != nil || !canManage {
		utils.Forbidden(w, "Only the assigned teacher or an admin can delete this class.")
		return
	}

	if err := h.service.DeleteClass(classID); err != nil {
		utils.InternalServerError(w, "Failed to delete class")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Class deleted successfully"})
}

func (h *Handler) ListMembers(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	classID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	hasAccess, err := h.accessController.CheckClassAccess(user.ID, user.Role, classID)
	if err != nil || !hasAccess {
		utils.Forbidden(w, "You are not enrolled in this class.")
		return
	}

	members, err := h.service.ListClassMembers(classID)
	if err != nil {
		utils.InternalServerError(w, "Failed to list members")
		return
	}

	utils.JSON(w, http.StatusOK, members)
}

func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	classID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, classID)
	if err != nil || !canManage {
		utils.Forbidden(w, "Only the assigned teacher or an admin can add members.")
		return
	}

	var req struct {
		UserID int64 `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID <= 0 {
		utils.BadRequest(w, "Valid user_id is required")
		return
	}

	if err := h.service.AddStudentToClass(classID, req.UserID); err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, map[string]string{"message": "Member added successfully"})
}

func (h *Handler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	classID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	targetUserID, _ := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, classID)
	if err != nil || !canManage {
		utils.Forbidden(w, "Only the assigned teacher or an admin can remove members.")
		return
	}

	if err := h.service.RemoveStudentFromClass(classID, targetUserID); err != nil {
		utils.InternalServerError(w, "Failed to remove member")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Member removed successfully"})
}
