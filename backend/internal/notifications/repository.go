package notifications

import (
	"lms/internal/database"
	"lms/internal/models"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(n *models.Notification) error {
	if r.db.Driver == "postgres" {
		query := `
			INSERT INTO notifications (user_id, type, title, message, is_read, created_at)
			VALUES ($1, $2, $3, $4, FALSE, CURRENT_TIMESTAMP)
			RETURNING id, created_at;
		`
		return r.db.QueryRow(query, n.UserID, n.Type, n.Title, n.Message).Scan(&n.ID, &n.CreatedAt)
	}

	query := `
		INSERT INTO notifications (user_id, type, title, message, is_read, created_at)
		VALUES (?, ?, ?, ?, 0, CURRENT_TIMESTAMP);
	`
	res, err := r.db.Exec(query, n.UserID, n.Type, n.Title, n.Message)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	n.ID = id
	return nil
}

func (r *Repository) ListByUser(userID int64) ([]models.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, is_read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 50;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT id, user_id, type, title, message, is_read, created_at
			FROM notifications
			WHERE user_id = ?
			ORDER BY created_at DESC
			LIMIT 50;
		`
	}

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Notification
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, n)
	}
	return list, nil
}

func (r *Repository) MarkAsRead(id, userID int64) error {
	query := "UPDATE notifications SET is_read = TRUE WHERE id = $1 AND user_id = $2;"
	if r.db.Driver != "postgres" {
		query = "UPDATE notifications SET is_read = 1 WHERE id = ? AND user_id = ?;"
	}
	_, err := r.db.Exec(query, id, userID)
	return err
}

func (r *Repository) MarkAllRead(userID int64) error {
	query := "UPDATE notifications SET is_read = TRUE WHERE user_id = $1;"
	if r.db.Driver != "postgres" {
		query = "UPDATE notifications SET is_read = 1 WHERE user_id = ?;"
	}
	_, err := r.db.Exec(query, userID)
	return err
}
