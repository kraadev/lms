package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"lms/internal/assignments"
	"lms/internal/auth"
	"lms/internal/chat"
	"lms/internal/classes"
	"lms/internal/config"
	"lms/internal/dashboard"
	"lms/internal/database"
	"lms/internal/materials"
	"lms/internal/meetings"
	"lms/internal/middleware"
	"lms/internal/notifications"
	"lms/internal/quizzes"
	"lms/internal/storage"
	"lms/internal/users"
	"lms/internal/utils"
)

func main() {
	seedFlag := flag.Bool("seed", false, "Run database seeding upon startup")
	flag.Parse()

	// 1. Config
	cfg := config.Load()

	// 2. Database Connection
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("[FATAL] Could not connect to database: %v", err)
	}
	defer db.Close()

	// 3. Run Migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("[FATAL] Migration failure: %v", err)
	}

	// 4. Seeding (auto-seed if empty or -seed flag provided)
	if *seedFlag || cfg.Env == "development" {
		if err := database.SeedData(db); err != nil {
			log.Printf("[WARNING] Seeding encountered issue: %v", err)
		}
	}

	// 5. Storage Service
	storageService, err := storage.NewStorageService(cfg.UploadDir)
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize storage: %v", err)
	}

	// 6. Repositories
	userRepo := users.NewRepository(db)
	classRepo := classes.NewRepository(db)
	materialRepo := materials.NewRepository(db)
	assignmentRepo := assignments.NewRepository(db)
	quizRepo := quizzes.NewRepository(db)
	meetingRepo := meetings.NewRepository(db)
	chatRepo := chat.NewRepository(db)
	notifRepo := notifications.NewRepository(db)

	// 7. Core Access Control & Services
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

	// 8. Realtime Chat Hub
	chatHub := chat.NewHub()
	go chatHub.Run()

	// 9. Handlers
	authHandler := auth.NewHandler(authService)
	userHandler := users.NewHandler(userService)
	classHandler := classes.NewHandler(classService, accessController)
	materialHandler := materials.NewHandler(materialService, accessController, storageService)
	assignmentHandler := assignments.NewHandler(assignmentService, accessController, storageService)
	quizHandler := quizzes.NewHandler(quizService, accessController)
	meetingHandler := meetings.NewHandler(meetingService, accessController)
	chatHandler := chat.NewHandler(chatHub, chatRepo, accessController, cfg)
	notifHandler := notifications.NewHandler(notifService)

	dashboardHandler := dashboard.NewHandler(db)

	// 10. HTTP Router Setup
	r := chi.NewRouter()

	// Global Middlewares
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// Global CORS Setup - allow all local network, IP, and tunnel origins
	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.AuthMiddleware(cfg))

	// Base Health & Info Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.JSON(w, http.StatusOK, map[string]interface{}{
			"name":    "LMS Backend API",
			"version": "1.0.0",
			"status":  "running",
			"driver":  db.Driver,
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.JSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Realtime WebSocket Endpoint
	r.Get("/ws", chatHandler.ServeWS)

	// Safe File Serving Endpoint
	r.Get("/api/files/*", storageService.ServeFileHandler())

	// Public Auth Routes
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/logout", authHandler.Logout)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAuth)
			r.Get("/me", authHandler.Me)
		})
	})

	// Protected Application API Routes
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.RequireAuth)

		// Dashboard Endpoints
		r.Route("/dashboard", func(r chi.Router) {
			r.Get("/student", dashboardHandler.StudentDashboard)
			r.Get("/teacher", dashboardHandler.TeacherDashboard)
			r.Get("/admin", dashboardHandler.AdminDashboard)
		})

		// Users Management
		r.Route("/users", func(r chi.Router) {
			r.Get("/", userHandler.List)
			r.Get("/{id}", userHandler.GetByID)

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAdmin)
				r.Post("/", userHandler.Create)
				r.Put("/{id}", userHandler.Update)
				r.Patch("/{id}", userHandler.Update)
				r.Delete("/{id}", userHandler.Delete)
			})
		})

		// Admin Aliases
		r.Route("/admin", func(r chi.Router) {
			r.Get("/users", userHandler.List)

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAdmin)
				r.Post("/users", userHandler.Create)
				r.Get("/users/{id}", userHandler.GetByID)
				r.Put("/users/{id}", userHandler.Update)
				r.Patch("/users/{id}", userHandler.Update)
				r.Delete("/users/{id}", userHandler.Delete)

				r.Get("/classes", classHandler.List)
				r.Post("/classes", classHandler.Create)
				r.Get("/classes/{id}", classHandler.GetByID)
				r.Put("/classes/{id}", classHandler.Update)
				r.Patch("/classes/{id}", classHandler.Update)
				r.Delete("/classes/{id}", classHandler.Delete)
			})
		})

		// Classes Management
		r.Route("/classes", func(r chi.Router) {
			r.Get("/", classHandler.List)
			r.Post("/", classHandler.Create)
			r.Get("/{id}", classHandler.GetByID)
			r.Put("/{id}", classHandler.Update)
			r.Patch("/{id}", classHandler.Update)
			r.Delete("/{id}", classHandler.Delete)

			// Class Members
			r.Get("/{id}/members", classHandler.ListMembers)
			r.Post("/{id}/members", classHandler.AddMember)
			r.Delete("/{id}/members/{userId}", classHandler.RemoveMember)

			// Nested Class Resources
			r.Get("/{classId}/materials", materialHandler.ListByClass)
			r.Post("/{classId}/materials", materialHandler.Create)

			r.Get("/{classId}/assignments", assignmentHandler.ListByClass)
			r.Post("/{classId}/assignments", assignmentHandler.Create)

			r.Get("/{classId}/quizzes", quizHandler.ListByClass)
			r.Post("/{classId}/quizzes", quizHandler.Create)

			r.Get("/{classId}/meetings", meetingHandler.ListByClass)
			r.Post("/{classId}/meetings", meetingHandler.Create)

			r.Get("/{classId}/messages", chatHandler.ListMessages)
			r.Post("/{classId}/messages", chatHandler.SendMessage)

			// Announcements (mock / class notice)
			r.Get("/{classId}/announcements", func(w http.ResponseWriter, r *http.Request) {
				utils.JSON(w, http.StatusOK, []interface{}{})
			})
			r.Post("/{classId}/announcements", func(w http.ResponseWriter, r *http.Request) {
				utils.JSON(w, http.StatusCreated, map[string]string{"message": "Announcement created"})
			})
		})

		// Materials Management
		r.Route("/materials", func(r chi.Router) {
			r.Get("/{id}", materialHandler.GetByID)
			r.Put("/{id}", materialHandler.Update)
			r.Patch("/{id}", materialHandler.Update)
			r.Delete("/{id}", materialHandler.Delete)
		})

		// Assignments Management
		r.Route("/assignments", func(r chi.Router) {
			r.Get("/", assignmentHandler.ListAll)
			r.Get("/{id}", assignmentHandler.GetByID)
			r.Put("/{id}", assignmentHandler.Update)
			r.Patch("/{id}", assignmentHandler.Update)
			r.Delete("/{id}", assignmentHandler.Delete)

			r.Post("/{id}/submissions", assignmentHandler.Submit)
			r.Post("/{id}/submit", assignmentHandler.Submit)
			r.Get("/{id}/my-submission", assignmentHandler.GetMySubmission)
			r.Get("/{id}/submissions", assignmentHandler.ListSubmissions)
		})

		// Submissions Grading
		r.Patch("/submissions/{id}/grade", assignmentHandler.Grade)
		r.Post("/submissions/{id}/grade", assignmentHandler.Grade)

		// Quizzes Management & Attempts
		r.Route("/quizzes", func(r chi.Router) {
			r.Get("/", quizHandler.ListAll)
			r.Get("/{id}", quizHandler.GetByID)
			r.Delete("/{id}", quizHandler.Delete)
			r.Post("/{id}/attempts", quizHandler.StartAttempt)
			r.Get("/{id}/attempts", quizHandler.ListAttempts)

			r.Route("/attempts", func(r chi.Router) {
				r.Post("/{id}/submit", quizHandler.SubmitAttempt)
				r.Get("/{id}", quizHandler.GetAttempt)
			})
		})

		// Standalone attempts route alias for frontend
		r.Route("/attempts", func(r chi.Router) {
			r.Get("/{id}", quizHandler.GetAttempt)
			r.Post("/{id}/submit", quizHandler.SubmitAttempt)
		})

		// Meetings Management
		r.Route("/meetings", func(r chi.Router) {
			r.Get("/{id}", meetingHandler.GetByID)
			r.Post("/{id}/join", meetingHandler.Join)
			r.Post("/{id}/end", meetingHandler.End)
		})

		// Notifications Management
		r.Route("/notifications", func(r chi.Router) {
			r.Get("/", notifHandler.List)
			r.Patch("/{id}/read", notifHandler.MarkRead)
			r.Patch("/read-all", notifHandler.MarkAllRead)
			r.Post("/read-all", notifHandler.MarkAllRead)
		})
	})

	// Start Server Gracefully on 0.0.0.0
	serverAddr := fmt.Sprintf("0.0.0.0:%s", cfg.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("==================================================")
		log.Printf("🚀 LMS Backend running at http://%s (listening on all interfaces)", serverAddr)
		log.Printf("📡 WebSocket available at ws://%s/ws", serverAddr)
		log.Printf("📦 Storage Directory: %s", cfg.UploadDir)
		log.Printf("==================================================")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] Server startup failed: %v", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[INFO] Shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[FATAL] Server forced to shutdown: %v", err)
	}

	log.Println("[INFO] Server exited cleanly.")
}
