package repository

import (
	"digital-wallet/internal/model"
	"digital-wallet/pkg/errs"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type UserRepo interface {
	CreateUser(user *model.User) error
	GetUserByID(userId string) (*model.User, error)
	SetUserTier(userId string, tierId string) error
	GetUsersByTierID(tierId string, page int, limit int) ([]model.User, error)
	GetTotalUsersByTierID(tierId string) (int64, error)
	GetUsers(page int, limit int) ([]model.User, error)
	GetTotalUsers() (int64, error)
	DeleteUser(userId string) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) DeleteUser(userId string) error {
	return r.db.Where("id = ?", userId).Delete(&model.User{}).Error
}

func (r *userRepo) GetUsersAccounts(userIds []string) ([]model.Account, error) {
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
		return nil, errs.NewInternalError("failed to get accounts", err)
	}

	return accounts, nil
}

func (r *userRepo) GetUsersList(users []model.User) ([]model.User, error) {
	if len(users) == 0 {
		return users, nil
	}
	// Get the list of user IDs for the query
	userIds := make([]string, len(users))
	for i, user := range users {
		userIds[i] = fmt.Sprintf("'%s'", user.ID) // Prepare for IN clause
	}

	accounts, err := r.GetUsersAccounts(userIds)
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

func (r *userRepo) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) GetUserByID(userId string) (*model.User, error) {
	var user model.User

	// Fetch the user from the public schema
	err := r.db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}

	accounts, err := r.GetUsersAccounts([]string{userId})
	if err != nil {
		return nil, err
	}

	user.Accounts = accounts

	return &user, nil
}

func (r *userRepo) SetUserTier(userId string, tierId string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userId).Update("tier_id", tierId).Error
}

func (r *userRepo) GetUsersByTierID(tierId string, page int, limit int) ([]model.User, error) {
	var users []model.User

	// Fetch the users who belong to the specified tier
	err := r.db.Where("tier_id = ?", tierId).Offset((page - 1) * limit).Limit(page).Find(&users).Error
	if err != nil {
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

	accounts, err := r.GetUsersAccounts(userIds)
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

func (r *userRepo) GetTotalUsersByTierID(tierId string) (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("tier_id = ?", tierId).Count(&count).Error
	return count, err
}

func (r *userRepo) GetUsers(page int, limit int) ([]model.User, error) {
	var users []model.User

	err := r.db.Offset((page - 1) * limit).Limit(page).Find(&users).Error
	if err != nil {
		return nil, err
	}

	users, err = r.GetUsersList(users)

	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepo) GetTotalUsers() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	return count, err
}
