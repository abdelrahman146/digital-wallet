package model

import "time"

type Wallet struct {
	ID        string    `gorm:"column:id" json:"id"`
	UserID    string    `gorm:"column:user_id" json:"userId"`
	Amount    float64   `gorm:"column:amount" json:"amount"`
	Version   int64     `gorm:"column:version" json:"version"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (Wallet) TableName() string {
	return "wallet"
}
