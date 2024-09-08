package model

import (
	"time"
)

type Account struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	UserID    string    `gorm:"column:user_id" json:"userId"`
	User      User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Balance   uint64    `gorm:"column:balance;" json:"balance"`
	TierID    string    `gorm:"column:tier_id" json:"tierId"`
	Version   uint64    `gorm:"column:version" json:"version"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (Account) TableName() string {
	return "accounts"
}
