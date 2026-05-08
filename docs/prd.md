# PRD - Operra v0.1: AI-Assisted Purchase Request Approval

## Document metadata

| Field | Value |
|---|---|
| Product | Operra |
| Version | v0.1 |
| Status | Draft for implementation |
| Owner | Taqi |
| Primary user segment | Growing companies with messy internal approvals |
| First workflow | Purchase Request approval |
| Deployment model | Self-hosted with Docker Compose, managed hosting later |
| AI scope | Workflow generation assistance only |

## 1. Product summary

Operra is a multi-tenant, self-hosted approval workflow platform for growing businesses.

In v0.1, Operra helps teams replace WhatsApp, email, and spreadsheet-based purchase approvals with structured, auditable approval workflows.

Users can submit purchase requests, route them through multi-level role-based approval, attach supporting documents, comment on requests, export records to CSV, and track every action through audit logs.

Operra v0.1 includes an AI workflow builder that lets admins describe a purchase approval process in natural language. The AI generates workflow JSON and a workflow diagram, but the backend validates and executes the workflow deterministically.

Core promise:

> Create a purchase approval workflow by chatting with AI, deploy it with Docker, store documents securely, export records to CSV, and keep every approval auditable.

## 2. Problem statement

Many operational teams manage purchase approvals through WhatsApp, email, spreadsheets, manual signatures, or verbal confirmation.

This creates several problems:

- Approval status is unclear.
- Requests get buried in chat.
- Finance lacks complete supporting documents.
- Business owners cannot easily see pending spending requests.
- Audit evidence is scattered across apps and people.
- SOPs are easy to bypass.
- Procurement data is unstructured before it reaches ERP or accounting systems.
- Approval disputes are hard to resolve because there is no reliable history.

Operra v0.1 focuses on one painful workflow: purchase request approval.

## 3. Target users

### Primary users

#### Requester

An employee, branch staff, or operational team member who needs to request a purchase.

Needs:

- Submit requests quickly.
- Attach supporting documents.
- Track approval status.
- See why a request was rejected or returned for revision.

#### Approver

A manager, finance staff, director, or procurement PIC who must review and act on requests.

Needs:

- See pending approvals.
- Review request details and attachments.
- Approve, reject, or request revision.
- Leave comments.

#### Admin

A user who configures users, roles, departments, and workflows.

Needs:

- Manage organization setup.
- Configure workflow JSON.
- Use AI to generate workflow config.
- Validate and activate workflow versions.

### Secondary users

#### Owner or Director

Needs visibility and control over spending requests, approval bottlenecks, and high-value purchases.

#### Finance

Needs complete evidence, approval history, and exportable reports.

#### Internal Auditor

Needs read-only access to requests, attachments, approval history, audit logs, and CSV exports.

## 4. v0.1 decisions

| Decision | Choice |
|---|---|
| Tenancy | Multi-tenant from day one |
| Tenancy model | Shared database, shared schema, scoped by `organization_id` |
| User organization membership | One organization per user in v0.1 |
| Approvers | Role-based |
| Workflow config | JSON-first |
| File storage | S3-compatible from day one |
| Local/self-hosted storage | MinIO included in Docker Compose |
| AI | Included in v0.1 for workflow generation only |
| Launch path | Private beta first, OSS after private beta |
| First workflow | Purchase Request approval |
| Export | CSV export required in v0.1 |

## 5. Goals

Operra v0.1 must allow a company to:

1. Create an organization.
2. Add users, departments, and roles.
3. Configure a Purchase Request approval workflow using JSON.
4. Generate workflow JSON and a diagram from AI chat.
5. Submit a purchase request.
6. Route the request through role-based approval steps.
7. Approve, reject, comment, or request revision.
8. Upload and download attachments using S3-compatible storage.
9. View audit logs for every important action.
10. Export purchase requests, approval history, and audit logs to CSV.
11. Deploy the app using Docker Compose.

## 6. Non-goals

Operra v0.1 will not include:

