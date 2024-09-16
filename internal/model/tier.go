package model

import (
	"gorm.io/gorm"
	"time"
)

type Tier struct {
	Auditable
	ID          string    `gorm:"column:id;primaryKey;" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	Description *string   `gorm:"column:description" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (m *Tier) TableName() string {
	return "tiers"
}

func (m *Tier) AfterCreate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(AuditOperationCreate, m.ID, m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *Tier) AfterUpdate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(AuditOperationUpdate, m.ID, m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *Tier) AfterDelete(tx *gorm.DB) error {
	audit, err := m.CreateAudit(AuditOperationDelete, m.ID, nil)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}
