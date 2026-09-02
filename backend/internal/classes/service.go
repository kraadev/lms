package classes

import (
	"errors"
	"strings"

	"lms/internal/models"
	"lms/internal/users"
)

type Service struct {
	repo     *Repository
	userRepo *users.Repository
}

func NewService(repo *Repository, userRepo *users.Repository) *Service {
	return &Service{
		repo:     repo,
		userRepo: userRepo,
	}
}

type CreateClassRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	TeacherID    int64  `json:"teacher_id"`
	AcademicYear string `json:"academic_year"`
}

func (s *Service) CreateClass(req CreateClassRequest, currentUserID int64, role models.Role) (*models.Class, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, errors.New("class name is required")
	}

	teacherID := req.TeacherID
	if role == models.RoleTeacher {
		teacherID = currentUserID
	} else if teacherID <= 0 {
		teacherID = currentUserID
	}

	teacher, err := s.userRepo.FindByID(teacherID)
	if err != nil || teacher == nil {
		return nil, errors.New("assigned teacher not found")
	}
	if teacher.Role != models.RoleTeacher && teacher.Role != models.RoleAdmin {
		return nil, errors.New("assigned user must have teacher role")
	}

	class := &models.Class{
		Name:         req.Name,
		Description:  req.Description,
		TeacherID:    teacherID,
		AcademicYear: req.AcademicYear,
		Status:       "active",
	}

	if err := s.repo.Create(class); err != nil {
		return nil, err
	}
	class.TeacherName = teacher.Name
	return class, nil
}

func (s *Service) ListUserClasses(userID int64, role models.Role) ([]models.Class, error) {
	return s.repo.ListForUser(userID, role)
}

func (s *Service) GetClassByID(id int64) (*models.Class, error) {
	return s.repo.FindByID(id)
}

func (s *Service) UpdateClass(class *models.Class) error {
	return s.repo.Update(class)
}

func (s *Service) DeleteClass(id int64) error {
	return s.repo.Delete(id)
}

func (s *Service) AddStudentToClass(classID, studentID int64) error {
	student, err := s.userRepo.FindByID(studentID)
	if err != nil || student == nil {
		return errors.New("student not found")
	}
	return s.repo.AddMember(classID, studentID)
}

func (s *Service) RemoveStudentFromClass(classID, studentID int64) error {
	return s.repo.RemoveMember(classID, studentID)
}

func (s *Service) ListClassMembers(classID int64) ([]models.ClassMember, error) {
	return s.repo.ListMembers(classID)
}
