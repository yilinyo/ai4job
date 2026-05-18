## 1. 响应 DTO 设计与转换

- [x] 1.1 盘点现有 `/api/v1` handler 的响应点，确认哪些接口仍直接返回领域实体或 GORM 模型。
- [x] 1.2 在 `backend/internal/application/dto` 中新增简历仓库 response DTO 与转换函数，字段使用前端小驼峰格式。
- [x] 1.3 在 `backend/internal/application/dto` 中新增简历版本 response DTO 与转换函数，字段使用前端小驼峰格式。
- [x] 1.4 复核用户 response DTO，确保登录、获取个人资料和更新个人资料响应不暴露密码或内部字段。

## 2. Handler 响应改造

- [x] 2.1 修改简历仓库列表接口，返回 `data.list` 中的 DTO 列表而不是 `[]domainrepo.ResumeRepo`。
- [x] 2.2 修改简历仓库详情接口，返回单个 `ResumeRepoResponse` 而不是 `*domainrepo.ResumeRepo`。
- [x] 2.3 修改简历版本列表接口，返回 DTO 列表而不是 `[]domainversion.ResumeVersion`。
- [x] 2.4 修改简历版本详情接口，返回单个 `ResumeVersionResponse` 而不是 `*domainversion.ResumeVersion`。
- [x] 2.5 检查创建、更新、删除、设置默认版本等接口的空响应仍保持现有统一 envelope。

## 3. 后端测试

- [x] 3.1 为简历仓库 handler 补充响应 JSON 字段断言，验证包含 `id/userId/name/createdAt/updatedAt` 且不包含 `ID/UserID/DeletedAt`。
- [x] 3.2 为简历版本 handler 补充响应 JSON 字段断言，验证包含 `id/repoId/title/content/versionNum/isDefault` 且不包含 `ID/RepoID/DeletedAt`。
- [x] 3.3 为用户 handler 响应补充或复核断言，验证资料字段保持小驼峰且不暴露密码。
- [x] 3.4 运行后端相关单元测试，确认 DTO 改造未破坏 handler、router、service 现有行为。

## 4. 前端联调检查

- [x] 4.1 使用真实后端响应验证简历仓库列表可读取 `repo.id` 并跳转到 `/resume/repo/:id`。
- [x] 4.2 验证简历仓库详情接口收到的 id 不再是 `undefined`。
- [x] 4.3 验证简历版本列表和详情页面可以读取版本 `id/title/content/isDefault` 等字段。
- [x] 4.4 验证登录、获取个人资料和更新个人资料仍能按前端现有模型读取用户信息。

## 5. 验证

- [x] 5.1 验证 `proposal.md`、`design.md`、specs 与最终实现范围一致，没有把未实现的新业务字段伪造成响应数据。
- [x] 5.2 验证后端测试通过：`cd backend && go test ./...`。
- [x] 5.3 验证核心 API JSON 结构正确：简历仓库、简历版本、用户资料响应均为小驼峰 DTO 字段。
- [x] 5.4 验证主用户流程走通：登录后进入简历仓库列表，查看仓库详情，再查看版本列表/详情。
- [x] 5.5 验证关键约束：响应中不包含 `DeletedAt`、密码、GORM 内部字段或 Go 导出字段名。
- [x] 5.6 验证前端视觉与交互无回归：卡片、详情入口和错误提示保持原有交互，不因 DTO 改造出现空白或 `undefined` 文案。
