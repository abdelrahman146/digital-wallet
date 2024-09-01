package service

import (
	"digital-wallet/internal/model"
	"digital-wallet/internal/repository"
	"digital-wallet/pkg/api"
)

type WalletService interface {
	CreateWallet(userId string) (*model.Wallet, error)
	GetWalletByUserID(userId string) (*model.Wallet, error)
	GetWallets(page int, limit int) (*api.List[model.Wallet], error)
	GetWalletsSum() (float64, error)
}

type walletService struct {
	repos *repository.Repos
}

func NewWalletService(repos *repository.Repos) WalletService {
	return &walletService{repos: repos}
}

func (s *walletService) CreateWallet(userId string) (*model.Wallet, error) {
	wallet := &model.Wallet{
		UserID: userId,
	}
	err := s.repos.Wallet.CreateWallet(wallet)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *walletService) GetWalletByUserID(userId string) (*model.Wallet, error) {
	wallet, err := s.repos.Wallet.GetWalletByUserID(userId)
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

func (s *walletService) GetWallets(page int, limit int) (*api.List[model.Wallet], error) {
	wallets, err := s.repos.Wallet.GetWallets(page, limit)
	if err != nil {
		return nil, err
	}
	total, err := s.repos.Wallet.GetTotalWallets()
	if err != nil {
		return nil, err
	}
	return &api.List[model.Wallet]{Items: wallets, Page: page, Limit: limit, Total: total}, nil
}
