#!/usr/bin/env bash
# A worker or dispatcher closes its own tab once its ticket waits for a person
# and its report has gone to the session that started it — at once when the
# ticket is done, else once its prompt cache has gone cold (reap.sh asks).
# Usage: close-tab.sh [--dry-run] --parent <session> <ticket-id> [dir]
#        close-tab.sh [--dry-run] --reap --pane <pane> <ticket-id> <dir>
set -euo pipefail

usage() {
  cat <<'USAGE'
close-tab.sh [--dry-run] --parent <session> <ticket-id> [dir]
close-tab.sh [--dry-run] --reap --pane <pane> <ticket-id> <dir>

The last thing a worker or dispatcher does, after its report reached
<session> (the --parent it was started with): close the Herdr tab it runs in.

A ticket in the terminal lane closes the tab at once. A ticket in any other
lane a person owns (signoff, human) may still be given back, and a worker
whose prompt cache is warm answers that round for a fraction of what a new
one costs reading everything again. So there the tab stays: the script writes
parked=<epoch seconds> into the tab's record (jaira/panes/ in the git common
directory, written by spawn.sh) and prints
'kept <ticket>: cache warm until HH:MM, reap.sh closes it after'.

reap.sh calls it again with --reap for every record that says parked=, and
it closes the tab once the ticket is done or no longer on the board, or the
park is older than the cache window. A ticket given back to an agent's lane
keeps the tab: the same worker carries it again and parks it anew when it
finishes.

It closes only when all of these hold, and otherwise prints
'kept <ticket>: <reason>' and closes nothing:
  - --parent names a session: without one the report is in this tab, and
    closing it would take the report with it. A uds: address must exist.
  - the ticket sits in a lane a person owns (one that takes a question, or
    one only a person moves work out of) or in the terminal lane,
  - the ticket has no open question: a worker waiting for an answer stays,
  - the directory has no uncommitted changes, as uncommitted.sh beside this
    file counts them - the same rule reap.sh asks (ticket and milestone files
    ride with the next commit and are on the board's refs),
  - this runs in a tab spawn.sh labelled <ticket>/<lane>.

The tab closes in the background, detached from this process: closing it ends
the claude that ran this script.

  --dry-run     Print 'would close <ticket> <tab>' instead of closing, and
                write no parked= line.
  --reap        Asked by reap.sh, not by the worker: no --parent, the tab is
                --pane's, and a parked tab closes once the window is over.
                A tab whose agent is working stays. A ticket 'jaira show' does
                not find any more has been filed away and counts as done.
  --pane <id>   The pane whose tab to judge; only with --reap.
  -h, --help    This text.

Output: 'closing <ticket> <tab>', or 'kept <ticket>: <reason>'. Exit 0 in both
cases — keeping is an answer, not a failure; 2 on a usage error.

Environment: HERDR_ENV, HERDR_BIN_PATH, HERDR_PANE_ID (Herdr), HC_TERMINAL_ID,
HC_BIN (Herdr App, through herdr-hc.sh), JAIRA_CACHE_WINDOW (seconds a parked
tab stays, default 3600: the prompt cache's lifetime on a subscription).
USAGE
}

dry=""; parent=""; reap=""; rpane=""
while [ $# -gt 0 ]; do
  case "$1" in
    --dry-run) dry=1; shift ;;
    --parent)  parent="${2-}"; shift 2 || shift $# ;;
    --reap)    reap=1; shift ;;
    --pane)    rpane="${2-}"; shift 2 || shift $# ;;
    -h|--help) usage; exit 0 ;;
    --)        shift; break ;;
    -*)        echo "unknown flag: $1" >&2; usage >&2; exit 2 ;;
    *)         break ;;
  esac
