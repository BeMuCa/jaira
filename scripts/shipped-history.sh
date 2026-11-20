#!/usr/bin/env bash
# Prints core/shipped/shipped.sha256: every version of every file jaira
# installs that was ever committed, plus the working tree's own, as
# "sha256  path in this repository" under the path it is embedded from today.
#
#   scripts/shipped-history.sh > core/shipped/shipped.sha256
#   scripts/shipped-history.sh --check
#
# --check writes nothing: it holds the committed manifest against every v* tag
# and fails, printing the lines to append, when a version of an installed file
# that a release tag carried is not in it. jaira update --roles takes such a
# file, untouched since that release, for one a person edited and skips it.
#
# Needs a full clone (a shallow one silently drops versions) and sha256sum
# (coreutils; Linux, WSL, Git Bash). The output is append-only in practice: a
# line that disappears turns a file jaira wrote back into "edited".
set -eu
# One byte order on every machine: a manifest sorted under another locale
# reorders lines, and two such branches merge into duplicates.
export LC_ALL=C
cd "$(git rev-parse --show-toplevel)"
if [ "$(git rev-parse --is-shallow-repository)" != false ]; then
	echo "shipped-history: a shallow clone has no full history; run 'git fetch --unshallow' first" >&2
	exit 1
fi

# What the binary installs. The cockpit mod's tests/ folder is embedded for its
# own tests but never written to a skills directory (core/mod/mod.go).
embedded() { git ls-files core/role/builtin core/mod/builtin core/lane/builtin core/role/jaira | grep -v '/tests/'; }

# The directories a version of an installed file was ever committed under.
sources='core/role/builtin core/mod/builtin core/lane/builtin core/role/jaira lanes .claude/skills/jaira'

# Reads "<blob> <path>" lines and prints "<sha256>  <path>" under the path the
# file is embedded from today. The blob's bytes are hashed as git stores them,
# CRLF included: that is what a build of that commit embedded.
hashed() {
	while read -r blob path; do
		case $path in
		# critique, optimize and testing were catalogue lanes before 2fb8f177 moved them into the binary
		lanes/*.md) path=$(ls core/lane/builtin/*-"${path#lanes/}" 2>/dev/null | head -n 1) || true ;;
		# the jaira skill is embedded from a copy of the repository's own
		.claude/skills/jaira/SKILL.md) path=core/role/jaira/SKILL.md ;;
		esac
		[ -n "$path" ] || continue
		printf '%s  %s\n' "$(git cat-file blob "$blob" | sha256sum | cut -d' ' -f1)" "$path"
	done
}

# Keeps the lines whose path jaira installs today; a version of a file it no
# longer ships is never asked about.
installed() { awk 'NR==FNR { keep[$1]=1; next } keep[$2]' <(embedded) -; }

if [ "${1:-}" = --check ]; then
	manifest=core/shipped/shipped.sha256
	tags=$(git tag -l 'v*')
	if [ -z "$tags" ]; then
		echo "shipped-history: this clone has no v* tags, so nothing could be checked; fetch them (git fetch --tags)" >&2
		exit 1
	fi
	# A line that is no entry — a conflict marker left from a merge, say — is
	# a manifest nobody can trust, whatever else it holds.
	bad=$(tr -d '\r' <"$manifest" | grep -nvE '^(#.*|[0-9a-f]{64}  [^ ]+|)$' || true)
	if [ -n "$bad" ]; then
		printf 'shipped-history: %s has lines that are no entry:\n%s\n' "$manifest" "$bad" >&2
		exit 1
	fi
	# Listed apart from the pipeline below: a tag git cannot read must stop the
	# check, not leave it comparing fewer versions and passing.
	listing=$(
		for tag in $tags; do
			# a tag without one of these directories simply lists nothing for it
			# shellcheck disable=SC2086
			git -c core.quotePath=false ls-tree -r "$tag" -- $sources || exit 1
		done
	)
	missing=$(
		printf '%s\n' "$listing" | awk -F'\t' '{ split($1, m, " ") } m[2] == "blob" { print m[3], $2 }' |
			sort -u | hashed | installed | sort -u |
			awk 'NR==FNR { have[$0]=1; next } !have[$0]' <(tr -d '\r' <"$manifest") -
	)
	if [ -n "$missing" ]; then
		n=$(printf '%s\n' "$missing" | wc -l)
		printf 'shipped-history: %s lacks %d version(s) a release tag shipped; run scripts/shipped-history.sh > %s, or append:\n%s\n' \
			"$manifest" "$n" "$manifest" "$missing" >&2
		exit 1
	fi
	echo "shipped-history: every version in $(printf '%s\n' "$tags" | wc -l) v* tag(s) is in $manifest"
	exit 0
fi

cat <<'HEADER'
# Every version of every file jaira installs — the roles, the cockpit mod, the jaira
# skill and the built-in lanes — as "sha256  path in this repository".
# A file on someone's disk equal to one of these was written by jaira and may be
# replaced by a newer one; anything else was edited by a person and is left alone.
# Append-only: a line removed here turns a file jaira wrote back into "edited".
# Regenerate (full clone, not shallow): scripts/shipped-history.sh > core/shipped/shipped.sha256
HEADER

{
	# -m: a conflict resolution on a merge commit is a version that shipped too.
	# shellcheck disable=SC2086
	git -c core.quotePath=false log --full-history -m --raw --no-abbrev --no-renames --format= HEAD -- $sources |
		awk -F'\t' '{ split($1, m, " ") } m[4] !~ /^0+$/ { print m[4], $2 }' | sort -u | hashed
	# the working tree's own versions, so this script alone makes the tests pass
	embedded | while read -r f; do
		printf '%s  %s\n' "$(sha256sum <"$f" | cut -d' ' -f1)" "$f"
	done
} | installed | sort -k2,2 -k1,1 -u
