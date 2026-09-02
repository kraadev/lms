package quizzes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"lms/internal/middleware"
	"lms/internal/models"
	"lms/internal/utils"
)

type Handler struct {
	service          *Service
	accessController *middleware.AccessController
}

func NewHandler(service *Service, accessController *middleware.AccessController) *Handler {
	return &Handler{
		service:          service,
		accessController: accessController,
	}
}

func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	classIDStr := r.URL.Query().Get("class_id")
	if classIDStr != "" {
		classID, _ := strconv.ParseInt(classIDStr, 10, 64)
		hasAccess, _ := h.accessController.CheckClassAccess(user.ID, user.Role, classID)
		if !hasAccess {
			utils.Forbidden(w, "Forbidden")
			return
		}
		quizzes, _ := h.service.ListByClass(classID)
		utils.JSON(w, http.StatusOK, quizzes)
		return
	}

	var query string
	if user.Role == "admin" {
		query = `SELECT id, class_id, teacher_id, title, description, duration_minutes, start_at, end_at, max_attempts, status, created_at, updated_at FROM quizzes ORDER BY start_at DESC;`
	} else if user.Role == "teacher" {
		query = `SELECT id, class_id, teacher_id, title, description, duration_minutes, start_at, end_at, max_attempts, status, created_at, updated_at FROM quizzes WHERE teacher_id = $1 ORDER BY start_at DESC;`
		if h.service.repo.db.Driver != "postgres" {
			query = `SELECT id, class_id, teacher_id, title, description, duration_minutes, start_at, end_at, max_attempts, status, created_at, updated_at FROM quizzes WHERE teacher_id = ? ORDER BY start_at DESC;`
		}
	} else {
		query = `
			SELECT q.id, q.class_id, q.teacher_id, q.title, q.description, q.duration_minutes, q.start_at, q.end_at, q.max_attempts, q.status, q.created_at, q.updated_at
			FROM quizzes q
			JOIN class_members cm ON q.class_id = cm.class_id
			WHERE cm.user_id = $1
			ORDER BY q.start_at DESC;
		`
		if h.service.repo.db.Driver != "postgres" {
			query = `
				SELECT q.id, q.class_id, q.teacher_id, q.title, q.description, q.duration_minutes, q.start_at, q.end_at, q.max_attempts, q.status, q.created_at, q.updated_at
				FROM quizzes q
				JOIN class_members cm ON q.class_id = cm.class_id
				WHERE cm.user_id = ?
				ORDER BY q.start_at DESC;
			`
		}
	}

	var rows *sql.Rows
	var err error
	if user.Role == "admin" {
		rows, err = h.service.repo.db.Query(query)
	} else {
		rows, err = h.service.repo.db.Query(query, user.ID)
	}

	var quizzes []models.Quiz
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var q models.Quiz
			var sRaw, eRaw, cRaw, uRaw interface{}
			if err := rows.Scan(&q.ID, &q.ClassID, &q.TeacherID, &q.Title, &q.Description, &q.DurationMinutes, &sRaw, &eRaw, &q.MaxAttempts, &q.Status, &cRaw, &uRaw); err == nil {
				q.StartAt = utils.ParseTime(sRaw)
				q.EndAt = utils.ParseTime(eRaw)
				q.CreatedAt = utils.ParseTime(cRaw)
				q.UpdatedAt = utils.ParseTime(uRaw)
				quizzes = append(quizzes, q)
			}
		}
	}

	utils.JSON(w, http.StatusOK, quizzes)
}

