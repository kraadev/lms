package assignments

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

type CreateAssignmentRequest struct {
	ClassID        int64     `json:"class_id"`
	TeacherID      int64     `json:"teacher_id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	AttachmentPath *string   `json:"attachment_path,omitempty"`
	Deadline       time.Time `json:"deadline"`
	MaxScore       float64   `json:"max_score"`
}

func (s *Service) CreateAssignment(req CreateAssignmentRequest) (*models.Assignment, error) {
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return nil, errors.New("assignment title is required")
	}
	if req.MaxScore <= 0 {
		req.MaxScore = 100.0
	}
	if req.Deadline.IsZero() {
		return nil, errors.New("valid deadline is required")
	}

	a := &models.Assignment{
		ClassID:        req.ClassID,
		TeacherID:      req.TeacherID,
		Title:          req.Title,
		Description:    req.Description,
		AttachmentPath: req.AttachmentPath,
		Deadline:       req.Deadline,
		MaxScore:       req.MaxScore,
	}

	if err := s.repo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) ListByClass(classID int64) ([]models.Assignment, error) {
	return s.repo.ListByClass(classID)
}

func (s *Service) GetByID(id int64) (*models.Assignment, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Update(a *models.Assignment) error {
	return s.repo.Update(a)
}

func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *Service) SubmitAssignment(assignmentID, studentID int64, textAnswer, filePath *string) (*models.AssignmentSubmission, error) {
	assignment, err := s.repo.FindByID(assignmentID)
	if err != nil || assignment == nil {
		return nil, errors.New("assignment not found")
	}

	status := "submitted"
	if time.Now().After(assignment.Deadline) {
		status = "late"
	}

	sub := &models.AssignmentSubmission{
		AssignmentID: assignmentID,
		StudentID:    studentID,
		TextAnswer:   textAnswer,
		FilePath:     filePath,
		Status:       status,
	}

	if err := s.repo.UpsertSubmission(sub); err != nil {
		return nil, err
	}

	return s.repo.GetSubmission(assignmentID, studentID)
}

func (s *Service) GetStudentSubmission(assignmentID, studentID int64) (*models.AssignmentSubmission, error) {
	return s.repo.GetSubmission(assignmentID, studentID)
}

func (s *Service) ListSubmissions(assignmentID int64) ([]models.AssignmentSubmission, error) {
	return s.repo.ListSubmissions(assignmentID)
}

func (s *Service) GradeSubmission(submissionID int64, score float64, feedback string) (*models.AssignmentSubmission, error) {
	sub, err := s.repo.GetSubmissionByID(submissionID)
	if err != nil || sub == nil {
		return nil, errors.New("submission not found")
	}

	assignment, err := s.repo.FindByID(sub.AssignmentID)
	if err != nil || assignment == nil {
		return nil, errors.New("associated assignment not found")
	}

	if score < 0 || score > assignment.MaxScore {
		return nil, errors.New("score must be between 0 and maximum assignment score")
	}

	if err := s.repo.GradeSubmission(submissionID, score, feedback); err != nil {
		return nil, err
	}

	return s.repo.GetSubmissionByID(submissionID)
}
