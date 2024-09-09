package service

import (
	"digital-wallet/internal/model"
	"digital-wallet/pkg/types"
	"github.com/shopspring/decimal"
)

type ExchangeResponse struct {
	FromTransaction model.Transaction `json:"fromTransaction"`
	ToTransaction   model.Transaction `json:"toTransaction"`
}

type CreateExchangeRateRequest struct {
	FromWalletID string          `json:"fromWalletId,omitempty" validate:"required"`
	ToWalletID   string          `json:"toWalletId,omitempty" validate:"required"`
	TierID       string          `json:"tierId,omitempty" validate:"required"`
	ExchangeRate decimal.Decimal `json:"exchangeRate,omitempty" validate:"required"`
}

type CreateUserRequest struct {
	ID     string `json:"id,omitempty" validate:"required,min=1,max=20"`
	TierID string `json:"tierId,omitempty" validate:"required"`
}

type CreateTierRequest struct {
	ID   string `json:"id,omitempty" validate:"required,min=1,max=20"`
	Name string `json:"name,omitempty" validate:"required,min=1,max=100"`
}

type CreateWalletRequest struct {
	ID                string  `json:"id,omitempty" validate:"required,min=1,max=4"`
	Name              string  `json:"name,omitempty" validate:"required,min=1,max=100"`
	Description       string  `json:"description,omitempty" validate:"required,min=1,max=255"`
	Currency          string  `json:"currency,omitempty" validate:"required,min=1,max=4"`
	IsMonetary        bool    `json:"isMonetary,omitempty"`
	PointsExpireAfter *int    `json:"pointsExpireAfter,omitempty"`
	LimitPerUser      *uint64 `json:"limitPerUser,omitempty"`
	LimitGlobal       *uint64 `json:"limitGlobal,omitempty"`
}

type UpdateWalletRequest struct {
	Name              string  `json:"name,omitempty" validate:"required,min=1,max=100"`
	Description       string  `json:"description,omitempty" validate:"required,min=1,max=255"`
	Currency          string  `json:"currency,omitempty" validate:"required,min=1,max=4"`
	IsMonetary        *bool   `json:"isMonetary,omitempty"`
	PointsExpireAfter *int64  `json:"pointsExpireAfter,omitempty"`
	LimitPerUser      *uint64 `json:"limitPerUser,omitempty"`
	LimitGlobal       *uint64 `json:"limitGlobal,omitempty"`
}

type TransactionRequest struct {
	Type      string      `json:"type,omitempty" validate:"required,oneof=credit debit"`
	Amount    uint64      `json:"amount,omitempty" validate:"required,gt=0"`
	Metadata  types.JSONB `json:"metadata,omitempty"`
	ProgramID *string     `json:"programId,omitempty" validate:"omitempty"`
}
