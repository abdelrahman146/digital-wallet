package repository

import (
	"digital-wallet/internal/model"
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
	return r.db.Create(tier).Error
}

func (r *tierRepo) GetTierByID(tierId string) (*model.Tier, error) {
	var tier model.Tier
	err := r.db.Where("id = ?", tierId).First(&tier).Error
	if err != nil {
		return nil, err
	}
	return &tier, nil
}

func (r *tierRepo) GetTiers(page int, limit int) ([]model.Tier, error) {
	var tiers []model.Tier
	err := r.db.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&tiers).Error
	if err != nil {
		return nil, err
	}
	return tiers, nil
}

func (r *tierRepo) GetTotalTiers() (int64, error) {
	var total int64
	err := r.db.Model(&model.Tier{}).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *tierRepo) DeleteTier(tierId string) error {
	return r.db.Where("id = ?", tierId).Delete(&model.Tier{}).Error
}
