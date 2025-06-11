package services

import (
	"fmt"
	"math"

	"github.com/chand-magar/SolidBaseGoStructure/internal/dto"
	"github.com/chand-magar/SolidBaseGoStructure/internal/interfaces"
	"github.com/google/uuid"
)

type userService struct {
	repo interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) interfaces.UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(data dto.RequestDTO) (uuid.UUID, error) {

	if data.UserFullName == "" {
		return uuid.Nil, fmt.Errorf("User full name is required")
	}
	return s.repo.Create(data)
}

func (s *userService) GetAll(params dto.PaginationParams) ([]dto.ResponseDTO, int64, int, error) {
	users, totalRecords, err := s.repo.GetAll(params)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(params.Size)))
	return users, totalRecords, totalPages, nil
}

func (s *userService) FindOne(id uuid.UUID) (*dto.ResponseDTO, error) {
	return s.repo.FindOne(id)
}

func (s *userService) Update(id uuid.UUID, data dto.RequestDTO) error {
	return s.repo.Update(id, data)
}
