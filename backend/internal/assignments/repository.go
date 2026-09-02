package assignments

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

func (r *Repository) ListByClass(classID int64) ([]models.Assignment, error) {
	query := `
		SELECT a.id, a.class_id, a.teacher_id, a.title, a.description, a.attachment_path, a.deadline, a.max_score,
		       (SELECT COUNT(*) FROM assignment_submissions s WHERE s.assignment_id = a.id) as sub_count,
		       a.created_at, a.updated_at
		FROM assignments a
		WHERE a.class_id = $1
		ORDER BY a.deadline ASC;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT a.id, a.class_id, a.teacher_id, a.title, a.description, a.attachment_path, a.deadline, a.max_score,
			       (SELECT COUNT(*) FROM assignment_submissions s WHERE s.assignment_id = a.id) as sub_count,
			       a.created_at, a.updated_at
			FROM assignments a
			WHERE a.class_id = ?
			ORDER BY a.deadline ASC;
		`
	}

	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Assignment
	for rows.Next() {
		var a models.Assignment
		var attach sql.NullString
		if err := rows.Scan(
			&a.ID, &a.ClassID, &a.TeacherID, &a.Title, &a.Description, &attach,
			&a.Deadline, &a.MaxScore, &a.SubmissionCount, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if attach.Valid {
			a.AttachmentPath = &attach.String
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *Repository) FindByID(id int64) (*models.Assignment, error) {
	query := `
		SELECT a.id, a.class_id, a.teacher_id, a.title, a.description, a.attachment_path, a.deadline, a.max_score,
		       (SELECT COUNT(*) FROM assignment_submissions s WHERE s.assignment_id = a.id) as sub_count,
		       a.created_at, a.updated_at
		FROM assignments a
		WHERE a.id = $1;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT a.id, a.class_id, a.teacher_id, a.title, a.description, a.attachment_path, a.deadline, a.max_score,
			       (SELECT COUNT(*) FROM assignment_submissions s WHERE s.assignment_id = a.id) as sub_count,
			       a.created_at, a.updated_at
			FROM assignments a
			WHERE a.id = ?;
		`
	}

	var a models.Assignment
	var attach sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&a.ID, &a.ClassID, &a.TeacherID, &a.Title, &a.Description, &attach,
		&a.Deadline, &a.MaxScore, &a.SubmissionCount, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if attach.Valid {
		a.AttachmentPath = &attach.String
	}
	return &a, nil
}

func (r *Repository) Create(a *models.Assignment) error {
	if r.db.Driver == "postgres" {
		query := `
			INSERT INTO assignments (class_id, teacher_id, title, description, attachment_path, deadline, max_score, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			RETURNING id, created_at, updated_at;
		`
		return r.db.QueryRow(query, a.ClassID, a.TeacherID, a.Title, a.Description, a.AttachmentPath, a.Deadline, a.MaxScore).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
	}

	query := `
		INSERT INTO assignments (class_id, teacher_id, title, description, attachment_path, deadline, max_score, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
	`
	res, err := r.db.Exec(query, a.ClassID, a.TeacherID, a.Title, a.Description, a.AttachmentPath, a.Deadline, a.MaxScore)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	a.ID = id
	return nil
}

func (r *Repository) Update(a *models.Assignment) error {
	query := `
		UPDATE assignments
		SET title = $1, description = $2, attachment_path = $3, deadline = $4, max_score = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6;
	`
	if r.db.Driver != "postgres" {
		query = `
			UPDATE assignments
			SET title = ?, description = ?, attachment_path = ?, deadline = ?, max_score = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?;
		`
	}
	_, err := r.db.Exec(query, a.Title, a.Description, a.AttachmentPath, a.Deadline, a.MaxScore, a.ID)
	return err
}

func (r *Repository) Delete(id int64) error {
	query := "DELETE FROM assignments WHERE id = $1;"
	if r.db.Driver != "postgres" {
		query = "DELETE FROM assignments WHERE id = ?;"
	}
	_, err := r.db.Exec(query, id)
	return err
}

func (r *Repository) UpsertSubmission(s *models.AssignmentSubmission) error {
	if r.db.Driver == "postgres" {
		query := `
			INSERT INTO assignment_submissions (assignment_id, student_id, text_answer, file_path, submitted_at, updated_at, status)
			VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $5)
			ON CONFLICT (assignment_id, student_id)
			DO UPDATE SET text_answer = EXCLUDED.text_answer, file_path = EXCLUDED.file_path, updated_at = CURRENT_TIMESTAMP, status = EXCLUDED.status
			RETURNING id, submitted_at, updated_at;
		`
		return r.db.QueryRow(query, s.AssignmentID, s.StudentID, s.TextAnswer, s.FilePath, s.Status).Scan(&s.ID, &s.SubmittedAt, &s.UpdatedAt)
	}

	query := `
		INSERT INTO assignment_submissions (assignment_id, student_id, text_answer, file_path, submitted_at, updated_at, status)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?)
		ON CONFLICT (assignment_id, student_id)
		DO UPDATE SET text_answer = excluded.text_answer, file_path = excluded.file_path, updated_at = CURRENT_TIMESTAMP, status = excluded.status;
	`
	_, err := r.db.Exec(query, s.AssignmentID, s.StudentID, s.TextAnswer, s.FilePath, s.Status)
	return err
}

func (r *Repository) GetSubmission(assignmentID, studentID int64) (*models.AssignmentSubmission, error) {
	query := `
		SELECT s.id, s.assignment_id, s.student_id, u.name as student_name, u.email as student_email,
		       s.text_answer, s.file_path, s.submitted_at, s.updated_at, s.score, s.feedback, s.status
		FROM assignment_submissions s
		JOIN users u ON s.student_id = u.id
		WHERE s.assignment_id = $1 AND s.student_id = $2;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT s.id, s.assignment_id, s.student_id, u.name as student_name, u.email as student_email,
			       s.text_answer, s.file_path, s.submitted_at, s.updated_at, s.score, s.feedback, s.status
			FROM assignment_submissions s
			JOIN users u ON s.student_id = u.id
			WHERE s.assignment_id = ? AND s.student_id = ?;
		`
	}

	var s models.AssignmentSubmission
	var textAns, filePath, feedback sql.NullString
	var score sql.NullFloat64
	err := r.db.QueryRow(query, assignmentID, studentID).Scan(
		&s.ID, &s.AssignmentID, &s.StudentID, &s.StudentName, &s.StudentEmail,
		&textAns, &filePath, &s.SubmittedAt, &s.UpdatedAt, &score, &feedback, &s.Status,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if textAns.Valid {
		s.TextAnswer = &textAns.String
	}
	if filePath.Valid {
		s.FilePath = &filePath.String
	}
	if score.Valid {
		s.Score = &score.Float64
	}
	if feedback.Valid {
		s.Feedback = &feedback.String
	}
	return &s, nil
}

func (r *Repository) GetSubmissionByID(id int64) (*models.AssignmentSubmission, error) {
	query := `
		SELECT s.id, s.assignment_id, s.student_id, u.name as student_name, u.email as student_email,
		       s.text_answer, s.file_path, s.submitted_at, s.updated_at, s.score, s.feedback, s.status
		FROM assignment_submissions s
		JOIN users u ON s.student_id = u.id
		WHERE s.id = $1;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT s.id, s.assignment_id, s.student_id, u.name as student_name, u.email as student_email,
			       s.text_answer, s.file_path, s.submitted_at, s.updated_at, s.score, s.feedback, s.status
			FROM assignment_submissions s
			JOIN users u ON s.student_id = u.id
			WHERE s.id = ?;
		`
	}

	var s models.AssignmentSubmission
	var textAns, filePath, feedback sql.NullString
	var score sql.NullFloat64
	err := r.db.QueryRow(query, id).Scan(
		&s.ID, &s.AssignmentID, &s.StudentID, &s.StudentName, &s.StudentEmail,
		&textAns, &filePath, &s.SubmittedAt, &s.UpdatedAt, &score, &feedback, &s.Status,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if textAns.Valid {
		s.TextAnswer = &textAns.String
	}
	if filePath.Valid {
		s.FilePath = &filePath.String
	}
	if score.Valid {
		s.Score = &score.Float64
	}
	if feedback.Valid {
		s.Feedback = &feedback.String
	}
	return &s, nil
}

func (r *Repository) ListSubmissions(assignmentID int64) ([]models.AssignmentSubmission, error) {
	query := `
		SELECT s.id, s.assignment_id, s.student_id, u.name as student_name, u.email as student_email,
		       s.text_answer, s.file_path, s.submitted_at, s.updated_at, s.score, s.feedback, s.status
		FROM assignment_submissions s
		JOIN users u ON s.student_id = u.id
		WHERE s.assignment_id = $1
		ORDER BY s.submitted_at ASC;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT s.id, s.assignment_id, s.student_id, u.name as student_name, u.email as student_email,
			       s.text_answer, s.file_path, s.submitted_at, s.updated_at, s.score, s.feedback, s.status
			FROM assignment_submissions s
			JOIN users u ON s.student_id = u.id
			WHERE s.assignment_id = ?
			ORDER BY s.submitted_at ASC;
		`
	}

	rows, err := r.db.Query(query, assignmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AssignmentSubmission
	for rows.Next() {
		var s models.AssignmentSubmission
		var textAns, filePath, feedback sql.NullString
		var score sql.NullFloat64
		if err := rows.Scan(
			&s.ID, &s.AssignmentID, &s.StudentID, &s.StudentName, &s.StudentEmail,
			&textAns, &filePath, &s.SubmittedAt, &s.UpdatedAt, &score, &feedback, &s.Status,
		); err != nil {
			return nil, err
		}
		if textAns.Valid {
			s.TextAnswer = &textAns.String
		}
		if filePath.Valid {
			s.FilePath = &filePath.String
		}
		if score.Valid {
			s.Score = &score.Float64
		}
		if feedback.Valid {
			s.Feedback = &feedback.String
		}
		list = append(list, s)
	}
	return list, nil
}

func (r *Repository) GradeSubmission(submissionID int64, score float64, feedback string) error {
	query := `
		UPDATE assignment_submissions
		SET score = $1, feedback = $2, status = 'graded', updated_at = CURRENT_TIMESTAMP
		WHERE id = $3;
	`
	if r.db.Driver != "postgres" {
		query = `
			UPDATE assignment_submissions
			SET score = ?, feedback = ?, status = 'graded', updated_at = CURRENT_TIMESTAMP
			WHERE id = ?;
		`
	}
	_, err := r.db.Exec(query, score, feedback, submissionID)
	return err
}
