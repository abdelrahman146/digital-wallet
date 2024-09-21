package service

type Services struct {
	Audit        AuditService
	Transaction  TransactionService
	Account      AccountService
	Wallet       WalletService
	Tier         TierService
	User         UserService
	ExchangeRate ExchangeRateService
	Trigger      TriggerService
	Program      ProgramService
}
