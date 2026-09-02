package materials

import (
	"errors"
	"strings"

	"lms/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateMaterialRequest struct {
	ClassID     int64   `json:"class_id"`
	TeacherID   int64   `json:"teacher_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Content     string  `json:"content"`
	FilePath    *string `json:"file_path,omitempty"`
	ExternalURL *string `json:"external_url,omitempty"`
}

func (s *Service) CreateMaterial(req CreateMaterialRequest) (*models.Material, error) {
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return nil, errors.New("material title is required")
	}

	m := &models.Material{
		ClassID:     req.ClassID,
		TeacherID:   req.TeacherID,
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		FilePath:    req.FilePath,
		ExternalURL: req.ExternalURL,
	}

	if err := s.repo.Create(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) ListByClass(classID int64) ([]models.Material, error) {
	return s.repo.ListByClass(classID)
}

func (s *Service) GetByID(id int64) (*models.Material, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Update(m *models.Material) error {
	return s.repo.Update(m)
}

func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}
