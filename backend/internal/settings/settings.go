package settings

import (
	"strconv"
	"sync"
)

type Store struct {
	mu    sync.RWMutex
	cache map[string]string
}

func NewStore() *Store {
	return &Store{
		cache: map[string]string{
			"app_name":       "LMS Portal",
			"max_upload_mb":  "50",
			"allow_register": "true",
			"theme_default":  "system",
		},
	}
}

func (s *Store) GetString(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cache[key]
}

func (s *Store) GetInt(key string, def int) int {
	val := s.GetString(key)
	if n, err := strconv.Atoi(val); err == nil {
		return n
	}
	return def
}

func (s *Store) GetBool(key string) bool {
	val := s.GetString(key)
	return val == "true" || val == "1"
}

func (s *Store) Set(key, val string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[key] = val
}
