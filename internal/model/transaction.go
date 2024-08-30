package model

type Transaction struct {
	ID              string  `gorm:"column:id" json:"id"`
	WalletID        string  `gorm:"column:wallet_id" json:"walletId"`
	Amount          float64 `gorm:"column:amount" json:"amount"`
	Type            string  `gorm:"column:type" json:"type"`
	ReferenceID     *string `gorm:"column:reference_id" json:"referenceId"`
	ReferenceType   *string `gorm:"column:reference_type" json:"referenceType"`
	InitiatedBy     string  `gorm:"column:initiated_by" json:"initiatedBy"`
	PreviousBalance float64 `gorm:"column:previous_balance" json:"previousBalance"`
	NewBalance      float64 `gorm:"column:new_balance" json:"newBalance"`
	Version         int64   `gorm:"column:version" json:"version"`
	CreatedAt       string  `gorm:"column:created_at" json:"createdAt"`
}

func (Transaction) TableName() string {
	return "transactions"
}
