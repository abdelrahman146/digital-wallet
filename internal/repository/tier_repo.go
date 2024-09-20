package repository

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/internal/resource"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
)

type TierRepo interface {
	CreateTier(ctx context.Context, tier *model.Tier) error
	GetTierByID(ctx context.Context, tierId string) (*model.Tier, error)
	GetTiers(ctx context.Context, page int, limit int) ([]model.Tier, error)
	GetTotalTiers(ctx context.Context) (int64, error)
	DeleteTier(ctx context.Context, tier *model.Tier) error
}

type tierRepo struct {
	resources *resource.Resources
}

func NewTierRepo(resources *resource.Resources) TierRepo {
	return &tierRepo{resources: resources}
}

func (r *tierRepo) CreateTier(ctx context.Context, tier *model.Tier) error {
	if err := r.resources.DB.Create(tier).Error; err != nil {
		api.GetLogger(ctx).Error("failed to create tier", logger.Field("error", err), logger.Field("tier", tier))
		return err
	}
	return nil
}

func (r *tierRepo) GetTierByID(ctx context.Context, tierId string) (*model.Tier, error) {
	var tier model.Tier
	err := r.resources.DB.Where("id = ?", tierId).First(&tier).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get tier by id", logger.Field("error", err), logger.Field("tierId", tierId))
		return nil, err
	}
	return &tier, nil
}

func (r *tierRepo) GetTiers(ctx context.Context, page int, limit int) ([]model.Tier, error) {
	var tiers []model.Tier
	err := r.resources.DB.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&tiers).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get tiers", logger.Field("error", err))
		return nil, err
	}
	return tiers, nil
}

func (r *tierRepo) GetTotalTiers(ctx context.Context) (int64, error) {
	var total int64
	err := r.resources.DB.Model(&model.Tier{}).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("failed to get total tiers", logger.Field("error", err))
		return 0, err
	}
	return total, nil
}

func (r *tierRepo) DeleteTier(ctx context.Context, tier *model.Tier) error {
	if err := r.resources.DB.Delete(tier).Error; err != nil {
		api.GetLogger(ctx).Error("failed to delete tier", logger.Field("error", err), logger.Field("tier", tier))
		return err
	}
	return nil
}
