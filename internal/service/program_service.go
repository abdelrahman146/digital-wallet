package service

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/internal/repository"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/errs"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
)

type ProgramService interface {
	CreateProgram(ctx context.Context, req CreateProgramRequest) (*model.Program, error)
	UpdateProgram(ctx context.Context, id uint64, req UpdateProgramRequest) (*model.Program, error)
	DeleteProgram(ctx context.Context, id uint64) error
	GetProgram(ctx context.Context, id uint64) (*model.Program, error)
	ListPrograms(ctx context.Context, page, limit int) (*api.List[model.Program], error)
}

type programService struct {
	repos *repository.Repos
}

func NewProgramService(repos *repository.Repos) ProgramService {
	return &programService{repos: repos}
}

func (s *programService) CreateProgram(ctx context.Context, req CreateProgramRequest) (*model.Program, error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("Unauthorized access", logger.Field("error", err))
		return nil, err
	}
	program := &model.Program{
		Name:         req.Name,
		WalletID:     req.WalletID,
		TriggerSlug:  req.TriggerSlug,
		Condition:    req.Condition,
		Effect:       req.Effect,
		ValidFrom:    req.ValidFrom,
		ValidUntil:   req.ValidUntil,
		IsActive:     req.IsActive,
		LimitPerUser: req.LimitPerUser,
	}
	program.SetActor(api.GetActor(ctx), api.GetActorID(ctx))
	program.SetRemarks("Program created")
	if err := s.repos.Program.CreateProgram(ctx, program); err != nil {
		return nil, err
	}
	return program, nil
}

func (s *programService) UpdateProgram(ctx context.Context, id uint64, req UpdateProgramRequest) (*model.Program, error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("Unauthorized access", logger.Field("error", err))
		return nil, err
	}
	program, err := s.repos.Program.FetchProgramByID(ctx, id)
	if program == nil {
		return nil, errs.NewNotFoundError("Program not found", "PROGRAM_NOT_FOUND", err)
	}
	program.SetOldRecord(*program)
	if req.Name != nil {
		program.Name = *req.Name
	}
	if req.WalletID != nil {
		program.WalletID = *req.WalletID
	}
	if req.TriggerSlug != nil {
		program.TriggerSlug = *req.TriggerSlug
	}
	if req.Condition != nil {
		program.Condition = *req.Condition
	}
	if req.Effect != nil {
		program.Effect = *req.Effect
	}
	if req.ValidFrom != nil {
		program.ValidFrom = *req.ValidFrom
	}
	if req.ValidUntil != nil {
		program.ValidUntil = req.ValidUntil
	}
	if req.IsActive != nil {
		program.IsActive = *req.IsActive
	}
	if req.LimitPerUser != nil {
		program.LimitPerUser = req.LimitPerUser
	}
	program.SetActor(api.GetActor(ctx), api.GetActorID(ctx))
	program.SetRemarks("Program updated")
	if err := s.repos.Program.UpdateProgram(ctx, program); err != nil {
		return nil, err
	}
	return program, nil
}

func (s *programService) DeleteProgram(ctx context.Context, id uint64) error {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("Unauthorized access", logger.Field("error", err))
		return err
	}
	program, err := s.repos.Program.FetchProgramByID(ctx, id)
	if program == nil {
		return errs.NewNotFoundError("Program not found", "PROGRAM_NOT_FOUND", err)
	}
	if err := s.repos.Program.DeleteProgram(ctx, program); err != nil {
		return err
	}
	return nil
}

func (s *programService) GetProgram(ctx context.Context, id uint64) (*model.Program, error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("Unauthorized access", logger.Field("error", err))
		return nil, err
	}
	program, err := s.repos.Program.FetchProgramByID(ctx, id)
	if program == nil {
		return nil, errs.NewNotFoundError("Program not found", "PROGRAM_NOT_FOUND", err)
	}
	return program, nil
}

func (s *programService) ListPrograms(ctx context.Context, page, limit int) (*api.List[model.Program], error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("Unauthorized access", logger.Field("error", err))
		return nil, err
	}
	programs, err := s.repos.Program.FetchPrograms(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	count, err := s.repos.Program.CountPrograms(ctx)
	if err != nil {
		return nil, err
	}
	return &api.List[model.Program]{Items: programs, Limit: limit, Page: page, Total: count}, nil
}
