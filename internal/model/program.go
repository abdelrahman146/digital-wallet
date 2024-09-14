package model

import (
	"digital-wallet/pkg/types"
	"time"
)

type Program struct {
	ID           uint64      `json:"id" gorm:"column:id;primaryKey"`
	Name         string      `json:"name" gorm:"column:name"`
	WalletID     string      `json:"walletId" gorm:"column:wallet_id"`
	TriggerSlug  string      `json:"triggerSlug" gorm:"column:trigger_slug"`
	Condition    types.JSONB `json:"condition" gorm:"column:condition"`
	Effect       types.JSONB `json:"effect" gorm:"column:effect"`
	ValidFrom    time.Time   `json:"validFrom" gorm:"column:valid_from"`
	ValidUntil   *time.Time  `json:"validUntil" gorm:"column:valid_until"`
	IsActive     bool        `json:"isActive" gorm:"column:is_active"`
	LimitPerUser *uint64     `json:"limitPerUser" gorm:"column:limit_per_user"`
	LimitGlobal  *uint64     `json:"limitGlobal" gorm:"column:limit_global"`
	CreatedAt    time.Time   `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time   `json:"updatedAt" gorm:"column:updated_at"`
}

func (Program) TableName() string {
	return "programs"
}
