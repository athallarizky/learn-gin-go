# Learn Gin Go — RESTful API & Architecture Lab

A robust, modular RESTful API service built with **Go**, **Gin Web Framework**, **GORM v2**, and **PostgreSQL 16**. This project serves as a hands-on architectural playground demonstrating production-grade patterns: Social OAuth2 authentication, stateless JWT issuance, Role-Based Access Control (RBAC), relational data modeling with cascade constraints, and idiomatic Go memory management.

---

## 🚀 Key Features

* **Authentication & Authorization:**
  * **OAuth2 Social Login:** Google & GitHub authentication powered by [`markbates/goth`](https://github.com/markbates/goth), with automatic user creation and email-based identity linking.
  * **Stateless JWT Tokens:** Custom HMAC-SHA256 (`HS256`) tokens issuing `user_id` and `user_role` claims via [`golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt/v5).
  * **Chained Middleware Pipeline:** Modular Bearer token authentication and Role-Based Access Control (`IsAdmin` middleware returning `403 Forbidden`).
* **Relational Data Modeling (PostgreSQL + GORM):**
  * Auto-migration of schemas for `users` and `articles`.
  * One-to-Many relational mapping with foreign key integrity and `CASCADE` updates/deletions.
  * Eager loading relationships using GORM `Preload("Articles")` to eliminate $N+1$ query overhead.
* **Content Management & Query Operations:**
  * Automatic URL slugification using [`gosimple/slug`](https://github.com/gosimple/slug) with collision avoidance via Unix timestamp suffixes.
  * Ownership verification: users can only update or mutate their own articles.
  * Partial updates (`PATCH`) with GORM model mutations.
  * Tag-based filtering via SQL `LIKE` pattern matching.
* **Developer Experience (DX):**
  * Containerized PostgreSQL 16 Alpine via Docker Compose (mapped to port `5433`).
  * Live-reloading during development via [Air](https://github.com/air-verse/air).
  * Centralized configuration loading using `.env` and [`joho/godotenv`](https://github.com/joho/godotenv).

---

## 🛠️ Tech Stack

| Layer | Technology |
|---|---|
| **Language** | [Go (1.22+)](https://go.dev/) |
| **HTTP Framework** | [Gin Web Framework](https://github.com/gin-gonic/gin) (`v1.10`) |
| **ORM** | [GORM v2](https://gorm.io/) with PostgreSQL Driver |
| **Database** | [PostgreSQL 16](https://www.postgresql.org/) (Docker Compose) |
| **Authentication** | [Goth](https://github.com/markbates/goth) (OAuth2) & [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) |
| **Dev Tooling** | [Air](https://github.com/air-verse/air) (Live Reloading), Docker / Colima |

---

## 📁 Project Structure

```text
learn-gin-go/
├── config/                  # Application configuration & initialization
│   ├── auth.go              # Goth OAuth2 provider setup (Google, GitHub)
│   ├── config.go            # Centralized environment variable loader
│   ├── db.go                # GORM database connection & schema auto-migration
│   └── jwt.go               # JWT generation, signing, and claim verification
├── middlewares/             # Request interceptors & guards
│   └── auth.go              # Bearer token validator & IsAdmin RBAC guard
├── models/                  # GORM entity models & database schemas
│   ├── article.go           # Article struct & relational foreign keys
│   └── user.go              # User struct & social profile fields
├── routes/                  # Controller handlers & business logic
│   ├── article.go           # Article CRUD handlers & tag search
│   ├── auth.go              # OAuth begin, callback, user registration, logout
│   └── user.go              # User profile & preloaded articles handler
├── docs/                    # Architectural findings (synced to engineering-handbook)
├── docker-compose.yml       # Local PostgreSQL service definition
├── main.go                  # Application entry point
├── server.go                # Gin router tree & route group configuration
├── .env.example             # Template for required environment variables
└── go.mod                   # Module dependencies
```

---

## 🏁 Getting Started

### 1. Prerequisites

* **Go**: Version `1.22` or later installed locally.
* **Docker & Docker Compose** (or Colima on macOS).
* **OAuth Credentials**: Google and/or GitHub OAuth App credentials (client ID & secret).

### 2. Environment Setup

Clone the repository and copy the environment configuration template:

```bash
git clone https://github.com/athallarizky/learn-gin-go.git
cd learn-gin-go
cp .env.example .env
```

Configure your `.env` values:

```env
PORT=8088

# PostgreSQL Database
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=learn_gin_db

# Security & Secrets
JWT_SECRET=your_super_secret_jwt_key_here

# Social OAuth (Google)
GOOGLE_KEY=your_google_client_id
GOOGLE_SECRET=your_google_client_secret

# Social OAuth (GitHub)
GITHUB_KEY=your_github_client_id
GITHUB_SECRET=your_github_client_secret
```

### 3. Start Database Container

Start the PostgreSQL 16 container in the background:

```bash
docker compose up -d
```

Verify that the database is accepting connections on port `5433`:

```bash
docker ps --filter "name=learn_gin_postgres"
```

### 4. Run the Application

#### Standard Run:
```bash
go run .
```

#### With Air Live-Reload (Recommended for Development):
```bash
air
```

The server will initialize database connections, run auto-migrations, and listen on `http://localhost:8088`.

---

## 📡 API Reference

All primary endpoints are nested under the `/api/v1/` prefix.

### 1. Public Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/ping` | Healthcheck endpoint (`{"message": "pong"}`) |
| `GET` | `/api/v1/articles/` | Retrieve all published articles |
| `GET` | `/api/v1/articles/:slug` | Retrieve a single article by its unique slug |
| `GET` | `/api/v1/articles/tag/:tag` | Filter articles matching a specific tag |

### 2. Authentication Routes (`/api/v1/auth`)

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/auth/:provider` | Initiate OAuth flow (`provider`: `google` or `github`) |
| `GET` | `/api/v1/auth/:provider/callback` | OAuth callback; auto-registers user and returns JWT token |
| `GET` | `/api/v1/auth/logout` | Clear OAuth session state |

### 3. Authenticated User Routes (`/api/v1/user` & `/api/v1/articles`)

*Requires header: `Authorization: Bearer <token>`*

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/user/articles` | Get authenticated user's profile with preloaded articles |
| `POST` | `/api/v1/articles/` | Create a new article associated with authenticated user |
| `PATCH` | `/api/v1/articles/:slug` | Partially update an article *(ownership verified)* |

### 4. Admin Routes (`/api/v1/admin`)

*Requires header: `Authorization: Bearer <token>` and `user_role: true`*

| Method | Endpoint | Description |
|---|---|---|
| `DELETE` | `/api/v1/admin/article/:slug` | Delete any article by slug (Admin privilege) |

---

## 🧠 Architectural Insights & Learning Highlights

During the development of this repository, in-depth architectural findings and design decisions were analyzed and documented:

1. **Go Pointers & Memory Semantics:** Understanding why GORM requires `&items` (address-of) rather than returning slices, reflection-driven struct hydration, and the difference between `var items []T` (`null` JSON) vs `items := []T{}` (`[]` JSON).
2. **Gin Context Lifecycle & Panic Recovery:** Flow contracts of `c.MustGet` vs defensive `c.Get`, type assertions `.(T)`, and how `gin.Recovery()` catches panics to return HTTP 500 without crashing the server process.
3. **HTTP Semantics & Bandwidth Efficiency:** Why `PATCH` with GORM `.Updates()` optimizes network payloads compared to strict `PUT` full-resource replacement.
4. **The Zero-Value Trap in GORM:** How GORM `.Updates(struct)` omits zero-value fields (`""`, `0`, `false`) by default, and how to safely override this behavior using maps or `.Select()`.
5. **Classless OOP & Constructor Patterns:** Why Go rejects inheritance hierarchies in favor of *Composition Over Inheritance*, struct embedding (`gorm.Model`), and factory constructor functions (`New...() *T`).

> Detailed documentation records are maintained in the central knowledge base under [`docs/`](docs/) and synced via portable relative symlinks.

---

## 📄 License

This repository is built for learning and demonstration purposes under the MIT License.
