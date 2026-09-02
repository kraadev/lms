package classes

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

func (r *Repository) ListForUser(userID int64, role models.Role) ([]models.Class, error) {
	var query string
	var args []interface{}

	if role == models.RoleAdmin {
		query = `
			SELECT c.id, c.name, c.description, c.teacher_id, u.name as teacher_name, c.academic_year, c.status,
			       (SELECT COUNT(*) FROM class_members cm WHERE cm.class_id = c.id) as member_count,
			       c.created_at, c.updated_at
			FROM classes c
			LEFT JOIN users u ON c.teacher_id = u.id
			ORDER BY c.id DESC;
		`
	} else if role == models.RoleTeacher {
		query = `
			SELECT c.id, c.name, c.description, c.teacher_id, u.name as teacher_name, c.academic_year, c.status,
			       (SELECT COUNT(*) FROM class_members cm WHERE cm.class_id = c.id) as member_count,
			       c.created_at, c.updated_at
			FROM classes c
			LEFT JOIN users u ON c.teacher_id = u.id
			WHERE c.teacher_id = $1
			ORDER BY c.id DESC;
		`
		if r.db.Driver != "postgres" {
			query = `
				SELECT c.id, c.name, c.description, c.teacher_id, u.name as teacher_name, c.academic_year, c.status,
				       (SELECT COUNT(*) FROM class_members cm WHERE cm.class_id = c.id) as member_count,
				       c.created_at, c.updated_at
				FROM classes c
				LEFT JOIN users u ON c.teacher_id = u.id
				WHERE c.teacher_id = ?
				ORDER BY c.id DESC;
			`
		}
		args = append(args, userID)
	} else {
		query = `
			SELECT c.id, c.name, c.description, c.teacher_id, u.name as teacher_name, c.academic_year, c.status,
			       (SELECT COUNT(*) FROM class_members cm WHERE cm.class_id = c.id) as member_count,
			       c.created_at, c.updated_at
			FROM classes c
			INNER JOIN class_members cm ON c.id = cm.class_id
			LEFT JOIN users u ON c.teacher_id = u.id
			WHERE cm.user_id = $1
			ORDER BY c.id DESC;
		`
		if r.db.Driver != "postgres" {
			query = `
				SELECT c.id, c.name, c.description, c.teacher_id, u.name as teacher_name, c.academic_year, c.status,
				       (SELECT COUNT(*) FROM class_members cm WHERE cm.class_id = c.id) as member_count,
				       c.created_at, c.updated_at
				FROM classes c
				INNER JOIN class_members cm ON c.id = cm.class_id
				LEFT JOIN users u ON c.teacher_id = u.id
				WHERE cm.user_id = ?
				ORDER BY c.id DESC;
			`
		}
		args = append(args, userID)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Class
	for rows.Next() {
		var c models.Class
		var teacherName sql.NullString
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Description, &c.TeacherID, &teacherName, &c.AcademicYear, &c.Status,
			&c.MemberCount, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if teacherName.Valid {
			c.TeacherName = teacherName.String
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *Repository) FindByID(id int64) (*models.Class, error) {
	query := `
		SELECT c.id, c.name, c.description, c.teacher_id, u.name as teacher_name, c.academic_year, c.status,
		       (SELECT COUNT(*) FROM class_members cm WHERE cm.class_id = c.id) as member_count,
		       c.created_at, c.updated_at
		FROM classes c
		LEFT JOIN users u ON c.teacher_id = u.id
		WHERE c.id = $1;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT c.id, c.name, c.description, c.teacher_id, u.name as teacher_name, c.academic_year, c.status,
			       (SELECT COUNT(*) FROM class_members cm WHERE cm.class_id = c.id) as member_count,
			       c.created_at, c.updated_at
			FROM classes c
			LEFT JOIN users u ON c.teacher_id = u.id
			WHERE c.id = ?;
		`
	}

	var c models.Class
	var teacherName sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &c.Name, &c.Description, &c.TeacherID, &teacherName, &c.AcademicYear, &c.Status,
		&c.MemberCount, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if teacherName.Valid {
		c.TeacherName = teacherName.String
	}
	return &c, nil
}

func (r *Repository) Create(class *models.Class) error {
	if r.db.Driver == "postgres" {
		query := `
			INSERT INTO classes (name, description, teacher_id, academic_year, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			RETURNING id, created_at, updated_at;
		`
		return r.db.QueryRow(query, class.Name, class.Description, class.TeacherID, class.AcademicYear, class.Status).Scan(&class.ID, &class.CreatedAt, &class.UpdatedAt)
	}

	query := `
		INSERT INTO classes (name, description, teacher_id, academic_year, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
	`
	res, err := r.db.Exec(query, class.Name, class.Description, class.TeacherID, class.AcademicYear, class.Status)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	class.ID = id
	return nil
}

func (r *Repository) Update(class *models.Class) error {
	query := `
		UPDATE classes
		SET name = $1, description = $2, academic_year = $3, status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5;
	`
	if r.db.Driver != "postgres" {
		query = `
			UPDATE classes
			SET name = ?, description = ?, academic_year = ?, status = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?;
		`
	}
	_, err := r.db.Exec(query, class.Name, class.Description, class.AcademicYear, class.Status, class.ID)
	return err
}

func (r *Repository) Delete(id int64) error {
	query := "DELETE FROM classes WHERE id = $1;"
	if r.db.Driver != "postgres" {
		query = "DELETE FROM classes WHERE id = ?;"
	}
	_, err := r.db.Exec(query, id)
	return err
}

func (r *Repository) AddMember(classID, userID int64) error {
	if r.db.Driver == "postgres" {
		_, err := r.db.Exec(`
			INSERT INTO class_members (class_id, user_id, joined_at)
			VALUES ($1, $2, CURRENT_TIMESTAMP)
			ON CONFLICT (class_id, user_id) DO NOTHING;
		`, classID, userID)
		return err
	}

	_, err := r.db.Exec(`
		INSERT OR IGNORE INTO class_members (class_id, user_id, joined_at)
		VALUES (?, ?, CURRENT_TIMESTAMP);
	`, classID, userID)
	return err
}

func (r *Repository) RemoveMember(classID, userID int64) error {
	query := "DELETE FROM class_members WHERE class_id = $1 AND user_id = $2;"
	if r.db.Driver != "postgres" {
		query = "DELETE FROM class_members WHERE class_id = ? AND user_id = ?;"
	}
	_, err := r.db.Exec(query, classID, userID)
	return err
}

func (r *Repository) ListMembers(classID int64) ([]models.ClassMember, error) {
	query := `
		SELECT cm.id, cm.class_id, cm.user_id, cm.joined_at, u.name, u.email, u.role, u.avatar_url
		FROM class_members cm
		JOIN users u ON cm.user_id = u.id
		WHERE cm.class_id = $1
		ORDER BY cm.joined_at ASC;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT cm.id, cm.class_id, cm.user_id, cm.joined_at, u.name, u.email, u.role, u.avatar_url
			FROM class_members cm
			JOIN users u ON cm.user_id = u.id
			WHERE cm.class_id = ?
			ORDER BY cm.joined_at ASC;
		`
	}

	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.ClassMember
	for rows.Next() {
		var cm models.ClassMember
		var u models.User
		var avatar sql.NullString
		if err := rows.Scan(
			&cm.ID, &cm.ClassID, &cm.UserID, &cm.JoinedAt,
			&u.Name, &u.Email, &u.Role, &avatar,
		); err != nil {
			return nil, err
		}
		u.ID = cm.UserID
		if avatar.Valid {
			u.AvatarURL = &avatar.String
		}
		cm.User = &u
		members = append(members, cm)
	}
	return members, nil
}
