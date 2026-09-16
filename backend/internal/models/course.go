package models

import "time"

type Course struct {
	ID          int64     `json:"id"`
	ClassID     int64     `json:"class_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Modules     []Module  `json:"modules,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Module struct {
	ID        int64     `json:"id"`
	CourseID  int64     `json:"course_id"`
	Title     string    `json:"title"`
	Order     int       `json:"order"`
	Lessons   []Lesson  `json:"lessons,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Lesson struct {
	ID          int64     `json:"id"`
	ModuleID    int64     `json:"module_id"`
	Title       string    `json:"title"`
	ContentType string    `json:"content_type"` // "markdown", "pdf", "video"
	ContentURL  string    `json:"content_url,omitempty"`
	Content     string    `json:"content,omitempty"`
	Order       int       `json:"order"`
	DurationSec int       `json:"duration_sec"`
	CreatedAt   time.Time `json:"created_at"`
}
