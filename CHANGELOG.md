# Changelog

All notable changes to Operra will be documented in this file.

## Unreleased

- Public OSS repository cleanup
- Public contributor, support, security, and governance docs
- Workflow chart and public roadmap

## 0.1.0

Initial v0.1 release focused on AI-assisted Purchase Request approval.

Includes:

- Multi-tenant organization setup
- Users, departments, and roles
- JSON-first workflow configuration
- AI-assisted workflow generation
- Purchase request drafting and submission
- Role-based approvals
- Comments and attachments
- Audit logs
- CSV exports
- Next.js web app
- Docker Compose deployment

Also included:

- Public contributor, support, security, and governance docs
- Public roadmap and workflow chart
- Public OSS repository cleanup

Notes:

- Workflow execution is deterministic.
- AI generates workflow JSON, but the backend validates and controls execution.
- Attachments use S3-compatible storage with MinIO for local/self-hosted deployment.
