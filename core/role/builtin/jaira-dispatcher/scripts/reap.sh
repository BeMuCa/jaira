#!/usr/bin/env bash
# Close what spawn.sh opened for a ticket once its pull request is merged: every
# tab it started, the ticket's worktree and its local branch.
# Usage: reap.sh [--dry-run] [repo-root]
set -euo pipefail

usage() {
  cat <<'USAGE'
reap.sh [--dry-run] [repo-root]

Reads the tabs spawn.sh recorded under the repository's git directory
(jaira/panes/) and, for every ticket whose pull request is merged, closes its
tabs, removes its worktree with 'git worktree remove' and deletes its local
branch. Nothing else is ever touched: a tab spawn.sh did not start has no
record, and neither the main worktree nor the one this runs in is removed.

A ticket whose branch has no merged pull request yet (gh on GitHub, glab on
GitLab; 'git config jaira.forge' settles a host that names neither) is left
alone as 'waiting'. A merged one is kept, with the reason printed, when
  - its local branch is not at the merged head: commits past it, or the head
    never fetched,
  - one of its worktrees has changes that are not committed (untracked files
    under .jaira/ do not count),
  - the ticket is back on the board in a lane other than done.

  --dry-run   Print what would be closed and removed, change nothing.
  -h, --help  This text.

Output: one line per thing done — 'closed <ticket>/<lane> <pane>', 'removed
<worktree>', 'deleted branch <branch>' — or 'waiting <ticket>: <reason>' and
'kept <ticket>: <reason>'; with
--dry-run the first three start with 'would'. Exit 0 also when something was
kept: keeping is the answer, not a failure.

Environment: HERDR_BIN_PATH, HERDR_PANE_ID (this pane is never closed),
JAIRA_BRANCH_PREFIX (default feat).
USAGE
}

dry=""
while [ $# -gt 0 ]; do
  case "$1" in
    --dry-run) dry=1; shift ;;
    -h|--help) usage; exit 0 ;;
    --)        shift; break ;;
    -*)        echo "unknown flag: $1" >&2; usage >&2; exit 2 ;;
    *)         break ;;
  esac
done

root="$(cd "${1:-$PWD}" && pwd -P)"
herdr="${HERDR_BIN_PATH:-herdr}"
common="$(git -C "$root" rev-parse --git-common-dir)"
panes="$(cd "$root" && cd "$common" && pwd)/jaira/panes"
[ -d "$panes" ] || exit 0
# Physical paths throughout: on macOS the temp directory is a symlink, and git
# reports the resolved path where cd keeps the one it was given.
main="$(cd "$(git -C "$root" worktree list --porcelain | sed -n '1s/^worktree //p')" && pwd -P)"

# The forge decides which tool asks about the pull request — the same ladder
# jaira-role-pr walks: jaira.forge wins, then the host of origin.
forge="$(git -C "$root" config jaira.forge 2>/dev/null || true)"
if [ -z "$forge" ]; then
  case "$(git -C "$root" remote get-url origin 2>/dev/null || true)" in
    *github.com*) forge=github ;;
    *gitlab*)     forge=gitlab ;;
  esac
fi

# The head commit of the branch's merged pull request, or nothing. 'gh pr view
# <branch>' finds nothing in a fork clone, where the pull request lives on the
# upstream; 'gh pr list --head' finds it there as well.
merged_head() {
  case "$forge" in
    github) (cd "$root" && gh pr list --head "$1" --state merged --limit 1 --json headRefOid) \
              | jaira json get 0.headRefOid --default "" ;;
    gitlab) (cd "$root" && glab mr list --source-branch "$1" --merged -F json) \
              | jaira json get 0.sha --default "" ;;
  esac
}

field() { sed -n "s/^$1=//p" "$2" | head -1; }

# A worker passes its slug, in lower case, where a teamlead passes the handle:
# both name the same ticket. No associative array: macOS still ships bash 3.2.
ticket_of() { field ticket "$1" | tr '[:lower:]' '[:upper:]'; }

