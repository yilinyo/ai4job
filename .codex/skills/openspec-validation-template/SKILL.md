---
name: openspec-validation-template
description: Use when creating, updating, or implementing OpenSpec changes in this repository and you need a reusable validation standard. Applies a consistent verification structure to every tasks.md and requires implementation, runtime, and acceptance validation to be explicit.
---

# OpenSpec Validation Template

This is a project-level skill for `/Users/yilin/project/ai-for`.

Use this skill whenever you:

- create a new OpenSpec change
- revise an existing `tasks.md`
- implement and validate an OpenSpec change
- prepare a change for archive

## Core Rule

Every OpenSpec change in this repository must include an explicit validation section in `tasks.md`, and validation must follow the same layered logic.

Do not leave validation as vague statements such as:

- “self test”
- “manual verify”
- “check if works”

Validation must be concrete, traceable, and tied to the actual behavior promised in the change.

## Required Validation Layers

For each OpenSpec change, validate in this order when applicable:

1. **Document alignment**
2. **Static validation**
3. **Data / API validation**
4. **User flow validation**
5. **Constraint / exception validation**
6. **Visual / interaction consistency validation**

Not every change needs every layer equally, but every change must be consciously evaluated against them.

## Layer Definitions

### 1. Document alignment

Check that:

- `proposal.md` still matches the true scope
- `design.md` still matches the final interaction and data model
- `tasks.md` matches the actual implementation and validation work

If the user changed requirements mid-stream, update artifacts before validation is considered complete.

### 2. Static validation

At minimum, run applicable static checks such as:

- lint
- type checking
- unit tests if present

For this project, default to:

```bash
cd frontend
npm run lint
```

If unit tests are directly relevant and available, also consider:

```bash
cd frontend
npm run test:unit
```

### 3. Data / API validation

Verify that the underlying data layer behaves as expected:

- Mock API routes are reachable
- returned JSON shape matches the documented model
- create/update/delete flows persist correctly
- linked fields and derived fields are consistent

Examples:

- status synchronization
- relation fields such as `jobPostingId`
- snapshot fields that should remain stable

### 4. User flow validation

Validate the real user path end-to-end, not just isolated fields.

Typical structure:

1. Enter source page
2. Trigger create/edit/view action
3. Complete required input
4. Save
5. Return to list/detail
6. Confirm state is reflected correctly

This should map directly to the primary user flow promised in the change.

### 5. Constraint / exception validation

Explicitly verify non-happy-path behavior, such as:

- deletion guards
- failure states
- required field validation
- status transitions
- blocked operations

If a change introduces a business rule, validation for that rule must appear in `tasks.md`.

### 6. Visual / interaction consistency validation

Validate the UI against agreed constraints, especially when the user has given explicit visual guidance.

In this repository, common examples include:

- same hierarchy level in navigation
- card size consistency
- spacing and container padding
- page vs dialog interaction choice
- merged vs split information regions
- color/state mapping for process nodes

If a visual requirement changed during discussion, it must appear in both `design.md` and validation tasks.

## Required tasks.md pattern

Every new change should include a final validation section using a structure close to this:

```md
## N. 验证

- [ ] N.1 验证静态检查通过（如 lint / type check / unit test）
- [ ] N.2 验证核心数据接口返回结构正确
- [ ] N.3 验证主用户流程可完整走通
- [ ] N.4 验证关键约束、异常或删除保护生效
- [ ] N.5 验证视觉与交互符合已确认设计
```

Adapt the wording to the change, but keep the layered logic.

## Project-specific validation defaults

When writing validation tasks in this repository, prefer covering:

- 页面入口是否正确
- 列表 / 详情 / 创建 / 编辑主流程是否可走通
- Mock API 是否返回预期结构
- 关系字段是否正确写入
- 关键状态是否同步
- 删除或限制规则是否生效
- UI 是否和既有模块风格保持一致

## During implementation

When applying a change:

1. Do not mark validation tasks complete just because code was written.
2. Complete validation tasks only after real checks were performed.
3. If environment constraints prevent validation, say so explicitly.
4. Do not archive a change with incomplete validation unless the scope is intentionally reduced and artifacts are updated.

## Archive readiness checklist

Before recommending archive, confirm:

- all implementation tasks are complete
- all validation tasks are complete
- docs and code are aligned
- no critical unresolved question remains

If any of the above is false, the change is not archive-ready.

## Quick template for future tasks.md files

When creating or normalizing a change, use this scaffold:

```md
## X. 验证

- [ ] X.1 验证文档、设计与实现保持一致
- [ ] X.2 验证静态检查通过
- [ ] X.3 验证核心接口或数据链路正确
- [ ] X.4 验证主用户流程走通
- [ ] X.5 验证关键约束或异常处理生效
- [ ] X.6 验证视觉与交互符合确认口径
```

This scaffold is the default unless the change has a strong reason to deviate.
