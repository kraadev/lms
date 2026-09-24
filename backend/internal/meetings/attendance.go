package meetings

import (
	"sync"
	"time"
)

type SessionAttendance struct {
	UserID       int64     `json:"user_id"`
	MeetingID    int64     `json:"meeting_id"`
	JoinedAt     time.Time `json:"joined_at"`
	LeftAt       time.Time `json:"left_at"`
	DurationSec  int       `json:"duration_sec"`
}

type AttendanceTracker struct {
	mu       sync.Mutex
	sessions map[string]time.Time
}

func NewAttendanceTracker() *AttendanceTracker {
	return &AttendanceTracker{
		sessions: make(map[string]time.Time),
	}
}

func (t *AttendanceTracker) sessionKey(meetingID, userID int64) string {
	return time.Now().Format("20060102") + "_" + string(rune(meetingID)) + "_" + string(rune(userID))
}

func (t *AttendanceTracker) RecordJoin(meetingID, userID int64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.sessions[t.sessionKey(meetingID, userID)] = time.Now()
}

func (t *AttendanceTracker) RecordLeave(meetingID, userID int64) *SessionAttendance {
	t.mu.Lock()
	defer t.mu.Unlock()

	key := t.sessionKey(meetingID, userID)
	joinedAt, exists := t.sessions[key]
	if !exists {
		return nil
	}
	delete(t.sessions, key)

	leftAt := time.Now()
	dur := int(leftAt.Sub(joinedAt).Seconds())

	return &SessionAttendance{
		UserID:      userID,
		MeetingID:   meetingID,
		JoinedAt:    joinedAt,
		LeftAt:      leftAt,
		DurationSec: dur,
	}
}
