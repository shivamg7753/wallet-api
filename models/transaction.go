package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypeTransfer TransactionType = "transfer"
)
type Transaction struct {
	ID              uint            `json:"id" gorm:"primaryKey"`
	SourceWalletID  *uint           `json:"source_wallet_id"`
	SourceWallet    *Wallet         `json:"source_wallet" gorm:"foreignKey:SourceWalletID"`
	TargetWalletID  uint            `json:"target_wallet_id" gorm:"not null"`
	TargetWallet    Wallet          `json:"target_wallet" gorm:"foreignKey:TargetWalletID"`
	Amount          int64           `json:"amount" gorm:"not null"` // Amount in smallest unit
	Type            TransactionType `json:"type" gorm:"not null"`
	ReferenceNumber string          `json:"reference_number" gorm:"size:50;index"`
	Status          string          `json:"status" gorm:"size:20;default:'completed'"` // pending, completed, failed
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `json:"deleted_at" gorm:"index"`
}


//DTO
type TransferResponse struct {
	ID              uint       `json:"id"`
	SourceWalletID  *uint      `json:"source_wallet_id"`  // Nullable
	TargetWalletID  uint      `json:"target_wallet_id"`  // Nullable
	Amount          int64      `json:"amount"`
	Type            TransactionType     `json:"type"`
	ReferenceNumber string     `json:"reference_number"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}