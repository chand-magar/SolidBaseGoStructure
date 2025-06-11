package interfaces

import (
	"github.com/chand-magar/SolidBaseGoStructure/internal/dto"
	"github.com/google/uuid"
)

type UserService interface {
	Create(data dto.RequestDTO) (uuid.UUID, error)
	GetAll(params dto.PaginationParams) ([]dto.ResponseDTO, int64, int, error)
	FindOne(id uuid.UUID) (*dto.ResponseDTO, error)
	Update(id uuid.UUID, data dto.RequestDTO) error
}

type UserRepository interface {
	Create(user dto.RequestDTO) (uuid.UUID, error)
	GetAll(params dto.PaginationParams) ([]dto.ResponseDTO, int64, error)
	FindOne(id uuid.UUID) (*dto.ResponseDTO, error)
	Update(id uuid.UUID, data dto.RequestDTO) error
}
