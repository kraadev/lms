package quizzes

import (
	"errors"
	"time"
)

type SubmittedAnswer struct {
	QuestionID int64 `json:"question_id"`
	OptionID   int64 `json:"option_id"`
}

type QuestionKey struct {
	ID              int64
	CorrectOptionID int64
	Points          int
}

type EvaluationResult struct {
	TotalQuestions int
	CorrectAnswers int
	EarnedPoints   int
	MaxPoints      int
	Percentage     float64
	IsLate         bool
}

// EvaluateSubmission grades quiz answers server-side with strict key isolation and timer verification
func EvaluateSubmission(startedAt time.Time, durationMinutes int, keys []QuestionKey, userAnswers []SubmittedAnswer) (*EvaluationResult, error) {
	if len(keys) == 0 {
		return nil, errors.New("no question keys provided for evaluation")
	}

	maxAllowedTime := startedAt.Add(time.Duration(durationMinutes+2) * time.Minute) // 2-min grace period for network jitter
	isLate := time.Now().After(maxAllowedTime)

	keyMap := make(map[int64]QuestionKey, len(keys))
	maxPoints := 0
	for _, k := range keys {
		keyMap[k.ID] = k
		maxPoints += k.Points
	}

	correctCount := 0
	earnedPoints := 0
	for _, ans := range userAnswers {
		if k, ok := keyMap[ans.QuestionID]; ok {
			if ans.OptionID == k.CorrectOptionID {
				correctCount++
				earnedPoints += k.Points
			}
		}
	}

	pct := 0.0
	if maxPoints > 0 {
		pct = (float64(earnedPoints) / float64(maxPoints)) * 100.0
	}

	return &EvaluationResult{
		TotalQuestions: len(keys),
		CorrectAnswers: correctCount,
		EarnedPoints:   earnedPoints,
		MaxPoints:      maxPoints,
		Percentage:     pct,
		IsLate:         isLate,
	}, nil
}
