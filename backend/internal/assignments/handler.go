package assignments

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"lms/internal/middleware"
	"lms/internal/storage"
	"lms/internal/utils"
)

type Handler struct {
	service          *Service
	accessController *middleware.AccessController
	storageService   *storage.StorageService
}

func NewHandler(service *Service, accessController *middleware.AccessController, storageService *storage.StorageService) *Handler {
	return &Handler{
		service:          service,
		accessController: accessController,
		storageService:   storageService,
	}
}

func (h *Handler) ListByClass(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	classID, _ := strconv.ParseInt(chi.URLParam(r, "classId"), 10, 64)

	hasAccess, err := h.accessController.CheckClassAccess(user.ID, user.Role, classID)
	if err != nil || !hasAccess {
		utils.Forbidden(w, "You do not have access to this class's assignments.")
		return
	}

	assignments, err := h.service.ListByClass(classID)
	if err != nil {
		utils.InternalServerError(w, "Failed to load assignments")
		return
	}

	utils.JSON(w, http.StatusOK, assignments)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	classID, _ := strconv.ParseInt(chi.URLParam(r, "classId"), 10, 64)

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, classID)
	if err != nil || !canManage {
		utils.Forbidden(w, "Only the assigned teacher or an admin can create assignments.")
		return
	}

	var req CreateAssignmentRequest
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			utils.BadRequest(w, "Failed to parse multipart form")
			return
		}

		req.Title = r.FormValue("title")
		req.Description = r.FormValue("description")
		if maxScore, err := strconv.ParseFloat(r.FormValue("max_score"), 64); err == nil {
			req.MaxScore = maxScore
		}
		if deadline, err := time.Parse(time.RFC3339, r.FormValue("deadline")); err == nil {
			req.Deadline = deadline
		}

		file, header, err := r.FormFile("file")
		if err == nil {
			_ = file.Close()
			filePath, err := h.storageService.SaveUploadedFile(header, "assignments")
			if err != nil {
				utils.BadRequest(w, "Upload error: "+err.Error())
				return
			}
			req.AttachmentPath = &filePath
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.BadRequest(w, "Invalid JSON payload")
			return
		}
	}

	req.ClassID = classID
	req.TeacherID = user.ID

	a, err := h.service.CreateAssignment(req)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, a)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	a, err := h.service.GetByID(id)
	if err != nil || a == nil {
		utils.NotFound(w, "Assignment not found")
		return
	}

	hasAccess, err := h.accessController.CheckClassAccess(user.ID, user.Role, a.ClassID)
	if err != nil || !hasAccess {
		utils.Forbidden(w, "Access denied to this assignment")
		return
	}

	utils.JSON(w, http.StatusOK, a)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	existing, err := h.service.GetByID(id)
	if err != nil || existing == nil {
		utils.NotFound(w, "Assignment not found")
		return
	}

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, existing.ClassID)
	if err != nil || !canManage {
		utils.Forbidden(w, "Only the assigned teacher or an admin can edit this assignment.")
		return
	}

	var req struct {
		Title       *string    `json:"title"`
		Description *string    `json:"description"`
		Deadline    *time.Time `json:"deadline"`
		MaxScore    *float64   `json:"max_score"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid JSON payload")
		return
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Deadline != nil {
		existing.Deadline = *req.Deadline
	}
	if req.MaxScore != nil {
		existing.MaxScore = *req.MaxScore
	}

	if err := h.service.Update(existing); err != nil {
		utils.InternalServerError(w, "Failed to update assignment")
		return
	}

	utils.JSON(w, http.StatusOK, existing)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	existing, err := h.service.GetByID(id)
	if err != nil || existing == nil {
		utils.NotFound(w, "Assignment not found")
		return
	}

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, existing.ClassID)
	if err != nil || !canManage {
		utils.Forbidden(w, "Only the assigned teacher or an admin can delete this assignment.")
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.InternalServerError(w, "Failed to delete assignment")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Assignment deleted successfully"})
}

func (h *Handler) Submit(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	assignment, err := h.service.GetByID(id)
	if err != nil || assignment == nil {
		utils.NotFound(w, "Assignment not found")
		return
	}

	hasAccess, err := h.accessController.CheckClassAccess(user.ID, user.Role, assignment.ClassID)
	if err != nil || !hasAccess {
		utils.Forbidden(w, "You are not enrolled in the class for this assignment.")
		return
	}

	var textAnswer *string
	var filePath *string
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			utils.BadRequest(w, "Failed to parse form")
			return
		}
		if txt := r.FormValue("text_answer"); txt != "" {
			textAnswer = &txt
		}
		file, header, err := r.FormFile("file")
		if err == nil {
			_ = file.Close()
			saved, err := h.storageService.SaveUploadedFile(header, "submissions")
			if err != nil {
				utils.BadRequest(w, "Upload error: "+err.Error())
				return
			}
			filePath = &saved
		}
	} else {
		var req struct {
			TextAnswer *string `json:"text_answer"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			textAnswer = req.TextAnswer
		}
	}

	if textAnswer == nil && filePath == nil {
		utils.BadRequest(w, "Either text answer or file upload is required for submission")
		return
	}

	sub, err := h.service.SubmitAssignment(id, user.ID, textAnswer, filePath)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, sub)
}

func (h *Handler) GetMySubmission(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	sub, err := h.service.GetStudentSubmission(id, user.ID)
	if err != nil {
		utils.InternalServerError(w, "Failed to load submission")
		return
	}
	if sub == nil {
		utils.NotFound(w, "No submission found for this assignment")
		return
	}

	utils.JSON(w, http.StatusOK, sub)
}

func (h *Handler) ListSubmissions(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	assignment, err := h.service.GetByID(id)
	if err != nil || assignment == nil {
		utils.NotFound(w, "Assignment not found")
		return
	}

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, assignment.ClassID)
	if err != nil || !canManage {
		utils.Forbidden(w, "Only the assigned teacher or an admin can view all submissions.")
		return
	}

	subs, err := h.service.ListSubmissions(id)
	if err != nil {
		utils.InternalServerError(w, "Failed to load submissions")
		return
	}

	utils.JSON(w, http.StatusOK, subs)
}

func (h *Handler) Grade(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	submissionID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	sub, err := h.service.repo.GetSubmissionByID(submissionID)
	if err != nil || sub == nil {
		utils.NotFound(w, "Submission not found")
		return
	}

	assignment, err := h.service.GetByID(sub.AssignmentID)
	if err != nil || assignment == nil {
		utils.NotFound(w, "Assignment not found")
		return
	}

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, assignment.ClassID)
	if err != nil || !canManage {
		utils.Forbidden(w, "Only the assigned teacher or an admin can grade this submission.")
		return
	}

	var req struct {
		Score    float64 `json:"score"`
		Feedback string  `json:"feedback"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid JSON payload")
		return
	}

	graded, err := h.service.GradeSubmission(submissionID, req.Score, req.Feedback)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, graded)
}
