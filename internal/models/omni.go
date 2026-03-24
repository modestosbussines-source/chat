package models

import "github.com/google/uuid"

// OmniCategory represents a group of quick-reply (Omni) scripts
type OmniCategory struct {
	BaseModel
	OrganizationID uuid.UUID      `gorm:"type:uuid;not null;index"`
	Name           string         `gorm:"type:varchar(255);not null"`
	Color          string         `gorm:"type:varchar(50);not null"`
	Scripts        []OmniScript   `gorm:"foreignKey:CategoryID"`
	Sequences      []OmniSequence `gorm:"foreignKey:CategoryID"`
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

// OmniSequence represents a flow of multiple messages to send in sequence
type OmniSequence struct {
	BaseModel
	OrganizationID uuid.UUID          `gorm:"type:uuid;not null;index"`
	CategoryID     uuid.UUID          `gorm:"type:uuid;not null;index"`
	Category       OmniCategory       `gorm:"foreignKey:CategoryID"`
	Title          string             `gorm:"type:varchar(255);not null"`
	Description    string             `gorm:"type:text"`
	Steps          []OmniSequenceStep `gorm:"foreignKey:SequenceID"`
}

// OmniSequenceStep represents a single step in an omni sequence
type OmniSequenceStep struct {
	BaseModel
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null;index"`
	SequenceID     uuid.UUID    `gorm:"type:uuid;not null;index"`
	Sequence       OmniSequence `gorm:"foreignKey:SequenceID"`
	StepOrder      int          `gorm:"type:int;not null"`
	Title          string       `gorm:"type:varchar(255);not null"`
	Content        string       `gorm:"type:text;not null"`
	MediaType      string       `gorm:"type:varchar(50);not null"` // text, audio, video, image, document
	FileName       *string      `gorm:"type:varchar(255)"`
	DelayMs        int          `gorm:"type:int;default:0"` // delay after sending this step
}
