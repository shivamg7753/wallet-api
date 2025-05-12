package services

import (
	"errors"

	"wallet-api/models"
	"wallet-api/repositories"
)

type UserServiceInterface interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
}

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Create(user *models.User) error {
	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("email already in use")
	}

	return s.userRepo.Create(user)
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}