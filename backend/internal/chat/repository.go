package chat

import (
	"database/sql"
	"time"

	"lms/internal/database"
	"lms/internal/models"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveMessage(classID, userID int64, message string) (*models.Message, error) {
	var msg models.Message
	msg.ClassID = classID
	msg.UserID = userID
	msg.Message = message
	msg.CreatedAt = time.Now()

	if r.db.Driver == "postgres" {
		query := `
			INSERT INTO messages (class_id, user_id, message, created_at)
			VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
			RETURNING id, created_at;
		`
		err := r.db.QueryRow(query, classID, userID, message).Scan(&msg.ID, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
	} else {
		query := `
			INSERT INTO messages (class_id, user_id, message, created_at)
			VALUES (?, ?, ?, CURRENT_TIMESTAMP);
		`
		res, err := r.db.Exec(query, classID, userID, message)
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		msg.ID = id
	}

	// Fetch user details for broadcast
	userQuery := "SELECT name, role, avatar_url FROM users WHERE id = $1;"
	if r.db.Driver != "postgres" {
		userQuery = "SELECT name, role, avatar_url FROM users WHERE id = ?;"
	}

	var avatar sql.NullString
	_ = r.db.QueryRow(userQuery, userID).Scan(&msg.UserName, &msg.UserRole, &avatar)
	if avatar.Valid {
		msg.UserAvatar = &avatar.String
	}

	return &msg, nil
}

func (r *Repository) ListMessages(classID int64, limit, offset int) ([]models.Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT m.id, m.class_id, m.user_id, u.name as user_name, u.role as user_role, u.avatar_url, m.message, m.created_at
		FROM messages m
		JOIN users u ON m.user_id = u.id
		WHERE m.class_id = $1
		ORDER BY m.created_at ASC
		LIMIT $2 OFFSET $3;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT m.id, m.class_id, m.user_id, u.name as user_name, u.role as user_role, u.avatar_url, m.message, m.created_at
			FROM messages m
			JOIN users u ON m.user_id = u.id
			WHERE m.class_id = ?
			ORDER BY m.created_at ASC
			LIMIT ? OFFSET ?;
		`
	}

	rows, err := r.db.Query(query, classID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		var avatar sql.NullString
		if err := rows.Scan(
			&msg.ID, &msg.ClassID, &msg.UserID, &msg.UserName, &msg.UserRole, &avatar,
			&msg.Message, &msg.CreatedAt,
		); err != nil {
			return nil, err
		}
		if avatar.Valid {
			msg.UserAvatar = &avatar.String
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
