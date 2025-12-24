# TourOS Platform - Architecture Overview

## High-Level Architecture

TourOS follows **Clean Architecture** principles with strict separation of concerns across three primary layers:

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP Layer (Gin)                      │
│  Handlers → Middleware (Auth, Logging, Metrics, Rate)   │
└──────────────────────┬──────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│                  Business Logic Layer                    │
│  Services (Auth, Guide, Agency, Permit, Safety)         │
└──────────────────────┬──────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│                  Data Access Layer                       │
│  Repositories (GORM) → PostgreSQL                        │
└──────────────────────────────────────────────────────────┘
```

## Layer Responsibilities

### 1. Domain Layer (`internal/domain/`)

Pure domain entities with no external dependencies:

- **User**: Authentication and authorization entity
- **Agency**: Tourism agency with verification workflow
- **Guide**: Trek guide profiles linked to users
- **Permit**: Trek permits with QR code generation
- **SafetyCheckIn**: Daily safety check-ins with GPS
- **Incident**: Safety incidents including SOS

**Key Principles:**
- No business logic
- Pure data structures
- Database-agnostic (GORM tags for persistence)

### 2. Repository Layer (`internal/repository/`)

Data access abstraction:

- **Interfaces**: Define contracts for data operations
- **Implementations**: GORM-based implementations
- **Separation**: Business logic never touches GORM directly

**Patterns:**
- Interface-based design for testability
- Single responsibility per repository
- Preloading relationships when needed

### 3. Service Layer (`internal/service/`)

Business logic orchestration:

- **AuthService**: JWT token generation/validation, password hashing
- **GuideService**: Guide lifecycle, verification, license expiry checks
- **AgencyService**: Agency registration, verification workflow
- **PermitService**: Permit issuance, QR code generation, validation
- **SafetyService**: Check-ins, incident reporting, SOS tracking

**Key Principles:**
- Transaction management
- Business rule enforcement
- Error handling and validation

### 4. Handler Layer (`internal/handler/`)

HTTP request/response handling:

- **Request binding**: Gin's JSON binding
- **Response formatting**: Consistent JSON responses
- **Error handling**: Proper HTTP status codes
- **Metric tracking**: Business metric increments

**Patterns:**
- Thin handlers (delegate to services)
- Input validation at handler level
- Consistent error responses

### 5. Middleware Layer (`internal/middleware/`)

Cross-cutting concerns:

- **AuthMiddleware**: JWT validation, user context injection
- **RequireRole**: RBAC enforcement
- **LoggerMiddleware**: Structured logging with correlation IDs
- **MetricsMiddleware**: Prometheus metrics collection
- **RateLimitMiddleware**: Per-IP rate limiting

## Database Design

### Schema Overview

```
users
├── id (UUID, PK)
├── email (unique)
├── password_hash
├── role (admin|agency|guide)
├── agency_id (FK, nullable)
└── is_active

agencies
├── id (UUID, PK)
├── registration_number (unique)
├── license_number (unique)
├── status (pending|verified|suspended|rejected)
├── license_expiry
├── verified_at
└── verified_by (FK)

guides
├── id (UUID, PK)
├── user_id (FK, unique)
├── agency_id (FK, nullable)
├── license_number (unique)
├── status (pending|verified|suspended|rejected)
├── license_expiry
├── last_check_in
└── verified_by (FK)

permits
├── id (UUID, PK)
├── permit_number (unique)
├── guide_id (FK)
├── client_id
├── start_date
├── end_date
├── route
├── status (active|expired|revoked)
├── qr_code
└── issued_by (FK)

safety_check_ins
├── id (UUID, PK)
├── guide_id (FK)
├── permit_id (FK, nullable)
├── latitude
├── longitude
├── location
└── check_in_time

