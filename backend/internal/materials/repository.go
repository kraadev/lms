package materials

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

func (r *Repository) ListByClass(classID int64) ([]models.Material, error) {
	query := `
		SELECT id, class_id, teacher_id, title, description, content, file_path, external_url, created_at, updated_at
		FROM materials
		WHERE class_id = $1
		ORDER BY id DESC;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT id, class_id, teacher_id, title, description, content, file_path, external_url, created_at, updated_at
			FROM materials
			WHERE class_id = ?
			ORDER BY id DESC;
		`
	}

	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Material
	for rows.Next() {
		var m models.Material
		var filePath, extURL sql.NullString
		if err := rows.Scan(
			&m.ID, &m.ClassID, &m.TeacherID, &m.Title, &m.Description, &m.Content,
			&filePath, &extURL, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if filePath.Valid {
			m.FilePath = &filePath.String
		}
		if extURL.Valid {
			m.ExternalURL = &extURL.String
		}
		list = append(list, m)
	}
	return list, nil
}

func (r *Repository) FindByID(id int64) (*models.Material, error) {
	query := `
		SELECT id, class_id, teacher_id, title, description, content, file_path, external_url, created_at, updated_at
		FROM materials
		WHERE id = $1;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT id, class_id, teacher_id, title, description, content, file_path, external_url, created_at, updated_at
			FROM materials
			WHERE id = ?;
		`
	}

	var m models.Material
	var filePath, extURL sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&m.ID, &m.ClassID, &m.TeacherID, &m.Title, &m.Description, &m.Content,
		&filePath, &extURL, &m.CreatedAt, &m.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if filePath.Valid {
		m.FilePath = &filePath.String
	}
	if extURL.Valid {
		m.ExternalURL = &extURL.String
	}
	return &m, nil
}

func (r *Repository) Create(m *models.Material) error {
	if r.db.Driver == "postgres" {
		query := `
			INSERT INTO materials (class_id, teacher_id, title, description, content, file_path, external_url, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			RETURNING id, created_at, updated_at;
		`
		return r.db.QueryRow(query, m.ClassID, m.TeacherID, m.Title, m.Description, m.Content, m.FilePath, m.ExternalURL).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	}

	query := `
		INSERT INTO materials (class_id, teacher_id, title, description, content, file_path, external_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
	`
	res, err := r.db.Exec(query, m.ClassID, m.TeacherID, m.Title, m.Description, m.Content, m.FilePath, m.ExternalURL)
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

func (r *Repository) Update(m *models.Material) error {
	query := `
		UPDATE materials
		SET title = $1, description = $2, content = $3, file_path = $4, external_url = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6;
	`
	if r.db.Driver != "postgres" {
		query = `
			UPDATE materials
			SET title = ?, description = ?, content = ?, file_path = ?, external_url = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?;
		`
	}
	_, err := r.db.Exec(query, m.Title, m.Description, m.Content, m.FilePath, m.ExternalURL, m.ID)
	return err
}

func (r *Repository) Delete(id int64) error {
	query := "DELETE FROM materials WHERE id = $1;"
	if r.db.Driver != "postgres" {
		query = "DELETE FROM materials WHERE id = ?;"
	}
	_, err := r.db.Exec(query, id)
	return err
}
