package database

import (
	"fmt"
	"log"
	"time"

	"lms/internal/models"
	"lms/internal/utils"
)

func SeedData(db *DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users;").Scan(&count)
	if err == nil && count > 0 {
		log.Println("[SEED] Data already exists. Skipping seed.")
		return nil
	}

	log.Println("[SEED] Seeding initial database data...")

	passHash, _ := utils.HashPassword("password123")
	adminPassHash, _ := utils.HashPassword("admin123")

	// 1. Users
	users := []struct {
		Name     string
		Email    string
		Password string
		Role     models.Role
	}{
		{"System Administrator", "admin@lms.local", adminPassHash, models.RoleAdmin},
		{"Budi Santoso (Teacher)", "teacher1@lms.local", passHash, models.RoleTeacher},
		{"Siti Rahma (Teacher)", "teacher2@lms.local", passHash, models.RoleTeacher},
		{"Andi Pratama", "student1@lms.local", passHash, models.RoleStudent},
		{"Bunga Citra", "student2@lms.local", passHash, models.RoleStudent},
		{"Candra Wijaya", "student3@lms.local", passHash, models.RoleStudent},
		{"Dewi Lestari", "student4@lms.local", passHash, models.RoleStudent},
	}

	for _, u := range users {
		_, err := db.Exec(`
			INSERT INTO users (name, email, password_hash, role, created_at, updated_at)
			VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			ON CONFLICT (email) DO NOTHING;
		`, u.Name, u.Email, u.Password, u.Role)
		if err != nil {
			// For SQLite fallback
			_, _ = db.Exec(`
				INSERT OR IGNORE INTO users (name, email, password_hash, role, created_at, updated_at)
				VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
			`, u.Name, u.Email, u.Password, u.Role)
		}
	}

	// Fetch IDs
	var teacher1ID, teacher2ID, student1ID, student2ID, student3ID, student4ID int64
	_ = db.QueryRow("SELECT id FROM users WHERE email = 'teacher1@lms.local'").Scan(&teacher1ID)
	_ = db.QueryRow("SELECT id FROM users WHERE email = 'teacher2@lms.local'").Scan(&teacher2ID)
	_ = db.QueryRow("SELECT id FROM users WHERE email = 'student1@lms.local'").Scan(&student1ID)
	_ = db.QueryRow("SELECT id FROM users WHERE email = 'student2@lms.local'").Scan(&student2ID)
	_ = db.QueryRow("SELECT id FROM users WHERE email = 'student3@lms.local'").Scan(&student3ID)
	_ = db.QueryRow("SELECT id FROM users WHERE email = 'student4@lms.local'").Scan(&student4ID)

	if teacher1ID == 0 {
		return fmt.Errorf("failed to retrieve seeded user IDs")
	}

	// 2. Classes
	execInsert := func(queryPG, querySQLite string, args ...interface{}) (int64, error) {
		if db.Driver == "postgres" {
			var id int64
			err := db.QueryRow(queryPG+" RETURNING id;", args...).Scan(&id)
			return id, err
		}
		res, err := db.Exec(querySQLite, args...)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}

	class1ID, err := execInsert(
		"INSERT INTO classes (name, description, teacher_id, academic_year, status) VALUES ($1, $2, $3, $4, $5)",
		"INSERT INTO classes (name, description, teacher_id, academic_year, status) VALUES (?, ?, ?, ?, ?)",
		"Kelas Backend Go & Microservices",
		"Mempelajari arsitektur backend Go, REST API, WebSocket, dan PostgreSQL.",
		teacher1ID,
		"2026/2027 Ganjil",
		"active",
	)
	if err != nil {
		return err
	}

	class2ID, err := execInsert(
		"INSERT INTO classes (name, description, teacher_id, academic_year, status) VALUES ($1, $2, $3, $4, $5)",
		"INSERT INTO classes (name, description, teacher_id, academic_year, status) VALUES (?, ?, ?, ?, ?)",
		"Kelas Frontend & UI Architecture",
		"Mempelajari arsitektur SPA modern, state management, dan styling clean.",
		teacher2ID,
		"2026/2027 Ganjil",
		"active",
	)
	if err != nil {
		return err
	}

	// 3. Class Members
	// Class 1 has Student 1 & 2
	// Class 2 has Student 3 & 4
	memberships := []struct {
		ClassID int64
		UserID  int64
	}{
		{class1ID, student1ID},
		{class1ID, student2ID},
		{class2ID, student3ID},
		{class2ID, student4ID},
	}

	for _, m := range memberships {
		if db.Driver == "postgres" {
			_, _ = db.Exec("INSERT INTO class_members (class_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;", m.ClassID, m.UserID)
		} else {
			_, _ = db.Exec("INSERT OR IGNORE INTO class_members (class_id, user_id) VALUES (?, ?);", m.ClassID, m.UserID)
		}
	}

	// 4. Learning Materials
	_, _ = execInsert(
		"INSERT INTO materials (class_id, teacher_id, title, description, content, external_url) VALUES ($1, $2, $3, $4, $5, $6)",
		"INSERT INTO materials (class_id, teacher_id, title, description, content, external_url) VALUES (?, ?, ?, ?, ?, ?)",
		class1ID,
		teacher1ID,
		"Modul 01: Pengenalan Concurrency Go",
		"Memahami Goroutines, Channels, dan WaitGroups.",
		"# Concurrency in Go\n\nGo memiliki model konkurensi bawaan bernama CSP (Communicating Sequential Processes).",
		"https://go.dev/tour/concurrency/1",
	)

	// 5. Assignments
	deadline := time.Now().Add(7 * 24 * time.Hour)
	assignmentID, _ := execInsert(
		"INSERT INTO assignments (class_id, teacher_id, title, description, deadline, max_score) VALUES ($1, $2, $3, $4, $5, $6)",
		"INSERT INTO assignments (class_id, teacher_id, title, description, deadline, max_score) VALUES (?, ?, ?, ?, ?, ?)",
		class1ID,
		teacher1ID,
		"Tugas 1: Implementasi WebSocket Hub di Go",
		"Buat concurrency-safe WebSocket room broker menggunakan channels dan mutex.",
		deadline,
		100.0,
	)

	// 6. Quizzes
	startAt := time.Now().Add(-1 * time.Hour)
	endAt := time.Now().Add(14 * 24 * time.Hour)
	quizID, _ := execInsert(
		"INSERT INTO quizzes (class_id, teacher_id, title, description, duration_minutes, start_at, end_at, max_attempts, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		"INSERT INTO quizzes (class_id, teacher_id, title, description, duration_minutes, start_at, end_at, max_attempts, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		class1ID,
		teacher1ID,
		"Kuis 1: Konsep Dasar Go",
		"Uji pemahaman dasar sintaksis, pointer, dan interface di Go.",
		30,
		startAt,
		endAt,
		3,
		"published",
	)

	// Question 1: Multiple Choice
	q1ID, _ := execInsert(
		"INSERT INTO quiz_questions (quiz_id, question, type, points, order_index) VALUES ($1, $2, $3, $4, $5)",
		"INSERT INTO quiz_questions (quiz_id, question, type, points, order_index) VALUES (?, ?, ?, ?, ?)",
		quizID,
		"Apa keyword untuk memulai sebuah Goroutine baru di Go?",
		models.QuestionMultipleChoice,
		10.0,
		1,
	)

	optionsQ1 := []struct {
		Text      string
		IsCorrect bool
	}{
		{"go", true},
		{"async", false},
		{"thread", false},
		{"spawn", false},
	}
	for _, opt := range optionsQ1 {
		_, _ = execInsert(
			"INSERT INTO quiz_options (question_id, option_text, is_correct) VALUES ($1, $2, $3)",
			"INSERT INTO quiz_options (question_id, option_text, is_correct) VALUES (?, ?, ?)",
			q1ID,
			opt.Text,
			opt.IsCorrect,
		)
	}

	// Question 2: True/False
	q2ID, _ := execInsert(
		"INSERT INTO quiz_questions (quiz_id, question, type, points, order_index) VALUES ($1, $2, $3, $4, $5)",
		"INSERT INTO quiz_questions (quiz_id, question, type, points, order_index) VALUES (?, ?, ?, ?, ?)",
		quizID,
		"Di Go, map aman untuk dibaca dan ditulis secara konkuren tanpa mutex.",
		models.QuestionTrueFalse,
		10.0,
		2,
	)
	optionsQ2 := []struct {
		Text      string
		IsCorrect bool
	}{
		{"True", false},
		{"False", true},
	}
	for _, opt := range optionsQ2 {
		_, _ = execInsert(
			"INSERT INTO quiz_options (question_id, option_text, is_correct) VALUES ($1, $2, $3)",
			"INSERT INTO quiz_options (question_id, option_text, is_correct) VALUES (?, ?, ?)",
			q2ID,
			opt.Text,
			opt.IsCorrect,
		)
	}

	// 7. Sample Messages
	_, _ = execInsert(
		"INSERT INTO messages (class_id, user_id, message, created_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP)",
		"INSERT INTO messages (class_id, user_id, message, created_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP)",
		class1ID,
		teacher1ID,
		"Selamat datang di kelas Backend Go! Silakan periksa materi Modul 01.",
	)

	// 8. Sample Notifications
	_, _ = execInsert(
		"INSERT INTO notifications (user_id, type, title, message) VALUES ($1, $2, $3, $4)",
		"INSERT INTO notifications (user_id, type, title, message) VALUES (?, ?, ?, ?)",
		student1ID,
		"assignment",
		"Tugas Baru Tersedia",
		"Guru telah menerbitkan Tugas 1: Implementasi WebSocket Hub di Go.",
	)

	log.Printf("[SEED] Database seeded successfully! (Class 1 ID: %d, Class 2 ID: %d, Assignment ID: %d, Quiz ID: %d)", class1ID, class2ID, assignmentID, quizID)
	return nil
}
