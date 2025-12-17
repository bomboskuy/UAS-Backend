package models

import "time"

type AchievementReference struct {
	ID                 string     `db:"id" json:"id"`
	StudentID          string     `db:"student_id" json:"student_id"`
	MongoAchievementID string     `db:"mongo_achievement_id" json:"mongo_achievement_id"`
	Status             string     `db:"status" json:"status"`
	SubmittedAt        *time.Time `db:"submitted_at" json:"submitted_at,omitempty"`
	VerifiedAt         *time.Time `db:"verified_at" json:"verified_at,omitempty"`
	VerifiedBy         *string    `db:"verified_by" json:"verified_by,omitempty"`
	RejectionNote      *string    `db:"rejection_note" json:"rejection_note,omitempty"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at" json:"updated_at"`
}

type CreateAchievementRequest struct {
	AchievementType string                 `json:"achievement_type"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Details         map[string]interface{} `json:"details"`
	Tags            []string               `json:"tags"`
	Points          int                    `json:"points"`
}

type AchievementResponse struct {
	ID                 string                 `json:"id"`
	StudentID          string                 `json:"student_id"`
	MongoAchievementID string                 `json:"mongo_achievement_id"`
	AchievementType    string                 `json:"achievement_type"`
	Title              string                 `json:"title"`
	Description        string                 `json:"description"`
	Details            map[string]interface{} `json:"details"`
	Tags               []string               `json:"tags"`
	Points             int                    `json:"points"`
	Status             string                 `json:"status"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

type StatisticsResponse struct {
	TotalAchievements        int                       `json:"total_achievements"`
	AchievementsByType       map[string]int            `json:"achievements_by_type"`
	AchievementsByStatus     map[string]int            `json:"achievements_by_status"`
	AchievementsByPeriod     map[string]int            `json:"achievements_by_period"`
	TopStudents              []TopStudentResponse      `json:"top_students"`
	CompetitionDistribution  map[string]int            `json:"competition_distribution"`
}

type TopStudentResponse struct {
	StudentID    string `json:"student_id"`
	StudentName  string `json:"student_name"`
	TotalPoints  int    `json:"total_points"`
	TotalAchievements int `json:"total_achievements"`
}
 