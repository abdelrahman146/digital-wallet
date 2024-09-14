package model

import (
	"time"
)

type Account struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	WalletID  string    `gorm:"column:wallet_id" json:"walletId"`
	UserID    string    `gorm:"column:user_id" json:"userId"`
	Balance   uint64    `gorm:"column:balance;" json:"balance"`
	Version   uint64    `gorm:"column:version" json:"version"`
	IsActive  bool      `gorm:"column:is_active" json:"isActive"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (Account) TableName() string {
	return "accounts"
}
