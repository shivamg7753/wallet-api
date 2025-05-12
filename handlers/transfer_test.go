package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"wallet-api/models"
)

// Mock TransferService
type MockTransferService struct {
	mock.Mock
}

func (m *MockTransferService) Transfer(sourceWalletID, targetWalletID uint, amount int64) error {
	args := m.Called(sourceWalletID, targetWalletID, amount)
	return args.Error(0)
}

func (m *MockTransferService) Deposit(walletID uint, amount int64) error {
	args := m.Called(walletID, amount)
	return args.Error(0)
}

func (m *MockTransferService) GetTransactionsByWalletID(walletID uint) ([]models.Transaction, error) {
	args := m.Called(walletID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func TestTransferHandler_Transfer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful transfer", func(t *testing.T) {
		mockService := new(MockTransferService)
		handler := NewTransferHandler(mockService)

		req := TransferRequest{
			SourceWalletID: 1,
			TargetWalletID: 2,
			Amount:         100,
		}

		mockService.On("Transfer", uint(1), uint(2), int64(100)).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		jsonReq, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/transfers", bytes.NewBuffer(jsonReq))
		c.Request.Header.Add("Content-Type", "application/json")

		handler.Transfer(c)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "transfer successful", response["message"])
		
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockService := new(MockTransferService)
		handler := NewTransferHandler(mockService)

		// Missing required fields
		req := map[string]interface{}{
			"source_wallet_id": 1,
			// Missing target_wallet_id and amount
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		jsonReq, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/transfers", bytes.NewBuffer(jsonReq))
		c.Request.Header.Add("Content-Type", "application/json")

		handler.Transfer(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "Transfer")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(MockTransferService)
		handler := NewTransferHandler(mockService)

		req := TransferRequest{
			SourceWalletID: 1,
			TargetWalletID: 2,
			Amount:         100,
		}

		mockService.On("Transfer", uint(1), uint(2), int64(100)).Return(errors.New("insufficient balance"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		jsonReq, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/transfers", bytes.NewBuffer(jsonReq))
		c.Request.Header.Add("Content-Type", "application/json")

		handler.Transfer(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestTransferHandler_Deposit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful deposit", func(t *testing.T) {
		mockService := new(MockTransferService)
		handler := NewTransferHandler(mockService)

		req := DepositRequest{
			WalletID: 1,
			Amount:   100,
		}

		mockService.On("Deposit", uint(1), int64(100)).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		jsonReq, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/deposits", bytes.NewBuffer(jsonReq))
		c.Request.Header.Add("Content-Type", "application/json")

		handler.Deposit(c)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "deposit successful", response["message"])
		
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockService := new(MockTransferService)
		handler := NewTransferHandler(mockService)

		// Invalid amount
		req := map[string]interface{}{
			"wallet_id": 1,
			"amount":    -100, // Negative amount should fail validation
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		jsonReq, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/deposits", bytes.NewBuffer(jsonReq))
		c.Request.Header.Add("Content-Type", "application/json")

		handler.Deposit(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "Deposit")
	})
}

func TestTransferHandler_GetTransactions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful retrieval", func(t *testing.T) {
		mockService := new(MockTransferService)
		handler := NewTransferHandler(mockService)

		transactions := []models.Transaction{
			{
				ID:             1,
				TargetWalletID: 1,
				Amount:         100,
				Type:           models.TransactionTypeDeposit,
				Status:         "completed",
			},
			{
				ID:             2,
				SourceWalletID: func() *uint { id := uint(1); return &id }(),
				TargetWalletID: 2,
				Amount:         50,
				Type:           models.TransactionTypeTransfer,
				Status:         "completed",
			},
		}

		mockService.On("GetTransactionsByWalletID", uint(1)).Return(transactions, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "walletID", Value: "1"}}
		
		handler.GetTransactions(c)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response []models.Transaction
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response))
		assert.Equal(t, uint(1), response[0].ID)
		assert.Equal(t, uint(2), response[1].ID)
		
		mockService.AssertExpectations(t)
	})

	t.Run("invalid wallet id", func(t *testing.T) {
		mockService := new(MockTransferService)
		handler := NewTransferHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "walletID", Value: "invalid"}}
		
		handler.GetTransactions(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "GetTransactionsByWalletID")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(MockTransferService)
		handler := NewTransferHandler(mockService)

		mockService.On("GetTransactionsByWalletID", uint(1)).Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "walletID", Value: "1"}}
		
		handler.GetTransactions(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}