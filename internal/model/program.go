package model

import (
	rule_engine "github.com/abdelrahman146/digital-wallet/pkg/rules_engine"
	"github.com/abdelrahman146/digital-wallet/pkg/types"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Program struct {
	Auditable
	ID           uint64           `json:"id" gorm:"column:id;primaryKey"`
	Name         string           `json:"name" gorm:"column:name"`
	WalletID     string           `json:"walletId" gorm:"column:wallet_id"`
	TriggerSlug  string           `json:"triggerSlug" gorm:"column:trigger_slug"`
	Condition    rule_engine.Rule `json:"condition" gorm:"column:condition"`
	Effect       types.JSONB      `json:"effect" gorm:"column:effect"`
	ValidFrom    time.Time        `json:"validFrom" gorm:"column:valid_from"`
	ValidUntil   *time.Time       `json:"validUntil" gorm:"column:valid_until"`
	IsActive     bool             `json:"isActive;default:false" gorm:"column:is_active"`
	LimitPerUser *uint64          `json:"limitPerUser" gorm:"column:limit_per_user"`
	LimitGlobal  *uint64          `json:"limitGlobal" gorm:"column:limit_global"`
	CreatedAt    time.Time        `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time        `json:"updatedAt" gorm:"column:updated_at"`
}

func (m *Program) TableName() string {
	return "programs"
}

func (m *Program) AfterCreate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(m.TableName(), AuditOperationCreate, strconv.FormatUint(m.ID, 10), m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *Program) AfterUpdate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(m.TableName(), AuditOperationUpdate, strconv.FormatUint(m.ID, 10), m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *Program) AfterDelete(tx *gorm.DB) error {
	audit, err := m.CreateAudit(m.TableName(), AuditOperationDelete, strconv.FormatUint(m.ID, 10), nil)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}
