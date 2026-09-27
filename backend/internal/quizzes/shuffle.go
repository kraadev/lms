package quizzes

import (
	"math/rand"
	"time"
)

func ShuffleQuestions[T any](items []T) []T {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]T, len(items))
	copy(result, items)
	r.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return result
}
