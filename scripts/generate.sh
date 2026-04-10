#!/usr/bin/env bash
set -euo pipefail

LAYER="${1:-all}"
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

run_backend() {
  echo "==> Generating backend code..."
  make -C "$REPO_ROOT/go-backend" generate
}

run_frontend() {
  echo "==> Generating frontend code..."
  pnpm --filter "./apps/react-frontend" run generate
}

case "$LAYER" in
  backend)
    run_backend
    ;;
  frontend)
    run_frontend
    ;;
  all)
    run_backend
    run_frontend
    ;;
  e2e|*)
    echo "==> No code generation needed for layer: $LAYER"
    ;;
esac
