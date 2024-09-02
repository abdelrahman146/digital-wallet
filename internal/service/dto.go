package service

import "github.com/shopspring/decimal"

type DepositRequest struct {
	UserID               string          `json:"userId,omitempty" validate:"required"`
	Amount               decimal.Decimal `json:"amount,omitempty" validate:"required,decimal2,gt=0"`
	PaymentTransactionId string          `json:"paymentTransactionId,omitempty" validate:"required"`
}

type WithdrawRequest struct {
	UserID               string          `json:"userId,omitempty" validate:"required"`
	Amount               decimal.Decimal `json:"amount,omitempty" validate:"required,decimal2,lt=0"`
	PaymentTransactionId string          `json:"paymentTransactionId,omitempty" validate:"required"`
}

type RefundRequest struct {
	UserID  string          `json:"userId,omitempty" validate:"required"`
	Amount  decimal.Decimal `json:"amount,omitempty" validate:"required,decimal2,gt=0"`
	OrderId string          `json:"orderId,omitempty" validate:"required"`
}

type PurchaseRequest struct {
	UserID  string          `json:"userId,omitempty" validate:"required"`
	Amount  decimal.Decimal `json:"amount,omitempty" validate:"required,decimal2,lt=0"`
	OrderId string          `json:"orderId,omitempty" validate:"required"`
}

type TransferRequest struct {
	FromUserID string          `json:"fromUserId,omitempty" validate:"required"`
	ToUserID   string          `json:"toUserId,omitempty" validate:"required"`
	Amount     decimal.Decimal `json:"amount,omitempty" validate:"required,decimal2,gt=0"`
}
