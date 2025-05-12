package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"wallet-api/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			Name:  "John Doe",
			Email: "john@example.com",
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

		user := &models.User{
			Name: "John Doe",
			// Missing email
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
			Name:  "John Doe",
			Email: "john@example.com",
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

func TestUserHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful retrieval", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		user := &models.User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockService.On("GetByID", uint(1)).Return(user, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		handler.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, response.ID)
		assert.Equal(t, user.Name, response.Name)
		assert.Equal(t, user.Email, response.Email)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		handler.GetByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "GetByID")
	})

	t.Run("user not found", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		mockService.On("GetByID", uint(1)).Return(nil, errors.New("not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		handler.GetByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func toString(v interface{}) string {
	switch v := v.(type) {
	case float64:
		return fmt.Sprintf("%d", int(v))
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
