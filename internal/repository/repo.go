package repository

type Repos struct {
	Transaction TransactionRepo
	Account     AccountRepo
	Wallet      WalletRepo
	Tier        TierRepo
	User        UserRepo
}
