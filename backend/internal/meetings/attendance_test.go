package meetings

import (
	"testing"
	"time"
)

func TestAttendanceTracker(t *testing.T) {
	tracker := NewAttendanceTracker()

	meetingID := int64(10)
	userID := int64(42)

	tracker.RecordJoin(meetingID, userID)
	time.Sleep(50 * time.Millisecond)

	rec := tracker.RecordLeave(meetingID, userID)
	if rec == nil {
		t.Fatalf("expected attendance record, got nil")
	}

	if rec.MeetingID != meetingID || rec.UserID != userID {
		t.Errorf("mismatched record identifiers")
	}

	if rec.DurationSec < 0 {
		t.Errorf("expected non-negative duration")
	}
}
