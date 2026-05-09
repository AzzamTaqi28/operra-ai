# Operra Security Requirements

This document defines the security rules for Operra v0.1.

Operra handles sensitive operational data: purchase requests, attachments, approvals, comments, and audit logs.

Security and tenant isolation are core requirements.

## Tenant isolation

Keep every tenant boundary explicit:

- Every tenant-owned table has `organization_id`.
- Every tenant-owned query filters by `organization_id`.
- Every mutation checks organization ownership.
- Users belong to one organization in v0.1.
- Do not allow users to access resources from another organization.

High-risk bug to avoid:

```text
GET /purchase-requests/{id}
```

must not fetch only by `id`. It must fetch by:

```text
id + organization_id
```

## Authentication

Authentication rules:

- Passwords must be hashed securely.
- Login errors should not reveal whether email exists.
- JWT/session secret must come from environment variable.
- Token expiration should be reasonable.
- Do not store raw passwords.
- Do not log passwords or tokens.

## Authorization

Backend must enforce permissions.

Frontend may hide UI, but the backend is the authority.

Important permission checks:

- Only owner/admin can manage users.
- Only owner/admin can manage workflows.
- Requester can create requests.
- Requester can view own requests.
- Manager can approve only relevant department steps.
- Finance can approve finance steps.
- Procurement can process procurement steps.
- Director can approve director steps.
- Auditor has read-only access.

## Approval safety

Approval safety rules:

- User cannot approve request from another organization.
- User cannot approve a step requiring a role they do not have.
- User cannot approve a non-current step.
- User cannot approve own request in v0.1.
- Required steps cannot be skipped.
- Rejected and completed requests are terminal in v0.1.

## AI safety

AI must not be trusted as authority.

AI can:

- Generate workflow JSON.
- Generate diagrams.
- Explain workflow.

AI cannot:

- Approve requests.
- Reject requests.
- Execute workflow actions.
- Save workflows without admin confirmation.
- Bypass validation.
- Bypass RBAC.

Backend must validate AI-generated workflow JSON.

## Attachment security

Attachment rules:

- File metadata stored in database.
- File content stored in S3-compatible storage.
- Do not expose public object URLs.
- v0.1 must use authenticated backend streaming for downloads.
- Future versions may optionally add short-lived signed URLs.
- Before serving or generating access, check organization and request permission.
- Limit maximum file size.
- Restrict dangerous file types if needed.
- Sanitize filenames.

Suggested v0.1 limits:

```text
Max file size: 10 MB or configurable
Allowed default: PDF, images, documents, spreadsheets
```

## Audit log integrity

Audit log rules:

- Normal users cannot edit audit logs.
- Audit logs should be append-only at application level.
- Every approval action logs audit entry.
- Every CSV export logs audit entry.
- Every workflow activation logs audit entry.
- Every attachment upload/download logs audit entry if useful.

## CSV export security

CSV export rules:

- Export endpoints require permission.
- Export must be scoped by organization.
- Export filters must be validated.
- Export itself must create audit log.
- CSV output must escape formulas to avoid spreadsheet injection if possible.

Spreadsheet formula injection protection:

If a CSV cell begins with one of these characters:

```text
= + - @
```

prefix with a single quote or otherwise escape safely.

## Secrets management

Never commit:

- JWT secret.
- Database password.
- S3 secret key.
- AI API key.
- SMTP password.

Use `.env.example` with placeholder values only.

## Logging safety

Never log:

- Passwords.
- Tokens.
- API keys.
- Secret keys.
- Full file contents.

## Minimum security tests

1. User from organization A cannot access organization B requests.
2. User without role cannot approve.
3. Requester cannot approve own request.
4. Auditor cannot mutate request.
5. CSV export requires permission.
6. Attachment download requires permission.
7. AI-generated invalid workflow is rejected.
8. Workflow update does not alter submitted request version.
