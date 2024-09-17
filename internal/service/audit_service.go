package service

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/internal/repository"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
)

type AuditService interface {
	// GetTableAuditLogs retrieves a paginated list of audit logs for a table
	GetTableAuditLogs(ctx context.Context, tableName string, page int, limit int) (*api.List[model.Audit], error)
	// GetRecordAuditLogs retrieves a paginated list of audit logs for a record
	GetRecordAuditLogs(ctx context.Context, tableName, recordId string, page int, limit int) (*api.List[model.Audit], error)
	// GetActorAuditLogs retrieves a paginated list of audit logs for an actor
	GetActorAuditLogs(ctx context.Context, actor, actorId string, page int, limit int) (*api.List[model.Audit], error)
}

type auditService struct {
	repos *repository.Repos
}

func NewAuditService(repos *repository.Repos) AuditService {
	return &auditService{repos: repos}
}

func (s *auditService) GetTableAuditLogs(ctx context.Context, tableName string, page int, limit int) (*api.List[model.Audit], error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return nil, err
	}
	audits, err := s.repos.Audit.FetchTableAuditLogs(ctx, tableName, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Audit.CountTableAuditLogs(ctx, tableName)
	if err != nil {
		return nil, err
	}
	return &api.List[model.Audit]{Items: audits, Page: page, Limit: limit, Total: total}, nil
}

func (s *auditService) GetRecordAuditLogs(ctx context.Context, tableName, recordId string, page int, limit int) (*api.List[model.Audit], error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return nil, err
	}
	audits, err := s.repos.Audit.FetchRecordAuditLogs(ctx, tableName, recordId, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Audit.CountRecordAuditLogs(ctx, tableName, recordId)
	if err != nil {
		return nil, err
	}
	return &api.List[model.Audit]{Items: audits, Page: page, Limit: limit, Total: total}, nil
}

func (s *auditService) GetActorAuditLogs(ctx context.Context, actor, actorId string, page int, limit int) (*api.List[model.Audit], error) {
	if err := api.IsAdmin(ctx); err != nil {
		api.GetLogger(ctx).Error("User not authorized")
		return nil, err
	}
	audits, err := s.repos.Audit.FetchActorAuditLogs(ctx, actor, actorId, page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Audit.CountActorAuditLogs(ctx, actor, actorId)
	if err != nil {
		return nil, err
	}
	return &api.List[model.Audit]{Items: audits, Page: page, Limit: limit, Total: total}, nil
}
