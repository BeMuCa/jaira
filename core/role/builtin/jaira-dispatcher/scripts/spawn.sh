#!/usr/bin/env bash
# Worktree + Herdr pane + claude + role prompt for one jaira lane.
# Usage: spawn.sh [--no-worktree] <slug> <ticket-id> <lane> [repo-root]
# Prints the pane id on stdout; everything else goes to stderr.
set -euo pipefail

usage() {
  cat <<'USAGE'
spawn.sh [--no-worktree] <slug> <ticket-id> <lane> [repo-root]

  --no-worktree   Start the worker in the repository directory itself: no
                  worktree, no branch of its own. For a one-lane job that
                  belongs on the branch already checked out. Two workers then
                  share one directory, so run only one at a time. The slug is
                  still required but unused: it names a worktree, and there is
                  none.
                  JAIRA_NO_WORKTREE=1 in the environment does the same.
  -h, --help      This text.

The lane <dispatch> is special: it starts the dispatcher itself in the tab
(/jaira-dispatcher <ticket-id>) instead of one lane's worker. Every other lane
name starts /jaira-role-lane <ticket-id> <lane>.

Environment: JAIRA_BRANCH_PREFIX (default feat), JAIRA_NO_WORKTREE,
HERDR_BIN_PATH, HERDR_WORKSPACE_ID.
USAGE
}

# Flags before the positionals, and -- ends them: a slug is free text and could
# one day start with a dash.
no_worktree="${JAIRA_NO_WORKTREE:-}"
while [ $# -gt 0 ]; do
  case "$1" in
    --no-worktree) no_worktree=1; shift ;;
    -h|--help)     usage; exit 0 ;;
    --)            shift; break ;;
    -*)            echo "unknown flag: $1" >&2; usage >&2; exit 2 ;;
    *)             break ;;
  esac
done

slug="${1:?slug}"; ticket="${2:?ticket id}"; lane="${3:?lane}"
root="${4:-$PWD}"
herdr="${HERDR_BIN_PATH:-herdr}"
[ "${HERDR_ENV:-}" = 1 ] || { echo "not inside a Herdr pane" >&2; exit 1; }

# Beside the repository, never inside it: a worktree under the repo is a second
# copy of the sources on a different branch, and grep -r / find / ls -R walk
# straight into it. Derived, so no path is baked in for one machine.
# JAIRA_NO_WORKTREE=1 runs the worker in $root itself, no worktree and no
# branch of its own. For a one-lane job on the branch that is already checked
# out — a doc line, a note, a lane that only reads — a worktree costs a clone,
# a branch and a merge for nothing. It is unsafe for two workers at once: they
# share the directory, and that is exactly what the worktree exists to prevent.
if [ "$no_worktree" = 1 ]; then
  wt="$root"
else
  wt="$(cd "$root/.." && pwd)/.worktrees/$(basename "$root")-$slug"
fi

if [ ! -d "$wt" ]; then
  # feat/, the prefix this board's own branches use. A repository that names
  # them differently sets JAIRA_BRANCH_PREFIX rather than editing this script.
  git -C "$root" worktree add "$wt" -b "${JAIRA_BRANCH_PREFIX:-feat}/$slug" >&2

  # Only a repo that carries a container stack needs its own ports. A repo
  # without a .env has nothing to offset, and must not fail here.
  if [ -f "$root/.env" ]; then
    # Deterministic per-slug port offset, so two workers never share a stack.
    # 0 stays free: it belongs to the main directory's own stack.
    off=$(( ( $(printf '%s' "$slug" | cksum | cut -d' ' -f1) % 40 ) + 1 ))
    cp "$root/.env" "$wt/.env"
    {
      echo
      echo "# worker stack: slug=$slug offset=$off"
      # Derived from the repository, never a name written in here: a prefix
      # baked into this script belongs to one project and silently names every
      # other project's stack after it. Docker takes lower case only, so the
      # name is folded before it is cleaned — a jaira slug is upper case, and
      # `docker compose` refuses the whole stack over a single capital. printf
      # rather than echo: the newline echo appends is a character like any
      # other to tr, and would come back as a trailing underscore.
      echo "COMPOSE_PROJECT_NAME=$(printf '%s_%s' "$(basename "$root")" "$slug" \
        | tr '[:upper:]' '[:lower:]' | tr -c 'a-z0-9_-' '_')"
      echo "HTTP_PORT=$((8080 + off))"
      echo "DB_PORT_HOST=$((5500 + off * 2))"
      echo "DB_PORT_TEST_HOST=$((5501 + off * 2))"
    } >> "$wt/.env"
  fi
fi

# Named, because without --workspace Herdr decides for itself where the tab
# lands: the worker can open in a window nobody is looking at, and a worker
# nobody sees is one nobody notices dying. Herdr exports its own workspace into
# every pane it starts, so this is the workspace the caller is sitting in.
ws=()
if [ -n "${HERDR_WORKSPACE_ID:-}" ]; then ws=(--workspace "$HERDR_WORKSPACE_ID"); fi

# A tab per worker, never a split. A split divides the height of one screen: at
# four workers each strip is a few lines, and nobody can read what any of them is
# doing — which is the whole reason a worker gets a surface of its own.
pane="$("$herdr" tab create ${ws[@]+"${ws[@]}"} --cwd "$wt" --label "$ticket/$lane" --no-focus \
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
# Only the two states the loop breaks on may pass. Anything else — above all
# Herdr's `blocked`, its state for a detected approval dialog — must stop here:
# the send-keys below would answer that dialog on the human's behalf, and this
# role is forbidden from ever doing that.
case "${st:-}" in
  claude\ idle|claude\ done) ;;
  claude\ blocked)
    echo "claude is up in $pane but an approval dialog is waiting:" \
         "report it to the human, let them answer it in that pane," \
         "then start this worker again" >&2
    exit 1 ;;
  *) echo "claude did not come up in $pane: ${st:-none}" >&2; exit 1 ;;
esac

# The lane name "dispatch" starts the dispatcher itself rather than one lane's
# worker — same worktree, same tab, same waiting for claude to come up, so it
# lives here rather than in a second script that drifts from this one.
if [ "$lane" = dispatch ]; then
  "$herdr" pane send-text "$pane" "/jaira-dispatcher $ticket" >/dev/null
else
  "$herdr" pane send-text "$pane" "/jaira-role-lane $ticket $lane" >/dev/null
fi
sleep 1
"$herdr" pane send-keys "$pane" enter >/dev/null

echo "$pane"
