# Operra Architecture

## 1. Architecture summary

Operra v0.1 is a modular monolith with a Go backend, Next.js frontend, PostgreSQL, S3-compatible object storage, and Docker Compose deployment.

The system is multi-tenant with a shared database and shared schema. Tenant-owned records are scoped with `organization_id`.

High-level architecture:

```text
Browser
  -> Next.js Web App
  -> Go REST API
  -> PostgreSQL
  -> S3-compatible Object Storage or MinIO
  -> AI Provider via OpenAI-compatible API
```

## 2. Core principles

- Modular monolith first.
- Multi-tenant from day one.
- Every tenant-owned table has `organization_id`.
- Workflow config is JSON-first.
- Workflow execution is deterministic.
- AI generates config only.
- Auditability is a first-class feature.
- S3-compatible storage from day one.
- Docker Compose must work for local and self-hosted deployment.
- Keep v0.1 focused on Purchase Request approval.

## 3. Recommended tech stack

| Layer | Choice | Notes |
|---|---|---|
| Backend | Go + Fiber | Fast and familiar for backend API |
| ORM | GORM | Good speed for v0.1 implementation |
| Database | PostgreSQL | Source of truth |
| Frontend | Next.js + TypeScript | Dashboard and admin UI |
| Styling | Tailwind CSS | Simple and fast UI development |
| Object storage | S3-compatible | AWS S3, Cloudflare R2, DO Spaces, MinIO |
| Local object storage | MinIO | Included in Docker Compose |
| AI | OpenAI-compatible chat completion API | Provider-abstracted client |
| Diagrams | Mermaid | Generated from workflow config |
| Deployment | Docker Compose | Required for v0.1 |

## 4. Monorepo structure

Suggested structure:

```text
operra/
  CONTRIBUTING.md
  README.md
  docker-compose.yml
  .env.example
  docs/
    prd.md
    architecture.md
    workflow-engine.md
    workflow-chart.md
    roadmap.md
    data-model.md
    api.md
    ui.md
    deployment.md
    security.md
    testing.md
  apps/
    api/
      cmd/
        server/
          main.go
      internal/
        auth/
        organizations/
        users/
        roles/
        departments/
        workflows/
        requests/
        approvals/
        attachments/
        audit/
        exports/
        ai/
        notifications/
        platform/
          config/
          database/
          middleware/
          storage/
          logger/
      migrations/
      go.mod
    web/
      app/
      components/
      lib/
      hooks/
      types/
      package.json
```

## 5. Backend module responsibilities

### auth

Responsible for:

- Login.
- Password hashing.
- Token/session generation.
- Current user endpoint.
- Auth middleware.

### organizations

Responsible for:

- Organization creation.
- Organization settings.
- Organization scoping.

### users

Responsible for:

- User CRUD.
- User activation/deactivation.
- User organization membership.

### roles

Responsible for:

- Built-in roles.
- Role assignment.
- Permission checks.

### departments

Responsible for:

- Department CRUD.
- Department assignment for users and requests.

### workflows

Responsible for:

- Workflow definitions.
- Workflow versions.
- JSON validation.
- Mermaid diagram generation.
- Workflow activation.

### requests

Responsible for:

- Purchase request CRUD.
- Request submission.
- Request status changes.
- Request detail and timeline.

### approvals

Responsible for:

- Approval step instances.
- Approver resolution.
- Approve/reject/revision actions.
- State transitions.

### attachments

Responsible for:

- File upload.
- File metadata.
- Signed downloads.
- S3/MinIO storage.

### audit

Responsible for:

- Audit log creation.
- Audit log queries.
- Immutable event history.

### exports

Responsible for:

- CSV exports.
- Export permission checks.
- Export audit logs.

### ai

Responsible for:

- AI provider abstraction.
- Workflow prompt templates.
- AI response parsing.
- AI output validation.

### notifications

Responsible for:

- Email notification later.
- Pending approval notifications.
- v0.1 can start with simple log or email stubs if SMTP is not configured.

## 6. Request lifecycle

```text
Draft created
-> Submitted
-> Workflow version is locked
-> Approval step instances are generated
-> Current step becomes pending
-> Approver acts
-> Engine advances to next required step
-> Request becomes approved when approval steps are complete
-> Procurement processes request
-> Request becomes completed
```

