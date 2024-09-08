package model

import "time"

type Tier struct {
	ID        string    `gorm:"column:id;primaryKey;" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (Tier) TableName() string {
	return "tiers"
}
