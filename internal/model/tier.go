package model

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Tier struct {
	Auditable
	ID          uint64    `gorm:"column:id;primaryKey;" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	Description *string   `gorm:"column:description" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (m *Tier) TableName() string {
	return "tiers"
}

func (m *Tier) AfterCreate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(AuditOperationCreate, strconv.FormatUint(m.ID, 10), m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *Tier) AfterUpdate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(AuditOperationUpdate, strconv.FormatUint(m.ID, 10), m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *Tier) AfterDelete(tx *gorm.DB) error {
	audit, err := m.CreateAudit(AuditOperationDelete, strconv.FormatUint(m.ID, 10), m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}
