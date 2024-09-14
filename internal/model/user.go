package model

import "time"

type User struct {
	ID        string    `json:"id" gorm:"column:id;primaryKey"`
	TierID    *string   `json:"tierId" gorm:"column:tier_id"`
	IsActive  bool      `json:"isActive" gorm:"column:is_active"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
	Accounts  []Account `json:"accounts,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

func (User) TableName() string {
	return "users"
}
