package dashboard

import (
	"database/sql"
	"net/http"
	"time"

	"lms/internal/database"
	"lms/internal/middleware"
	"lms/internal/models"
	"lms/internal/utils"
)

type Handler struct {
	db *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) StudentDashboard(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	// Current classes
	classQuery := `
		SELECT c.id, c.name, c.description, c.teacher_id, c.academic_year, c.status, c.created_at, c.updated_at
		FROM classes c
		JOIN class_members cm ON c.id = cm.class_id
		WHERE cm.user_id = $1 AND c.status = 'active';
	`
	if h.db.Driver != "postgres" {
		classQuery = `
			SELECT c.id, c.name, c.description, c.teacher_id, c.academic_year, c.status, c.created_at, c.updated_at
			FROM classes c
			JOIN class_members cm ON c.id = cm.class_id
			WHERE cm.user_id = ? AND c.status = 'active';
		`
	}

	rows, err := h.db.Query(classQuery, user.ID)
	var classes []models.Class
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var c models.Class
			var createdRaw, updatedRaw interface{}
			if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.TeacherID, &c.AcademicYear, &c.Status, &createdRaw, &updatedRaw); err == nil {
				c.CreatedAt = utils.ParseTime(createdRaw)
				c.UpdatedAt = utils.ParseTime(updatedRaw)
				classes = append(classes, c)
			}
		}
	}

	// Upcoming assignments
	assignQuery := `
		SELECT a.id, a.class_id, a.teacher_id, a.title, a.description, a.attachment_path, a.deadline, a.max_score, a.created_at, a.updated_at
		FROM assignments a
		JOIN class_members cm ON a.class_id = cm.class_id
		WHERE cm.user_id = $1 AND a.deadline > CURRENT_TIMESTAMP
		ORDER BY a.deadline ASC LIMIT 5;
	`
	if h.db.Driver != "postgres" {
		assignQuery = `
			SELECT a.id, a.class_id, a.teacher_id, a.title, a.description, a.attachment_path, a.deadline, a.max_score, a.created_at, a.updated_at
			FROM assignments a
			JOIN class_members cm ON a.class_id = cm.class_id
			WHERE cm.user_id = ? AND a.deadline > CURRENT_TIMESTAMP
			ORDER BY a.deadline ASC LIMIT 5;
		`
	}

	aRows, err := h.db.Query(assignQuery, user.ID)
	var assignments []models.Assignment
	if err == nil {
		defer aRows.Close()
		for aRows.Next() {
			var a models.Assignment
			var dRaw, cRaw, uRaw interface{}
			if err := aRows.Scan(&a.ID, &a.ClassID, &a.TeacherID, &a.Title, &a.Description, &a.AttachmentPath, &dRaw, &a.MaxScore, &cRaw, &uRaw); err == nil {
				a.Deadline = utils.ParseTime(dRaw)
				a.CreatedAt = utils.ParseTime(cRaw)
				a.UpdatedAt = utils.ParseTime(uRaw)
				assignments = append(assignments, a)
			}
		}
	}

	// Upcoming quizzes
	quizQuery := `
		SELECT q.id, q.class_id, q.teacher_id, q.title, q.description, q.duration_minutes, q.start_at, q.end_at, q.max_attempts, q.status, q.created_at, q.updated_at
		FROM quizzes q
		JOIN class_members cm ON q.class_id = cm.class_id
		WHERE cm.user_id = $1 AND q.end_at > CURRENT_TIMESTAMP AND q.status = 'published'
		ORDER BY q.start_at ASC LIMIT 5;
	`
	if h.db.Driver != "postgres" {
		quizQuery = `
			SELECT q.id, q.class_id, q.teacher_id, q.title, q.description, q.duration_minutes, q.start_at, q.end_at, q.max_attempts, q.status, q.created_at, q.updated_at
			FROM quizzes q
			JOIN class_members cm ON q.class_id = cm.class_id
			WHERE cm.user_id = ? AND q.end_at > CURRENT_TIMESTAMP AND q.status = 'published'
			ORDER BY q.start_at ASC LIMIT 5;
		`
	}

	qRows, err := h.db.Query(quizQuery, user.ID)
	var quizzes []models.Quiz
	if err == nil {
		defer qRows.Close()
		for qRows.Next() {
			var q models.Quiz
			var sRaw, eRaw, cRaw, uRaw interface{}
			if err := qRows.Scan(&q.ID, &q.ClassID, &q.TeacherID, &q.Title, &q.Description, &q.DurationMinutes, &sRaw, &eRaw, &q.MaxAttempts, &q.Status, &cRaw, &uRaw); err == nil {
				q.StartAt = utils.ParseTime(sRaw)
				q.EndAt = utils.ParseTime(eRaw)
				q.CreatedAt = utils.ParseTime(cRaw)
				q.UpdatedAt = utils.ParseTime(uRaw)
				quizzes = append(quizzes, q)
			}
		}
	}

	// Active meeting in student's classes
	meetQuery := `
		SELECT m.id, m.class_id, m.teacher_id, m.title, m.room_name, m.type, m.status, m.started_at, m.ended_at, m.created_at
		FROM meetings m
		JOIN class_members cm ON m.class_id = cm.class_id
		WHERE cm.user_id = $1 AND m.status = 'active'
		LIMIT 1;
	`
	if h.db.Driver != "postgres" {
		meetQuery = `
			SELECT m.id, m.class_id, m.teacher_id, m.title, m.room_name, m.type, m.status, m.started_at, m.ended_at, m.created_at
			FROM meetings m
			JOIN class_members cm ON m.class_id = cm.class_id
			WHERE cm.user_id = ? AND m.status = 'active'
			LIMIT 1;
		`
	}

	var activeMeeting *models.Meeting
	mRow := h.db.QueryRow(meetQuery, user.ID)
	var m models.Meeting
	var sMeetRaw, eMeetRaw, cMeetRaw interface{}
	if err := mRow.Scan(&m.ID, &m.ClassID, &m.TeacherID, &m.Title, &m.RoomName, &m.Type, &m.Status, &sMeetRaw, &eMeetRaw, &cMeetRaw); err == nil {
		if sMeetRaw != nil {
			tVal := utils.ParseTime(sMeetRaw)
			m.StartedAt = &tVal
		}
		if eMeetRaw != nil {
			tVal := utils.ParseTime(eMeetRaw)
			m.EndedAt = &tVal
		}
		m.CreatedAt = utils.ParseTime(cMeetRaw)
		activeMeeting = &m
	}

	greeting := "Selamat datang kembali!"
	hour := time.Now().Hour()
	if hour < 12 {
		greeting = "Selamat pagi!"
	} else if hour < 17 {
		greeting = "Selamat siang!"
	} else {
		greeting = "Selamat malam!"
	}

	resp := map[string]interface{}{
		"greeting":              greeting,
		"current_classes":       classes,
		"upcoming_assignments":  assignments,
		"upcoming_quizzes":      quizzes,
		"active_meeting":        activeMeeting,
		"recent_announcements":  []interface{}{},
		"recent_grades":         []interface{}{},
		"unread_chat_count":     0,
	}

	utils.JSON(w, http.StatusOK, resp)
}

