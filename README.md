# Go Gallery - Modular Monolith

A high-performance, scalable web gallery backend built with Go, following the Modular Monolith architectural pattern. This project is designed to manage events and their associated photo galleries efficiently.

## 🚀 Technology Stack

- **Language:** [Go 1.25+](https://go.dev/)
- **Web Framework:** [Fiber v3](https://github.com/gofiber/fiber) (Next-gen Go web framework)
- **Database:** [PostgreSQL](https://www.postgresql.org/)
- **DB Driver:** [pgx/v5](https://github.com/jackc/pgx) (Pure Go PostgreSQL driver)
- **Migrations:** [goose](https://github.com/pressly/goose)
- **Validation:** [validator/v10](https://github.com/go-playground/validator)
- **API Documentation:** [Swagger (swaggo)](https://github.com/swaggo/swag)
- **Monitoring:** [Prometheus](https://prometheus.io/) (via Fiber middleware)
- **Authentication:** [UniAuth SSO Client](https://github.com/strbagus/uniauth) (Asymmetric RS256 signature verification, auto refresh token rotation, and fallback key resolution for kid-less JWTs)

## 📁 Project Structure

```text
├── db/
│   ├── migrations/     # SQL migration files managed by goose
│   └── seeders/        # DB seeders
├── internal/           # Private library code
│   ├── database/       # DB pool initialization
│   ├── helper/         # Shared utilities
│   └── middleware/     # Middlewares (Prometheus, AdminMiddleware)
├── module/             # Business modules
│   ├── event/          # Event management module (CRUD & Tokens)
│   │   ├── event_handler.go
│   │   ├── event_model.go
│   │   ├── event_repository.go
│   │   └── event_routes.go
│   └── photo/          # Photo management module
├── pkg/                # Public library code (request/response formatters)
├── router/             # Centralized route registration
├── uniauth-client/     # SSO Integration middleware client
├── main.go             # Application entry point
└── Makefile            # Task automation
```

## 🛠️ Getting Started

### Prerequisites

- Go 1.25 or higher
- PostgreSQL
- `goose` installed locally (optional, run via Makefile)
- `swag` installed locally (for generating API documentation)

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd gallery
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Setup environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your database and Auth service configurations
   ```

### Database Management

The project uses `goose` for migrations. Use the provided `Makefile` commands:

- **Apply all migrations:**
  ```bash
  make migrate-up
  ```
- **Rollback last migration:**
  ```bash
  make migrate-down
  ```
- **Fresh migration (Reset & Up):**
  ```bash
  make migrate-fresh
  ```
- **Check migration status:**
  ```bash
  make migrate-status
  ```
- **Create new migration:**
  ```bash
  make migrate-create name=migration_name
  ```

### Running the Application

To start the server:
```bash
make run
```
The server will start on the port defined in your `.env` (default: `8083`).

---

## 🔒 Authentication & Security

Administrative endpoints are secured via **UniAuth SSO Integration** using asymmetric key signature verification (RS256).

### Features
1. **Cookie-Based JWT Resolving**: The server parses `access_token` and `refresh_token` from cookies or fallback `Authorization: Bearer <token>` headers.
2. **Transparent Refresh Flow**: Expired access tokens are automatically refreshed with the identity provider.
3. **Fallback KID Resolution**: The client middleware can verify signature credentials using the JWKS keys map even if a token lacks a `kid` header.
4. **CORS Credentials**: Configured with `AllowCredentials: true` to support cross-origin cookie-based authorization.

---

## 📚 API Documentation

This project uses Swagger for API documentation.

### Accessing Swagger UI

When the server is running, you can access the Swagger UI at:
`http://localhost:8083/api/gallery/swagger/index.html`

*(Note: The actual URL may vary based on your `APP_PORT` and `APP_PATH` settings in `.env`)*

Swagger is configured with `withCredentials: true` and a cookie authentication scheme (`access_token`), permitting secure endpoint testing directly from the Swagger UI panel.

### Generating Documentation

Swagger documentation is automatically generated on build/run. To manually regenerate, run:
```bash
swag init
```

---

## 🛠️ Development Guidelines

- **New Modules:** Follow the modular structure in `module/`. Register routes in `router/routes.go`.
- **API Responses:** Always use the `pkg/response` package for consistent JSON formatting.
- **Database:** Use `database.PgxPool` and prefer explicit column scanning over `SELECT *`.

## 📜 License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
