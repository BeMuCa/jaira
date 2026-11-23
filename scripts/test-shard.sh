#!/usr/bin/env bash
# Runs one shard of the test suite, for CI runners where the whole suite is too
# slow to wait for (windows-latest: internal/cli alone takes ~8 minutes).
#
# The two slow packages, internal/cli and internal/tui, cannot run their tests
# in parallel inside one process — they share global state — so they are split
# across processes instead: shard i of n runs every n-th of their top-level
# tests, counted from i. Shard 0 also runs every other package whole. Together
# the n shards run every test exactly once.
#
# The go test runs of one shard — every other package, internal/cli,
# internal/tui — start at once, each a process of its own, so the shard takes
# as long as its slowest run rather than their sum. Their output is held and
# printed after all have finished, one run after the other, so a failure reads
# as one block. The exit code is the first non-zero one.
#
# Usage: test-shard.sh <index> <count> [go test flags...]
set -euo pipefail

i=${1:?usage: $0 <index> <count> [go test flags...]}
n=${2:?usage: $0 <index> <count> [go test flags...]}
shift 2

split=(./internal/cli ./internal/tui)

out=$(mktemp -d)
trap 'rm -rf "$out"' EXIT
pids=()
logs=()

if [ "$i" = 0 ]; then
  others=$(go list ./... | grep -v -e '/internal/cli$' -e '/internal/tui$')
  # shellcheck disable=SC2086 # one package per word
  go test "$@" $others >"$out/others" 2>&1 &
  pids+=($!)
  logs+=("$out/others")
fi

for pkg in "${split[@]}"; do
  mine=$(go test -list . "$pkg" | grep -E '^(Test|Example|Fuzz)' | awk -v n="$n" -v i="$i" '(NR - 1) % n == i' | paste -sd'|' -)
  if [ -z "$mine" ]; then
    continue
  fi
  log="$out/$(basename "$pkg")"
  echo "shard $i/$n: $(echo "$mine" | tr '|' '\n' | wc -l | tr -d ' ') tests of $pkg" >"$log"
  go test "$@" -run "^($mine)\$" "$pkg" >>"$log" 2>&1 &
  pids+=($!)
  logs+=("$log")
done

status=0
for k in "${!pids[@]}"; do
  rc=0
  wait "${pids[$k]}" || rc=$?
  cat "${logs[$k]}"
  if [ "$status" = 0 ] && [ "$rc" != 0 ]; then
    status=$rc
  fi
done
exit "$status"
