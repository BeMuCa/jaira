#!/usr/bin/env bash
# Close what spawn.sh opened for a ticket once the ticket is done and its branch is
# pushed: the tabs and the worktree go then, the local branch once it is merged.
# Usage: reap.sh [--dry-run] [repo-root]
set -euo pipefail

usage() {
  cat <<'USAGE'
reap.sh [--dry-run] [repo-root]

Reads the tabs spawn.sh recorded under the repository's git directory
(jaira/panes/) and works through every ticket they name.

A ticket that is not done is 'waiting': it is still worked, or was given back.
A done one loses its tabs and its worktree, with 'git worktree remove', once
  - the worktree has no uncommitted changes (untracked files under .jaira/ and
    the ticket's own file do not count),
  - its HEAD is not detached,
  - its branch is pushed, and origin's copy is at the same commit.
The pull request is not waited for: the commits are on origin. One record stays
behind, naming the branch, so the next round still knows which branch it is.

The local branch goes once it is merged: its own pull request is merged at its
head (gh on GitHub, glab on GitLab; 'git config jaira.forge' settles a host that
names neither), or its head lies in origin's default branch (a milestone pull
request, merged into the integration branch). A squash merge the forge does not
report keeps the branch, which is safe: its commits are on origin.

Nothing else is ever touched: a tab spawn.sh did not start has no record, and
neither the main worktree nor the one this runs in is removed. A branch checked
out in either of them stays with it.

A thing is kept, with the reason printed, when
  - its worktree has changes that are not committed, or a detached HEAD,
  - its branch is not pushed, or origin's copy is at another commit,
  - its branch is not at the merged head,
  - herdr does not answer, so an open tab cannot be told from a closed one.

  --dry-run   Print what would be closed and removed, change nothing.
  -h, --help  This text.

Output: one line per thing done — 'closed <ticket>/<lane> <pane>', 'removed
<worktree>', 'deleted branch <branch>' — or 'waiting <ticket>: <reason>' and
'kept <ticket>: <reason>'; with --dry-run the first three start with 'would'.
Exit 0 also when something was kept: keeping is the answer, not a failure.

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

# The worktree's top, not the directory: started in a subdirectory, it would
# otherwise not know the worktree it stands in.
root="$(cd "$(git -C "${1:-$PWD}" rev-parse --show-toplevel)" && pwd -P)"
herdr="${HERDR_BIN_PATH:-herdr}"
common="$(git -C "$root" rev-parse --git-common-dir)"
panes="$(cd "$root" && cd "$common" && pwd)/jaira/panes"
# No record, nothing to ask the forge about: the round costs nothing.
ls "$panes"/* >/dev/null 2>&1 || exit 0
# A herdr that does not answer cannot tell a closed tab from an open one; read
# as closed, it would cost the worktree and the record of a tab still open.
herdr_up=""; "$herdr" tab list >/dev/null 2>&1 && herdr_up=1
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

# Where the pull requests are, as jaira-role-pr settles it: origin's own
# repository, or its parent when origin is a fork. Named, never left to gh: in
# a fork clone without an upstream remote gh asks the fork, which has none.
repo=""
if [ "$forge" = github ] && origin="$(git -C "$root" remote get-url origin 2>/dev/null)" \
  && view="$(gh repo view "$origin" --json isFork,parent,nameWithOwner 2>/dev/null)"; then
  if [ "$(jaira json get isFork --default false <<<"$view")" = true ]; then
    repo="$(jaira json get parent.owner.login <<<"$view")/$(jaira json get parent.name <<<"$view")"
  else
    repo="$(jaira json get nameWithOwner <<<"$view")"
  fi
fi

# The head commit of the branch's merged pull request, or nothing. 'gh pr view
# <branch>' finds nothing in a fork clone, where the pull request lives on the
# upstream; 'gh pr list --head' finds it there as well.
merged_head() {
  case "$forge" in
    github) (cd "$root" && gh pr list --head "$1" --state merged --limit 1 --json headRefOid ${repo:+--repo "$repo"}) \
              | jaira json get 0.headRefOid --default "" ;;
    gitlab) (cd "$root" && glab mr list --source-branch "$1" --merged -F json) \
              | jaira json get 0.sha --default "" ;;
  esac
}

field() { sed -n "s/^$1=//p" "$2" | head -1; }

# A worker passes its slug, in lower case, where a teamlead passes the handle:
# both name the same ticket. No associative array: macOS still ships bash 3.2.
ticket_of() { field ticket "$1" | tr '[:lower:]' '[:upper:]'; }

# Drops the ticket's records, but not this pane's own: its tab is never closed
# here, and its record goes when the pane does.
forget() {
  for f in "${recs[@]}"; do
    own="$(field pane "$f")"
    if [ -n "$own" ] && [ "$own" = "${HERDR_PANE_ID:-}" ]; then continue; fi
    rm -f "$f"
  done
  return 0
}

for t in $(for f in "$panes"/*; do [ -f "$f" ] && ticket_of "$f"; done | sort -u); do
  recs=()
  for f in "$panes"/*; do
    [ -f "$f" ] && [ "$(ticket_of "$f")" = "$t" ] && recs+=("$f")
  done

  # The lane comes first. A ticket that is still worked, or was given back to
  # the board, keeps everything, whatever its branch holds: its worker answers
  # there. Filed away reads as done; a ticket the board does not know is not
  # waiting for anything.
  rc=0; shown="$(jaira -C "$root" show "$t" --json 2>/dev/null)" || rc=$?
  if [ "$rc" = 0 ]; then
    now="$(jaira json get status --default "" <<<"$shown")"
    if [ "$now" != done ]; then echo "waiting $t: the ticket is in $now, not done"; continue; fi
  elif [ "$rc" != 5 ]; then
    echo "kept $t: jaira show $t failed (exit $rc)"; continue
  fi

  # Its worktrees, apart from the main one, and the branch they have out — or
  # the branch a record names, which outlives the worktree.
  wts=(); branch=""
  for f in "${recs[@]}"; do
    [ -n "$branch" ] || branch="$(field branch "$f")"
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

  # The tabs and the worktrees: only when there is one to close or remove.
  tabs=""
  for f in "${recs[@]}"; do
    pane="$(field pane "$f")"
    if [ -n "$pane" ] && [ "$pane" != "${HERDR_PANE_ID:-}" ]; then tabs=1; fi
  done
  if [ -n "${wts[*]-}" ] || [ -n "$tabs" ]; then
    dirty=""
    for w in "${wts[@]+"${wts[@]}"}"; do
      # A detached HEAD — a rebase paused, a commit checked out — may carry
      # commits no branch holds, and removing the worktree would lose them.
      if ! git -C "$w" symbolic-ref -q HEAD >/dev/null; then dirty="$w has a detached HEAD"; break; fi
      if ! st="$(git -C "$w" status --porcelain)"; then dirty="git status failed in $w"; break; fi
      # Untracked files under .jaira/ are the board's; the ticket's own file is
      # the one a worker leaves modified until its next commit that carries code.
      d="$(printf '%s' "$st" | grep -v -E "^(\?\? \.jaira/|.. \.jaira/tickets/[^/]*$t-)" | cut -c4- | head -3 | tr '\n' ' ' || true)"
      if [ -n "$d" ]; then dirty="$w has uncommitted changes: $d"; break; fi
    done
    if [ -n "$dirty" ]; then echo "kept $t: ${dirty% }"; continue; fi

    # Pushed means origin's copy sits at the local head: the commits are safe
    # there, and nothing is left that only this worktree holds. Read from the
    # local tracking ref, so no fetch and no network.
    tip="$(git -C "$root" rev-parse --verify -q "refs/heads/$branch" || true)"
    up="$(git -C "$root" rev-parse --verify -q "refs/remotes/origin/$branch" || true)"
    if [ -z "$up" ]; then echo "kept $t: $branch is not pushed"; continue; fi
    if [ "$up" != "$tip" ]; then
      echo "kept $t: $branch is not pushed at its head (origin at ${up:0:12}, local at ${tip:0:12})"; continue
    fi
    if [ -z "$herdr_up" ]; then echo "kept $t: herdr not reachable ($herdr)"; continue; fi

    failed=""
    for f in "${recs[@]}"; do
      pane="$(field pane "$f")"; lane="$(field lane "$f")"
      # A record with no pane names a branch only: there is no tab to close.
      if [ -z "$pane" ] || [ "$pane" = "${HERDR_PANE_ID:-}" ]; then continue; fi
      # The person may have closed the tab already; that is the job done. A
      # record outlives its tab by days, and herdr hands its ids out again after
      # a restart: a tab not labelled as spawn.sh labelled ours is somebody
      # else's, and ours is gone.
      tab="$("$herdr" pane get "$pane" 2>/dev/null | jaira json get result.pane.tab_id 2>/dev/null || true)"
      if [ -n "$tab" ] && [ "$("$herdr" tab get "$tab" 2>/dev/null | jaira json get result.tab.label 2>/dev/null || true)" \
        = "$(field ticket "$f")/$lane" ]; then
        if [ -n "$dry" ]; then echo "would close $t/$lane $pane"; continue; fi
        if ! "$herdr" tab close "$tab" >/dev/null 2>&1; then
          echo "kept $t: closing $pane failed"; failed=1; break
        fi
        echo "closed $t/$lane $pane"
      fi
    done
    # A tab still open has claude sitting in the worktree: leave it standing.
    [ -z "$failed" ] || continue

    for w in "${wts[@]+"${wts[@]}"}"; do
      if [ -n "$dry" ]; then echo "would remove $w"; continue; fi
      # --force for the untracked .jaira/ files and the ticket's own modified
      # file, which the check above lets through.
      if git -C "$root" worktree remove --force "$w" >/dev/null; then echo "removed $w"
      else echo "kept $t: removing $w failed"; failed=1; fi
    done
    [ -z "$failed" ] || continue
    # The tabs and worktrees are gone, so their records go too. One record stays
    # that names the branch: the rounds to come need it, and no tab in it.
    if [ -z "$dry" ]; then
      forget
      recs=("$panes/branch-$t")
      printf 'ticket=%s\nbranch=%s\n' "$t" "$branch" > "$panes/branch-$t"
    fi
  fi

  # The local branch. Gone already, or out in the main worktree or the one this
  # runs in: nothing to delete. The records go with a gone branch, and stay for
  # one that is checked out.
  if ! git -C "$root" show-ref -q --verify "refs/heads/$branch"; then
    [ -n "$dry" ] || forget
    continue
  fi
  if [ "$(git -C "$root" symbolic-ref -q --short HEAD)" = "$branch" ] \
    || [ "$(git -C "$main" symbolic-ref -q --short HEAD)" = "$branch" ]; then continue; fi

  # Merged: into origin's default branch, where a milestone pull request lands
  # it; or as its own pull request, merged at the head the branch has now. The
  # first needs no forge.
  tip="$(git -C "$root" rev-parse --verify -q "refs/heads/$branch" || true)"
  origin_head="$(git -C "$root" rev-parse --verify -q refs/remotes/origin/HEAD || true)"
  if [ -z "$origin_head" ] || ! git -C "$root" merge-base --is-ancestor "$tip" "$origin_head" 2>/dev/null; then
    if [ -z "$forge" ]; then
      echo "kept $t: cannot tell the forge from origin; set it with 'git config jaira.forge github' (or gitlab)"; continue
    fi
    if ! head="$(merged_head "$branch" 2>/dev/null)"; then
      echo "kept $t: asking $forge about $branch failed"; continue
    fi
    # Work still under way, the usual case: named, but not as something kept.
    if [ -z "$head" ]; then echo "waiting $t: $branch has no merged pull request"; continue; fi
    # A squash or rebase merge leaves the branch unmerged in git's eyes, so the
    # pull request's head is the measure: anything past it was never merged. A
    # head that was never fetched cannot be compared either, and is kept alike.
    if [ "$tip" != "$head" ] && ! git -C "$root" merge-base --is-ancestor "$tip" "$head" 2>/dev/null; then
      echo "kept $t: $branch is not at the merged head ${head:0:12}"; continue
    fi
  fi

  if [ -n "$dry" ]; then echo "would delete branch $branch"; continue; fi
  # -D, because a squash merge is no merge to git; the checks above stand in for
  # the one -d would make.
  if git -C "$root" branch -D "$branch" >/dev/null; then echo "deleted branch $branch"; forget
  else echo "kept $t: deleting branch $branch failed"; fi
done
exit 0
