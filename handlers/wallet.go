package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"wallet-api/models"
	"wallet-api/services"
)

type WalletHandler struct {
	walletService services.IWalletService
}

func NewWalletHandler(walletService services.IWalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}


func (h *WalletHandler) Create(c *gin.Context) {
	var wallet models.Wallet
	if err := c.ShouldBindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if wallet.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	err := h.walletService.Create(&wallet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response := gin.H{
		"id":         wallet.ID,
		"user_id":    wallet.UserID,
		"balance":    wallet.Balance,
		"created_at": wallet.CreatedAt,
		"updated_at": wallet.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *WalletHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet ID"})
		return
	}

	wallet, err := h.walletService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func (h *WalletHandler) GetByUserID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	wallets, err := h.walletService.GetByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wallets not found"})
		return
	}

	var response []models.WalletResponse
	for _, w := range wallets {
		response = append(response, models.WalletResponse{
			ID:        w.ID,
			UserID:    w.UserID,
			Balance:   w.Balance,
			CreatedAt: w.CreatedAt,
			UpdatedAt: w.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}