package models

import (
	"time"

	"github.com/google/uuid"
)

type UsersCredentials struct {
	CredentialNo uint32     `json:"credential_no" gorm:"primary_key;autoIncrement;"`
	CredentialId uuid.UUID  `json:"credential_id" gorm:"type:uuid;index"`
	ProfileNo    uint32     `json:"profile_no" gorm:"foreignKey:ProfileNo;index"` // Foreign key field referencing Users' primary key
	Username     string     `json:"username" gorm:"type:varchar(65);index"`
	Password     string     `json:"password" gorm:"type:varchar(255);index"`
	Status       StatusEnum `json:"status" gorm:"type:status_enum;default:'I';index"`
	CreatedAt    time.Time  `json:"created_at" gorm:"index;default:NULL"`
	CreatedBy    uint32     `json:"created_by" gorm:"index;default:NULL"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"index;default:NULL"`
	UpdatedBy    uint32     `json:"updated_by" gorm:"index;default:NULL"`
}

// TableName specifies the custom table name for the User model
func (UsersCredentials) TableName() string {
	return "master.user_credentials"
}
