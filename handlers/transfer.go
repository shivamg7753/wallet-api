package handlers

import (
	"net/http"
	"strconv"

	"wallet-api/models"
	"wallet-api/services"

	"github.com/gin-gonic/gin"
)

// handlers/transfer_handler.go

type TransferHandler struct {
	transferService services.ITransferService
}

func NewTransferHandler(service services.ITransferService) *TransferHandler {
	return &TransferHandler{transferService: service}
}

type TransferRequest struct {
	SourceWalletID uint  `json:"source_wallet_id" binding:"required"`
	TargetWalletID uint  `json:"target_wallet_id" binding:"required"`
	Amount         int64 `json:"amount" binding:"required,gt=0"`
}

func (h *TransferHandler) Transfer(c *gin.Context) {
	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.transferService.Transfer(req.SourceWalletID, req.TargetWalletID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transfer successful"})
}

type DepositRequest struct {
	WalletID uint  `json:"wallet_id" binding:"required"`
	Amount   int64 `json:"amount" binding:"required,gt=0"`
}

func (h *TransferHandler) Deposit(c *gin.Context) {
	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.transferService.Deposit(req.WalletID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deposit successful"})
}

func (h *TransferHandler) GetTransactions(c *gin.Context) {
	walletID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet ID"})
		return
	}

	transactions, err := h.transferService.GetTransactionsByWalletID(uint(walletID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transactions not found"})
		return
	}

	var response []models.TransferResponse
	for _, t := range transactions {
		response = append(response, models.TransferResponse{
			ID:              t.ID,
			SourceWalletID:  t.SourceWalletID,
			TargetWalletID:  t.TargetWalletID,
			Amount:          t.Amount,
			Type:            t.Type,
			ReferenceNumber: t.ReferenceNumber,
			Status:          t.Status,
			CreatedAt:       t.CreatedAt,
			UpdatedAt:       t.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}
