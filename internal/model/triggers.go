package model

import (
	"github.com/abdelrahman146/digital-wallet/pkg/types"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Trigger struct {
	Auditable
	ID         uint64      `json:"id" gorm:"column:id;primaryKey"`
	Name       string      `json:"name" gorm:"column:name"`
	Slug       string      `json:"slug" gorm:"column:slug"`
	Properties types.JSONB `json:"properties" gorm:"column:properties"`
	CreatedAt  time.Time   `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt  time.Time   `json:"updatedAt" gorm:"column:updated_at"`
}

func (m *Trigger) TableName() string {
	return "triggers"
}

func (m *Trigger) AfterCreate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(m.TableName(), AuditOperationCreate, strconv.FormatUint(m.ID, 10), m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *Trigger) AfterUpdate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(m.TableName(), AuditOperationUpdate, strconv.FormatUint(m.ID, 10), m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *Trigger) AfterDelete(tx *gorm.DB) error {
	audit, err := m.CreateAudit(m.TableName(), AuditOperationDelete, strconv.FormatUint(m.ID, 10), nil)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}
