# Go-Hand

A robust backend service implemented in Go, following **Clean Architecture** principles. This project provides a secure foundation for user authentication and session management using PostgreSQL and Redis.

## Features

- **Clean Architecture**: Clear separation of concerns into API, Service, Repository, and Domain layers.
- **Authentication**: Secure JWT-based authentication (Access & Refresh Tokens).
- **Session Management**: Redis-backed session storage.
- **Database**: PostgreSQL with schema migrations via `golang-migrate`.
- **Security**: Password hashing using `bcrypt`.
- **Standard Library**: Built using Go's standard `net/http` `ServeMux` for routing.

## Tech Stack

- **Language**: Go
- **Database**: PostgreSQL (`lib/pq`)
- **Cache**: Redis (`go-redis/v9`)
- **Auth**: JWT (`golang-jwt/jwt/v5`)
- **Config**: `godotenv`

## Project Structure

```
├── cmd/api          # Main entry point (main.go)
├── internal         # Application logic
│   ├── api          # API definitions
│   ├── config       # Configuration loading
│   ├── model        # Domain models
│   ├── platform     # Infrastructure (DB, Cache, Logger)
│   ├── repository   # Data access layer
│   ├── service      # Business logic
│   ├── transport    # HTTP handlers and middleware
│   └── ...
├── migrations       # Database migration SQL files
├── Makefile         # Build and migration commands
└── go.mod           # Dependencies
```

## Setup & Installation

### Prerequisites
- Go 1.22+
- PostgreSQL
- Redis
- Make (optional, for migrations)
- `golang-migrate` (for database migrations)

### 1. Clone the Repository
```bash
git clone https://github.com/razedwell/go-hand.git
cd go-hand
```

### 2. Configure Environment
Copy the example environment file and update it with your credentials:
```bash
cp .env.example .env
```
Ensure your `.env` contains the correct database and Redis credentials.

### 3. Run Migrations
Initialize the database schema:
```bash
make migrate-up
```

### 4. Run the Application
Start the server:
```bash
go run cmd/api/main.go
```
The server will start on the port specified in your `.env` (default is usually `8080`).

## 🔌 API Endpoints

### Authentication
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/register` | Register a new user | ✗ |
| `POST` | `/login` | Login and receive tokens | ✗ |
| `POST` | `/refresh` | Refresh access token (requires cookie) | ✗ |
| `GET` | `/logout` | Invalidate session | ✓ |

### General
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `GET` | `/` | Protected home route | ✓ |
| `GET` | `/test` | Serves a test HTML page | ✗ |

## Testing

You can test the endpoints using the provided `/test` page or tools like curl/Postman.

```bash
# Example Health/Home Check (requires token)
curl -H "Authorization: Bearer <your_token>" http://localhost:8080/
```

## License

This project is licensed under the MIT License.
