#!/usr/bin/env bash
# Land an accepted pull request: bring the base into the ticket's branch,
# regenerate what the board names, push, wait for CI, then merge — only where
# the board lets an agent merge. Otherwise hand the merge line back.
# Usage: land.sh <pr-number>
set -euo pipefail

usage() {
  cat <<'USAGE'
land.sh <pr-number>

Run in the ticket's worktree, on the branch the pull request comes from.
GitHub only: gh does the asking, the waiting and the merging.

  0. Stops (exit 2) when origin's copy of the branch has commits this
     worktree lacks: pull them first.
  1. Merges the pull request's base into the branch. A conflict in a file
     .jaira/land-regenerate names is taken from the base and regenerated; any
     other conflict aborts the merge and stops (exit 1), nothing pushed.
  2. Runs every command .jaira/land-regenerate names, merge or not, stages
     its file and commits the merge and the regenerated files as one commit
     naming the ticket: chore(<ID>): merge <base>, regenerate ... — or
     chore(<ID>): regenerate ... where the branch held the base already.
  3. Pushes to origin and waits for the pull request's checks (gh pr checks
     --watch). A red check stops (exit 1).
  4. Asks 'jaira roles pr-gate --merge <base>': .jaira/land in the base says
     agent and every ticket on the branch is in done. Yes: gh pr merge
     --merge, then reap.sh closes the ticket's tabs and worktree. No: prints
     the reason and the merge line for a person, merges nothing (exit 3).
     gh refusing the merge prints the merge line too (exit 1).

.jaira/land-regenerate — read as the base holds it, like .jaira/land — has
one generated file per line, '<path>: <command>'; the command runs from the
repository root and rewrites <path>. Blank lines and lines starting with #
are skipped. Example:
  core/shipped/shipped.sha256: scripts/shipped-history.sh --write

JAIRA_LAND_SETTLE seconds (default 30) pass between a push and the first look
at the checks, so the checks of the new head have registered.

Exit 0 merged, 1 stopped on a conflict, a failed command or red CI, 2 usage
or a precondition, 3 the board leaves the merge to a person.
USAGE
}

case "${1:-}" in
  -h|--help) usage; exit 0 ;;
  ''|*[!0-9]*) usage >&2; exit 2 ;;
esac
pr="$1"
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

die() { echo "land: $2" >&2; exit "$1"; }

root="$(git rev-parse --show-toplevel 2>/dev/null)" || die 2 "not in a git repository"
cd "$root"
branch="$(git symbolic-ref -q --short HEAD)" || die 2 "detached HEAD — check out the pull request's branch"
id="$(printf '%s' "${branch##*/}" | tr '[:lower:]' '[:upper:]')"
[[ "$id" =~ ^[0-9A-Z]{6}$ ]] || die 2 "branch $branch names no ticket (want feat/<id>)"
if [ -n "$(git status --porcelain --untracked-files=no)" ]; then
  die 2 "uncommitted changes in $root — commit them first: $(git status --porcelain --untracked-files=no | head -3 | tr '\n' ' ')"
fi

# Where the pull request lives: origin's repository, or its parent when
# origin is a fork — named, never left to gh, which asks the fork otherwise.
origin="$(git remote get-url origin 2>/dev/null)" || die 2 "no origin remote"
case "$origin" in *github.com*) ;; *) die 2 "origin is not on GitHub — land.sh speaks gh only; merge by hand" ;; esac
view="$(gh repo view "$origin" --json isFork,parent,nameWithOwner)" || die 2 "gh repo view $origin did not answer"
if [ "$(jaira json get isFork --default false <<<"$view")" = true ]; then
  repo="$(jaira json get parent.owner.login <<<"$view")/$(jaira json get parent.name <<<"$view")"
else
  repo="$(jaira json get nameWithOwner <<<"$view")"
fi
mergeline="gh pr merge $pr --repo $repo --merge"

prv="$(gh pr view "$pr" --repo "$repo" --json state,headRefName,baseRefName)" || die 2 "gh pr view $pr --repo $repo did not answer"
[ "$(jaira json get state <<<"$prv")" = OPEN ] || die 2 "pull request #$pr is not open"
[ "$(jaira json get headRefName <<<"$prv")" = "$branch" ] || die 2 "pull request #$pr comes from $(jaira json get headRefName <<<"$prv"), not $branch"
basebranch="$(jaira json get baseRefName <<<"$prv")"

# The local remote that points at $repo: origin when it is no fork, else the
# remote whose URL names the parent.
remote=""
pat="[:/]$(printf '%s' "$repo" | sed 's/[.]/[.]/g')(\.git)?/?\$"
for r in $(git remote); do
  if grep -qiE "$pat" <<<"$(git remote get-url "$r" 2>/dev/null || true)"; then remote="$r"; break; fi
