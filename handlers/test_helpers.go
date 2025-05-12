package handlers

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"wallet-api/models"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// var testDB *gorm.DB

func setupTestDB(t *testing.T) *gorm.DB {
	// Use test database URL or fallback to default (connect to postgres (default) database)
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	// Connect to the database
	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)

	
	_, err = db.Exec("CREATE DATABASE wallet_db")
	if err != nil && err.Error() != "pq: database \"wallet_db\" already exists" {
		assert.NoError(t, err)
	}

	// Close the connection to the default database
	db.Close()


	dsnWallet := os.Getenv("TEST_DATABASE_URL")
	if dsnWallet == "" {
		dsnWallet = "postgres://postgres:postgres@localhost:5432/wallet_db?sslmode=disable"
	}
	gormDB, err := gorm.Open(postgres.Open(dsnWallet), &gorm.Config{})
	assert.NoError(t, err)

	// Run migrations for User, Wallet, and Transaction models
	err = gormDB.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.Transaction{},
	)
	assert.NoError(t, err)

	// Run migrations
	// TODO: Add your migrations here

	return gormDB
}

func teardownTestDB(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	assert.NoError(t, err)
	sqlDB.Close()
}

func truncateTables(t *testing.T, db *gorm.DB) {
	tables := []string{"transactions", "wallets", "users"}
	for _, table := range tables {
		err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error
		assert.NoError(t, err)
	}
}
