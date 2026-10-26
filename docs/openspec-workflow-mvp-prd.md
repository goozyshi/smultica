# Multica OpenSpec Workflow MVP

## 1. Summary

Add an opt-in, OpenSpec-backed workflow to Multica without changing the existing Issue, Agent, Runtime, Run, Review, or Inbox experience.

The workflow helps a traditional product and engineering team move from discussion to an agreed specification, implementation, developer self-test, review, and QA handoff. OpenSpec files in the repository remain the source of truth; Multica provides the collaboration, orchestration, visibility, and audit layer.

## 2. Goals

- Reuse existing `openspec/` artifacts and OpenSpec CLI behavior.
- Allow a discussion to mature into an OpenSpec Change.
- Require explicit agreement before implementation begins.
- Reuse Multica Agents, Issues, Runtimes, Runs, Reviews, and Inbox.
- Support controlled requirement changes without destroying history.
- Provide developer self-test evidence while keeping full testing independent.
- Keep the UI visually and behaviorally consistent with Multica.

## 3. Non-goals for MVP

- Full test-case management or replacement of the independent QA system.
- Feishu, Slack, or other IM integrations.
- A replacement OpenSpec CLI or custom OpenSpec schema editor.
- Mandatory Worktree support.
- Mobile editing of complex Markdown or Git diffs.

## 4. Core concepts

### Discuss Room

A focused conversation for exploring a requirement, resolving design questions, or discussing a change during implementation. A room may exist before a Change and can later be converted into one.

The room stores conversation, decisions, unresolved questions, source references, and a link to its Change. Full conversation remains in Multica; only confirmed decisions and relevant summaries are written to OpenSpec artifacts.

### Change

The unit of OpenSpec work. A Change maps to one OpenSpec change directory and may contain multiple versions caused by approved requirement changes.

```text
openspec/changes/<change-id>/
  proposal.md
  specs/
  design.md
  tasks.md
```

### Change version

A versioned snapshot of the Change's intent and plan. An approved change request creates a new version or a linked Change when the scope or external contract changes materially.

### QA Handoff

The delivery package created after implementation and developer self-test. It summarizes the requirement, changed files, OpenSpec acceptance scenarios, self-test evidence, risks, and suggested independent test scope.

## 5. Workflow state machine

```text
Discuss
  ↓
Proposal → Specs → Design → Tasks
  ↓                         ↓
Approved                 Apply / Implement
                            ↓
                       Self-test
                            ↓
                         Review
                            ↓
                       QA Handoff
                            ↓
                  Independent testing
                            ↓
                         Archive
```

Every transition records actor, timestamp, source artifact version, and decision. A stage may be marked `Blocked`, `Parked`, or `Needs change` without losing its history.

### Default gates

1. Proposal and relevant Specs must be approved before Design/Tasks are considered executable.
2. Design and Tasks must be approved before the Agent can start affected implementation tasks.
3. Self-test evidence is required before Review.
4. Review approval is required before QA Handoff.
5. Archive requires completed required tasks, no unresolved artifact conflicts, and explicit user action.

Projects may enable a lightweight workflow that skips selected stages, but every skip records who skipped it and why.

## 6. Requirement change protocol

When a requirement changes, Multica creates a Change Request instead of silently editing the active plan.

```text
Change Request
  ↓
Impact analysis
  ↓
Minor: amend current version
Moderate: create new version
Major: create linked Change
  ↓
Approval
  ↓
Re-plan affected tasks
```

Impact analysis should identify affected Specs, Design sections, tasks, code areas, running Agent tasks, and acceptance scenarios.

Completed tasks are not automatically reverted. They are marked `Needs revalidation`, `Needs rework`, or `Obsolete`; a user decides whether to reopen them.

Archived Changes are immutable. Follow-up work creates a linked Change such as `Fixes archived Change: <id>`.

## 7. Mapping to existing Multica objects

| OpenSpec concept | Multica representation |
| --- | --- |
| Discuss Room | New discussion entity using existing chat, comments, mentions, and Inbox patterns |
| Change | New workflow entity linked to Project and repository |
| Proposal / Specs / Design | Versioned Markdown artifacts backed by repository files |
| Task | Existing Issue with stable OpenSpec task ID |
| Agent execution | Existing Agent Run |
| Approval | Existing review/permission model plus workflow gate |
| Developer self-test | Run result and evidence attached to Change/QA Handoff |
| QA Handoff | New delivery node linked to Change and Issues |
| Archive | OpenSpec archive action plus immutable Multica history |

