# Operra

Operra is a multi-tenant, self-hosted approval workflow platform.

Version v0.1 focuses on AI-assisted Purchase Request approval.

## Repository Layout

```text
operra/
  AGENTS.md
  README.md
  .env.example
  apps/
    api/
    web/
  docs/
```

## Current Scope

This repository is being initialized following the project docs and Codex task plan.

Initial setup targets:

1. Monorepo scaffold for `apps/api` and `apps/web`.
2. Deployment environment variables in `.env.example`.
3. Docker Compose and application implementation in later tasks.

## Reference Docs

Read these in order before implementing features:

1. `docs/prd.md`
2. `docs/architecture.md`
3. `docs/data-model.md`
4. `docs/workflow-engine.md`
5. `docs/api.md`
6. `docs/security.md`
7. `docs/testing.md`

Use `docs/codex-task-plan.md` as the implementation sequence.
