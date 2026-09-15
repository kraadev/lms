package middleware

import (
	"net/http"

	"lms/internal/models"
	"lms/internal/utils"
)

type Resource string
type Action string

const (
	ResourceUsers       Resource = "users"
	ResourceClasses     Resource = "classes"
	ResourceMaterials   Resource = "materials"
	ResourceAssignments Resource = "assignments"
	ResourceQuizzes     Resource = "quizzes"
	ResourceSettings    Resource = "settings"

	ActionRead   Action = "read"
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
	ActionGrade  Action = "grade"
)

// Can evaluates whether a given role has permission to perform an action on a resource
func Can(role models.Role, res Resource, act Action) bool {
	if role == models.RoleAdmin {
		return true // Admins have superuser rights
	}

	switch role {
	case models.RoleTeacher:
		if res == ResourceUsers && (act == ActionCreate || act == ActionDelete) {
			return false
		}
		if res == ResourceSettings && (act == ActionUpdate || act == ActionDelete) {
			return false
		}
		return true

	case models.RoleStudent:
		if act == ActionRead {
			return true
		}
		if res == ResourceAssignments && act == ActionCreate {
			return true // Can submit assignment
		}
		return false
	}

	return false
}

// RequirePermission verifies the authenticated user possesses the specific permission
func RequirePermission(res Resource, act Action) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUser(r)
			if user == nil {
				utils.Unauthorized(w, "")
				return
			}

			if !Can(user.Role, res, act) {
				utils.Forbidden(w, "Insufficient role permissions for this operation.")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
