# Axiom - Financial Services Static Data System

A modular monolith financial services static data management system built with Go, Next.js, and PostgreSQL.

## Overview

Axiom is a comprehensive financial services platform for managing static data including countries, currencies,
entities, financial instruments, accounts, and settlement instructions.

## Architecture

### Tech Stack

**Backend:**
- **Gin Framework** - Fast HTTP web framework for REST APIs
- **Fiber Framework** - Ultra-fast web framework for high-throughput services
- **GORM** - ORM for database operations
- **PostgreSQL** - Primary database
- **RabbitMQ** - Message queue for async operations
- **Beego** - Notification service framework

**Frontend:**
- **Next.js** - React framework for server-rendered applications
- **React** - UI component library
- **Tailwind CSS** - Utility-first CSS framework
- **shadcn/ui** - React UI components

**Features:**
- JWT Authentication
- CORS Configuration
- CQRS Pattern (Command Query Responsibility Segregation)
- Input Validation
- Error Handling Middleware
- Rate Limiting
- WebSocket Support
- Comprehensive Logging
- Prometheus Metrics
- Swagger API Documentation

## Project Structure

```text
axiom/
├── backend/
│   ├── cmd/                    # Application entry points
│   │   ├── api/               # Main API server (Gin)
│   │   ├── worker/            # Background worker (Fiber)
│   │   └── notification/      # Notification service (Beego)
│   ├── internal/              # Private application code
│   │   ├── domain/           # Domain models
│   │   ├── repository/       # Data access layer
│   │   ├── service/          # Business logic
│   │   ├── handler/          # HTTP handlers
│   │   ├── middleware/       # Custom middleware
│   │   ├── cqrs/            # CQRS implementation
│   │   └── config/          # Configuration
│   ├── pkg/                   # Public libraries
│   │   ├── auth/            # JWT authentication
│   │   ├── validator/       # Input validation
│   │   ├── logger/          # Logging utilities
│   │   └── queue/           # RabbitMQ client
│   ├── migrations/           # Database migrations
│   ├── docs/                 # Swagger documentation
│   └── tests/               # Integration tests
├── frontend/
│   ├── app/                  # Next.js app directory
│   ├── components/          # React components
│   ├── lib/                 # Utilities
│   └── public/             # Static assets
├── docker/                   # Docker configurations
├── scripts/                 # Build and deployment scripts
└── docs/                    # Documentation
```

## Development Phases

### Phase 1: Data Acquisition and Storage
- Implement domain data models (Countries, Currencies, Entities, Instruments, Accounts, SSI's)
- Set up PostgreSQL schemas and migrations
- Create API endpoints for CRUD operations
- Implement data acquisition from third-party sources (CSV, XML, JSON)
- Set up RabbitMQ for async processing

### Phase 2: Scheduled Updates
- Implement scheduled jobs for data refresh
- Add change detection and tracking
- Create webhook endpoints for real-time updates

### Phase 3: Auditing
- Implement audit logging for all data changes
- Add audit trail UI
- Create reporting endpoints

### Phase 4 and Beyond
- Advanced analytics
- Data visualization
- Export capabilities
- API rate limiting by client
- Multi-tenancy support

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- RabbitMQ 3.12+
- Docker & Docker Compose

### Multi-Environment Support

Axiom supports running multiple environments simultaneously on the same machine. Each environment uses a unique port
prefix to avoid conflicts:

- **Development (dev)**: Port prefix 1 (e.g., 18080, 13000, 15432)
- **UAT**: Port prefix 2 (e.g., 28080, 23000, 25432)
- **Production (prod)**: Port prefix 3 (e.g., 38080, 33000, 35432)

#### Starting a Specific Environment

```bash
# Start development environment
make docker-dev-up

# Start UAT environment
make docker-uat-up

# Start production environment
make docker-prod-up

# Start all environments simultaneously
make docker-all-up
```

#### Stopping Environments

```bash
# Stop development environment
make docker-dev-down

# Stop UAT environment
make docker-uat-down

# Stop production environment
make docker-prod-down

# Stop all environments
make docker-all-down
```

#### Checking Environment Status

```bash
# View status of all environments
make docker-all-status

# View logs for specific environment
make docker-dev-logs    # Development
make docker-uat-logs    # UAT
make docker-prod-logs   # Production
```

#### Environment-Specific URLs

Once started, each environment is accessible at:

**Development Environment:**
- Frontend: http://localhost:13000
- Backend API: http://localhost:18080
- Swagger UI: http://localhost:18080/swagger/index.html
- PostgreSQL: localhost:15432
- RabbitMQ Management: http://localhost:115672

**UAT Environment:**
- Frontend: http://localhost:23000
- Backend API: http://localhost:28080
- Swagger UI: http://localhost:28080/swagger/index.html
- PostgreSQL: localhost:25432
- RabbitMQ Management: http://localhost:25673

**Production Environment:**
- Frontend: http://localhost:33000
- Backend API: http://localhost:38080
- Swagger UI: http://localhost:38080/swagger/index.html
- PostgreSQL: localhost:35432
- RabbitMQ Management: http://localhost:315672

### Local Development with Docker Compose (Legacy)

```bash
# Start all services
docker-compose up -d

# Run migrations
make migrate-up

# Start backend
cd backend
go run cmd/api/main.go

# Start frontend (in another terminal)
cd frontend
npm install
npm run dev
```

### Running Database Migrations

```bash
# Run migrations on development environment
make migrate-dev-up

# Run migrations on UAT environment
make migrate-uat-up

# Run migrations on production environment
make migrate-prod-up

# Rollback migrations on specific environment
make migrate-dev-down
make migrate-uat-down
make migrate-prod-down
```

### Running Tests

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

## API Documentation

API documentation is available via Swagger UI at:
- Development: http://localhost:18080/swagger/index.html
- UAT: http://localhost:28080/swagger/index.html
- Production: http://localhost:38080/swagger/index.html
- Legacy/Local: http://localhost:8080/swagger/index.html

## Configuration

Configuration is managed through environment variables and config files:

```yaml
# config/config.yaml
database:
  host: localhost
  port: 5432
  name: axiom
  user: axiom
  password: ${DB_PASSWORD}

rabbitmq:
  url: amqp://guest:guest@localhost:5672/

jwt:
  secret: ${JWT_SECRET}
  expiry: 24h

server:
  port: 8080
  cors:
    allowed_origins:
      - http://localhost:3000
```

## Performance Optimization

- PostgreSQL caching for frequently accessed data
- Database query optimization with proper indexing
- Connection pooling for database connections
- Horizontal scaling with stateless services
- Request monitoring with Prometheus
- Structured logging with request tracing

## Security

- JWT-based authentication
- CORS configuration
- Input validation on all endpoints
- SQL injection prevention via ORM
- Rate limiting to prevent abuse
- Error messages don't expose sensitive information

## Contributing

Please read [CONTRIBUTING.md](docs/CONTRIBUTING.md) for details on our development workflow and code of conduct.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Documentation

Detailed documentation is available in the `docs/` directory:
- [Architecture Overview](docs/architecture.md)
- [API Reference](docs/api-reference.md)
- [Database Schema](docs/database-schema.md)
- [Deployment Guide](docs/deployment.md)
- [Development Workflow](docs/development-workflow.md)
