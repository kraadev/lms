package quizzes

import (
	"database/sql"
	"time"

	"lms/internal/database"
	"lms/internal/models"
	"lms/internal/utils"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListByClass(classID int64) ([]models.Quiz, error) {
	query := `
		SELECT id, class_id, teacher_id, title, description, duration_minutes, start_at, end_at, max_attempts, status, created_at, updated_at
		FROM quizzes
		WHERE class_id = $1
		ORDER BY id DESC;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT id, class_id, teacher_id, title, description, duration_minutes, start_at, end_at, max_attempts, status, created_at, updated_at
			FROM quizzes
			WHERE class_id = ?
			ORDER BY id DESC;
		`
	}

	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Quiz
	for rows.Next() {
		var q models.Quiz
		var startRaw, endRaw, createdRaw, updatedRaw interface{}
		if err := rows.Scan(
			&q.ID, &q.ClassID, &q.TeacherID, &q.Title, &q.Description,
			&q.DurationMinutes, &startRaw, &endRaw, &q.MaxAttempts,
			&q.Status, &createdRaw, &updatedRaw,
		); err != nil {
			return nil, err
		}
		q.StartAt = utils.ParseTime(startRaw)
		q.EndAt = utils.ParseTime(endRaw)
		q.CreatedAt = utils.ParseTime(createdRaw)
		q.UpdatedAt = utils.ParseTime(updatedRaw)
		list = append(list, q)
	}
	return list, nil
}

func (r *Repository) FindByID(id int64) (*models.Quiz, error) {
	query := `
		SELECT id, class_id, teacher_id, title, description, duration_minutes, start_at, end_at, max_attempts, status, created_at, updated_at
		FROM quizzes
		WHERE id = $1;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT id, class_id, teacher_id, title, description, duration_minutes, start_at, end_at, max_attempts, status, created_at, updated_at
			FROM quizzes
			WHERE id = ?;
		`
	}

	var q models.Quiz
	var startRaw, endRaw, createdRaw, updatedRaw interface{}
	err := r.db.QueryRow(query, id).Scan(
		&q.ID, &q.ClassID, &q.TeacherID, &q.Title, &q.Description,
		&q.DurationMinutes, &startRaw, &endRaw, &q.MaxAttempts,
		&q.Status, &createdRaw, &updatedRaw,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	q.StartAt = utils.ParseTime(startRaw)
	q.EndAt = utils.ParseTime(endRaw)
	q.CreatedAt = utils.ParseTime(createdRaw)
	q.UpdatedAt = utils.ParseTime(updatedRaw)

	// Load questions & options
	questions, err := r.GetQuestionsByQuizID(id)
	if err != nil {
		return nil, err
	}
	q.Questions = questions

	return &q, nil
}

func (r *Repository) GetQuestionsByQuizID(quizID int64) ([]models.QuizQuestion, error) {
	query := `
		SELECT id, quiz_id, question, type, points, order_index
		FROM quiz_questions
		WHERE quiz_id = $1
		ORDER BY order_index ASC, id ASC;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT id, quiz_id, question, type, points, order_index
			FROM quiz_questions
			WHERE quiz_id = ?
			ORDER BY order_index ASC, id ASC;
		`
	}

	rows, err := r.db.Query(query, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.QuizQuestion
	for rows.Next() {
		var qq models.QuizQuestion
		var pointsRaw interface{}
		if err := rows.Scan(&qq.ID, &qq.QuizID, &qq.Question, &qq.Type, &pointsRaw, &qq.OrderIndex); err != nil {
			return nil, err
		}
		qq.Points = utils.ParseFloat(pointsRaw)

		options, err := r.GetOptionsByQuestionID(qq.ID)
		if err != nil {
			return nil, err
		}
		qq.Options = options
		questions = append(questions, qq)
	}

	return questions, nil
}

