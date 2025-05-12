package services

import (
	"errors"
	"fmt"
	"time"

	"wallet-api/models"
	"wallet-api/repositories"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ITransferService defines methods for wallet transactions
type ITransferService interface {
	Transfer(sourceWalletID, targetWalletID uint, amount int64) error
	Deposit(walletID uint, amount int64) error
	GetTransactionsByWalletID(walletID uint) ([]models.Transaction, error)
}

// TransferService implements ITransferService
type TransferService struct {
	transactionRepo *repositories.TransactionRepository
	walletRepo      *repositories.WalletRepository
	db              *gorm.DB
}

// ✅ Compile-time assertion to ensure TransferService implements ITransferService
var _ ITransferService = &TransferService{}

// NewTransferService creates a new TransferService instance
func NewTransferService(
	transactionRepo *repositories.TransactionRepository,
	walletRepo *repositories.WalletRepository,
	db *gorm.DB,
) *TransferService {
	return &TransferService{
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
		db:              db,
	}
}

func (s *TransferService) Transfer(sourceWalletID, targetWalletID uint, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	if sourceWalletID == targetWalletID {
		return errors.New("source and target wallets cannot be the same")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		var sourceWallet, targetWallet models.Wallet

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&sourceWallet, sourceWalletID).Error; err != nil {
			return err
		}

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&targetWallet, targetWalletID).Error; err != nil {
			return err
		}

		if sourceWallet.Balance < amount {
			return errors.New("insufficient balance")
		}

		sourceWallet.Balance -= amount
		if err := tx.Save(&sourceWallet).Error; err != nil {
			return err
		}

		targetWallet.Balance += amount
		if err := tx.Save(&targetWallet).Error; err != nil {
			return err
		}

		transaction := models.Transaction{
			SourceWalletID:  &sourceWalletID,
			TargetWalletID:  targetWalletID,
			Amount:          amount,
			Type:            models.TransactionTypeTransfer,
			ReferenceNumber: fmt.Sprintf("TRF-%d", time.Now().UnixNano()),
			Status:          "completed",
		}

		return tx.Create(&transaction).Error
	})
}

func (s *TransferService) Deposit(walletID uint, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&wallet, walletID).Error; err != nil {
			return err
		}

		wallet.Balance += amount
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		transaction := models.Transaction{
			TargetWalletID:  walletID,
			Amount:          amount,
			Type:            models.TransactionTypeDeposit,
			ReferenceNumber: fmt.Sprintf("DEP-%d", time.Now().UnixNano()),
			Status:          "completed",
		}

		return tx.Create(&transaction).Error
	})
}

func (s *TransferService) GetTransactionsByWalletID(walletID uint) ([]models.Transaction, error) {
	return s.transactionRepo.GetByWalletID(walletID)
}
