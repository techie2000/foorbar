---
post_title: "ADR-0005: RabbitMQ for Asynchronous Processing"
author1: "techie2000"
post_slug: "adr-0005-rabbitmq-async-processing"
microsoft_alias: "techie2000"
featured_image: "https://placehold.co/1200x630.png"
categories: ["backend"]
tags: ["adr", "backend", "messaging", "rabbitmq"]
ai_note: "AI-assisted draft based on repository state and user request."
summary: "Records the decision to use RabbitMQ for asynchronous workflows and background jobs."
post_date: "2026-02-10"
title: "ADR-0005: RabbitMQ for Asynchronous Processing"
status: "Accepted"
date: "2026-02-10"
authors: "techie2000"
supersedes: ""
superseded_by: ""
---

## Status

Accepted

## Context

Axiom needs background processing for data import, export, and future integrations. The configuration and
Docker Compose stack include RabbitMQ to support asynchronous job handling.

## Decision Drivers

- **DRV-001**: Decouple long-running work from API request latency.
- **DRV-002**: Reliable delivery semantics for background jobs.
- **DRV-003**: Operational familiarity and local dev support in Docker.
- **DRV-004**: Support for multiple consumers and routing patterns.

## Decision

Use RabbitMQ as the message broker for asynchronous workflows and job processing.

## Decision Outcome

**Chosen Option:** RabbitMQ for async processing and job queues.

## Consequences

### Positive

- **POS-001**: Background work no longer blocks API response times.
- **POS-002**: Flexible routing for different job types.
- **POS-003**: Durable queues support retry strategies.

### Negative

- **NEG-001**: Additional infrastructure component to monitor and maintain.
- **NEG-002**: Requires message schema versioning and evolution strategy.
- **NEG-003**: Operational overhead for scaling consumers.

### Mitigation

- **MIT-001**: Define message contracts and versioning rules early.
- **MIT-002**: Add health checks and monitoring for broker availability.
- **MIT-003**: Start with small consumer pools and scale as needed.

## Alternatives Considered

### Kafka

- **ALT-001**: **Description**: Distributed log with high throughput.
- **ALT-002**: **Rejection Reason**: Operational complexity is higher for current needs.

### AWS SQS

- **ALT-003**: **Description**: Managed queue service.
- **ALT-004**: **Rejection Reason**: Adds cloud lock-in and requires cloud credentials for local dev.

### In-process goroutines

- **ALT-005**: **Description**: Background work in the API process.
- **ALT-006**: **Rejection Reason**: Lacks durability and decoupling.

## Implementation Notes

- **IMP-001**: Configure broker connection via [backend/internal/config/config.go](../../backend/internal/config/config.go).
- **IMP-002**: Keep broker container in [docker-compose.yml](../../docker-compose.yml).
- **IMP-003**: Define queues and routing keys per job type.

## References

- **REF-001**: [docker-compose.yml](../../docker-compose.yml)
- **REF-002**: [backend/internal/config/config.go](../../backend/internal/config/config.go)
- **REF-003**: [adr-0001-modular-monolith-architecture.md](adr-0001-modular-monolith-architecture.md)
