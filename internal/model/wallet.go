package model

import "time"

type Wallet struct {
	ID                string         `gorm:"column:id;primaryKey" json:"id"`
	Name              string         `gorm:"column:name" json:"name"`
	Description       string         `gorm:"column:description" json:"description"`
	Currency          string         `gorm:"column:currency" json:"currency"`
	IsMonetary        bool           `gorm:"column:is_monetary" json:"isMonetary"`
	PointsExpireAfter *time.Duration `gorm:"column:points_expire_after" json:"pointsExpireAfter"`
	LimitPerUser      *uint64        `gorm:"column:limit_per_user" json:"limitPerUser"`
	LimitGlobal       *uint64        `gorm:"column:limit_global" json:"limitGlobal"`
	CreatedAt         time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt         time.Time      `gorm:"column:updated_at" json:"updatedAt"`
}

func (Wallet) TableName() string {
	return "wallets"
}
