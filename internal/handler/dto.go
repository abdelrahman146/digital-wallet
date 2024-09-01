package handler

type UserDepositRequest struct {
	UserID            string  `json:"userId,omitempty" validate:"required"`
	Amount            float64 `json:"amount,omitempty" validate:"required,numeric,gt=0"`
	InitiatedBy       string  `json:"initiatedBy,omitempty" validate:"required,oneof=SYSTEM USER BACKOFFICE"`
	BankTransactionId string  `json:"bankTransactionId,omitempty" validate:"required"`
}
