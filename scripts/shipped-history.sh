#!/usr/bin/env bash
# Prints core/shipped/shipped.sha256: every version of every file jaira
# installs that was ever committed, plus the working tree's own, as
# "sha256  path in this repository" under the path it is embedded from today.
#
#   scripts/shipped-history.sh > core/shipped/shipped.sha256
#
# Needs a full clone (a shallow one silently drops versions) and sha256sum
# (coreutils; Linux, WSL, Git Bash). The output is append-only in practice: a
# line that disappears turns a file jaira wrote back into "edited".
set -eu
cd "$(git rev-parse --show-toplevel)"
if [ "$(git rev-parse --is-shallow-repository)" != false ]; then
	echo "shipped-history: a shallow clone has no full history; run 'git fetch --unshallow' first" >&2
	exit 1
fi

# What the binary installs. The cockpit mod's tests/ folder is embedded for its
# own tests but never written to a skills directory (core/mod/mod.go).
embedded() { git ls-files core/role/builtin core/mod/builtin core/lane/builtin core/role/jaira | grep -v '/tests/'; }

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
	git log --full-history -m --raw --no-abbrev --no-renames --format= HEAD -- \
		core/role/builtin core/mod/builtin core/lane/builtin core/role/jaira lanes .claude/skills/jaira |
		awk '$4 !~ /^0+$/ { print $4, $6 }' | sort -u |
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
	# the working tree's own versions, so this script alone makes the tests pass
	embedded | while read -r f; do
		printf '%s  %s\n' "$(sha256sum <"$f" | cut -d' ' -f1)" "$f"
	done
} | awk 'NR==FNR { keep[$1]=1; next } keep[$2]' <(embedded) - | sort -k2,2 -k1,1 -u
