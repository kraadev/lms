package users

import (
	"database/sql"
	"fmt"
	"strings"

	"lms/internal/database"
	"lms/internal/models"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByEmail(email string) (*models.User, error) {
	query := "SELECT id, name, email, password_hash, role, avatar_url, created_at, updated_at FROM users WHERE LOWER(email) = LOWER($1);"
	if r.db.Driver != "postgres" {
		query = "SELECT id, name, email, password_hash, role, avatar_url, created_at, updated_at FROM users WHERE LOWER(email) = LOWER(?);"
	}

	var u models.User
	var avatar sql.NullString
	err := r.db.QueryRow(query, strings.TrimSpace(email)).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &avatar, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if avatar.Valid {
		u.AvatarURL = &avatar.String
	}
	return &u, nil
}

func (r *Repository) FindByID(id int64) (*models.User, error) {
	query := "SELECT id, name, email, password_hash, role, avatar_url, created_at, updated_at FROM users WHERE id = $1;"
	if r.db.Driver != "postgres" {
		query = "SELECT id, name, email, password_hash, role, avatar_url, created_at, updated_at FROM users WHERE id = ?;"
	}

	var u models.User
	var avatar sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &avatar, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if avatar.Valid {
		u.AvatarURL = &avatar.String
	}
	return &u, nil
}

func (r *Repository) Create(user *models.User) error {
	if r.db.Driver == "postgres" {
		query := `
			INSERT INTO users (name, email, password_hash, role, avatar_url, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			RETURNING id, created_at, updated_at;
		`
		return r.db.QueryRow(query, user.Name, user.Email, user.PasswordHash, user.Role, user.AvatarURL).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	}

	query := `
		INSERT INTO users (name, email, password_hash, role, avatar_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
	`
	res, err := r.db.Exec(query, user.Name, user.Email, user.PasswordHash, user.Role, user.AvatarURL)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *Repository) List(roleFilter, search string) ([]models.User, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if roleFilter != "" {
		if r.db.Driver == "postgres" {
			conditions = append(conditions, fmt.Sprintf("role = $%d", argIdx))
		} else {
			conditions = append(conditions, "role = ?")
		}
		args = append(args, roleFilter)
		argIdx++
	}

	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		if r.db.Driver == "postgres" {
			conditions = append(conditions, fmt.Sprintf("(LOWER(name) LIKE $%d OR LOWER(email) LIKE $%d)", argIdx, argIdx))
		} else {
			conditions = append(conditions, "(LOWER(name) LIKE ? OR LOWER(email) LIKE ?)")
			args = append(args, searchTerm)
		}
		args = append(args, searchTerm)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf("SELECT id, name, email, role, avatar_url, created_at, updated_at FROM users %s ORDER BY id ASC;", whereClause)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.User
	for rows.Next() {
		var u models.User
		var avatar sql.NullString
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &avatar, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		if avatar.Valid {
			u.AvatarURL = &avatar.String
		}
		result = append(result, u)
	}
	return result, nil
}

func (r *Repository) Update(user *models.User) error {
	query := "UPDATE users SET name = $1, email = $2, role = $3, avatar_url = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5;"
	if r.db.Driver != "postgres" {
		query = "UPDATE users SET name = ?, email = ?, role = ?, avatar_url = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;"
	}
	_, err := r.db.Exec(query, user.Name, user.Email, user.Role, user.AvatarURL, user.ID)
	return err
}

func (r *Repository) Delete(id int64) error {
	query := "DELETE FROM users WHERE id = $1;"
	if r.db.Driver != "postgres" {
		query = "DELETE FROM users WHERE id = ?;"
	}
	_, err := r.db.Exec(query, id)
	return err
}
