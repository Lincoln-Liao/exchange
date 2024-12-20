# Exchange Wallet System

## Introduction

Exchange Wallet System is a backend service for managing user wallets and transactions. It provides the ability to deposit, withdraw, transfer, and query balances and transaction history, and ensures that all operations are performed within transactions to ensure data consistency.

## Deployment

### Pre-Request
- [Docker](https://docs.docker.com/get-docker/)
- Go 1.23+

### steps
1. clone project
```bash
    git clone https://github.com/Lincoln-Liao/exchange.git
    cd exchange-wallet-system
```
2. Install dependencies
```bash
    go mod download
```
3. Deploy infrastructure
```bash
    sh scripts/server_developing.sh
```
4. Database migrations
- install golan-migrate
```bash
    brew install golang-migrate
```
- migrate
```bash
    migrate -source file://internal/ports/persistence/migrations -database "postgres://root:example@localhost:5432/exchange?sslmode=disable" up
```
- Reference:
  [golang-migrate](https://github.com/golang-migrate/migrate)
5. run project
```bash
    go run cmd/server/main.go
```

## API Document
./doc/postman/wallet/wallet.postman_collection.json