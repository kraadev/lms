package meetings

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"lms/internal/config"
	"lms/internal/models"
)

type Service struct {
	repo       *Repository
	livekit    *LiveKitService
	cfg        *config.Config
}

func NewService(repo *Repository, livekit *LiveKitService, cfg *config.Config) *Service {
	return &Service{
		repo:    repo,
		livekit: livekit,
		cfg:     cfg,
	}
}

type CreateMeetingRequest struct {
	ClassID   int64              `json:"class_id"`
	TeacherID int64              `json:"teacher_id"`
	Title     string             `json:"title"`
	Type      models.MeetingType `json:"type"`
}

func (s *Service) CreateMeeting(req CreateMeetingRequest) (*models.Meeting, error) {
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		req.Title = "Kelas Online"
	}
	if req.Type != models.MeetingTypeAudio && req.Type != models.MeetingTypeVideo {
		req.Type = models.MeetingTypeVideo
	}

	roomName := fmt.Sprintf("class-%d-room-%s", req.ClassID, uuid.New().String()[:8])

	meeting := &models.Meeting{
		ClassID:   req.ClassID,
		TeacherID: req.TeacherID,
		Title:     req.Title,
		RoomName:  roomName,
		Type:      req.Type,
		Status:    models.MeetingStatusActive,
	}

	if err := s.repo.Create(meeting); err != nil {
		return nil, err
	}

	return meeting, nil
}

type JoinMeetingResponse struct {
	URL     string          `json:"url"`
	Token   string          `json:"token"`
	Meeting *models.Meeting `json:"meeting"`
}

func (s *Service) JoinMeeting(meetingID int64, user *models.User) (*JoinMeetingResponse, error) {
	meeting, err := s.repo.FindByID(meetingID)
	if err != nil || meeting == nil {
		return nil, errors.New("meeting not found")
	}

	if meeting.Status != models.MeetingStatusActive {
		return nil, errors.New("meeting is no longer active")
	}

	token, err := s.livekit.GenerateToken(meeting.RoomName, user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate meeting token: %w", err)
	}

	return &JoinMeetingResponse{
		URL:     s.cfg.LiveKitURL,
		Token:   token,
		Meeting: meeting,
	}, nil
}

func (s *Service) EndMeeting(meetingID int64) error {
	return s.repo.EndMeeting(meetingID)
}

func (s *Service) GetByID(meetingID int64) (*models.Meeting, error) {
	return s.repo.FindByID(meetingID)
}

func (s *Service) ListByClass(classID int64) ([]models.Meeting, error) {
	return s.repo.ListByClass(classID)
}