- Full procurement suite.
- Vendor scoring.
- Purchase order generation.
- Payment processing.
- Budgeting module.
- Inventory module.
- ERP integration.
- WhatsApp integration.
- Mobile app.
- Plugin marketplace.
- Drag-and-drop workflow builder.
- Complex SSO.
- Custom domain management.
- Billing/subscription system.
- Multi-organization user switching.
- Autonomous AI agents.
- AI approval or rejection decisions.

## 7. Roles and permissions

Initial roles:

| Role | Description | Main permissions |
|---|---|---|
| Owner | Organization owner | Manage organization, users, workflows, all requests, exports |
| Admin | Internal admin | Manage users, departments, roles, workflows |
| Requester | Request creator | Create and view own purchase requests |
| Manager | Department approver | Approve department-scoped requests |
| Finance | Finance approver | Approve finance steps, view financial requests, export reports |
| Procurement | Procurement processor | Process approved purchase requests |
| Director | High-value approver | Approve high-value requests |
| Auditor | Read-only auditor | View requests, attachments, audit logs, exports |

Backend must enforce permissions. Frontend-only permission checks are not enough.

## 8. Purchase Request workflow

### Default flow

```text
Requester submits PR
-> Manager approval
-> Finance approval if amount > 5,000,000
-> Director approval if amount > 25,000,000
-> Procurement processing
-> Completed
```

### Amount rules

```text
Amount <= IDR 5,000,000:
Requester -> Manager -> Procurement

Amount > IDR 5,000,000:
Requester -> Manager -> Finance -> Procurement

Amount > IDR 25,000,000:
Requester -> Manager -> Finance -> Director -> Procurement
```

### Request statuses

```text
draft
submitted
in_review
revision_requested
approved
rejected
processing
completed
cancelled
```

### Required request fields

```text
title
department_id
item_name
description
quantity
estimated_amount
currency
urgency
expected_date
```

### Optional request fields

```text
vendor_name
notes
attachments
```

## 9. Workflow configuration

Workflow config is JSON-first.

Admin can:

- View workflow config.
- Edit workflow JSON.
- Validate workflow JSON.
- Preview workflow steps.
- Generate Mermaid diagram.
- Save workflow as a new version.
- Activate/deactivate workflow versions.

Submitted requests must use the workflow version that was active at the time of submission. Updating a workflow must not modify active or completed requests.

Example config:

```json
{
  "name": "Purchase Request Approval",
  "type": "purchase_request",
  "version": 1,
  "steps": [
    {
      "key": "manager_approval",
      "name": "Manager Approval",
      "approver_role": "manager",
      "scope": "requester_department",
      "required": true
    },
    {
      "key": "finance_approval",
      "name": "Finance Approval",
      "approver_role": "finance",
      "scope": "organization",
      "required": true,
      "condition": {
        "field": "estimated_amount",
        "operator": ">",
        "value": 5000000
      }
    },
    {
      "key": "director_approval",
      "name": "Director Approval",
      "approver_role": "director",
      "scope": "organization",
      "required": true,
      "condition": {
        "field": "estimated_amount",
        "operator": ">",
        "value": 25000000
      }
    },
    {
      "key": "procurement_processing",
      "name": "Procurement Processing",
      "approver_role": "procurement",
      "scope": "organization",
      "required": true
    }
  ]
}
```

## 10. AI workflow builder

AI is included in v0.1 because it is a key differentiator.

### AI can

- Generate workflow JSON from natural language.
- Generate a Mermaid diagram.
- Explain the workflow.
- Suggest missing details.
- Validate whether the prompt is ambiguous.

### AI cannot

- Approve requests.
- Reject requests.
- Bypass RBAC.
- Execute workflow steps.
- Modify active workflows without admin confirmation.
- Make financial or procurement decisions.

Product principle:

> AI creates and explains configuration. Operra executes deterministic workflows.

### AI flow

```text
Admin opens AI Workflow Builder
-> Admin describes desired approval flow
-> AI generates workflow JSON
-> Backend validates workflow JSON
-> UI shows JSON, explanation, and diagram
-> Admin confirms and saves workflow version
```

### Example prompt

```text
Create a purchase request workflow.
Requests below 5 million only need manager approval.
Requests above 5 million need manager and finance approval.
Requests above 25 million need manager, finance, and director approval.
After approval, procurement should process the request.
```

