#!/usr/bin/env bash
set -euo pipefail

frontend_changed="${FRONTEND_CHANGED:-false}"
backend_changed="${BACKEND_CHANGED:-false}"
infrastructure_changed="${INFRASTRUCTURE_CHANGED:-false}"
env_file="${AI_FOR_ENV_FILE:-/etc/ai-for/production.env}"
compose_file="${COMPOSE_FILE:-deploy/docker-compose.prod.yml}"
log_dir="${DEPLOY_LOG_DIR:-/var/lib/ai-for}"

[[ -f "$env_file" ]] || { echo "Missing production environment file: $env_file" >&2; exit 1; }
set -a
# shellcheck disable=SC1090
. "$env_file"
set +a

base_version="$(tr -d '[:space:]' < VERSION)"
short_sha="$(git rev-parse --short HEAD)"
export GIT_SHA="$(git rev-parse HEAD)"
export APP_VERSION="v${base_version}-run.${GITHUB_RUN_NUMBER:-manual}-${short_sha}"

compose=(docker compose --env-file "$env_file" -f "$compose_file")
mkdir -p "$log_dir"
record() {
  printf '{"version":"%s","git_sha":"%s","services":"%s","status":"%s","time":"%s"}\n' \
    "$APP_VERSION" "$GIT_SHA" "$services" "$1" "$(date -u +%FT%TZ)" >> "$log_dir/deployments.jsonl"
}

services=""
if [[ "$infrastructure_changed" == true ]]; then
  frontend_changed=true
  backend_changed=true
fi
if [[ "$frontend_changed" != true && "$backend_changed" != true ]]; then
  echo "No deployable files changed."
  exit 0
fi

trap 'record failed' ERR
"${compose[@]}" up -d postgres redis

if [[ "$backend_changed" == true ]]; then
  services="api"
  "${compose[@]}" run --rm migrate
  "${compose[@]}" up -d --build --no-deps api
fi
if [[ "$frontend_changed" == true ]]; then
  services="${services:+$services,}web"
  "${compose[@]}" up -d --build --no-deps web
fi

AI_FOR_ENV_FILE="$env_file" COMPOSE_FILE="$compose_file" WEB_PORT="${WEB_PORT:-80}" deploy/scripts/healthcheck.sh
record success
docker image prune -f
echo "Deployed $APP_VERSION ($services)"
