package repository

import "github.com/abdelrahman146/digital-wallet/internal/model"

type ExchangeRequest struct {
	WalletID       string
	Transaction    *model.Transaction
	AccountVersion uint64
}
