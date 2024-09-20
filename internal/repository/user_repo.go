package repository

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/internal/resource"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
)

type UserRepo interface {
	// CreateUser Creates a new user
	CreateUser(ctx context.Context, user *model.User) error
	// UpdateUserTier Sets the user's tier by user ID and tier ID
	UpdateUserTier(ctx context.Context, userId string, tierId string) error
	// DeleteUser Deletes a user by user ID
	DeleteUser(ctx context.Context, user *model.User) error
	// FetchUserByID Retrieves a user by user ID
	FetchUserByID(ctx context.Context, userId string) (*model.User, error)
	// FetchUsersByTierID Retrieves users by their tier ID with pagination
	FetchUsersByTierID(ctx context.Context, tierId string, page int, limit int) ([]model.User, error)
	// CountUsersByTierID Retrieves the total number of users in a given tier
	CountUsersByTierID(ctx context.Context, tierId string) (int64, error)
	// FetchUsers Retrieves all users with pagination
	FetchUsers(ctx context.Context, page int, limit int) ([]model.User, error)
	// CountTotalUsers Retrieves the total number of users
	CountTotalUsers(ctx context.Context) (int64, error)
}

type userRepo struct {
	resources *resource.Resources
}

// NewUserRepo initializes the user repository
func NewUserRepo(resources *resource.Resources) UserRepo {
	return &userRepo{resources: resources}
}

// CreateUser creates a new user in the database
func (r *userRepo) CreateUser(ctx context.Context, user *model.User) error {
	if err := r.resources.DB.Create(user).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to create user", logger.Field("error", err), logger.Field("user", user))
		return err
	}
	return nil
}

// FetchUserByID retrieves a user by user ID and preloads related accounts
func (r *userRepo) FetchUserByID(ctx context.Context, userId string) (*model.User, error) {
	var user model.User
	err := r.resources.DB.Where("id = ?", userId).Preload("Accounts").First(&user).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve user by ID", logger.Field("error", err), logger.Field("userId", userId))
		return nil, err
	}
	return &user, nil
}

// UpdateUserTier sets the user's tier by user ID and tier ID
func (r *userRepo) UpdateUserTier(ctx context.Context, userId string, tierId string) error {
	if err := r.resources.DB.Model(&model.User{}).Where("id = ?", userId).Update("tier_id", tierId).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to update user tier", logger.Field("error", err), logger.Field("userId", userId), logger.Field("tierId", tierId))
		return err
	}
	return nil
}

// FetchUsersByTierID retrieves users by their tier ID with pagination and preloads related accounts
func (r *userRepo) FetchUsersByTierID(ctx context.Context, tierId string, page int, limit int) ([]model.User, error) {
	var users []model.User
	err := r.resources.DB.Where("tier_id = ?", tierId).Offset((page - 1) * limit).Limit(limit).Preload("Accounts").Find(&users).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve users by tier ID", logger.Field("error", err), logger.Field("tierId", tierId))
		return nil, err
	}
	return users, nil
}

// CountUsersByTierID retrieves the total number of users in a given tier
func (r *userRepo) CountUsersByTierID(ctx context.Context, tierId string) (int64, error) {
	var count int64
	err := r.resources.DB.Model(&model.User{}).Where("tier_id = ?", tierId).Count(&count).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to count users by tier ID", logger.Field("error", err), logger.Field("tierId", tierId))
		return 0, err
	}
	return count, nil
}

// FetchUsers retrieves all users with pagination and preloads related accounts
func (r *userRepo) FetchUsers(ctx context.Context, page int, limit int) ([]model.User, error) {
	var users []model.User
	err := r.resources.DB.Offset((page - 1) * limit).Limit(limit).Preload("Accounts").Find(&users).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve users", logger.Field("error", err))
		return nil, err
	}
	return users, nil
}

// CountTotalUsers retrieves the total number of users
func (r *userRepo) CountTotalUsers(ctx context.Context) (int64, error) {
	var count int64
	err := r.resources.DB.Model(&model.User{}).Count(&count).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to count total users", logger.Field("error", err))
		return 0, err
	}
	return count, nil
}

// DeleteUser deletes a user by user ID
func (r *userRepo) DeleteUser(ctx context.Context, user *model.User) error {
	if err := r.resources.DB.Delete(user).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to delete user", logger.Field("error", err), logger.Field("user", user))
		return err
	}
	return nil
}
