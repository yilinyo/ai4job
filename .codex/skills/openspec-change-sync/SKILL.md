---
name: openspec-change-sync
description: Use when working on OpenSpec changes in this repository and the user adds, clarifies, narrows, or revises requirements during the process. Ensures every meaningful mid-stream requirement change is immediately synced back into the relevant OpenSpec artifacts before or alongside implementation.
---

# OpenSpec Change Sync

This is a project-level skill for `/Users/yilin/project/ai-for`.

Use this skill whenever you are:

- proposing an OpenSpec change
- implementing an OpenSpec change
- reviewing or refining an existing OpenSpec change
- discussing product details that materially affect an existing OpenSpec change

Especially trigger this skill when the user adds new requirements or corrects previous assumptions in the middle of the workflow.

## Core Rule

If the user makes a meaningful requirement supplement during an OpenSpec workflow, you must update the corresponding OpenSpec documents promptly.

Do not leave the change only in:

- chat responses
- temporary reasoning
- code implementation
- a PRD outside `openspec/changes/`

The OpenSpec artifacts must remain aligned with the latest agreed product behavior.

## What counts as a “meaningful requirement supplement”

Treat the following as requiring document sync:

- adding a new user flow or entry point
- changing navigation level or page hierarchy
- changing a page layout structure
- changing interaction style
  Example: page -> dialog, dialog -> page
- changing status definitions, enums, or lifecycle nodes
- changing data ownership or entity relationships
- changing what fields are required, optional, or auto-filled
- changing implementation boundaries
  Example: frontend mock first, backend later
- changing acceptance expectations or validation criteria
- changing visual consistency requirements that affect scope
  Example: must match resume list card size and spacing

Small wording polish that does not affect product behavior does not require forced sync.

## Required sync behavior

When a meaningful supplement appears:

1. Identify which active OpenSpec change it belongs to.
2. Update the relevant artifacts before closing the turn, or in the same turn as implementation.
3. Prefer updating the minimum correct set of files, not every file by default.

Typical mapping:

- `proposal.md`
  Use when the change affects scope, goals, user value, module boundaries, or impact.
- `design.md`
  Use when the change affects data model, interaction model, route hierarchy, page composition, visual rules, state logic, or implementation approach.
- `tasks.md`
  Use when the change introduces new implementation work, validation work, or shifts task granularity.

## Sync heuristics

Use these heuristics when deciding what to edit:

### Update `proposal.md` if:

- the “what changes” section is no longer accurate
- the module boundary changed
- the user added a new first-class capability
- the navigation or product positioning changed materially

### Update `design.md` if:

- the user clarified page structure
- the interaction model changed
- the route level changed
- a status system or process flow changed
- the UI must follow a specific visual pattern
- a new field relationship or autofill behavior was agreed

### Update `tasks.md` if:

- there is new implementation work
- there is new validation work
- there is a newly agreed visual or behavior constraint
- completed work no longer matches the documented task list

## Project-specific expectations

For this repository, pay special attention to keeping OpenSpec docs updated when the user changes:

- `简历管理` vs `岗位列表` navigation hierarchy
- card layout consistency with existing resume list UI
- job posting detail page structure
- whether an action is a page, dialog, or inline section
- job posting statuses
- interview process node display rules
- whether frontend uses mock parsing first and backend parses later
- whether a relation is version-scoped, job-scoped, or application-scoped

These details have already changed during discussion in this project, so do not assume the earlier artifact text is still correct.

## During implementation

If you implement behavior that was agreed in chat but not yet reflected in OpenSpec:

1. Pause long enough to update the artifacts.
2. Then continue implementation.

Do not treat document sync as optional cleanup for later.

## Completion checklist

Before ending an OpenSpec-related turn, quickly check:

- Does the latest agreed behavior exist in code?
- Does it also exist in the correct OpenSpec artifact(s)?
- If tasks changed, is `tasks.md` updated?
- If a visual or structural decision changed, is `design.md` updated?
- If scope changed, is `proposal.md` updated?

If any answer is “no”, fix that before finishing the turn.
