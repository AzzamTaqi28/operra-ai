# Operra Testing Plan

## 1. Testing goals

v0.1 must prove:

- Multi-tenancy is safe.
- Workflow execution is correct.
- Role-based approvals work.
- Attachments are permission-protected.
- CSV exports are correct and audited.
- AI output is validated before saving.

## 2. Backend unit tests

### Workflow validator

Test cases:

1. Valid purchase request workflow passes.
2. Missing workflow name fails.
3. Empty steps fails.
4. Duplicate step keys fail.
5. Unsupported approver role fails.
6. Unsupported condition field fails.
7. Unsupported condition operator fails.
8. Condition value type mismatch fails.

### Condition evaluator

Test cases:

1. `estimated_amount > 5000000` returns true for 10000000.
2. `estimated_amount > 5000000` returns false for 1000000.
3. `urgency == urgent` works.
4. Unknown field returns validation error, not silent false.

### Approval engine

Test cases:

1. Amount 1,000,000 generates manager + procurement.
2. Amount 10,000,000 generates manager + finance + procurement.
3. Amount 30,000,000 generates manager + finance + director + procurement.
4. Approving current step advances to next step.
5. Final approval marks request completed or approved based on procurement step rule.
6. Reject marks request rejected.
7. Request revision marks request revision_requested.
8. Resubmission restarts approval.

## 3. Permission tests

Test cases:

1. User cannot access another organization's purchase request.
2. User cannot approve another organization's request.
3. Requester cannot approve own request.
4. User without manager role cannot approve manager step.
5. Finance cannot approve manager step.
6. Manager from wrong department cannot approve requester_department step.
7. Auditor cannot approve or edit.
8. Admin can manage workflows.

## 4. API integration tests

Critical flows:

### Organization setup

1. Register organization.
2. Login as owner.
3. Verify built-in roles exist.

### Purchase request flow below threshold

1. Create requester, manager, procurement users.
2. Create default workflow.
3. Submit request amount 1,000,000.
4. Manager approves.
5. Procurement approves/processes.
6. Request becomes completed.
7. Audit logs exist.

### Purchase request flow high amount

1. Submit request amount 30,000,000.
2. Manager approves.
3. Finance approves.
4. Director approves.
5. Procurement processes.
6. Request becomes completed.

### CSV export

1. Create requests.
2. Export purchase requests CSV.
3. Verify headers.
4. Verify row count.
5. Verify audit log `csv.exported`.

### Attachments

1. Upload attachment.
2. Metadata saved.
3. File exists in MinIO/S3.
4. Authorized user can download.
5. Unauthorized user cannot download.
6. Audit log exists.

### AI workflow generation

1. Submit valid AI prompt.
2. AI returns JSON.
3. Backend validates JSON.
4. Save workflow after confirmation.
5. Invalid JSON cannot be saved.

## 5. Frontend tests

Minimum UI tests:

1. Login form works.
2. Purchase request form validates required fields.
3. Request list displays status.
4. Request detail displays approval timeline.
5. Approver can approve from detail page.
6. Workflow JSON editor displays validation errors.
7. AI workflow builder shows generated JSON and diagram.
8. CSV export button downloads file for permitted users.

## 6. Manual acceptance tests

Before private beta, manually test:

1. Fresh `docker compose up -d` works.
2. Register organization.
3. Create departments.
4. Create users and roles.
5. Generate workflow with AI.
6. Validate and activate workflow.
7. Submit purchase request with attachment.
8. Approve through all required roles.
9. Export CSV reports.
10. Review audit logs.
11. Confirm no cross-organization data leak with a second demo organization.

## 7. Definition of done for v0.1

- All critical backend tests pass.
- Docker Compose works from clean setup.
- No known tenant isolation bug.
- Workflow engine passes amount threshold cases.
- CSV exports work and are audited.
- Attachments work with MinIO.
- AI workflow generation is validated before saving.
- README setup instructions are accurate.
