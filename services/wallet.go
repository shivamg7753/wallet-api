package services

import (
	"wallet-api/models"
	"wallet-api/repositories"
)

type IWalletService interface {
	Create(wallet *models.Wallet) error
	GetByID(id uint) (*models.Wallet, error)
	GetByUserID(userID uint) ([]models.Wallet, error)
}

type WalletService struct {
	walletRepo *repositories.WalletRepository
	userRepo   *repositories.UserRepository
}

func NewWalletService(
	walletRepo *repositories.WalletRepository,
	userRepo *repositories.UserRepository,
) *WalletService {
	return &WalletService{
		walletRepo: walletRepo,
		userRepo:   userRepo,
	}
}

func (s *WalletService) Create(wallet *models.Wallet) error {
	// Verify user exists
	_, err := s.userRepo.GetByID(wallet.UserID)
	if err != nil {
		return err
	}

	return s.walletRepo.Create(wallet)
}

func (s *WalletService) GetByID(id uint) (*models.Wallet, error) {
	return s.walletRepo.GetByID(id)
}

func (s *WalletService) GetByUserID(userID uint) ([]models.Wallet, error) {
	return s.walletRepo.GetByUserID(userID)
}