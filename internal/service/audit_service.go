package service

import (
	"context"
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
)

type AuditService interface {
	// GetTableAuditLogs retrieves a paginated list of audit logs for a table
	GetTableAuditLogs(ctx context.Context, tableName string, page int, limit int) ([]model.Audit, error)
	// GetRecordAuditLogs retrieves a paginated list of audit logs for a record
	GetRecordAuditLogs(ctx context.Context, tableName, recordId string, page int, limit int) ([]model.Audit, error)
	// GetActorAuditLogs retrieves a paginated list of audit logs for an actor
	GetActorAuditLogs(ctx context.Context, actor, actorId string, page int, limit int) ([]model.Audit, error)
}

type auditService struct {
	repos *repository.Repos
}

func NewAuditService(repos *repository.Repos) AuditService {
	return &auditService{repos: repos}
}
