package service

import (
	"github.com/abdelrahman146/digital-wallet/internal/model"
	rule_engine "github.com/abdelrahman146/digital-wallet/pkg/rules_engine"
	"github.com/abdelrahman146/digital-wallet/pkg/types"
	"github.com/shopspring/decimal"
	"time"
)

type CreateProgramRequest struct {
	Name         string           `json:"name,omitempty" validate:"required,min=1,max=100"`
	WalletID     string           `json:"walletId,omitempty" validate:"required"`
	TriggerSlug  string           `json:"triggerSlug,omitempty" validate:"required"`
	Condition    rule_engine.Rule `json:"condition,omitempty" validate:"required,json"`
	Effect       types.JSONB      `json:"effect,omitempty" validate:"required,json"`
	ValidFrom    time.Time        `json:"validFrom,omitempty" validate:"required"`
	ValidUntil   *time.Time       `json:"validUntil,omitempty"`
	IsActive     bool             `json:"isActive,omitempty"`
	LimitPerUser *uint64          `json:"limitPerUser,omitempty"`
}

type UpdateProgramRequest struct {
	Name         *string           `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	WalletID     *string           `json:"walletId,omitempty" validate:"omitempty"`
	TriggerSlug  *string           `json:"triggerSlug,omitempty" validate:"omitempty"`
	Condition    *rule_engine.Rule `json:"condition,omitempty" validate:"omitempty,json"`
	Effect       *types.JSONB      `json:"effect,omitempty" validate:"omitempty,json"`
	ValidFrom    *time.Time        `json:"validFrom,omitempty" validate:"omitempty"`
	ValidUntil   *time.Time        `json:"validUntil,omitempty" validate:"omitempty"`
	IsActive     *bool             `json:"isActive,omitempty"`
	LimitPerUser *uint64           `json:"limitPerUser,omitempty"`
}

type CreateTriggerRequest struct {
	Name       string                 `json:"name,omitempty" validate:"required,min=1,max=100"`
	Slug       string                 `json:"slug,omitempty" validate:"required,slug"`
	Properties map[string]interface{} `json:"properties,omitempty" validate:"required,json"`
}

type UpdateTriggerRequest struct {
	Name       *string                 `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Slug       *string                 `json:"slug,omitempty" validate:"omitempty,slug"`
	Properties *map[string]interface{} `json:"properties,omitempty" validate:"omitempty,json"`
}

type ExchangeResponse struct {
	FromTransaction model.Transaction `json:"fromTransaction"`
	ToTransaction   model.Transaction `json:"toTransaction"`
}

type CreateExchangeRateRequest struct {
	FromWalletID string          `json:"fromWalletId,omitempty" validate:"required"`
	ToWalletID   string          `json:"toWalletId,omitempty" validate:"required"`
	TierID       *string         `json:"tierId,omitempty"`
	ExchangeRate decimal.Decimal `json:"exchangeRate,omitempty" validate:"required"`
}

type CreateUserRequest struct {
	ID     string  `json:"id,omitempty" validate:"required,min=1,max=20"`
	TierID *string `json:"tierId,omitempty"`
}

type CreateTierRequest struct {
	ID   string `json:"id,omitempty" validate:"required,min=1,max=20"`
	Name string `json:"name,omitempty" validate:"required,min=1,max=100"`
}

type CreateWalletRequest struct {
	ID                string  `json:"id,omitempty" validate:"required"`
	Name              string  `json:"name,omitempty" validate:"required,min=1,max=100"`
	Description       *string `json:"description,omitempty" validate:"max=255"`
	Currency          string  `json:"currency,omitempty" validate:"required,min=1,max=4"`
	IsMonetary        bool    `json:"isMonetary,omitempty"`
	PointsExpireAfter *int    `json:"pointsExpireAfter,omitempty"`
	LimitPerUser      *uint64 `json:"limitPerUser,omitempty"`
	LimitGlobal       *uint64 `json:"limitGlobal,omitempty"`
}

type UpdateWalletRequest struct {
	Name              string  `json:"name,omitempty" validate:"required,min=1,max=100"`
	Description       *string `json:"description,omitempty" validate:"max=255"`
	Currency          string  `json:"currency,omitempty" validate:"required"`
	IsMonetary        *bool   `json:"isMonetary,omitempty"`
	PointsExpireAfter *int64  `json:"pointsExpireAfter,omitempty"`
	LimitPerUser      *uint64 `json:"limitPerUser,omitempty"`
	LimitGlobal       *uint64 `json:"limitGlobal,omitempty"`
}

type TransactionRequest struct {
	Type      string      `json:"type,omitempty" validate:"required,oneof=CREDIT DEBIT"`
	Amount    uint64      `json:"amount,omitempty" validate:"required,gt=0"`
	Reason    string      `json:"reason,omitempty" validate:"required,oneof=REWARD PURCHASE REDEEM PENALTY EXPIRED WITHDRAWAL DEPOSIT"`
	Metadata  types.JSONB `json:"metadata,omitempty"`
	ProgramID *string     `json:"programId,omitempty" validate:"omitempty"`
}
