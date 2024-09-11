package service

import (
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/validator"
)

type UserService interface {
	CreateUser(req *CreateUserRequest) (*model.User, error)
	GetUserByID(userId string) (*model.User, error)
	SetUserTier(userId string, tierId string) (*model.User, error)
	GetUsersByTierID(tierId string, page int, limit int) (*api.List[model.User], error)
	GetUsers(page int, limit int) (*api.List[model.User], error)
	DeleteUser(userId string) error
}

type userService struct {
	repos *repository.Repos
}

func NewUserService(repos *repository.Repos) UserService {
	return &userService{repos: repos}
}

func (s *userService) CreateUser(req *CreateUserRequest) (*model.User, error) {
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		return nil, errs.NewValidationError("invalid user request", fields)
	}
	user, _ := s.repos.User.GetUserByID(req.ID)
	if user != nil {
		return nil, errs.NewConflictError("user already exists", nil)
	}
	user = &model.User{
		ID:     req.ID,
		TierID: req.TierID,
	}
	if err := s.repos.User.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByID(userId string) (*model.User, error) {
	user, err := s.repos.User.GetUserByID(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) SetUserTier(userId string, tierId string) (*model.User, error) {
	user, err := s.repos.User.GetUserByID(userId)
	if err != nil {
		return nil, err
	}
	user.TierID = tierId
	if err := s.repos.User.SetUserTier(userId, tierId); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUsersByTierID(tierId string, page int, limit int) (*api.List[model.User], error) {
	users, err := s.repos.User.GetUsersByTierID(tierId, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.User.GetTotalUsersByTierID(tierId)
	if err != nil {
		return nil, err
	}

	return &api.List[model.User]{Items: users, Total: total, Page: page, Limit: limit}, nil
}

func (s *userService) GetUsers(page int, limit int) (*api.List[model.User], error) {
	users, err := s.repos.User.GetUsers(page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.User.GetTotalUsers()
	if err != nil {
		return nil, err
	}

	return &api.List[model.User]{Items: users, Total: total, Page: page, Limit: limit}, nil
}

func (s *userService) DeleteUser(userId string) error {
	return s.repos.User.DeleteUser(userId)
}
