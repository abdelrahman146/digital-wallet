package model

import (
	"digital-wallet/pkg/types"
	"time"
)

const (
	TransactionTypeDebit  = "DEBIT"
	TransactionTypeCredit = "CREDIT"
)

type Transaction struct {
	ID              string      `gorm:"column:id;primaryKey" json:"id"`
	Type            string      `gorm:"column:type" json:"type"`
	WalletID        string      `gorm:"column:wallet_id" json:"walletId"`
	AccountID       string      `gorm:"column:account_id" json:"accountId"`
	Metadata        types.JSONB `gorm:"column:metadata;type:jsonb" json:"metadata"`
	ProgramID       *string     `gorm:"column:program_id" json:"programId"`
	Amount          uint64      `gorm:"column:amount" json:"amount"`
	ExpireAt        *time.Time  `gorm:"column:expire_at" json:"expireAt"`
	PreviousBalance uint64      `gorm:"column:previous_balance" json:"previousBalance"`
	NewBalance      uint64      `gorm:"column:new_balance" json:"newBalance"`
	Version         uint64      `gorm:"column:version" json:"version"`
	CreatedAt       time.Time   `gorm:"column:created_at" json:"createdAt"`
}

func (Transaction) TableName() string {
	return "transactions"
}
