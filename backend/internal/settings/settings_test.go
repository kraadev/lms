package settings

import "testing"

func TestSettingsStore(t *testing.T) {
	store := NewStore()

	if store.GetString("app_name") != "LMS Portal" {
		t.Errorf("expected default app_name LMS Portal")
	}

	if store.GetInt("max_upload_mb", 10) != 50 {
		t.Errorf("expected 50MB max upload")
	}

	if !store.GetBool("allow_register") {
		t.Errorf("expected allow_register true")
	}

	store.Set("allow_register", "false")
	if store.GetBool("allow_register") {
		t.Errorf("expected allow_register false after Set")
	}
}
