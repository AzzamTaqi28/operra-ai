CREATE TABLE workflows (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID NOT NULL REFERENCES organizations(id),
  name TEXT NOT NULL,
  type TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'draft',
  active_version_id UUID NULL,
  created_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE workflow_versions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID NOT NULL REFERENCES organizations(id),
  workflow_id UUID NOT NULL REFERENCES workflows(id),
  version_number INTEGER NOT NULL,
  config_json JSONB NOT NULL,
  mermaid_diagram TEXT NULL,
  explanation TEXT NULL,
  created_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (organization_id, workflow_id, version_number)
);

CREATE INDEX idx_workflows_org_type ON workflows (organization_id, type);
CREATE INDEX idx_workflows_org_status ON workflows (organization_id, status);
CREATE INDEX idx_workflow_versions_org_workflow ON workflow_versions (organization_id, workflow_id);
