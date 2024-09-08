package service

type CreateUserRequest struct {
	ID     string `json:"id,omitempty" validate:"required,min=1,max=20"`
	TierID string `json:"tierId,omitempty"`
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
	PointsExpireAfter *int    `json:"pointsExpireAfter,omitempty"`
	LimitPerUser      *uint64 `json:"limitPerUser,omitempty"`
	LimitGlobal       *uint64 `json:"limitGlobal,omitempty"`
}

type UpdateWalletRequest struct {
	Name              string  `json:"name,omitempty" validate:"required,min=1,max=100"`
	Description       string  `json:"description,omitempty" validate:"required,min=1,max=255"`
	Currency          string  `json:"currency,omitempty" validate:"required,min=1,max=4"`
	PointsExpireAfter *int64  `json:"pointsExpireAfter,omitempty"`
	LimitPerUser      *uint64 `json:"limitPerUser,omitempty"`
	LimitGlobal       *uint64 `json:"limitGlobal,omitempty"`
}

type DepositRequest struct {
	UserID               string  `json:"userId,omitempty" validate:"required"`
	Amount               float64 `json:"amount,omitempty" validate:"required,decimal2,gt=0"`
	PaymentTransactionId string  `json:"paymentTransactionId,omitempty" validate:"required"`
}

type WithdrawRequest struct {
	UserID               string  `json:"userId,omitempty" validate:"required"`
	Amount               float64 `json:"amount,omitempty" validate:"required,decimal2,lt=0"`
	PaymentTransactionId string  `json:"paymentTransactionId,omitempty" validate:"required"`
}

type RefundRequest struct {
	UserID  string  `json:"userId,omitempty" validate:"required"`
	Amount  float64 `json:"amount,omitempty" validate:"required,decimal2,gt=0"`
	OrderId string  `json:"orderId,omitempty" validate:"required"`
}

type PurchaseRequest struct {
	UserID  string  `json:"userId,omitempty" validate:"required"`
	Amount  float64 `json:"amount,omitempty" validate:"required,decimal2,lt=0"`
	OrderId string  `json:"orderId,omitempty" validate:"required"`
}

type TransferRequest struct {
	FromUserID string  `json:"fromUserId,omitempty" validate:"required"`
	ToUserID   string  `json:"toUserId,omitempty" validate:"required"`
	Amount     float64 `json:"amount,omitempty" validate:"required,decimal2,gt=0"`
}
