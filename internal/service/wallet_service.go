package service

import (
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
)

type WalletService interface {
	CreateWallet(userID string) (*model.Wallet, error)
	GetWalletByUserID(userID string) (*model.Wallet, error)
	GetWallets(page int, limit int) (*api.List[model.Wallet], int64, error)
	GetWalletsSum() (float64, error)
}

type walletService struct {
	repos *repository.Repos
}

func NewWalletService(repos *repository.Repos) WalletService {
	return &walletService{repos: repos}
}

func (s *walletService) CreateWallet(userID string) (*model.Wallet, error) {
	wallet := &model.Wallet{
		UserID: userID,
	}
	err := s.repos.Wallet.CreateWallet(wallet)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *walletService) GetWalletByUserID(userID string) (*model.Wallet, error) {
	wallet, err := s.repos.Wallet.GetWalletByUserID(userID)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *walletService) GetWalletsSum() (float64, error) {
	sum, err := s.repos.Wallet.GetWalletsSum()
	if err != nil {
		return 0, err
	}
	return sum, nil
}

func (s *walletService) GetWallets(page int, limit int) (*api.List[model.Wallet], int64, error) {
	wallets, err := s.repos.Wallet.GetWallets(page, limit)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repos.Wallet.GetTotalWallets()
	if err != nil {
		return nil, 0, err
	}
	return &api.List[model.Wallet]{Items: wallets, Page: page, Limit: limit, Total: total}, total, nil
}
