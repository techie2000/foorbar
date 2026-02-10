---
post_title: "ADR-0007: Docker Compose for Local Development"
author1: "techie2000"
post_slug: "adr-0007-docker-compose-local-dev"
microsoft_alias: "techie2000"
featured_image: "https://placehold.co/1200x630.png"
categories: ["devops"]
tags: ["adr", "devops", "docker", "compose"]
ai_note: "AI-assisted draft based on repository state and user request."
summary: "Records the decision to use Docker Compose for local development and service orchestration."
post_date: "2026-02-10"
title: "ADR-0007: Docker Compose for Local Development"
status: "Accepted"
date: "2026-02-10"
authors: "techie2000"
supersedes: ""
superseded_by: ""
---

## Status

Accepted

## Context

Axiom relies on PostgreSQL, RabbitMQ, the Go backend, and the Next.js frontend. Local development benefits
from repeatable, one-command setup for all services. The repository includes Dockerfiles and a Compose file.

## Decision Drivers

- **DRV-001**: Consistent local environments across contributors.
- **DRV-002**: Simple service orchestration for multiple dependencies.
- **DRV-003**: Reproducible configuration for onboarding.
- **DRV-004**: Alignment with production containerization strategy.

## Decision

Use Docker Compose to orchestrate local development services and provide a single entry point for running
the stack.

## Decision Outcome

**Chosen Option:** Docker Compose for local development.

## Consequences

### Positive

- **POS-001**: One-command startup for the full stack.
- **POS-002**: Consistent service versions across the team.
- **POS-003**: Supports local health checks and dependencies.

### Negative

- **NEG-001**: Requires Docker availability on developer machines.
- **NEG-002**: Compose logs can be noisy for troubleshooting.
- **NEG-003**: Some debug workflows are slower in containers.

### Mitigation

- **MIT-001**: Document local overrides and non-container workflows.
- **MIT-002**: Provide guidance for filtering Compose logs.
- **MIT-003**: Use live reload for frontend development.

## Alternatives Considered

### Manual local installation

- **ALT-001**: **Description**: Developers install Postgres and RabbitMQ locally.
- **ALT-002**: **Rejection Reason**: Leads to inconsistent versions and setup drift.

### Kubernetes for local dev

- **ALT-003**: **Description**: Use a local K8s cluster for development.
- **ALT-004**: **Rejection Reason**: Higher operational cost for daily dev tasks.

### Remote shared dev environment

- **ALT-005**: **Description**: Centralized dev stack shared by all contributors.
- **ALT-006**: **Rejection Reason**: Slower feedback loops and less isolated testing.

## Implementation Notes

- **IMP-001**: Keep service definitions in [docker-compose.yml](../../docker-compose.yml).
- **IMP-002**: Use Dockerfiles in [docker](../../docker) for backend and frontend images.
- **IMP-003**: Maintain health checks for dependency readiness.

## References

- **REF-001**: [docker-compose.yml](../../docker-compose.yml)
- **REF-002**: [docker/Dockerfile.backend](../../docker/Dockerfile.backend)
- **REF-003**: [adr-0005-rabbitmq-async-processing.md](adr-0005-rabbitmq-async-processing.md)