func (h *Handler) TeacherDashboard(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	// Classes taught
	classQuery := `
		SELECT id, name, description, teacher_id, academic_year, status, created_at, updated_at
		FROM classes
		WHERE teacher_id = $1 AND status = 'active';
	`
	if h.db.Driver != "postgres" {
		classQuery = `
			SELECT id, name, description, teacher_id, academic_year, status, created_at, updated_at
			FROM classes
			WHERE teacher_id = ? AND status = 'active';
		`
	}

	rows, err := h.db.Query(classQuery, user.ID)
	var classes []models.Class
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var c models.Class
			var createdRaw, updatedRaw interface{}
			if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.TeacherID, &c.AcademicYear, &c.Status, &createdRaw, &updatedRaw); err == nil {
				c.CreatedAt = utils.ParseTime(createdRaw)
				c.UpdatedAt = utils.ParseTime(updatedRaw)
				classes = append(classes, c)
			}
		}
	}

	// Active meetings
	meetQuery := `
		SELECT id, class_id, teacher_id, title, room_name, type, status, started_at, ended_at, created_at
		FROM meetings
		WHERE teacher_id = $1 AND status = 'active';
	`
	if h.db.Driver != "postgres" {
		meetQuery = `
			SELECT id, class_id, teacher_id, title, room_name, type, status, started_at, ended_at, created_at
			FROM meetings
			WHERE teacher_id = ? AND status = 'active';
		`
	}

	mRows, err := h.db.Query(meetQuery, user.ID)
	var activeMeetings []models.Meeting
	if err == nil {
		defer mRows.Close()
		for mRows.Next() {
			var m models.Meeting
			var sMeetRaw, eMeetRaw, cMeetRaw interface{}
			if err := mRows.Scan(&m.ID, &m.ClassID, &m.TeacherID, &m.Title, &m.RoomName, &m.Type, &m.Status, &sMeetRaw, &eMeetRaw, &cMeetRaw); err == nil {
				if sMeetRaw != nil {
					tVal := utils.ParseTime(sMeetRaw)
					m.StartedAt = &tVal
				}
				if eMeetRaw != nil {
					tVal := utils.ParseTime(eMeetRaw)
					m.EndedAt = &tVal
				}
				m.CreatedAt = utils.ParseTime(cMeetRaw)
				activeMeetings = append(activeMeetings, m)
			}
		}
	}

	resp := map[string]interface{}{
		"classes_taught":        classes,
		"pending_grading_count": 0,
		"pending_grading":       []interface{}{},
		"upcoming_assignments":  []interface{}{},
		"recent_submissions":    []interface{}{},
		"active_meetings":       activeMeetings,
		"recent_chats":          []interface{}{},
		"quiz_overview": map[string]interface{}{
			"total_quizzes":  1,
			"active_quizzes": 1,
			"total_attempts": 0,
		},
	}

	utils.JSON(w, http.StatusOK, resp)
}

