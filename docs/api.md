# Operra API Contract

## 1. API principles

- REST API first.
- JSON request and response bodies.
- Base path: `/api/v1`.
- Authenticated endpoints require bearer token or secure session.
- Every tenant-owned response must be scoped by organization.
- Backend is the authority for permissions.

## 2. Common response formats

### Success object

```json
{
  "data": {}
}
```

### Success list

```json
{
  "data": [],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100
  }
}
```

### Error

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "estimated_amount is required",
    "details": []
  }
}
```

## 3. Error codes

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

## 4. Auth endpoints

### POST /api/v1/auth/register-organization

Creates organization and initial owner user.

Request:

```json
{
  "organization_name": "Demo Company",
  "organization_slug": "demo-company",
  "owner_name": "Owner Name",
  "owner_email": "owner@example.com",
  "password": "change-me"
}
```

Response:

```json
{
  "data": {
    "organization": {
      "id": "uuid",
      "name": "Demo Company",
      "slug": "demo-company"
    },
    "user": {
      "id": "uuid",
      "name": "Owner Name",
      "email": "owner@example.com"
    },
    "token": "jwt-token"
  }
}
```

### POST /api/v1/auth/login

Request:

```json
{
  "email": "owner@example.com",
  "password": "change-me"
}
```

Response:

```json
{
  "data": {
    "token": "jwt-token",
    "user": {
      "id": "uuid",
      "organization_id": "uuid",
      "name": "Owner Name",
      "email": "owner@example.com",
      "roles": ["owner"]
    }
  }
}
```

### GET /api/v1/auth/me

Returns current user.

## 5. Organization endpoints

### GET /api/v1/organization

Returns current organization.

### PATCH /api/v1/organization

Owner/admin updates organization settings.

## 6. User endpoints

### GET /api/v1/users

Query params:

```text
page
page_size
status
role
department_id
search
```

### POST /api/v1/users

Creates user.

Request:

```json
{
  "name": "Finance User",
  "email": "finance@example.com",
  "password": "temporary-password",
  "department_id": "uuid",
  "role_keys": ["finance"]
}
```

### GET /api/v1/users/{id}

### PATCH /api/v1/users/{id}

### POST /api/v1/users/{id}/roles

Assign roles.

Request:

```json
{
  "role_keys": ["manager", "requester"]
}
```

### DELETE /api/v1/users/{id}/roles/{role_key}

Remove role.

## 7. Department endpoints

### GET /api/v1/departments

### POST /api/v1/departments

Request:

```json
{
  "name": "Finance",
  "code": "FIN"
}
```

### PATCH /api/v1/departments/{id}

### DELETE /api/v1/departments/{id}

## 8. Workflow endpoints

### GET /api/v1/workflows

Query params:

```text
type
status
```

### POST /api/v1/workflows

Creates logical workflow and first version.

Request:

```json
{
  "name": "Purchase Request Approval",
  "type": "purchase_request",
  "config_json": {}
}
```

### GET /api/v1/workflows/{id}

### GET /api/v1/workflows/{id}/versions

### POST /api/v1/workflows/{id}/versions

Creates new workflow version.

Request:

```json
{
  "config_json": {}
}
```

Response:

```json
{
  "data": {
    "id": "uuid",
    "version_number": 2,
    "validation": {
      "valid": true,
      "errors": []
    },
    "mermaid_diagram": "flowchart TD ..."
  }
}
```

### POST /api/v1/workflows/validate

Validates workflow JSON without saving.

Request:

```json
{
  "type": "purchase_request",
  "config_json": {}
}
```

### POST /api/v1/workflows/{id}/versions/{version_id}/activate

Activates a workflow version.

### POST /api/v1/workflows/generate-with-ai

Generates workflow config from prompt.

Request:

```json
{
  "prompt": "Create purchase approval flow where above 5 million needs finance approval and above 25 million needs director approval."
}
```

Response:

```json
{
  "data": {
    "workflow_json": {},
    "explanation": "This workflow routes requests based on estimated amount.",
    "mermaid_diagram": "flowchart TD ...",
    "validation": {
      "valid": true,
      "errors": []
    },
    "warnings": []
  }
}
```

## 9. Purchase request endpoints

### GET /api/v1/purchase-requests

Query params:

```text
page
page_size
status
department_id
requester_id
from_date
to_date
search
```

Role behavior:

- Requester sees own requests.
- Manager sees department-related requests and pending approvals.
- Finance, procurement, director, owner, admin, auditor see according to permission rules.

### POST /api/v1/purchase-requests

Creates draft purchase request.

Request:

```json
{
  "department_id": "uuid",
  "title": "Buy laptops for new engineers",
  "item_name": "Laptop",
  "description": "Need laptops for new engineering hires",
  "quantity": 3,
  "estimated_amount": 45000000,
  "currency": "IDR",
  "urgency": "high",
  "expected_date": "2026-07-01",
  "vendor_name": "Vendor ABC",
  "notes": "Prefer business warranty"
}
```

### GET /api/v1/purchase-requests/{id}

Returns request detail including:

- Request fields.
- Attachments.
- Approval steps.
- Approval actions.
- Comments.
- Audit timeline if user has permission.

### PATCH /api/v1/purchase-requests/{id}

Updates draft or revision_requested request.

### POST /api/v1/purchase-requests/{id}/submit

Submits request and generates approval step instances.

### POST /api/v1/purchase-requests/{id}/cancel

Cancels request if allowed.

## 10. Approval endpoints

### GET /api/v1/approvals/pending

Returns current user's pending approvals.

### POST /api/v1/purchase-requests/{id}/approval-actions

Approves, rejects, or requests revision for current pending step.

Request:

```json
{
  "action": "approve",
  "comment": "Approved. Proceed."
}
```

Supported actions:

```text
approve
reject
request_revision
```

Reject and request_revision should require comment.

## 11. Comment endpoints

### POST /api/v1/purchase-requests/{id}/comments

Request:

```json
{
  "body": "Please attach updated quotation."
}
```

### GET /api/v1/purchase-requests/{id}/comments

## 12. Attachment endpoints

### POST /api/v1/purchase-requests/{id}/attachments

Multipart upload.

Response:

```json
{
  "data": {
    "id": "uuid",
    "file_name": "quotation.pdf",
    "file_size": 123456,
    "mime_type": "application/pdf"
  }
}
```

### GET /api/v1/purchase-requests/{id}/attachments/{attachment_id}/download

Returns signed URL or streams file after permission check.

Recommended response if signed URL:

```json
{
  "data": {
    "download_url": "https://signed-url",
    "expires_at": "2026-06-01T10:00:00Z"
  }
}
```

## 13. Audit endpoints

### GET /api/v1/audit-logs

Query params:

```text
page
page_size
action
entity_type
entity_id
from_date
to_date
actor_user_id
```

Only owner, admin, finance, and auditor should access broad audit logs. Requester may see request-level timeline for their own requests.

## 14. Export endpoints

### GET /api/v1/exports/purchase-requests.csv

Query params:

```text
status
department_id
from_date
to_date
```

Returns CSV.

### GET /api/v1/exports/approval-history.csv

Query params:

```text
from_date
to_date
request_id
```

Returns CSV.

### GET /api/v1/exports/audit-logs.csv

Query params:

```text
action
entity_type
from_date
to_date
```

Returns CSV.

Every export endpoint must create audit log `csv.exported`.

## 15. Health endpoints

### GET /health

Basic health check.

### GET /ready

Readiness check including database and storage connectivity if possible.

## 16. Pagination

Use:

```text
page=1&page_size=20
```

Response:

```json
{
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100
  }
}
```

## 17. API acceptance tests

Minimum API tests:

1. Organization registration creates owner and roles.
2. Login returns token.
3. User cannot access another organization's purchase request.
4. Workflow validation rejects invalid config.
5. Purchase request submission creates approval steps.
6. Manager can approve manager step.
7. Finance cannot approve manager step.
8. Requester cannot approve own request.
9. Rejection marks request rejected.
10. CSV export creates audit log.
