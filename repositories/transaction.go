package repositories

import (
	"wallet-api/models"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	return r.DB.Create(transaction).Error
}

func (r *TransactionRepository) GetByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.DB.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) GetByWalletID(walletID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.DB.Where("source_wallet_id = ? OR target_wallet_id = ?", walletID, walletID).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}