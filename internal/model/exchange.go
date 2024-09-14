package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type ExchangeRate struct {
	ID            uint64          `gorm:"column:id;primary_key" json:"id"`
	FromWalletID  string          `gorm:"column:from_wallet_id" json:"fromWalletId"`
	ToWalletID    string          `gorm:"column:to_wallet_id" json:"toWalletId"`
	TierID        *string         `gorm:"column:tier_id" json:"tierId"`
	ExchangeRate  decimal.Decimal `gorm:"column:exchange_rate" json:"exchangeRate"`
	MinimumAmount *uint64         `gorm:"column:minimum_amount" json:"minimumAmount"`
	CreatedAt     time.Time       `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt     time.Time       `gorm:"column:updated_at" json:"updatedAt"`
}

func (ExchangeRate) TableName() string {
	return "exchange_rates"
}
