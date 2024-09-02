package main

import (
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/resources"
	"fmt"
	"github.com/shopspring/decimal"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	startDateString         = "2024-01-01"
	endDateString           = "2024-08-30"
	minWalletsPerDay        = 100
	maxWalletsPerDay        = 200
	minTransactionsPerDay   = 200
	maxTransactionsPerDay   = 1000
	minTransactionAmount    = 10.0
	maxTransactionAmount    = 1000.0
	paymentTransactionIds   = 1000000
	orderIds                = 1000000
	maxConcurrentOperations = 50
	maxRetries              = 5
	retryDelay              = 100 * time.Millisecond
)

var walletLocks sync.Map

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Parse start and end dates
	startDate, _ := time.Parse("2006-01-02", startDateString)
	endDate, _ := time.Parse("2006-01-02", endDateString)

	// Initialize database connection
	db := resources.InitDB()

	// Initialize repositories
	walletRepo := repository.NewWalletRepo(db)
	transactionRepo := repository.NewTransactionRepo(db)

	// Seed data for each day between startDate and endDate
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		log.Printf("Seeding data for date: %s", currentDate.Format("2006-01-02"))

		// Seed wallets
		numWallets := getRandomInt(minWalletsPerDay, maxWalletsPerDay)
		wallets, err := generateWallets(walletRepo, numWallets, currentDate)
		if err != nil {
			log.Printf("Failed to generate wallets for date %s: %v", currentDate.Format("2006-01-02"), err)
			continue
		}

		// Seed transactions
		numTransactions := getRandomInt(minTransactionsPerDay, maxTransactionsPerDay)
		err = generateTransactions(transactionRepo, wallets, numTransactions, currentDate)
		if err != nil {
			log.Printf("Failed to generate transactions for date %s: %v", currentDate.Format("2006-01-02"), err)
			continue
		}
	}

	log.Println("Seeding completed successfully.")
}

func generateWallets(walletRepo repository.WalletRepo, count int, createdAt time.Time) ([]*model.Wallet, error) {
	var (
		wallets   = make([]*model.Wallet, 0, count)
		userIDs   = generateUniqueUserIDs(count)
		wg        sync.WaitGroup
		walletMux sync.Mutex
		sem       = make(chan struct{}, maxConcurrentOperations)
	)

	for i := 0; i < count; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(userID string) {
			defer wg.Done()
			defer func() { <-sem }()

			// Initialize wallet with a random balance to avoid insufficient balance issues
			initialBalance := getRandomAmount(minTransactionAmount, maxTransactionAmount)
			wallet := &model.Wallet{
				UserID:    userID,
				Balance:   initialBalance,
				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			}

			err := walletRepo.CreateWallet(wallet)
			if err != nil {
				log.Printf("Failed to create wallet for user %s: %v", userID, err)
				return
			}

			walletMux.Lock()
			wallets = append(wallets, wallet)
			walletMux.Unlock()
		}(userIDs[i])
	}

	wg.Wait()

	return wallets, nil
}

func generateTransactions(transactionRepo repository.TransactionRepo, wallets []*model.Wallet, count int, createdAt time.Time) error {
	var (
		wg  sync.WaitGroup
		sem = make(chan struct{}, maxConcurrentOperations)
	)

	for i := 0; i < count; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			err := createRandomTransactionWithRetry(transactionRepo, wallets, createdAt)
			if err != nil {
				log.Printf("Failed to create transaction after retries: %v", err)
			}
		}()
	}

	wg.Wait()

	return nil
}

func createRandomTransactionWithRetry(transactionRepo repository.TransactionRepo, wallets []*model.Wallet, createdAt time.Time) error {
	var err error
	for retries := 0; retries < maxRetries; retries++ {
		err = createRandomTransaction(transactionRepo, wallets, createdAt)
		if err == nil {
			return nil
		}
		log.Printf("Transaction failed, retrying (%d/%d): %v", retries+1, maxRetries, err)
		time.Sleep(retryDelay)
	}
	return err
}

