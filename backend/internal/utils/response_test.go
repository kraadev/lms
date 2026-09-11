package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONResponse(t *testing.T) {
	w := httptest.NewRecorder()
	payload := map[string]string{"message": "success"}

	JSON(w, http.StatusOK, payload)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success true, got %v", resp.Success)
	}
}

func TestErrorResponseRFC7807(t *testing.T) {
	w := httptest.NewRecorder()
	details := []ErrorDetail{
		{Field: "due_date", Issue: "Batas waktu tidak boleh di masa lalu"},
	}

	ErrorWithDetails(w, http.StatusUnprocessableEntity, "INVALID_ARGUMENT", "Validasi gagal", details)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", w.Code)
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Success {
		t.Errorf("expected success false, got true")
	}

	if resp.Error == nil {
		t.Fatalf("expected error object, got nil")
	}

	if resp.Error.Code != "INVALID_ARGUMENT" {
		t.Errorf("expected code INVALID_ARGUMENT, got %s", resp.Error.Code)
	}

	if resp.Error.Timestamp == "" {
		t.Errorf("expected non-empty RFC-7807 timestamp")
	}

	if len(resp.Error.Details) != 1 || resp.Error.Details[0].Field != "due_date" {
		t.Errorf("expected 1 detail for field due_date, got %+v", resp.Error.Details)
	}
}

func TestStandardErrorHelpers(t *testing.T) {
	tests := []struct {
		name         string
		fn           func(w http.ResponseWriter)
		expectedCode int
		expectedErr  string
	}{
		{
			name:         "BadRequest",
			fn:           func(w http.ResponseWriter) { BadRequest(w, "invalid id") },
			expectedCode: http.StatusBadRequest,
			expectedErr:  "BAD_REQUEST",
		},
		{
			name:         "Unauthorized",
			fn:           func(w http.ResponseWriter) { Unauthorized(w, "") },
			expectedCode: http.StatusUnauthorized,
			expectedErr:  "UNAUTHENTICATED",
		},
		{
			name:         "Forbidden",
			fn:           func(w http.ResponseWriter) { Forbidden(w, "") },
			expectedCode: http.StatusForbidden,
			expectedErr:  "PERMISSION_DENIED",
		},
		{
			name:         "NotFound",
			fn:           func(w http.ResponseWriter) { NotFound(w, "") },
			expectedCode: http.StatusNotFound,
			expectedErr:  "NOT_FOUND",
		},
		{
			name:         "Conflict",
			fn:           func(w http.ResponseWriter) { Conflict(w, "") },
			expectedCode: http.StatusConflict,
			expectedErr:  "RESOURCE_CONFLICT",
		},
		{
			name:         "RateLimitExceeded",
			fn:           func(w http.ResponseWriter) { RateLimitExceeded(w, "") },
			expectedCode: http.StatusTooManyRequests,
			expectedErr:  "RATE_LIMIT_EXCEEDED",
		},
		{
			name:         "InternalServerError",
			fn:           func(w http.ResponseWriter) { InternalServerError(w, "") },
			expectedCode: http.StatusInternalServerError,
			expectedErr:  "INTERNAL_SERVER_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.fn(w)

			if w.Code != tt.expectedCode {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedCode, w.Code)
			}

			var resp APIResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("%s: failed to decode response: %v", tt.name, err)
			}

			if resp.Error == nil || resp.Error.Code != tt.expectedErr {
				t.Errorf("%s: expected error code %s, got %+v", tt.name, tt.expectedErr, resp.Error)
			}
		})
	}
}
