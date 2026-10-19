#!/usr/bin/env bash
set -euo pipefail

compose_file="${COMPOSE_FILE:-deploy/docker-compose.prod.yml}"
env_file="${AI_FOR_ENV_FILE:?AI_FOR_ENV_FILE is required}"
timeout="${HEALTHCHECK_TIMEOUT_SECONDS:-120}"
start_time="$(date +%s)"

compose=(docker compose --env-file "$env_file" -f "$compose_file")

until "${compose[@]}" exec -T api wget -q -O /dev/null http://127.0.0.1:8080/healthz && curl -fsS http://127.0.0.1:"${WEB_PORT:-80}"/healthz >/dev/null; do
  if (( $(date +%s) - start_time >= timeout )); then
    "${compose[@]}" ps
    "${compose[@]}" logs --tail=100 api web || true
    exit 1
  fi
  sleep 3
done
