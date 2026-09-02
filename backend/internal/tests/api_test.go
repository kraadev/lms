package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"lms/internal/assignments"
	"lms/internal/auth"
	"lms/internal/chat"
	"lms/internal/classes"
	"lms/internal/config"
	"lms/internal/database"
	"lms/internal/materials"
	"lms/internal/meetings"
	"lms/internal/middleware"
	"lms/internal/notifications"
	"lms/internal/quizzes"
	"lms/internal/storage"
	"lms/internal/users"
)

func setupTestServer(t *testing.T) (http.Handler, *database.DB, *config.Config) {
	tempDir := t.TempDir()
	cfg := &config.Config{
		Port:           "8080",
		Env:            "testing",
		DBDriver:       "sqlite",
		DBSQLitePath:   tempDir + "/test_lms.db",
		JWTSecret:      "test-secret-key-for-lms-jwt-testing-32chars",
		JWTExpiryHours: 2,
		LiveKitURL:     "ws://localhost:7880",
		LiveKitAPIKey:  "devkey",
		LiveKitSecret:  "secret",
		UploadDir:      tempDir + "/uploads",
	}

	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test db: %v", err)
	}

	if err := database.RunMigrations(db); err != nil {
		t.Fatalf("Failed to run test migrations: %v", err)
	}

	if err := database.SeedData(db); err != nil {
		t.Fatalf("Failed to seed test db: %v", err)
	}

	storageService, err := storage.NewStorageService(cfg.UploadDir)
	if err != nil {
		t.Fatalf("Failed to init test storage: %v", err)
	}

	userRepo := users.NewRepository(db)
	classRepo := classes.NewRepository(db)
	materialRepo := materials.NewRepository(db)
	assignmentRepo := assignments.NewRepository(db)
	quizRepo := quizzes.NewRepository(db)
	meetingRepo := meetings.NewRepository(db)
	chatRepo := chat.NewRepository(db)
	notifRepo := notifications.NewRepository(db)

	accessController := middleware.NewAccessController(db)
	authService := auth.NewService(userRepo, cfg)
	userService := users.NewService(userRepo)
	classService := classes.NewService(classRepo, userRepo)
	materialService := materials.NewService(materialRepo)
	assignmentService := assignments.NewService(assignmentRepo)
	quizService := quizzes.NewService(quizRepo)
	liveKitService := meetings.NewLiveKitService(cfg)
	meetingService := meetings.NewService(meetingRepo, liveKitService, cfg)
	notifService := notifications.NewService(notifRepo)

	chatHub := chat.NewHub()
	go chatHub.Run()

	authHandler := auth.NewHandler(authService)
	userHandler := users.NewHandler(userService)
	classHandler := classes.NewHandler(classService, accessController)
	materialHandler := materials.NewHandler(materialService, accessController, storageService)
	assignmentHandler := assignments.NewHandler(assignmentService, accessController, storageService)
	quizHandler := quizzes.NewHandler(quizService, accessController)
	meetingHandler := meetings.NewHandler(meetingService, accessController)
	chatHandler := chat.NewHandler(chatHub, chatRepo, accessController, cfg)
	notifHandler := notifications.NewHandler(notifService)

	r := chi.NewRouter()
	r.Use(middleware.AuthMiddleware(cfg))

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/logout", authHandler.Logout)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAuth)
			r.Get("/me", authHandler.Me)
		})
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.RequireAuth)

		r.Route("/users", func(r chi.Router) {
			r.Get("/{id}", userHandler.GetByID)
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAdmin)
				r.Get("/", userHandler.List)
				r.Post("/", userHandler.Create)
			})
		})

		r.Route("/classes", func(r chi.Router) {
			r.Get("/", classHandler.List)
			r.Post("/", classHandler.Create)
			r.Get("/{id}", classHandler.GetByID)
			r.Patch("/{id}", classHandler.Update)
			r.Delete("/{id}", classHandler.Delete)
			r.Get("/{id}/members", classHandler.ListMembers)
			r.Post("/{id}/members", classHandler.AddMember)

			r.Get("/{classId}/materials", materialHandler.ListByClass)
			r.Post("/{classId}/materials", materialHandler.Create)

			r.Get("/{classId}/assignments", assignmentHandler.ListByClass)
			r.Post("/{classId}/assignments", assignmentHandler.Create)

			r.Get("/{classId}/quizzes", quizHandler.ListByClass)
			r.Post("/{classId}/quizzes", quizHandler.Create)

			r.Get("/{classId}/meetings", meetingHandler.ListByClass)
			r.Post("/{classId}/meetings", meetingHandler.Create)
			r.Get("/{classId}/messages", chatHandler.ListMessages)
		})

		r.Route("/materials", func(r chi.Router) {
			r.Get("/{id}", materialHandler.GetByID)
		})

		r.Route("/assignments", func(r chi.Router) {
			r.Get("/{id}", assignmentHandler.GetByID)
			r.Post("/{id}/submissions", assignmentHandler.Submit)
			r.Get("/{id}/my-submission", assignmentHandler.GetMySubmission)
			r.Get("/{id}/submissions", assignmentHandler.ListSubmissions)
		})

		r.Patch("/submissions/{id}/grade", assignmentHandler.Grade)

		r.Route("/quizzes", func(r chi.Router) {
			r.Get("/{id}", quizHandler.GetByID)
			r.Post("/{id}/attempts", quizHandler.StartAttempt)
			r.Route("/attempts", func(r chi.Router) {
				r.Post("/{id}/submit", quizHandler.SubmitAttempt)
				r.Get("/{id}", quizHandler.GetAttempt)
			})
		})

		r.Route("/meetings", func(r chi.Router) {
			r.Get("/{id}", meetingHandler.GetByID)
			r.Post("/{id}/join", meetingHandler.Join)
			r.Post("/{id}/end", meetingHandler.End)
		})

		r.Route("/notifications", func(r chi.Router) {
			r.Get("/", notifHandler.List)
			r.Patch("/{id}/read", notifHandler.MarkRead)
		})
	})

	return r, db, cfg
}

