package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"lms/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	search := r.URL.Query().Get("search")

	users, err := h.service.ListUsers(role, search)
	if err != nil {
		utils.InternalServerError(w, "Failed to fetch users")
		return
	}

	utils.JSON(w, http.StatusOK, users)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid JSON payload")
		return
	}

	user, err := h.service.CreateUser(req)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, user)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(w, "Invalid user ID")
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		utils.InternalServerError(w, "Error finding user")
		return
	}
	if user == nil {
		utils.NotFound(w, "User not found")
		return
	}

	utils.JSON(w, http.StatusOK, user)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(w, "Invalid user ID")
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		utils.InternalServerError(w, "Failed to delete user")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
