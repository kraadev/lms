package config

import (
	"sync"
)

type FeatureFlags struct {
	mu    sync.RWMutex
	flags map[string]bool
}

func NewFeatureFlags() *FeatureFlags {
	return &FeatureFlags{
		flags: map[string]bool{
			"video_meetings": true,
			"chat_reactions": true,
			"ai_grader":      false,
			"audit_logging":  true,
		},
	}
}

func (ff *FeatureFlags) IsEnabled(key string) bool {
	ff.mu.RLock()
	defer ff.mu.RUnlock()
	return ff.flags[key]
}

func (ff *FeatureFlags) Set(key string, enabled bool) {
	ff.mu.Lock()
	defer ff.mu.Unlock()
	ff.flags[key] = enabled
}

func (ff *FeatureFlags) All() map[string]bool {
	ff.mu.RLock()
	defer ff.mu.RUnlock()
	copyMap := make(map[string]bool, len(ff.flags))
	for k, v := range ff.flags {
		copyMap[k] = v
	}
	return copyMap
}
