package repository

import (
	"context"
	"digital-wallet/internal/model"
)

type AuditRepo interface {
	Log(ctx context.Context, audit *model.Audit) error
}
