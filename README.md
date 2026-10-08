# EICAP-BANK API Server

REST API backend for **EICAP-BANK**, a bank/cash-desk management system. It handles
authentication, clients, accounts, denominations, cash sessions (opening/closing with
cash counts), bank operations and dashboard reports.

Built with **Go**, **Fiber v3**, **GORM**, and **PostgreSQL**.

---

## Table of Contents

- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Quick Start](#quick-start)
- [Configuration (.env)](#configuration-env)
- [Database Migrations](#database-migrations)
- [Seeders](#seeders)
- [Code Generation](#code-generation)
- [Running the Server](#running-the-server)
- [Taskfile Commands](#taskfile-commands)
- [API Reference](#api-reference)
- [Authentication](#authentication)

---

## Tech Stack

| Component    | Technology                                        |
| ------------ | ------------------------------------------------- |
| Language     | Go `1.26.5`                                       |
| HTTP Server  | [Fiber v3](https://docs.gofiber.io/)              |
| Database     | PostgreSQL                                        |
| ORM          | [GORM](https://gorm.io/)                          |
| Migrations   | [goose](https://github.com/pressly/goose/v3)      |
| Codegen      | [gorm.io/cli/gorm](https://gorm.io/)              |
| Auth / JWT   | `gofiber/contrib/jwt` + `golang-jwt/jwt/v5`       |
| Validation   | `go-playground/validator/v10`                     |
| Money        | `shopspring/decimal`                              |
| Config       | `.env` via `joho/godotenv`                        |
| Task Runner  | [go-task](https://taskfile.dev/) (optional)       |

---

## Prerequisites

- **Go 1.26 or higher** - see [Installing Go](#installing-go)
- **PostgreSQL** (running locally or remotely)
- **go-task** (optional, but recommended) - used to run the code generation task

### Installing Go

If you don't have Go installed:

1. Download the installer for your OS from <https://go.dev/dl/>.
   - On Windows, use the `.msi` archive; on macOS the `.pkg`; on Linux extract the `.tar.gz`.
2. Install it and make sure `go` is available on your `PATH`:

   ```bash
   go version
   # go version go1.26.5 windows/amd64 (example)
   ```

3. (Optional) Set up your module cache:

   ```bash
   go env -w GOFLAGS=-mod=mod
   ```

### Installing go-task (optional)

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

Verify with `task --version`.

### Installing the GORM CLI (required only for code generation)

```bash
go install gorm.io/cli/gorm@latest
```

Verify it is on your `PATH`:

```bash
gorm version
```

---

## Project Structure

```
server/
├── cmd/
│   └── api/
│       └── main.go                 # Application entrypoint
├── internal/
│   ├── boottrap/                   # Bootstrap: container, routes, validator, error handler
│   ├── config/                     # .env config loading
│   ├── constants/                  # Shared constants (e.g. well-known UUIDs)
│   ├── database/
│   │   ├── db_connection.go        # Postgres connection + migration runner
│   │   ├── migration/              # goose SQL migrations
│   │   └── seed/                   # Seeders (users, denominations, type accounts, type operations)
│   ├── enum/                       # Permissions (admin / student)
│   ├── generated/                  # GORM code generated from models (DO NOT EDIT)
│   ├── middleware/                 # JWT + role-based permission middlewares
│   ├── model/                      # GORM models
│   ├── module/                     # Feature modules (handler/service/repo/dto per module)
│   │   ├── account/
│   │   ├── auth/
│   │   ├── bank_operation/
│   │   ├── cash_count/
│   │   ├── cash_session/
│   │   ├── client/
│   │   ├── dashboard/
│   │   ├── denomination/
│   │   ├── type_account/
│   │   ├── type_operation/
│   │   └── user/
│   └── response/                   # Uniform API responses + error types
├── pkg/                            # Shared utilities (hash, jwt, cookie, pagination)
├── .env.example                    # Environment template
├── Taskfile.yml                    # Task definitions (gorm gen)
├── go.mod / go.sum
└── README.md
```

### Module conventions

Every feature module in `internal/module/<name>/` follows the same pattern:

- `handler.go` - Fiber handlers + route registration (`RegisterRoutes`)
- `service.go` - Business logic
- `repo.go` - GORM data access
- `dto.go` - Request/response DTOs and validation tags

---

## Quick Start

1. **Clone and enter the project**

   ```bash
   git clone <repository-url> eicap-bank
   cd eicap-bank/server
   ```

2. **Install Go dependencies**

   ```bash
   go mod download
   ```

3. **Create your environment file**

   ```bash
   cp .env.example .env
   ```

   Edit `.env` and fill in at least:

   ```bash
   DATABASE_URL=postgres://postgres:postgres@localhost:5432/eicap_bank?sslmode=disable
   JWT_SECRET=change_me_to_a_long_random_string
   ADMIN_PASSWORD_ONE=your_admin_password
   ```

4. **Create the database** (once)

   ```bash
   createdb eicap_bank
   # or with psql
   psql -U postgres -c "CREATE DATABASE eicap_bank;"
   ```

5. **Enable migrations and seeders on first run**

   ```bash
   RUN_MIGRATION=true
   RUN_SEEDER=true
   ```

6. **Run the server**

   ```bash
   go run ./cmd/api
   ```

   The API is now available at <http://localhost:8000> (base path `/api/v1`).

7. **Login with the seeded admin user**

   ```text
   Email:    admin1@eicap.com
   Password: <value of ADMIN_PASSWORD_ONE>
   ```

> The server itself applies migrations and seeders on startup **before** serving requests,
> so no extra CLI steps are required beyond setting the env flags.

---

## Configuration (.env)

Copy `.env.example` to `.env` and adjust the values.

| Variable                | Required | Default                 | Description                                                                 |
| ----------------------- | -------- | ----------------------- | --------------------------------------------------------------------------- |
| `APP_PORT`              | No       | `8000`                  | Port the HTTP server listens on                                             |
| `APP_ENV`               | No       | `development`           | Use `production` to enable secure cookies (`Secure` flag)                   |
| `DATABASE_URL`          | **Yes**  | —                       | PostgreSQL connection string, e.g. `postgres://user:pass@host:5432/db?sslmode=disable` |
| `RUN_MIGRATION`         | No       | `false`                 | Set to `true` to run goose migrations on startup                           |
| `RUN_SEEDER`            | No       | `false`                 | Set to `true` to run seeders on startup                                    |
| `ALLOW_ORIGINS`         | No       | `http://localhost:5173` | Comma/space separated CORS allowed origins (frontend URL)                  |
| `JWT_SECRET`            | **Yes**  | —                       | Secret used to sign JWT tokens (use a long random string)                  |
| `JWT_EXPIRATION`        | No       | `24h`                   | Token lifetime, any Go duration (e.g. `30m`, `24h`, `7d`)                  |
| `ADMIN_PASSWORD_ONE`    | *¹        | —                       | Password for the seeded admin user (`admin1@eicap.com`)                    |

¹ Only required if you enable the seeder (`RUN_SEEDER=true`). If the seed runs with an
empty password, that user is skipped.

---

## Database Migrations

Migrations are managed with **[goose](https://github.com/pressly/goose)** and stored as
`.sql` files in `internal/database/migration/`.

They are **automatically applied on application startup** when
`RUN_MIGRATION=true` is set in `.env`:

```bash
RUN_MIGRATION=true
go run ./cmd/api
```

The migration runner is wired in `internal/database/db_connection.go:48`

### Running migrations manually with the goose CLI (optional)

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
cd internal/database/migration
goose postgres "$DATABASE_URL" up       # apply
goose postgres "$DATABASE_URL" down     # rollback one step
goose postgres "$DATABASE_URL" status   # show status
```

### Creating a new migration

```bash
cd internal/database/migration
goose create add_something sql
```

Then edit the generated file and keep the `-- +goose Up` / `-- +goose Down` markers:

```sql
-- +goose Up
CREATE TABLE example (...);

-- +goose Down
DROP TABLE IF EXISTS example;
```

> The initial migration requires the `uuid-ossp` extension. If you use the automatic
> runner it is created for you; when applying manually ensure the extension exists.

---

## Seeders

Seeders run automatically on startup when `RUN_SEEDER=true`:

```bash
RUN_SEEDER=true
go run ./cmd/api
```

What gets seeded (all idempotent — safe to re-run):

| Seeder              | Contents                                                                  |
| ------------------- | ------------------------------------------------------------------------- |
| `users`             | One admin user `admin1@eicap.com` with role `admin` (password = `ADMIN_PASSWORD_ONE`) |
| `denominations`     | 8 bill/coin denominations (200, 100, 50, 20, 10, 5, 2, 1)                |
| `type_accounts`     | `Caja de Ahorro`, `Cuenta Corriente`, `Caja DPF`                          |
| `type_operations`   | `ING` (Ingreso), `EGR` (Egreso), `APC` (Apertura de Cuenta), `APCA` (Apertura de Caja), `CICA` (Cierre de Caja) |

---

## Code Generation

The package `internal/generated/` contains GORM column constants generated from the
models in `internal/model/`. These files are marked *DO NOT EDIT* and are regenerated
whenever a model changes.

### Generating

Install the GORM CLI if you haven't already:

```bash
go install gorm.io/cli/gorm@latest
```

Then, from the `server/` directory run:

```bash
# Via Taskfile (recommended)
task gen
```

or equivalently, the command it wraps:

```bash
gorm gen -i ./internal/model -o ./internal/generated
```

This scans the models in `-i ./internal/model` and writes generated types to
`-o ./internal/generated`.

> Always re-run `task gen` after adding or modifying a struct in `internal/model/`
> and `go build ./...` to get a fresh compile check.

### Regenerating after model changes

```bash
task gen && go build ./...
```

---

## Running the Server

### Development (hot reload with Air)

Install [Air](https://github.com/air-verse/air):

```bash
go install github.com/air-verse/air@latest
```

Then run:

```bash
air
```

### Run directly

```bash
go run ./cmd/api
```

### Build a binary

```bash
go build -o bin/eicap-bank-server ./cmd/api
./bin/eicap-bank-server
```

### Verify the build without running

```bash
go vet ./... && go build ./...
```

### Health check

```bash
# Root endpoint (no auth)
curl http://localhost:8000/
# Welcome to EICAP-BANK API
```

---

## Taskfile Commands

The project includes a [`Taskfile.yml`](./Taskfile.yml) (go-task). Available tasks:

```bash
task --list
```

| Task  | Description                | Command                            |
| ----- | -------------------------- | ---------------------------------- |
| `gen` | Regenerate GORM code       | `gorm gen -i ./internal/model -o ./internal/generated` |

---

## API Reference

All endpoints are under the base path **`/api/v1`**. JSON request/response bodies.
Every endpoint except `POST /auth/login` and `POST /auth/logout` requires a valid JWT
cookie (see [Authentication](#authentication)).

Specific mutation routes are restricted to the `admin` role (marked **🔒 Admin**).

### Health Check

| Method | Path              | Auth     | Description          |
| ------ | ----------------- | -------- | -------------------- |
| GET    | `/`               | No       | Welcome message      |

### Auth

| Method | Path              | Auth     | Description          |
| ------ | ----------------- | -------- | -------------------- |
| POST   | `/auth/login`     | No       | Log in, sets `token` cookie |
| POST   | `/auth/logout`    | No       | Clears auth cookie   |

Login request body:

```json
{
  "email": "admin1@eicap.com",
  "password": "your_password"
}
```

(`username` is accepted as an alternative to `email`.)

### Users

| Method | Path                | Auth          | Description                  |
| ------ | ------------------- | ------------- | ---------------------------- |
| POST   | `/users`            | 🔒 Admin      | Create user                 |
| GET    | `/users/me`         | JWT           | Current user's profile      |
| GET    | `/users/:id`        | 🔒 Admin      | Find user by ID             |
| GET    | `/users`            | 🔒 Admin      | List users (paginated)      |
| PUT    | `/users/:id`        | 🔒 Admin      | Update user                 |
| DELETE | `/users/:id`        | 🔒 Admin      | Delete user (not yourself)  |

### Clients

| Method | Path                     | Auth     | Description                          |
| ------ | ------------------------ | -------- | ------------------------------------ |
| POST   | `/clients`               | JWT      | Create client                        |
| GET    | `/clients/mine`          | JWT      | Current user's clients (paginated)   |
| PUT    | `/clients/:id`           | JWT      | Update client                        |
| GET    | `/clients/:id`           | JWT      | Find client by ID                    |
| GET    | `/clients`               | JWT      | List clients                         |
| GET    | `/clients/search/:data`  | JWT      | Search client by C.I. / data         |
| DELETE | `/clients/:id`           | JWT      | Delete client                        |

### Type Accounts

| Method | Path                    | Auth          | Description              |
| ------ | ----------------------- | ------------- | ------------------------ |
| POST   | `/type-accounts`        | 🔒 Admin      | Create type account      |
| PUT    | `/type-accounts/:id`    | 🔒 Admin      | Update type account      |
| GET    | `/type-accounts/:id`    | JWT           | Find by ID               |
| GET    | `/type-accounts`        | JWT           | List type accounts       |
| DELETE | `/type-accounts/:id`    | 🔒 Admin      | Delete type account      |

### Accounts

| Method | Path                          | Auth     | Description              |
| ------ | ----------------------------- | -------- | ------------------------ |
| POST   | `/accounts`                   | JWT      | Create account           |
| PUT    | `/accounts/:id`               | JWT      | Update account           |
| GET    | `/accounts/client/:clientId`  | JWT      | Accounts of a client     |
| GET    | `/accounts/:id`               | JWT      | Find account by ID       |
| GET    | `/accounts`                   | JWT      | List accounts            |
| DELETE | `/accounts/:id`               | JWT      | Delete account           |

### Type Operations

| Method | Path                          | Auth          | Description              |
| ------ | ----------------------------- | ------------- | ------------------------ |
| POST   | `/type-operations`            | 🔒 Admin      | Create type operation    |
| PUT    | `/type-operations/:id`        | 🔒 Admin      | Update type operation    |
| GET    | `/type-operations/:id`        | JWT           | Find by ID               |
| GET    | `/type-operations`            | JWT           | List type operations     |
| DELETE | `/type-operations/:id`        | 🔒 Admin      | Delete type operation    |

### Denominations

| Method | Path                    | Auth          | Description              |
| ------ | ----------------------- | ------------- | ------------------------ |
| POST   | `/denominations`        | 🔒 Admin      | Create denomination      |
| PUT    | `/denominations/:id`    | 🔒 Admin      | Update denomination      |
| GET    | `/denominations/:id`    | JWT           | Find by ID               |
| GET    | `/denominations`        | JWT           | List denominations       |
| DELETE | `/denominations/:id`    | 🔒 Admin      | Delete denomination      |

### Cash Sessions

| Method | Path                      | Auth     | Description                     |
| ------ | ------------------------- | -------- | ------------------------------- |
| POST   | `/cash-sessions/open`     | JWT      | Open a cash session with a cash count |
| PUT    | `/cash-sessions/:id/close`| JWT      | Close a cash session (own session only) |
| GET    | `/cash-sessions/mine/open`| JWT      | Current user's open session     |
| GET    | `/cash-sessions/:id`      | JWT      | Find session by ID              |
| GET    | `/cash-sessions`          | JWT      | List sessions (paginated)       |

Opening/closing a session automatically records `APCA` / `CICA` bank operations and
computes expected vs counted amounts (difference).

### Cash Counts

| Method | Path                               | Auth     | Description                      |
| ------ | ---------------------------------- | -------- | -------------------------------- |
| GET    | `/cash-counts/session/:sessionId`  | JWT      | Cash counts of a session         |

### Bank Operations

| Method | Path                                | Auth     | Description                          |
| ------ | ----------------------------------- | -------- | ------------------------------------ |
| POST   | `/bank-operations`                  | JWT      | Register a bank operation            |
| GET    | `/bank-operations/active-session`   | JWT      | Operations of the active session     |
| GET    | `/bank-operations/mine`             | JWT      | Current user's operations            |
| GET    | `/bank-operations/client/:clientId` | JWT      | Operations of a client               |
| GET    | `/bank-operations/user/:userId`     | JWT      | Operations of a user                 |
| GET    | `/bank-operations/:id`              | JWT      | Find operation by ID                 |
| GET    | `/bank-operations`                  | JWT      | List operations (paginated)          |

### Dashboard

| Method | Path                | Auth     | Description              |
| ------ | ------------------- | -------- | ------------------------ |
| GET    | `/dashboard/summary`| JWT      | Dashboard summary report (role-aware) |

---

## Authentication

Authentication is **cookie-based**. Logging in via `POST /auth/login` sets an HTTP-only
cookie named `token` containing a signed JWT.

- The application uses a **same-site none** cookie and, when `APP_ENV=production`, the
  `Secure` flag is enabled (HTTPS only).
- The JWT payload contains the user's `user_id` and `role` (see `pkg/jwt.go`).
- Role-based authorization uses the `middleware.RequirePermission(...)` middleware
  (`internal/middleware/permission.go`). Available roles (`internal/enum/permission.go`):
  - `admin`
  - `student`

Example login with `curl` (store the cookie jar):

```bash
curl -c cookies.txt -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin1@eicap.com","password":"your_password"}'
```

Then call protected endpoints with the cookie:

```bash
curl -b cookies.txt http://localhost:8000/api/v1/users/me
```

Logout clears the cookie:

```bash
curl -b cookies.txt -X POST http://localhost:8000/api/v1/auth/logout
```

---

## Common Issues & Tips

- **`DATABASE_URL not set`** - you forgot to create `.env`; copy `.env.example`.
- **Connection retries logged** - the server retries the DB connection 10 times (2s
  apart) before failing; make sure PostgreSQL is up.
- **Migrations not applied** - set `RUN_MIGRATION=true` and restart. The runner resolves
  the migration dir relative to the working directory, so run the server from the
  `server/` folder.
- **`gorm` command not found** - the GORM CLI is not on your `PATH`; run
  `go install gorm.io/cli/gorm@latest` and ensure your `GOBIN`/`GOPATH` bin is set, or
  invoke the codegen with `go run gorm.io/cli/gorm gen ...`.
- **CORS errors from the frontend** - update `ALLOW_ORIGINS` to include your frontend
  origin (default `http://localhost:5173` for Vite).