func (r *Repository) GetOptionsByQuestionID(questionID int64) ([]models.QuizOption, error) {
	query := `
		SELECT id, question_id, option_text, is_correct
		FROM quiz_options
		WHERE question_id = $1
		ORDER BY id ASC;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT id, question_id, option_text, is_correct
			FROM quiz_options
			WHERE question_id = ?
			ORDER BY id ASC;
		`
	}

	rows, err := r.db.Query(query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []models.QuizOption
	for rows.Next() {
		var opt models.QuizOption
		var isCorrectRaw interface{}
		if err := rows.Scan(&opt.ID, &opt.QuestionID, &opt.OptionText, &isCorrectRaw); err != nil {
			return nil, err
		}
		isCorrect := utils.ParseBool(isCorrectRaw)
		opt.IsCorrect = &isCorrect
		options = append(options, opt)
	}
	return options, nil
}

func (r *Repository) CreateQuiz(q *models.Quiz) error {
	if r.db.Driver == "postgres" {
		query := `
			INSERT INTO quizzes (class_id, teacher_id, title, description, duration_minutes, start_at, end_at, max_attempts, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			RETURNING id, created_at, updated_at;
		`
		return r.db.QueryRow(query, q.ClassID, q.TeacherID, q.Title, q.Description, q.DurationMinutes, q.StartAt, q.EndAt, q.MaxAttempts, q.Status).Scan(&q.ID, &q.CreatedAt, &q.UpdatedAt)
	}

	query := `
		INSERT INTO quizzes (class_id, teacher_id, title, description, duration_minutes, start_at, end_at, max_attempts, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
	`
	res, err := r.db.Exec(query, q.ClassID, q.TeacherID, q.Title, q.Description, q.DurationMinutes, q.StartAt, q.EndAt, q.MaxAttempts, q.Status)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	q.ID = id
	return nil
}

func (r *Repository) CreateQuestion(qq *models.QuizQuestion) error {
	if r.db.Driver == "postgres" {
		query := `
			INSERT INTO quiz_questions (quiz_id, question, type, points, order_index)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id;
		`
		return r.db.QueryRow(query, qq.QuizID, qq.Question, qq.Type, qq.Points, qq.OrderIndex).Scan(&qq.ID)
	}

	query := `
		INSERT INTO quiz_questions (quiz_id, question, type, points, order_index)
		VALUES (?, ?, ?, ?, ?);
	`
	res, err := r.db.Exec(query, qq.QuizID, qq.Question, qq.Type, qq.Points, qq.OrderIndex)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	qq.ID = id
	return nil
}

func (r *Repository) CreateOption(opt *models.QuizOption) error {
	isCorrect := false
	if opt.IsCorrect != nil {
		isCorrect = *opt.IsCorrect
	}

	if r.db.Driver == "postgres" {
		query := `
			INSERT INTO quiz_options (question_id, option_text, is_correct)
			VALUES ($1, $2, $3)
			RETURNING id;
		`
		return r.db.QueryRow(query, opt.QuestionID, opt.OptionText, isCorrect).Scan(&opt.ID)
	}

	query := `
		INSERT INTO quiz_options (question_id, option_text, is_correct)
		VALUES (?, ?, ?);
	`
	res, err := r.db.Exec(query, opt.QuestionID, opt.OptionText, isCorrect)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	opt.ID = id
	return nil
}

func (r *Repository) UpdateQuiz(q *models.Quiz) error {
	query := `
		UPDATE quizzes
		SET title = $1, description = $2, duration_minutes = $3, start_at = $4, end_at = $5, max_attempts = $6, status = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8;
	`
	if r.db.Driver != "postgres" {
		query = `
			UPDATE quizzes
			SET title = ?, description = ?, duration_minutes = ?, start_at = ?, end_at = ?, max_attempts = ?, status = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?;
		`
	}
	_, err := r.db.Exec(query, q.Title, q.Description, q.DurationMinutes, q.StartAt, q.EndAt, q.MaxAttempts, q.Status, q.ID)
	return err
}

