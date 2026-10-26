# OpenSpec Workflow MVP — Implementation Plan

## 1. Architecture decision

OpenSpec is a local-repository concern. The API server must not assume it can read a user's checkout.

```text
Web UI
  ↕ HTTP / realtime
Multica API
  ├── Change metadata, gates, versions, links, audit
  └── Issue/Agent/Run orchestration
       ↕ daemon task / runtime protocol
Local Daemon
  ├── project resource / work directory
  ├── OpenSpec CLI discovery
  ├── artifact scan and Git diff
  └── OpenSpec apply/archive/self-test commands
       ↕
Git repository
  └── openspec/changes/<change-id>/...
```

The repository remains the source of truth for OpenSpec files. The server stores workflow state and a digest/version summary, not an authoritative copy of the artifacts.

## 2. Existing code to reuse

### Existing server capabilities

- `project` and `project_resource` already associate a Project with a repository or local directory.
- `agent_task_queue`, Issue handlers, task runs, comments, and realtime events already provide execution and progress infrastructure.
- `chat_session` and chat task execution can provide the first Agent-assisted discussion path.
- Project workflow code under `packages/core/issue-workflows` and `packages/views/workflows` is an existing Issue status/handoff system. It should not be reused as the domain model for OpenSpec; it can only inspire the stage rail UI.
- `project_resource.local_directory` already handles daemon-local paths and optional worktree execution.

### Existing client capabilities

- `packages/core/api/client.ts`, `schemas.ts`, and `types/` are the API contract layer.
- `packages/core/projects`, `packages/core/chat`, and `packages/core/issues` provide query/mutation/cache patterns.
- `packages/views/projects/components/project-detail.tsx` already has a project detail surface, resource section, right sidebar, and Issue board.
- `packages/views/editor` provides the existing content editor and Markdown-related UI.
- `packages/views/common/task-transcript` provides reusable Agent Run evidence rendering.

## 3. Naming

Avoid `issue_workflow` and `workflow` as the primary server domain names because they already mean Issue status/handoff workflows.

Recommended names:

- SQL tables: `spec_change`, `spec_change_version`, `spec_artifact`, `spec_gate`, `spec_task_link`, `spec_self_test`
- API namespace: `/api/spec-changes`
- TypeScript feature: `packages/core/spec-changes` and `packages/views/spec-changes`
- UI label: `OpenSpec Workflow` or `Changes`
- Existing Issue workflow remains `Issue Workflow` internally and in code.

## 4. MVP vertical slices

### Slice 0 — Feature flag and read-only discovery

Goal: prove the repository/runtime boundary without changing existing Issue behavior.

Server:

- Add a workspace/project feature flag for OpenSpec Workflow.
- Add a minimal `spec_change` record linked to workspace, project, and project resource.
- Add an API to list registered/discovered Changes.

Daemon/runtime:

- Add a read-only `spec discover` operation to the existing runtime task protocol.
- Given a `local_directory` resource, detect `openspec/`, `openspec/changes/`, and known artifact files.
- Return change IDs, artifact paths, file digests, and basic parse status.

Web:

- Add a `Changes` section to Project Detail.
- Show “OpenSpec detected”, “No OpenSpec directory”, or “Runtime unavailable”.
- Show existing Changes without editing them.

Acceptance: opening a Project with existing OpenSpec artifacts shows them in Multica, while normal Issues and the existing Issue Workflow remain unchanged.

### Slice 1 — Change detail and artifact viewer

Server:

- Add `spec_change_version` and `spec_artifact` metadata.
- Store path, artifact kind, digest, detected status, and last scan time.
- Add versioned audit events for discovery and user actions.

Daemon:

- Add read-only artifact fetch and Git diff operations.
- Never allow the API server to read arbitrary local paths.

Web:

- Add a Change detail route under the existing workspace/project navigation.
- Add stage rail: Discuss, Proposal, Specs, Design, Tasks, Implement, Self-test, Review, QA, Archive.
- Add Markdown preview, file path, digest, Git status, and diff view.

Acceptance: a user can inspect an existing Change and see whether the local checkout differs from the last known scan.

### Slice 2 — Create Change and OpenSpec generation

Server:

- Add `POST /api/spec-changes` with project/resource and requested name.
- Add a Change command/task request, not direct filesystem writes.
- Add gate state and actor/audit records.

Daemon:

