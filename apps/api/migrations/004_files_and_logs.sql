CREATE TABLE attachments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE audit_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID NOT NULL REFERENCES organizations(id),
  actor_user_id UUID NULL REFERENCES users(id),
  action TEXT NOT NULL,
  entity_type TEXT NOT NULL,
  entity_id UUID NULL,
  old_value JSONB NULL,
  new_value JSONB NULL,
  ip_address TEXT NULL,
  user_agent TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE export_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID NOT NULL REFERENCES organizations(id),
  requested_by UUID NOT NULL REFERENCES users(id),
  export_type TEXT NOT NULL,
  filters_json JSONB NULL,
  row_count INTEGER NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE ai_generation_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID NOT NULL REFERENCES organizations(id),
  actor_user_id UUID NOT NULL REFERENCES users(id),
  purpose TEXT NOT NULL,
  provider TEXT NOT NULL,
  model TEXT NULL,
  input_prompt TEXT NOT NULL,
  output_json JSONB NULL,
  validation_status TEXT NOT NULL,
  validation_errors JSONB NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID NOT NULL REFERENCES organizations(id),
  user_id UUID NOT NULL REFERENCES users(id),
  type TEXT NOT NULL,
  title TEXT NOT NULL,
  body TEXT NOT NULL,
  read_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_attachments_org_request ON attachments (organization_id, purchase_request_id);
CREATE INDEX idx_audit_logs_org_created ON audit_logs (organization_id, created_at);
CREATE INDEX idx_audit_logs_org_entity ON audit_logs (organization_id, entity_type, entity_id);
