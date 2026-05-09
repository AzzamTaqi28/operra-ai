# Operra UI Requirements

This document describes the public UI expectations for Operra v0.1.

## UI principles

- Clean operational dashboard, not a consumer social app
- Make status obvious
- Make approval actions fast
- Make audit history easy to inspect
- Avoid complex UI in v0.1
- JSON-first workflow config is acceptable for v0.1
- AI workflow builder should make workflow creation feel modern

## Main navigation

Suggested sidebar:

```text
Dashboard
Purchase Requests
Pending Approvals
Workflows
AI Workflow Builder
Users
Departments
Audit Logs
Exports
Settings
```

Role-based visibility should apply.

## Required screens

### Login

Fields:

- Email
- Password

Actions:

- Login

### Organization setup

For first-time setup.

Fields:

- Organization name
- Organization slug
- Owner name
- Owner email
- Password

### Dashboard

Cards:

- Total purchase requests.
- Pending approvals.
- Approved this month.
- Rejected this month.
- Average approval time.

Tables:

- Recent requests.
- My pending approvals.

### Purchase request list

Columns:

```text
Request ID
Title
Requester
Department
Amount
Status
Current Step
Created At
```

Filters:

```text
Status
Department
Date range
Search
```

Actions:

- Create request
- View detail
- Export CSV if permitted

### Create purchase request

Fields:

```text
Title
Department
Item name
Description
Quantity
Estimated amount
Currency
Urgency
Expected date
Vendor name optional
Notes optional
Attachment optional
```

Actions:

- Save draft
- Submit request

### Purchase request detail

Sections:

- Request summary.
- Request fields.
- Attachments.
- Current approval step.
- Approval actions.
- Comments.
- Approval timeline.
- Audit history if permitted.

Actions by role:

- Requester: edit draft, resubmit revision, cancel if allowed.
- Approver: approve, reject, request revision, comment.
- Auditor: read-only.

### Pending approvals

List of approval items assigned to the current user's roles.

Columns:

```text
Request ID
Title
Requester
Department
Amount
Current Step
Waiting Since
```

### Workflow list

Columns:

```text
Name
Type
Status
Active Version
Updated At
```

Actions:

- View
- Create
- Activate/deactivate

### Workflow JSON editor

Features:

- JSON textarea/editor.
- Validate button.
- Save as new version.
- Activate version.
- Diagram preview.
- Validation errors.

### AI Workflow Builder

Layout:

- Chat/prompt input on left.
- Generated JSON on right.
- Mermaid diagram preview.
- Explanation and warnings.
- Save workflow button.

Important UI rule:

> Saving AI-generated workflow must require explicit admin confirmation.

### Users page

Columns:

```text
Name
Email
Department
Roles
Status
```

Actions:

- Create user.
- Edit user.
- Assign roles.
- Deactivate user.

### Departments page

Columns:

```text
Name
Code
User Count
```

### Audit logs page

Columns:

```text
Timestamp
Actor
Action
Entity Type
Entity ID
```

Filters:

```text
Action
Entity Type
Actor
Date range
```

### Exports page

Export options:

- Purchase requests CSV.
- Approval history CSV.
- Audit logs CSV.

Filters:

- Date range.
- Status.
- Department.

## 4. Status labels

Use clear status chips.

Request statuses:

```text
Draft
Submitted
In Review
Revision Requested
Approved
Rejected
Processing
Completed
Cancelled
```

Approval step statuses:

```text
Waiting
Pending
Approved
Rejected
Revision Requested
Skipped
```

## 5. Empty states

Examples:

- No purchase requests yet. Create your first purchase request.
- No pending approvals. You are all caught up.
- No workflow configured. Create one manually or generate with AI.

## 6. MVP UI shortcuts

Allowed in v0.1:

- JSON editor instead of visual workflow builder.
- Basic tables instead of advanced dashboards.
- Simple upload input instead of drag-and-drop.
- Mermaid diagram text rendering instead of custom graph engine.

## 7. Suggested frontend structure

```text
apps/web/
  app/
    login/
    setup/
    dashboard/
    purchase-requests/
      page.tsx
      new/page.tsx
      [id]/page.tsx
    approvals/
    workflows/
      page.tsx
      [id]/page.tsx
    ai-workflow-builder/
    users/
    departments/
    audit-logs/
    exports/
    settings/
  components/
    layout/
    tables/
    forms/
    status-chip.tsx
    workflow-diagram.tsx
    json-editor.tsx
  lib/
    api-client.ts
    auth.ts
    formatting.ts
  types/
    api.ts
```
