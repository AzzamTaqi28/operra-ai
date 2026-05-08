# Operra Data Model

## 1. Database principles

- Use PostgreSQL.
- Use UUID primary keys.
- Use `organization_id` on every tenant-owned table.
- Use `created_at` and `updated_at` where appropriate.
- Use soft delete only where useful; do not overuse it.
- Use foreign keys for critical relationships.
- Add indexes for common filters.

## 2. Multi-tenancy rule

Every tenant-owned table must include:

```text
organization_id UUID NOT NULL
```

Every query for tenant-owned data must include `organization_id`.

## 3. Core tables

### organizations

Purpose: tenant/company records.

```sql
CREATE TABLE organizations (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  status TEXT NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
```

### users

Purpose: users belonging to an organization.

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id),
  department_id UUID NULL,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  UNIQUE (organization_id, email)
);
```

Recommended statuses:

```text
active
inactive
invited
```

### departments

```sql
CREATE TABLE departments (
  id UUID PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id),
  name TEXT NOT NULL,
  code TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  UNIQUE (organization_id, name)
);
```

### roles

Roles can be built-in per organization.

```sql
CREATE TABLE roles (
  id UUID PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id),
  key TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT NULL,
  is_system BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  UNIQUE (organization_id, key)
);
```

Built-in role keys:

```text
owner
admin
requester
manager
finance
procurement
director
auditor
```

### user_roles

```sql
CREATE TABLE user_roles (
  id UUID PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id),
  user_id UUID NOT NULL REFERENCES users(id),
  role_id UUID NOT NULL REFERENCES roles(id),
  created_at TIMESTAMPTZ NOT NULL,
  UNIQUE (organization_id, user_id, role_id)
);
```

## 4. Workflow tables

### workflows

Logical workflow record.

```sql
CREATE TABLE workflows (
  id UUID PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id),
  name TEXT NOT NULL,
  type TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'draft',
  active_version_id UUID NULL,
  created_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
```

Workflow statuses:

```text
draft
active
inactive
```

### workflow_versions

Immutable config version.

```sql
CREATE TABLE workflow_versions (
  id UUID PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id),
  workflow_id UUID NOT NULL REFERENCES workflows(id),
  version_number INTEGER NOT NULL,
  config_json JSONB NOT NULL,
  mermaid_diagram TEXT NULL,
  explanation TEXT NULL,
  created_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL,
  UNIQUE (organization_id, workflow_id, version_number)
);
```

Note:

- Do not update `config_json` after creation.
- Create a new version when workflow config changes.

## 5. Purchase request tables

### purchase_requests

```sql
CREATE TABLE purchase_requests (
  id UUID PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id),
  workflow_id UUID NULL REFERENCES workflows(id),
  workflow_version_id UUID NULL REFERENCES workflow_versions(id),
  requester_id UUID NOT NULL REFERENCES users(id),
  department_id UUID NOT NULL REFERENCES departments(id),
  title TEXT NOT NULL,
  item_name TEXT NOT NULL,
  description TEXT NOT NULL,
  quantity NUMERIC(18,2) NOT NULL,
  estimated_amount NUMERIC(18,2) NOT NULL,
  currency TEXT NOT NULL DEFAULT 'IDR',
  urgency TEXT NOT NULL DEFAULT 'normal',
  expected_date DATE NULL,
  vendor_name TEXT NULL,
  notes TEXT NULL,
  status TEXT NOT NULL DEFAULT 'draft',
  current_step_instance_id UUID NULL,
  submitted_at TIMESTAMPTZ NULL,
  completed_at TIMESTAMPTZ NULL,
  cancelled_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
```

Urgency values:

```text
low
normal
high
urgent
```

Status values:

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

### approval_step_instances

Runtime approval steps for a request.

```sql
CREATE TABLE approval_step_instances (
  id UUID PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id),
  purchase_request_id UUID NOT NULL REFERENCES purchase_requests(id),
  workflow_version_id UUID NOT NULL REFERENCES workflow_versions(id),
  step_key TEXT NOT NULL,
  step_name TEXT NOT NULL,
  sequence_number INTEGER NOT NULL,
  approver_role_key TEXT NOT NULL,
  scope TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'waiting',
  started_at TIMESTAMPTZ NULL,
  completed_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  UNIQUE (organization_id, purchase_request_id, step_key)
);
```

Step statuses:

```text
waiting
pending
approved
rejected
revision_requested
skipped
```

### approval_actions

Action history for approvals.

```sql
CREATE TABLE approval_actions (
  id UUID PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id),
  purchase_request_id UUID NOT NULL REFERENCES purchase_requests(id),
  approval_step_instance_id UUID NOT NULL REFERENCES approval_step_instances(id),
  actor_user_id UUID NOT NULL REFERENCES users(id),
  action TEXT NOT NULL,
  comment TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL
);
```

Actions:

```text
approve
reject
request_revision
```

### comments

General request comments.

```sql
CREATE TABLE comments (
  id UUID PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id),
  purchase_request_id UUID NOT NULL REFERENCES purchase_requests(id),
  actor_user_id UUID NOT NULL REFERENCES users(id),
  body TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
