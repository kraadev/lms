package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"lms/internal/models"
)

func TestEndToEndAcceptanceScenario(t *testing.T) {
	router, db, _ := setupTestServer(t)
	defer db.Close()

	// ==========================================
	// 1. ADMIN FLOW
	// ==========================================
	// Admin logs in
	tokenAdmin := loginUser(t, router, "admin@lms.local", "admin123")

	// Admin creates a new user
	newUserPayload, _ := json.Marshal(map[string]interface{}{
		"name":     "Eko Kurniawan",
		"email":    "eko@lms.local",
		"password": "password123",
		"role":     "student",
	})
	req := httptest.NewRequest("POST", "/api/users", bytes.NewReader(newUserPayload))
	req.Header.Set("Authorization", "Bearer "+tokenAdmin)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Admin failed to create user: %d, body: %s", w.Code, w.Body.String())
	}

	// ==========================================
	// 2. TEACHER A FLOW
	// ==========================================
	// Teacher A logs in
	tokenTeacherA := loginUser(t, router, "teacher1@lms.local", "password123")

	// Teacher A creates Class Alpha
	createClassPayload, _ := json.Marshal(map[string]interface{}{
		"name":         "Kelas Algoritma & Struktur Data",
		"description":  "Belajar algoritma tingkat lanjut.",
		"academic_year": "2026/2027",
	})
	req = httptest.NewRequest("POST", "/api/classes", bytes.NewReader(createClassPayload))
	req.Header.Set("Authorization", "Bearer "+tokenTeacherA)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Teacher A failed to create class: %d, body: %s", w.Code, w.Body.String())
	}

	var classRes struct {
		Data models.Class `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &classRes)
	classAID := classRes.Data.ID

	// Teacher A adds Student A (ID 4 - student1@lms.local) to Class Alpha
	addStudentPayload, _ := json.Marshal(map[string]int64{
		"user_id": 4,
	})
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/classes/%d/members", classAID), bytes.NewReader(addStudentPayload))
	req.Header.Set("Authorization", "Bearer "+tokenTeacherA)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Teacher A failed to add student: %d, body: %s", w.Code, w.Body.String())
	}

	// Teacher A creates Assignment in Class Alpha
	assignmentPayload, _ := json.Marshal(map[string]interface{}{
		"title":       "Tugas 1: Binary Search Tree",
		"description": "Implementasikan tree traversal di Go.",
		"deadline":    time.Now().Add(48 * time.Hour).Format(time.RFC3339),
		"max_score":   100.0,
	})
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/classes/%d/assignments", classAID), bytes.NewReader(assignmentPayload))
	req.Header.Set("Authorization", "Bearer "+tokenTeacherA)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Teacher A failed to create assignment: %d, body: %s", w.Code, w.Body.String())
	}

	var assignRes struct {
		Data models.Assignment `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &assignRes)
	assignmentAID := assignRes.Data.ID

	// Teacher A creates Quiz in Class Alpha
	quizPayload, _ := json.Marshal(map[string]interface{}{
		"title":            "Kuis 1: Tree & Graph",
		"description":      "Uji pemahaman tree.",
		"duration_minutes": 45,
		"start_at":         time.Now().Add(-10 * time.Minute).Format(time.RFC3339),
		"end_at":           time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"max_attempts":     2,
		"questions": []map[string]interface{}{
			{
				"question":    "Berapa kompleksitas waktu pencarian rata-rata di BST?",
				"type":        "multiple_choice",
				"points":      20.0,
				"order_index": 1,
				"options": []map[string]interface{}{
					{"option_text": "O(log n)", "is_correct": true},
					{"option_text": "O(n^2)", "is_correct": false},
					{"option_text": "O(1)", "is_correct": false},
				},
			},
		},
	})
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/classes/%d/quizzes", classAID), bytes.NewReader(quizPayload))
	req.Header.Set("Authorization", "Bearer "+tokenTeacherA)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Teacher A failed to create quiz: %d, body: %s", w.Code, w.Body.String())
	}

	var quizCreatedRes struct {
		Data models.Quiz `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &quizCreatedRes)
	quizAID := quizCreatedRes.Data.ID

	// Teacher A starts meeting in Class Alpha
	meetingPayload, _ := json.Marshal(map[string]string{
		"title": "Live Diskusi BST",
		"type":  "video",
	})
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/classes/%d/meetings", classAID), bytes.NewReader(meetingPayload))
	req.Header.Set("Authorization", "Bearer "+tokenTeacherA)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Teacher A failed to start meeting: %d, body: %s", w.Code, w.Body.String())
	}

	var meetingCreatedRes struct {
		Data models.Meeting `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &meetingCreatedRes)
	meetingAID := meetingCreatedRes.Data.ID

	// ==========================================
	// 3. STUDENT A FLOW (Enrolled in Class Alpha)
	// ==========================================
	tokenStudentA := loginUser(t, router, "student1@lms.local", "password123")

	// Student A views Class Alpha
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/classes/%d", classAID), nil)
	req.Header.Set("Authorization", "Bearer "+tokenStudentA)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Student A should be able to view Class Alpha, got %d", w.Code)
	}

	// Student A submits assignment
	ansTxt := "Penyelesaian BST berhasil diimplementasikan."
	subPayload, _ := json.Marshal(map[string]*string{
		"text_answer": &ansTxt,
	})
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/assignments/%d/submissions", assignmentAID), bytes.NewReader(subPayload))
	req.Header.Set("Authorization", "Bearer "+tokenStudentA)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Student A should submit assignment, got %d", w.Code)
	}

	// Student A starts quiz attempt
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/quizzes/%d/attempts", quizAID), nil)
	req.Header.Set("Authorization", "Bearer "+tokenStudentA)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Student A should start quiz attempt, got %d (body: %s)", w.Code, w.Body.String())
	}

	var attemptRes struct {
		Data models.QuizAttempt `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &attemptRes)
	attemptID := attemptRes.Data.ID

	// Student A submits correct answer for Q1
	q1 := quizCreatedRes.Data.Questions[0]
	var correctOptionID int64
	for _, opt := range q1.Options {
		if opt.IsCorrect != nil && *opt.IsCorrect {
			correctOptionID = opt.ID
		}
	}

	submitQuizPayload, _ := json.Marshal(map[string]interface{}{
		"answers": []map[string]interface{}{
			{
				"question_id":        q1.ID,
				"selected_option_id": correctOptionID,
			},
		},
	})
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/quizzes/attempts/%d/submit", attemptID), bytes.NewReader(submitQuizPayload))
	req.Header.Set("Authorization", "Bearer "+tokenStudentA)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Student A submit quiz failed: %d, body: %s", w.Code, w.Body.String())
	}

	var gradedRes struct {
		Data models.QuizAttempt `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &gradedRes)
	if gradedRes.Data.Score == nil || *gradedRes.Data.Score != 20.0 {
		t.Fatalf("Expected score 20.0, got: %v", gradedRes.Data.Score)
	}

	// Student A joins Class Alpha meeting
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/meetings/%d/join", meetingAID), nil)
	req.Header.Set("Authorization", "Bearer "+tokenStudentA)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Student A should join meeting, got %d", w.Code)
	}

	// ==========================================
	// 4. STUDENT B FLOW (Belongs only to Class 2, NOT Class Alpha)
	// ==========================================
	tokenStudentB := loginUser(t, router, "student3@lms.local", "password123")

	// 4.1 Student B tries to view Class Alpha -> FORBIDDEN
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/classes/%d", classAID), nil)
	req.Header.Set("Authorization", "Bearer "+tokenStudentB)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 when Student B accesses Class Alpha, got %d", w.Code)
	}

	// 4.2 Student B tries to access Class Alpha assignments -> FORBIDDEN
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/classes/%d/assignments", classAID), nil)
	req.Header.Set("Authorization", "Bearer "+tokenStudentB)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 when Student B accesses Class Alpha assignments, got %d", w.Code)
	}

	// 4.3 Student B tries to submit to Class Alpha assignment -> FORBIDDEN
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/assignments/%d/submissions", assignmentAID), bytes.NewReader(subPayload))
	req.Header.Set("Authorization", "Bearer "+tokenStudentB)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 when Student B submits to Class Alpha assignment, got %d", w.Code)
	}

	// 4.4 Student B tries to view or start Class Alpha quiz -> FORBIDDEN
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/quizzes/%d", quizAID), nil)
	req.Header.Set("Authorization", "Bearer "+tokenStudentB)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 when Student B accesses Class Alpha quiz, got %d", w.Code)
	}

	req = httptest.NewRequest("POST", fmt.Sprintf("/api/quizzes/%d/attempts", quizAID), nil)
	req.Header.Set("Authorization", "Bearer "+tokenStudentB)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 when Student B attempts Class Alpha quiz, got %d", w.Code)
	}

	// 4.5 Student B tries to join Class Alpha meeting -> FORBIDDEN
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/meetings/%d/join", meetingAID), nil)
	req.Header.Set("Authorization", "Bearer "+tokenStudentB)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 when Student B joins Class Alpha meeting, got %d", w.Code)
	}

	t.Log("End-to-end acceptance test passed with 100% security & authorization compliance!")
}
