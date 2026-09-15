package middleware

import (
	"testing"

	"lms/internal/models"
)

func TestRBACPermissions(t *testing.T) {
	// Admin can do anything
	if !Can(models.RoleAdmin, ResourceUsers, ActionDelete) {
		t.Errorf("expected admin to have delete user permission")
	}

	// Teacher can grade assignments
	if !Can(models.RoleTeacher, ResourceAssignments, ActionGrade) {
		t.Errorf("expected teacher to have grading permission")
	}

	// Teacher cannot delete users
	if Can(models.RoleTeacher, ResourceUsers, ActionDelete) {
		t.Errorf("expected teacher NOT to have delete user permission")
	}

	// Student can read classes
	if !Can(models.RoleStudent, ResourceClasses, ActionRead) {
		t.Errorf("expected student to have read class permission")
	}

	// Student cannot grade assignments
	if Can(models.RoleStudent, ResourceAssignments, ActionGrade) {
		t.Errorf("expected student NOT to have grade assignment permission")
	}
}
