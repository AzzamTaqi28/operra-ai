# Operra Workflow Chart

This document shows how the repository is organized and how the app is used end to end.

## Repo Workflow

Operra follows the public roadmap in `docs/roadmap.md`:

```mermaid
flowchart TD
  A[Docs: PRD, architecture, data model, workflow engine, API, security, testing] --> B[docs/roadmap.md]
  B --> C[apps/api]
  B --> D[apps/web]
  C --> E[PostgreSQL migrations]
  C --> F[REST API modules]
  C --> G[Workflow validation and approval engine]
  C --> H[Audit, attachments, exports, AI]
  D --> I[Next.js app shell]
  D --> J[shadcn-style UI]
  D --> K[Auth, dashboard, requests, workflows, users, departments]
  C --> L[Docker Compose]
  L --> M[PostgreSQL]
  L --> N[MinIO]
  L --> O[Web]
  L --> P[API]
  C --> Q[Tests]
  D --> Q
  Q --> R[Release ready v0.1]
```

What this means:

- Docs define the product, architecture, data model, and constraints.
- The roadmap turns those docs into a public build overview.
- The backend owns tenancy, permissions, workflow execution, and audit rules.
- Attachment downloads are served through authenticated backend checks in v0.1.
- The frontend is a thin operational layer over the API.
- Docker Compose ties the stack together for local and self-hosted use.

## App Usage

This is the main user journey in v0.1:

```mermaid
flowchart TD
  A[First-time setup] --> B[Create organization]
  B --> C[Create departments, users, and roles]
  C --> D[Create or generate workflow]
  D --> E[Validate workflow JSON]
  E --> F[Activate workflow version]
  F --> G[Requester creates purchase request]
  G --> H[Requester attaches files and saves draft]
  H --> I[Requester submits request]
  I --> J[Backend locks active workflow version]
  J --> K[Workflow engine generates approval steps]
  K --> L[Approver reviews request]
  L --> M{Approve, reject, or request revision?}
  M -->|Approve| N[Advance to next step]
  M -->|Reject| O[Mark request rejected]
  M -->|Request revision| P[Return to requester]
  N --> Q{More steps left?}
  Q -->|Yes| L
  Q -->|No| R[Mark request approved or completed]
  R --> S[Export CSV or review audit logs]
```

## Runtime Flow

This is how a request moves through the system:

```mermaid
sequenceDiagram
  participant U as User
  participant W as Next.js Web App
  participant A as Go API
  participant D as PostgreSQL
  participant S as MinIO/S3 Storage

  U->>W: Login / create request / approve
  W->>A: Authenticated request
  A->>D: Read or write tenant-scoped data
  A->>S: Store or fetch attachment object
  A->>D: Write audit log
  A-->>W: JSON response
  W-->>U: Updated UI state
```

## Request State Flow

The backend uses the request lifecycle defined in `docs/workflow-engine.md`:

```text
draft -> submitted -> in_review -> approved -> processing -> completed
                   -> revision_requested -> submitted
                   -> rejected
draft -> cancelled
submitted -> cancelled
revision_requested -> cancelled
```

Notes:

- The request stores the workflow version that was active when it was submitted.
- Workflow edits do not rewrite already submitted requests.
- Approvers are resolved by role, organization, and step scope.
- AI can propose workflow JSON, but the backend validates and executes it.

## How To Read The Repo

- Start with `docs/prd.md` for product intent.
- Use `docs/architecture.md` for system shape.
- Use `docs/workflow-engine.md` for approval rules and state transitions.
- Use `docs/api.md` for endpoint contracts.
- Use `docs/ui.md` for the screen list and UX expectations.
- Use `docs/roadmap.md` for the public implementation overview.
