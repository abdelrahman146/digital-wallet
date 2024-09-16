package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Auditable
	ID        string    `json:"id" gorm:"column:id;primaryKey"`
	TierID    *string   `json:"tierId" gorm:"column:tier_id"`
	IsActive  bool      `json:"isActive" gorm:"column:is_active"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
	Accounts  []Account `json:"accounts,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

func (m *User) TableName() string {
	return "users"
}

func (m *User) AfterCreate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(AuditOperationCreate, m.ID, m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *User) AfterUpdate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(AuditOperationUpdate, m.ID, m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *User) AfterDelete(tx *gorm.DB) error {
	audit, err := m.CreateAudit(AuditOperationDelete, m.ID, nil)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}
