#!/usr/bin/env bash
# Worktree + Herdr pane + claude + role prompt for one jaira lane.
# Usage: spawn.sh <slug> <ticket-id> <lane> [repo-root]
# Prints the pane id on stdout; everything else goes to stderr.
set -euo pipefail

slug="${1:?slug}"; ticket="${2:?ticket id}"; lane="${3:?lane}"
root="${4:-$PWD}"
herdr="${HERDR_BIN_PATH:-herdr}"
[ "${HERDR_ENV:-}" = 1 ] || { echo "not inside a Herdr pane" >&2; exit 1; }

# Beside the repository, never inside it: a worktree under the repo is a second
# copy of the sources on a different branch, and grep -r / find / ls -R walk
# straight into it. Derived, so no path is baked in for one machine.
wt="$(cd "$root/.." && pwd)/.worktrees/$(basename "$root")-$slug"

if [ ! -d "$wt" ]; then
  git -C "$root" worktree add "$wt" -b "feature/$slug" >&2

  # Only a repo that carries a container stack needs its own ports. A repo
  # without a .env has nothing to offset, and must not fail here.
  if [ -f "$root/.env" ]; then
    # Deterministic per-slug port offset, so two workers never share a stack.
    # 0 stays free: it belongs to the main directory (80, 5432, 5433, 5173, 8000).
    off=$(( ( $(printf '%s' "$slug" | cksum | cut -d' ' -f1) % 40 ) + 1 ))
    cp "$root/.env" "$wt/.env"
    {
      echo
      echo "# worker stack: slug=$slug offset=$off"
      echo "COMPOSE_PROJECT_NAME=rg_$slug"
      echo "HTTP_PORT=$((8080 + off))"
      echo "DB_PORT_HOST=$((5500 + off * 2))"
      echo "DB_PORT_TEST_HOST=$((5501 + off * 2))"
      echo "VITE_PORT_HOST=$((5200 + off))"
      echo "BACKEND_PORT_HOST=$((8100 + off))"
    } >> "$wt/.env"
  fi
fi

pane="$("$herdr" pane split --current --direction down --no-focus \
  | python3 -c 'import sys,json;print(json.load(sys.stdin)["result"]["pane"]["pane_id"])')"

"$herdr" pane run "$pane" "cd '$wt' && claude" >/dev/null

# The state hook reports idle a few seconds after startup.
for _ in $(seq 20); do
  sleep 2
  st="$("$herdr" pane get "$pane" \
    | python3 -c 'import sys,json;p=json.load(sys.stdin)["result"]["pane"];print(p.get("agent","-"),p.get("agent_status","-"))')"
  case "$st" in claude\ idle|claude\ done) break ;; esac
done
case "${st:-}" in claude*) ;; *) echo "claude did not come up in $pane: ${st:-none}" >&2; exit 1 ;; esac

"$herdr" pane send-text "$pane" "/jaira-role-lane $ticket $lane" >/dev/null
sleep 1
"$herdr" pane send-keys "$pane" enter >/dev/null

echo "$pane"
