package services

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/bomboskuy/UAS-Backend/app/models"
	"github.com/bomboskuy/UAS-Backend/app/repositories"
	"github.com/bomboskuy/UAS-Backend/helper"
	"github.com/bomboskuy/UAS-Backend/utils"
)

type UserService struct {
	userRepo     repositories.UserRepository
	roleRepo     repositories.RoleRepository
	studentRepo  repositories.StudentRepository
	lecturerRepo repositories.LecturerRepository
}

func NewUserService(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	studentRepo repositories.StudentRepository,
	lecturerRepo repositories.LecturerRepository,
) *UserService {
	return &UserService{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		studentRepo:  studentRepo,
		lecturerRepo: lecturerRepo,
	}
}

//
// ==========================
// USERS
// ==========================
//
func (s *UserService) GetAll(c *fiber.Ctx) error {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return helper.InternalServerError(c, "Gagal mengambil data user")
	}

	var res []models.UserResponse
	for _, u := range users {
		role, _ := s.roleRepo.FindByID(u.RoleID)

		res = append(res, models.UserResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			FullName:  u.FullName,
			Role:      role.Name,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt,
		})
	}

	return helper.Success(c, "OK", res)
}

func (s *UserService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return helper.NotFound(c, "User tidak ditemukan")
	}

	role, _ := s.roleRepo.FindByID(user.RoleID)

	return helper.Success(c, "OK", models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      role.Name,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	})
}

func (s *UserService) Create(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest(c, "Invalid request body", nil)
	}

	role, err := s.roleRepo.FindByName(req.RoleName)
	if err != nil {
		return helper.BadRequest(c, "Role tidak valid", nil)
	}

	// =========================
	// VALIDASI MAHASISWA
	// =========================
	if role.Name == "Mahasiswa" {
		if req.AdvisorID == "" {
			return helper.BadRequest(c, "advisor_id wajib diisi", nil)
		}

		// cek dosen wali ada atau tidak
		if _, err := s.lecturerRepo.FindByID(req.AdvisorID); err != nil {
			return helper.BadRequest(c, "advisor_id tidak valid", nil)
		}
	}

	hashedPassword, _ := utils.HashPassword(req.Password)

	userID := uuid.NewString()
	user := &models.User{
		ID:           userID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		RoleID:       role.ID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// =========================
	// CREATE USER
	// =========================
	if err := s.userRepo.Create(user); err != nil {
		return helper.InternalServerError(c, "Gagal membuat user")
	}

	// =========================
	// CREATE STUDENT
	// =========================
	if role.Name == "Mahasiswa" {
		student := &models.Student{
			ID:           uuid.NewString(),
			UserID:       userID,
			StudentID:    req.StudentID,
			ProgramStudy: req.ProgramStudy,
			AcademicYear: req.AcademicYear,
			AdvisorID:    req.AdvisorID,
			CreatedAt:    time.Now(),
		}

		if err := s.studentRepo.Create(student); err != nil {
			return helper.InternalServerError(c, "Gagal membuat data mahasiswa")
		}
	}

	// =========================
	// CREATE LECTURER
	// =========================
	if role.Name == "Dosen Wali" {
		lecturer := &models.Lecturer{
			ID:         uuid.NewString(),
			UserID:     userID,
			LecturerID: req.LecturerID,
			Department: req.Department,
			CreatedAt:  time.Now(),
		}

		if err := s.lecturerRepo.Create(lecturer); err != nil {
			return helper.InternalServerError(c, "Gagal membuat data dosen")
		}
	}

	return helper.Success(c, "User berhasil dibuat", user)
}

//
// ==========================
// UPDATE & DELETE USER
// ==========================
func (s *UserService) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest(c, "Invalid request body", nil)
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return helper.NotFound(c, "User tidak ditemukan")
	}

	user.Username = req.Username
	user.Email = req.Email
	user.FullName = req.FullName
	user.IsActive = req.IsActive
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return helper.InternalServerError(c, "Gagal update user")
	}

	return helper.Success(c, "User berhasil diupdate", user)
}

func (s *UserService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := s.userRepo.Delete(id); err != nil {
		return helper.InternalServerError(c, "Gagal menghapus user")
	}

	return helper.Success(c, "User berhasil dihapus", nil)
}

//
// ==========================
// STUDENTS
// ==========================
func (s *UserService) GetStudents(c *fiber.Ctx) error {
	students, err := s.studentRepo.FindAll()
	if err != nil {
		return helper.InternalServerError(c, "Gagal mengambil data mahasiswa")
	}

	return helper.Success(c, "OK", students)
}

func (s *UserService) GetStudentByID(c *fiber.Ctx) error {
	id := c.Params("id")

	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		return helper.NotFound(c, "Mahasiswa tidak ditemukan")
	}

	return helper.Success(c, "OK", student)
}

func (s *UserService) AssignAdvisor(c *fiber.Ctx) error {
	studentID := c.Params("id")

	var req struct {
		AdvisorID string `json:"advisor_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest(c, "Invalid request body", nil)
	}

	// validasi dosen wali
	if _, err := s.lecturerRepo.FindByID(req.AdvisorID); err != nil {
		return helper.BadRequest(c, "advisor_id tidak valid", nil)
	}

	if err := s.studentRepo.AssignAdvisor(studentID, req.AdvisorID); err != nil {
		return helper.InternalServerError(c, "Gagal assign dosen wali")
	}

	return helper.Success(c, "Dosen wali berhasil ditetapkan", nil)
}

//
// ==========================
// LECTURERS
// ==========================
func (s *UserService) GetLecturers(c *fiber.Ctx) error {
	lecturers, err := s.lecturerRepo.FindAll()
	if err != nil {
		return helper.InternalServerError(c, "Gagal mengambil data dosen")
	}

	return helper.Success(c, "OK", lecturers)
}

func (s *UserService) GetAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")

	students, err := s.lecturerRepo.FindAdvisees(lecturerID)
	if err != nil {
		return helper.InternalServerError(c, "Gagal mengambil mahasiswa bimbingan")
	}

	return helper.Success(c, "OK", students)
}

