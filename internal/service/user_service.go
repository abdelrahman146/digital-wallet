package service

import (
	"context"
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/logger"
	"digital-wallet/pkg/validator"
)

type UserService interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*model.User, error)
	GetUserByID(ctx context.Context, userId string) (*model.User, error)
	SetUserTier(ctx context.Context, userId string, tierId string) (*model.User, error)
	GetUsersByTierID(ctx context.Context, tierId string, page int, limit int) (*api.List[model.User], error)
	GetUsers(ctx context.Context, page int, limit int) (*api.List[model.User], error)
	DeleteUser(ctx context.Context, userId string) error
}

type userService struct {
	repos *repository.Repos
}

func NewUserService(repos *repository.Repos) UserService {
	return &userService{repos: repos}
}

func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*model.User, error) {
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		api.GetLogger(ctx).Error("Invalid user request", logger.Field("fields", fields), logger.Field("request", req))
		return nil, errs.NewValidationError("Invalid user request", "", fields)
	}
	user, _ := s.repos.User.FetchUserByID(ctx, req.ID)
	if user != nil {
		api.GetLogger(ctx).Error("User already exists", logger.Field("userId", req.ID))
		return nil, errs.NewConflictError("User already exists", "USER_ALREADY_EXISTS", nil)
	}
	user = &model.User{
		ID:     req.ID,
		TierID: req.TierID,
	}
	if err := s.repos.User.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, userId string) (*model.User, error) {
	user, err := s.repos.User.FetchUserByID(ctx, userId)
	if user == nil {
		api.GetLogger(ctx).Error("User not found", logger.Field("userId", userId), logger.Field("error", err))
		return nil, errs.NewNotFoundError("User not found", "USER_NOT_FOUND", err)
	}
	return user, nil
}

func (s *userService) SetUserTier(ctx context.Context, userId string, tierId string) (*model.User, error) {
	user, err := s.repos.User.FetchUserByID(ctx, userId)
	if user == nil {
		api.GetLogger(ctx).Error("User not found", logger.Field("userId", userId), logger.Field("error", err))
		return nil, errs.NewNotFoundError("User not found", "USER_NOT_FOUND", err)
	}
	if err := s.repos.User.UpdateUserTier(ctx, userId, tierId); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUsersByTierID(ctx context.Context, tierId string, page int, limit int) (*api.List[model.User], error) {
	users, err := s.repos.User.FetchUsersByTierID(ctx, tierId, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.User.CountUsersByTierID(ctx, tierId)
	if err != nil {
		return nil, err
	}

	return &api.List[model.User]{Items: users, Total: total, Page: page, Limit: limit}, nil
}

func (s *userService) GetUsers(ctx context.Context, page int, limit int) (*api.List[model.User], error) {
	users, err := s.repos.User.FetchUsers(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.User.CountTotalUsers(ctx)
	if err != nil {
		return nil, err
	}

	return &api.List[model.User]{Items: users, Total: total, Page: page, Limit: limit}, nil
}

func (s *userService) DeleteUser(ctx context.Context, userId string) error {
	return s.repos.User.RemoveUser(ctx, userId)
}
