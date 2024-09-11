package service

import (
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
	"digital-wallet/pkg/errs"
	"digital-wallet/pkg/validator"
)

type TierService interface {
	CreateTier(req *CreateTierRequest) (*model.Tier, error)
	GetTierByID(tierId string) (*model.Tier, error)
	GetTiers(page int, limit int) (*api.List[model.Tier], error)
	DeleteTier(tierId string) error
}

type tierService struct {
	repos *repository.Repos
}

func NewTierService(repos *repository.Repos) TierService {
	return &tierService{repos: repos}
}

func (s *tierService) CreateTier(req *CreateTierRequest) (*model.Tier, error) {
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		return nil, errs.NewValidationError("invalid tier request", fields)
	}
	tier := &model.Tier{
		ID:   req.ID,
		Name: req.Name,
	}
	if err := s.repos.Tier.CreateTier(tier); err != nil {
		return nil, err
	}
	return tier, nil
}

func (s *tierService) GetTierByID(tierId string) (*model.Tier, error) {
	tier, err := s.repos.Tier.GetTierByID(tierId)
	if err != nil {
		return nil, err
	}
	return tier, nil
}

func (s *tierService) GetTiers(page int, limit int) (*api.List[model.Tier], error) {
	tiers, err := s.repos.Tier.GetTiers(page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Tier.GetTotalTiers()
	if err != nil {
		return nil, err
	}

	return &api.List[model.Tier]{Items: tiers, Total: total, Page: page, Limit: limit}, nil
}

func (s *tierService) DeleteTier(tierId string) error {
	return s.repos.Tier.DeleteTier(tierId)
}
