# Lucky Mart — Grocery Ordering Platform

A full-stack ordering platform I'm building for my family's grocery store
(*toko sembako*). The goal is to move day-to-day work — orders, stock, deliveries
and customers who buy on credit (*kasbon*) — out of notebooks and into one system.

> **Status:** in development, not deployed yet. Core backend modules work; the
> customer web store and admin dashboard are still being finished.

## What's inside

| Part | Path | Stack |
| --- | --- | --- |
| REST API | `cmd/`, `internal/`, `config/` | Go 1.22, Gin, GORM, PostgreSQL, JWT, Zerolog, Viper, Swagger |
| Customer web store | `customer-web/` | React, TypeScript, Vite, Tailwind CSS, TanStack Query, Zustand |
| Admin dashboard | `admin-dashboard/` | React, TypeScript, Vite, Tailwind CSS, TanStack Query/Table, Recharts |
| Local infrastructure | `docker-compose.yml`, `deploy/` | Docker Compose, Nginx, Prometheus/Grafana and k6 configs (work in progress) |

## Features

**Backend (85 REST endpoints across 10 route groups)**
- Product catalog with categories, brands and units
- Inventory and stock movements
- Cart and checkout
- Orders with status history (pending payment → paid → packing → delivery → completed)
- Shipping methods, rates and shipment tracking
- Payment methods and payment transactions
- Customer credit ledger (*kasbon*) for B2B / trusted customers
- Notifications with templates
- Auth with access and refresh tokens (JWT), role-based access for admin staff

**Customer web store:** browse products, cart, checkout, order list.
**Admin dashboard:** overview, products and categories, inventory and stock
movements, orders, payments and credit ledger, shipments, users, audit log.

## Architecture

```
customer-web ─┐
              ├─► REST API (Gin) ─► services ─► repositories (GORM) ─► PostgreSQL
admin-dashboard┘
```

The backend follows a layered structure: **handlers** (HTTP) → **services**
(business rules) → **repositories** (database).

## Running locally

Requirements: Go 1.22+, Node.js 18+, Docker Desktop.

```bash
# 1. Database (PostgreSQL in Docker, exposed on port 5433)
docker compose up -d postgres

# 2. Backend API — http://localhost:8080 (Swagger: /swagger/index.html)
cp .env.example .env.development   # then fill in your own values
go run cmd/api/main.go             # or: make run

# 3. Customer web store — http://localhost:3001
cd customer-web && npm install && npm run dev

# 4. Admin dashboard — http://localhost:3002
cd admin-dashboard && npm install && npm run dev
```

The seeder creates a development super-admin account and sample data on first
start (see `config/seed.go`). Set `SEED_ADMIN_PASSWORD` to choose the admin
password, and always change it outside local development.

Other commands: `make build`, `make test`, `make lint`, `make swagger`,
`make docker-up`, `make docker-down`.

## Configuration

All settings come from environment variables — see `.env.example`. Real `.env*`
files are git-ignored. `docker compose` reads `POSTGRES_PASSWORD` from a local
`.env` file (defaults to `postgres`).

## Roadmap

- Finish the customer checkout flow and admin pages
- Deploy (API + PostgreSQL + both frontends) so the store can use it
- Finish CI (GitHub Actions) and monitoring (Prometheus / Grafana)

## Author

**Lewi Lucky Siagian** — [GitHub](https://github.com/LuckySiagian) ·
[LinkedIn](https://www.linkedin.com/in/lewi-lucky-siagian-185573319)
