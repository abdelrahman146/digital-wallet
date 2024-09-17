package repository

type Repos struct {
	Audit        AuditRepo
	Transaction  TransactionRepo
	Account      AccountRepo
	Wallet       WalletRepo
	Tier         TierRepo
	User         UserRepo
	ExchangeRate ExchangeRateRepo
}
