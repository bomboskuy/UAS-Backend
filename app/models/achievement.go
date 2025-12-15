package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Achievement struct {
	ID              primitive.ObjectID        `bson:"_id,omitempty" json:"id"`
	StudentID       string        `bson:"studentId" json:"student_id"`
	AchievementType string        `bson:"achievementType" json:"achievement_type"`
	Title           string        `bson:"title" json:"title"`
	Description     string        `bson:"description" json:"description"`
	Details         AchievementDetail `bson:"details" json:"details"`
	Attachments     []Attachment  `bson:"attachments" json:"attachments"`
	Tags            []string      `bson:"tags" json:"tags"`
	Points          int           `bson:"points" json:"points"`
	CreatedAt       time.Time     `bson:"createdAt" json:"created_at"`
	UpdatedAt       time.Time     `bson:"updatedAt" json:"updated_at"`
}

type AchievementDetail struct {
	CompetitionName   *string   `bson:"competitionName,omitempty" json:"competition_name,omitempty"`
	CompetitionLevel  *string   `bson:"competitionLevel,omitempty" json:"competition_level,omitempty"`
	Rank              *int      `bson:"rank,omitempty" json:"rank,omitempty"`
	MedalType         *string   `bson:"medalType,omitempty" json:"medal_type,omitempty"`

	PublicationType   *string   `bson:"publicationType,omitempty" json:"publication_type,omitempty"`
	PublicationTitle  *string   `bson:"publicationTitle,omitempty" json:"publication_title,omitempty"`
	Authors           []string  `bson:"authors,omitempty" json:"authors,omitempty"`

	OrganizationName  *string   `bson:"organizationName,omitempty" json:"organization_name,omitempty"`
	Position          *string   `bson:"position,omitempty" json:"position,omitempty"`

	EventDate         *time.Time `bson:"eventDate,omitempty" json:"event_date,omitempty"`
	Location          *string    `bson:"location,omitempty" json:"location,omitempty"`

	CustomFields      map[string]interface{} `bson:"customFields,omitempty" json:"custom_fields,omitempty"`
}

type Attachment struct {
	FileName   string    `bson:"fileName" json:"file_name"`
	FileURL    string    `bson:"fileUrl" json:"file_url"`
	FileType   string    `bson:"fileType" json:"file_type"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploaded_at"`
}
