package courses

import (
	"context"
	"database/sql"

	"lms/internal/database"
	"lms/internal/models"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateCourse(ctx context.Context, c *models.Course) error {
	query := `INSERT INTO courses (class_id, title, description) VALUES (?, ?, ?)`
	if r.db.Driver == "postgres" {
		query = `INSERT INTO courses (class_id, title, description) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
		return r.db.QueryRowContext(ctx, query, c.ClassID, c.Title, c.Description).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	}

	res, err := r.db.ExecContext(ctx, query, c.ClassID, c.Title, c.Description)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (r *Repository) GetCourseByClassID(ctx context.Context, classID int64) (*models.Course, error) {
	query := `SELECT id, class_id, title, description, created_at, updated_at FROM courses WHERE class_id = ?`
	if r.db.Driver == "postgres" {
		query = `SELECT id, class_id, title, description, created_at, updated_at FROM courses WHERE class_id = $1`
	}

	var c models.Course
	err := r.db.QueryRowContext(ctx, query, classID).Scan(
		&c.ID, &c.ClassID, &c.Title, &c.Description, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}
