# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project at a glance
- This repo is a **Vue 2 + TypeScript + Element UI** admin app, adapted into a **resume management system** (resume repos + resume versions + user profile).
- It is built on top of `vue-typescript-admin-template`, so both template demo modules and resume-specific modules coexist.

## Common commands

### Install dependencies
```bash
npm install --legacy-peer-deps
```

### Local development
```bash
npm run serve
```
- Runs frontend dev server + mock API server concurrently.
- Frontend dev server port is set in `vue.config.js` (`9527`).
- Mock API server runs on `9528` (`mock/mock-server.ts`).

Run mock API only:
```bash
npm run mock
```

Run frontend only (without `concurrently` wrapper):
```bash
npx vue-cli-service serve
```

### Lint
```bash
npm run lint
```

### Unit tests
Run full unit test suite:
```bash
npm run test:unit
```

Run a single test file:
```bash
npx vue-cli-service test:unit --testPathPattern=tests/unit/utils/validate.spec.ts
```

### Build
Production build:
```bash
npm run build:prod
```

Staging build:
```bash
npm run build:stage
```

## Architecture (high-level)

### 1) Frontend data flow
- UI pages/components live under `src/views` and `src/components`.
- Business state is managed by class-based Vuex modules in `src/store/modules/*` (using `vuex-module-decorators`).
- Vuex modules call typed API wrappers in `src/api/*`.
- API wrappers use a shared axios instance in `src/utils/request.ts`.

Core convention: `request.ts` expects backend responses shaped like:
```ts
{ code: number, message: string, data: any }
```
and treats `code === 20000` as success.

### 2) Auth + permission routing model
- Token is stored via cookie helpers and injected as `X-Access-Token` header in `src/utils/request.ts`.
- Route guards in `src/permission.ts`:
  - fetch user info/roles (`UserModule.GetUserInfo`) if needed,
  - generate role-filtered routes (`PermissionModule.GenerateRoutes`),
  - dynamically register routes via `router.addRoute`.
- Router split:
  - `constantRoutes` (always available),
  - `asyncRoutes` (role-gated; includes resume module).

### 3) Resume domain modules
Resume features were added as a cohesive vertical slice:
- Route module: `src/router/modules/resume.ts`
- Models: `src/models/resume-repo.ts`, `src/models/resume-version.ts`
- API clients: `src/api/resume-repos.ts`, `src/api/resume-versions.ts`
- Vuex modules: `src/store/modules/resume-repo.ts`, `src/store/modules/resume-version.ts`
- Views:
  - repo list/detail: `src/views/resume-repo/*`
  - version editor/view: `src/views/resume-version/*`
  - user profile/register: `src/views/user/*`

When extending resume functionality, keep these layers aligned (model → api → store → view → route).

### 4) Mock backend and environment wiring
- Dev frontend calls `VUE_APP_BASE_API` (`/dev-api` in `.env.development`).
- `vue.config.js` proxies that prefix to `http://127.0.0.1:9528` and strips prefix.
- Mock server (`mock/mock-server.ts`) exposes resume/user REST endpoints under `/api/*`.
- Mock handlers live in:
  - `mock/users.ts`
  - `mock/resume-repos.ts`
  - `mock/resume-versions.ts`

If you add or change API behavior for local development, update both frontend API wrappers and corresponding mock handlers.

## Important project-specific notes from current docs/code
- Default local login credentials in mock data include `admin / 111111` (also `user / 111111`).
- README still contains template-era references (ports/URLs) in places; rely on `package.json`, `vue.config.js`, and `mock/mock-server.ts` as source of truth for local run behavior.
