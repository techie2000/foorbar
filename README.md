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

### Local Development with Docker Compose

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
- Local: http://localhost:8080/swagger/index.html
- Production: https://api.axiom.example.com/swagger/index.html

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