- Detect the installed OpenSpec CLI and report version/profile.
- Run the official propose operation in the selected project directory.
- Return created paths and digests.
- If CLI is missing, return an actionable error; do not silently pretend to be OpenSpec.

Web:

- Add `New Change` beside existing Project/Issue creation actions.
- Preview title, scope, repository, runtime, and generated artifact paths before execution.
- Reuse existing Agent Run/task transcript UI for command output.

Acceptance: a new Change creates real OpenSpec artifacts in the user's checkout and is visible in Multica.

### Slice 3 — Tasks to Issues

Parser/daemon:

- Parse `tasks.md` into stable task IDs, labels, dependencies, and checked state.
- Treat task IDs as stable identifiers; never use list position alone.

Server:

- Add `spec_task_link(change_id, task_id, issue_id, artifact_version, state)`.
- Create/link Issues idempotently.
- Store the OpenSpec task ID in Issue metadata until a first-class field is justified.

Web:

- Show Tasks using existing Issue cards/board and link back to the Change.
- Display dependency and gate state.
- Reuse Agent assignment and Run controls.

Acceptance: each OpenSpec task maps to one primary Multica Issue, repeated sync does not duplicate Issues, and completed Issue state can be reconciled with `tasks.md`.

### Slice 4 — Gates, approvals, and self-test

Server:

- Implement the Change state machine and gate transitions.
- Enforce permissions server-side; Agent actors cannot approve their own artifacts.
- Add `spec_self_test` records for command, runtime, task/run ID, exit status, output reference, and commit/digest.

Runtime:

- Read configured self-test commands from project configuration.
- Execute tests through the existing Agent/Daemon task path.
- Return structured results and bounded logs.

Web:

- Add approval actions to the stage rail and artifact panel.
- Add self-test results with explicit `Developer self-test` labeling.
- Keep `Independent testing: Pending` separate from self-test status.

Acceptance: implementation cannot pass the Review gate without configured or explicitly waived self-test evidence.

### Slice 5 — Change requests and QA Handoff

Server:

- Add `spec_change_request` and impact result records.
- Support Minor amend, Moderate new version, and Major linked Change.
- Mark affected completed tasks as `needs_revalidation`, `needs_rework`, or `obsolete`; never silently revert them.
- Add QA Handoff data and status.

Runtime:

- Scan the repository and OpenSpec artifacts for an impact-analysis input.
- Produce references to affected artifact paths and optional code areas.

Web:

- Add Change Request creation from any active stage.
- Display impact summary before approval.
- Generate QA Handoff with requirements, changed files, acceptance scenarios, self-test evidence, and risks.

Acceptance: a requirement change creates an auditable version/linked Change and a tester receives a clear handoff without being forced into the implementation workflow.

## 5. Suggested first migrations

Keep the first schema additive and narrow:

```sql
spec_change
  id, workspace_id, project_id, resource_id,
  change_key, title, status, current_version,
  created_by, created_at, updated_at

spec_artifact
  id, change_id, version, kind, relative_path,
  digest, status, scanned_at

spec_task_link
  id, change_id, task_id, issue_id,
  artifact_version, state, updated_at
```

Add gates, self-tests, versions, and change requests in later migrations once the read/create/task vertical slice is proven.

## 6. First API surface

```text
GET    /api/projects/{projectId}/spec-changes
POST   /api/spec-changes
GET    /api/spec-changes/{id}
POST   /api/spec-changes/{id}/scan
GET    /api/spec-changes/{id}/artifacts
GET    /api/spec-changes/{id}/artifacts/{artifactId}/content
GET    /api/spec-changes/{id}/artifacts/{artifactId}/diff
POST   /api/spec-changes/{id}/propose
POST   /api/spec-changes/{id}/sync-tasks
POST   /api/spec-changes/{id}/gates/{gate}/approve
POST   /api/spec-changes/{id}/self-test
```

The first implementation should keep endpoints explicit rather than exposing a generic “run arbitrary OpenSpec command” API.

## 7. Recommended implementation order

1. Add feature flag, schema names, types, and API contract tests.
2. Build read-only daemon discovery and project Change list.
3. Build Change detail with artifact metadata/content/diff.
4. Add Change creation through the official OpenSpec CLI.
5. Sync `tasks.md` into existing Issues.
6. Add approval gates and task-start guards.
7. Add self-test evidence and QA Handoff.
8. Add Change Request/versioning and Park/Resume.
9. Add optional Worktree/parallel tasks.
10. Add Feishu/other IM adapters only after the native Discuss Room and Change lifecycle are stable.

