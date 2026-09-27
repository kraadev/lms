package quizzes

import "testing"

func TestShuffleQuestions(t *testing.T) {
	orig := []int{1, 2, 3, 4, 5, 6, 7, 8}
	shuffled := ShuffleQuestions(orig)
	if len(shuffled) != len(orig) {
		t.Errorf("length mismatch")
	}
}
