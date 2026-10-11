#!/usr/bin/env bash
# Lays out one isolated run of the lane eval: a git clone with no remote, at a
# fixed base commit, holding the ticket in a state nobody has worked yet.
#
#   clone.sh <repo> <out> <arm> <ticket> <run> [base]
#
# The clone lands at <out>/<arm>/<ticket>/r<run> and is never removed. Its
# metadata is printed as JSON and written next to it (r<run>.json), outside the
# clone, so nothing in the tree an agent works tells it which arm it is in.
#
# Where the ticket comes from:
#   - on the board at <base>: used as it is, assigned to whoever runs the eval
#   - in .jaira/logbook/ at the repo's HEAD (a golden ticket): created afresh
#     from its title, goal, context and definition of done, in todo, so its
#     plan, notes, outcome and review do not leak into the run. That gives it
#     a new handle, which the metadata records. <base> defaults to the parent of
#     the oldest commit naming the ticket, and the accepted diff is written to
#     <out>/reference/<ticket>.diff for the judge.
# Running it again for a clone that already exists prints the stored metadata.
set -euo pipefail

if [ $# -lt 5 ]; then
  sed -n '2,20p' "$0" >&2
  exit 2
fi
repo=$(cd "$1" && pwd); out=$2; arm=$3; ticket=$4; run=$5; base=${6:-}
mkdir -p "$out"; out=$(cd "$out" && pwd)
# judge-prep.sh puts the judge folders in the temp root and refuses an <out>
# there; refusing it here already saves the hour the arms would run first.
case "$(cd "$out" && pwd -P)/" in
  "$(cd "${TMPDIR:-/tmp}" && pwd -P)"/*) echo "clone.sh: <out> $out lies under the temp root, where the judge folders go; use an out under \$HOME" >&2; exit 2 ;;
esac
dir="$out/$arm/$ticket/r$run"
meta="$out/$arm/$ticket/r$run.json"

if [ -f "$meta" ]; then cat "$meta"; exit 0; fi
if [ -e "$dir" ]; then echo "clone.sh: $dir exists without $meta; remove it by hand" >&2; exit 1; fi

who=$(git -C "$repo" config user.name)
refdir="$out/reference"; mkdir -p "$refdir"

# A golden ticket: the logbook file at the source repo's HEAD.
logbook=$(git -C "$repo" ls-tree -r --name-only HEAD .jaira/logbook | grep -E "/[0-9A-Z]*${ticket}-[^/]*\.md$" | head -1 || true)

if [ -z "$base" ]; then
  [ -n "$logbook" ] || { echo "clone.sh: $ticket is not in the logbook; pass a base commit" >&2; exit 2; }
  first=$(git -C "$repo" log --all --reverse --format=%H --grep="($ticket)" | head -1)
  [ -n "$first" ] || { echo "clone.sh: no commit names $ticket; pass a base commit" >&2; exit 2; }
  base=$(git -C "$repo" rev-parse "$first^1")
fi
base=$(git -C "$repo" rev-parse "$base^{commit}")

# jaira in the clone, with its real errors on stderr but without the notice
# that the board predates this binary: every board at an old base says that.
j() { jaira -C "$dir" "$@" 2> >(grep -v 'set up by an older version of jaira' >&2); }

git clone -q --no-checkout "$repo" "$dir"
git -C "$dir" remote remove origin
git -C "$dir" checkout -q -b "lane-eval/$arm-$ticket-r$run" "$base"

on_board=$(git -C "$dir" ls-tree -r --name-only HEAD .jaira/tickets | grep -E "/[0-9A-Z]*${ticket}-[^/]*\.md$" || true)
if [ -n "$on_board" ]; then
  handle=$ticket
  j set "$handle" "assignee=$who" >/dev/null
elif [ -n "$logbook" ]; then
  # Bring the filed ticket back only to read it, then file a clean one.
  mkdir -p "$dir/.jaira/logbook/lane-eval"
  src="$dir/.jaira/logbook/lane-eval/$(basename "$logbook")"
  git -C "$repo" show "HEAD:$logbook" > "$src"
  j restore "$src" >/dev/null
  old=$(j show "$ticket" --json)
  j delete "$ticket" --force >/dev/null
  rm -rf "$dir/.jaira/logbook/lane-eval"
  args=(create "$(jq -r .title <<<"$old")" --lane todo --assignee "$who"
        --goal "$(jq -r .goal <<<"$old")"
        --context "$(jq -r .context <<<"$old")")
  while IFS= read -r d; do args+=(--dod "$d"); done < <(jq -r '.dod_items[].text' <<<"$old")
  handle=$(j "${args[@]}" --json | jq -r .handle)
  [ -n "$handle" ] && [ "$handle" != null ] || { echo "clone.sh: could not create the ticket in $dir" >&2; exit 1; }
  if [ ! -f "$refdir/$ticket.diff" ]; then
    for c in $(jq -r '.commits[]' <<<"$old"); do
      git -C "$repo" show --format='commit %H%n%s%n' "$c" -- . ':(exclude).jaira'
    done > "$refdir/$ticket.diff"
  fi
else
  echo "clone.sh: $ticket is neither on the board at $base nor in the logbook" >&2; exit 2
fi

# What the judge reads about the ticket: the demand, not the record of the work.
if [ ! -f "$refdir/$ticket.ticket.md" ]; then
  j show "$handle" --json | jq -r '"# " + .title + "\n\n## Goal\n" + .goal + "\n\n## Context\n" + .context + "\n\n## Definition of Done\n" + ([.dod_items[] | "- " + .text] | join("\n"))' > "$refdir/$ticket.ticket.md"
fi

# Where the ticket stands, and which lanes have an agent step: the workflow
# decides from these which lane a lane-per-agent worker gets next.
lanes=$(j lanes --json | jq -c '[.lanes[] | select(.agentic) | .id]')
shown=$(j show "$handle" --json)
# A ticket still in backlog is walked to the lane right before the first agent
# step, so every arm starts from the same place.
for _ in 1 2 3; do
  st=$(jq -r .status <<<"$shown"); nx=$(jq -r .next_lane <<<"$shown")
  case $st in backlog|todo) ;; *) break ;; esac
  jq -e --arg l "$nx" 'index($l) != null' <<<"$lanes" >/dev/null && break
  j move "$handle" --to "$nx" >/dev/null || break
  shown=$(j show "$handle" --json)
done
jq -n --arg arm "$arm" --arg ticket "$ticket" --argjson run "$run" --arg dir "$dir" \
  --arg base "$base" --arg handle "$handle" --argjson shown "$shown" --argjson lanes "$lanes" \
  '{arm:$arm, ticket:$ticket, run:$run, dir:$dir, base:$base, handle:$handle,
    lane:$shown.status, next_lane:$shown.next_lane, agent_lanes:$lanes}' | tee "$meta"
