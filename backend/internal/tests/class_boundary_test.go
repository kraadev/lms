package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClassAccessBoundaryIsolation(t *testing.T) {
	router, db, _ := setupTestServer(t)
	defer db.Close()

	// 1. Student logs in
	tokenStudent := loginUser(t, router, "student1@lms.local", "password123")

	// 2. Student attempts to access materials of a class they haven't joined (Class ID 9999)
	req := httptest.NewRequest("GET", "/api/classes/9999/materials", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStudent)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound && w.Code != http.StatusForbidden {
		t.Errorf("expected 403 or 404 for unauthorized class access, got %d", w.Code)
	}

	// 3. Student attempts to grade assignment (must be rejected)
	reqGrade := httptest.NewRequest("POST", "/api/assignments/1/grade", nil)
	reqGrade.Header.Set("Authorization", "Bearer "+tokenStudent)
	wGrade := httptest.NewRecorder()
	router.ServeHTTP(wGrade, reqGrade)

	if wGrade.Code != http.StatusForbidden && wGrade.Code != http.StatusNotFound {
		t.Errorf("expected forbidden status when student attempts to grade, got %d", wGrade.Code)
	}
}
