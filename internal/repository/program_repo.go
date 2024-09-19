package repository

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
	"gorm.io/gorm"
)

type ProgramRepo interface {
	// CreateProgram creates a new program
	CreateProgram(ctx context.Context, program *model.Program) error
	// UpdateProgram updates an existing program
	UpdateProgram(ctx context.Context, program *model.Program) error
	// DeleteProgram deletes a program
	DeleteProgram(ctx context.Context, program *model.Program) error
	// FetchProgramByID retrieves a program by its ID
	FetchProgramByID(ctx context.Context, id uint64) (*model.Program, error)
	// FetchProgramsByTriggerID retrieves programs for a specific trigger
	FetchProgramsByTriggerID(ctx context.Context, triggerID uint64) ([]*model.Program, error)
	// CountProgramsByTriggerID retrieves the total number of programs for a specific trigger
	CountProgramsByTriggerID(ctx context.Context, triggerID uint64) (int64, error)
	// FetchProgramsByWalletID retrieves programs for a specific wallet
	FetchProgramsByWalletID(ctx context.Context, walletID uint64) ([]*model.Program, error)
	// CountProgramsByWalletID retrieves the total number of programs for a specific wallet
	CountProgramsByWalletID(ctx context.Context, walletID uint64) (int64, error)
	// FetchPrograms retrieves a paginated list of programs
	FetchPrograms(ctx context.Context, page int, limit int) ([]*model.Program, error)
	// CountPrograms retrieves the total number of programs
	CountPrograms(ctx context.Context) (int64, error)
}

type programRepo struct {
	db *gorm.DB
}

func NewProgramRepo(db *gorm.DB) ProgramRepo {
	return &programRepo{db: db}
}

func (r *programRepo) CreateProgram(ctx context.Context, program *model.Program) error {
	if err := r.db.Create(program).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to create program", logger.Field("error", err), logger.Field("program", program))
		return err
	}
	return nil
}

func (r *programRepo) UpdateProgram(ctx context.Context, program *model.Program) error {
	if err := r.db.Save(program).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to update program", logger.Field("error", err), logger.Field("program", program))
		return err
	}
	return nil
}

func (r *programRepo) DeleteProgram(ctx context.Context, program *model.Program) error {
	if err := r.db.Delete(program).Error; err != nil {
		api.GetLogger(ctx).Error("Failed to delete program", logger.Field("error", err), logger.Field("program", program))
		return err
	}
	return nil
}

func (r *programRepo) FetchProgramByID(ctx context.Context, id uint64) (*model.Program, error) {
	var program model.Program
	err := r.db.Where("id = ?", id).First(&program).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve program by ID", logger.Field("error", err), logger.Field("id", id))
		return nil, err
	}
	return &program, nil
}

func (r *programRepo) FetchProgramsByTriggerID(ctx context.Context, triggerID uint64) ([]*model.Program, error) {
	var programs []*model.Program
	err := r.db.Where("trigger_id = ?", triggerID).Find(&programs).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve programs by trigger ID", logger.Field("error", err), logger.Field("triggerId", triggerID))
		return nil, err
	}
	return programs, nil
}

func (r *programRepo) CountProgramsByTriggerID(ctx context.Context, triggerID uint64) (int64, error) {
	var total int64
	err := r.db.Model(&model.Program{}).Where("trigger_id = ?", triggerID).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve total programs count by trigger ID", logger.Field("error", err), logger.Field("triggerId", triggerID))
		return 0, err
	}
	return total, nil
}

func (r *programRepo) FetchProgramsByWalletID(ctx context.Context, walletID uint64) ([]*model.Program, error) {
	var programs []*model.Program
	err := r.db.Where("wallet_id = ?", walletID).Find(&programs).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve programs by wallet ID", logger.Field("error", err), logger.Field("walletId", walletID))
		return nil, err
	}
	return programs, nil
}

func (r *programRepo) CountProgramsByWalletID(ctx context.Context, walletID uint64) (int64, error) {
	var total int64
	err := r.db.Model(&model.Program{}).Where("wallet_id = ?", walletID).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve total programs count by wallet ID", logger.Field("error", err), logger.Field("walletId", walletID))
		return 0, err
	}
	return total, nil
}

func (r *programRepo) FetchPrograms(ctx context.Context, page int, limit int) ([]*model.Program, error) {
	var programs []*model.Program
	err := r.db.Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&programs).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve programs", logger.Field("error", err))
		return nil, err
	}
	return programs, nil
}

func (r *programRepo) CountPrograms(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.Model(&model.Program{}).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to retrieve total programs count", logger.Field("error", err))
		return 0, err
	}
	return total, nil
}