## 8. Review checkpoints

After each slice run:

- Go handler/service/migration tests.
- API schema and query tests.
- Core query/mutation tests.
- Web component and route tests.
- A Docker self-host smoke test with a real local repository containing a small `openspec/` fixture.

No slice should change the behavior of ordinary Issues, existing Issue Workflows, or existing Agent Runs when the OpenSpec feature flag is off.

## 9. Concrete PR breakdown

### PR 0 — Contract and terminology

Documentation and types only:

- Confirm `spec_change` naming does not collide with `issue_workflow`.
- Add shared TypeScript enums/types for Change status, artifact kind, and gate state.
- Add API schemas with empty/fallback response fixtures.
- Add a feature-flag key with the default disabled.
- Add fixture repositories containing existing and malformed OpenSpec changes.

Primary areas:

```text
packages/core/types/
packages/core/api/schemas.ts
packages/core/feature-flags/
docs/
```

### PR 1 — Discovery backend contract

Server-only read path:

- Add the minimal `spec_change` and `spec_artifact` migrations.
- Add sqlc queries and generated models.
- Add workspace/project authorization helpers.
- Add `GET /api/projects/{projectId}/spec-changes`.
- Add the daemon request/response contract for discovery.
- Add handler, service, and authorization tests.

Primary areas:

```text
server/migrations/
server/pkg/db/queries/
server/internal/handler/
server/internal/service/
server/cmd/server/router.go
server/pkg/protocol/
```

### PR 2 — Local discovery in the Daemon

- Resolve the project's `local_directory` resource.
- Scan only the repository root and `openspec/` paths.
- Return change IDs, artifact paths, file digests, and parse errors.
- Reject paths outside the selected project resource.
- Make discovery idempotent and safe when the repository has no OpenSpec directory.

Primary areas:

```text
server/internal/daemon/
server/internal/handler/daemon_*.go
server/pkg/protocol/
```

### PR 3 — Project UI and read-only Change detail

- Add a `Changes` section beside Project resources and Issues.
- Add a Change detail route under the existing workspace/project navigation.
- Add stage rail and artifact list.
- Add loading, empty, unavailable-runtime, malformed-artifact, and permission states.
- Reuse existing Markdown/editor, Sheet/sidebar, Badge, Timeline, and Run components.

Primary areas:

```text
packages/core/spec-changes/
packages/core/api/client.ts
packages/views/spec-changes/
packages/views/projects/components/project-detail.tsx
apps/web/app/[workspaceSlug]/(dashboard)/projects/
packages/views/locales/*/projects.json
```

### PR 4 — Artifact content and diff

- Add daemon-backed artifact content and Git diff requests.
- Add digest mismatch handling.
- Render Markdown preview and raw source.
- Never silently overwrite a changed file.

### PR 5 — Create Change through OpenSpec

- Add `New Change` modal and preview.
- Dispatch the official OpenSpec propose command through the Daemon.
- Store returned paths/digests and refresh the Change view.
- Surface missing CLI/profile as an actionable runtime error.

### PR 6 — Task synchronization

- Parse stable task IDs from `tasks.md`.
- Add `spec_task_link` and idempotent Issue creation/linking.
- Reuse Issue assignment, status, comments, Agent Runs, and realtime events.
- Add drift detection when `tasks.md` changes outside Multica.

### PR 7 — Gates and delivery quality

- Add approval gates and server-side permissions.
- Add self-test configuration and evidence.
- Add QA Handoff.
- Add Review/Archive guards.

### PR 8 — Requirement changes

- Add Change Request and impact analysis.
- Add versioning and linked Change relationships.
- Add Park/Resume.
- Add revalidation/rework markers for completed tasks.

## 10. First implementation checklist

The first coding slice is complete only when all of the following are true:

- [ ] Feature is disabled by default and ordinary Multica behavior is unchanged.
- [ ] A Project with a local directory can be selected as an OpenSpec source.
- [ ] Existing `openspec/changes/<id>` directories are discovered by the local Daemon.
- [ ] Discovery returns deterministic digests and relative paths only.
- [ ] The API returns discovered Changes with workspace/project authorization.
- [ ] The Project page shows a read-only Changes section.
- [ ] Existing malformed or incomplete artifacts render a warning instead of failing the page.
- [ ] No OpenSpec CLI is invoked in the discovery-only slice.
- [ ] Go, TypeScript, API contract, and UI tests cover the happy path and missing-runtime path.
