# ![banner](https://i.imgur.com/8benIPj.png)

# BookArena

**BookArena** features a RESTful API, server-side rendered web interface with Alpine.js + HTMX, real-time booking conflict detection, automated payment invoice generation, and enterprise-grade security with AES-256-GCM encryption.

<div align="center">

**[Live Demo](https://bookarena.1001api.com)** •
**[API Documentation](https://vifvpe310e.apidog.io)** •
**[Postman Collection](https://github.com/1001api/bookarena/releases/download/resources/BookArenaPostman.zip)** •
**[OpenAPI Spec](https://github.com/1001api/bookarena/releases/download/resources/BookArena.openapi.json)**

[Features](#-features) •
[Quick Start](#-quick-start) •
[Documentation](#-complete-api-reference) •
[Architecture](#-architecture--design) •
[Demo Accounts](#-demo-accounts)

</div>

---

## 📖 Project Overview

BookArena **built with Go 1.24 and PostgreSQL**, it delivers good performance, security, and scalability while maintaining clean, maintainable code architecture.

>
> **🌐 Live Application:** [https://bookarena.1001api.com](https://bookarena.1001api.com)
>
> **📚 Live API Documentation:** [https://vifvpe310e.apidog.io](https://vifvpe310e.apidog.io)
>
> **📝 Postman Collection:** [See Resources](https://github.com/1001api/bookarena/releases/download/resources/BookArenaPostman.zip)
>
> **📄 OpenAPI Spec:** [See Resources](https://github.com/1001api/bookarena/releases/download/resources/BookArena.openapi.json)
>

### Use Cases

**For Venue Owners:**
- Manage multiple sports fields/courts with individual pricing
- Track all bookings in real-time with visual dashboard
- Monitor revenue through automated payment invoices
- Control field availability and operating hours
- Access detailed analytics on field utilization

**For Customers:**
- Browse available fields with images, descriptions, and pricing
- Book time slots instantly with real-time availability checking
- View booking history and payment status
- Receive automatic payment invoices
- Manage profile and booking preferences

**For Administrators:**
- Full system oversight with multi-level access control
- User management (create, update, delete, search)
- Field management with geolocation support
- Booking oversight and conflict resolution
- Payment tracking and invoice management

### Key Capabilities

#### 🏟️ **Sports Field Management**
- Complete CRUD operations for field resources
- Support for multiple field types (soccer, basketball, futsal, etc.)
- Per-hour pricing with automatic calculation
- Image upload support for field visualization
- Location tracking with latitude/longitude coordinates
- Description and metadata management
- Soft deletion for data integrity

#### 📅 **Advanced Booking System**
- Real-time time-slot availability checking
- Intelligent booking conflict detection
- Automatic price calculation based on duration
- Timezone-aware booking (Asia/Jakarta)
- User-specific booking history with pagination
- Field-specific booking views for availability planning
- Booking status tracking (linked to payment status)
- Admin-level booking management and cancellation

#### 💳 **Payment & Invoicing System**
- Automatic invoice generation on booking creation
- Payment status tracking (pending, paid, failed)
- User-specific invoice history
- Admin invoice management dashboard
- Dummy payment endpoint for testing (simulates payment gateway)
- Payment-booking relationship tracking
- Total price calculation and validation

#### 👥 **Enterprise User Management**
- Secure user registration with email validation
- Three-tier role-based access control:
  - **User**: Standard booking and payment access
  - **Admin**: Field and user management capabilities
  - **SuperAdmin**: Full system access including admin creation
- Email encryption (AES-256-GCM) with SHA-256 hashing
- Profile management (update username, email, password)
- User search functionality for administrators
- Account status tracking (deleted_at for soft deletes)
- Login tracking and failed attempt monitoring
- Account lockout mechanism for security

#### 🔐 **Enterprise-Grade Security**
- JWT-based dual token authentication (access + refresh)
- Encrypted HTTPOnly cookies for token storage
- Bcrypt password hashing with salt
- AES-256-GCM encryption for PII (email addresses)
- SHA-256 hashing for email lookups
- Rate limiting (60 requests/minute per IP)
- CORS protection with configurable origins
- Cookie encryption middleware
- SQL injection protection via prepared statements (SQLC)
- Input validation on all endpoints
- Secure token expiry and rotation

#### 🌐 **Modern Web Interface**
- Server-side rendering with Templ templates
- Reactive UI using Alpine.js (lightweight, no build step)
- HTMX for seamless partial updates
- TailwindCSS for responsive design
- Dark mode support
- Real-time field browsing dashboard
- Interactive booking form with date/time pickers
- Live booking status display
- Authentication flows (login/register pages)
- Mobile-responsive design

#### 🚀 **API Features**
- RESTful API design with versioning (`/api/v1`)
- Consistent JSON response format
- Request duration tracking in metadata
- Paginated list endpoints (limit + offset)
- Comprehensive error handling with HTTP status codes
- Structured logging with Zerolog
- Health monitoring endpoint (`/health`)
- Request/response logging for debugging
- OpenAPI documentation available

---

## ✨ Features

### Authentication & Security

<details>
<summary><b>Click to expand security features</b></summary>

**Token-Based Authentication:**
- **Dual JWT Strategy**: Separate access tokens (60min) and refresh tokens (3 days)
- **Token Storage**: Encrypted HTTPOnly cookies (`__hk_asid`, `__hk_rsid`)
- **Token Rotation**: Automatic refresh mechanism prevents expired token issues
- **Token Validation**: Middleware validates tokens on every protected route
- **User Existence Check**: Validates that token user still exists in database

**Encryption & Hashing:**
- **Password Security**: Bcrypt hashing with automatic salt generation
- **Email Encryption**: AES-256-GCM encryption for storing email addresses
- **Email Lookup**: SHA-256 hashing for fast, secure email queries
- **Cookie Encryption**: Fiber's encryptcookie middleware protects all cookies
- **Encryption Key Management**: Configurable via environment variables

**Access Control:**
- **Role-Based Authorization**: Three-tier system (user/admin/superadmin)
- **Middleware Protection**: Route-level role enforcement
- **Context-Based Auth**: User info injected into request context
- **Forbidden Responses**: Clear 403 errors for insufficient permissions

**Attack Prevention:**
- **Rate Limiting**: 60 requests/minute per IP with sliding window
- **CORS Protection**: Configurable allowed origins
- **SQL Injection Prevention**: All queries use prepared statements via SQLC
- **Account Lockout**: Failed login attempt tracking (future feature)
- **Input Validation**: Schema validation on all POST/PATCH requests

</details>

### Booking Management

<details>
<summary><b>Click to expand booking features</b></summary>

**Conflict Detection:**
- Prevents double-booking of the same field at overlapping times
- Checks for existing reservations before confirming new bookings
- Returns HTTP 409 Conflict with descriptive error messages

**Automatic Pricing:**
- Calculates total cost based on time duration (end_time - start_time)
- Uses field's `price_per_hour` rate
- Supports hourly, half-hourly, or custom duration bookings
- Price displayed in Indonesian Rupiah (Rp) format

**Booking Views:**
- **User View** (`/api/v1/bookings/me`): Personal booking history
- **Field View** (`/api/v1/bookings/field/:id`): All bookings for specific field
- **Admin View** (`/api/v1/bookings/list`): System-wide booking oversight
- **Detail View** (`/api/v1/bookings/:id`): Individual booking with all details

**Booking Lifecycle:**
1. User submits booking request with field_id, start_time, end_time
2. System validates field exists and time slot is available
3. System calculates total price
4. Booking record created with status
5. Payment invoice automatically generated
6. User receives booking confirmation with invoice

**Pagination:**
- All list endpoints support `?limit=10&page=1` query parameters
- Default: 10 items per page
- Efficient offset-based pagination

</details>

### Payment System

<details>
<summary><b>Click to expand payment features</b></summary>

**Invoice Generation:**
- Automatic invoice creation when booking is confirmed
- Links payment to specific booking and user
- Tracks total_price from booking
- Status: `pending`, `paid`, or `failed`

**Payment Tracking:**
- User-specific invoice list (`/api/v1/payments/me`)
- Admin payment oversight (`/api/v1/payments/list`)
- Individual invoice details (`/api/v1/payments/:id`)
- Payment-booking relationship maintained via foreign keys

**Dummy Payment Gateway:**
- Test endpoint (`/api/v1/payments/pay`) simulates payment processing
- Validates invoice exists and is in `pending` status
- Updates status to `paid` on success
- Ready for real gateway integration (Stripe, Midtrans, etc.)

**Payment Flow:**
1. Booking created → Invoice auto-generated with status `pending`
2. User views invoice in their payment list
3. User initiates payment (via dummy endpoint or real gateway)
4. System validates payment and updates status to `paid`
5. Booking status reflects payment completion

</details>

### User Management

<details>
<summary><b>Click to expand user management features</b></summary>

**User Registration:**
- Public endpoint for new user signup
- Validates email format and password strength
- Automatically assigns `user` role
- Encrypts email with AES-256-GCM
- Hashes email with SHA-256 for lookups
- Bcrypt hashes password

**Profile Management:**
- Users can update their own profile (`PATCH /api/v1/users/me`)
- Supports username, email, password updates
- Password change triggers re-hash

**Admin Capabilities:**
- **User Listing** (`GET /api/v1/users/list`): Paginated user directory
- **User Search** (`GET /api/v1/users/search?query=`): Find users by username/email
- **User Details** (`GET /api/v1/users/:id`): View individual user data
- **User Update** (`PATCH /api/v1/users/:id`): Modify any user's data
- **User Deletion** (`DELETE /api/v1/users/:id`): Soft delete (sets deleted_at)
- **Admin Creation** (`POST /api/v1/users/create`): Create admin accounts

**User Data Model:**
```go
- id (UUID)
- username (unique, indexed)
- email_enc (AES-256-GCM encrypted)
- email_hash (SHA-256 hash, unique, indexed)
- role (user/admin/superadmin)
- password_hash (bcrypt)
- last_login_at
- failed_login_attempts
- locked_until
- created_at, last_updated_at, deleted_at (soft delete)
```

</details>

### Web Interface

<details>
<summary><b>Click to expand web interface features</b></summary>

**Technology Stack:**
- **Templ**: Type-safe Go HTML templating (compile-time checked)
- **Alpine.js**: Lightweight reactive JavaScript framework (~15KB)
- **HTMX**: Server-driven interactions without heavy JavaScript
- **TailwindCSS**: Utility-first CSS via CDN (no build step required)

**Pages:**

**1. Dashboard (`/`)**
- Grid layout of all available fields
- Field cards with image, name, type, price/hour
- Real-time data fetching via `/api/v1/fields/list`
- Loading spinner during fetch
- Click to view field details
- Auto-redirect to `/login` if unauthenticated

**2. Field Detail (`/d/:id`)**
- Full field information display
- Image, name, type, description, location
- Real-time booking calendar
- Interactive booking form (start_time, end_time pickers)
- Live booking list for the field
- Instant booking creation with validation
- Success/error message display

**3. Login (`/login`)**
- Username/email + password input
- Form validation
- Error message display
- Redirect to dashboard on success
- Link to registration page

**4. Register (`/register`)**
- Email, username, password, confirm password fields
- Client-side password matching validation
- Success message with auto-redirect
- Link to login page

**UI/UX Features:**
- Dark mode by default (gray-900 background)
- Responsive design (mobile, tablet, desktop)
- Loading states on all async operations
- Indonesian Rupiah currency formatting
- Date/time localization (Asia/Jakarta timezone)
- Smooth transitions and hover effects
- Accessible form controls

</details>

### Multi-Role System

| Role | Booking | Field Management | User Management | Payments | Admin Creation |
|------|---------|-----------------|-----------------|----------|---------------|
| **User** | ✅ Create/View Own | ❌ | ❌ | ✅ View Own | ❌ |
| **Admin** | ✅ View All | ✅ Full CRUD | ✅ Full CRUD | ✅ View All | ❌ |
| **SuperAdmin** | ✅ View All | ✅ Full CRUD | ✅ Full CRUD | ✅ View All | ✅ Create Admins |

**Role Enforcement:**
- Implemented via `RoleMiddleware` in routing layer
- Checks `ctx.Locals("role")` set by JWT middleware
- Returns 403 Forbidden if role insufficient
- Multiple roles can be allowed per endpoint

### API Highlights

**Design Principles:**
- **Versioned**: All endpoints under `/api/v1` for future compatibility
- **RESTful**: Standard HTTP methods (GET, POST, PATCH, DELETE)
- **Consistent**: Uniform JSON response structure
- **Descriptive**: Clear error messages with HTTP status codes
- **Performant**: Request duration tracking in metadata

**Response Format:**
```json
{
  "data": { /* response payload */ },
  "meta": {
    "duration": "45.2ms",
    "page": 1,
    "limit": 10
  }
}
```

**Error Format:**
```json
{
  "error": "descriptive error message",
  "meta": {
    "duration": "12.1ms"
  }
}
```

**Features:**
- Pagination on all list endpoints (`?limit=&page=`)
- Structured logging (Zerolog) with request/response tracking
- Health monitoring (`/health`) with system metrics
- Automatic recovery middleware (prevents crashes)
- Request ID tracking for debugging

---

## 🏗️ Architecture & Design

BookArena implements **Clean Architecture** principles with clear separation of concerns, making the codebase maintainable, testable, and scalable.

### System Architecture

```mermaid
flowchart LR
 subgraph ClientLayer["CLIENT LAYER"]
        Clients["Web <br> Mobile <br> Postman"]
  end
 subgraph AppLayer["APP LAYER - Fiber"]
        CTR["Controllers"]
        R["Routes"]
        MW["Middleware"]
  end
 subgraph BusinessLogic["BUSINESS LOGIC"]
        Services["AuthSvc <br> FieldSvc <br> BookingSvc <br> UserSvc <br> PaymentSvc"]
  end
 subgraph DataAccess["DATA ACCESS - SQLC"]
        Repo["Repositories"]
  end
 subgraph Database["DATABASE - PostgreSQL"]
        Tables["users <br> fields <br> bookings <br> payments"]
  end
    MW --> R
    R --> CTR
    Clients --> MW
    CTR --> Services
    Services --> Repo
    Repo --> Tables
```

### Clean Architecture Layers

**1. Presentation Layer** (`internal/controllers/` + `views/`)
- **Responsibility**: Handle HTTP requests, validate input, return responses
- **Components**:
  - `auth.go` - Authentication endpoints (register, login, refresh)
  - `bookings.go` - Booking management endpoints
  - `fields.go` - Field CRUD endpoints
  - `payments.go` - Payment/invoice endpoints
  - `users.go` - User management endpoints
  - `web.go` - Server-side rendered pages
- **Dependencies**: Services (via dependency injection)
- **No direct database access**

**2. Business Logic Layer** (`internal/modules/`)
- **Responsibility**: Implement use cases, business rules, validations
- **Module Structure** (each module follows same pattern):
  ```
  modules/
  ├── auth/
  │   ├── service.go       # Business logic
  │   ├── types.go         # DTOs, requests, responses
  │   └── jwt.go           # JWT utilities
  ├── bookings/
  │   ├── service.go       # Booking logic, conflict detection
  │   ├── repository.go    # Data access interface
  │   └── types.go
  ├── fields/
  │   ├── service.go
  │   ├── repository.go
  │   └── types.go
  └── [users, payments...]
  ```
- **Key Responsibilities**:
  - Business validation (e.g., booking conflicts)
  - Data transformation
  - Orchestration between repositories
  - Encryption/decryption logic
  - Price calculations

**3. Data Access Layer** (`internal/database/`)
```
database/
├── migrations/              # SQL migrations
│   ├── 000001_ext.up.sql   # Enable UUID extension
│   ├── 000002_users.up.sql
│   ├── 000003_fields.up.sql
│   ├── 000004_bookings.up.sql
│   └── 000005_payments.up.sql
├── queries/                 # SQLC query definitions
│   ├── users.sql
│   ├── fields.sql
│   ├── bookings.sql
│   └── payments.sql
├── generated/               # SQLC generated Go code
│   ├── db_gen.go           # Database interface
│   ├── models_gen.go       # Struct definitions
│   ├── querier_gen.go      # Query interface
│   └── [entity]_gen.go
└── database.go             # Connection management
```
- **Type-Safe SQL**: All queries defined in `.sql` files, validated at compile-time
- **No ORM Overhead**: Direct SQL with Go struct mapping
- **Migration Management**: Sequential, reversible migrations

### Database Schema

<details>
<summary><b>Click to view complete database schema</b></summary>

**Users Table:**
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email_hash TEXT UNIQUE NOT NULL,        -- SHA-256 for lookups
    email_enc BYTEA NOT NULL,               -- AES-256-GCM encrypted
    username VARCHAR(50) UNIQUE NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin', 'superadmin')),
    password_hash VARCHAR(255),             -- Bcrypt
    password_changed_at TIMESTAMPTZ,
    last_login_at TIMESTAMPTZ,
    failed_login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ                  -- Soft delete
);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email_hash ON users(email_hash);
```

**Fields Table:**
```sql
CREATE TABLE fields (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    image_url TEXT,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(100) NOT NULL,
    description TEXT,
    location VARCHAR(255),
    location_lat DOUBLE PRECISION,
    location_lon DOUBLE PRECISION,
    price_per_hour BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_fields_name ON fields(name);
```

**Bookings Table:**
```sql
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    field_id UUID NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    total_price BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (field_id) REFERENCES fields(id)
);
CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_field_id ON bookings(field_id);
```

**Payments Table:**
```sql
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id UUID NOT NULL,
    user_id UUID NOT NULL,
    total_price BIGINT NOT NULL,
    status VARCHAR(255) NOT NULL DEFAULT 'pending' 
        CHECK (status IN ('pending', 'paid', 'failed')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (booking_id) REFERENCES bookings(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_booking_id ON payments(booking_id);
```

**Relationships:**
- User → Bookings (1:N)
- User → Payments (1:N)
- Field → Bookings (1:N)
- Booking → Payment (1:1)

</details>

### Design Patterns

**1. Repository Pattern**
- Abstracts database operations behind interfaces
- Each entity has its own repository
- Enables easy testing with mock repositories
- Example:
  ```go
  type UserRepository interface {
      GetUserByID(id uuid.UUID) (*User, error)
      CreateUser(params CreateUserParams) (*User, error)
      UpdateUser(id uuid.UUID, params UpdateUserParams) error
      DeleteUser(id uuid.UUID) error
  }
  ```

**2. Service Layer Pattern**
- Encapsulates business logic
- Orchestrates multiple repository calls
- Handles cross-cutting concerns (validation, encryption)
- Example: `BookingService` checks conflicts before creating booking

**3. Dependency Injection**
- Constructor-based injection
- No globals or singletons
- Easy to test and mock
- Example:
  ```go
  func NewAuthController(
      authService auth.AuthService,
      validator *validator.Validate,
  ) *AuthController
  ```

**4. Middleware Chain**
- Request processing pipeline
- Order: Recovery → Rate Limit → Logger → CORS → Cookie → Auth
- Each middleware can short-circuit the chain
- Example: `JWTMiddleware` validates token before reaching controller

**5. Factory Pattern**
- Used in service creation
- Centralizes dependency wiring
- See `internal/routes/index.go` for initialization

**6. DTO Pattern**
- Separate request/response types from domain models
- Validation tags on DTOs
- Example: `CreateBookingRequest`, `BookingResponse`

---

## 🛠️ Tech Stack

### Backend Technologies

| Technology | Version | Purpose | Why Chosen |
|------------|---------|---------|------------|
| **[Go](https://golang.org/)** | 1.24.0 | Primary language | High performance, built-in concurrency, strong typing, excellent tooling |
| **[Fiber](https://gofiber.io/)** | v2.52.9 | Web framework | Express-inspired API, blazing fast (built on Fasthttp), low memory footprint |
| **[PostgreSQL](https://www.postgresql.org/)** | 15+ | Database | ACID compliance, rich feature set, excellent JSON support, proven reliability |
| **[SQLC](https://sqlc.dev/)** | Latest | SQL code generation | Type-safe database access, compile-time query validation, no ORM overhead |
| **[pgx](https://github.com/jackc/pgx)** | v5.7.5 | PostgreSQL driver | High-performance, native Go driver with excellent connection pooling |
| **[Viper](https://github.com/spf13/viper)** | v1.20.1 | Configuration | Environment variable management, .env file support, type-safe config access |
| **[golang-jwt](https://github.com/golang-jwt/jwt)** | v5.2.3 | JWT tokens | Industry-standard JWT implementation, supports multiple algorithms |
| **[Zerolog](https://github.com/rs/zerolog)** | v1.34.0 | Structured logging | Zero-allocation JSON logger, high performance, contextual logging |
| **[go-playground/validator](https://github.com/go-playground/validator)** | v10.27.0 | Input validation | Extensive validation tags, custom validators, cross-field validation |

### Frontend Technologies

| Technology | Purpose | Benefits |
|------------|---------|----------|
| **[Templ](https://templ.guide/)** | Go template engine | Type-safe templates, compile-time checking, IDE support, hot reload |
| **[Alpine.js](https://alpinejs.dev/)** | Reactive JavaScript | Lightweight (15KB), Vue-like syntax, no build step, perfect for SSR |
| **[HTMX](https://htmx.org/)** | Dynamic HTML | Server-driven interactions, partial page updates, no heavy JS frameworks |
| **[TailwindCSS](https://tailwindcss.com/)** | CSS framework | Utility-first, rapid prototyping, consistent design system, via CDN |

### Security & Cryptography

| Component | Implementation | Purpose |
|-----------|----------------|---------|
| **Password Hashing** | Bcrypt | Adaptive hashing with automatic salt |
| **Email Encryption** | AES-256-GCM | Symmetric encryption for PII |
| **Email Hashing** | SHA-256 | Fast lookups without exposing emails |
| **Token Signing** | HMAC-SHA256 | JWT signature verification |
| **Cookie Encryption** | Fiber middleware | Prevents cookie tampering |

### Development Tools

| Tool | Purpose |
|------|---------|
| **[Air](https://github.com/cosmtrek/air)** | Hot reload for Go applications during development |
| **[golang-migrate](https://github.com/golang-migrate/migrate)** | Database migration management (up/down migrations) |
| **[Postman](https://www.postman.com/)** | API testing and documentation |
| **[ApiDog](https://apidog.com/)** | API documentation and collaboration platform |

### Infrastructure & Deployment

- **Database Connection Pooling**: pgx connection pool for optimal performance
- **Rate Limiting**: Sliding window algorithm (60 req/min per IP)
- **Timezone**: Asia/Jakarta (customizable via environment)
- **Logging**: Structured JSON logs with timestamps
- **Health Monitoring**: `/health` endpoint with Fiber Monitor

---

## 🚀 Quick Start

Get BookArena running in under 5 minutes:

```bash
# 1. Clone repository
git clone https://github.com/1001api/bookarena.git
cd bookarena

# 2. Install dependencies
go mod download

# 3. Setup database (PostgreSQL must be running)
createdb bookarena

# 4. Run migrations
./scripts/migrate.sh up

# 5. Generate code (if needed)
sqlc generate
templ generate

# 6. Start server
go run ./cmd/server

# Server running at http://localhost:8181
# Dashboard: http://localhost:8181
# API: http://localhost:8181/api/v1
# Health: http://localhost:8181/health
```

**Test the API:**
```bash
# Register a new user
curl -X POST http://localhost:8181/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"Pass123!@#"}'

# Login
curl -X POST http://localhost:8181/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"Pass123!@#"}' \
  -c cookies.txt

# Get your profile
curl http://localhost:8181/api/v1/users/me -b cookies.txt
```

---

## 📦 Installation & Setup

### Prerequisites

**Required:**
- **Go**: 1.24 or higher ([Download](https://golang.org/dl/))
- **PostgreSQL**: 15+ ([Download](https://www.postgresql.org/download/))

**Optional (but recommended):**
- **Air**: For hot reload during development
  ```bash
  go install github.com/cosmtrek/air@latest
  ```
- **golang-migrate**: For manual migration management
  ```bash
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  ```
- **sqlc**: If modifying database queries
  ```bash
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
  ```
- **templ**: If modifying templates
  ```bash
  go install github.com/a-h/templ/cmd/templ@latest
  ```

### Step-by-Step Installation

#### 1. Clone Repository

```bash
git clone https://github.com/1001api/bookarena.git
cd bookarena
```

#### 2. Install Go Dependencies

```bash
go mod download
go mod verify
```

#### 3. Configure Environment Variables

The project includes a `.env` file for easy development setup. For production, use proper environment variable management.

**Copy and edit the .env file:**
```bash
# For development, the existing .env file should work
# For production, create a new one:
cp .env .env.production
```

**Required Environment Variables:**

```env
# Server Configuration
PORT=8181                                    # HTTP port

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/bookarena?sslmode=disable

# Application Domain
APP_DOMAIN=localhost                         # Cookie domain
CLIENT_DOMAIN=http://localhost:5173          # CORS origin
ALLOWED_DOMAINS="['localhost','yourdomain.com']"

# Security Keys (MUST CHANGE IN PRODUCTION!)
COOKIE_ENC_KEY=your-32-character-key-here    # Must be exactly 32 bytes
ENC_KEY=your-32-character-key-here           # Must be exactly 32 bytes  
JWT_SECRET=your-long-random-secret-key       # Any length

# JWT Token Lifetimes (in minutes)
JWT_EXPIRY=60                                # Access token: 60 minutes
JWT_REFRESH_EXPIRY=4320                      # Refresh token: 3 days

# Root User (auto-created on first run)
ROOT_EMAIL=admin@yourdomain.com
ROOT_USERNAME=admin
ROOT_PASSWORD='SecurePassword123!@#'
```

**Generate Secure Keys:**
```bash
# Generate 32-byte encryption keys
openssl rand -base64 32

# Generate JWT secret
openssl rand -base64 64
```

#### 4. Database Setup

**Create Database:**
```bash
# Using psql
createdb bookarena

# Or manually
psql -U postgres -c "CREATE DATABASE bookarena;"
```

**Run Migrations:**
```bash
# Using the provided script (recommended)
./scripts/migrate.sh up

# Or manually with migrate CLI
migrate -path internal/database/migrations \
        -database "postgres://user:pass@localhost:5432/bookarena?sslmode=disable" \
        up
```

**Verify Migration:**
```bash
psql -U postgres -d bookarena -c "\dt"
# Should show: users, fields, bookings, payments tables
```

#### 5. Generate Code (Optional)

Only needed if you've modified SQL queries or templates:

```bash
# Generate SQLC code (if queries modified)
sqlc generate

# Generate Templ templates (if templates modified)
templ generate
# Or use the script
./scripts/templ.sh
```

#### 6. Run the Server

**Development Mode (with hot reload):**
```bash
air

# Server will restart automatically on code changes
# Watches: *.go, *.templ files
```

**Production Mode:**
```bash
# Build binary
go build -o bin/bookarena -ldflags="-s -w" ./cmd/server

# Run binary
./bin/bookarena
```

**Direct Run (no build):**
```bash
go run ./cmd/server
```

#### 7. Verify Installation

**Check Server Health:**
```bash
curl http://localhost:8181/health
```

**Access Web Interface:**
- Dashboard: http://localhost:8181
- Login: http://localhost:8181/login  
- Register: http://localhost:8181/register

**Test API:**
```bash
# Health check
curl http://localhost:8181/health

# Register test user
curl -X POST http://localhost:8181/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"demo","email":"demo@test.com","password":"Pass123!@#"}'
```

## 📡 Complete API Reference

**Base URL:** `https://bookarena.1001api.com/api/v1` (Production)  
**Base URL:** `http://localhost:8181/api/v1` (Development)

**External Documentation:**
- **ApiDog Interactive Docs:** [https://vifvpe310e.apidog.io](https://vifvpe310e.apidog.io)
- **Postman Collection:** [Download ZIP](https://github.com/1001api/bookarena/releases/download/resources/BookArenaPostman.zip)
- **OpenAPI/Swagger Spec:** [Download JSON](https://github.com/1001api/bookarena/releases/download/resources/BookArena.openapi.json)

### Authentication Endpoints

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json
```

**Request Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "Pass123!@#"
}
```

**Response (200 OK):**
```json
{
  "data": "550e8400-e29b-41d4-a716-446655440000",
  "meta": {
    "duration": "120.5ms"
  }
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json
```

**Request Body:**
```json
{
  "identifier": "john_doe",  // username or email
  "password": "Pass123!@#"
}
```

**Response (200 OK):**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  },
  "meta": {
    "duration": "150.2ms"
  }
}
```
**Sets Cookies:** `__hk_asid` (access token), `__hk_rsid` (refresh token)

#### Refresh Token
```http
POST /api/v1/auth/refresh
Cookie: __hk_rsid=<refresh_token>
```

**Response (200 OK):**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  },
  "meta": {
    "duration": "50.1ms"
  }
}
```

### User Endpoints

#### Get Current User Profile
```http
GET /api/v1/users/me
Cookie: __hk_asid=<access_token>
```

**Response (200 OK):**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "john_doe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2024-10-27T14:30:00Z"
  },
  "meta": {
    "duration": "25.3ms"
  }
}
```

#### Update Current User
```http
PATCH /api/v1/users/me
Cookie: __hk_asid=<access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "username": "new_username",
  "email": "newemail@example.com",
  "password": "NewPass123!@#"
}
```

#### List Users (Admin Only)
```http
GET /api/v1/users/list?limit=10&page=1
Cookie: __hk_asid=<admin_access_token>
```

#### Search Users (Admin Only)
```http
GET /api/v1/users/search?query=john&limit=10&page=1
Cookie: __hk_asid=<admin_access_token>
```

#### Get User by ID (Admin Only)
```http
GET /api/v1/users/:id
Cookie: __hk_asid=<admin_access_token>
```

#### Update User (Admin Only)
```http
PATCH /api/v1/users/:id
Cookie: __hk_asid=<admin_access_token>
Content-Type: application/json
```

#### Delete User (Admin Only)
```http
DELETE /api/v1/users/:id
Cookie: __hk_asid=<admin_access_token>
```

#### Create Admin User (SuperAdmin Only)
```http
POST /api/v1/users/create
Cookie: __hk_asid=<superadmin_access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "username": "new_admin",
  "email": "admin@example.com",
  "password": "AdminPass123!@#",
  "role": "admin"
}
```

### Field Endpoints

#### List Fields
```http
GET /api/v1/fields/list?limit=10&page=1
Cookie: __hk_asid=<access_token>
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "field-uuid-1",
      "name": "Main Soccer Field",
      "type": "Soccer",
      "description": "Professional grass field",
      "image_url": "https://example.com/field.jpg",
      "location": "Jakarta Stadium",
      "location_lat": -6.2088,
      "location_lon": 106.8456,
      "price_per_hour": 200000,
      "created_at": "2024-10-27T10:00:00Z"
    }
  ],
  "meta": {
    "duration": "35.2ms"
  }
}
```

#### Get Field by ID
```http
GET /api/v1/fields/:id
Cookie: __hk_asid=<access_token>
```

#### Create Field (Admin Only)
```http
POST /api/v1/fields/create
Cookie: __hk_asid=<admin_access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Basketball Court A",
  "type": "Basketball",
  "description": "Indoor professional court",
  "image_url": "https://example.com/court.jpg",
  "location": "Sports Center Jakarta",
  "location_lat": -6.2088,
  "location_lon": 106.8456,
  "price_per_hour": 150000
}
```

#### Update Field (Admin Only)
```http
PATCH /api/v1/fields/:id
Cookie: __hk_asid=<admin_access_token>
Content-Type: application/json
```

#### Delete Field (Admin Only)
```http
DELETE /api/v1/fields/:id
Cookie: __hk_asid=<admin_access_token>
```

### Booking Endpoints

#### Create Booking (User Only)
```http
POST /api/v1/bookings/create
Cookie: __hk_asid=<user_access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "field_id": "550e8400-e29b-41d4-a716-446655440000",
  "start_time": "2024-10-28T14:00:00Z",
  "end_time": "2024-10-28T16:00:00Z"
}
```

**Response (200 OK):**
```json
{
  "data": {
    "booking_id": "booking-uuid",
    "invoice_id": "invoice-uuid",
    "total_price": 400000
  },
  "meta": {
    "duration": "85.4ms"
  }
}
```

**Error Response (409 Conflict):**
```json
{
  "error": "Booking conflict with already existing booking, please choose another timeframe",
  "meta": {
    "duration": "45.1ms"
  }
}
```

#### Get User Bookings
```http
GET /api/v1/bookings/me?limit=10&page=1
Cookie: __hk_asid=<access_token>
```

#### Get Field Bookings
```http
GET /api/v1/bookings/field/:id?limit=10&page=1
Cookie: __hk_asid=<access_token>
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "booking-uuid",
      "field_name": "Main Soccer Field",
      "user_username": "john_doe",
      "start_time": "2024-10-28T14:00:00Z",
      "end_time": "2024-10-28T16:00:00Z",
      "total_price": 400000,
      "invoice_status": "pending",
      "created_at": "2024-10-27T10:30:00Z"
    }
  ],
  "meta": {
    "duration": "40.5ms"
  }
}
```

#### Get Booking by ID
```http
GET /api/v1/bookings/:id
Cookie: __hk_asid=<access_token>
```

#### List All Bookings (Admin Only)
```http
GET /api/v1/bookings/list?limit=10&page=1
Cookie: __hk_asid=<admin_access_token>
```

#### Delete Booking (Admin Only)
```http
DELETE /api/v1/bookings/:id
Cookie: __hk_asid=<admin_access_token>
```

### Payment Endpoints

#### Get User Invoices
```http
GET /api/v1/payments/me?limit=10&page=1
Cookie: __hk_asid=<access_token>
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "invoice-uuid",
      "booking_id": "booking-uuid",
      "total_price": 400000,
      "status": "pending",
      "created_at": "2024-10-27T10:30:00Z"
    }
  ],
  "meta": {
    "duration": "30.2ms"
  }
}
```

#### Get Invoice by ID
```http
GET /api/v1/payments/:id
Cookie: __hk_asid=<access_token>
```

#### List All Invoices (Admin Only)
```http
GET /api/v1/payments/list?limit=10&page=1
Cookie: __hk_asid=<admin_access_token>
```

#### Pay Invoice (Dummy)
```http
POST /api/v1/payments/pay
Cookie: __hk_asid=<access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "invoice_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response (200 OK):**
```json
{
  "data": {
    "id": "invoice-uuid",
    "message": "Payment paid successfully"
  },
  "meta": {
    "duration": "55.7ms"
  }
}
```

