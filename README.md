# TourOS Platform API

Production-ready B2B + GovTech SaaS backend for Tourism Operating System for Nepal.

## Architecture

TourOS follows Clean Architecture principles with strict separation of concerns:

```
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── domain/          # Domain entities
│   ├── repository/      # Data access layer (GORM)
│   ├── service/         # Business logic layer
│   ├── handler/         # HTTP handlers (Gin)
│   ├── middleware/      # Auth, logging, metrics, rate limiting
│   ├── router/          # Route setup
│   ├── database/        # Database connection & migrations
│   └── observability/   # Logging, tracing, metrics
├── config/              # Prometheus, Grafana configs
└── docker-compose.yml   # Local development stack
```

## Tech Stack

- **Backend**: Go 1.21 with Gin framework
- **Database**: PostgreSQL 15
- **ORM**: GORM
- **Auth**: JWT (access + refresh tokens) with RBAC
- **Observability**: 
  - Metrics: Prometheus
  - Logging: Structured JSON logs (Loki-compatible)
  - Tracing: OpenTelemetry (Jaeger)
- **Containerization**: Docker & Docker Compose

## Features

### Core Modules

1. **Authentication Service**
   - JWT-based authentication
   - Access + refresh token pattern
   - Role-based access control (Admin, Agency, Guide)

2. **Guide & Agency Management**
   - Guide profile management
   - Agency registration and verification
   - License expiry tracking
   - Status management (pending, verified, suspended)

3. **Trek Permit Service**
   - Permit issuance with QR code generation
   - Permit validation
   - Permit revocation
   - Active permit tracking

4. **Safety Service**
   - Daily safety check-ins with GPS coordinates
   - SOS incident reporting
   - Incident management workflow
   - Last-seen tracking for guides

### Observability

- **Metrics**: HTTP latency, error rates, request counts, business metrics
- **Logging**: Structured JSON logs with request_id, user_id, role
- **Tracing**: Distributed tracing with correlation IDs
- **Health Checks**: `/health` and `/ready` endpoints

### Security

- JWT access + refresh tokens
- Role-based access control (RBAC)
- Rate limiting middleware
- Input validation
- SQL injection protection via GORM
- Secrets via environment variables

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15 (or use Docker Compose)

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone <repository-url>
cd touros-platform
```

2. Create `.env` file:
```bash
cp .env.example .env
```

3. Update `.env` with your JWT secrets (required):
```env
JWT_ACCESS_SECRET=your-min-32-character-access-secret
JWT_REFRESH_SECRET=your-min-32-character-refresh-secret
```

4. Start services:
```bash
docker-compose up -d
```

This starts:
- API server on `http://localhost:8080`
- PostgreSQL on `localhost:5432`
- Prometheus on `http://localhost:9090`
- Grafana on `http://localhost:3000` (admin/admin)
- Jaeger on `http://localhost:16686`

### Manual Setup

1. Install dependencies:
```bash
go mod download
```

2. Start PostgreSQL and ensure it's accessible

3. Set environment variables (see `.env.example`)

4. Run migrations (auto-migrates on startup):
```bash
go run cmd/api/main.go
```

## API Endpoints

### Authentication

- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh access token

### Guides

- `POST /api/v1/guides` - Create guide profile
- `GET /api/v1/guides` - List guides (with filters)
- `GET /api/v1/guides/:id` - Get guide by ID
- `PUT /api/v1/guides/:id` - Update guide
- `POST /api/v1/guides/:id/verify` - Verify guide (admin only)
- `POST /api/v1/guides/:id/suspend` - Suspend guide (admin only)

### Agencies

- `POST /api/v1/agencies` - Create agency
- `GET /api/v1/agencies` - List agencies (with filters)
- `GET /api/v1/agencies/:id` - Get agency by ID
- `PUT /api/v1/agencies/:id` - Update agency
- `POST /api/v1/agencies/:id/verify` - Verify agency (admin only)
- `POST /api/v1/agencies/:id/suspend` - Suspend agency (admin only)

### Permits

- `POST /api/v1/permits` - Issue permit
- `GET /api/v1/permits` - List permits (with filters)
- `GET /api/v1/permits/:id` - Get permit by ID
- `POST /api/v1/permits/:id/revoke` - Revoke permit (admin only)
- `GET /api/v1/permits/validate/:number` - Validate permit (public)

### Safety

- `POST /api/v1/safety/check-ins` - Create check-in
- `GET /api/v1/safety/check-ins/:id` - Get check-in by ID
- `GET /api/v1/safety/guides/:guide_id/check-ins` - List check-ins for guide
- `POST /api/v1/safety/incidents` - Report incident
- `GET /api/v1/safety/incidents` - List incidents (with filters)
- `GET /api/v1/safety/incidents/:id` - Get incident by ID
- `PUT /api/v1/safety/incidents/:id` - Update incident
- `GET /api/v1/safety/guides/:guide_id/sos` - Get active SOS for guide

### Health & Monitoring

- `GET /health` - Health check
- `GET /ready` - Readiness check (includes DB ping)
- `GET /metrics` - Prometheus metrics

## Authentication

All protected endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer <access_token>
```

### Example Login Request

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900,
  "token_type": "Bearer"
}
```

## Database Schema

The application uses GORM auto-migration. Key tables:

- `users` - User accounts with roles
- `agencies` - Tourism agencies
- `guides` - Trek guides linked to users
- `permits` - Trek permits with QR codes
- `safety_check_ins` - Daily check-ins
- `incidents` - Safety incidents including SOS

## Development

### Run Tests

```bash
make test
```

### Build

```bash
make build
```

### Run Locally

```bash
make run
```

### Database Migrations

Migrations run automatically on startup. For manual control, see `cmd/migrate/` (to be implemented).

## Production Deployment

1. Set strong JWT secrets (32+ characters)
2. Use proper SSL/TLS for database connections
3. Configure proper logging levels
4. Enable OpenTelemetry tracing
5. Set up proper rate limiting based on traffic
6. Configure Grafana dashboards
7. Set up alerts in Prometheus
8. Use secrets management (e.g., AWS Secrets Manager, HashiCorp Vault)

## Observability

### Metrics

Prometheus metrics are available at `/metrics`. Key metrics:

- `http_requests_total` - Total HTTP requests
- `http_request_duration_seconds` - Request latency
- `permits_issued_total` - Business metric: permits issued
- `check_ins_total` - Business metric: safety check-ins
- `sos_incidents_total` - Business metric: SOS incidents

### Logging

Structured JSON logs include:
- `request_id` - Correlation ID
- `user_id` - Authenticated user ID
- `role` - User role
- `service_name` - Service identifier
- Standard HTTP fields (method, path, status, latency)

### Tracing

OpenTelemetry traces are sent to Jaeger when enabled. Configure `OTEL_ENABLED=true` and set `OTEL_ENDPOINT`.

## License

[Your License Here]

## Contributing

[Your Contributing Guidelines Here]

