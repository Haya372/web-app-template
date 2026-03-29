#!/usr/bin/env bash
# cspell:ignore worktree worktrees
set -euo pipefail

GIT_COMMON_DIR=$(git rev-parse --git-common-dir)
REGISTRY="$GIT_COMMON_DIR/worktree-ports"

if [ ! -f "$REGISTRY" ] || [ ! -s "$REGISTRY" ]; then
  echo "No worktrees registered."
  exit 0
fi

echo "Registered worktrees:"
while IFS='=' read -r path offset; do
  echo "  [offset=$offset] APP_PORT=$((8080+offset)) VITE_PORT=$((3000+offset)) DB_PORT=$((55432+offset))  $path"
done < "$REGISTRY"