func (r *Repository) DeleteQuiz(id int64) error {
	query := "DELETE FROM quizzes WHERE id = $1;"
	if r.db.Driver != "postgres" {
		query = "DELETE FROM quizzes WHERE id = ?;"
	}
	_, err := r.db.Exec(query, id)
	return err
}

func (r *Repository) CountAttempts(quizID, studentID int64) (int, error) {
	query := "SELECT COUNT(*) FROM quiz_attempts WHERE quiz_id = $1 AND student_id = $2;"
	if r.db.Driver != "postgres" {
		query = "SELECT COUNT(*) FROM quiz_attempts WHERE quiz_id = ? AND student_id = ?;"
	}
	var count int
	err := r.db.QueryRow(query, quizID, studentID).Scan(&count)
	return count, err
}

func (r *Repository) CreateAttempt(qa *models.QuizAttempt) error {
	if r.db.Driver == "postgres" {
		query := `
			INSERT INTO quiz_attempts (quiz_id, student_id, started_at, status)
			VALUES ($1, $2, CURRENT_TIMESTAMP, 'in_progress')
			RETURNING id, started_at;
		`
		return r.db.QueryRow(query, qa.QuizID, qa.StudentID).Scan(&qa.ID, &qa.StartedAt)
	}

	query := `
		INSERT INTO quiz_attempts (quiz_id, student_id, started_at, status)
		VALUES (?, ?, CURRENT_TIMESTAMP, 'in_progress');
	`
	res, err := r.db.Exec(query, qa.QuizID, qa.StudentID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	qa.ID = id
	qa.StartedAt = time.Now()
	return nil
}

func (r *Repository) GetAttemptByID(id int64) (*models.QuizAttempt, error) {
	query := `
		SELECT qa.id, qa.quiz_id, qa.student_id, u.name as student_name, qa.started_at, qa.submitted_at, qa.score, qa.status
		FROM quiz_attempts qa
		JOIN users u ON qa.student_id = u.id
		WHERE qa.id = $1;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT qa.id, qa.quiz_id, qa.student_id, u.name as student_name, qa.started_at, qa.submitted_at, qa.score, qa.status
			FROM quiz_attempts qa
			JOIN users u ON qa.student_id = u.id
			WHERE qa.id = ?;
		`
	}

	var qa models.QuizAttempt
	var startRaw, subRaw, scoreRaw interface{}
	err := r.db.QueryRow(query, id).Scan(
		&qa.ID, &qa.QuizID, &qa.StudentID, &qa.StudentName, &startRaw, &subRaw, &scoreRaw, &qa.Status,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	qa.StartedAt = utils.ParseTime(startRaw)
	if subRaw != nil {
		t := utils.ParseTime(subRaw)
		if !t.IsZero() {
			qa.SubmittedAt = &t
		}
	}
	if scoreRaw != nil {
		if sc, ok := scoreRaw.(float64); ok {
			qa.Score = &sc
		}
	}

	answers, err := r.GetAnswersByAttempt(id)
	if err != nil {
		return nil, err
	}
	qa.Answers = answers

	return &qa, nil
}

func (r *Repository) SaveAnswer(ans *models.QuizAnswer) error {
	if r.db.Driver == "postgres" {
		query := `
			INSERT INTO quiz_answers (attempt_id, question_id, selected_option_id, text_answer, is_correct, earned_points)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id;
		`
		return r.db.QueryRow(query, ans.AttemptID, ans.QuestionID, ans.SelectedOptionID, ans.TextAnswer, ans.IsCorrect, ans.EarnedPoints).Scan(&ans.ID)
	}

	query := `
		INSERT INTO quiz_answers (attempt_id, question_id, selected_option_id, text_answer, is_correct, earned_points)
		VALUES (?, ?, ?, ?, ?, ?);
	`
	res, err := r.db.Exec(query, ans.AttemptID, ans.QuestionID, ans.SelectedOptionID, ans.TextAnswer, ans.IsCorrect, ans.EarnedPoints)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	ans.ID = id
	return nil
}

func (r *Repository) CompleteAttempt(attemptID int64, score float64, status string) error {
	query := `
		UPDATE quiz_attempts
		SET score = $1, status = $2, submitted_at = CURRENT_TIMESTAMP
		WHERE id = $3;
	`
	if r.db.Driver != "postgres" {
		query = `
			UPDATE quiz_attempts
			SET score = ?, status = ?, submitted_at = CURRENT_TIMESTAMP
			WHERE id = ?;
		`
	}
	_, err := r.db.Exec(query, score, status, attemptID)
	return err
}

func (r *Repository) GetAnswersByAttempt(attemptID int64) ([]models.QuizAnswer, error) {
	query := `
		SELECT id, attempt_id, question_id, selected_option_id, text_answer, is_correct, earned_points
		FROM quiz_answers
		WHERE attempt_id = $1
		ORDER BY id ASC;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT id, attempt_id, question_id, selected_option_id, text_answer, is_correct, earned_points
			FROM quiz_answers
			WHERE attempt_id = ?
			ORDER BY id ASC;
		`
	}

	rows, err := r.db.Query(query, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.QuizAnswer
	for rows.Next() {
		var a models.QuizAnswer
		var optID sql.NullInt64
		var txtAns sql.NullString
		var isCorrRaw, pointsRaw interface{}
		if err := rows.Scan(&a.ID, &a.AttemptID, &a.QuestionID, &optID, &txtAns, &isCorrRaw, &pointsRaw); err != nil {
			return nil, err
		}
		if optID.Valid {
			a.SelectedOptionID = &optID.Int64
		}
		if txtAns.Valid {
			a.TextAnswer = &txtAns.String
		}
		if isCorrRaw != nil {
			b := utils.ParseBool(isCorrRaw)
			a.IsCorrect = &b
		}
		if pointsRaw != nil {
			if pt, ok := pointsRaw.(float64); ok {
				a.EarnedPoints = &pt
			}
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *Repository) ListAttemptsByQuiz(quizID int64) ([]models.QuizAttempt, error) {
	query := `
		SELECT qa.id, qa.quiz_id, qa.student_id, u.name as student_name, qa.started_at, qa.submitted_at, qa.score, qa.status
		FROM quiz_attempts qa
		JOIN users u ON qa.student_id = u.id
		WHERE qa.quiz_id = $1
		ORDER BY qa.started_at DESC;
	`
	if r.db.Driver != "postgres" {
		query = `
			SELECT qa.id, qa.quiz_id, qa.student_id, u.name as student_name, qa.started_at, qa.submitted_at, qa.score, qa.status
			FROM quiz_attempts qa
			JOIN users u ON qa.student_id = u.id
			WHERE qa.quiz_id = ?
			ORDER BY qa.started_at DESC;
		`
	}

	rows, err := r.db.Query(query, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.QuizAttempt
	for rows.Next() {
		var qa models.QuizAttempt
		var startRaw, subRaw, scoreRaw interface{}
		if err := rows.Scan(
			&qa.ID, &qa.QuizID, &qa.StudentID, &qa.StudentName, &startRaw, &subRaw, &scoreRaw, &qa.Status,
		); err != nil {
			return nil, err
		}
		qa.StartedAt = utils.ParseTime(startRaw)
		if subRaw != nil {
			t := utils.ParseTime(subRaw)
			if !t.IsZero() {
				qa.SubmittedAt = &t
			}
		}
		if scoreRaw != nil {
			if sc, ok := scoreRaw.(float64); ok {
				qa.Score = &sc
			}
		}
		list = append(list, qa)
	}
	return list, nil
}
