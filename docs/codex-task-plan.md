# Operra Codex Task Plan

This file defines a safe implementation sequence for Codex.

Do not ask Codex to build the entire app in one prompt. Use small tasks.

## Sprint 0 - Repository setup

### Task 0.1 - Initialize repo structure

Prompt:

```text
Read AGENTS.md and docs/architecture.md. Initialize the Operra monorepo structure with apps/api and apps/web. Do not implement business logic yet. Add README placeholders and .env.example based on docs/deployment.md.
```

Acceptance:

- Folder structure exists.
- `.env.example` exists.
- No unnecessary dependencies.

### Task 0.2 - Docker Compose skeleton

Prompt:

```text
Create Docker Compose setup for web, api, postgres, and minio based on docs/deployment.md. Include volumes and default environment variables. Do not add Redis or worker yet.
```

Acceptance:

- `docker compose up -d postgres minio` works.
- Ports are documented.

## Sprint 1 - Backend foundation

### Task 1.1 - API server foundation

Prompt:

```text
Create the Go API server foundation using Fiber. Add config loading, database connection placeholder, logging, health endpoint, and error response format based on docs/api.md and docs/architecture.md.
```

Acceptance:

- `GET /health` returns OK.
- Error format matches docs.

### Task 1.2 - Database migrations foundation

Prompt:

```text
Implement initial PostgreSQL migrations for organizations, users, departments, roles, and user_roles based on docs/data-model.md. Do not implement workflow tables yet.
```

Acceptance:

- Migrations run cleanly.
- Tables include organization_id where required.

### Task 1.3 - Auth and organization registration

Prompt:

```text
Implement organization registration, owner user creation, password hashing, login, and GET /auth/me. Seed built-in roles for each organization. Follow docs/api.md and docs/security.md.
```

Acceptance:

- Can register organization.
- Can login.
- Owner role exists.
- Passwords are hashed.

## Sprint 2 - Users, roles, departments

### Task 2.1 - Department CRUD

Prompt:

```text
Implement department CRUD endpoints. All queries must be scoped by organization_id. Follow docs/api.md and docs/data-model.md.
```

Acceptance:

- Admin can create/list/update/delete departments.
- Cross-organization access is impossible.

### Task 2.2 - User and role management

Prompt:

```text
Implement user list/create/update and role assignment endpoints. Enforce owner/admin permissions. Follow AGENTS.md and docs/security.md.
```

Acceptance:

- Admin can create users.
- Admin can assign roles.
- Non-admin cannot manage users.

## Sprint 3 - Workflow foundation

### Task 3.1 - Workflow data model

Prompt:

```text
Implement workflow and workflow_versions migrations and models based on docs/data-model.md. Implement create/list/get workflow endpoints without AI yet.
```

Acceptance:

- Admin can create workflow.
- Workflow versions are immutable.

### Task 3.2 - Workflow JSON validator

Prompt:

```text
Implement workflow JSON validation based on docs/workflow-engine.md. Add POST /workflows/validate. Include unit tests for invalid and valid configs.
```

Acceptance:

- Invalid configs return structured validation errors.
- Valid default PR workflow passes.

### Task 3.3 - Workflow activation and Mermaid generation

Prompt:

```text
Implement workflow version activation and simple Mermaid diagram generation from workflow config. Follow docs/workflow-engine.md.
```

Acceptance:

- Admin can activate version.
- Diagram string is returned.
- Activation creates audit log if audit module exists, or TODO hook if not yet built.

## Sprint 4 - Purchase requests

### Task 4.1 - Purchase request migrations and CRUD

Prompt:

```text
Implement purchase request migrations, models, create/list/detail/update endpoints based on docs/data-model.md and docs/api.md. Only draft creation and viewing in this task.
```

Acceptance:

- Requester can create draft.
- Requester can view own request.
- Tenant scoping enforced.

### Task 4.2 - Submit request and generate approval steps

Prompt:

```text
Implement purchase request submission. On submit, lock active workflow_version_id and generate approval_step_instances based on workflow config conditions. Follow docs/workflow-engine.md.
```

Acceptance:

- Amount thresholds generate correct steps.
- Request status becomes in_review.
- Current step is set.

## Sprint 5 - Approval engine and audit

### Task 5.1 - Audit log service

Prompt:

```text
Implement audit_logs migration, model, service, and basic list endpoint. Add helper function to create audit entries for mutations.
```

Acceptance:

- Audit logs are organization-scoped.
- Admin/auditor can list logs.

### Task 5.2 - Approval actions

