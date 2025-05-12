# Cushon Technical Task

A Go application implementing a hexagonal architecture for managing users and transactions, with a simple React frontend.

## Prerequisites

- Go 1.21 or later
- MySQL 8.0 or later
- Node.js 16 or later
- npm or yarn

## Backend Setup

1. Create a MySQL database named `cushon`:
```sql
CREATE DATABASE cushon;
```

2. Run the schema file to create the necessary tables:
```bash
mysql -u root -p cushon < internal/adapters/secondary/persistence/mysql/schema.sql
```

3. Install Go dependencies:
```bash
go mod tidy
```

4. Set up environment variable for use with your database, in this instance, MySQL:
```bash
# Windows
set DB_PASSWORD=your_password

# Linux/MacOS
export DB_PASSWORD=your_password
```

5. Update the database configuration in `cmd/api/main.go` if needed:
```go
dbConfig := mysql.Config{
    Host:     "localhost",
    Port:     3306,
    User:     "root",
    Password: os.Getenv("DB_PASSWORD"),
    Database: "cushon",
}
```

## Frontend Setup

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
# or
yarn install
```

## Running the Application

### Backend

1. Start the Go server:
```bash
go run cmd/api/main.go
```
The backend server will start on port 8080.

### Frontend

1. Start the React development server:
```bash
cd frontend
npm start
# or
yarn start
```
The frontend will start on port 3000 and automatically open in your browser.

## API Endpoints

### Health Check
- `GET /health` - Check if the service is running

### Direct Users
- `POST /direct-users` - Create a new direct user
- `GET /direct-users/:id` - Get a direct user by ID
- `PUT /direct-users/:id` - Update a direct user
- `DELETE /direct-users/:id` - Delete a direct user

### Transactions
- `POST /transactions` - Create a new transaction
  ```json
  {
    "user_id": "uuid",
    "amount": "100.50",
    "fund_name": "Cushon Equities Fund"
  }
  ```
- `GET /transactions/:id` - Get a transaction by ID
- `GET /transactions/user/:userID` - Get all transactions for a user
- `PUT /transactions/:id` - Update a transaction
  ```json
  {
    "amount": "150.75",
    "fund_name": "Cushon Equities Fund"
  }
  ```
- `DELETE /transactions/:id` - Delete a transaction

### Fund Names
- `GET /fund-names` - Get list of available fund names

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   ├── App.js
│   │   └── index.js
│   ├── package.json
│   └── README.md
├── internal/
│   ├── adapters/
│   │   ├── primary/
│   │   │   └── http/
│   │   │       ├── direct_user_handler.go
│   │   │       └── transaction_handler.go
│   │   └── secondary/
│   │       └── persistence/
│   │           └── mysql/
│   │               ├── direct_user_repository.go
│   │               ├── transaction_repository.go
│   │               ├── connection.go
│   │               └── schema.sql
│   ├── core/
│   │   ├── domain/
│   │   │   ├── direct_user.go
│   │   │   ├── transaction.go
│   │   │   └── fund.go
│   │   ├── ports/
│   │   │   ├── input/
│   │   │   │   ├── direct_user_service.go
│   │   │   │   └── transaction_service.go
│   │   │   └── output/
│   │   │       ├── direct_user_repository.go
│   │   │       └── transaction_repository.go
│   │   └── services/
│   │       ├── direct_user_service.go
│   │       └── transaction_service.go
│   └── config/
├── go.mod
└── README.md
```

## Architecture Overview

The application follows hexagonal architecture principles:

1. **Domain Layer** (`internal/core/domain`)
   - Contains business entities and rules
   - Independent of external frameworks

2. **Ports** (`internal/core/ports`)
   - Input ports: Define use cases (what the application can do)
   - Output ports: Define interfaces for external services

3. **Adapters** (`internal/adapters`)
   - Primary adapters: Handle incoming requests (HTTP, CLI, etc.)
   - Secondary adapters: Implement output ports (databases, external services)

## Development Guidelines

- Keep domain logic independent of external frameworks
- Use interfaces (ports) to define boundaries
- Implement adapters to connect with external systems
- Write tests for each layer independently
- Frontend components should be reusable and maintainable
- Use proper error handling and loading states in the frontend
- Keep API calls centralized and consistent 

## Security Scanning

 `govulncheck` is used to scan for known vulnerabilities in Go dependencies. To use it:

1. Install govulncheck:
```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
```

2. Run the scan:
```bash
govulncheck ./...
```

No vulnerabilities found.