done
t="${1:-}"
[ -n "$t" ] || { usage >&2; exit 2; }
t="$(printf '%s' "$t" | tr '[:lower:]' '[:upper:]')"
dir="${2:-$PWD}"
if [ -n "$reap" ]; then
  [ -n "$rpane" ] && [ -z "$parent" ] && [ -n "${2:-}" ] \
    || { echo "--reap takes --pane and a directory, and no --parent" >&2; usage >&2; exit 2; }
elif [ -n "$rpane" ]; then
  echo "--pane goes with --reap" >&2; usage >&2; exit 2
fi
window="${JAIRA_CACHE_WINDOW:-3600}"
case "$window" in ''|*[!0-9]*) echo "JAIRA_CACHE_WINDOW must be seconds, not '$window'" >&2; exit 2 ;; esac

kept() { echo "kept $t: $*"; exit 0; }
# What this script did not foresee - a 'jaira json get' that refuses (an older
# jaira without --default), a tool that is gone - keeps the tab too, instead of
# ending it under set -e with that command's exit code against the exit 0
# promised above. Broken JSON is not such a case: --default answers it.
trap 'kept "line $LINENO failed (exit $?)"' ERR

if [ -z "$reap" ]; then
  case "$parent" in
    "")     kept "no parent: the report is in this tab" ;;
    uds:*)  [ -e "${parent#uds:}" ] || kept "parent $parent does not listen: the report did not leave this tab" ;;
  esac
fi

root="$(git -C "$dir" rev-parse --show-toplevel 2>/dev/null)" || kept "$dir is not in a git repository"

rc=0; shown="$(jaira -C "$root" show "$t" --json 2>/dev/null)" || rc=$?
# Filed away by 'jaira logbook': off the board, which is done for reap.
filed=""
if [ "$rc" = 5 ] && [ -n "$reap" ]; then filed=1; shown='{}'
elif [ "$rc" != 0 ]; then kept "jaira show $t failed (exit $rc)"; fi
lane="$(jaira json get status --default "" <<<"$shown")"
board="$(jaira json get board --default "" <<<"$shown")"
question="$(jaira json get question --default "" <<<"$shown")"
[ -z "$question" ] || kept "the ticket has an open question"

# The lane's kind, off the ticket's own board.
lanes="$(jaira -C "$root" lanes ${board:+--board "$board"} --json 2>/dev/null)" || kept "jaira lanes failed"
persons=""; terminal="$filed"
for i in $(seq 0 99); do
  [ -z "$filed" ] || break
  id="$(jaira json get "lanes.$i.id" --default "" <<<"$lanes")"
  [ -n "$id" ] || break
  [ "$id" = "$lane" ] || continue
  for k in human_exit requires_question terminal; do
    [ "$(jaira json get "lanes.$i.$k" --default false <<<"$lanes")" = true ] || continue
    persons=1; [ "$k" != terminal ] || terminal=1
  done
  break
done
[ -n "$persons$terminal" ] || kept "the ticket is in $lane, not with a person"

here="$(cd "$(dirname "$0")" && pwd)"
d="$("$here/uncommitted.sh" "$root")" || kept "git status failed in $root"
[ -z "$d" ] || kept "uncommitted changes: $d"

herdr="${HERDR_BIN_PATH:-herdr}"
if [ "${HERDR_ENV:-}" != 1 ] && [ -n "${HC_TERMINAL_ID:-}" ] \
  && command -v "${HC_BIN:-hc}" >/dev/null 2>&1; then
  herdr="$here/herdr-hc.sh"
fi
pane="${rpane:-${HERDR_PANE_ID:-${HC_TERMINAL_ID:-}}}"
{ [ "${HERDR_ENV:-}" = 1 ] || [ "${herdr##*/}" = herdr-hc.sh ]; } && [ -n "$pane" ] \
  || kept "not in a Herdr tab"
