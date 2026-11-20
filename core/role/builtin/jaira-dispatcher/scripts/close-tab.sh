#!/usr/bin/env bash
# A worker or dispatcher closes its own tab once its ticket waits for a person
# and its report has gone to the session that started it.
# Usage: close-tab.sh [--dry-run] --parent <session> <ticket-id> [dir]
set -euo pipefail

usage() {
  cat <<'USAGE'
close-tab.sh [--dry-run] --parent <session> <ticket-id> [dir]

The last thing a worker or dispatcher does, after its report reached
<session> (the --parent it was started with): close the Herdr tab it runs in.
Nobody needs its context any more — a ticket given back gets a new worker, and
everything the old one knew is on the ticket.

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

  --dry-run   Print 'would close <ticket> <tab>' instead of closing.
  -h, --help  This text.

Output: 'closing <ticket> <tab>', or 'kept <ticket>: <reason>'. Exit 0 in both
cases — keeping is an answer, not a failure; 2 on a usage error.

Environment: HERDR_ENV, HERDR_BIN_PATH, HERDR_PANE_ID (Herdr), HC_TERMINAL_ID,
HC_BIN (Herdr App, through herdr-hc.sh).
USAGE
}

dry=""; parent=""
while [ $# -gt 0 ]; do
  case "$1" in
    --dry-run) dry=1; shift ;;
    --parent)  parent="${2-}"; shift 2 || shift $# ;;
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

kept() { echo "kept $t: $*"; exit 0; }
# What this script did not foresee - a 'jaira json get' that refuses (an older
# jaira without --default), a tool that is gone - keeps the tab too, instead of
# ending it under set -e with that command's exit code against the exit 0
# promised above. Broken JSON is not such a case: --default answers it.
trap 'kept "line $LINENO failed (exit $?)"' ERR

case "$parent" in
  "")     kept "no parent: the report is in this tab" ;;
  uds:*)  [ -e "${parent#uds:}" ] || kept "parent $parent does not listen: the report did not leave this tab" ;;
esac

root="$(git -C "$dir" rev-parse --show-toplevel 2>/dev/null)" || kept "$dir is not in a git repository"

rc=0; shown="$(jaira -C "$root" show "$t" --json 2>/dev/null)" || rc=$?
[ "$rc" = 0 ] || kept "jaira show $t failed (exit $rc)"
lane="$(jaira json get status --default "" <<<"$shown")"
board="$(jaira json get board --default "" <<<"$shown")"
question="$(jaira json get question --default "" <<<"$shown")"
[ -z "$question" ] || kept "the ticket has an open question"

# The lane's kind, off the ticket's own board.
lanes="$(jaira -C "$root" lanes ${board:+--board "$board"} --json 2>/dev/null)" || kept "jaira lanes failed"
persons=""
for i in $(seq 0 99); do
  id="$(jaira json get "lanes.$i.id" --default "" <<<"$lanes")"
  [ -n "$id" ] || break
  [ "$id" = "$lane" ] || continue
  for k in human_exit requires_question terminal; do
    [ "$(jaira json get "lanes.$i.$k" --default false <<<"$lanes")" = true ] && persons=1
  done
  break
done
[ -n "$persons" ] || kept "the ticket is in $lane, not with a person"

here="$(cd "$(dirname "$0")" && pwd)"
d="$("$here/uncommitted.sh" "$root")" || kept "git status failed in $root"
[ -z "$d" ] || kept "uncommitted changes: $d"

herdr="${HERDR_BIN_PATH:-herdr}"
if [ "${HERDR_ENV:-}" != 1 ] && [ -n "${HC_TERMINAL_ID:-}" ] \
  && command -v "${HC_BIN:-hc}" >/dev/null 2>&1; then
  herdr="$here/herdr-hc.sh"
fi
pane="${HERDR_PANE_ID:-${HC_TERMINAL_ID:-}}"
{ [ "${HERDR_ENV:-}" = 1 ] || [ "${herdr##*/}" = herdr-hc.sh ]; } && [ -n "$pane" ] \
  || kept "not in a Herdr tab"
tab="$("$herdr" pane get "$pane" 2>/dev/null | jaira json get result.pane.tab_id 2>/dev/null || true)"
[ -n "$tab" ] || kept "herdr does not know pane $pane"
label="$("$herdr" tab get "$tab" 2>/dev/null | jaira json get result.tab.label 2>/dev/null || true)"
# spawn.sh labels the tab with the id as it was handed in: a handle or a slug.
case "$(printf '%s' "$label" | tr '[:lower:]' '[:upper:]')" in
  "$t"/*) ;;
  *) kept "tab $tab is labelled '$label', not $t/<lane>: spawn.sh did not open it" ;;
esac

if [ -n "$dry" ]; then echo "would close $t $tab"; exit 0; fi
echo "closing $t $tab"
# A session of its own where there is setsid, so the end of this tab's
# processes does not take the close down half way (herdr-hc.sh kills, then
# removes the terminal).
if command -v setsid >/dev/null 2>&1; then
  setsid "$herdr" tab close "$tab" </dev/null >/dev/null 2>&1 &
else
  nohup "$herdr" tab close "$tab" </dev/null >/dev/null 2>&1 &
fi
exit 0
