#!/usr/bin/env bash
# Prepares the blind comparison of one ticket's run in two arms.
#
#   judge-prep.sh <out> <ticket> <run> <armX> <armY>
#
# Writes two judge folders that differ only in which diff is called X:
#   p1  X = <armX>, Y = <armY>
#   p2  X = <armY>, Y = <armX>
# Each holds X.diff, Y.diff, ticket.md and, for a golden ticket, reference.diff
# (the accepted change), plus X.ticket.md and Y.ticket.md: what each arm wrote
# onto its ticket - the definition of done with its proofs, the outcome - and
# the commits it made, with their subjects. Notes, author and lane history stay
# out; they would tell the arms apart. The folders are fresh temp directories outside <out>
# (mktemp, under $TMPDIR or /tmp), so neither their path nor anything next to
# them names an arm, <out>, the clones or KEY.txt; clone.sh refuses an <out>
# under that temp root. Which diff is which, and
# where each folder is, stands only in <out>/KEY.txt, one line per folder.
# Prints "p1 <folder>" and "p2 <folder>". Running it again for the same ticket
# and run writes new folders and replaces their lines in KEY.txt.
# A diff is the clone's working tree against its base, uncommitted files
# included, without .jaira/ (the ticket's bookkeeping is not the work).
set -euo pipefail
if [ $# -ne 5 ]; then sed -n '2,21p' "$0" >&2; exit 2; fi
out=$(cd "$1" && pwd); ticket=$2; run=$3; ax=$4; ay=$5

armdiff() {
  local meta="$out/$1/$ticket/r$run.json" dir base
  dir=$(jq -r .dir "$meta"); base=$(jq -r .base "$meta")
  git -C "$dir" add -A -N
  git -C "$dir" diff "$base" -- . ':(exclude).jaira'
}

# The arm's ticket as a judge needs it: what it claims and what proves it. The
# commits are the clone's own since base: jaira derives a ticket's commit list
# only when it is filed away, so on the board the field is still empty.
armticket() {
  local meta="$out/$1/$ticket/r$run.json" dir handle base
  dir=$(jq -r .dir "$meta"); handle=$(jq -r .handle "$meta"); base=$(jq -r .base "$meta")
  jaira -C "$dir" show "$handle" --json 2>/dev/null | jq -r '
    "# What this attempt recorded on its ticket\n",
    "## Definition of Done\n",
    (.dod_items[] | "- [\(if .done then "x" else " " end)] \(.text)\n  proof: \(if .proof == "" then "(none)" else .proof end)"),
    "\n## Outcome\n", "what: \(.outcome.what)", "why: \(.outcome.why)", "resolves: \(.outcome.resolves)",
    "\n## Commits\n"'
  git -C "$dir" log --reverse --format='- %h %s' "$base..HEAD"
}

touch "$out/KEY.txt"
grep -v "^$ticket r$run p[12] " "$out/KEY.txt" > "$out/KEY.txt.new" || true
mv "$out/KEY.txt.new" "$out/KEY.txt"
for p in p1 p2; do
  if [ $p = p1 ]; then x=$ax; y=$ay; else x=$ay; y=$ax; fi
  j=$(mktemp -d "${TMPDIR:-/tmp}/judge.XXXXXXXX")
  armdiff "$x" > "$j/X.diff"
  armdiff "$y" > "$j/Y.diff"
  armticket "$x" > "$j/X.ticket.md"
  armticket "$y" > "$j/Y.ticket.md"
  cp "$out/reference/$ticket.ticket.md" "$j/ticket.md"
  [ -s "$out/reference/$ticket.diff" ] && cp "$out/reference/$ticket.diff" "$j/reference.diff"
  echo "$ticket r$run $p X=$x Y=$y dir=$j" >> "$out/KEY.txt"
  echo "$p $j"
done
