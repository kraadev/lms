package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"lms/internal/config"
	"lms/internal/database"
)

func TestHealthProbes(t *testing.T) {
	cfg := &config.Config{
		DBDriver:     "sqlite",
		DBSQLitePath: ":memory:",
	}
	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("failed to connect in-memory db: %v", err)
	}
	defer db.Close()

	handler := NewHandler(db)

	// Test /healthz
	reqHealthz := httptest.NewRequest("GET", "/healthz", nil)
	wHealthz := httptest.NewRecorder()
	handler.Healthz(wHealthz, reqHealthz)
	if wHealthz.Code != http.StatusOK {
		t.Errorf("expected 200 for healthz, got %d", wHealthz.Code)
	}

	// Test /readyz
	reqReadyz := httptest.NewRequest("GET", "/readyz", nil)
	wReadyz := httptest.NewRecorder()
	handler.Readyz(wReadyz, reqReadyz)
	if wReadyz.Code != http.StatusOK {
		t.Errorf("expected 200 for readyz, got %d", wReadyz.Code)
	}
}
