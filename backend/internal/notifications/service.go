package notifications

import (
	"lms/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) NotifyUser(userID int64, notifType, title, message string) error {
	n := &models.Notification{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Message: message,
	}
	return s.repo.Create(n)
}

func (s *Service) ListUserNotifications(userID int64) ([]models.Notification, error) {
	return s.repo.ListByUser(userID)
}

func (s *Service) MarkAsRead(id, userID int64) error {
	return s.repo.MarkAsRead(id, userID)
}

func (s *Service) MarkAllRead(userID int64) error {
	return s.repo.MarkAllRead(userID)
}
