package repository

import "digital-wallet/internal/model"

type ExchangeRequest struct {
	WalletID       string
	Transaction    *model.Transaction
	AccountVersion uint64
}
