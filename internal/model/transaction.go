package model

import (
	"digital-wallet/pkg/types"
	"time"
)

var (
	TransactionTypeDebit  = "DEBIT"
	TransactionTypeCredit = "CREDIT"
)

var (
	TransactionActorTypeSystem     = "SYSTEM"
	TransactionActorTypeBackoffice = "BACKOFFICE"
	TransactionActorTypeUser       = "USER"
)

type Transaction struct {
	ID              string      `gorm:"column:id;primaryKey" json:"id"`
	Type            string      `gorm:"column:type" json:"type"`
	AccountID       string      `gorm:"column:account_id" json:"accountId"`
	ActorType       string      `gorm:"column:actor_type" json:"actorType"`
	ActorID         string      `gorm:"column:actor_id" json:"actorId"`
	Metadata        types.JSONB `gorm:"column:metadata;type:jsonb" json:"metadata"`
	ProgramID       *string     `gorm:"column:program_id" json:"programId"`
	Amount          uint64      `gorm:"column:amount" json:"amount"`
	PreviousBalance uint64      `gorm:"column:previous_balance" json:"previousBalance"`
	NewBalance      uint64      `gorm:"column:new_balance" json:"newBalance"`
	Version         uint64      `gorm:"column:version" json:"version"`
	CreatedAt       time.Time   `gorm:"column:created_at" json:"createdAt"`
}

func (Transaction) TableName() string {
	return "transactions"
}
