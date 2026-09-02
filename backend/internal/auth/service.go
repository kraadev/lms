package auth

import (
	"errors"
	"strings"

	"lms/internal/config"
	"lms/internal/models"
	"lms/internal/users"
	"lms/internal/utils"
)

type Service struct {
	userRepo *users.Repository
	cfg      *config.Config
}

func NewService(userRepo *users.Repository, cfg *config.Config) *Service {
	return &Service{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

type LoginResult struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (s *Service) Login(email, password string) (*LoginResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, string(user.Role), user.Name, s.cfg.JWTSecret, s.cfg.JWTExpiryHours)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Token: token,
		User:  user,
	}, nil
}

func (s *Service) GetMe(userID int64) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}
