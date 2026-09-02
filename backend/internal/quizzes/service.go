package quizzes

import (
	"errors"
	"strings"
	"time"

	"lms/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateQuizRequest struct {
	ClassID         int64                         `json:"class_id"`
	TeacherID       int64                         `json:"teacher_id"`
	Title           string                        `json:"title"`
	Description     string                        `json:"description"`
	DurationMinutes int                           `json:"duration_minutes"`
	StartAt         time.Time                     `json:"start_at"`
	EndAt           time.Time                     `json:"end_at"`
	MaxAttempts     int                           `json:"max_attempts"`
	Questions       []CreateQuizQuestionRequest   `json:"questions"`
}

type CreateQuizQuestionRequest struct {
	Question   string                    `json:"question"`
	Type       models.QuestionType       `json:"type"`
	Points     float64                   `json:"points"`
	OrderIndex int                       `json:"order_index"`
	Options    []CreateQuizOptionRequest `json:"options"`
}

type CreateQuizOptionRequest struct {
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct"`
}

func (s *Service) CreateQuiz(req CreateQuizRequest) (*models.Quiz, error) {
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return nil, errors.New("quiz title is required")
	}
	if req.DurationMinutes <= 0 {
		req.DurationMinutes = 60
	}
	if req.MaxAttempts <= 0 {
		req.MaxAttempts = 1
	}
	if req.EndAt.Before(req.StartAt) {
		return nil, errors.New("end_at must be after start_at")
	}

	quiz := &models.Quiz{
		ClassID:         req.ClassID,
		TeacherID:       req.TeacherID,
		Title:           req.Title,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		StartAt:         req.StartAt,
		EndAt:           req.EndAt,
		MaxAttempts:     req.MaxAttempts,
		Status:          "published",
	}

	if err := s.repo.CreateQuiz(quiz); err != nil {
		return nil, err
	}

	for _, qReq := range req.Questions {
		q := &models.QuizQuestion{
			QuizID:     quiz.ID,
			Question:   qReq.Question,
			Type:       qReq.Type,
			Points:     qReq.Points,
			OrderIndex: qReq.OrderIndex,
		}
		if q.Points <= 0 {
			q.Points = 10.0
		}
		if err := s.repo.CreateQuestion(q); err != nil {
			return nil, err
		}

		for _, optReq := range qReq.Options {
			isCorrect := optReq.IsCorrect
			opt := &models.QuizOption{
				QuestionID: q.ID,
				OptionText: optReq.OptionText,
				IsCorrect:  &isCorrect,
			}
			if err := s.repo.CreateOption(opt); err != nil {
				return nil, err
			}
		}
	}

	return s.repo.FindByID(quiz.ID)
}

func (s *Service) ListByClass(classID int64) ([]models.Quiz, error) {
	return s.repo.ListByClass(classID)
}

func (s *Service) GetByIDForTeacher(id int64) (*models.Quiz, error) {
	return s.repo.FindByID(id)
}

// GetByIDForStudent sanitizes is_correct from options to strictly prevent leaking answer keys
func (s *Service) GetByIDForStudent(id int64) (*models.Quiz, error) {
	quiz, err := s.repo.FindByID(id)
	if err != nil || quiz == nil {
		return nil, err
	}

	// Sanitize answer keys
	sanitizedQuestions := make([]models.QuizQuestion, len(quiz.Questions))
	for i, q := range quiz.Questions {
		sanitizedOptions := make([]models.QuizOption, len(q.Options))
		for j, opt := range q.Options {
			sanitizedOptions[j] = models.QuizOption{
				ID:         opt.ID,
				QuestionID: opt.QuestionID,
				OptionText: opt.OptionText,
				IsCorrect:  nil, // Mask is_correct
			}
		}
		sanitizedQuestions[i] = q
		sanitizedQuestions[i].Options = sanitizedOptions
	}
	quiz.Questions = sanitizedQuestions

	return quiz, nil
}

