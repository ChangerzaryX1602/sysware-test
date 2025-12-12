## 🚀 Features

- **Authentication**: JWT-based authentication (ECDSA).
- **RBAC**: Role-Based Access Control for granular permission management.
- **Todo Management**: Create, read, update, and delete todos.


## 🛠 Tech Stack

- **Language**: [Go](https://golang.org/)
- **Framework**: [Fiber v2](https://gofiber.io/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: PostgreSQL
- **Docs**: Swagger / OpenAPI

## 📋 Prerequisites

Ensure you have the following installed:
- [Go](https://golang.org/dl/) (1.20+)
- [Docker](https://www.docker.com/)
- [OpenSSL](https://www.openssl.org/) (for key generation)
- [Make](https://www.gnu.org/software/make/) (optional, for using Makefile)

## ⚡️ Getting Started

### 1. Database Setup

Start the PostgreSQL container:

```bash
make docker-init-db
# OR
docker run -p 15432:5432 --name sysware-test-postgres -e POSTGRES_PASSWORD=SyswareTestPassword -d postgres:16-alpine
```
### 2. Security Setup (JWT Keys)

Generate the ECDSA private and public keys for JWT signing:

```bash
# Generate Private Key
make privkey-generate
# OR
openssl ecparam -name prime256v1 -genkey -noout -out internal/assets/dev/jwt/privkey.pem

# Generate Public Key
make pubkey-generate
# OR
openssl ec -in internal/assets/dev/jwt/privkey.pem -pubout -out internal/assets/dev/jwt/pubkey.pem
```
### 3. Running the Application

You can run the application in development mode or build it for production.

**Development:**
```bash
make run
# OR
go run cmd/server/main.go
```

**Production Build:**
```bash
make build
./dist/server
```

## 📚 API Documentation

Once the server is running, you can access the Swagger API documentation at:

```
http://localhost:8080/swagger/index.html
```
*(Note: Port may vary based on your configuration)*

## 📂 Project Structure

```
.
├── cmd/                # Application entry points
├── configs/            # Configuration files (dev, uat, prd)
├── internal/           # Private application code
│   ├── handlers/       # HTTP handlers
│   ├── infrastructure/ # Core infrastructure (Server, Router)
│   └── datasources/    # Database and external connections
├── pkg/                # Public library code (Business Logic)
│   ├── auth/           # Authentication logic
│   ├── domain/         # Domain interfaces
│   ├── models/         # Data models
│   └── todo/           # Todo feature implementation
└── uploads/            # File uploads directory
```
