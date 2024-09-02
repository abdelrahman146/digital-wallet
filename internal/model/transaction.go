package model

import (
	"github.com/shopspring/decimal"
	"time"
)

var (
	TransactionTypeDeposit     = "DEPOSIT"
	TransactionTypeWithdraw    = "WITHDRAW"
	TransactionTypeRefund      = "REFUND"
	TransactionTypePurchase    = "PURCHASE"
	TransactionTypeTransferIn  = "TRANSFER_IN"
	TransactionTypeTransferOut = "TRANSFER_OUT"
)

var (
	TransactionReferenceTypeBankTransaction = "PAYMENT_TRANSACTION"
	TransactionReferenceTypeOrder           = "ORDER"
	TransactionReferenceTypeTransfer        = "TRANSFER"
)

var (
	TransactionInitiatedBySystem     = "SYSTEM"
	TransactionInitiatedByBackoffice = "BACKOFFICE"
	TransactionInitiatedByUser       = "USER"
)

type Transaction struct {
	ID              string          `gorm:"column:id;primaryKey;default:uuid_generate_v4()" json:"id"`
	WalletID        string          `gorm:"column:wallet_id" json:"walletId"`
	Amount          decimal.Decimal `gorm:"column:amount;numeric(18,2)" json:"amount"`
	Type            string          `gorm:"column:type" json:"type"`
	ReferenceID     *string         `gorm:"column:reference_id" json:"referenceId"`
	ReferenceType   *string         `gorm:"column:reference_type" json:"referenceType"`
	InitiatedBy     string          `gorm:"column:initiated_by" json:"initiatedBy"`
	PreviousBalance decimal.Decimal `gorm:"column:previous_balance;numeric(18,2)" json:"previousBalance"`
	NewBalance      decimal.Decimal `gorm:"column:new_balance;numeric(18,2)" json:"newBalance"`
	Version         int64           `gorm:"column:version" json:"version"`
	CreatedAt       time.Time       `gorm:"column:created_at" json:"createdAt"`
}

func (Transaction) TableName() string {
	return "transactions"
}
