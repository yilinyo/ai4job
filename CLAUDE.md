# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project at a glance
- This repo is split into `frontend/`, `backend/`, `docs/`, and root-level OpenSpec files.
- The frontend is a **Vue 2 + TypeScript + Element UI** admin app under `frontend/`, adapted into a **resume management system**.
- It is built on top of `vue-typescript-admin-template`, so both template demo modules and resume-specific modules coexist.

## Common commands

### Install dependencies
```bash
cd frontend
npm install --legacy-peer-deps
```

### Local development
```bash
cd frontend
npm run serve
```
- Runs frontend dev server + mock API server concurrently.
- Frontend dev server port is set in `frontend/vue.config.js` (`9527`).
- Mock API server runs on `9528` (`frontend/mock/mock-server.ts`).

Run mock API only:
```bash
cd frontend
npm run mock
```

Run frontend only (without `concurrently` wrapper):
```bash
cd frontend
npx vue-cli-service serve
```

### Lint
```bash
cd frontend
npm run lint
```

### Unit tests
Run full unit test suite:
```bash
cd frontend
npm run test:unit
```

Run a single test file:
```bash
cd frontend
npx vue-cli-service test:unit --testPathPattern=tests/unit/utils/validate.spec.ts
```

### Build
Production build:
```bash
cd frontend
npm run build:prod
```

Staging build:
```bash
cd frontend
npm run build:stage
```

## Architecture (high-level)

### 1) Frontend data flow
- UI pages/components live under `frontend/src/views` and `frontend/src/components`.
- Business state is managed by class-based Vuex modules in `frontend/src/store/modules/*` (using `vuex-module-decorators`).
- Vuex modules call typed API wrappers in `frontend/src/api/*`.
- API wrappers use a shared axios instance in `frontend/src/utils/request.ts`.

Core convention: `request.ts` expects backend responses shaped like:
```ts
{ code: number, message: string, data: any }
```
and treats `code === 20000` as success.

### 2) Auth + permission routing model
- Token is stored via cookie helpers and injected as `X-Access-Token` header in `frontend/src/utils/request.ts`.
- Route guards in `frontend/src/permission.ts`:
  - fetch user info/roles (`UserModule.GetUserInfo`) if needed,
  - generate role-filtered routes (`PermissionModule.GenerateRoutes`),
  - dynamically register routes via `router.addRoute`.
- Router split:
  - `constantRoutes` (always available),
  - `asyncRoutes` (role-gated; includes resume module).

### 3) Resume domain modules
Resume features were added as a cohesive vertical slice:
- Route module: `frontend/src/router/modules/resume.ts`
- Models: `frontend/src/models/resume-repo.ts`, `frontend/src/models/resume-version.ts`
- API clients: `frontend/src/api/resume-repos.ts`, `frontend/src/api/resume-versions.ts`
- Vuex modules: `frontend/src/store/modules/resume-repo.ts`, `frontend/src/store/modules/resume-version.ts`
- Views:
  - repo list/detail: `frontend/src/views/resume-repo/*`
  - version editor/view: `frontend/src/views/resume-version/*`
  - user profile/register: `frontend/src/views/user/*`

When extending resume functionality, keep these layers aligned (model → api → store → view → route).

### 4) Mock backend and environment wiring
- Dev frontend calls `VUE_APP_BASE_API` (`/dev-api` in `.env.development`).
- `frontend/vue.config.js` proxies that prefix to `http://127.0.0.1:9528` and strips prefix.
- Mock server (`frontend/mock/mock-server.ts`) exposes resume/user REST endpoints under `/api/*`.
- Mock handlers live in:
  - `frontend/mock/users.ts`
  - `frontend/mock/resume-repos.ts`
  - `frontend/mock/resume-versions.ts`

If you add or change API behavior for local development, update both frontend API wrappers and corresponding mock handlers.

## Important project-specific notes from current docs/code
- Default local login credentials in mock data include `admin / 111111` (also `user / 111111`).
- Root README describes the workspace; PRDs and visual drafts live in `docs/PRD文档/`, while project notes, API docs, and legacy frontend docs live in `docs/项目文档/`.
