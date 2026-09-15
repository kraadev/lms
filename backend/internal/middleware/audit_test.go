package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuditLogger(t *testing.T) {
	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	})

	wrapped := AuditLogger(handler)

	req := httptest.NewRequest("POST", "/api/classes", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if !called {
		t.Errorf("expected handler to be invoked")
	}
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}
