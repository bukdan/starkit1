package service

import (
	"errors"
	"user-service/model"
	"user-service/repository"
)

type UserService interface {
	GetUserByID(userID string) (*model.User, error)
	UpdateProfile(user *model.User) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetUserByID(userID string) (*model.User, error) {
	return s.repo.FindByID(userID)
}

func (s *userService) UpdateProfile(user *model.User) error {
	if user.ID == "" {
		return errors.New("user id required")
	}
	return s.repo.Update(user)
}
