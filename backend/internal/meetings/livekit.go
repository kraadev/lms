package meetings

import (
	"time"

	"github.com/livekit/protocol/auth"

	"lms/internal/config"
	"lms/internal/models"
)

type LiveKitService struct {
	cfg *config.Config
}

func NewLiveKitService(cfg *config.Config) *LiveKitService {
	return &LiveKitService{cfg: cfg}
}

func boolPtr(b bool) *bool {
	return &b
}

func (l *LiveKitService) GenerateToken(roomName string, user *models.User) (string, error) {
	apiKey := l.cfg.LiveKitAPIKey
	apiSecret := l.cfg.LiveKitSecret

	at := auth.NewAccessToken(apiKey, apiSecret)

	grant := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         roomName,
		CanPublish:   boolPtr(true),
		CanSubscribe: boolPtr(true),
	}

	if user.Role == models.RoleTeacher || user.Role == models.RoleAdmin {
		grant.RoomAdmin = true
		grant.RoomCreate = true
	}

	at.AddGrant(grant).
		SetIdentity(user.Email).
		SetName(user.Name).
		SetValidFor(6 * time.Hour)

	return at.ToJWT()
}
