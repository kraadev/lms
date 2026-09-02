package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"lms/internal/database"
	"lms/internal/models"
	"lms/internal/utils"
)

type AccessController struct {
	db *database.DB
}

func NewAccessController(db *database.DB) *AccessController {
	return &AccessController{db: db}
}

// CheckClassMembership verifies if the authenticated user has access to read the class
func (ac *AccessController) CheckClassAccess(userID int64, role models.Role, classID int64) (bool, error) {
	if role == models.RoleAdmin {
		return true, nil
	}

	if role == models.RoleTeacher {
		var count int
		query := "SELECT COUNT(*) FROM classes WHERE id = $1 AND teacher_id = $2;"
		if ac.db.Driver != "postgres" {
			query = "SELECT COUNT(*) FROM classes WHERE id = ? AND teacher_id = ?;"
		}
		err := ac.db.QueryRow(query, classID, userID).Scan(&count)
		if err != nil {
			return false, err
		}
		if count > 0 {
			return true, nil
		}
		// Also check if teacher was added as member
	}

	// For student (or teacher added as member)
	var count int
	query := "SELECT COUNT(*) FROM class_members WHERE class_id = $1 AND user_id = $2;"
	if ac.db.Driver != "postgres" {
		query = "SELECT COUNT(*) FROM class_members WHERE class_id = ? AND user_id = ?;"
	}
	err := ac.db.QueryRow(query, classID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckClassManagement verifies if the authenticated user has rights to modify the class
func (ac *AccessController) CheckClassManagement(userID int64, role models.Role, classID int64) (bool, error) {
	if role == models.RoleAdmin {
		return true, nil
	}
	if role != models.RoleTeacher {
		return false, nil
	}

	var count int
	query := "SELECT COUNT(*) FROM classes WHERE id = $1 AND teacher_id = $2;"
	if ac.db.Driver != "postgres" {
		query = "SELECT COUNT(*) FROM classes WHERE id = ? AND teacher_id = ?;"
	}
	err := ac.db.QueryRow(query, classID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ClassAccessMiddleware enforces class reading access based on URL param {classId} or {id}
func (ac *AccessController) RequireClassAccess(paramName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUser(r)
			if user == nil {
				utils.Unauthorized(w, "Authentication required")
				return
			}

			idStr := chi.URLParam(r, paramName)
			classID, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil || classID <= 0 {
				utils.BadRequest(w, "Invalid class ID")
				return
			}

			hasAccess, err := ac.CheckClassAccess(user.ID, user.Role, classID)
			if err != nil {
				utils.InternalServerError(w, "Error validating class access")
				return
			}

			if !hasAccess {
				utils.Forbidden(w, "You are not enrolled in this class.")
				return
			}

			ctx := context.WithValue(r.Context(), "validated_class_id", classID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireClassManagement enforces teacher ownership or admin rights on {classId} or {id}
func (ac *AccessController) RequireClassManagement(paramName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUser(r)
			if user == nil {
				utils.Unauthorized(w, "Authentication required")
				return
			}

			idStr := chi.URLParam(r, paramName)
			classID, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil || classID <= 0 {
				utils.BadRequest(w, "Invalid class ID")
				return
			}

			canManage, err := ac.CheckClassManagement(user.ID, user.Role, classID)
			if err != nil {
				utils.InternalServerError(w, "Error validating class management permissions")
				return
			}

			if !canManage {
				utils.Forbidden(w, "Only the assigned teacher of this class or an administrator can perform this action.")
				return
			}

			ctx := context.WithValue(r.Context(), "validated_class_id", classID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
