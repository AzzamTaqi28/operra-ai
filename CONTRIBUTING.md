# Contributing to Operra

Operra is a multi-tenant, self-hosted approval workflow platform.

## Project principles

- Keep the backend modular and tenant-scoped.
- Keep workflow execution deterministic.
- Keep the frontend operational and simple.
- Do not add drag-and-drop workflow builders in v0.1.
- Do not add autonomous AI behavior.

## Before contributing

Read these docs first:

1. `docs/prd.md`
2. `docs/architecture.md`
3. `docs/data-model.md`
4. `docs/workflow-engine.md`
5. `docs/api.md`
6. `docs/security.md`
7. `docs/testing.md`
8. `docs/roadmap.md`

## Code standards

- Prefer readable code over clever abstractions.
- Keep tenant checks explicit in backend queries.
- Add tests for permission rules and workflow transitions.
- Do not commit secrets or real credentials.
- Keep docs and code aligned.

## Pull request expectations

- Describe the user-facing change.
- Mention any schema or API changes.
- Include tests when behavior changes.
- Update docs if the product or public setup changes.

