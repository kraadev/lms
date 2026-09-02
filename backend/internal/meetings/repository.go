package meetings

import (
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

func (r *Repository) ListByClass(classID int64) ([]models.Meeting, error) {
	query := `
		SELECT m.id, m.class_id, m.teacher_id, u.name as teacher_name, m.title, m.room_name, m.type, m.status,
		       m.started_at, m.ended_at, m.created_at
		FROM meetings m
		JOIN users u ON m.teacher_id = u.id
		WHERE m.class_id = $1
		ORDER BY m.id DESC;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT m.id, m.class_id, m.teacher_id, u.name as teacher_name, m.title, m.room_name, m.type, m.status,
			       m.started_at, m.ended_at, m.created_at
			FROM meetings m
			JOIN users u ON m.teacher_id = u.id
			WHERE m.class_id = ?
			ORDER BY m.id DESC;
		`
	}

	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Meeting
	for rows.Next() {
		var m models.Meeting
		var startedAt, endedAt sql.NullTime
		if err := rows.Scan(
			&m.ID, &m.ClassID, &m.TeacherID, &m.TeacherName, &m.Title, &m.RoomName, &m.Type, &m.Status,
			&startedAt, &endedAt, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		if startedAt.Valid {
			m.StartedAt = &startedAt.Time
		}
		if endedAt.Valid {
			m.EndedAt = &endedAt.Time
		}
		list = append(list, m)
	}
	return list, nil
}

func (r *Repository) FindByID(id int64) (*models.Meeting, error) {
	query := `
		SELECT m.id, m.class_id, m.teacher_id, u.name as teacher_name, m.title, m.room_name, m.type, m.status,
		       m.started_at, m.ended_at, m.created_at
		FROM meetings m
		JOIN users u ON m.teacher_id = u.id
		WHERE m.id = $1;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT m.id, m.class_id, m.teacher_id, u.name as teacher_name, m.title, m.room_name, m.type, m.status,
			       m.started_at, m.ended_at, m.created_at
			FROM meetings m
			JOIN users u ON m.teacher_id = u.id
			WHERE m.id = ?;
		`
	}

	var m models.Meeting
	var startedAt, endedAt sql.NullTime
	err := r.db.QueryRow(query, id).Scan(
		&m.ID, &m.ClassID, &m.TeacherID, &m.TeacherName, &m.Title, &m.RoomName, &m.Type, &m.Status,
		&startedAt, &endedAt, &m.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if startedAt.Valid {
		m.StartedAt = &startedAt.Time
	}
	if endedAt.Valid {
		m.EndedAt = &endedAt.Time
	}
	return &m, nil
}

func (r *Repository) Create(m *models.Meeting) error {
	if r.db.Driver == "postgres" {
		query := `
			INSERT INTO meetings (class_id, teacher_id, title, room_name, type, status, started_at, created_at)
			VALUES ($1, $2, $3, $4, $5, 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			RETURNING id, created_at;
		`
		return r.db.QueryRow(query, m.ClassID, m.TeacherID, m.Title, m.RoomName, m.Type).Scan(&m.ID, &m.CreatedAt)
	}

	query := `
		INSERT INTO meetings (class_id, teacher_id, title, room_name, type, status, started_at, created_at)
		VALUES (?, ?, ?, ?, ?, 'active', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
	`
	res, err := r.db.Exec(query, m.ClassID, m.TeacherID, m.Title, m.RoomName, m.Type)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	m.ID = id
	return nil
}

func (r *Repository) EndMeeting(id int64) error {
	query := `
		UPDATE meetings
		SET status = 'ended', ended_at = CURRENT_TIMESTAMP
		WHERE id = $1;
	`
	if r.db.Driver != "postgres" {
		query = `
			UPDATE meetings
			SET status = 'ended', ended_at = CURRENT_TIMESTAMP
			WHERE id = ?;
		`
	}
	_, err := r.db.Exec(query, id)
	return err
}
