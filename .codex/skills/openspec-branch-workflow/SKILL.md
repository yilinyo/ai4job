---
name: openspec-branch-workflow
description: Use when proposing, applying, or archiving OpenSpec changes in this repository. Enforces project-specific rules for change granularity, PRD requirements, branch creation, implementation on non-master branches, and archive-before-merge discipline.
---

# OpenSpec Branch Workflow

This is a project-level skill for `/Users/yilin/project/ai-for`.

Use this skill whenever you are:

- proposing a new OpenSpec change
- deciding whether a new request should become a new change
- applying an OpenSpec change
- preparing a change for archive
- discussing whether scope should be expanded or split

## Core Rule

In this repository:

- OpenSpec demand definition happens on `master`
- Code implementation must happen on a new non-`master` branch
- A change is only truly complete when:
  - OpenSpec tasks are complete
  - validation is complete
  - code is committed
  - branch is merged back to `master`
  - then archive is appropriate

## When a new change is required

Create a new change when the request is one of:

1. a brand-new module
2. a new feature inside an existing module
3. a bug fix for one function or one behavior
4. removal of one specific function

Do not use a single change for:

- multiple unrelated modules
- multiple unrelated feature groups
- broad mixed scope such as feature + refactor + unrelated bug
- deleting an entire module

## Granularity rules

A good change in this project should usually correspond to one of:

- one module
- one feature
- one bug
- one deletion

If the request spans multiple first-class modules, recommend splitting it.

If deletion affects more than one feature, recommend splitting it.

## PRD and visual requirements

Use this default policy:

- Lightweight change:
  - small bug
  - small interaction fix
  - small wording correction
  - simple requirement note is acceptable

- Standard change:
  - existing module feature
  - flow adjustment
  - page structure change
  - should have a PRD

- Module-level change:
  - new module
  - should have a fuller PRD
  - provide visual support when needed

Do not force a heavy PRD for trivial fixes, but do require clear documentation.

## Apply rules

Before implementing:

1. Confirm the change exists and scope is stable enough.
2. Confirm the current branch is not `master`.
3. If currently on `master`, stop and create a new branch first.
4. If the current branch does not match the permission for the current action, explicitly check and report branch state before doing anything else.

Recommended branch naming:

- `feat/<change-name>`
- `fix/<change-name>`
- `chore/remove-<feature-name>`
- `docs/<change-name>`

Do not implement change code directly on `master`.

Do not commit implementation code or merge a feature branch back to `master` while the corresponding OpenSpec `tasks.md` still has unchecked tasks. If tasks are incomplete, continue implementation or validation first; do not use commit or merge to bypass the task checklist.

## Commit and merge rules

Before committing implementation work on a feature branch:

1. Confirm the corresponding OpenSpec change is known.
2. Confirm `tasks.md` has no unchecked implementation or validation tasks.
3. If any task remains unchecked, stop and finish or explicitly de-scope the task in the artifacts before committing.

Before merging a feature branch back to `master`:

1. Re-check the corresponding `tasks.md`.
2. Confirm all tasks are complete.
3. If any task remains unchecked, do not merge to `master`.

## Branch permission enforcement

This repository requires explicit branch-state checks when branch context and action type may conflict.

If the current branch does not satisfy the current action's permission, you must:

1. explicitly check current branch state
2. explicitly tell the user what branch you are on
3. explicitly explain why this branch does not match the required workflow
4. stop the operation until the branch issue is resolved or the user gives a clear exception

Do not silently continue.

Do not assume branch state is acceptable.

Do not move directly into implementation, archive, or merge-oriented actions without branch confirmation when the workflow requires it.

Do not move into commit or merge actions for an OpenSpec apply if the change tasks are incomplete.

Examples:

- On `master`, user asks to apply a change:
  - you must check branch state
  - you must state that `master` is not allowed for code implementation
  - you must request or perform branch creation before continuing

- On a feature branch, user asks to finalize archive:
  - you must check whether merge-back to `master` has happened
  - if not, do not present the change as fully workflow-complete

- On an unknown branch, user asks to continue OpenSpec work:
  - check branch first
  - do not guess

## Archive rules

Do not recommend archive unless all of the following are true:

1. `tasks.md` implementation tasks are complete
2. validation tasks are complete
3. OpenSpec documents reflect the final agreed behavior
4. code is committed
5. the implementation branch has been merged back to `master`
6. the current branch is `master`

Archive must be performed on `master`. If the current branch is not `master`, stop before moving files and switch back to `master` first.

After archive files are moved, immediately commit the archive changes on `master` and push `master` to the remote. Do not present archive as workflow-complete until the archive commit has been pushed successfully.

If OpenSpec work is complete but branch merge has not happened yet, say that the change is document-complete but not workflow-complete.

## Deletion rules

If the request deletes a feature:

- force a dedicated change
- describe:
  - why it is being removed
  - what pages are affected
  - what APIs are affected
  - what data structures are affected
  - whether there is user data impact
  - whether rollback is possible

Do not allow deletion of a whole module in one direct step.

## Decision prompts to apply

When a user asks for a change, always reason through:

1. Is this one module, one feature, one bug, or one delete?
2. Is the scope too broad for one change?
3. Does this level of work require PRD only, or PRD + visual?
4. If implementation starts now, are we on the correct branch?
5. If archive is requested, has merge-back to `master` happened?

## Project-specific expectations

For this repository, be especially strict about splitting changes when the request mixes:

- 简历管理
- 投递记录
- 面试进展
- 岗位库
- 导航层级调整

These are often related, but they should not automatically become one change unless the user explicitly wants one tightly bounded feature slice.