# The record spawn.sh wrote for this pane, named as spawn.sh names it.
common="$(git -C "$root" rev-parse --git-common-dir)"
rec="$(cd "$root" && cd "$common" && pwd)/jaira/panes/$(printf '%s' "$pane" | tr -c 'A-Za-z0-9_.-' '_')"
# The record without its parked= line, or with a new one ($1). Written to a
# dot file beside it, which reap.sh's glob does not see, and moved over, so a
# reap reading it at the same moment sees the old record or the new one.
park() {
  [ -z "$dry" ] && [ -f "$rec" ] || return 0
  local tmp="${rec%/*}/.${rec##*/}.$$"
  { grep -v '^parked=' "$rec" || true; [ -z "${1:-}" ] || echo "parked=$1"; } >"$tmp" \
    && mv -f "$tmp" "$rec"
}
# A parked tab that is gone — closed by the person, or its id handed out
# again after a restart — is not asked about again on the next round.
gone() { [ -z "$reap" ] || park; kept "$@"; }
# Only herdr's own not_found (hc's kind, old herdr's pane_not_found) says the
# tab is gone. A herdr that gave no answer keeps the park: read as gone, the
# record would lose it and no later round would close the tab.
ask() {
  local e rc=0
  e="$(mktemp)"
  out="$("$herdr" "$@" 2>"$e")" || rc=$?
  err="$(tr -d '\r' <"$e")"; rm -f "$e"
  [ "$rc" != 0 ] || return 0
  case "$out $err" in *not_found*) out=""; return 0 ;; esac
  err="$(printf '%s\n%s' "$out" "$err" | sed '/^$/d' | tail -1)"
  kept "herdr gave no answer for $1 $3: ${err:-exit $rc, no message}"
}
ask pane get "$pane"
paned="$out"
tab="$(jaira json get result.pane.tab_id --default "" <<<"$paned" 2>/dev/null || true)"
[ -n "$tab" ] || gone "herdr does not know pane $pane"
ask tab get "$tab"
label="$(jaira json get result.tab.label --default "" <<<"$out" 2>/dev/null || true)"
# spawn.sh labels the tab with the id as it was handed in: a handle or a slug.
case "$(printf '%s' "$label" | tr '[:lower:]' '[:upper:]')" in
  "$t"/*) ;;
  *) gone "tab $tab is labelled '$label', not $t/<lane>: spawn.sh did not open it" ;;
esac

now="$(date +%s)"
if [ -z "$terminal" ]; then
  if [ -z "$reap" ]; then
    # Parked: the worker's last answer is now, so its cache lasts a window.
    [ -f "$rec" ] || kept "no record of pane $pane under jaira/panes: spawn.sh did not open it"
    until="$(date -d "@$((now + window))" +%H:%M 2>/dev/null || date -r "$((now + window))" +%H:%M 2>/dev/null || echo "+${window}s")"
    park "$now"
    kept "cache warm until $until, reap.sh closes it after"
  fi
  parked="$( (sed -n 's/^parked=//p' "$rec" 2>/dev/null || true) | head -1)"
  case "$parked" in ''|*[!0-9]*) kept "not parked: its worker did not finish in $lane" ;; esac
  left=$((parked + window - now))
  [ "$left" -le 0 ] || kept "cache warm for $(( (left + 59) / 60 )) more min"
fi
if [ -n "$reap" ] && [ "$(jaira json get result.pane.agent_status --default "" <<<"$paned" 2>/dev/null || true)" = working ]; then
  kept "its agent is working"
fi

if [ -n "$dry" ]; then echo "would close $t $tab"; exit 0; fi
echo "closing $t $tab"
[ -z "$reap" ] || park
# A session of its own where there is setsid, so the end of this tab's
# processes does not take the close down half way (herdr-hc.sh kills, then
# removes the terminal).
if command -v setsid >/dev/null 2>&1; then
  setsid "$herdr" tab close "$tab" </dev/null >/dev/null 2>&1 &
else
  nohup "$herdr" tab close "$tab" </dev/null >/dev/null 2>&1 &
fi
exit 0
