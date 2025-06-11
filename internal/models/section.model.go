package models

import (
	"time"

	"github.com/google/uuid"
)

type Section struct {
	SectionNo    uint8      `json:"section_no" gorm:"primary_key; autoIncrement;"`
	SectionId    uuid.UUID  `json:"section_id" gorm:"type:uuid;index;"`
	SectionName  string     `json:"section_name" gorm:"type:varchar(65)"`
	SectionPath  string     `json:"section_path" gorm:"type:varchar(65)"`
	SectionIcon  string     `json:"section_icon" gorm:"type:varchar(65)"`
	SectionOrder uint8      `json:"section_order" gorm:"index"`
	Status       StatusEnum `json:"status" gorm:"type:status_enum;default:'A';index"`
	CreatedAt    time.Time  `json:"created_at" gorm:"index;default:NULL"`
	CreatedBy    uint32     `json:"created_by" gorm:"index;default:NULL"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"index;default:NULL"`
	UpdatedBy    uint32     `json:"updated_by" gorm:"index;default:NULL"`
}

// TableName specifies the custom table name for the Section model
func (Section) TableName() string {
	return "master.sections"
}
