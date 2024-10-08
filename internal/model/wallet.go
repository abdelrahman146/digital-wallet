package model

import (
	"github.com/abdelrahman146/digital-wallet/pkg/types"
	"gorm.io/gorm"
	"time"
)

type Wallet struct {
	Auditable
	ID          string  `gorm:"column:id;primaryKey" json:"id"`
	Name        string  `gorm:"column:name" json:"name"`
	Description *string `gorm:"column:description" json:"description"`
	Currency    string  `gorm:"column:currency" json:"currency"`
	// @swaggertype string
	PointsExpireAfter *types.Interval `gorm:"column:points_expire_after" json:"pointsExpireAfter"`
	LimitPerUser      *uint64         `gorm:"column:limit_per_user" json:"limitPerUser"`
	LimitGlobal       *uint64         `gorm:"column:limit_global" json:"limitGlobal"`
	MinimumWithdrawal *uint64         `gorm:"column:minimum_withdrawal" json:"minimumWithdrawal"`
	IsMonetary        bool            `gorm:"column:is_monetary;default:false" json:"isMonetary"`
	IsActive          bool            `gorm:"column:is_active;default:true" json:"isActive"`
	IsArchived        bool            `gorm:"column:is_archived;default:false" json:"isArchived"`
	CreatedAt         time.Time       `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt         time.Time       `gorm:"column:updated_at" json:"updatedAt"`
}

func (m *Wallet) TableName() string {
	return "wallets"
}

func (m *Wallet) AfterCreate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(m.TableName(), AuditOperationCreate, m.ID, m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *Wallet) AfterUpdate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(m.TableName(), AuditOperationUpdate, m.ID, m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *Wallet) AfterDelete(tx *gorm.DB) error {
	audit, err := m.CreateAudit(m.TableName(), AuditOperationDelete, m.ID, nil)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}
