package service

type Services struct {
	Transaction TransactionService
	Account     AccountService
	Wallet      WalletService
	Tier        TierService
	User        UserService
}
