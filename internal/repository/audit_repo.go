package repository

import (
	"context"
	"github.com/abdelrahman146/digital-wallet/internal/model"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
	"gorm.io/gorm"
)

type AuditRepo interface {
	// FetchTableAuditLogs retrieves a paginated list of audit logs for a table
	FetchTableAuditLogs(ctx context.Context, tableName string, page int, limit int) ([]model.Audit, error)
	// CountTableAuditLogs retrieves the total number of audit logs for a table
	CountTableAuditLogs(ctx context.Context, tableName string) (int64, error)
	// FetchRecordAuditLogs retrieves a paginated list of audit logs for a record
	FetchRecordAuditLogs(ctx context.Context, tableName, recordId string, page int, limit int) ([]model.Audit, error)
	// CountRecordAuditLogs retrieves the total number of audit logs for a record
	CountRecordAuditLogs(ctx context.Context, tableName, recordId string) (int64, error)
	// FetchActorAuditLogs retrieves a paginated list of audit logs for an actor
	FetchActorAuditLogs(ctx context.Context, actor, actorId string, page int, limit int) ([]model.Audit, error)
	// CountActorAuditLogs retrieves the total number of audit logs for an actor
	CountActorAuditLogs(ctx context.Context, actor, actorId string) (int64, error)
}

type auditRepo struct {
	db *gorm.DB
}

func NewAuditRepo(db *gorm.DB) AuditRepo {
	return &auditRepo{db: db}
}

// FetchTableAuditLogs retrieves a paginated list of audit logs for a table
func (r *auditRepo) FetchTableAuditLogs(ctx context.Context, tableName string, page int, limit int) ([]model.Audit, error) {
	var audits []model.Audit
	err := r.db.Where("table_name = ?", tableName).Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&audits).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to fetch table audit logs", logger.Field("error", err), logger.Field("tableName", tableName))
		return nil, err
	}
	return audits, nil
}

// CountTableAuditLogs retrieves the total number of audit logs for a table
func (r *auditRepo) CountTableAuditLogs(ctx context.Context, tableName string) (int64, error) {
	var total int64
	err := r.db.Model(&model.Audit{}).Where("table_name = ?", tableName).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to count table audit logs", logger.Field("error", err), logger.Field("tableName", tableName))
		return 0, err
	}
	return total, nil
}

// FetchRecordAuditLogs retrieves a paginated list of audit logs for a record
func (r *auditRepo) FetchRecordAuditLogs(ctx context.Context, tableName, recordId string, page int, limit int) ([]model.Audit, error) {
	var audits []model.Audit
	err := r.db.Where("table_name = ? AND record_id = ?", tableName, recordId).Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&audits).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to fetch record audit logs", logger.Field("error", err), logger.Field("tableName", tableName), logger.Field("recordId", recordId))
		return nil, err
	}
	return audits, nil
}

// CountRecordAuditLogs retrieves the total number of audit logs for a record
func (r *auditRepo) CountRecordAuditLogs(ctx context.Context, tableName, recordId string) (int64, error) {
	var total int64
	err := r.db.Model(&model.Audit{}).Where("table_name = ? AND record_id = ?", tableName, recordId).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to count record audit logs", logger.Field("error", err), logger.Field("tableName", tableName), logger.Field("recordId", recordId))
		return 0, err
	}
	return total, nil
}

// FetchActorAuditLogs retrieves a paginated list of audit logs for an actor
func (r *auditRepo) FetchActorAuditLogs(ctx context.Context, actor, actorId string, page int, limit int) ([]model.Audit, error) {
	var audits []model.Audit
	err := r.db.Where("actor = ? AND actor_id = ?", actor, actorId).Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&audits).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to fetch actor audit logs", logger.Field("error", err), logger.Field("actorType", actor), logger.Field("actorId", actorId))
		return nil, err
	}
	return audits, nil
}

// CountActorAuditLogs retrieves the total number of audit logs for an actor
func (r *auditRepo) CountActorAuditLogs(ctx context.Context, actor, actorId string) (int64, error) {
	var total int64
	err := r.db.Model(&model.Audit{}).Where("actor = ? AND actor_id = ?", actor, actorId).Count(&total).Error
	if err != nil {
		api.GetLogger(ctx).Error("Failed to count actor audit logs", logger.Field("error", err), logger.Field("actorType", actor), logger.Field("actorId", actorId))
		return 0, err
	}
	return total, nil
}