func (h *Handler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	var totalStudents, totalTeachers, totalClasses, activeClasses int

	_ = h.db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'student'").Scan(&totalStudents)
	_ = h.db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'teacher'").Scan(&totalTeachers)
	_ = h.db.QueryRow("SELECT COUNT(*) FROM classes").Scan(&totalClasses)
	_ = h.db.QueryRow("SELECT COUNT(*) FROM classes WHERE status = 'active'").Scan(&activeClasses)

	userQuery := "SELECT id, name, email, role, avatar_url, created_at, updated_at FROM users ORDER BY id DESC LIMIT 5"
	rows, err := h.db.Query(userQuery)
	var recentUsers []models.User
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var u models.User
			var cRaw, uRaw interface{}
			var avatar sql.NullString
			if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &avatar, &cRaw, &uRaw); err == nil {
				if avatar.Valid {
					u.AvatarURL = &avatar.String
				}
				u.CreatedAt = utils.ParseTime(cRaw)
				u.UpdatedAt = utils.ParseTime(uRaw)
				recentUsers = append(recentUsers, u)
			}
		}
	}

	resp := map[string]interface{}{
		"total_students":   totalStudents,
		"total_teachers":   totalTeachers,
		"total_classes":    totalClasses,
		"active_classes":   activeClasses,
		"system_activity":  []interface{}{},
		"recent_users":     recentUsers,
	}

	utils.JSON(w, http.StatusOK, resp)
}
