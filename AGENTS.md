# Operra Agent Instructions

These instructions are for Codex and other coding agents working on Operra.

## Project identity

Operra is a multi-tenant, self-hosted approval workflow platform.

Version v0.1 focuses on AI-assisted Purchase Request approval.

The app must allow a company to create a purchase approval workflow, submit purchase requests, route them through role-based approval steps, store attachments, export CSV reports, and maintain a complete audit trail.

## Primary docs to read

Before implementing anything, read these files in order:

1. `docs/prd.md`
2. `docs/architecture.md`
3. `docs/data-model.md`
4. `docs/workflow-engine.md`
5. `docs/api.md`
6. `docs/security.md`
7. `docs/testing.md`

Use `docs/codex-task-plan.md` to decide the implementation order.

## Core engineering rules

- Build a modular monolith, not microservices.
- Use Go for the backend API.
- Use Next.js with TypeScript for the frontend.
- Use PostgreSQL as the source of truth.
- Use S3-compatible object storage for attachments.
- Include MinIO in Docker Compose for local/self-hosted storage.
- Use Docker Compose for local development and self-hosted deployment.
- Keep abstractions simple and explicit.
- Prefer readable code over clever generic frameworks.
- Do not build plugin architecture in v0.1.
- Do not build a drag-and-drop workflow builder in v0.1.

## Multi-tenancy rules

- v0.1 must support multiple organizations.
- Use shared database and shared schema.
- Every tenant-owned table must include `organization_id`.
- Every API request that reads or writes tenant-owned data must be scoped by `organization_id`.
- Never fetch tenant-owned records by ID without also checking `organization_id`.
- Users belong to exactly one organization in v0.1.
- Do not implement multi-organization user switching in v0.1.

## Workflow rules

- Workflow config is JSON-first.
- Workflow execution must be deterministic.
- A request must store the `workflow_version_id` that was active when it was submitted.
- Updating a workflow must not change already submitted requests.
- Approvers are role-based in v0.1.
- Do not allow required workflow steps to be skipped.
- Do not allow requesters to approve their own request unless a future explicit setting allows it.

## AI rules

- AI is included in v0.1 only for workflow creation assistance.
- AI may generate workflow JSON.
- AI may generate Mermaid diagrams.
- AI may explain the workflow.
- AI may suggest missing fields.
- AI must never approve a request.
- AI must never reject a request.
- AI must never bypass RBAC.
- AI must never modify an active workflow without admin confirmation.
- AI output must be validated by the backend before saving.
- AI-generated workflow JSON must pass the same validation as manually written JSON.

## Audit rules

- Every meaningful mutation must create an audit log.
- Every approval action must create an audit log.
- Every rejection must create an audit log.
- Every attachment upload must create an audit log.
- Every CSV export must create an audit log.
- Every AI workflow generation must create an audit log or AI generation log.
- Audit logs must be scoped by `organization_id`.
- Do not allow normal users to edit audit logs.

## Attachment rules

- Store file metadata in PostgreSQL.
- Store file content in S3-compatible object storage.
- Use MinIO for local/self-hosted development.
- Do not expose public object URLs.
- Use signed URLs or authenticated download endpoints.
- Always enforce organization and request permissions before file download.

## CSV export rules

- CSV export must respect organization scope.
- CSV export must respect role permissions.
- CSV exports must support date range filters where applicable.
- CSV exports must use UTF-8.
- CSV exports must safely escape commas, quotes, and newlines.
- CSV export itself must be recorded in audit logs.

## Security rules

- Passwords must be hashed with a secure password hashing function.
- JWT/session secrets must come from environment variables.
- Never commit secrets.
- Do not log raw API keys, passwords, or tokens.
- Enforce role permissions on the backend, not only on frontend.
- Protect against tenant data leakage.
- Validate user input server-side.

## Implementation style

- Keep feature modules separated by domain.
- Suggested backend modules: auth, organizations, users, roles, departments, workflows, requests, approvals, attachments, audit, exports, ai, notifications.
- Suggested frontend areas: dashboard, requests, approvals, workflows, users, audit logs, exports, settings.
- Write tests for workflow transitions, tenancy checks, and permission rules.
- Keep v0.1 focused on Purchase Request approval.

## Avoid in v0.1

- No microservices.
- No Kubernetes-specific deployment.
- No ERP replacement features.
- No vendor scoring.
- No payment processing.
- No purchase order generation.
- No full procurement suite.
- No visual drag-and-drop workflow builder.
- No plugin marketplace.
- No autonomous AI agents.
- No complex billing system.
- No multi-org user switching.
