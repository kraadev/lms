package models

import (
	"time"
)

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleTeacher Role = "teacher"
	RoleStudent Role = "student"
)

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	AvatarURL    *string   `json:"avatar_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Class struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	TeacherID    int64     `json:"teacher_id"`
	TeacherName  string    `json:"teacher_name,omitempty"`
	AcademicYear string    `json:"academic_year"`
	Status       string    `json:"status"` // "active", "archived"
	MemberCount  int       `json:"member_count,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ClassMember struct {
	ID       int64     `json:"id"`
	ClassID  int64     `json:"class_id"`
	UserID   int64     `json:"user_id"`
	User     *User     `json:"user,omitempty"`
	JoinedAt time.Time `json:"joined_at"`
}

type Material struct {
	ID          int64     `json:"id"`
	ClassID     int64     `json:"class_id"`
	TeacherID   int64     `json:"teacher_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	FilePath    *string   `json:"file_path,omitempty"`
	ExternalURL *string   `json:"external_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Assignment struct {
	ID              int64     `json:"id"`
	ClassID         int64     `json:"class_id"`
	TeacherID       int64     `json:"teacher_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	AttachmentPath  *string   `json:"attachment_path,omitempty"`
	Deadline        time.Time `json:"deadline"`
	MaxScore        float64   `json:"max_score"`
	SubmissionCount int       `json:"submission_count,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type AssignmentSubmission struct {
	ID           int64      `json:"id"`
	AssignmentID int64      `json:"assignment_id"`
	StudentID    int64      `json:"student_id"`
	StudentName  string     `json:"student_name,omitempty"`
	StudentEmail string     `json:"student_email,omitempty"`
	TextAnswer   *string    `json:"text_answer,omitempty"`
	FilePath     *string    `json:"file_path,omitempty"`
	SubmittedAt  time.Time  `json:"submitted_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Score        *float64   `json:"score,omitempty"`
	Feedback     *string    `json:"feedback,omitempty"`
	Status       string     `json:"status"` // "submitted", "late", "graded"
}

type Quiz struct {
	ID              int64          `json:"id"`
	ClassID         int64          `json:"class_id"`
	TeacherID       int64          `json:"teacher_id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	StartAt         time.Time      `json:"start_at"`
	EndAt           time.Time      `json:"end_at"`
	MaxAttempts     int            `json:"max_attempts"`
	Status          string         `json:"status"` // "draft", "published", "closed"
	Questions       []QuizQuestion `json:"questions,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type QuestionType string

const (
	QuestionMultipleChoice QuestionType = "multiple_choice"
	QuestionTrueFalse      QuestionType = "true_false"
	QuestionShortAnswer    QuestionType = "short_answer"
)

type QuizQuestion struct {
	ID         int64        `json:"id"`
	QuizID     int64        `json:"quiz_id"`
	Question   string       `json:"question"`
	Type       QuestionType `json:"type"`
	Points     float64      `json:"points"`
	OrderIndex int          `json:"order_index"`
	Options    []QuizOption `json:"options,omitempty"`
}

type QuizOption struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"question_id"`
	OptionText string `json:"option_text"`
	IsCorrect  *bool  `json:"is_correct,omitempty"` // Omitted for students prior to grading/completion
}

type QuizAttempt struct {
	ID          int64        `json:"id"`
	QuizID      int64        `json:"quiz_id"`
	StudentID   int64        `json:"student_id"`
	StudentName string       `json:"student_name,omitempty"`
	StartedAt   time.Time    `json:"started_at"`
	SubmittedAt *time.Time   `json:"submitted_at,omitempty"`
	Score       *float64     `json:"score,omitempty"`
	Status      string       `json:"status"` // "in_progress", "submitted", "graded", "expired"
	Answers     []QuizAnswer `json:"answers,omitempty"`
}

type QuizAnswer struct {
	ID               int64    `json:"id"`
	AttemptID        int64    `json:"attempt_id"`
	QuestionID       int64    `json:"question_id"`
	SelectedOptionID *int64   `json:"selected_option_id,omitempty"`
	TextAnswer       *string  `json:"text_answer,omitempty"`
	IsCorrect        *bool    `json:"is_correct,omitempty"`
	EarnedPoints     *float64 `json:"earned_points,omitempty"`
}

type MeetingType string

const (
	MeetingTypeVideo MeetingType = "video"
	MeetingTypeAudio MeetingType = "audio"
)

type MeetingStatus string

const (
	MeetingStatusScheduled MeetingStatus = "scheduled"
	MeetingStatusActive    MeetingStatus = "active"
	MeetingStatusEnded     MeetingStatus = "ended"
)

type Meeting struct {
	ID          int64         `json:"id"`
	ClassID     int64         `json:"class_id"`
	TeacherID   int64         `json:"teacher_id"`
	TeacherName string        `json:"teacher_name,omitempty"`
	Title       string        `json:"title"`
	RoomName    string        `json:"room_name"`
	Type        MeetingType   `json:"type"`
	Status      MeetingStatus `json:"status"`
	StartedAt   *time.Time    `json:"started_at,omitempty"`
	EndedAt     *time.Time    `json:"ended_at,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
}

type Message struct {
	ID         int64     `json:"id"`
	ClassID    int64     `json:"class_id"`
	UserID     int64     `json:"user_id"`
	UserName   string    `json:"user_name"`
	UserRole   string    `json:"user_role"`
	UserAvatar *string   `json:"user_avatar,omitempty"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}

type Notification struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Type      string    `json:"type"` // "assignment", "grade", "quiz", "meeting", "system"
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
