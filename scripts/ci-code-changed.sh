#!/usr/bin/env bash
# Says whether a pull request changes anything the tests or the binary read,
# for ci.yaml's 'changes' job. Prints code=true or code=false, and appends the
# same line to $GITHUB_OUTPUT when it is set.
#
# Most commits on this repository's board touch nothing but .jaira/ — a ticket
# changes lane, a note is written — and running the matrix on three platforms
# for them buys nothing. A pull request whose every changed file is on the list
# below runs no test job; one changed file anywhere else, and it runs them all.
# A ticket file rides in the same commit as its code, so such a pull request
# runs the full matrix.
#
# The list is what no build and no test reads. Not every markdown file is on
# it: core/** and internal/** are go:embed inputs, lanes/*.md is read by
# core/lane/shipped_test.go, .claude/skills/jaira/SKILL.md by
# core/role/jairaskill_test.go, and README.md and docs/COMMANDS.md by
# internal/cli/commands_doc_test.go. A test that starts reading a file on the
# list takes it off the list in the same change.
#
# .jaira/ is on it although core/ticket/decode_once_test.go reads the ticket
# files there: a ticket the CLI wrote is not what that test guards, and the next
# pull request that carries code runs it over the same files.
#
# Usage: ci-code-changed.sh <base> [<head>]   (head defaults to HEAD)
set -euo pipefail

base=${1:?usage: $0 <base> [<head>]}
head=${2:-HEAD}

# The merge base, not base itself: a master that moved on since the branch was
# cut would otherwise count master's own changes as this pull request's.
# NUL-separated, because without -z git quotes a path holding a non-ASCII byte
# ("\"...\""), and logbook file names carry umlauts: the quoted name would match
# no pattern below and count as code. The diff runs before the loop so that a
# failing git stops the script instead of reading as an empty diff.
# A name holding a newline splits into pieces that match nothing: tests run.
files=$(git diff --name-only -z --no-renames "$base...$head" | tr '\0' '\n')

code=false
while IFS= read -r f; do
	[ -z "$f" ] && continue
	case "$f" in
	docs/COMMANDS.md) ;; # read by a test: falls through to code=true below
	.jaira/* | .planning/* | docs/* | AGENTS.md | CLAUDE.md | lanes/README.md) continue ;;
	esac
	echo "needs tests: $f" >&2
	code=true
	break
done <<<"$files"

echo "code=$code"
if [ -n "${GITHUB_OUTPUT:-}" ]; then
	echo "code=$code" >>"$GITHUB_OUTPUT"
fi
