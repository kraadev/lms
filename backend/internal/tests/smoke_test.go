package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlatformCoreHealthSmoke(t *testing.T) {
	router, db, _ := setupTestServer(t)
	defer db.Close()

	// 1. Unauthenticated request to protected endpoint must return 401
	reqUnauth := httptest.NewRequest("GET", "/api/auth/me", nil)
	wUnauth := httptest.NewRecorder()
	router.ServeHTTP(wUnauth, reqUnauth)
	if wUnauth.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for missing token, got %d", wUnauth.Code)
	}

	// 2. User login validation smoke
	tokenTeacher := loginUser(t, router, "teacher1@lms.local", "password123")
	if tokenTeacher == "" {
		t.Fatalf("teacher login token generation failed")
	}

	// 3. Authenticated profile endpoint smoke
	reqMe := httptest.NewRequest("GET", "/api/auth/me", nil)
	reqMe.Header.Set("Authorization", "Bearer "+tokenTeacher)
	wMe := httptest.NewRecorder()
	router.ServeHTTP(wMe, reqMe)
	if wMe.Code != http.StatusOK {
		t.Fatalf("me profile inspection failed with code %d", wMe.Code)
	}

	// 4. Authenticated classes list smoke
	reqClasses := httptest.NewRequest("GET", "/api/classes", nil)
	reqClasses.Header.Set("Authorization", "Bearer "+tokenTeacher)
	wClasses := httptest.NewRecorder()
	router.ServeHTTP(wClasses, reqClasses)
	if wClasses.Code != http.StatusOK {
		t.Fatalf("classes list inspection failed with code %d", wClasses.Code)
	}
}
