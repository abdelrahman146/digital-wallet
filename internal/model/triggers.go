package model

import (
	"digital-wallet/pkg/types"
	"time"
)

type Trigger struct {
	ID         uint64      `json:"id" gorm:"column:id;primaryKey"`
	Name       string      `json:"name" gorm:"column:name"`
	Slug       string      `json:"slug" gorm:"column:slug"`
	Properties types.JSONB `json:"properties" gorm:"column:properties"`
	CreatedAt  time.Time   `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt  time.Time   `json:"updatedAt" gorm:"column:updated_at"`
}

func (Trigger) TableName() string {
	return "triggers"
}
