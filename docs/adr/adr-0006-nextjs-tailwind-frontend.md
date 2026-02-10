---
post_title: "ADR-0006: Next.js and Tailwind CSS for the Frontend"
author1: "techie2000"
post_slug: "adr-0006-nextjs-tailwind-frontend"
microsoft_alias: "techie2000"
featured_image: "https://placehold.co/1200x630.png"
categories: ["frontend"]
tags: ["adr", "frontend", "nextjs", "tailwind"]
ai_note: "AI-assisted draft based on repository state and user request."
summary: "Records the decision to use Next.js with Tailwind CSS for the Axiom frontend."
post_date: "2026-02-10"
title: "ADR-0006: Next.js and Tailwind CSS for the Frontend"
status: "Accepted"
date: "2026-02-10"
authors: "techie2000"
supersedes: ""
superseded_by: ""
---

## Status

Accepted

## Context

Axiom requires a web UI for managing reference data. The repository includes a Next.js app using the App
Router and Tailwind CSS, indicating a decision to standardize on this frontend stack.

## Decision Drivers

- **DRV-001**: Support server-rendered pages and modern React patterns.
- **DRV-002**: Fast iteration with a single framework for routing and builds.
- **DRV-003**: Utility-first styling for consistent UI primitives.
- **DRV-004**: TypeScript support for safer UI code.

## Decision

Use Next.js as the frontend framework with Tailwind CSS for styling and React for UI composition.

## Decision Outcome

**Chosen Option:** Next.js + Tailwind CSS.

## Consequences

### Positive

- **POS-001**: Unified routing, data fetching, and build pipeline.
- **POS-002**: Tailwind provides consistent design tokens and rapid styling.
- **POS-003**: React ecosystem supports component reuse and tooling.

### Negative

- **NEG-001**: Tailwind utility classes can be verbose in complex components.
- **NEG-002**: Next.js upgrades can require coordinated migrations.
- **NEG-003**: SSR adds build-time complexity compared to static SPA.

### Mitigation

- **MIT-001**: Use shared components and class composition utilities.
- **MIT-002**: Keep dependencies current and review release notes.
- **MIT-003**: Prefer static rendering where SSR is not needed.

## Alternatives Considered

### Vite with React SPA

- **ALT-001**: **Description**: Client-rendered React app with Vite tooling.
- **ALT-002**: **Rejection Reason**: Lacks integrated SSR and routing features.

### Vue with Nuxt

- **ALT-003**: **Description**: Vue-based SSR framework.
- **ALT-004**: **Rejection Reason**: Team stack aligns more with React tooling.

### Server-rendered Go templates

- **ALT-005**: **Description**: HTML templates rendered in the Go backend.
- **ALT-006**: **Rejection Reason**: Less flexible for rich client interactions.

## Implementation Notes

- **IMP-001**: Keep the app in [frontend/app](../../frontend/app).
- **IMP-002**: Configure Tailwind in [frontend/tailwind.config.ts](../../frontend/tailwind.config.ts).
- **IMP-003**: Define shared UI utilities in [frontend/package.json](../../frontend/package.json).

## References

- **REF-001**: [frontend/package.json](../../frontend/package.json)
- **REF-002**: [frontend/app/layout.tsx](../../frontend/app/layout.tsx)
- **REF-003**: [adr-0001-modular-monolith-architecture.md](adr-0001-modular-monolith-architecture.md)
