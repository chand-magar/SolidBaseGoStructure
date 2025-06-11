package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	RoleNo      uint8      `json:"role_no" gorm:"primary_key; autoIncrement;"`
	RoleId      uuid.UUID  `json:"role_id" gorm:"type:uuid;index;"`
	RoleName    string     `json:"role_name" gorm:"type:varchar(65)"`
	RoleDetails string     `json:"role_details" gorm:"type:jsonb;default:'[]'"`
	Status      StatusEnum `json:"status" gorm:"type:status_enum;default:'A';index"`
	CreatedAt   time.Time  `json:"created_at" gorm:"index;default:NULL"`
	CreatedBy   uint32     `json:"created_by" gorm:"index;default:NULL"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"index;default:NULL"`
	UpdatedBy   uint32     `json:"updated_by" gorm:"index;default:NULL"`
}

// TableName specifies the custom table name for the Role model
func (Role) TableName() string {
	return "master.roles"
}
