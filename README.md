![banner](https://i.imgur.com/8benIPj.png)

# BookArena

**BookArena** is a high-performance sports field booking management system built with Go. It provides a complete RESTful API and web interface for managing field rentals, bookings, payments, and user administration with enterprise-grade security features.

[![Live Demo](https://img.shields.io/badge/demo-live-success)](https://bookarena.1001api.com)
[![API Docs](https://img.shields.io/badge/API-docs-blue)](https://vifvpe310e.apidog.io)
[![Go Version](https://img.shields.io/badge/Go-1.24.0-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

---

## 🎯 Project Overview

BookArena is designed for sports facility managers, arena owners, and booking platforms. It enables users to browse available sports fields, create time-slot bookings, process payments, and manage facilities through a secure, role-based system.

**Live Application:** [https://bookarena.1001api.com](https://bookarena.1001api.com)

### Key Capabilities

- 🏟️ **Field Management** — CRUD operations for sports fields with pricing, location, and availability
- 📅 **Booking System** — Time-slot reservations with conflict detection and automatic pricing
- 💳 **Payment Processing** — Invoice generation and payment tracking (includes dummy payment endpoint for testing)
- 👥 **User Management** — Registration, authentication, and profile management
- 🔐 **Role-Based Access Control** — Three-tier authorization (User, Admin, SuperAdmin)
- 🔒 **Data Encryption** — AES-256-GCM encryption for sensitive user data (emails)
- 🌐 **Web Interface** — Server-side rendered dashboard using HTMX and Templ
- 📊 **Admin Dashboard** — Comprehensive management interface for administrators

---

## ✨ Features

### Authentication & Security
- JWT-based authentication with access and refresh tokens
- HTTPOnly encrypted cookies for token storage
- Bcrypt password hashing
- Email encryption (AES-256-GCM) and SHA-256 hashing
- Account lockout after failed login attempts
- Rate limiting (60 requests/minute per IP)
- CORS protection with configurable origins
- Secure cookie encryption

### Booking Management
- Availability checking
- Booking conflict detection
- Automatic price calculation based on duration
- User-specific booking history
- Field-specific booking views
- Admin booking oversight

### Multi-Role System
- **User** — Create bookings, view own bookings and payments
- **Admin** — Manage fields, view all bookings, manage users
- **SuperAdmin** — Full system access including admin creation

### API Features
- RESTful API design with versioning (`/api/v1`)
- Paginated list endpoints
- Comprehensive error handling
- Request duration tracking in responses
- Structured logging with Zerolog
- Health monitoring endpoint

---

## 🏗️ Architecture

BookArena follows **Clean Architecture** principles with clear separation of concerns:

```
bookarena/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── controllers/     # HTTP handlers (presentation layer)
│   ├── modules/         # Business logic (use cases)
│   │   ├── auth/
│   │   ├── bookings/
│   │   ├── fields/
│   │   ├── payments/
│   │   └── users/
│   ├── database/        # Data access layer
│   │   ├── migrations/  # SQL migrations
│   │   ├── queries/     # SQLC queries
│   │   └── generated/   # SQLC generated code
│   └── routes/          # Route definitions & middleware
├── views/               # Templ templates (SSR)
├── pkg/                 # Shared utilities
└── scripts/             # Build & deployment scripts
```

### Design Patterns
- **Repository Pattern** — Data access abstraction
- **Service Layer** — Business logic encapsulation
- **Dependency Injection** — Loose coupling between layers
- **Middleware Chain** — Authentication, authorization, logging, rate limiting

### Database Schema
- **users** — Encrypted email storage, role-based access, login tracking
- **fields** — Sports field details with geolocation support
- **bookings** — Time-slot reservations with foreign key constraints
- **payments** — Invoice and payment status tracking

---

## 🛠️ Tech Stack

### Backend
- **[Go 1.24](https://golang.org/)** — Primary language
- **[Fiber v2](https://gofiber.io/)** — Fast Express-inspired web framework
- **[PostgreSQL](https://www.postgresql.org/)** — Primary database
- **[SQLC](https://sqlc.dev/)** — Type-safe SQL code generation
- **[golang-migrate](https://github.com/golang-migrate/migrate)** — Database migrations
- **[Viper](https://github.com/spf13/viper)** — Configuration management

### Frontend
- **[Templ](https://templ.guide/)** — Type-safe Go templating
- **[HTMX](https://htmx.org/)** — Modern interactions without JavaScript frameworks

### Security & Authentication
- **[JWT (golang-jwt/jwt)](https://github.com/golang-jwt/jwt)** — Token-based authentication
- **[bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)** — Password hashing
- **AES-256-GCM** — Data encryption

### Logging & Monitoring
- **[Zerolog](https://github.com/rs/zerolog)** — High-performance structured logging
- **Fiber Monitor** — Built-in health monitoring

### Validation
- **[go-playground/validator](https://github.com/go-playground/validator)** — Request validation

---

## 🚀 Installation & Setup

### Prerequisites
- Go 1.24 or higher
- PostgreSQL 15+

### 1. Clone the Repository
```bash
git clone https://github.com/1001api/bookarena.git
cd bookarena
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Configure Environment
Copy the `.env` file and configure your settings:

```bash
cp .env.example .env
```

**Required Environment Variables:**
```env
# Server
PORT=8181

# Database
DATABASE_URL=postgres://admin:admin@localhost:5432/bookarena?sslmode=disable

# Application
APP_DOMAIN=localhost
CLIENT_DOMAIN=http://localhost:5173
ALLOWED_DOMAINS="['localhost', 'yourdomain.com']"

# Security Keys (CHANGE IN PRODUCTION!)
COOKIE_ENC_KEY=your-32-character-key-here
ENC_KEY=your-32-character-key-here
JWT_SECRET=your-jwt-secret-key-here

# JWT Configuration
JWT_EXPIRY=60              # Access token expiry in minutes
JWT_REFRESH_EXPIRY=4320    # Refresh token expiry in minutes (3 days)

# Root User (created automatically on first run)
ROOT_EMAIL=admin@yourdomain.com
ROOT_USERNAME=admin
ROOT_PASSWORD=SecurePassword123!
```

### 4. Database Setup

**Create Database:**
```bash
createdb bookarena
```

**Run Migrations:**
```bash
./scripts/migrate.sh up
# or manually:
migrate -path internal/database/migrations -database "$DATABASE_URL" up
```

### 5. Generate SQLC Code (if modified queries)
```bash
sqlc generate
```

### 6. Generate Templ Templates (if modified)
```bash
./scripts/templ.sh
# or:
templ generate
```

### 7. Run the Server

**Development Mode with Air (hot reload):**
```bash
air
```

**Production Mode:**
```bash
go build -o bin/bookarena ./cmd/server
./bin/bookarena
```

The server will start at `http://localhost:8181`

### 8. Access the Application

- **Web Interface:** http://localhost:8181
- **Health Check:** http://localhost:8181/health
- **API Base:** http://localhost:8181/api/v1

**Default Root User:**
- Email: `root@1001api.com` (or as configured)
- Username: `root`
- Password: `Pass123!@#` (or as configured)

---

## 📡 API Documentation

### Base URL
```
https://bookarena.1001api.com/api/v1
```

### Interactive API Documentation
- **ApiDog:** [https://vifvpe310e.apidog.io](https://vifvpe310e.apidog.io)
- **Postman Collection:** [Download ZIP](https://github.com/1001api/bookarena/releases/download/resources/BookArenaPostman.zip)
- **OpenAPI Spec:** [Download JSON](https://github.com/1001api/bookarena/releases/download/resources/BookArena.openapi.json)

### Core Endpoints

#### Authentication
| Endpoint | Method | Description | Access |
|----------|--------|-------------|--------|
| `/api/v1/auth/register` | POST | Register new user | Public |
| `/api/v1/auth/login` | POST | User login | Public |
| `/api/v1/auth/refresh` | POST | Refresh access token | Public |

#### Users
| Endpoint | Method | Description | Access |
|----------|--------|-------------|--------|
| `/api/v1/users/me` | GET | Get current user profile | Authenticated |
| `/api/v1/users/me` | PATCH | Update current user | Authenticated |
| `/api/v1/users/list` | GET | List all users | Admin |
| `/api/v1/users/search` | GET | Search users | Admin |
| `/api/v1/users/:id` | GET | Get user by ID | Admin |
| `/api/v1/users/:id` | PATCH | Update user | Admin |
| `/api/v1/users/:id` | DELETE | Delete user | Admin |
| `/api/v1/users/create` | POST | Create admin user | Admin |

#### Fields
| Endpoint | Method | Description | Access |
|----------|--------|-------------|--------|
| `/api/v1/fields/list` | GET | List all fields | Authenticated |
| `/api/v1/fields/:id` | GET | Get field details | Authenticated |
| `/api/v1/fields/create` | POST | Create new field | Admin |
| `/api/v1/fields/:id` | PATCH | Update field | Admin |
| `/api/v1/fields/:id` | DELETE | Delete field | Admin |

#### Bookings
| Endpoint | Method | Description | Access |
|----------|--------|-------------|--------|
| `/api/v1/bookings/create` | POST | Create booking | User |
| `/api/v1/bookings/me` | GET | Get user bookings | Authenticated |
| `/api/v1/bookings/field/:id` | GET | Get field bookings | Authenticated |
| `/api/v1/bookings/:id` | GET | Get booking details | Authenticated |
| `/api/v1/bookings/list` | GET | List all bookings | Admin |
| `/api/v1/bookings/:id` | DELETE | Delete booking | Admin |

#### Payments
| Endpoint | Method | Description | Access |
|----------|--------|-------------|--------|
| `/api/v1/payments/me` | GET | Get user invoices | Authenticated |
| `/api/v1/payments/:id` | GET | Get invoice details | Authenticated |
| `/api/v1/payments/list` | GET | List all invoices | Admin |
| `/api/v1/payments/:id` | DELETE | Delete invoice | Admin |
| `/api/v1/payments/pay` | POST | Pay invoice (dummy) | Authenticated |

### Query Parameters
- `limit` — Number of results per page (default: 10)
- `page` — Page number (default: 1)

---

## 🔧 Examples / Quickstart

### Register a New User
```bash
curl -X POST https://bookarena.1001api.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

**Response:**
```json
{
  "data": "550e8400-e29b-41d4-a716-446655440000",
  "meta": {
    "duration": "45.2ms"
  }
}
```

### Login
```bash
curl -X POST https://bookarena.1001api.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "SecurePass123!"
  }' \
  -c cookies.txt
```

**Response:**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  },
  "meta": {
    "duration": "120.5ms"
  }
}
```

### List Available Fields
```bash
curl -X GET "https://bookarena.1001api.com/api/v1/fields/list?limit=10&page=1" \
  -b cookies.txt
```

### Create a Booking
```bash
curl -X POST https://bookarena.1001api.com/api/v1/bookings/create \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "field_id": "550e8400-e29b-41d4-a716-446655440000",
    "start_time": "2024-10-28T14:00:00Z",
    "end_time": "2024-10-28T16:00:00Z"
  }'
```

### Get User Bookings
```bash
curl -X GET "https://bookarena.1001api.com/api/v1/bookings/me?limit=10&page=1" \
  -b cookies.txt
```

### Refresh Token
```bash
curl -X POST https://bookarena.1001api.com/api/v1/auth/refresh \
  -b cookies.txt \
  -c cookies.txt
```

---

## 🔐 Authentication

BookArena uses **JWT-based authentication** with a dual-token system for enhanced security.

### Token Types
1. **Access Token** (`__hk_asid`)
   - Short-lived (60 minutes default)
   - Used for API requests
   - Stored in HTTPOnly encrypted cookie

2. **Refresh Token** (`__hk_rsid`)
   - Long-lived (3 days default)
   - Used to obtain new access tokens
   - Stored in HTTPOnly encrypted cookie

### Authentication Flow
1. User logs in with username/password
2. Server validates credentials
3. Server generates access + refresh tokens
4. Tokens stored in HTTPOnly cookies
5. Client automatically sends cookies with requests
6. Server validates access token on protected routes
7. When access token expires, use refresh token to get new one

### Using Authentication in API Requests

**With cURL (cookies):**
```bash
# Login and save cookies
curl -X POST https://bookarena.1001api.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"pass"}' \
  -c cookies.txt

# Use cookies in subsequent requests
curl -X GET https://bookarena.1001api.com/api/v1/users/me \
  -b cookies.txt
```

**With JavaScript (Fetch API):**
```javascript
// Login
const response = await fetch('https://bookarena.1001api.com/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: 'user', password: 'pass' }),
  credentials: 'include' // Important: includes cookies
});

// Subsequent requests automatically include cookies
const userData = await fetch('https://bookarena.1001api.com/api/v1/users/me', {
  credentials: 'include'
});
```

### Role-Based Authorization

Endpoints are protected based on user roles:

| Role | Permissions |
|------|-------------|
| **user** | Create bookings, view own data |
| **admin** | Manage fields, view all bookings, manage users |
| **superadmin** | Full system access, create admins |

---

## 🧪 Development

### Project Structure
```
internal/
├── controllers/          # HTTP request handlers
│   ├── auth.go          # Authentication endpoints
│   ├── bookings.go      # Booking management
│   ├── fields.go        # Field management
│   ├── payments.go      # Payment processing
│   ├── users.go         # User management
│   └── web.go           # Web page rendering
├── modules/              # Business logic (services)
│   ├── auth/            # Authentication service
│   ├── bookings/        # Booking service
│   ├── fields/          # Field service
│   ├── payments/        # Payment service
│   └── users/           # User service
├── database/
│   ├── migrations/      # SQL migrations
│   ├── queries/         # SQLC query definitions
│   └── generated/       # Auto-generated SQLC code
└── routes/
    └── index.go         # Route definitions & middleware
```

### Adding a New Feature

1. **Define Database Schema** — Create migration in `internal/database/migrations/`
2. **Write SQL Queries** — Add queries to `internal/database/queries/`
3. **Generate Code** — Run `sqlc generate`
4. **Create Service** — Implement business logic in `internal/modules/`
5. **Create Controller** — Add HTTP handlers in `internal/controllers/`
6. **Register Routes** — Update `internal/routes/index.go`

### Code Generation

**SQLC (database code):**
```bash
sqlc generate
```

**Templ (templates):**
```bash
templ generate
# or
./scripts/templ.sh
```

### Database Migrations

**Create new migration:**
```bash
migrate create -ext sql -dir internal/database/migrations -seq your_migration_name
```

**Apply migrations:**
```bash
./scripts/migrate.sh up
```

**Rollback:**
```bash
./scripts/migrate.sh down
```