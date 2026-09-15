package models

import "time"

type AuditLog struct {
	ID           int64     `json:"id"`
	ActorID      *int64    `json:"actor_id,omitempty"`
	ActorEmail   string    `json:"actor_email,omitempty"`
	Action       string    `json:"action"`
	Resource     string    `json:"resource"`
	ResourceID   string    `json:"resource_id,omitempty"`
	IPAddress    string    `json:"ip_address,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
	StatusCode   int       `json:"status_code"`
	DurationMs   int64     `json:"duration_ms"`
	CreatedAt    time.Time `json:"created_at"`
}
