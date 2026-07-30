# coffee-site-server

Go backend API for the [coffee-site](../coffee-site) Vue storefront. SQLite storage (pure-Go driver, no CGO needed), JWT auth, and a REST API covering products, cart, wishlist, orders, and account flows.

## Run

```sh
go run ./cmd/server
```

Server listens on `:8080` by default. Config is read from environment variables (see `.env.example`); a `coffee-site.db` SQLite file is created and seeded automatically on first run.

Seeded accounts:
- `admin@coffeehouse.example.com` / `admin123` (admin)
- `member@coffeehouse.example.com` / `member123` (member)

## Config

| Env var | Default | Purpose |
| --- | --- | --- |
| `PORT` | `8080` | HTTP port |
| `DATABASE_PATH` | `coffee-site.db` | SQLite file path |
| `JWT_SECRET` | `dev-secret-change-me` | JWT signing secret — set a real value outside dev |
| `FRONTEND_ORIGIN` | `http://localhost:5173` | Allowed CORS origin for the Vite dev server |

## Email is stubbed

Registration and password-reset verification codes are **not** emailed. They're logged to the server console (`[stub-email] ...`) and also returned in the API response body as `devCode`, so the frontend can surface them during development. Swap in a real mail provider in `internal/handlers/auth.go` before shipping to production.

## API overview

All routes are under `/api`. Public GETs for `beans`, `utensils`, and their category lists; admin-only (`Authorization: Bearer <token>` with an admin account) for create/update/delete/category management. Cart and wishlist are scoped by an `X-Guest-Id` header (any client-generated stable ID — no login required). Orders and checkout require a logged-in user.

See `internal/handlers/router.go` for the full route table.