func (s *Service) StartAttempt(quizID, studentID int64) (*models.QuizAttempt, error) {
	quiz, err := s.repo.FindByID(quizID)
	if err != nil || quiz == nil {
		return nil, errors.New("quiz not found")
	}

	now := time.Now()
	if now.Before(quiz.StartAt) {
		return nil, errors.New("quiz has not started yet")
	}
	if now.After(quiz.EndAt) {
		return nil, errors.New("quiz schedule has expired")
	}

	count, err := s.repo.CountAttempts(quizID, studentID)
	if err != nil {
		return nil, err
	}
	if count >= quiz.MaxAttempts {
		return nil, errors.New("maximum quiz attempts exceeded")
	}

	attempt := &models.QuizAttempt{
		QuizID:    quizID,
		StudentID: studentID,
		Status:    "in_progress",
	}

	if err := s.repo.CreateAttempt(attempt); err != nil {
		return nil, err
	}

	return attempt, nil
}

type SubmitAnswerItem struct {
	QuestionID       int64   `json:"question_id"`
	SelectedOptionID *int64  `json:"selected_option_id,omitempty"`
	TextAnswer       *string `json:"text_answer,omitempty"`
}

type SubmitAttemptRequest struct {
	Answers []SubmitAnswerItem `json:"answers"`
}

func (s *Service) SubmitAttempt(attemptID, studentID int64, req SubmitAttemptRequest) (*models.QuizAttempt, error) {
	attempt, err := s.repo.GetAttemptByID(attemptID)
	if err != nil || attempt == nil {
		return nil, errors.New("attempt not found")
	}
	if attempt.StudentID != studentID {
		return nil, errors.New("unauthorized attempt submission")
	}
	if attempt.Status != "in_progress" {
		return nil, errors.New("attempt has already been submitted or completed")
	}

	quiz, err := s.repo.FindByID(attempt.QuizID)
	if err != nil || quiz == nil {
		return nil, errors.New("associated quiz not found")
	}

	// Map questions and options for quick grading
	type questionMeta struct {
		points    float64
		qType     models.QuestionType
		correctID int64
	}
	qMap := make(map[int64]questionMeta)

	for _, q := range quiz.Questions {
		var correctOptID int64
		for _, opt := range q.Options {
			if opt.IsCorrect != nil && *opt.IsCorrect {
				correctOptID = opt.ID
				break
			}
		}
		qMap[q.ID] = questionMeta{
			points:    q.Points,
			qType:     q.Type,
			correctID: correctOptID,
		}
	}

	var totalScore float64

	for _, ansItem := range req.Answers {
		meta, exists := qMap[ansItem.QuestionID]
		if !exists {
			continue
		}

		isCorrect := false
		var earnedPoints float64

		if (meta.qType == models.QuestionMultipleChoice || meta.qType == models.QuestionTrueFalse) && ansItem.SelectedOptionID != nil {
			if *ansItem.SelectedOptionID == meta.correctID {
				isCorrect = true
				earnedPoints = meta.points
				totalScore += earnedPoints
			}
		}

		ansRecord := &models.QuizAnswer{
			AttemptID:        attemptID,
			QuestionID:       ansItem.QuestionID,
			SelectedOptionID: ansItem.SelectedOptionID,
			TextAnswer:       ansItem.TextAnswer,
			IsCorrect:        &isCorrect,
			EarnedPoints:     &earnedPoints,
		}

		_ = s.repo.SaveAnswer(ansRecord)
	}

	if err := s.repo.CompleteAttempt(attemptID, totalScore, "graded"); err != nil {
		return nil, err
	}

	return s.repo.GetAttemptByID(attemptID)
}

func (s *Service) GetAttemptByID(attemptID int64) (*models.QuizAttempt, error) {
	return s.repo.GetAttemptByID(attemptID)
}

func (s *Service) ListAttemptsByQuiz(quizID int64) ([]models.QuizAttempt, error) {
	return s.repo.ListAttemptsByQuiz(quizID)
}

func (s *Service) DeleteQuiz(id int64) error {
	return s.repo.DeleteQuiz(id)
}
