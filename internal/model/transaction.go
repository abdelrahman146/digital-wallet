package model

import (
	"digital-wallet/pkg/types"
	"gorm.io/gorm"
	"time"
)

const (
	TransactionTypeDebit  = "DEBIT"
	TransactionTypeCredit = "CREDIT"
)

const (
	TransactionReasonDeposit    = "DEPOSIT"
	TransactionReasonWithdrawal = "WITHDRAWAL"
	TransactionReasonExchange   = "EXCHANGE"
	TransactionReasonPurchase   = "PURCHASE"
	TransactionReasonRedeem     = "REDEEM"
	TransactionReasonPenalty    = "PENALTY"
	TransactionReasonExpired    = "EXPIRED"
)

type Transaction struct {
	Auditable
	ID              string      `gorm:"column:id;primaryKey" json:"id"`
	Type            string      `gorm:"column:type" json:"type"`
	WalletID        string      `gorm:"column:wallet_id" json:"walletId"`
	AccountID       string      `gorm:"column:account_id" json:"accountId"`
	Reason          string      `gorm:"column:reason" json:"reason"`
	Metadata        types.JSONB `gorm:"column:metadata;type:jsonb" json:"metadata"`
	ProgramID       *string     `gorm:"column:program_id" json:"programId"`
	Amount          uint64      `gorm:"column:amount" json:"amount"`
	AvailableAmount uint64      `gorm:"column:available_amount" json:"availableAmount"`
	ExpireAt        *time.Time  `gorm:"column:expire_at" json:"expireAt"`
	PreviousBalance uint64      `gorm:"column:previous_balance" json:"previousBalance"`
	NewBalance      uint64      `gorm:"column:new_balance" json:"newBalance"`
	Version         uint64      `gorm:"column:version" json:"version"`
	CreatedAt       time.Time   `gorm:"column:created_at" json:"createdAt"`
}

func (m *Transaction) TableName() string {
	return "transactions"
}

func (m *Transaction) AfterCreate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(m.TableName(), AuditOperationCreate, m.ID, m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *Transaction) AfterUpdate(tx *gorm.DB) error {
	audit, err := m.CreateAudit(m.TableName(), AuditOperationUpdate, m.ID, m)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}

func (m *Transaction) AfterDelete(tx *gorm.DB) error {
	audit, err := m.CreateAudit(m.TableName(), AuditOperationDelete, m.ID, nil)
	if err != nil {
		return err
	}
	return tx.Create(audit).Error
}
