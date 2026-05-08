CREATE TABLE purchase_requests (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE approval_step_instances (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (organization_id, purchase_request_id, step_key)
);

CREATE TABLE approval_actions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID NOT NULL REFERENCES organizations(id),
  purchase_request_id UUID NOT NULL REFERENCES purchase_requests(id),
  approval_step_instance_id UUID NOT NULL REFERENCES approval_step_instances(id),
  actor_user_id UUID NOT NULL REFERENCES users(id),
  action TEXT NOT NULL,
  comment TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE comments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID NOT NULL REFERENCES organizations(id),
  purchase_request_id UUID NOT NULL REFERENCES purchase_requests(id),
  actor_user_id UUID NOT NULL REFERENCES users(id),
  body TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_purchase_requests_org_status ON purchase_requests (organization_id, status);
CREATE INDEX idx_purchase_requests_org_requester ON purchase_requests (organization_id, requester_id);
CREATE INDEX idx_purchase_requests_org_department ON purchase_requests (organization_id, department_id);
CREATE INDEX idx_purchase_requests_org_created ON purchase_requests (organization_id, created_at);
CREATE INDEX idx_approval_steps_org_request ON approval_step_instances (organization_id, purchase_request_id);
CREATE INDEX idx_approval_steps_org_status ON approval_step_instances (organization_id, status);
CREATE INDEX idx_approval_actions_org_request ON approval_actions (organization_id, purchase_request_id);
