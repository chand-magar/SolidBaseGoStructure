package models

import (
	"time"

	"github.com/google/uuid"
)

type Page struct {
	PageNo    uint8      `json:"page_no" gorm:"primary_key; autoIncrement;"`
	PageId    uuid.UUID  `json:"page_id" gorm:"type:uuid;index;"`
	SectionNo uint8      `json:"section_no" gorm:"index;"`
	PageName  string     `json:"page_name" gorm:"type:varchar(65)"`
	PagePath  string     `json:"page_path" gorm:"type:varchar(255)"`
	PageOrder uint8      `json:"page_order" gorm:"index"`
	Status    StatusEnum `json:"status" gorm:"type:status_enum;default:'A';index"`
	CreatedAt time.Time  `json:"created_at" gorm:"index;default:NULL"`
	CreatedBy uint32     `json:"created_by" gorm:"index;default:NULL"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"index;default:NULL"`
	UpdatedBy uint32     `json:"updated_by" gorm:"index;default:NULL"`
}

// TableName specifies the custom table name for the Page model
func (Page) TableName() string {
	return "master.pages"
}
