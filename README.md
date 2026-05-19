# Go Gallery - Modular Monolith

A high-performance, scalable web gallery backend built with Go, following the Modular Monolith architectural pattern. This project is designed to manage events and their associated photo galleries efficiently.

## 🚀 Technology Stack

- **Language:** [Go 1.25+](https://go.dev/)
- **Web Framework:** [Fiber v3](https://github.com/gofiber/fiber) (Next-gen Go web framework)
- **Database:** [PostgreSQL](https://www.postgresql.org/)
- **DB Driver:** [pgx/v5](https://github.com/jackc/pgx) (Pure Go PostgreSQL driver)
- **Migrations:** [goose](https://github.com/pressly/goose)
- **Validation:** [validator/v10](https://github.com/go-playground/validator)
- **Monitoring:** [Prometheus](https://prometheus.io/) (via Fiber middleware)

## 📁 Project Structure

```text
├── internal/           # Private library code (database, middleware)
├── migrations/         # SQL migration files
├── module/             # Business modules (event, photo)
│   └── {module}/
│       ├── {name}_handler.go    # HTTP Logic
│       ├── {name}_model.go      # Domain Models
│       ├── {name}_repository.go # DB Logic
│       └── {name}_routes.go     # Route Registration
├── pkg/                # Public library code (request/response formatters)
├── router/             # Centralized route entry point
├── main.go             # Application entry point
└── Makefile            # Task automation
```

## 🛠️ Getting Started

### Prerequisites

- Go 1.25 or higher
- PostgreSQL
- `goose` installed locally (optional, can be run via Makefile)

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
   # Edit .env with your database credentials
   ```

### Database Management

The project uses `goose` for migrations. Use the provided `Makefile` commands for convenience:

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
The server will start on the port defined in your `.env` (default: `8055`).

## 🛠️ Development Guidelines

- **New Modules:** Follow the existing structure in `module/`. Register routes in `router/routes.go`.
- **API Responses:** Always use the `pkg/response` package for consistent JSON formatting.
- **Database:** Use `database.PgxPool` and prefer explicit column scanning over `SELECT *`.

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
