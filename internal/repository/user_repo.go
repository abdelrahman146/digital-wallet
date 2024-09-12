package repository

import (
	"context"
	"digital-wallet/internal/model"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/logger"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, userId string) (*model.User, error)
	SetUserTier(ctx context.Context, userId string, tierId string) error
	GetUsersByTierID(ctx context.Context, tierId string, page int, limit int) ([]model.User, error)
	GetTotalUsersByTierID(ctx context.Context, tierId string) (int64, error)
	GetUsers(ctx context.Context, page int, limit int) ([]model.User, error)
	GetTotalUsers(ctx context.Context) (int64, error)
	DeleteUser(ctx context.Context, userId string) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) DeleteUser(ctx context.Context, userId string) error {
	return r.db.Where("id = ?", userId).Delete(&model.User{}).Error
}

func (r *userRepo) GetUsersAccounts(ctx context.Context, userIds []string) ([]model.Account, error) {
	var err error
	userIdsString := strings.Join(userIds, ", ")

	// Fetch all wallets to build the UNION query for accounts
	var wallets []model.Wallet
	err = r.db.Table("wallets").Find(&wallets).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets: %w", err)
	}

	// Build the UNION query to fetch accounts from different schemas
	var unionQueries []string
	for _, wallet := range wallets {
		schemaName := fmt.Sprintf("%s_wallet", wallet.ID)
		query := fmt.Sprintf("SELECT id, user_id, balance, version, '%s' as wallet FROM %s.accounts WHERE user_id IN (%s)", wallet.ID, schemaName, userIdsString)
		unionQueries = append(unionQueries, query)
	}

	// Join all UNION queries
	finalQuery := strings.Join(unionQueries, " UNION ALL ")

	// Fetch accounts with a single query
	var accounts []model.Account
	err = r.db.Raw(finalQuery).Scan(&accounts).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get accounts", logger.Field("error", err), logger.Field("userIds", userIds), logger.Field("finalQuery", finalQuery))
		return nil, err
	}

	return accounts, nil
}

func (r *userRepo) GetUsersList(ctx context.Context, users []model.User) ([]model.User, error) {
	if len(users) == 0 {
		return users, nil
	}
	// Get the list of user IDs for the query
	userIds := make([]string, len(users))
	for i, user := range users {
		userIds[i] = fmt.Sprintf("'%s'", user.ID) // Prepare for IN clause
	}

	accounts, err := r.GetUsersAccounts(ctx, userIds)
	if err != nil {
		return nil, err
	}
	// Create a map of user IDs to their accounts
	accountsMap := make(map[string][]model.Account)
	for _, account := range accounts {
		accountsMap[account.UserID] = append(accountsMap[account.UserID], account)
	}

	// Attach the accounts to the respective users
	for i := range users {
		users[i].Accounts = accountsMap[users[i].ID]
	}

	return users, nil
}

func (r *userRepo) CreateUser(ctx context.Context, user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		api.GetLogger(ctx).Error("failed to create user", logger.Field("error", err), logger.Field("user", user))
		return err
	}
	return nil
}

func (r *userRepo) GetUserByID(ctx context.Context, userId string) (*model.User, error) {
	var user model.User

	// Fetch the user from the public schema
	err := r.db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get user by id", logger.Field("error", err), logger.Field("userId", userId))
		return nil, err
	}

	accounts, err := r.GetUsersAccounts(ctx, []string{fmt.Sprintf("'%s'", userId)})
	if err != nil {
		api.GetLogger(ctx).Error("failed to get user accounts", logger.Field("error", err), logger.Field("userId", userId), logger.Field("accounts", accounts))
		return nil, err
	}

	user.Accounts = accounts

	return &user, nil
}

func (r *userRepo) SetUserTier(ctx context.Context, userId string, tierId string) error {
	if err := r.db.Model(&model.User{}).Where("id = ?", userId).Update("tier_id", tierId).Error; err != nil {
		api.GetLogger(ctx).Error("failed to set user tier", logger.Field("error", err), logger.Field("userId", userId), logger.Field("tierId", tierId))
	}
	return nil
}

func (r *userRepo) GetUsersByTierID(ctx context.Context, tierId string, page int, limit int) ([]model.User, error) {
	var users []model.User

	// Fetch the users who belong to the specified tier
	err := r.db.Where("tier_id = ?", tierId).Offset((page - 1) * limit).Limit(limit).Find(&users).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get users by tier id", logger.Field("error", err), logger.Field("tierId", tierId))
		return nil, err
	}

	// If no users were found, return early
	if len(users) == 0 {
		return users, nil
	}

	// Get the list of user IDs for the query
	userIds := make([]string, len(users))
	for i, user := range users {
		userIds[i] = fmt.Sprintf("'%s'", user.ID) // Prepare for IN clause
	}

	accounts, err := r.GetUsersAccounts(ctx, userIds)
	if err != nil {
		return nil, err
	}

	// Create a map of user IDs to their accounts
	accountsMap := make(map[string][]model.Account)
	for _, account := range accounts {
		accountsMap[account.UserID] = append(accountsMap[account.UserID], account)
	}

	// Attach the accounts to the respective users
	for i := range users {
		users[i].Accounts = accountsMap[users[i].ID]
	}

	return users, nil
}

func (r *userRepo) GetTotalUsersByTierID(ctx context.Context, tierId string) (int64, error) {
	var count int64
	if err := r.db.Model(&model.User{}).Where("tier_id = ?", tierId).Count(&count).Error; err != nil {
		api.GetLogger(ctx).Error("failed to get total users by tier id", logger.Field("error", err), logger.Field("tierId", tierId))
	}
	return count, nil
}

func (r *userRepo) GetUsers(ctx context.Context, page int, limit int) ([]model.User, error) {
	var users []model.User

	err := r.db.Offset((page - 1) * limit).Limit(limit).Find(&users).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get users", logger.Field("error", err))
		return nil, err
	}

	users, err = r.GetUsersList(ctx, users)

	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepo) GetTotalUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.Model(&model.User{}).Count(&count).Error; err != nil {
		api.GetLogger(ctx).Error("failed to get total users", logger.Field("error", err))
	}
	return count, nil
}
