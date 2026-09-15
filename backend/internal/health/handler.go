package health

import (
	"context"
	"net/http"
	"time"

	"lms/internal/database"
	"lms/internal/utils"
)

type Handler struct {
	db        *database.DB
	startTime time.Time
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{
		db:        db,
		startTime: time.Now(),
	}
}

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"status":    "pass",
		"uptime":    time.Since(h.startTime).String(),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if err := h.db.DB.PingContext(ctx); err != nil {
		utils.Error(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Database connection ping failed")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]interface{}{
		"status":    "pass",
		"database":  "healthy",
		"driver":    h.db.Driver,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