func createRandomTransaction(transactionRepo repository.TransactionRepo, wallets []*model.Wallet, createdAt time.Time) error {
	// Select wallets that are not currently in use
	fromWallet, toWallet := selectWallets(wallets)
	if fromWallet == nil {
		return fmt.Errorf("could not find available wallets for transaction")
	}

	amount := getRandomAmount(minTransactionAmount, maxTransactionAmount)
	initiatedBy := getRandomInitiator()

	// Mark wallets as in use
	lockWallet(fromWallet.ID)
	defer unlockWallet(fromWallet.ID)

	if toWallet != nil {
		lockWallet(toWallet.ID)
		defer unlockWallet(toWallet.ID)
	}

	// Create transfer if both wallets are selected
	if toWallet != nil {
		transactionOut := &model.Transaction{
			WalletID:    fromWallet.ID,
			Amount:      amount.Neg(),
			Type:        model.TransactionTypeTransferOut,
			InitiatedBy: initiatedBy,
			CreatedAt:   createdAt,
		}
		transactionIn := &model.Transaction{
			WalletID:      toWallet.ID,
			Amount:        amount,
			Type:          model.TransactionTypeTransferIn,
			ReferenceType: &model.TransactionReferenceTypeTransfer,
			InitiatedBy:   initiatedBy,
			CreatedAt:     createdAt,
		}
		return transactionRepo.Transfer(transactionOut, fromWallet.Version, transactionIn, toWallet.Version)
	}

	// Ensure wallet has enough balance for withdrawals, purchases, or transfer outs
	if transactionType := rand.Intn(4); transactionType != 0 { // DEPOSIT is transactionType 0
		if fromWallet.Balance.LessThan(amount) {
			log.Printf("Skipping transaction due to insufficient balance in wallet %s", fromWallet.ID)
			return nil
		}
	}

	// Create single transaction (Deposit, Withdraw, Purchase, Refund)
	transaction := &model.Transaction{
		WalletID:    fromWallet.ID,
		CreatedAt:   createdAt,
		InitiatedBy: initiatedBy,
	}

	transactionType := rand.Intn(4)
	switch transactionType {
	case 0: // DEPOSIT
		transaction.Type = model.TransactionTypeDeposit
		transaction.Amount = amount
		referenceId := fmt.Sprintf("payment_%d", rand.Intn(paymentTransactionIds))
		transaction.ReferenceID = &referenceId
		transaction.ReferenceType = &model.TransactionReferenceTypeBankTransaction
	case 1: // WITHDRAW
		transaction.Type = model.TransactionTypeWithdraw
		transaction.Amount = amount.Neg()
		referenceId := fmt.Sprintf("payment_%d", rand.Intn(paymentTransactionIds))
		transaction.ReferenceID = &referenceId
		transaction.ReferenceType = &model.TransactionReferenceTypeBankTransaction
	case 2: // PURCHASE
		transaction.Type = model.TransactionTypePurchase
		transaction.Amount = amount.Neg()
		referenceId := fmt.Sprintf("order_%d", rand.Intn(orderIds))
		transaction.ReferenceID = &referenceId
		transaction.ReferenceType = &model.TransactionReferenceTypeOrder
	case 3: // REFUND
		transaction.Type = model.TransactionTypeRefund
		transaction.Amount = amount
		referenceId := fmt.Sprintf("order_%d", rand.Intn(orderIds))
		transaction.ReferenceID = &referenceId
		transaction.ReferenceType = &model.TransactionReferenceTypeOrder
	}

	return transactionRepo.Create(transaction, fromWallet.Version)
}

func lockWallet(walletID string) {
	walletLocks.Store(walletID, struct{}{})
}

func unlockWallet(walletID string) {
	walletLocks.Delete(walletID)
}

func selectWallets(wallets []*model.Wallet) (*model.Wallet, *model.Wallet) {
	var fromWallet, toWallet *model.Wallet

	for _, wallet := range wallets {
		if _, loaded := walletLocks.Load(wallet.ID); !loaded {
			if fromWallet == nil {
				fromWallet = wallet
			} else if toWallet == nil && rand.Intn(10) < 4 {
				toWallet = wallet
				break
			}
		}
	}
	return fromWallet, toWallet
}

func generateUniqueUserIDs(count int) []string {
	userIDs := make(map[string]struct{}, count)
	for len(userIDs) < count {
		userID := fmt.Sprintf("user_%d", rand.Intn(count*10))
		userIDs[userID] = struct{}{}
	}

	ids := make([]string, 0, count)
	for id := range userIDs {
		ids = append(ids, id)
	}
	return ids
}

func getRandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func getRandomAmount(min, max float64) decimal.Decimal {
	amount := min + rand.Float64()*(max-min)
	return decimal.NewFromFloat(math.Round(amount*100) / 100)
}

func getRandomInitiator() string {
	switch rand.Intn(3) {
	case 0:
		return model.TransactionInitiatedBySystem
	case 1:
		return model.TransactionInitiatedByUser
	case 2:
		return model.TransactionInitiatedByBackoffice
	default:
		return model.TransactionInitiatedBySystem
	}
}
