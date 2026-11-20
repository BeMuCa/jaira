#!/usr/bin/env bash
# What a worktree holds that no commit and no ref does: the one rule close-tab.sh
# and reap.sh both ask before they close a tab or remove a worktree.
# Usage: uncommitted.sh <dir>
set -euo pipefail

case "${1:-}" in
  ""|-h|--help)
    cat <<'USAGE'
uncommitted.sh <dir>

Prints the first three paths of <dir>'s worktree that are not committed and
that no ref holds, separated by spaces; nothing when there are none.

Not counted: files under .jaira/tickets/ and .jaira/milestones/ that exist,
whether untracked, staged or modified. The board's refs hold them, and a lane
that changed no code leaves the ticket file modified until the next commit.

Counted: everything else, .jaira/ included. Lanes, tags and the logbook are on
no ref, and neither is a ticket deleted by a 'jaira logbook' not yet committed.

Paths are printed as they are, umlauts included (core.quotePath=false), so
ticket files with German titles match the rule.

Exit 0; 1 when git status fails; 2 without <dir>.
USAGE
    [ -n "${1:-}" ] && exit 0 || exit 2 ;;
esac

st="$(git -C "$1" -c core.quotePath=false status --porcelain -uall)" || exit 1
d="$(printf '%s\n' "$st" | grep -v -E '^$|^(\?\?|[ MA][ M]) \.jaira/(tickets|milestones)/' | cut -c4- | head -3 | tr '\n' ' ' || true)"
printf '%s\n' "${d% }"
