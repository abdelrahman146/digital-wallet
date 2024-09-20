package service

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/internal/repository"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/errs"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
	"github.com/abdelrahman146/digital-wallet/pkg/validator"
)

type TriggerService interface {
	CreateTrigger(ctx context.Context, req CreateTriggerRequest) (*model.Trigger, error)
	UpdateTrigger(ctx context.Context, triggerId uint64, req UpdateTriggerRequest) (*model.Trigger, error)
	DeleteTrigger(ctx context.Context, triggerId uint64) error
	GetTrigger(ctx context.Context, triggerId uint64) (*model.Trigger, error)
	ListTriggers(ctx context.Context, page, limit int) (*api.List[model.Trigger], error)
}

type triggerService struct {
	repos *repository.Repos
}

func NewTriggerService(repos *repository.Repos) TriggerService {
	return &triggerService{repos: repos}
}

func (s *triggerService) CreateTrigger(ctx context.Context, req CreateTriggerRequest) (*model.Trigger, error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("Unauthorized access", logger.Field("error", err))
		return nil, err
	}
	if err := validator.GetValidator().ValidateStruct(req); err != nil {
		fields := validator.GetValidator().GetValidationErrors(err)
		api.GetLogger(ctx).Error("Invalid request", logger.Field("fields", fields), logger.Field("request", req))
		return nil, errs.NewValidationError("Invalid request", "", fields)
	}
	trigger := &model.Trigger{
		Name:       req.Name,
		Slug:       req.Slug,
		Properties: req.Properties,
	}
	trigger.SetActor(api.GetActor(ctx), api.GetActorID(ctx))
	trigger.SetRemarks("Trigger created")
	if err := s.repos.Trigger.CreateTrigger(ctx, trigger); err != nil {
		return nil, err
	}
	return trigger, nil
}

func (s *triggerService) UpdateTrigger(ctx context.Context, triggerId uint64, req UpdateTriggerRequest) (*model.Trigger, error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("Unauthorized access", logger.Field("error", err))
		return nil, err
	}
	trigger, err := s.repos.Trigger.FetchTriggerByID(ctx, triggerId)
	if trigger == nil {
		return nil, errs.NewNotFoundError("Trigger not found", "TRIGGER_NOT_FOUND", err)
	}
	trigger.SetActor(api.GetActor(ctx), api.GetActorID(ctx))
	trigger.SetRemarks("Trigger updated")
	trigger.SetOldRecord(*trigger)
	if req.Name != nil {
		trigger.Name = *req.Name
	}
	if req.Slug != nil {
		trigger.Slug = *req.Slug
	}
	if req.Properties != nil {
		trigger.Properties = *req.Properties
	}
	if err := s.repos.Trigger.UpdateTrigger(ctx, trigger); err != nil {
		return nil, err
	}
	return trigger, nil
}

func (s *triggerService) DeleteTrigger(ctx context.Context, triggerId uint64) error {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("Unauthorized access", logger.Field("error", err))
		return err
	}
	trigger, err := s.repos.Trigger.FetchTriggerByID(ctx, triggerId)
	if trigger == nil {
		return errs.NewNotFoundError("Trigger not found", "TRIGGER_NOT_FOUND", err)
	}
	if err := s.repos.Trigger.DeleteTrigger(ctx, trigger); err != nil {
		return err
	}
	return nil
}

func (s *triggerService) GetTrigger(ctx context.Context, triggerId uint64) (*model.Trigger, error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("Unauthorized access", logger.Field("error", err))
		return nil, err
	}
	trigger, err := s.repos.Trigger.FetchTriggerByID(ctx, triggerId)
	if trigger == nil {
		return nil, errs.NewNotFoundError("Trigger not found", "TRIGGER_NOT_FOUND", err)
	}
	return trigger, nil
}

func (s *triggerService) ListTriggers(ctx context.Context, page, limit int) (*api.List[model.Trigger], error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("Unauthorized access", logger.Field("error", err))
		return nil, err
	}
	triggers, err := s.repos.Trigger.FetchTriggers(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Trigger.CountTriggers(ctx)
	if err != nil {
		return nil, err
	}
	return &api.List[model.Trigger]{Items: triggers, Total: total, Page: page, Limit: limit}, nil
}