## 11. Attachment storage

Operra v0.1 must support S3-compatible attachment storage.

Default self-hosted setup uses MinIO.

Supported storage modes:

| Mode | Description |
|---|---|
| S3-compatible | AWS S3, Cloudflare R2, DigitalOcean Spaces, MinIO, or compatible provider |
| MinIO | Default Docker Compose option for local/self-hosted usage |
| Local filesystem | Development fallback only, not recommended for production |

Attachment requirements:

- Store file metadata in PostgreSQL.
- Store file content in object storage.
- Enforce organization and request permissions before download.
- Use signed URLs or authenticated download endpoints.
- Do not expose public bucket URLs by default.
- Attachment upload must create audit log entry.

## 12. CSV export

CSV export is required in v0.1.

### Export types

#### Purchase request export

Fields:

```text
request_id
title
department
requester_name
estimated_amount
currency
status
current_step
created_at
submitted_at
completed_at
total_approval_duration
vendor_name
urgency
```

#### Approval history export

Fields:

```text
request_id
step_name
approver_role
approver_name
action
comment
acted_at
duration_from_previous_step
```

#### Audit log export

Fields:

```text
timestamp
actor
action
entity_type
entity_id
old_value
new_value
ip_address
```

### Export requirements

- Exports must respect organization scope.
- Exports must respect role permissions.
- Exports must support date range filters.
- Exports must support status filters where relevant.
- Exports must use UTF-8 encoding.
- Exports must escape commas, quotes, and newlines correctly.
- Every export must create an audit log entry.

## 13. Audit logs

Auditability is a core product feature.

Actions to log:

```text
organization.created
user.created
user.updated
role.assigned
role.removed
workflow.created
workflow.updated
workflow.activated
workflow.deactivated
request.created
request.submitted
request.updated
request.approved
request.rejected
request.revision_requested
request.cancelled
request.completed
comment.created
attachment.uploaded
attachment.downloaded
csv.exported
ai.workflow_generated
```

Audit log fields:

```text
id
organization_id
actor_user_id
entity_type
entity_id
action
old_value
new_value
ip_address
user_agent
created_at
```

## 14. Screens required in v0.1

Minimum frontend screens:

1. Login.
2. Organization setup.
3. Dashboard.
4. Purchase request list.
5. Create purchase request.
6. Purchase request detail.
7. Approval action panel.
8. Request timeline and audit history.
9. Attachment upload/download section.
10. Admin users page.
11. Admin departments page.
12. Admin workflow JSON editor.
13. AI workflow builder.
14. CSV export page.
15. Audit log page.
16. Settings page.

## 15. Success metrics

### Product success

- User can deploy app with Docker Compose.
- Admin can create first organization and user.
- Admin can configure Purchase Request workflow.
- Admin can generate workflow JSON using AI.
- Requester can submit a purchase request.
- Approver can approve/reject/comment.
- Request status updates correctly.
- Attachments upload and download correctly.
- CSV exports are generated correctly.
- Audit logs are generated for required actions.

### Business validation

- 10 user interviews completed.
- 3 design partners interested.
- 1 real company tests a pilot workflow.
- 1 paid pilot secured.
- OSS launch happens after private beta.

## 16. Acceptance checklist

Operra v0.1 is complete when:

- Multi-tenant organization scoping works.
- Users cannot access other organizations' data.
- Purchase request workflow works end-to-end.
- Role-based approval works.
- Workflow JSON validation works.
- AI-generated workflow config can be validated and saved.
- Attachments are stored through S3-compatible storage.
- MinIO works in Docker Compose.
- CSV export works for requests, approval history, and audit logs.
- Every approval and export creates audit log entry.
- Docker Compose setup works from a clean machine.
- README and setup docs are accurate.

## 17. v0.2 backlog

Possible v0.2 features:

- Expense reimbursement workflow.
- VendorVault module.
- Workflow template marketplace.
- Slack notifications.
- WhatsApp notification integration.
- Visual workflow builder.
- Advanced analytics dashboard.
- SSO.
- Audit evidence packs.
- Cloud managed hosting.
- Local model/Ollama full support if not completed in v0.1.
- Plugin/module system.
