#!/usr/bin/env bash
# cspell:ignore worktree worktrees WORKTREE
set -euo pipefail

GIT_COMMON_DIR=$(git rev-parse --git-common-dir) || {
  echo "Fatal: not a git repository" >&2
  exit 1
}

REGISTRY="$GIT_COMMON_DIR/worktree-ports"
PWD_ESCAPED=$(printf '%s\n' "$PWD" | sed 's/[[\.*^$()+?{|]/\\&/g')

if ! touch "$REGISTRY" 2>/dev/null; then
  echo "Error: Cannot write to registry: $REGISTRY" >&2
  exit 1
fi

if grep -q "^${PWD_ESCAPED}=" "$REGISTRY" 2>/dev/null; then
  OFFSET=$(grep "^${PWD_ESCAPED}=" "$REGISTRY" | head -1 | cut -d= -f2)
  echo "Already registered: WORKTREE_OFFSET=$OFFSET"
else
  USED=$(grep -oE '=[0-9]+' "$REGISTRY" 2>/dev/null | sed 's/=//g' || true)
  OFFSET=0
  MAX_OFFSET=900
  while echo "$USED" | grep -q "^${OFFSET}$" 2>/dev/null; do
    OFFSET=$((OFFSET + 100))
    if [ "$OFFSET" -gt "$MAX_OFFSET" ]; then
      echo "Error: Too many worktrees registered (max 10)" >&2
      exit 1
    fi
  done
  printf '%s=%d\n' "$PWD" "$OFFSET" >> "$REGISTRY"
  echo "New offset assigned: WORKTREE_OFFSET=$OFFSET"
fi

ENV_FILE=".env.local"
ENV_TMP="${ENV_FILE}.tmp.$$"

if [ -f "$ENV_FILE" ]; then
  grep -v '^WORKTREE_OFFSET=\|^APP_PORT=\|^APP_GRPC_PORT=\|^VITE_PORT=\|^DB_PORT=\|^CORS_ALLOW_ORIGINS=' "$ENV_FILE" > "$ENV_TMP" || true
else
  > "$ENV_TMP"
fi

{
  echo "WORKTREE_OFFSET=$OFFSET"
  echo "APP_PORT=$((8080 + OFFSET))"
  echo "APP_GRPC_PORT=$((8081 + OFFSET))"
  echo "VITE_PORT=$((3000 + OFFSET))"
  echo "DB_PORT=$((55432 + OFFSET))"
  echo "CORS_ALLOW_ORIGINS=http://localhost:$((3000 + OFFSET))"
} >> "$ENV_TMP"

mv "$ENV_TMP" "$ENV_FILE"
echo "Done: APP_PORT=$((8080+OFFSET)) APP_GRPC_PORT=$((8081+OFFSET)) VITE_PORT=$((3000+OFFSET)) DB_PORT=$((55432+OFFSET))"
