#!/usr/bin/env bash
# cspell:ignore worktree worktrees
set -euo pipefail

GIT_COMMON_DIR=$(git rev-parse --git-common-dir)
REGISTRY="$GIT_COMMON_DIR/worktree-ports"
PWD_ESCAPED=$(printf '%s\n' "$PWD" | sed 's/[[\.*^$()+?{|]/\\&/g')

if [ -f "$REGISTRY" ]; then
  grep -v "^${PWD_ESCAPED}=" "$REGISTRY" > "$REGISTRY.tmp.$$" && mv "$REGISTRY.tmp.$$" "$REGISTRY"
  echo "Cleaned: $PWD"
fi