func (h *Handler) ListByClass(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	classID, _ := strconv.ParseInt(chi.URLParam(r, "classId"), 10, 64)

	hasAccess, err := h.accessController.CheckClassAccess(user.ID, user.Role, classID)
	if err != nil || !hasAccess {
		utils.Forbidden(w, "You do not have access to this class's quizzes.")
		return
	}

	quizzes, err := h.service.ListByClass(classID)
	if err != nil {
		utils.InternalServerError(w, "Failed to load quizzes")
		return
	}

	utils.JSON(w, http.StatusOK, quizzes)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	classID, _ := strconv.ParseInt(chi.URLParam(r, "classId"), 10, 64)

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, classID)
	if err != nil || !canManage {
		utils.Forbidden(w, "Only the assigned teacher or an admin can create quizzes.")
		return
	}

	var req CreateQuizRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid JSON payload")
		return
	}

	req.ClassID = classID
	req.TeacherID = user.ID

	quiz, err := h.service.CreateQuiz(req)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, quiz)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	teacherQuiz, err := h.service.GetByIDForTeacher(id)
	if err != nil {
		log.Printf("[DEBUG QUIZ] GetByIDForTeacher error for id %d: %v", id, err)
		utils.NotFound(w, "Quiz not found: "+err.Error())
		return
	}
	if teacherQuiz == nil {
		utils.NotFound(w, "Quiz not found")
		return
	}

	hasAccess, err := h.accessController.CheckClassAccess(user.ID, user.Role, teacherQuiz.ClassID)
	if err != nil || !hasAccess {
		utils.Forbidden(w, "Access denied to this quiz")
		return
	}

	if user.Role == models.RoleTeacher || user.Role == models.RoleAdmin {
		utils.JSON(w, http.StatusOK, teacherQuiz)
		return
	}

	// For student, hide correct answers
	studentQuiz, err := h.service.GetByIDForStudent(id)
	if err != nil {
		utils.InternalServerError(w, "Failed to load quiz")
		return
	}

	utils.JSON(w, http.StatusOK, studentQuiz)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	existing, err := h.service.GetByIDForTeacher(id)
	if err != nil || existing == nil {
		utils.NotFound(w, "Quiz not found")
		return
	}

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, existing.ClassID)
	if err != nil || !canManage {
		utils.Forbidden(w, "Only the assigned teacher or an admin can delete this quiz.")
		return
	}

	if err := h.service.DeleteQuiz(id); err != nil {
		utils.InternalServerError(w, "Failed to delete quiz")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Quiz deleted successfully"})
}

func (h *Handler) StartAttempt(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	quiz, err := h.service.GetByIDForTeacher(id)
	if err != nil || quiz == nil {
		utils.NotFound(w, "Quiz not found")
		return
	}

	hasAccess, err := h.accessController.CheckClassAccess(user.ID, user.Role, quiz.ClassID)
	if err != nil || !hasAccess {
		utils.Forbidden(w, "You are not enrolled in the class for this quiz.")
		return
	}

	attempt, err := h.service.StartAttempt(id, user.ID)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, attempt)
}

func (h *Handler) SubmitAttempt(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	attemptID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	var req SubmitAttemptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid JSON payload")
		return
	}

	result, err := h.service.SubmitAttempt(attemptID, user.ID, req)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, result)
}

func (h *Handler) GetAttempt(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	attemptID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	attempt, err := h.service.GetAttemptByID(attemptID)
	if err != nil || attempt == nil {
		utils.NotFound(w, "Attempt not found")
		return
	}

	quiz, err := h.service.GetByIDForTeacher(attempt.QuizID)
	if err != nil || quiz == nil {
		utils.NotFound(w, "Quiz not found")
		return
	}

	if user.Role == models.RoleStudent && attempt.StudentID != user.ID {
		utils.Forbidden(w, "You can only view your own quiz attempts")
		return
	}

	hasAccess, err := h.accessController.CheckClassAccess(user.ID, user.Role, quiz.ClassID)
	if err != nil || !hasAccess {
		utils.Forbidden(w, "Access denied")
		return
	}

	utils.JSON(w, http.StatusOK, attempt)
}

func (h *Handler) ListAttempts(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	quiz, err := h.service.GetByIDForTeacher(id)
	if err != nil || quiz == nil {
		utils.NotFound(w, "Quiz not found")
		return
	}

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, quiz.ClassID)
	if err != nil || !canManage {
		utils.Forbidden(w, "Only the assigned teacher or an admin can view all quiz attempts.")
		return
	}

	attempts, err := h.service.ListAttemptsByQuiz(id)
	if err != nil {
		utils.InternalServerError(w, "Failed to load attempts")
		return
	}

	utils.JSON(w, http.StatusOK, attempts)
}
