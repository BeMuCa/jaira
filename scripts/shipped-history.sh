#!/usr/bin/env bash
# Prints core/shipped/shipped.sha256: every version of every file jaira
# installs that was ever committed, plus the working tree's own, as
# "sha256  path in this repository" under the path it is embedded from today.
#
#   scripts/shipped-history.sh --write
#   scripts/shipped-history.sh --check
#   scripts/shipped-history.sh > somewhere-else
#
# --write replaces core/shipped/shipped.sha256 only once the whole manifest is
# generated. Never redirect into the manifest itself: the shell empties it
# before the script runs, and a run that then fails leaves it cut short.
#
# --check writes nothing: it holds the committed manifest against every v* tag
# and fails, printing the lines to append, when a version of an installed file
# that a release tag carried is not in it. jaira update --roles takes such a
# file, untouched since that release, for one a person edited and skips it.
#
# Needs a full clone (a shallow one silently drops versions) and sha256sum
# (coreutils; Linux, WSL, Git Bash). The output is append-only in practice: a
# line that disappears turns a file jaira wrote back into "edited".
set -euo pipefail
# One byte order on every machine: a manifest sorted under another locale
# reorders lines, and two such branches merge into duplicates.
export LC_ALL=C
cd "$(git rev-parse --show-toplevel)"
if [ "$(git rev-parse --is-shallow-repository)" != false ]; then
	echo "shipped-history: a shallow clone has no full history; run 'git fetch --unshallow' first" >&2
	exit 1
fi

manifest=core/shipped/shipped.sha256

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
		# Apart from the printf: a blob git cannot read (a partial clone offline,
		# a damaged object) must stop the script, not be hashed as the empty file.
		# Said with || exit: --check runs this inside $(...), where set -e is off.
		sum=$(git cat-file blob "$blob" | sha256sum) || exit 1
		printf '%s  %s\n' "${sum%% *}" "$path"
	done
}

# Keeps the lines whose path jaira installs today; a version of a file it no
# longer ships is never asked about.
installed() { awk 'NR==FNR { keep[$1]=1; next } keep[$2]' <(embedded) -; }

if [ "${1:-}" = --check ]; then
	tags=$(git tag -l 'v*')
	if [ -z "$tags" ]; then
		echo "shipped-history: this clone has no v* tags, so nothing could be checked; fetch them (git fetch --tags)" >&2
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
		printf 'shipped-history: %s lacks %d version(s) a release tag shipped; run scripts/shipped-history.sh --write, or append:\n%s\n' \
			"$manifest" "$n" "$missing" >&2
		exit 1
	fi
	echo "shipped-history: every version in $(printf '%s\n' "$tags" | wc -l) v* tag(s) is in $manifest"
	exit 0
fi

generate() {
	cat <<'HEADER'
# Every version of every file jaira installs — the roles, the cockpit mod, the jaira
# skill and the built-in lanes — as "sha256  path in this repository".
# A file on someone's disk equal to one of these was written by jaira and may be
# replaced by a newer one; anything else was edited by a person and is left alone.
# Append-only: a line removed here turns a file jaira wrote back into "edited".
# Regenerate (full clone, not shallow): scripts/shipped-history.sh --write
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
}

if [ "${1:-}" = --write ]; then
	# Beside the manifest, so the mv below is a rename on one file system.
	tmp=$(mktemp "$manifest.XXXXXX")
	trap 'rm -f "$tmp"' EXIT
	# On a line of its own: set -e does not reach a command left of && or ||.
	generate >"$tmp"
	chmod 644 "$tmp"
	mv "$tmp" "$manifest"
	echo "shipped-history: wrote $manifest ($(grep -cv '^#' "$manifest") versions)"
	exit 0
fi

generate
