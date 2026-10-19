#!/usr/bin/env bash
set -euo pipefail

base_ref="${1:-${GITHUB_EVENT_BEFORE:-}}"
head_ref="${2:-${GITHUB_SHA:-HEAD}}"

frontend=false
backend=false
infrastructure=false

if [[ -z "$base_ref" || "$base_ref" =~ ^0+$ ]]; then
  frontend=true
  backend=true
  infrastructure=true
else
  changed_files="$(git diff --name-only "$base_ref" "$head_ref")"
  grep -qE '^frontend/' <<<"$changed_files" && frontend=true || true
  grep -qE '^backend/' <<<"$changed_files" && backend=true || true
  grep -qE '^(deploy/|scripts/ci/|\.github/workflows/|VERSION$|frontend/Dockerfile$|backend/Dockerfile$|\.dockerignore$)' <<<"$changed_files" && infrastructure=true || true
fi

if [[ "$infrastructure" == true ]]; then
  frontend=true
  backend=true
fi

emit() {
  local key="$1" value="$2"
  if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
    echo "$key=$value" >> "$GITHUB_OUTPUT"
  else
    echo "$key=$value"
  fi
}

emit frontend "$frontend"
emit backend "$backend"
emit infrastructure "$infrastructure"
