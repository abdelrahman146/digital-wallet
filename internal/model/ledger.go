package model

import (
	"time"
)

var (
	LedgerTypeDebit  = "DEBIT"
	LedgerTypeCredit = "CREDIT"
)

var (
	LedgerActorTypeSystem     = "SYSTEM"
	LedgerActorTypeBackoffice = "BACKOFFICE"
	LedgerActorTypeUser       = "USER"
)

type Transaction struct {
	ID              string    `gorm:"column:id;primaryKey" json:"id"`
	Type            string    `gorm:"column:type" json:"type"`
	AccountID       string    `gorm:"column:account_id" json:"accountId"`
	ActorType       string    `gorm:"column:actor_type" json:"actorType"`
	ActorID         *string   `gorm:"column:actor_id" json:"actorId"`
	Amount          uint64    `gorm:"column:amount" json:"amount"`
	Metadata        string    `gorm:"column:metadata;type:jsonb" json:"metadata"`
	PreviousBalance uint64    `gorm:"column:previous_balance" json:"previousBalance"`
	NewBalance      uint64    `gorm:"column:new_balance" json:"newBalance"`
	Version         uint64    `gorm:"column:version" json:"version"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (Transaction) TableName() string {
	return "transactions"
}
