package materials

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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
		utils.Forbidden(w, "You do not have access to this class's materials.")
		return
	}

	materials, err := h.service.ListByClass(classID)
	if err != nil {
		utils.InternalServerError(w, "Failed to load materials")
		return
	}

	utils.JSON(w, http.StatusOK, materials)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	classID, _ := strconv.ParseInt(chi.URLParam(r, "classId"), 10, 64)

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, classID)
	if err != nil || !canManage {
		utils.Forbidden(w, "Only the assigned teacher or an admin can add materials.")
		return
	}

	var req CreateMaterialRequest
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Parse multipart form (max 32MB)
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			utils.BadRequest(w, "Failed to parse multipart form")
			return
		}

		req.Title = r.FormValue("title")
		req.Description = r.FormValue("description")
		req.Content = r.FormValue("content")
		if extURL := r.FormValue("external_url"); extURL != "" {
			req.ExternalURL = &extURL
		}

		// Handle file upload
		file, header, err := r.FormFile("file")
		if err == nil {
			_ = file.Close()
			filePath, err := h.storageService.SaveUploadedFile(header, "materials")
			if err != nil {
				utils.BadRequest(w, "Upload error: "+err.Error())
				return
			}
			req.FilePath = &filePath
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.BadRequest(w, "Invalid JSON payload")
			return
		}
	}

	req.ClassID = classID
	req.TeacherID = user.ID

	mat, err := h.service.CreateMaterial(req)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, mat)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	mat, err := h.service.GetByID(id)
	if err != nil || mat == nil {
		utils.NotFound(w, "Material not found")
		return
	}

	hasAccess, err := h.accessController.CheckClassAccess(user.ID, user.Role, mat.ClassID)
	if err != nil || !hasAccess {
		utils.Forbidden(w, "Access denied to this material")
		return
	}

	utils.JSON(w, http.StatusOK, mat)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	existing, err := h.service.GetByID(id)
	if err != nil || existing == nil {
		utils.NotFound(w, "Material not found")
		return
	}

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, existing.ClassID)
	if err != nil || !canManage {
		utils.Forbidden(w, "Only the assigned teacher or an admin can edit this material.")
		return
	}

	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Content     *string `json:"content"`
		ExternalURL *string `json:"external_url"`
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
	if req.Content != nil {
		existing.Content = *req.Content
	}
	if req.ExternalURL != nil {
		existing.ExternalURL = req.ExternalURL
	}

	if err := h.service.Update(existing); err != nil {
		utils.InternalServerError(w, "Failed to update material")
		return
	}

	utils.JSON(w, http.StatusOK, existing)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	existing, err := h.service.GetByID(id)
	if err != nil || existing == nil {
		utils.NotFound(w, "Material not found")
		return
	}

	canManage, err := h.accessController.CheckClassManagement(user.ID, user.Role, existing.ClassID)
	if err != nil || !canManage {
		utils.Forbidden(w, "Only the assigned teacher or an admin can delete this material.")
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.InternalServerError(w, "Failed to delete material")
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"message": "Material deleted successfully"})
}
