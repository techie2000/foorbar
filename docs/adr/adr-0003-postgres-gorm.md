---
post_title: "ADR-0003: PostgreSQL with GORM"
author1: "techie2000"
post_slug: "adr-0003-postgres-gorm"
microsoft_alias: "techie2000"
featured_image: "https://placehold.co/1200x630.png"
categories: ["backend"]
tags: ["adr", "backend", "database", "gorm", "postgresql"]
ai_note: "AI-assisted draft based on repository state and user request."
summary: "Records the decision to use PostgreSQL as the primary database with GORM as the ORM."
post_date: "2026-02-10"
title: "ADR-0003: PostgreSQL with GORM"
status: "Accepted"
date: "2026-02-10"
authors: "techie2000"
supersedes: ""
superseded_by: ""
---

## Status

Accepted

## Context

Axiom requires a relational database to manage reference and static data with integrity constraints and
transactional consistency. The backend connects to PostgreSQL using GORM and Docker Compose provisions the
database for local development.

## Decision Drivers

- **DRV-001**: Relational integrity and strong transactional guarantees.
- **DRV-002**: Mature tooling for migrations and observability.
- **DRV-003**: ORM support to reduce boilerplate data access code.
- **DRV-004**: Operational familiarity and Docker-friendly deployment.

## Decision

Use PostgreSQL as the primary database and GORM as the ORM for data access in the Go backend.

## Decision Outcome

**Chosen Option:** PostgreSQL + GORM.

## Consequences

### Positive

- **POS-001**: Strong SQL capabilities and data integrity features.
- **POS-002**: GORM reduces manual query boilerplate for common operations.
- **POS-003**: Works well with containerized development workflows.

### Negative

- **NEG-001**: ORM abstractions can obscure performance issues.
- **NEG-002**: GORM migrations require explicit management outside the ORM.
- **NEG-003**: Advanced SQL features may need raw queries.

### Mitigation

- **MIT-001**: Review generated SQL for performance-critical paths.
- **MIT-002**: Maintain explicit migration files in [backend/migrations](../../backend/migrations).
- **MIT-003**: Use raw SQL when necessary for complex queries.

## Alternatives Considered

### MySQL

- **ALT-001**: **Description**: Widely used relational database with broad support.
- **ALT-002**: **Rejection Reason**: PostgreSQL offers richer features for data integrity and queries.

### SQLite for early phases

- **ALT-003**: **Description**: Embedded database for fast local setup.
- **ALT-004**: **Rejection Reason**: Lacks features needed for multi-user and production environments.

### SQLC with database/sql

- **ALT-005**: **Description**: Generate type-safe SQL with minimal ORM behavior.
- **ALT-006**: **Rejection Reason**: Higher upfront query authoring cost for CRUD-heavy domains.

## Implementation Notes

- **IMP-001**: Keep the connection setup in [backend/cmd/api/main.go](../../backend/cmd/api/main.go).
- **IMP-002**: Store migrations in [backend/migrations](../../backend/migrations).
- **IMP-003**: Use connection pooling settings aligned with service load.

## References

- **REF-001**: [backend/go.mod](../../backend/go.mod)
- **REF-002**: [backend/cmd/api/main.go](../../backend/cmd/api/main.go)
- **REF-003**: [docker-compose.yml](../../docker-compose.yml)
