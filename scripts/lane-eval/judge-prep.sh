#!/usr/bin/env bash
# Prepares the blind comparison of one ticket's run in two arms.
#
#   judge-prep.sh <out> <ticket> <run> <armX> <armY>
#
# Writes two judge folders that differ only in which diff is called X:
#   <out>/judge/<ticket>-r<run>-p1/  X = <armX>, Y = <armY>
#   <out>/judge/<ticket>-r<run>-p2/  X = <armY>, Y = <armX>
# Each holds X.diff, Y.diff, ticket.md and, for a golden ticket, reference.diff
# (the accepted change). Nothing in a judge folder names an arm: which diff is
# which stands only in <out>/KEY.txt, one line per folder.
# A diff is the clone's working tree against its base, uncommitted files
# included, without .jaira/ (the ticket's bookkeeping is not the work).
set -euo pipefail
if [ $# -ne 5 ]; then sed -n '2,13p' "$0" >&2; exit 2; fi
out=$(cd "$1" && pwd); ticket=$2; run=$3; ax=$4; ay=$5

armdiff() {
  local meta="$out/$1/$ticket/r$run.json" dir base
  dir=$(jq -r .dir "$meta"); base=$(jq -r .base "$meta")
  git -C "$dir" add -A -N
  git -C "$dir" diff "$base" -- . ':(exclude).jaira'
}

for p in p1 p2; do
  if [ $p = p1 ]; then x=$ax; y=$ay; else x=$ay; y=$ax; fi
  j="$out/judge/$ticket-r$run-$p"
  rm -rf "$j"; mkdir -p "$j"
  armdiff "$x" > "$j/X.diff"
  armdiff "$y" > "$j/Y.diff"
  cp "$out/reference/$ticket.ticket.md" "$j/ticket.md"
  [ -s "$out/reference/$ticket.diff" ] && cp "$out/reference/$ticket.diff" "$j/reference.diff"
  key="$ticket r$run $p X=$x Y=$y"
  touch "$out/KEY.txt"
  grep -qxF "$key" "$out/KEY.txt" || echo "$key" >> "$out/KEY.txt"
  echo "$j"
done
