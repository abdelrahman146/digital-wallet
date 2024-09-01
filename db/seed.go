package main

import (
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/internal/service"
	"digital-wallet/pkg/resources"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	numWallets             = 10000
	numTransactions        = 100000
	maxConcurrentWallets   = 10
	maxConcurrentTransfers = 15
	transactionTypes       = 5 // DEPOSIT, WITHDRAW, PURCHASE, REFUND, TRANSFER_IN, TRANSFER_OUT
	startDate              = "2024-01-01"
	endDate                = "2024-08-30"
	minTransactionAmount   = 10.0
	maxTransactionAmount   = 1000.0
	paymentTransactionIds  = 1000000
	orderIds               = 1000000
	initiators             = 3 // SYSTEM, USER, BACKOFFICE
)

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Initialize the database connection
	db := resources.InitDB()

	// Initialize repositories and services
	repos := &repository.Repos{
		Wallet:      repository.NewWalletRepo(db),
		Transaction: repository.NewTransactionRepo(db),
	}
	transactionService := service.NewTransactionService(repos)
	walletService := service.NewWalletService(repos)

	// Generate wallets
	log.Println("Starting wallet generation...")
	wallets, err := generateWallets(walletService, numWallets)
	if err != nil {
		log.Fatalf("Failed to generate wallets: %v", err)
	}
	log.Printf("Successfully generated %d wallets.", len(wallets))

	// Generate transactions
	log.Println("Starting transaction generation...")
	err = generateTransactions(transactionService, wallets, numTransactions)
	if err != nil {
		log.Fatalf("Failed to generate transactions: %v", err)
	}
	log.Printf("Successfully generated %d transactions.", numTransactions)

	// Verify data consistency
	log.Println("Verifying data consistency...")
	err = verifyDataConsistency(repos)
	if err != nil {
		log.Fatalf("Data inconsistency detected: %v", err)
	}
	log.Println("Data consistency verified successfully.")
}

func generateWallets(walletService service.WalletService, count int) ([]*model.Wallet, error) {
	var (
		wallets   = make([]*model.Wallet, 0, count)
		userIDs   = generateUniqueUserIDs(count)
		wg        sync.WaitGroup
		errChan   = make(chan error, count)
		walletMux sync.Mutex
		sem       = make(chan struct{}, maxConcurrentWallets)
	)

	for i := 0; i < count; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(userID string) {
			defer wg.Done()
			defer func() { <-sem }()

			wallet, err := walletService.CreateWallet(userID)
			if err != nil {
				errChan <- fmt.Errorf("failed to create wallet for user %s: %v", userID, err)
				return
			}

			walletMux.Lock()
			wallets = append(wallets, wallet)
			walletMux.Unlock()
		}(userIDs[i])
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return wallets, nil
}

func generateUniqueUserIDs(count int) []string {
	userIDs := make(map[string]struct{}, count)
	for len(userIDs) < count {
		userID, _ := uuid.NewUUID()
		userIDs[userID.String()] = struct{}{}
	}
	ids := make([]string, 0, count)
	for id := range userIDs {
		ids = append(ids, id)
	}
	return ids
}

func generateTransactions(transactionService service.TransactionService, wallets []*model.Wallet, count int) error {
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, count)
		sem     = make(chan struct{}, maxConcurrentTransfers)
	)

	for i := 0; i < count; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			err := createRandomTransaction(transactionService, wallets)
			if err != nil {
				errChan <- err
			}
		}()
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

func createRandomTransaction(transactionService service.TransactionService, wallets []*model.Wallet) error {
	wallet := wallets[rand.Intn(len(wallets))]
	transactionType := rand.Intn(transactionTypes)
	initiatedBy := getRandomInitiator()

	switch transactionType {
	case 0: // DEPOSIT
		req := &service.DepositRequest{
			UserID:               wallet.UserID,
			Amount:               getRandomAmount(minTransactionAmount, maxTransactionAmount),
			PaymentTransactionId: fmt.Sprintf("payment_%d", rand.Intn(paymentTransactionIds)),
		}
		_, err := transactionService.Deposit(req, initiatedBy)
		return err
	case 1: // WITHDRAW
		req := &service.WithdrawRequest{
			UserID:               wallet.UserID,
			Amount:               -getRandomAmount(minTransactionAmount, maxTransactionAmount),
			PaymentTransactionId: fmt.Sprintf("payment_%d", rand.Intn(paymentTransactionIds)),
		}
		_, err := transactionService.Withdraw(req, initiatedBy)
		return err
	case 2: // PURCHASE
		req := &service.PurchaseRequest{
			UserID:  wallet.UserID,
			Amount:  -getRandomAmount(minTransactionAmount, maxTransactionAmount),
			OrderId: fmt.Sprintf("order_%d", rand.Intn(orderIds)),
		}
		_, err := transactionService.Purchase(req, initiatedBy)
		return err
	case 3: // REFUND
		req := &service.RefundRequest{
			UserID:  wallet.UserID,
			Amount:  getRandomAmount(minTransactionAmount, maxTransactionAmount),
			OrderId: fmt.Sprintf("order_%d", rand.Intn(orderIds)),
		}
		_, err := transactionService.Refund(req, initiatedBy)
		return err
	case 4: // TRANSFER
		// To simulate TRANSFER_IN, we need to perform a transfer from another wallet
		return createRandomTransfer(transactionService, wallets, initiatedBy)
	default:
		return errors.New("unknown transaction type")
	}
}

func createRandomTransfer(transactionService service.TransactionService, wallets []*model.Wallet, initiatedBy string) error {
	fromIndex := rand.Intn(len(wallets))
	toIndex := rand.Intn(len(wallets))
	// Ensure from and to wallets are different
	for fromIndex == toIndex {
		toIndex = rand.Intn(len(wallets))
	}

	req := &service.TransferRequest{
		FromUserID: wallets[fromIndex].UserID,
		ToUserID:   wallets[toIndex].UserID,
		Amount:     getRandomAmount(minTransactionAmount, maxTransactionAmount),
	}
	_, err := transactionService.Transfer(req, initiatedBy)
	return err
}

func getRandomAmount(min, max float64) float64 {
	amount := min + rand.Float64()*(max-min)
	return math.Round(amount*100) / 100
}

func getRandomInitiator() string {
	switch rand.Intn(initiators) {
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

func verifyDataConsistency(repos *repository.Repos) error {
	// Get sum of all wallet balances
	walletsSum, err := repos.Wallet.GetWalletsSum()
	if err != nil {
		return fmt.Errorf("failed to get wallets sum: %v", err)
	}

	// Get sum of all transactions
	transactionsSum, err := repos.Transaction.GetTransactionsSum()
	if err != nil {
		return fmt.Errorf("failed to get transactions sum: %v", err)
	}

	// Allow a small delta due to potential rounding errors
	delta := 0.0001
	if abs(walletsSum-transactionsSum) > delta {
		return fmt.Errorf("sum mismatch: wallets sum = %f, transactions sum = %f", walletsSum, transactionsSum)
	}

	return nil
}

func abs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}
