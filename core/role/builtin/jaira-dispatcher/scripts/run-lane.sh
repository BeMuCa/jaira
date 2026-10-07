#!/usr/bin/env bash
# One lane from start to finish: spawn.sh starts a worker on it, this waits
# until the ticket has left the lane and the worker has finished its turn, then
# closes the worker's tab. Start it with run_in_background: its exit is what
# wakes the dispatcher.
# Usage: run-lane.sh [--no-worktree] [--keep] [--timeout <minutes>] <ticket-id> <lane> [repo-root]
set -euo pipefail

usage() {
  cat <<'USAGE'
run-lane.sh [--no-worktree] [--keep] [--timeout <minutes>] <ticket-id> <lane> [repo-root]

  --no-worktree   Passed on to spawn.sh; see its --help for when to take it.
                  JAIRA_NO_WORKTREE=1 in the environment does the same.
  --keep          Leave the worker's tab open when the lane is finished.
  --timeout <m>   Give up waiting after m minutes (default 240, or
                  JAIRA_LANE_TIMEOUT); 0 waits for ever. The worker is left
                  running and its tab open.
  -h, --help      This text.

Exit codes: 0 the lane is finished and the tab closed (or kept), 3 timed out,
4 the worker is at an approval dialog — its tab stays open for the human,
5 the worker stopped on a question — it is in the ticket's question field and
its notes; the tab is closed (or kept), and the lane runs again once answered.
USAGE
}

no_worktree="${JAIRA_NO_WORKTREE:-}"
keep=""
timeout="${JAIRA_LANE_TIMEOUT:-240}"
while [ $# -gt 0 ]; do
  case "$1" in
    --no-worktree) no_worktree=1; shift ;;
    --keep)        keep=1; shift ;;
    --timeout)     timeout="${2:?--timeout needs minutes}"; shift 2 ;;
    -h|--help)     usage; exit 0 ;;
    --)            shift; break ;;
    -*)            echo "unknown flag: $1" >&2; usage >&2; exit 2 ;;
    *)             break ;;
  esac
done

ticket="${1:?ticket id}"; lane="${2:?lane}"
root="$(cd "${3:-$PWD}" && pwd)"
herdr="${HERDR_BIN_PATH:-herdr}"
# Beside this file, not under ~/.claude/skills: roles are installed into a
# project's .claude/skills as well, and a hard-coded home path misses those.
here="$(cd "$(dirname "$0")" && pwd)"
slug="$(printf '%s' "$ticket" | tr '[:upper:]' '[:lower:]')"

# The worker writes the ticket in the directory it runs in, so that is where its
# status has to be read. Same derivation as spawn.sh. A board that is not
# shared yet is gitignored and missing from a fresh worktree; then the ticket
# only exists in the repository itself.
if [ "$no_worktree" = 1 ]; then
  wt="$root"
else
  wt="$(cd "$root/.." && pwd)/.worktrees/$(basename "$root")-$slug"
fi

flags=()
[ "$no_worktree" = 1 ] && flags+=(--no-worktree)
# spawn.sh types the lane command itself once claude is idle; sending a prompt
# of our own after it would hand the worker the lane a second time.
pane="$("$here/spawn.sh" ${flags[@]+"${flags[@]}"} "$slug" "$ticket" "$lane" "$root")"
echo "pane $pane"

board="$wt"; [ -d "$wt/.jaira" ] || board="$root"
status() {
  (cd "$board" && jaira show "$ticket" --json 2>/dev/null) \
    | python3 -c 'import sys,json;print(json.load(sys.stdin).get("status",""))' 2>/dev/null || true
}
# A worker never asks the person: it writes its question onto the ticket and
# stops in the lane. Without this the wait below sees a lane never left and
# sits there until the timeout.
question() {
  (cd "$board" && jaira show "$ticket" --json 2>/dev/null) \
    | python3 -c 'import sys,json;print(json.load(sys.stdin).get("question","").strip())' 2>/dev/null || true
}
# A finished worker reports done, not only idle. A loop that waits for
# idle|blocked alone never returns, and the tab never closes.
agent() {
  "$herdr" pane get "$pane" 2>/dev/null \
    | python3 -c 'import sys,json;print(json.load(sys.stdin)["result"]["pane"].get("agent_status","-"))' 2>/dev/null || true
}

deadline=$(( timeout > 0 ? SECONDS + timeout * 60 : 0 ))
check() {
  if [ "$deadline" -gt 0 ] && [ "$SECONDS" -ge "$deadline" ]; then
    echo "timed out after $timeout min: $ticket is in $(status), worker in $pane is $(agent), its tab left open" >&2
    exit 3
  fi
  # blocked is Herdr's state for an approval dialog. Answering it is the human's
  # call, never this script's, so stop here and leave the tab to them.
  if [ "$(agent)" = blocked ]; then
    echo "worker in $pane is at an approval dialog: report it to the human, the tab stays open" >&2
    exit 4
  fi
}

# Where the ticket stood when the worker started; the wait below needs it. An
# empty one is a failed read like any other, and taken as the start it would
# make the first status before the lane look like a worker that moved past it.
until start="$(status)"; [ -n "$start" ]; do check; sleep 5; done
# jaira move never clears the question field: a ticket back from human still
# carries the question answered there. Only a question other than this one is
# the worker's own.
q0="$(question)"

# The lane is finished once the ticket has been in it and left it. A worker is
# started before its ticket is moved into the lane, so a status other than the
# lane is not yet an exit: until the lane has been seen, only a status that is
# neither the starting one nor the lane counts — a worker that moved the ticket
# straight on. Once seen, any other status counts, the starting one included:
# that is critique sending work back. An empty status is a read that failed,
# never an exit — taking it for one would close the tab of a worker at work.
seen=""
asked=""
while :; do
  s="$(status)"
  if [ -n "$s" ]; then
    [ "$s" = "$lane" ] && seen=1
    if [ "$s" != "$lane" ] && { [ -n "$seen" ] || [ "$s" != "$start" ]; }; then break; fi
    if [ "$s" = "$lane" ] && { q="$(question)"; [ -n "$q" ] && [ "$q" != "$q0" ]; } \
      && case "$(agent)" in idle|done) true ;; *) false ;; esac; then
      asked=1; break
    fi
  fi
  check; sleep 20
done
until case "$(agent)" in idle|done) true ;; *) false ;; esac; do check; sleep 10; done

(cd "$board" && jaira show "$ticket" --json 2>/dev/null) | python3 -c '
import sys, json
d = json.load(sys.stdin)
print("status:", d.get("status"))
for k in ("outcome", "review", "question"):
    if d.get(k):
        print(k + ":", json.dumps(d[k], ensure_ascii=False)[:900])
' || true

if [ -z "$keep" ]; then
  # The person may have closed the tab already; that is the job done, not an error.
  tab="$("$herdr" pane get "$pane" 2>/dev/null \
    | python3 -c 'import sys,json;print(json.load(sys.stdin)["result"]["pane"]["tab_id"])' 2>/dev/null || true)"
  if [ -n "$tab" ]; then
    "$herdr" tab close "$tab" >/dev/null
    echo "closed $pane"
  else
    echo "$pane is already gone"
  fi
fi
git -C "$wt" log --oneline -3 2>/dev/null || true
if [ -n "$asked" ]; then
  echo "the worker stopped on a question; it is on the ticket — answer it, clear it with 'jaira set $ticket question=', then run the lane again" >&2
  exit 5
fi