func loginUser(t *testing.T, router http.Handler, email, password string) string {
	payload, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Login failed for %s: code %d, body: %s", email, w.Code, w.Body.String())
	}

	var res struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil || res.Data.Token == "" {
		t.Fatalf("Failed to parse login token: %v", err)
	}

	return res.Data.Token
}

func TestUnauthenticatedAccess(t *testing.T) {
	router, db, _ := setupTestServer(t)
	defer db.Close()

	req := httptest.NewRequest("GET", "/api/classes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized, got %d", w.Code)
	}
}

func TestCrossClassAccessControl(t *testing.T) {
	router, db, _ := setupTestServer(t)
	defer db.Close()

	// Student 1 belongs to Class 1 (ID 1)
	// Student 3 belongs to Class 2 (ID 2)
	tokenStudent3 := loginUser(t, router, "student3@lms.local", "password123")

	// Student 3 attempts to access Class 1 materials
	req := httptest.NewRequest("GET", "/api/classes/1/materials", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStudent3)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 Forbidden when Student 3 accesses Class 1 materials, got %d (body: %s)", w.Code, w.Body.String())
	}

	// Student 3 attempts to access Class 1 assignments
	req = httptest.NewRequest("GET", "/api/classes/1/assignments", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStudent3)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 Forbidden when Student 3 accesses Class 1 assignments, got %d", w.Code)
	}
}

func TestTeacherIsolation(t *testing.T) {
	router, db, _ := setupTestServer(t)
	defer db.Close()

	// Teacher 1 owns Class 1
	// Teacher 2 owns Class 2
	tokenTeacher2 := loginUser(t, router, "teacher2@lms.local", "password123")

	// Teacher 2 attempts to modify Class 1
	updatePayload, _ := json.Marshal(map[string]string{
		"name": "Hacked Class Name",
	})
	req := httptest.NewRequest("PATCH", "/api/classes/1", bytes.NewReader(updatePayload))
	req.Header.Set("Authorization", "Bearer "+tokenTeacher2)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 Forbidden when Teacher 2 updates Teacher 1's class, got %d", w.Code)
	}
}