incidents
├── id (UUID, PK)
├── incident_type (check_in|sos|medical|weather|other)
├── guide_id (FK)
├── permit_id (FK, nullable)
├── status (open|in_progress|resolved|closed)
├── latitude
├── longitude
├── description
└── resolved_at
```

### Key Design Decisions

1. **Soft Deletes**: Using GORM's `DeletedAt` for audit trail
2. **UUID Primary Keys**: Better for distributed systems
3. **Status Enums**: Type-safe status management
4. **Timestamps**: Automatic `created_at`, `updated_at` tracking
5. **Relationships**: Proper foreign keys with indexes

## Authentication & Authorization

### JWT Token Flow

1. **Login**: User provides email/password → Access + Refresh tokens
2. **Access Token**: Short-lived (15min), contains user_id, email, role
3. **Refresh Token**: Long-lived (7 days), used to get new access tokens
4. **Authorization**: Bearer token in `Authorization` header

### RBAC Implementation

Three roles with different capabilities:

- **Admin**: Full access, can verify/suspend guides/agencies, revoke permits
- **Agency**: Manage own agency, view associated guides
- **Guide**: Manage own profile, create check-ins, report incidents

### Security Measures

- Password hashing: bcrypt with default cost
- Token secrets: Separate secrets for access/refresh tokens
- Rate limiting: 100 req/s per IP with burst of 200
- Input validation: Gin's binding validation
- SQL injection: Protected by GORM parameterization

## Observability

### Metrics (Prometheus)

**HTTP Metrics:**
- `http_requests_total` - Total requests by method, path, status
- `http_request_duration_seconds` - Request latency histogram

**Business Metrics:**
- `permits_issued_total` - Permit issuance counter
- `check_ins_total` - Safety check-in counter
- `sos_incidents_total` - SOS incident counter

### Logging (Structured JSON)

**Log Fields:**
- `request_id` - Correlation ID (from header or generated)
- `user_id` - Authenticated user ID
- `role` - User role
- Standard HTTP fields (method, path, status, latency, IP, user-agent)

**Log Levels:**
- INFO: Normal operations
- WARN: Client errors (4xx)
- ERROR: Server errors (5xx)

### Tracing (OpenTelemetry)

- **Provider**: Jaeger
- **Instrumentation**: Request-level spans
- **Correlation**: Request IDs linked to traces

## Configuration Management

12-factor app principles:

- **Environment Variables**: All configuration via env vars
- **No Secrets in Code**: JWT secrets must be provided via env
- **Defaults**: Sensible defaults for development
- **Validation**: Config validation on startup

### Required Environment Variables

```
JWT_ACCESS_SECRET (required)
JWT_REFRESH_SECRET (required)
DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
```

### Optional Environment Variables

```
SERVER_PORT (default: 8080)
APP_ENV (default: development)
LOG_LEVEL (default: info)
OTEL_ENABLED (default: false)
```

## API Design

### RESTful Principles

- **Resources**: Nouns (guides, agencies, permits)
- **HTTP Methods**: GET (read), POST (create), PUT (update)
- **Status Codes**: Proper HTTP status codes
- **JSON**: Consistent JSON request/response format

### Pagination

Standard offset/limit pagination:

```
GET /api/v1/guides?limit=20&offset=0
```

Response includes metadata:

```json
{
  "data": [...],
  "total": 100,
  "limit": 20,
  "offset": 0
}
```

### Filtering

Query parameters for filtering:

```
GET /api/v1/guides?status=verified&agency_id=...
GET /api/v1/permits?guide_id=...&status=active
```

## Error Handling

Consistent error response format:

```json
{
  "error": "Error message here"
}
```

HTTP status codes:
- `200 OK`: Success
- `201 Created`: Resource created
- `400 Bad Request`: Validation errors
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server errors

## Deployment Considerations

### Database Migrations

- Auto-migration on startup (development)
- For production: Use dedicated migration tool or manual SQL

### Health Checks

- `/health`: Basic health check
- `/ready`: Readiness check (includes DB ping)

### Graceful Shutdown

- Signal handling (SIGINT, SIGTERM)
- 30-second graceful shutdown window
- Connection cleanup

### Scaling Considerations

- Stateless API (can horizontal scale)
- Database connection pooling (25 max connections)
- Rate limiting per-instance (consider distributed rate limiting for multi-instance)

## Testing Strategy (Future)

Recommended testing layers:

1. **Unit Tests**: Service layer with mocked repositories
2. **Integration Tests**: Repository layer with test database
3. **API Tests**: Handler layer with test HTTP server
4. **End-to-End Tests**: Full stack with test environment

## Performance Considerations

1. **Database Indexes**: On foreign keys, status fields, dates
2. **Connection Pooling**: Configurable pool sizes
3. **Query Optimization**: Eager loading where needed, lazy loading otherwise
4. **Caching**: Consider Redis for frequently accessed data (future)

## Security Checklist

- [x] JWT token authentication
- [x] Password hashing (bcrypt)
- [x] RBAC implementation
- [x] Rate limiting
- [x] Input validation
- [x] SQL injection protection
- [x] Secrets via environment variables
- [ ] HTTPS/TLS (deployment requirement)
- [ ] CORS configuration (if frontend added)
- [ ] API key rotation strategy
- [ ] Audit logging

