package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"wallet-api/models"
	"wallet-api/repositories"
	"wallet-api/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestServer(t *testing.T) (*gin.Engine, *gorm.DB) {
	// Setup test database
	db := setupTestDB(t)
	truncateTables(t, db)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	walletRepo := repositories.NewWalletRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	walletService := services.NewWalletService(walletRepo, userRepo)
	transferService := services.NewTransferService(transactionRepo, walletRepo, db)

	// Initialize handlers
	userHandler := NewUserHandler(userService)
	walletHandler := NewWalletHandler(walletService)
	transferHandler := NewTransferHandler(transferService)

	// Setup router
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		// User routes
		api.POST("/users", userHandler.Create)
		api.GET("/users/:id", userHandler.GetByID)

		// Wallet routes
		api.POST("/wallets", walletHandler.Create)
		api.GET("/wallets/:id", walletHandler.GetByID)
		api.GET("/users/:id/wallets", walletHandler.GetByUserID)

		// Transfer routes
		api.POST("/transfers", transferHandler.Transfer)
		api.POST("/deposits", transferHandler.Deposit)
		api.GET("/wallets/:id/transactions", transferHandler.GetTransactions)
	}

	return router, db
}

func TestAPI_CompleteFlow(t *testing.T) {
	router, db := setupTestServer(t)
	defer teardownTestDB(t, db)

	// Step 1: Create a user
	userPayload := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
	}
	jsonUser, _ := json.Marshal(userPayload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var userResponse models.User
	err := json.Unmarshal(w.Body.Bytes(), &userResponse)
	assert.NoError(t, err)
	assert.NotZero(t, userResponse.ID)
	userID := userResponse.ID

	// Step 2: Create a wallet for the user
	walletPayload := map[string]interface{}{
		"user_id": userID,
	}
	jsonWallet, _ := json.Marshal(walletPayload)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/wallets", bytes.NewBuffer(jsonWallet))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var walletResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &walletResponse)
	assert.NoError(t, err)
	assert.NotZero(t, walletResponse["id"])
	walletID := uint(walletResponse["id"].(float64))

	// Step 3: Deposit money into the wallet
	depositPayload := map[string]interface{}{
		"wallet_id": walletID,
		"amount":    1000,
	}
	jsonDeposit, _ := json.Marshal(depositPayload)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/deposits", bytes.NewBuffer(jsonDeposit))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var depositResponse map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &depositResponse)
	assert.NoError(t, err)
	assert.Equal(t, "deposit successful", depositResponse["message"])

	// Step 4: Create another user and wallet for transfer
	user2Payload := map[string]interface{}{
		"name":  "Jane Doe",
		"email": "jane@example.com",
	}
	jsonUser2, _ := json.Marshal(user2Payload)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonUser2))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var user2Response models.User
	err = json.Unmarshal(w.Body.Bytes(), &user2Response)
	assert.NoError(t, err)
	assert.NotZero(t, user2Response.ID)

	// Create wallet for second user
	wallet2Payload := map[string]interface{}{
		"user_id": user2Response.ID,
	}
	jsonWallet2, _ := json.Marshal(wallet2Payload)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/wallets", bytes.NewBuffer(jsonWallet2))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var wallet2Response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &wallet2Response)
	assert.NoError(t, err)
	assert.NotZero(t, wallet2Response["id"])
	wallet2ID := uint(wallet2Response["id"].(float64))

	// Step 5: Transfer money between wallets
	transferPayload := map[string]interface{}{
		"source_wallet_id": walletID,
		"target_wallet_id": wallet2ID,
		"amount":           500,
	}
	jsonTransfer, _ := json.Marshal(transferPayload)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/transfers", bytes.NewBuffer(jsonTransfer))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var transferResponse map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &transferResponse)
	assert.NoError(t, err)
	assert.Equal(t, "transfer successful", transferResponse["message"])

	// Step 6: Verify wallet balances
	// Check source wallet
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/wallets/%d", walletID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var sourceWallet models.Wallet
	err = json.Unmarshal(w.Body.Bytes(), &sourceWallet)
	assert.NoError(t, err)
	assert.Equal(t, int64(500), sourceWallet.Balance) // 1000 - 500

	// Check target wallet
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/wallets/%d", wallet2ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var targetWallet models.Wallet
	err = json.Unmarshal(w.Body.Bytes(), &targetWallet)
	assert.NoError(t, err)
	assert.Equal(t, int64(500), targetWallet.Balance) // 0 + 500

	// Step 7: Check transaction history
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/wallets/%d/transactions", walletID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var transactions []models.Transaction
	err = json.Unmarshal(w.Body.Bytes(), &transactions)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(transactions)) // Should have deposit and transfer transactions
}

func TestAPI_ErrorCases(t *testing.T) {
	router, db := setupTestServer(t)
	defer teardownTestDB(t, db)

	t.Run("create user with invalid email", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":  "John Doe",
			"email": "invalid-email",
		}
		jsonData, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("create wallet without user_id", func(t *testing.T) {
		payload := map[string]interface{}{}
		jsonData, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallets", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("transfer with insufficient balance", func(t *testing.T) {
		// Create two users and wallets
		user1 := createTestUser(t, router, "John Doe", "john@example.com")
		user2 := createTestUser(t, router, "Jane Doe", "jane@example.com")
		wallet1 := createTestWallet(t, router, user1.ID)
		wallet2 := createTestWallet(t, router, user2.ID)

		// Try to transfer more than balance
		payload := map[string]interface{}{
			"source_wallet_id": wallet1.ID,
			"target_wallet_id": wallet2.ID,
			"amount":           1000, // Try to transfer 1000 when balance is 0
		}
		jsonData, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/transfers", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("get non-existent user", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/999999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("get non-existent wallet", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/wallets/999999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// Helper functions for creating test data
func createTestUser(t *testing.T, router *gin.Engine, name, email string) models.User {
	payload := map[string]interface{}{
		"name":  name,
		"email": email,
	}
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	return response
}

func createTestWallet(t *testing.T, router *gin.Engine, userID uint) models.Wallet {
	payload := map[string]interface{}{
		"user_id": userID,
	}
	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallets", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response models.Wallet
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	return response
}