func TestQuizSecurityAndNoLeakage(t *testing.T) {
	router, db, _ := setupTestServer(t)
	defer db.Close()

	tokenStudent1 := loginUser(t, router, "student1@lms.local", "password123")

	// Student 1 views Quiz 1 (in Class 1)
	req := httptest.NewRequest("GET", "/api/quizzes/1", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStudent1)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for enrolled student, got %d", w.Code)
	}

	// Verify is_correct is NOT leaked to student
	var res struct {
		Success bool `json:"success"`
		Data    struct {
			Questions []struct {
				Options []struct {
					IsCorrect *bool `json:"is_correct"`
				} `json:"options"`
			} `json:"questions"`
		} `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to parse quiz response: %v", err)
	}

	for _, q := range res.Data.Questions {
		for _, opt := range q.Options {
			if opt.IsCorrect != nil {
				t.Fatalf("CRITICAL SECURITY FLAW: is_correct leaked to student: %v", *opt.IsCorrect)
			}
		}
	}
}

func TestMeetingLifecycleAndEndedRestriction(t *testing.T) {
	router, db, _ := setupTestServer(t)
	defer db.Close()

	tokenTeacher1 := loginUser(t, router, "teacher1@lms.local", "password123")
	tokenStudent1 := loginUser(t, router, "student1@lms.local", "password123")
	tokenStudent3 := loginUser(t, router, "student3@lms.local", "password123") // Not in Class 1

	// Teacher 1 creates meeting in Class 1
	meetPayload, _ := json.Marshal(map[string]string{
		"title": "Live Session Concurrency",
		"type":  "video",
	})
	req := httptest.NewRequest("POST", "/api/classes/1/meetings", bytes.NewReader(meetPayload))
	req.Header.Set("Authorization", "Bearer "+tokenTeacher1)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create meeting: %d (body: %s)", w.Code, w.Body.String())
	}

	var meetRes struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &meetRes)
	meetingID := meetRes.Data.ID

	// Student 3 (wrong class) tries to join Class 1 meeting -> must fail 403
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/meetings/%d/join", meetingID), nil)
	req.Header.Set("Authorization", "Bearer "+tokenStudent3)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 when Student 3 joins Class 1 meeting, got %d", w.Code)
	}

	// Student 1 (enrolled in Class 1) joins Class 1 meeting -> must succeed 200 with LiveKit token
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/meetings/%d/join", meetingID), nil)
	req.Header.Set("Authorization", "Bearer "+tokenStudent1)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 when Student 1 joins meeting, got %d", w.Code)
	}

	// Teacher 1 ends the meeting
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/meetings/%d/end", meetingID), nil)
	req.Header.Set("Authorization", "Bearer "+tokenTeacher1)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 when Teacher ends meeting, got %d", w.Code)
	}

	// Student 1 tries to join ENDED meeting -> must be rejected 400 Bad Request
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/meetings/%d/join", meetingID), nil)
	req.Header.Set("Authorization", "Bearer "+tokenStudent1)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 when joining ended meeting, got %d", w.Code)
	}
}

func TestAssignmentSubmissionAndGrading(t *testing.T) {
	router, db, _ := setupTestServer(t)
	defer db.Close()

	tokenStudent1 := loginUser(t, router, "student1@lms.local", "password123")
	tokenTeacher1 := loginUser(t, router, "teacher1@lms.local", "password123")

	// Student 1 submits text answer for Assignment 1
	txt := "Goroutine adalah lightweight execution thread yang dikelola oleh Go runtime."
	subPayload, _ := json.Marshal(map[string]*string{
		"text_answer": &txt,
	})
	req := httptest.NewRequest("POST", "/api/assignments/1/submissions", bytes.NewReader(subPayload))
	req.Header.Set("Authorization", "Bearer "+tokenStudent1)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Failed to submit assignment: %d (body: %s)", w.Code, w.Body.String())
	}

	// Teacher 1 grades submission
	gradePayload, _ := json.Marshal(map[string]interface{}{
		"score":    95.0,
		"feedback": "Penjelasan sangat bagus dan tepat!",
	})
	req = httptest.NewRequest("PATCH", "/api/submissions/1/grade", bytes.NewReader(gradePayload))
	req.Header.Set("Authorization", "Bearer "+tokenTeacher1)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Failed to grade submission: %d (body: %s)", w.Code, w.Body.String())
	}
}
