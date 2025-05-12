# Simple Wallet Management REST API

A RESTful backend service built in GoLang to manage user wallets and basic transactions between them.

## Features

- User management (create users)
- Wallet management (create wallets, check balances)
- Transaction processing (transfer funds between wallets)
- Transaction history (view all transactions for a wallet)

## Tech Stack

- **Language**: GoLang
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **Containerization**: Docker & Docker Compose
- **Unit Testing**: Go testing package with testify

## Project Structure

```
wallet-api/
├── cmd/
│   └── main.go            # Application entry point
├── docker-compose.yml     # Docker compose configuration
├── Dockerfile             # Docker build instructions
├── go.mod                 # Go modules definition
├── go.sum                 # Go modules checksums
├── handlers/              # HTTP request handlers
│   ├── transfer.go
│   ├── user.go
│   ├── wallet.go
|   ├── transfer_test.go
|   ├── user_test.go
|   └── wallet_test.go
├── models/                # Data models
│   ├── transaction.go
│   ├── user.go
│   └── wallet.go
├── repositories/          # Database interactions
│   ├── transaction.go
│   ├── user.go
│   └── wallet.go
├── services/              # Business logic
│   ├── transfer.go
│   ├── user.go
│   └── wallet.go
└── README.md              # This file
```

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- Git

### Installation & Setup

#### Using Docker (Recommended)

1. Clone the repository:
   ```bash
   git clone https://github.com/shivamg7753/wallet-api.git
   cd wallet-api
   ```

2. Run the application using Docker Compose:
   ```bash
   docker-compose up --build
   ```

3. The API will be available at `http://localhost:8080`

#### Manual Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/shivamg7753/wallet-api.git
   cd wallet-api
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up a PostgreSQL database and update the connection string in `cmd/main.go` or set the `DATABASE_URL` environment variable.

4. Run the application:
   ```bash
   go run cmd/main.go
   ```

5. The API will be available at `http://localhost:8080`

### Running Tests

```bash
go test ./... -v
```

## API Endpoints

### Users

#### Create a user

- **URL**: `/api/v1/users`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "name": "John Doe",
    "email": "john@example.com"
  }
  ```
- **Response**: 
  ```json
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2025-05-12T12:00:00Z",
    "updated_at": "2025-05-12T12:00:00Z"
  }
  ```

#### Get a user by ID

- **URL**: `/api/v1/users/:id`
- **Method**: `GET`
- **Response**: 
  ```json
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2025-05-12T12:00:00Z",
    "updated_at": "2025-05-12T12:00:00Z"
  }
  ```

### Wallets

#### Create a wallet

- **URL**: `/api/v1/wallets`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "user_id": 1
  }
  ```
- **Response**: 
  ```json
  {
    "id": 1,
    "user_id": 1,
    "balance": 0,
    "created_at": "2025-05-12T12:00:00Z",
    "updated_at": "2025-05-12T12:00:00Z"
  }
  ```

#### Get a wallet by ID

- **URL**: `/api/v1/wallets/:id`
- **Method**: `GET`
- **Response**: 
  ```json
  {
    "id": 1,
    "user_id": 1,
    "balance": 1000,
    "created_at": "2025-05-12T12:00:00Z",
    "updated_at": "2025-05-12T12:00:00Z"
  }
  ```

#### Get wallets by user ID

- **URL**: `/api/v1/users/:userID/wallets`
- **Method**: `GET`
- **Response**: 
  ```json
  [
    {
      "id": 1,
      "user_id": 1,
      "balance": 1000,
      "created_at": "2025-05-12T12:00:00Z",
      "updated_at": "2025-05-12T12:00:00Z"
    },
    {
      "id": 2,
      "user_id": 1,
      "balance": 500,
      "created_at": "2025-05-12T12:30:00Z",
      "updated_at": "2025-05-12T12:30:00Z"
    }
  ]
  ```

### Transfers

#### Transfer funds between wallets

- **URL**: `/api/v1/transfers`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "source_wallet_id": 1,
    "target_wallet_id": 2,
    "amount": 500
  }
  ```
- **Response**: 
  ```json
  {
    "message": "transfer successful"
  }
  ```

#### Deposit funds to a wallet

- **URL**: `/api/v1/deposits`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "wallet_id": 1,
    "amount": 1000
  }
  ```
- **Response**: 
  ```json
  {
    "message": "deposit successful"
  }
  ```

#### Get transaction history for a wallet

- **URL**: `/api/v1/wallets/:walletID/transactions`
- **Method**: `GET`
- **Response**: 
  ```json
  [
    {
      "id": 1,
      "source_wallet_id": null,
      "target_wallet_id": 1,
      "amount": 1000,
      "type": "deposit",
      "reference_number": "DEP-1715433600000000000",
      "status": "completed",
      "created_at": "2025-05-12T12:00:00Z",
      "updated_at": "2025-05-12T12:00:00Z"
    },
    {
      "id": 2,
      "source_wallet_id": 1,
      "target_wallet_id": 2,
      "amount": 500,
      "type": "transfer",
      "reference_number": "TRF-1715435400000000000",
      "status": "completed",
      "created_at": "2025-05-12T12:30:00Z",
      "updated_at": "2025-05-12T12:30:00Z"
    }
  ]
  ```

## Error Handling

- All endpoints return error messages in JSON format with relevant status codes.

 ## Common Responses

**400 Bad Request** – Invalid or missing parameters

**404 Not Found** – Resource not found

**409 Conflict** – Wallet has insufficient funds or duplicate entry

**500 Internal Server Error** – Server-side issue