Prompt:

```text
Implement approval action endpoint for approve, reject, and request_revision. Enforce role-based approver resolution, current step only, organization scoping, and requester cannot approve own request. Follow docs/workflow-engine.md and docs/security.md.
```

Acceptance:

- Approvals advance steps.
- Rejection works.
- Revision works.
- Audit logs are created.

### Task 5.3 - Comments

Prompt:

```text
Implement comments for purchase requests. Comments must be organization-scoped and logged in audit logs.
```

Acceptance:

- Authorized users can comment.
- Audit log created.

## Sprint 6 - Attachments

### Task 6.1 - S3/MinIO storage adapter

Prompt:

```text
Implement S3-compatible storage adapter based on docs/deployment.md and docs/security.md. Use environment variables for endpoint, bucket, access key, secret key, region, and path style. Include MinIO compatibility.
```

Acceptance:

- File can be uploaded to MinIO.
- File metadata saved.

### Task 6.2 - Attachment endpoints

Prompt:

```text
Implement purchase request attachment upload and download endpoints. Enforce permissions before upload/download. Do not expose public object URLs. Create audit logs.
```

Acceptance:

- Authorized upload works.
- Authorized download works.
- Unauthorized download blocked.

## Sprint 7 - CSV exports

### Task 7.1 - Purchase requests CSV export

Prompt:

```text
Implement GET /exports/purchase-requests.csv with filters and permission checks. Escape CSV safely and create audit log csv.exported.
```

Acceptance:

- CSV has correct headers.
- Filters work.
- Export audit log exists.

### Task 7.2 - Approval history and audit logs CSV export

Prompt:

```text
Implement approval-history and audit-logs CSV exports based on docs/api.md and docs/prd.md. Enforce permissions and audit export activity.
```

Acceptance:

- Both CSV exports work.
- Permission checks work.

## Sprint 8 - AI workflow builder

### Task 8.1 - AI provider abstraction

Prompt:

```text
Implement AI provider abstraction using OpenAI-compatible chat completion API settings from env. Do not hardcode provider. Do not expose API keys in logs.
```

Acceptance:

- AI client can be configured through env.
- Provider errors are handled.

### Task 8.2 - Generate workflow JSON from prompt

Prompt:

```text
Implement POST /workflows/generate-with-ai. It should accept a prompt, call AI provider, parse workflow JSON, validate it using the workflow validator, generate Mermaid diagram, and return validation results. It must not save workflow automatically.
```

Acceptance:

- Valid prompt returns valid workflow JSON.
- Invalid AI output returns validation errors.
- AI generation log or audit log created.

## Sprint 9 - Frontend

### Task 9.1 - Web app foundation

Prompt:

```text
Create Next.js frontend foundation with layout, auth pages, sidebar, API client, and protected routes based on docs/ui.md.
```

Acceptance:

- Login works.
- Dashboard shell works.

### Task 9.2 - Purchase request UI

Prompt:

```text
Implement purchase request list, create form, detail page, comments, attachments, and approval action panel based on docs/ui.md.
```

Acceptance:

- End-to-end purchase request flow works from UI.

### Task 9.3 - Workflow admin UI

Prompt:

```text
Implement workflow list, JSON editor, validate button, activate version, and Mermaid preview. Follow docs/ui.md and docs/workflow-engine.md.
```

Acceptance:

- Admin can validate and activate workflow from UI.

### Task 9.4 - AI workflow builder UI

Prompt:

```text
Implement AI Workflow Builder page with prompt input, generated JSON, explanation, warnings, Mermaid preview, and save workflow button requiring confirmation.
```

Acceptance:

- AI-generated config can be reviewed before saving.

### Task 9.5 - Users, departments, audit logs, exports UI

Prompt:

```text
Implement users, departments, audit logs, and exports pages based on docs/ui.md.
```

Acceptance:

- Admin can manage setup.
- Auditor/admin can view logs and export CSV.

## Sprint 10 - Private beta polish

### Task 10.1 - Seed data and demo workflow

Prompt:

```text
Add seed data for Demo Company, departments, users, roles, and default purchase request workflow. Make it easy to run locally.
```

### Task 10.2 - Documentation polish

Prompt:

```text
Update README with local setup, Docker Compose setup, env variables, default demo credentials, and troubleshooting.
```

### Task 10.3 - Tests and hardening

Prompt:

```text
Add tests for workflow transitions, tenant isolation, permission checks, attachments, CSV exports, and AI workflow validation based on docs/testing.md.
```
