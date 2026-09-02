package users

import (
	"errors"
	"strings"

	"lms/internal/models"
	"lms/internal/utils"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateUserRequest struct {
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Password  string      `json:"password"`
	Role      models.Role `json:"role"`
	AvatarURL *string     `json:"avatar_url,omitempty"`
}

func (s *Service) CreateUser(req CreateUserRequest) (*models.User, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	if req.Email == "" || req.Name == "" || req.Password == "" {
		return nil, errors.New("name, email, and password are required")
	}

	if req.Role != models.RoleAdmin && req.Role != models.RoleTeacher && req.Role != models.RoleStudent {
		return nil, errors.New("invalid role: must be admin, teacher, or student")
	}

	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email is already registered")
	}

	passHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passHash,
		Role:         req.Role,
		AvatarURL:    req.AvatarURL,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUserByID(id int64) (*models.User, error) {
	return s.repo.FindByID(id)
}

type UpdateUserRequest struct {
	Name      *string      `json:"name,omitempty"`
	Email     *string      `json:"email,omitempty"`
	Role      *models.Role `json:"role,omitempty"`
	AvatarURL *string      `json:"avatar_url,omitempty"`
	Password  *string      `json:"password,omitempty"`
}

func (s *Service) UpdateUser(id int64, req UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if req.Name != nil && strings.TrimSpace(*req.Name) != "" {
		user.Name = strings.TrimSpace(*req.Name)
	}
	if req.Email != nil && strings.TrimSpace(*req.Email) != "" {
		user.Email = strings.TrimSpace(strings.ToLower(*req.Email))
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}
	if req.Password != nil && len(*req.Password) >= 6 {
		passHash, err := utils.HashPassword(*req.Password)
		if err == nil {
			user.PasswordHash = passHash
		}
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) ListUsers(role, search string) ([]models.User, error) {
	return s.repo.List(role, search)
}

func (s *Service) DeleteUser(id int64) error {
	return s.repo.Delete(id)
}