done
[ -n "$remote" ] || die 2 "no remote here points at $repo — git remote add upstream <its url>"
# origin's copy of the branch first: a push somebody made there would make
# ours fail non-fast-forward after the merge, with nothing said about why.
if git fetch -q origin "$branch" 2>/dev/null &&
  ! git merge-base --is-ancestor "refs/remotes/origin/$branch" HEAD; then
  die 2 "origin's $branch has commits this worktree lacks — git pull origin $branch, then run again; nothing merged or pushed"
fi
git fetch -q "$remote" "$basebranch"
base="$remote/$basebranch"
git rev-parse -q --verify "refs/remotes/$base" >/dev/null || die 2 "$base was not fetched"

# Generated files, as the base names them: "<path>: <command>".
genpath=(); gencmd=()
while IFS= read -r line || [ -n "$line" ]; do
  line="${line%$'\r'}"
  case "$line" in ''|'#'*) continue ;; esac
  case "$line" in *': '*) ;; *) die 2 ".jaira/land-regenerate in $base: '$line' is not '<path>: <command>'" ;; esac
  genpath+=("${line%%: *}"); gencmd+=("${line#*: }")
done < <(git show "$base:.jaira/land-regenerate" 2>/dev/null || true)

generated() {
  local p
  for p in "${genpath[@]+"${genpath[@]}"}"; do [ "$p" = "$1" ] && return 0; done
  return 1
}

before="$(git rev-parse HEAD)"
if ! git merge --no-commit --no-ff -q "$base" >/dev/null 2>&1; then
  git rev-parse -q --verify MERGE_HEAD >/dev/null || die 1 "git merge $base failed before any conflict — run it by hand to see why"
  others=()
  while IFS= read -r f; do
    # The base's side, to regenerate over — unless the base deleted or
    # renamed the file: then there is no side to take, and it is a conflict
    # like any other.
    if generated "$f" && git checkout -q --theirs -- "$f" 2>/dev/null; then :; else others+=("$f"); fi
  done < <(git diff --name-only --diff-filter=U)
  if [ "${#others[@]}" -gt 0 ]; then
    git merge --abort
    die 1 "merging $base conflicts in ${others[*]} — resolve by hand; nothing pushed"
  fi
fi
# Guarded: bash 3.2 (macOS) calls an empty array under set -u unbound.
[ "${#genpath[@]}" -gt 0 ] && for i in "${!genpath[@]}"; do
  if ! bash -c "${gencmd[$i]}" >/dev/null; then
    git rev-parse -q --verify MERGE_HEAD >/dev/null && git merge --abort
    git reset -q --hard "$before"
    die 1 "regenerating ${genpath[$i]} failed: ${gencmd[$i]} — nothing pushed"
  fi
  git add -- "${genpath[$i]}"
done
# One commit naming the ticket: the merge and what was regenerated, or the
# regeneration alone where the branch held the base already.
if git rev-parse -q --verify MERGE_HEAD >/dev/null; then
  msg="chore($id): merge $base"
  [ "${#genpath[@]}" -gt 0 ] && msg="$msg, regenerate ${genpath[*]}"
  git commit -q --no-edit -m "$msg"
elif ! git diff --cached --quiet; then
  git commit -q -m "chore($id): regenerate ${genpath[*]}"
fi

if [ "$(git rev-parse HEAD)" != "$(git rev-parse -q --verify "refs/remotes/origin/$branch" || true)" ]; then
  if ! git push -q origin "HEAD:$branch"; then
    git reset -q --hard "$before"
    die 1 "push to origin's $branch refused — somebody pushed in between; git pull origin $branch, then run again; the local merge is undone"
  fi
  sleep "${JAIRA_LAND_SETTLE:-30}"
fi
if ! gh pr checks "$pr" --repo "$repo" --watch --interval 30 >/dev/null; then
  die 1 "#$pr: CI is not green — $(gh pr checks "$pr" --repo "$repo" 2>/dev/null | awk -F'\t' '$2!="pass"&&$2!="skipping"{printf "%s=%s ", $1, $2}')"
fi

set +e
why="$(jaira roles pr-gate --merge "$base" 2>&1)"
gate=$?
set -e
case "$gate" in
  0) ;;
  3) echo "land: #$pr is up to date with $base and green, but a person merges it — ${why#jaira: }"
     echo "$mergeline"
     exit 3 ;;
  *) die 2 "pr-gate --merge did not answer: $why" ;;
esac
# Merge exactly the head the gate and CI judged: a push to the branch while
# this waited makes gh refuse instead of merging commits nobody checked.
if ! gh pr merge "$pr" --repo "$repo" --merge --match-head-commit "$(git rev-parse HEAD)" >/dev/null; then
  echo "land: gh refused to merge #$pr (gh's reason above: the head moved, or the repository takes no merge commits) — a person merges it"
  echo "$mergeline"
  exit 1
fi
echo "land: $id #$pr merged into $base"
main="$(dirname "$(git rev-parse --path-format=absolute --git-common-dir)")"
(cd "$main" && "$here/../../jaira-dispatcher/scripts/reap.sh" "$main" 2>&1 | grep -v '^waiting' || true)
