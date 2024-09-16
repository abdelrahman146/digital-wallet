package model

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type ExchangeRate struct {
	Auditable
	ID            uint64          `gorm:"column:id;primary_key" json:"id"`
	FromWalletID  string          `gorm:"column:from_wallet_id" json:"fromWalletId"`
	ToWalletID    string          `gorm:"column:to_wallet_id" json:"toWalletId"`
	TierID        *string         `gorm:"column:tier_id" json:"tierId"`
	ExchangeRate  decimal.Decimal `gorm:"column:exchange_rate" json:"exchangeRate"`
	MinimumAmount *uint64         `gorm:"column:minimum_amount" json:"minimumAmount"`
	CreatedAt     time.Time       `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt     time.Time       `gorm:"column:updated_at" json:"updatedAt"`
}

func (m *ExchangeRate) TableName() string {
	return "exchange_rates"
}

func (m *ExchangeRate) AfterCreate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(AuditOperationCreate, strconv.FormatUint(m.ID, 10), m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *ExchangeRate) AfterUpdate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(AuditOperationUpdate, strconv.FormatUint(m.ID, 10), m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *ExchangeRate) AfterDelete(tx *gorm.DB) error {
	audit, err := m.CreateAudit(AuditOperationDelete, strconv.FormatUint(m.ID, 10), nil)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}
