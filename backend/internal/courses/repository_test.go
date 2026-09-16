package courses

import (
	"context"
	"testing"

	"lms/internal/config"
	"lms/internal/database"
	"lms/internal/models"
)

func TestCourseRepository(t *testing.T) {
	cfg := &config.Config{
		DBDriver:     "sqlite",
		DBSQLitePath: ":memory:",
	}
	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Insert parent records to satisfy foreign key constraints
	_, err = db.Exec("INSERT INTO users (id, name, email, password_hash, role) VALUES (1, 'Guru', 'guru@lms.local', 'hash', 'teacher')")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}
	_, err = db.Exec("INSERT INTO classes (id, name, teacher_id) VALUES (1, 'Kelas Algoritma', 1)")
	if err != nil {
		t.Fatalf("failed to insert class: %v", err)
	}

	repo := NewRepository(db)
	ctx := context.Background()

	// 1. Create Course
	course := &models.Course{
		ClassID:     1,
		Title:       "Algoritma Pemrograman Lanjut",
		Description: "Silabus lengkap kurikulum 2026",
	}
	if err := repo.CreateCourse(ctx, course); err != nil {
		t.Fatalf("failed to create course: %v", err)
	}
	if course.ID == 0 {
		t.Errorf("expected positive course ID, got 0")
	}

	// 2. Fetch Course
	found, err := repo.GetCourseByClassID(ctx, 1)
	if err != nil {
		t.Fatalf("failed to fetch course: %v", err)
	}
	if found == nil || found.Title != course.Title {
		t.Errorf("expected course '%s', got %+v", course.Title, found)
	}
}
