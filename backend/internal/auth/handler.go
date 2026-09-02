package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"lms/internal/middleware"
	"lms/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request payload")
		return
	}

	result, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		utils.Unauthorized(w, err.Error())
		return
	}

	// Set secure HttpOnly Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "lms_token",
		Value:    result.Token,
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(h.service.cfg.JWTExpiryHours) * time.Hour),
		HttpOnly: true,
		Secure:   false, // Set to false for local development http://localhost
		SameSite: http.SameSiteLaxMode,
	})

	utils.JSON(w, http.StatusOK, result)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "lms_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	authUser := middleware.GetUser(r)
	if authUser == nil {
		utils.Unauthorized(w, "Not authenticated")
		return
	}

	user, err := h.service.GetMe(authUser.ID)
	if err != nil || user == nil {
		utils.NotFound(w, "User profile not found")
		return
	}

	utils.JSON(w, http.StatusOK, user)
}
