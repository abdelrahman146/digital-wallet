package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type Wallet struct {
	ID        string          `gorm:"column:id;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID    string          `gorm:"column:user_id" json:"userId"`
	Balance   decimal.Decimal `gorm:"column:balance;numeric(18,2)" json:"balance"`
	Version   int64           `gorm:"column:version" json:"version"`
	CreatedAt time.Time       `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time       `gorm:"column:updated_at" json:"updatedAt"`
}

func (Wallet) TableName() string {
	return "wallets"
}
