package services

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/bomboskuy/UAS-Backend/app/models"
	"github.com/bomboskuy/UAS-Backend/app/repositories"
	"github.com/bomboskuy/UAS-Backend/helper"
)

type AchievementService struct {
	achievementRepo          repositories.AchievementRepository
	achievementReferenceRepo repositories.AchievementReferenceRepository
	studentRepo              repositories.StudentRepository
	lecturerRepo             repositories.LecturerRepository
}

func NewAchievementService(
	a repositories.AchievementRepository,
	r repositories.AchievementReferenceRepository,
	s repositories.StudentRepository,
	l repositories.LecturerRepository,
) *AchievementService {
	return &AchievementService{a, r, s, l}
}


func (s *AchievementService) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req models.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	student, err := s.studentRepo.FindByUserID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Student not found",
		})
	}

	achievement := &models.Achievement{
		StudentID:       student.ID,
		AchievementType: req.AchievementType,
		Title:           req.Title,
		Description:     req.Description,
		Details:         parseDetailsToAchievementDetail(req.Details),
		Tags:            req.Tags,
		Points:          req.Points,
	}

	mongoID, err := s.achievementRepo.Create(achievement)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to save achievement",
		})
	}

	ref := &models.AchievementReference{
		ID:                 uuid.New().String(),
		UserID:             userID,
		StudentID:          student.ID,
		MongoAchievementID: mongoID,
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.achievementReferenceRepo.Create(ref); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to save reference",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Achievement created",
		"data":    ref,
	})
}


func (s *AchievementService) GetAll(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	roleID := c.Locals("role_id").(string)

	var refs []models.AchievementReference
	var err error

	switch roleID {
	case "student":
		student, _ := s.studentRepo.FindByUserID(userID)
		refs, err = s.achievementReferenceRepo.FindByStudentID(student.ID)

	case "lecturer":
		lecturer, _ := s.lecturerRepo.FindByUserID(userID)
		refs, err = s.achievementReferenceRepo.FindByAdvisorID(lecturer.ID)

	default:
		refs, err = s.achievementReferenceRepo.FindAll()
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch achievements",
		})
	}

	return c.JSON(fiber.Map{
		"data": refs,
	})
}

func (s *AchievementService) Submit(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)

	ref, err := s.achievementReferenceRepo.FindByID(id)
	if err != nil || ref.Status != "draft" {
		return c.Status(403).JSON(fiber.Map{
			"message": "Cannot submit",
		})
	}

	student, _ := s.studentRepo.FindByUserID(userID)
	if ref.StudentID != student.ID {
		return c.Status(403).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}

	err = s.achievementReferenceRepo.UpdateStatus(id, "submitted", nil, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Submit failed",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Submitted",
	})
}

func (s *AchievementService) Verify(c *fiber.Ctx) error {
	id := c.Params("id")
	lecturerID := c.Locals("user_id").(string)

	err := s.achievementReferenceRepo.UpdateStatus(
		id,
		"verified",
		&lecturerID,
		nil,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Verify failed",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Verified",
	})
}

func (s *AchievementService) Reject(c *fiber.Ctx) error {
	id := c.Params("id")
	lecturerID := c.Locals("user_id").(string)

	var body struct {
		Note string `json:"note"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	err := s.achievementReferenceRepo.UpdateStatus(
		id,
		"rejected",
		&lecturerID,
		&body.Note,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Reject failed",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Rejected",
	})
}

// Statistik achievement (ADMIN / DOSEN)
func (s *AchievementService) Statistics(c *fiber.Ctx) error {
	// TODO: implementasi real (count by status, dsb)

	data := map[string]int{
		"draft":     0,
		"submitted": 0,
		"verified":  0,
		"rejected":  0,
	}

	return helper.Success(c, "OK", data)
}

// Laporan achievement per mahasiswa
func (s *AchievementService) StudentReport(c *fiber.Ctx) error {
	userID := c.Params("id")

	// =========================
	// 1. AMBIL STUDENT DARI USER
	// =========================
	student, err := s.studentRepo.FindByUserID(userID)
	if err != nil {
		return helper.NotFound(c, "Mahasiswa tidak ditemukan")
	}

	// =========================
	// 2. AMBIL ACHIEVEMENT VERIFIED (POSTGRES)
	// =========================
	refs, err := s.achievementReferenceRepo.FindVerifiedByStudentID(student.ID)
	if err != nil {
		return helper.InternalServerError(c, "Gagal mengambil referensi prestasi")
	}

	// =========================
	// 3. AMBIL DETAIL ACHIEVEMENT (MONGO)
	// =========================
	var achievements []models.Achievement
	for _, ref := range refs {
		achievement, err := s.achievementRepo.FindByID(ref.MongoAchievementID)
		if err != nil {
			continue // skip kalau mongo hilang
		}

		achievements = append(achievements, *achievement)
	}

	// =========================
	// 4. RESPONSE
	// =========================
	return helper.Success(c, "OK", fiber.Map{
		"student_id":   student.ID,
		"total":        len(achievements),
		"achievements": achievements,
	})
}


func parseDetailsToAchievementDetail(
	raw map[string]interface{},
) models.AchievementDetail {

	detail := models.AchievementDetail{
		CustomFields: make(map[string]interface{}),
	}

	for k, v := range raw {
		switch k {

		case "competition_name":
			if s, ok := v.(string); ok {
				detail.CompetitionName = &s
			}

		case "competition_level":
			if s, ok := v.(string); ok {
				detail.CompetitionLevel = &s
			}

		case "rank":
			if f, ok := v.(float64); ok { // JSON number → float64
				r := int(f)
				detail.Rank = &r
			}

		case "medal_type":
			if s, ok := v.(string); ok {
				detail.MedalType = &s
			}

		case "publication_type":
			if s, ok := v.(string); ok {
				detail.PublicationType = &s
			}

		case "publication_title":
			if s, ok := v.(string); ok {
				detail.PublicationTitle = &s
			}

		case "authors":
			if arr, ok := v.([]interface{}); ok {
				for _, a := range arr {
					if s, ok := a.(string); ok {
						detail.Authors = append(detail.Authors, s)
					}
				}
			}

		case "organization_name":
			if s, ok := v.(string); ok {
				detail.OrganizationName = &s
			}

		case "position":
			if s, ok := v.(string); ok {
				detail.Position = &s
			}

		case "location":
			if s, ok := v.(string); ok {
				detail.Location = &s
			}

		default:
			// semua field tambahan masuk ke custom_fields
			detail.CustomFields[k] = v
		}
	}

	return detail
}

