package quizzes

import (
	"testing"
	"time"
)

func TestEvaluateSubmission(t *testing.T) {
	keys := []QuestionKey{
		{ID: 1, CorrectOptionID: 101, Points: 10},
		{ID: 2, CorrectOptionID: 202, Points: 10},
		{ID: 3, CorrectOptionID: 303, Points: 10},
	}

	answers := []SubmittedAnswer{
		{QuestionID: 1, OptionID: 101}, // Correct
		{QuestionID: 2, OptionID: 999}, // Wrong
		{QuestionID: 3, OptionID: 303}, // Correct
	}

	startedAt := time.Now().Add(-10 * time.Minute)
	durationMinutes := 30

	res, err := EvaluateSubmission(startedAt, durationMinutes, keys, answers)
	if err != nil {
		t.Fatalf("unexpected evaluation error: %v", err)
	}

	if res.CorrectAnswers != 2 {
		t.Errorf("expected 2 correct answers, got %d", res.CorrectAnswers)
	}
	if res.EarnedPoints != 20 || res.MaxPoints != 30 {
		t.Errorf("expected 20/30 points, got %d/%d", res.EarnedPoints, res.MaxPoints)
	}
	if res.IsLate {
		t.Errorf("expected submission to be within valid time window")
	}
}
