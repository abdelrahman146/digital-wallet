package service

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/internal/repository"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/errs"
	"github.com/abdelrahman146/digital-wallet/pkg/validator"
)

type TierService interface {
	CreateTier(ctx context.Context, req *CreateTierRequest) (*model.Tier, error)
	GetTierByID(ctx context.Context, tierId string) (*model.Tier, error)
	GetTiers(ctx context.Context, page int, limit int) (*api.List[model.Tier], error)
	DeleteTier(ctx context.Context, tierId string) error
}

type tierService struct {
	repos *repository.Repos
}

func NewTierService(repos *repository.Repos) TierService {
	return &tierService{repos: repos}
}

func (s *tierService) CreateTier(ctx context.Context, req *CreateTierRequest) (*model.Tier, error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return nil, err
	}
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		return nil, errs.NewValidationError("Invalid tier request", "", fields)
	}
	tier := &model.Tier{
		ID:   req.ID,
		Name: req.Name,
	}
	tier.SetActor(api.GetActor(ctx), api.GetActorID(ctx))
	tier.SetRemarks("Tier created")
	if err := s.repos.Tier.CreateTier(ctx, tier); err != nil {
		return nil, err
	}
	return tier, nil
}

func (s *tierService) GetTierByID(ctx context.Context, tierId string) (*model.Tier, error) {
	tier, err := s.repos.Tier.FetchTierByID(ctx, tierId)
	if err != nil {
		return nil, err
	}
	return tier, nil
}

func (s *tierService) GetTiers(ctx context.Context, page int, limit int) (*api.List[model.Tier], error) {
	tiers, err := s.repos.Tier.FetchTiers(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Tier.CountTiers(ctx)
	if err != nil {
		return nil, err
	}

	return &api.List[model.Tier]{Items: tiers, Total: total, Page: page, Limit: limit}, nil
}

func (s *tierService) DeleteTier(ctx context.Context, tierId string) error {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return err
	}
	tier, err := s.repos.Tier.FetchTierByID(ctx, tierId)
	if err != nil {
		return err
	}
	tier.SetActor(api.GetActor(ctx), api.GetActorID(ctx))
	tier.SetRemarks("Tier deleted")
	tier.SetOldRecord(tier)
	return s.repos.Tier.DeleteTier(ctx, tier)
}