#### Delete Invoice (Admin Only)
```http
DELETE /api/v1/payments/:id
Cookie: __hk_asid=<admin_access_token>
```

### Web Routes

| Route | Description |
|-------|-------------|
| `GET /` | Dashboard page (field listing) |
| `GET /d/:id` | Field detail page with booking form |
| `GET /login` | Login page |
| `GET /register` | Registration page |
| `GET /health` | System health monitoring |

### Common Response Codes

| Code | Meaning | Description |
|------|---------|-------------|
| **200** | OK | Request succeeded |
| **201** | Created | Resource created successfully |
| **400** | Bad Request | Invalid request body or parameters |
| **401** | Unauthorized | Missing or invalid authentication |
| **403** | Forbidden | Insufficient permissions |
| **404** | Not Found | Resource doesn't exist |
| **409** | Conflict | Resource conflict (e.g., booking overlap) |
| **500** | Internal Server Error | Server error (check logs) |

---

## 🔧 Examples / Quickstart

```bash
# Register
curl -X POST https://bookarena.1001api.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"Pass123!@#"}'
```

```bash
# Login
curl -X POST https://bookarena.1001api.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"Pass123!@#"}' \
  -c cookies.txt
```

```bash
# List Fields
curl -X GET "https://bookarena.1001api.com/api/v1/fields/list?limit=10&page=1" -b cookies.txt
```

