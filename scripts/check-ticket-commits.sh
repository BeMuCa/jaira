#!/usr/bin/env bash
# Refuses a commit that adds or changes a ticket file under .jaira/tickets/
# without changing anything outside .jaira/.
#
# A ticket nobody is working lives on its ref (refs/jaira/tickets/<id>), not
# in a branch. A commit that carries only ticket files puts such a ticket into
# master, where it lies as a stale copy beside its ref — that is how master
# came to hold 32 backlog tickets. A ticket rides with the code it belongs to,
# so a commit that carries a ticket file must carry code too.
#
# Removing a ticket file is always allowed: that is 'jaira logbook' and
# 'jaira archive' filing it away, and it is how stale copies are cleaned up.
#
# Usage: check-ticket-commits.sh <base> [<head>]   (head defaults to HEAD)
set -euo pipefail

base=${1:?usage: $0 <base> [<head>]}
head=${2:-HEAD}

bad=0
for c in $(git rev-list --no-merges "$base..$head"); do
	files=$(git diff-tree --no-commit-id --name-status -r --no-renames --root "$c")
	tickets=$(awk -F'\t' '$1 != "D" && $2 ~ /^\.jaira\/tickets\// {print $2}' <<<"$files")
	[ -z "$tickets" ] && continue
	code=$(awk -F'\t' '$2 !~ /^\.jaira\// {print $2}' <<<"$files")
	[ -n "$code" ] && continue
	echo "$(git log -1 --format='%h %s' "$c")"
	sed 's/^/  carries only ticket files: /' <<<"$tickets"
	bad=1
done

if [ "$bad" = 1 ]; then
	echo
	echo "A ticket file goes into a commit only together with code outside .jaira/."
	echo "A ticket nobody works yet lives on its ref; 'jaira pull <id>' brings it here when work starts."
	exit 1
fi
