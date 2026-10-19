## Why

登录页和注册页虽已在平板宽度下改为单列，但手机窄屏仍会显示占用大量纵向空间的品牌介绍、三列信息块和大尺寸卡片。用户需要滚动较长距离才能进入或完成认证表单，验证码操作在较窄设备上也不够紧凑。

## What Changes

- 为登录页和注册页补充面向手机窄屏的布局规则，收敛品牌信息、缩小页面与卡片留白，并确保表单、按钮和入口在不横向滚动的前提下可见、可操作。
- 调整注册页的验证码输入和发送控件，使其在常见手机宽度上保持稳定的单行布局与足够的触控尺寸。
- 保持现有登录、注册、邮箱验证码、表单校验和路由行为不变。

## Capabilities

### New Capabilities

- `auth-page-mobile-layout`: 登录与注册认证页面在手机窄屏下的可读、可操作布局要求。

### Modified Capabilities

- `auth-page-enterprise-layout`: 认证页蓝白企业风格需在手机窄屏布局中保持连续，同时避免装饰性内容压缩核心认证流程。

## Impact

- 前端：修改 `frontend/src/views/login/index.vue` 与 `frontend/src/views/user/register.vue` 的响应式样式。
- 文档：新增认证页手机端布局的 OpenSpec 要求和 PRD。
- API、路由、认证数据模型及后端逻辑均不变。