## 7. Multi-tenancy model

Use shared database and shared schema.

Rules:

- Every tenant-owned table includes `organization_id`.
- Every authenticated request has an organization context.
- Every query must include `organization_id` for tenant-owned records.
- A user belongs to exactly one organization in v0.1.
- Do not implement tenant switching in v0.1.

Example backend rule:

```text
Bad:
SELECT * FROM purchase_requests WHERE id = ?

Good:
SELECT * FROM purchase_requests WHERE id = ? AND organization_id = ?
```

## 8. Permission model

Use role-based access control.

Backend must check:

- Is the user authenticated?
- Does the user belong to the same organization as the resource?
- Does the user have the required role or permission?
- Is the user allowed to act on this workflow step?

Frontend can hide UI, but backend is the authority.

## 9. Workflow versioning

Important rule:

> Submitted requests must use the workflow version active at the time of submission.

Implementation:

- `workflows` stores the logical workflow.
- `workflow_versions` stores immutable versions of workflow JSON.
- `purchase_requests.workflow_version_id` locks the version.
- Updating a workflow creates a new version.
- Activating a workflow version affects only future submissions.

## 10. AI architecture

AI is a separate service module.

AI module responsibilities:

- Accept a user prompt.
- Construct a strict system prompt.
- Call an OpenAI-compatible API.
- Parse JSON output.
- Validate output using workflow validator.
- Return JSON, explanation, warnings, and Mermaid diagram.

AI module must not:

- Save workflow without explicit admin confirmation.
- Approve requests.
- Reject requests.
- Bypass permissions.

## 11. Storage architecture

File content goes to S3-compatible object storage.

File metadata goes to PostgreSQL.

Storage options:

- S3-compatible provider in production.
- MinIO in Docker Compose.
- Local filesystem only as a development fallback if needed.

Object key format:

```text
organizations/{organization_id}/requests/{request_id}/attachments/{attachment_id}/{safe_filename}
```

## 12. CSV export architecture

For v0.1, CSV exports can be synchronous if data size is small.

Rules:

- Apply filters.
- Check permissions.
- Generate CSV stream.
- Log `csv.exported` audit action.

Future versions can support async export jobs.

## 13. Notification architecture

v0.1 may implement basic email notification or stubbed notification events.

Minimum useful notifications:

- Approval needed.
- Request approved.
- Request rejected.
- Revision requested.
- Request completed.

Notifications should not block workflow transitions. If notification fails, the request action should still succeed and the failure should be logged.

## 14. Error handling

Use consistent API errors:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "estimated_amount is required",
    "details": []
  }
}
```

Common error codes:

```text
UNAUTHORIZED
FORBIDDEN
NOT_FOUND
VALIDATION_ERROR
CONFLICT
WORKFLOW_INVALID
WORKFLOW_TRANSITION_INVALID
STORAGE_ERROR
AI_PROVIDER_ERROR
INTERNAL_ERROR
```

## 15. Logging

Log server events without exposing secrets.

Never log:

- Passwords.
- Raw access tokens.
- API keys.
- S3 secrets.
- Full file content.

## 16. Environment variables

Minimum environment variables:

```env
APP_ENV=development
APP_URL=http://localhost:3000
API_URL=http://localhost:8080
DATABASE_URL=postgres://your-db-user:your-db-password@postgres:5432/operra?sslmode=disable
JWT_SECRET=change-me

STORAGE_DRIVER=s3
S3_ENDPOINT=http://minio:9000
S3_BUCKET=operra
S3_ACCESS_KEY=your-s3-access-key
S3_SECRET_KEY=change-me
S3_REGION=us-east-1
S3_FORCE_PATH_STYLE=true

AI_PROVIDER=openai
AI_BASE_URL=https://api.openai.com/v1
AI_API_KEY=
AI_MODEL=

SMTP_HOST=
SMTP_PORT=
SMTP_USER=
SMTP_PASSWORD=
SMTP_FROM=
```

## 17. v0.1 deployment topology

```text
Docker Compose services:

- web
- api
- postgres
- minio
```

Optional later:

- redis
- worker
- ollama