```

## 6. Attachment tables

### attachments

```sql
CREATE TABLE attachments (
  id UUID PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id),
  purchase_request_id UUID NOT NULL REFERENCES purchase_requests(id),
  uploaded_by UUID NOT NULL REFERENCES users(id),
  file_name TEXT NOT NULL,
  file_size BIGINT NOT NULL,
  mime_type TEXT NOT NULL,
  storage_driver TEXT NOT NULL DEFAULT 's3',
  storage_bucket TEXT NOT NULL,
  storage_key TEXT NOT NULL,
  checksum TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL
);
```

## 7. Audit and export tables

### audit_logs

```sql
CREATE TABLE audit_logs (
  id UUID PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id),
  actor_user_id UUID NULL REFERENCES users(id),
  action TEXT NOT NULL,
  entity_type TEXT NOT NULL,
  entity_id UUID NULL,
  old_value JSONB NULL,
  new_value JSONB NULL,
  ip_address TEXT NULL,
  user_agent TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL
);
```

### export_logs

Optional table for detailed export tracking. Audit logs are still required.

```sql
CREATE TABLE export_logs (
  id UUID PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id),
  requested_by UUID NOT NULL REFERENCES users(id),
  export_type TEXT NOT NULL,
  filters_json JSONB NULL,
  row_count INTEGER NULL,
  created_at TIMESTAMPTZ NOT NULL
);
```

Export types:

```text
purchase_requests
approval_history
audit_logs
```

## 8. AI tables

### ai_generation_logs

```sql
CREATE TABLE ai_generation_logs (
  id UUID PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id),
  actor_user_id UUID NOT NULL REFERENCES users(id),
  purpose TEXT NOT NULL,
  provider TEXT NOT NULL,
  model TEXT NULL,
  input_prompt TEXT NOT NULL,
  output_json JSONB NULL,
  validation_status TEXT NOT NULL,
  validation_errors JSONB NULL,
  created_at TIMESTAMPTZ NOT NULL
);
```

Purpose values:

```text
workflow_generation
workflow_explanation
diagram_generation
```

Validation statuses:

```text
valid
invalid
error
```

## 9. Notifications table

Optional for v0.1 but useful.

```sql
CREATE TABLE notifications (
  id UUID PRIMARY KEY,
  organization_id UUID NOT NULL REFERENCES organizations(id),
  user_id UUID NOT NULL REFERENCES users(id),
  type TEXT NOT NULL,
  title TEXT NOT NULL,
  body TEXT NOT NULL,
  read_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL
);
```

## 10. Recommended indexes

```sql
CREATE INDEX idx_users_org ON users (organization_id);
CREATE INDEX idx_users_org_email ON users (organization_id, email);
CREATE INDEX idx_departments_org ON departments (organization_id);
CREATE INDEX idx_roles_org_key ON roles (organization_id, key);
CREATE INDEX idx_user_roles_org_user ON user_roles (organization_id, user_id);

CREATE INDEX idx_workflows_org_type ON workflows (organization_id, type);
CREATE INDEX idx_workflow_versions_org_workflow ON workflow_versions (organization_id, workflow_id);

CREATE INDEX idx_purchase_requests_org_status ON purchase_requests (organization_id, status);
CREATE INDEX idx_purchase_requests_org_requester ON purchase_requests (organization_id, requester_id);
CREATE INDEX idx_purchase_requests_org_department ON purchase_requests (organization_id, department_id);
CREATE INDEX idx_purchase_requests_org_created ON purchase_requests (organization_id, created_at);

CREATE INDEX idx_approval_steps_org_request ON approval_step_instances (organization_id, purchase_request_id);
CREATE INDEX idx_approval_steps_org_status ON approval_step_instances (organization_id, status);
CREATE INDEX idx_approval_actions_org_request ON approval_actions (organization_id, purchase_request_id);

CREATE INDEX idx_attachments_org_request ON attachments (organization_id, purchase_request_id);
CREATE INDEX idx_audit_logs_org_created ON audit_logs (organization_id, created_at);
CREATE INDEX idx_audit_logs_org_entity ON audit_logs (organization_id, entity_type, entity_id);
```

## 11. Seed data

For local development, seed:

- One organization: Demo Company.
- Departments: Engineering, Finance, Operations, Procurement.
- Roles: owner, admin, requester, manager, finance, procurement, director, auditor.
- Users for each role.
- Default purchase request workflow.

## 12. Migration rules

- All schema changes must use migrations.
- Do not rely on automatic schema sync in production.
- Migration filenames should be ordered and descriptive.

Example:

```text
000001_create_organizations.up.sql
000001_create_organizations.down.sql
000002_create_users_and_roles.up.sql
000002_create_users_and_roles.down.sql
```

If using GORM AutoMigrate for early dev, still create real migrations before beta.
