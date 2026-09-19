package config

import "testing"

func TestFeatureFlags(t *testing.T) {
	ff := NewFeatureFlags()

	if !ff.IsEnabled("video_meetings") {
		t.Errorf("expected video_meetings to be enabled by default")
	}

	ff.Set("video_meetings", false)
	if ff.IsEnabled("video_meetings") {
		t.Errorf("expected video_meetings to be disabled after Set")
	}

	all := ff.All()
	if _, exists := all["chat_reactions"]; !exists {
		t.Errorf("expected chat_reactions flag in All() map")
	}
}
