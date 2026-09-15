package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := SecurityHeaders(dummyHandler)

	req := httptest.NewRequest("GET", "/any", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	expectedHeaders := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"X-XSS-Protection":       "1; mode=block",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
	}

	for key, expected := range expectedHeaders {
		got := w.Header().Get(key)
		if got != expected {
			t.Errorf("header %s: expected '%s', got '%s'", key, expected, got)
		}
	}
}