---

## 🔐 Authentication

BookArena uses dual JWT tokens (`access` + `refresh`) stored in encrypted HTTPOnly cookies.

| Token                       | Lifetime | Purpose                  |
| --------------------------- | -------- | ------------------------ |
| Access Token (`__hk_asid`)  | 60 min   | Authenticated requests   |
| Refresh Token (`__hk_rsid`) | 3 days   | Obtain new access tokens |

### Flow

1. User logs in
2. Server issues both tokens
3. Client stores them in cookies
4. Server validates on each request

---

## 💻 Code Examples

See comprehensive examples with cURL, JavaScript, Go, and Python in the [examples](./examples) directory.

---

## 🔐 Authentication Deep Dive

**JWT Token System:**
- Access tokens expire in 60 minutes
- Refresh tokens valid for 3 days
- Tokens stored in encrypted HTTPOnly cookies
- Automatic token rotation on refresh

**Security Features:**
- Bcrypt password hashing (cost: 10)
- AES-256-GCM email encryption
- SHA-256 email hashing for lookups
- Cookie encryption via Fiber middleware
- Rate limiting: 60 req/min per IP

---

## 👤 Demo Accounts

| Username   | Role       | Password     | Capabilities |
| ---------- | ---------- | ------------ | ------------ |
| **root**   | superadmin | `Pass123!@#` | Full system access |
| **admin1** | admin      | `Pass123!@#` | Manage fields & users |
| **user1**  | user       | `Pass123!@#` | Create bookings |

