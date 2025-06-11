package dto

import (
	"time"

	"github.com/chand-magar/SolidBaseGoStructure/internal/models"
	"github.com/google/uuid"
)

type RequestDTO struct {
	ProfileId    uuid.UUID         `json:"profile_id"`
	RoleId       uuid.UUID         `json:"role_id" validate:"required"`
	UserFullName string            `json:"user_fullname" validate:"required"`
	Username     string            `json:"username" validate:"required"`
	Password     string            `json:"password" validate:"required"`
	EmailId      string            `json:"email_id" validate:"required"`
	Gender       string            `json:"gender"`
	Dob          *time.Time        `json:"dob"`
	MobileNo     string            `json:"mobile_no"`
	Address      string            `json:"address"`
	XApiKey      string            `json:"x_api_key"`
	SecretKey    string            `json:"secret_key"`
	Status       models.StatusEnum `json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
	CreatedBy    uint32            `json:"created_by"`
	UpdatedAt    time.Time         `json:"updated_at"`
	UpdatedBy    uint32            `json:"updated_by"`
}

type ResponseDTO struct {
	ProfileId    uuid.UUID         `json:"profile_id"`
	RoleId       uuid.UUID         `json:"role_id"`
	UserFullName string            `json:"user_fullname"`
	EmailId      string            `json:"email_id"`
	Gender       string            `json:"gender"`
	Dob          *time.Time        `json:"dob"`
	MobileNo     string            `json:"mobile_no"`
	Address      string            `json:"address"`
	XApiKey      string            `json:"x_api_key"`
	SecretKey    string            `json:"secret_key"`
	Status       models.StatusEnum `json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
	CreatedBy    uint32            `json:"created_by"`
	UpdatedAt    time.Time         `json:"updated_at"`
	UpdatedBy    uint32            `json:"updated_by"`
}
