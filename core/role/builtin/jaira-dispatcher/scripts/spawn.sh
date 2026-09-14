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
  # feat/, the prefix this board's own branches use. A repository that names
  # them differently sets JAIRA_BRANCH_PREFIX rather than editing this script.
  git -C "$root" worktree add "$wt" -b "${JAIRA_BRANCH_PREFIX:-feat}/$slug" >&2

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
      # Derived from the repository, never a name written in here: a prefix
      # baked into this script belongs to one project and silently names every
      # other project's stack after it.
      echo "COMPOSE_PROJECT_NAME=$(basename "$root" | tr -c 'a-zA-Z0-9' '_')_$slug"
      echo "HTTP_PORT=$((8080 + off))"
      echo "DB_PORT_HOST=$((5500 + off * 2))"
      echo "DB_PORT_TEST_HOST=$((5501 + off * 2))"
    } >> "$wt/.env"
  fi
fi

# A tab per worker, never a split. A split divides the height of one screen: at
# four workers each strip is a few lines, and nobody can read what any of them is
# doing — which is the whole reason a worker gets a surface of its own.
pane="$("$herdr" tab create --cwd "$wt" --label "$ticket/$lane" --no-focus \
  | python3 -c 'import sys,json;print(json.load(sys.stdin)["result"]["root_pane"]["pane_id"])')"

# The tab's shell runs on the machine Herdr itself runs on. When that is Windows
# and this script runs in WSL, that shell has no /home/... at all: --cwd is
# resolved against Windows and dropped, `cd "$wt"` fails, and claude comes up in
# the Windows home in front of its trust dialog — which the send-keys below would
# then answer on the human's behalf. Cross back into WSL instead: wsl.exe --cd
# sets the directory before any shell starts, so no cd function and no wrong home
# can intervene, and `bash -lic` is what puts claude on PATH.
case "$herdr" in
  /mnt/*|*.exe) start="wsl.exe --cd '$wt' -- bash -lic claude" ;;
  *)            start="cd '$wt' && claude" ;;
esac

"$herdr" pane run "$pane" "$start" >/dev/null

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
