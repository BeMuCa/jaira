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
alone as 'waiting'. On a board whose .jaira/pr says none there are no pull
requests: there the branch counts as merged once the default branch of the
remote work lands on contains its tip — origin, or its parent where origin is
a GitHub fork — and no forge is asked about a pull request. A merged one is kept, with the reason printed, when
  - its local branch is not at the merged head: commits past it, or the head
    never fetched,
  - one of its worktrees has changes that are not committed, as
    uncommitted.sh beside this file counts them (ticket and milestone files,
    which the board's refs hold, do not count),
  - one of its worktrees has a detached HEAD,
  - the ticket is back on the board in a lane other than done,
  - herdr does not answer, so an open tab cannot be told from a closed one.

  --dry-run   Print what would be closed and removed, change nothing.
  -h, --help  This text.

One run at a time per repository: a second one while the first still runs
prints 'skipped: reap.sh runs already' and exits 0. The lock is flock on
reap.lock beside jaira/panes/; a run that dies lets go of it.

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

# The worktree's top, not the directory: started in a subdirectory, it would
# otherwise not know the worktree it stands in.
root="$(cd "$(git -C "${1:-$PWD}" rev-parse --show-toplevel)" && pwd -P)"
here="$(cd "$(dirname "$0")" && pwd)"
herdr="${HERDR_BIN_PATH:-herdr}"
# Herdr App (herdr-core) sets no HERDR_ENV and has no herdr binary: its hc runs
# the same calls through herdr-hc.sh beside this file. HC_TERMINAL_ID is what
# the app exports into every terminal it starts.
if [ "${HERDR_ENV:-}" != 1 ] && [ -n "${HC_TERMINAL_ID:-}" ] \
  && command -v "${HC_BIN:-hc}" >/dev/null 2>&1; then
  herdr="$here/herdr-hc.sh"
fi
common="$(git -C "$root" rev-parse --git-common-dir)"
panes="$(cd "$root" && cd "$common" && pwd)/jaira/panes"
# No record, nothing to ask the forge about: the round costs nothing.
ls "$panes"/* >/dev/null 2>&1 || exit 0
# One run at a time per repository: every person's cockpit and every teamlead
# start this, and two runs ask the forge twice and race each other over
# 'git worktree remove'. The lock is flock on reap.lock, held through fd 9 for
# as long as this shell lives: the kernel lets go of it when the run dies, so
# no lock is ever left over to take over. flock(1) where there is one, perl's
# flock on macOS, which ships no flock(1); with neither the run goes unlocked,
# as before. A second run while the first lives skips the round, the next one
# comes. --dry-run changes nothing and takes no lock.
if [ -z "$dry" ]; then
  exec 9>>"$(dirname "$panes")/reap.lock"
  if command -v flock >/dev/null 2>&1; then
    flock -n 9 || { echo "skipped: reap.sh runs already"; exit 0; }
  elif command -v perl >/dev/null 2>&1; then
    rc=0; perl -MFcntl=:flock -e 'open(my $f, ">>&=", 9) or exit 2; flock($f, LOCK_EX | LOCK_NB) or exit 1' || rc=$?
    [ "$rc" != 1 ] || { echo "skipped: reap.sh runs already"; exit 0; }
  fi
fi
# Everything below runs with fd 9 closed: bash keeps the lock in a copy no
# child inherits. An inherited fd 9 would keep the lock held by whatever
# outlives this run — 'jaira' starts a detached fetch and snapshot when they
# are due — and the next run would skip with nobody reaping.
{
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
repo=""; fork=""
if [ "$forge" = github ] && origin="$(git -C "$root" remote get-url origin 2>/dev/null)" \
  && view="$(gh repo view "$origin" --json isFork,parent,nameWithOwner 2>/dev/null)"; then
  if [ "$(jaira json get isFork --default false <<<"$view")" = true ]; then
    repo="$(jaira json get parent.owner.login <<<"$view")/$(jaira json get parent.name <<<"$view")"
    fork=1
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

# How work arrives on this board ('jaira roles list'): under none there is no
# pull request to ask about, and the default branch containing the branch's
# tip is what merged means.
# When it cannot be told, nothing is reaped and every ticket says why: read as
# person instead, a none board would wait for a pull request that never comes
# and say 'no merged pull request', which names the wrong cause.
prmode=""; prwhy=""
prerr="$(mktemp)"
if roles="$( (cd "$root" && jaira roles list --json) 2>"$prerr")"; then
  prmode="$(jaira json get pr --default person <<<"$roles" 2>/dev/null)" \
    || prwhy="'jaira json get pr' cannot read what 'jaira roles list --json' printed — an older jaira? run 'jaira self upgrade'"
else
  prwhy="'jaira roles list' failed: $( (grep -v 'older version' "$prerr" || true) | tail -1)"
  [ "${prwhy%: }" != "$prwhy" ] && prwhy="${prwhy}no message (exit status only)"
fi
rm -f "$prerr"

# The remote work lands on under none, as jaira-role-pr settles it: origin,
# or — where origin is a fork — the remote that points at its parent, so a
# fork clone measures the upstream it lands on and not its own fork, which
# never receives the landing. The parent is matched by URL first, then by
# asking gh what each remote's URL resolves to: a repository renamed or moved
# since the remote was added keeps its old URL, which GitHub redirects.
# Nothing when origin is a fork and no remote here points at its parent.
land_remote() {
  local r url pat
  if [ -z "$fork" ]; then echo origin; return; fi
  pat="[:/]$(printf '%s' "$repo" | sed 's/[.]/[.]/g')(\.git)?/?\$"
  for r in $(git -C "$root" remote); do
    url="$(git -C "$root" remote get-url "$r" 2>/dev/null || true)"
    if grep -qiE "$pat" <<<"$url"; then echo "$r"; return; fi
  done
  for r in $(git -C "$root" remote); do
    [ "$r" = origin ] && continue
    url="$(git -C "$root" remote get-url "$r" 2>/dev/null || true)"
    [ -n "$url" ] || continue
    if [ "$( (gh repo view "$url" --json nameWithOwner 2>/dev/null || true) | jaira json get nameWithOwner --default "" 2>/dev/null)" \
      = "$repo" ]; then echo "$r"; return; fi
  done
}

# That remote's default branch: <remote>/HEAD, which a clone sets, or
# <remote>/main or <remote>/master where it is missing (git init + remote add)
# or points at a branch that is gone; nothing when none of them is there.
default_base() {
  local r="$1" b
  [ -n "$r" ] || return 0
  # Full ref names throughout: --short would spell origin/X as remotes/origin/X
  # where a local branch is called origin/X.
  b="$(git -C "$root" symbolic-ref -q "refs/remotes/$r/HEAD" 2>/dev/null || true)"
  # <remote>/HEAD can dangle: the remote renamed its default branch and a
  # 'fetch --prune' took the old one away. Then fall back as if it were unset.
  if [ -n "$b" ] && ! git -C "$root" rev-parse --verify -q "$b" >/dev/null 2>&1; then b=""; fi
  if [ -z "$b" ]; then
    for b in "refs/remotes/$r/main" "refs/remotes/$r/master" ""; do
      [ -n "$b" ] && git -C "$root" rev-parse --verify -q "$b" >/dev/null 2>&1 && break
    done
  fi
  printf '%s' "$b"
}

# Under none: the branch's tip once the default branch ($3) contains it and a
# commit there names the ticket ($2), or nothing — not landed yet, or a branch
# with no work of its own yet, whose tip the default branch contains from the
# start. git does the search (--grep, in C) and hands over only the commits
# with a line naming the ticket; the subject check then reads those few,
# instead of every subject of the whole history passing through the shell on
# every run. --grep also matches a body line, so the subject is checked again.
landed_head() {
  local base="$3" tip pat
  tip="$(git -C "$root" rev-parse --verify -q "refs/heads/$1" 2>/dev/null || true)"
  pat="^((Revert|Reapply) \")*[a-z]+\\(([^)]*[, ])?$2([, ][^)]*)?\\)!?:"
  if [ -n "$base" ] && [ -n "$tip" ] && git -C "$root" merge-base --is-ancestor "$tip" "$base" 2>/dev/null \
    && grep -qE "$pat" <<<"$(git -C "$root" log -E --grep="$pat" --format=%s "$tip" --)"; then
    echo "$tip"
  fi
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

  if [ -n "$prwhy" ]; then
    echo "kept $t: cannot tell how work arrives on this board — $prwhy"; continue
  fi
  if [ "$prmode" = none ]; then
    lr="$(land_remote)"
    if [ -z "$lr" ]; then
      echo "waiting $t: no remote here points at $repo, where work lands; add one ('git remote add upstream <its URL>')"; continue
    fi
    base="$(default_base "$lr")"
    if [ -z "$base" ]; then
      echo "waiting $t: cannot tell $lr's default branch; run 'git remote set-head $lr --auto'"; continue
    fi
    head="$(landed_head "$branch" "$t" "$base")"
    if [ -z "$head" ]; then echo "waiting $t: $branch has not landed on $lr's default branch"; continue; fi
  else
    if [ -z "$forge" ]; then
      echo "kept $t: cannot tell the forge from origin; set it with 'git config jaira.forge github' (or gitlab)"; continue
    fi
    if ! head="$(merged_head "$branch" 2>/dev/null)"; then
      echo "kept $t: asking $forge about $branch failed"; continue
    fi
    # Work still under way, the usual case: named, but not as something kept.
    if [ -z "$head" ]; then echo "waiting $t: $branch has no merged pull request"; continue; fi
  fi

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
    # A detached HEAD — a rebase paused, a commit checked out — may carry
    # commits no branch holds, and removing the worktree would lose them.
    if ! git -C "$w" symbolic-ref -q HEAD >/dev/null; then dirty="$w has a detached HEAD"; break; fi
    # What no commit and no ref holds: the rule close-tab.sh asks as well.
    # --force below would delete it.
    if ! d="$("$here/uncommitted.sh" "$w")"; then dirty="git status failed in $w"; break; fi
    if [ -n "$d" ]; then dirty="$w has uncommitted changes: $d"; break; fi
  done
  if [ -n "$dirty" ]; then echo "kept $t: ${dirty% }"; continue; fi
  if [ -z "$herdr_up" ]; then echo "kept $t: herdr not reachable ($herdr)"; continue; fi

  failed=""; gone=()
  for f in "${recs[@]}"; do
    pane="$(field pane "$f")"; lane="$(field lane "$f")"
    if [ "$pane" = "${HERDR_PANE_ID:-${HC_TERMINAL_ID:-}}" ]; then continue; fi
    gone+=("$f")
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
    # --force for the .jaira/ files the check above lets through.
    if git -C "$root" worktree remove --force "$w" >/dev/null; then echo "removed $w"
    else echo "kept $t: removing $w failed"; failed=1; fi
  done
  # A branch checked out in a worktree that stays — the main one, or the one
  # this runs in — stays with it.
  if [ -z "$failed" ] && git -C "$root" show-ref -q --verify "refs/heads/$branch" \
    && [ "$(git -C "$root" symbolic-ref -q --short HEAD)" != "$branch" ] \
    && [ "$(git -C "$main" symbolic-ref -q --short HEAD)" != "$branch" ]; then
    if [ -n "$dry" ]; then echo "would delete branch $branch"
    # -D, because a squash merge is no merge to git; the head check above
    # stands in for the one -d would make.
    elif git -C "$root" branch -D "$branch" >/dev/null; then echo "deleted branch $branch"
    else echo "kept $t: deleting branch $branch failed"; failed=1; fi
  fi
  # The records go last: one left behind brings the next round back to what
  # failed, where a record gone first left the worktree standing for good.
  [ -n "$dry$failed" ] || rm -f "${gone[@]+"${gone[@]}"}"
done
} 9>&-
exit 0