**Login URL:** https://bookarena.1001api.com/login

---

## 🧪 Development Guide

**Workflow:**
```bash
# 1. Create migration
./scripts/migrate.sh create add_new_feature

# 2. Write SQL in migrations/
# Edit: 000006_add_new_feature.up.sql

# 3. Apply migration
./scripts/migrate.sh up

# 4. Define SQLC queries in queries/
# 5. Generate code
sqlc generate && templ generate

# 6. Implement business logic
# 7. Add controllers & routes
# 8. Test with Air hot reload
air
```

---

## 🤝 Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open Pull Request

---

## 📄 License

This project is open source. See [LICENSE](./LICENSE) file for details.

---

## 💡 FAQ

**Q: How do I reset the database?**
```bash
./scripts/migrate.sh drop
./scripts/migrate.sh up
```

**Q: How do I add a real payment gateway?**
Replace the dummy endpoint in `internal/controllers/payments.go` with Stripe/Midtrans integration.

---

## Resources

- **Live Demo:** https://bookarena.1001api.com
- **API Docs:** https://vifvpe310e.apidog.io
- **Issues:** [GitHub Issues](https://github.com/1001api/bookarena/issues)


<br/>
<br/>


<div align="center">

**Built with ❤️ using Go by [@resqiar](https://github.com/resqiar) | dev of [TenOhOne](https://github.com/1001api)**

</div>
