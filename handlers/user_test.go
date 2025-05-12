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

// Mock UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestUserHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful creation", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		user := &models.User{
			Name:  "Test User",
			Email: "test@example.com",
		}

		mockService.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		jsonUser, _ := json.Marshal(user)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonUser))
		c.Request.Header.Add("Content-Type", "application/json")

		handler.Create(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		
		var response models.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, response.Name)
		assert.Equal(t, user.Email, response.Email)
		
		mockService.AssertExpectations(t)
	})

	t.Run("missing required fields", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		// User with missing email
		user := &models.User{
			Name: "Test User",
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		jsonUser, _ := json.Marshal(user)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonUser))
		c.Request.Header.Add("Content-Type", "application/json")

		handler.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "Create")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		user := &models.User{
			Name:  "Test User",
			Email: "test@example.com",
		}

		mockService.On("Create", mock.AnythingOfType("*models.User")).Return(errors.New("service error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		jsonUser, _ := json.Marshal(user)
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonUser))
		c.Request.Header.Add("Content-Type", "application/json")

		handler.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}