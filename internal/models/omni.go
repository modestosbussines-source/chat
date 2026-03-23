package models

import "github.com/google/uuid"

// OmniCategory represents a group of quick-reply (Omni) scripts
type OmniCategory struct {
	BaseModel
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name           string    `gorm:"type:varchar(255);not null"`
	Color          string    `gorm:"type:varchar(50);not null"`
	Scripts        []OmniScript `gorm:"foreignKey:CategoryID"`
}

// OmniScript represents a single point-and-click message/media payload
type OmniScript struct {
	BaseModel
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null;index"`
	CategoryID     uuid.UUID    `gorm:"type:uuid;not null;index"`
	Category       OmniCategory `gorm:"foreignKey:CategoryID"`
	Title          string       `gorm:"type:varchar(255);not null"`
	Content        string       `gorm:"type:text;not null"`
	MediaType      string       `gorm:"type:varchar(50);not null"` // text, audio, video, image
	FileName       *string      `gorm:"type:varchar(255)"`
	DelayMs        int          `gorm:"type:int;default:0"`
}
