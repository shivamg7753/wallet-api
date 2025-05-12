package repositories

import (
	"errors"

	"wallet-api/models"
	"gorm.io/gorm"
)

type WalletRepository struct {
	DB *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{DB: db}
}

func (r *WalletRepository) Create(wallet *models.Wallet) error {
	return r.DB.Create(wallet).Error
}


func (r *WalletRepository) GetByID(id uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.DB.Preload("User").First(&wallet, id).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) GetByUserID(userID uint) ([]models.Wallet, error) {
	var wallets []models.Wallet
	err := r.DB.Where("user_id = ?", userID).Find(&wallets).Error
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

func (r *WalletRepository) UpdateBalance(id uint, amount int64) error {
	result := r.DB.Model(&models.Wallet{}).Where("id = ?", id).Update("balance", gorm.Expr("balance + ?", amount))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("wallet not found")
	}
	return nil
}
