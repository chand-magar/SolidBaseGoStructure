package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ProfileNo    uint32     `json:"profile_no" gorm:"primaryKey;autoIncrement;"`
	ProfileId    uuid.UUID  `json:"profile_id" gorm:"type:uuid;index"`
	RoleNo       uint32     `json:"role_no" gorm:"index;"`
	UserFullName string     `json:"user_fullname" gorm:"type:varchar(65)"`
	EmailId      string     `json:"email_id" gorm:"type:varchar(65)"`
	Gender       string     `json:"gender" gorm:"type:varchar(65);default:NULL"`
	Dob          *time.Time `json:"dob" gorm:"type:date;default:NULL"`
	MobileNo     string     `json:"mobile_no" gorm:"type:varchar(15);default:NULL"`
	Address      string     `json:"address" gorm:"type:jsonb;default:'{}'"`
	XApiKey      string     `json:"x_api_key" gorm:"type:varchar(55);default:NULL"`
	SecretKey    string     `json:"secret_key" gorm:"type:varchar(55);default:NULL"`
	Status       StatusEnum `json:"status" gorm:"type:status_enum;default:'A';index"`
	CreatedAt    time.Time  `json:"created_at" gorm:"index;default:NULL"`
	CreatedBy    uint32     `json:"created_by" gorm:"index;default:NULL"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"index;default:NULL"`
	UpdatedBy    uint32     `json:"updated_by" gorm:"index;default:NULL"`
}

// TableName specifies the custom table name for the User model
func (User) TableName() string {
	return "master.users"
}
