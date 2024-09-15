package model

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	Auditable
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	WalletID  string    `gorm:"column:wallet_id" json:"walletId"`
	UserID    string    `gorm:"column:user_id" json:"userId"`
	Balance   uint64    `gorm:"column:balance;" json:"balance"`
	Version   uint64    `gorm:"column:version" json:"version"`
	IsActive  bool      `gorm:"column:is_active" json:"isActive"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (m *Account) TableName() string {
	return "accounts"
}

func (m *Account) AfterCreate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(AuditOperationCreate, m.ID, m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *Account) AfterUpdate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(AuditOperationUpdate, m.ID, m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *Account) AfterDelete(tx *gorm.DB) error {
	audit, err := m.CreateAudit(AuditOperationDelete, m.ID, m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}