One OpenSpec task maps to one primary Multica Issue. An Issue may have multiple Agent Runs, retries, comments, and review events.

## 8. UI information architecture

### Entry points

- Existing project page: `Workflow` section when `openspec/` is detected.
- Existing creation actions: add `New Change` beside `New Issue`.
- Existing chat: `Create OpenSpec Change` action from a Discuss Room.

### Change page

Use existing Multica layout, cards, badges, timeline, Inbox, and review patterns. Add:

- Stage rail: Discuss, Proposal, Specs, Design, Tasks, Implement, Self-test, Review, QA, Archive.
- Artifact panel: Markdown edit/preview, Git status, diff, version, source path.
- Task panel: existing Issues with OpenSpec IDs, dependencies, Agent assignment, and progress.
- Activity timeline: discussions, Agent Runs, approvals, test evidence, Git/PR links, and change requests.
- Gate panel: pending approvals, blocked reasons, skipped stages, and required actions.

### Discuss Room

- Focused topic and context scope.
- Code and Spec references.
- One active decision question at a time.
- Agent suggestions with alternatives and trade-offs.
- Decision records and unresolved questions.
- `Create OpenSpec Change` action with a preview of generated artifacts.
- `Park` and `Resume` actions.

## 9. OpenSpec integration

- Detect `openspec/`, project configuration, and `openspec/changes/`.
- Read existing artifacts without copying or rewriting them.
- Prefer the OpenSpec CLI installed in the selected Runtime.
- Fall back to an Agent skill only when the CLI is unavailable.
- Preserve project-selected OpenSpec version/schema.
- Do not silently overwrite files changed outside Multica; show a diff and require a merge decision.
- Use official OpenSpec operations for propose/apply/archive where available.

## 10. Agent behavior

Agents may:

- Scan relevant code and existing Specs during Discuss and impact analysis.
- Generate or update artifacts only within the active gate.
- Create or update linked Issues from `tasks.md`.
- Execute implementation tasks in the selected Runtime.
- Run configured developer self-tests.
- Produce evidence and report unresolved risks.

Agents may not:

- Approve their own artifacts or Review gate.
- Archive a Change.
- Silently change a requirement or overwrite an external Git edit.

## 11. Developer self-test

Self-test commands come from project configuration first, with a Multica workflow override available:

```yaml
workflow:
  self_test:
    - pnpm typecheck
    - pnpm lint
    - pnpm test
```

The system should also detect conventional commands from `package.json`, `Makefile`, `pyproject.toml`, `go.mod`, and similar project files.

The UI must distinguish:

```text
Developer self-test: Passed
Independent testing: Pending
```

## 12. Permissions

- Member: discuss, comment, edit assigned artifacts, and request Agent work.
- Developer/lead: approve implementation-ready artifacts and Review.
- Admin/Owner: configure workflow, skip gates, approve major changes, and Archive.
- Agent: propose and modify artifacts within its active task, but never approve itself.

## 13. MVP acceptance criteria

Given a repository with an existing `openspec/` directory, a user can:

1. Open the project and discover active Changes.
2. Create or open a Discuss Room.
3. Convert a discussion into an OpenSpec Change.
4. Review and approve generated Proposal, Specs, Design, and Tasks.
5. See tasks as Multica Issues and assign them to Agents.
6. Run implementation and developer self-tests.
7. Review changes and generate a QA Handoff.
8. Create a Change Request, inspect impact, and approve a new version.
9. Archive the completed Change through the OpenSpec-compatible flow.
10. Continue using normal Multica Issues without entering Workflow mode.

## 14. Delivery phases

### Phase 1 — Native Workflow foundation

- Change discovery and creation.
- Artifact viewer/editor with Git diff.
- Stage rail and approval gates.
- Task-to-Issue mapping.
- Existing Agent Run integration.

### Phase 2 — Discuss and change management

- Standalone Discuss Room.
- Decision records.
- Change Request and impact analysis.
- Versioning, Park/Resume, and conflict handling.

### Phase 3 — Delivery quality

- Self-test configuration and evidence.
- QA Handoff.
- PR/Commit/CI links.
- Optional Worktree and dependency-aware parallel tasks.

### Phase 4 — External conversation adapters

- Feishu first, then other IM platforms.
- Topic binding, mentions, approvals, and summaries.
- External messages remain an entry point; OpenSpec/Git remains the source of truth.
