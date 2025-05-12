package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"wallet-api/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock WalletService
type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) Create(wallet *models.Wallet) error {
	args := m.Called(wallet)
	return args.Error(0)
}

func (m *MockWalletService) GetByID(id uint) (*models.Wallet, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Wallet), args.Error(1)
}

func (m *MockWalletService) GetByUserID(userID uint) ([]models.Wallet, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Wallet), args.Error(1)
}

func TestWalletHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful creation", func(t *testing.T) {
		mockService := new(MockWalletService)
		handler := NewWalletHandler(mockService)

		wallet := &models.Wallet{
			UserID: 1,
		}

		mockService.On("Create", mock.AnythingOfType("*models.Wallet")).Return(nil).Run(func(args mock.Arguments) {
			w := args.Get(0).(*models.Wallet)
			w.ID = 1
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonWallet, _ := json.Marshal(wallet)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/wallets", bytes.NewBuffer(jsonWallet))
		c.Request.Header.Add("Content-Type", "application/json")

		handler.Create(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response["id"])
		assert.Equal(t, float64(0), response["balance"])
		assert.Equal(t, float64(1), response["user_id"])

		mockService.AssertExpectations(t)
	})

	t.Run("missing user_id", func(t *testing.T) {
		mockService := new(MockWalletService)
		handler := NewWalletHandler(mockService)

		wallet := &models.Wallet{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonWallet, _ := json.Marshal(wallet)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/wallets", bytes.NewBuffer(jsonWallet))
		c.Request.Header.Add("Content-Type", "application/json")

		handler.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "Create")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(MockWalletService)
		handler := NewWalletHandler(mockService)

		wallet := &models.Wallet{
			UserID: 1,
		}

		mockService.On("Create", mock.AnythingOfType("*models.Wallet")).Return(errors.New("service error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonWallet, _ := json.Marshal(wallet)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/wallets", bytes.NewBuffer(jsonWallet))
		c.Request.Header.Add("Content-Type", "application/json")

		handler.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestWalletHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful retrieval", func(t *testing.T) {
		mockService := new(MockWalletService)
		handler := NewWalletHandler(mockService)

		wallet := &models.Wallet{
			ID:      1,
			UserID:  1,
			Balance: 100,
		}

		mockService.On("GetByID", uint(1)).Return(wallet, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		handler.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Wallet
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, wallet.ID, response.ID)
		assert.Equal(t, wallet.UserID, response.UserID)
		assert.Equal(t, wallet.Balance, response.Balance)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockService := new(MockWalletService)
		handler := NewWalletHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		handler.GetByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "GetByID")
	})

	t.Run("wallet not found", func(t *testing.T) {
		mockService := new(MockWalletService)
		handler := NewWalletHandler(mockService)

		mockService.On("GetByID", uint(1)).Return(nil, errors.New("not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		handler.GetByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestWalletHandler_GetByUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful retrieval", func(t *testing.T) {
		mockService := new(MockWalletService)
		handler := NewWalletHandler(mockService)

		wallets := []models.Wallet{
			{
				ID:      1,
				UserID:  1,
				Balance: 100,
			},
			{
				ID:      2,
				UserID:  1,
				Balance: 200,
			},
		}

		mockService.On("GetByUserID", uint(1)).Return(wallets, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		handler.GetByUserID(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.WalletResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response))
		assert.Equal(t, wallets[0].ID, response[0].ID)
		assert.Equal(t, wallets[0].UserID, response[0].UserID)
		assert.Equal(t, wallets[0].Balance, response[0].Balance)
		assert.Equal(t, wallets[1].ID, response[1].ID)
		assert.Equal(t, wallets[1].UserID, response[1].UserID)
		assert.Equal(t, wallets[1].Balance, response[1].Balance)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid user id", func(t *testing.T) {
		mockService := new(MockWalletService)
		handler := NewWalletHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		handler.GetByUserID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "GetByUserID")
	})

	t.Run("no wallets found", func(t *testing.T) {
		mockService := new(MockWalletService)
		handler := NewWalletHandler(mockService)

		mockService.On("GetByUserID", uint(1)).Return(nil, errors.New("not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		handler.GetByUserID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}
