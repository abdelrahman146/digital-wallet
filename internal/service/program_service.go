package service

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
)

type ProgramService interface {
	GetProgram(ctx context.Context, id uint64) (*model.Program, error)
}
