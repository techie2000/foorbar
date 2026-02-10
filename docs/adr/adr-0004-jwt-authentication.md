---
post_title: "ADR-0004: JWT-Based API Authentication"
author1: "techie2000"
post_slug: "adr-0004-jwt-authentication"
microsoft_alias: "techie2000"
featured_image: "https://placehold.co/1200x630.png"
categories: ["backend"]
tags: ["adr", "backend", "security", "jwt"]
ai_note: "AI-assisted draft based on repository state and user request."
summary: "Records the decision to use JWTs for API authentication and authorization."
post_date: "2026-02-10"
title: "ADR-0004: JWT-Based API Authentication"
status: "Accepted"
date: "2026-02-10"
authors: "techie2000"
supersedes: ""
superseded_by: ""
---

## Status

Accepted

## Context

The API needs stateless authentication for protected endpoints with a straightforward developer experience.
The backend already includes JWT middleware and configuration for token secrets and expiry.

## Decision Drivers

- **DRV-001**: Stateless authentication compatible with horizontal scaling.
- **DRV-002**: Simple integration with HTTP headers and API clients.
- **DRV-003**: Standard, widely adopted token format.
- **DRV-004**: Clear separation of public and protected routes.

## Decision

Use JWTs for API authentication, validated in middleware for all protected routes.

## Decision Outcome

**Chosen Option:** JWT authentication with Bearer tokens.

## Consequences

### Positive

- **POS-001**: Stateless tokens simplify scaling across API instances.
- **POS-002**: Standard Bearer format works with common tooling.
- **POS-003**: Middleware enforcement centralizes auth logic.

### Negative

- **NEG-001**: Token revocation requires additional infrastructure.
- **NEG-002**: Secret rotation needs careful rollout.
- **NEG-003**: Large JWT payloads can increase request size.

### Mitigation

- **MIT-001**: Keep token payloads minimal and short-lived.
- **MIT-002**: Plan for revocation via deny lists or token versioning.
- **MIT-003**: Document rotation procedures for the signing secret.

## Alternatives Considered

### Session cookies

- **ALT-001**: **Description**: Server-stored sessions with cookie identifiers.
- **ALT-002**: **Rejection Reason**: Requires shared session storage across instances.

### OAuth2 token introspection

- **ALT-003**: **Description**: Centralized auth server validates tokens per request.
- **ALT-004**: **Rejection Reason**: Adds runtime dependencies and latency.

### API keys only

- **ALT-005**: **Description**: Static API keys for service access.
- **ALT-006**: **Rejection Reason**: Limited user context and weaker rotation model.

## Implementation Notes

- **IMP-001**: Use middleware in [backend/internal/middleware/middleware.go](../../backend/internal/middleware/middleware.go).
- **IMP-002**: Configure secrets and expiry via [backend/internal/config/config.go](../../backend/internal/config/config.go).
- **IMP-003**: Apply auth only to protected route groups.

## References

- **REF-001**: [backend/internal/middleware/middleware.go](../../backend/internal/middleware/middleware.go)
- **REF-002**: [backend/internal/config/config.go](../../backend/internal/config/config.go)
- **REF-003**: [adr-0002-go-gin-backend.md](adr-0002-go-gin-backend.md)