for t in $(for f in "$panes"/*; do [ -f "$f" ] && ticket_of "$f"; done | sort -u); do
  recs=()
  for f in "$panes"/*; do
    [ -f "$f" ] && [ "$(ticket_of "$f")" = "$t" ] && recs+=("$f")
  done

  # Its worktrees, apart from the main one, and the branch they have out.
  wts=(); branch=""
  for f in "${recs[@]}"; do
    w="$(field worktree "$f")"
    [ -d "$w" ] || continue
    w="$(cd "$w" && pwd -P)"
    # Neither the main worktree nor the one this runs in is ever removed.
    { [ "$w" = "$main" ] || [ "$w" = "$root" ]; } && continue
    case " ${wts[*]-} " in *" $w "*) continue ;; esac
    wts+=("$w")
    [ -n "$branch" ] || branch="$(git -C "$w" symbolic-ref --short HEAD 2>/dev/null || true)"
  done
  [ -n "$branch" ] || branch="${JAIRA_BRANCH_PREFIX:-feat}/$(printf '%s' "$t" | tr '[:upper:]' '[:lower:]')"

  if [ -z "$forge" ]; then
    echo "kept $t: cannot tell the forge from origin; set it with 'git config jaira.forge github' (or gitlab)"; continue
  fi
  if ! head="$(merged_head "$branch" 2>/dev/null)"; then
    echo "kept $t: asking $forge about $branch failed"; continue
  fi
  # Work still under way, the usual case: named, but not as something kept.
  if [ -z "$head" ]; then echo "waiting $t: $branch has no merged pull request"; continue; fi

  # A merged pull request is not the end when the ticket came back: a person
  # gave it back and its worker answers there. Filed away reads as done; a
  # ticket the board does not know is not waiting for anything.
  rc=0; shown="$(jaira -C "$root" show "$t" --json 2>/dev/null)" || rc=$?
  if [ "$rc" = 0 ]; then
    lane="$(jaira json get status --default "" <<<"$shown")"
    if [ "$lane" != done ]; then echo "kept $t: the ticket is in $lane, not done"; continue; fi
  elif [ "$rc" != 5 ]; then
    echo "kept $t: jaira show $t failed (exit $rc)"; continue
  fi

  # A squash or rebase merge leaves the branch unmerged in git's eyes, so the
  # pull request's head is the measure: anything past it was never merged. A
  # head that was never fetched cannot be compared either, and is kept alike.
  if tip="$(git -C "$root" rev-parse --verify -q "refs/heads/$branch")" \
    && [ "$tip" != "$head" ] && ! git -C "$root" merge-base --is-ancestor "$tip" "$head" 2>/dev/null; then
    echo "kept $t: $branch is not at the merged head ${head:0:12}"; continue
  fi

  dirty=""
  for w in "${wts[@]+"${wts[@]}"}"; do
    d="$(git -C "$w" status --porcelain | grep -v '^?? \.jaira/' | cut -c4- | head -3 | tr '\n' ' ' || true)"
    if [ -n "$d" ]; then dirty="$w has uncommitted changes: $d"; break; fi
  done
  if [ -n "$dirty" ]; then echo "kept $t: ${dirty% }"; continue; fi

  failed=""
  for f in "${recs[@]}"; do
    pane="$(field pane "$f")"; lane="$(field lane "$f")"
    if [ "$pane" = "${HERDR_PANE_ID:-}" ]; then continue; fi
    if [ -n "$dry" ]; then echo "would close $t/$lane $pane"; continue; fi
    # The person may have closed the tab already; that is the job done.
    tab="$("$herdr" pane get "$pane" 2>/dev/null | jaira json get result.pane.tab_id 2>/dev/null || true)"
    if [ -n "$tab" ]; then
      if ! "$herdr" tab close "$tab" >/dev/null 2>&1; then
        echo "kept $t: closing $pane failed"; failed=1; break
      fi
      echo "closed $t/$lane $pane"
    fi
    rm -f "$f"
  done
  # A tab still open has claude sitting in the worktree: leave it standing.
  [ -z "$failed" ] || continue

  for w in "${wts[@]+"${wts[@]}"}"; do
    if [ -n "$dry" ]; then echo "would remove $w"; continue; fi
    # --force for the untracked .jaira/ files the check above lets through.
    git -C "$root" worktree remove --force "$w" >/dev/null && echo "removed $w"
  done
  # A branch checked out in a worktree that stays — the main one, or the one
  # this runs in — stays with it.
  if git -C "$root" show-ref -q --verify "refs/heads/$branch" \
    && [ "$(git -C "$root" symbolic-ref -q --short HEAD)" != "$branch" ] \
    && [ "$(git -C "$main" symbolic-ref -q --short HEAD)" != "$branch" ]; then
    if [ -n "$dry" ]; then echo "would delete branch $branch"
    # -D, because a squash merge is no merge to git; the head check above
    # stands in for the one -d would make.
    else git -C "$root" branch -D "$branch" >/dev/null && echo "deleted branch $branch"; fi
  fi
done
exit 0
