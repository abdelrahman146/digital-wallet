package repository

import (
	"digital-wallet/internal/model"
	"digital-wallet/pkg/logger"
	"gorm.io/gorm"
)

type TierRepo interface {
	CreateTier(tier *model.Tier) error
	GetTierByID(tierId string) (*model.Tier, error)
	GetTiers(page int, limit int) ([]model.Tier, error)
	GetTotalTiers() (int64, error)
	DeleteTier(tierId string) error
}

type tierRepo struct {
	db *gorm.DB
}

func NewTierRepo(db *gorm.DB) TierRepo {
	return &tierRepo{db: db}
}

func (r *tierRepo) CreateTier(tier *model.Tier) error {
	if err := r.db.Create(tier).Error; err != nil {
		logger.GetLogger().Error("failed to create tier", logger.Field("error", err), logger.Field("tier", tier))
		return err
	}
	return nil
}

func (r *tierRepo) GetTierByID(tierId string) (*model.Tier, error) {
	var tier model.Tier
	err := r.db.Where("id = ?", tierId).First(&tier).Error
	if err != nil {
		logger.GetLogger().Error("failed to get tier by id", logger.Field("error", err), logger.Field("tierId", tierId))
		return nil, err
	}
	return &tier, nil
}

func (r *tierRepo) GetTiers(page int, limit int) ([]model.Tier, error) {
	var tiers []model.Tier
	err := r.db.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&tiers).Error
	if err != nil {
		logger.GetLogger().Error("failed to get tiers", logger.Field("error", err))
		return nil, err
	}
	return tiers, nil
}

func (r *tierRepo) GetTotalTiers() (int64, error) {
	var total int64
	err := r.db.Model(&model.Tier{}).Count(&total).Error
	if err != nil {
		logger.GetLogger().Error("failed to get total tiers", logger.Field("error", err))
		return 0, err
	}
	return total, nil
}

func (r *tierRepo) DeleteTier(tierId string) error {
	if err := r.db.Where("id = ?", tierId).Delete(&model.Tier{}).Error; err != nil {
		logger.GetLogger().Error("failed to delete tier", logger.Field("error", err), logger.Field("tierId", tierId))
	}
	return nil
}
