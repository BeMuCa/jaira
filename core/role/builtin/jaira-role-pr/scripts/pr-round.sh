#!/usr/bin/env bash
# One round of review comments on a GitHub pull request: list what is open,
# answer and resolve one thread, answer a comment that has no thread, and ask
# for the review again. Never merges, never approves.
# Usage: pr-round.sh <command> <owner/repo> <number> [args]
set -euo pipefail

usage() {
  cat <<'USAGE'
pr-round.sh <command> <owner/repo> <number> [args]
pr-round.sh open <owner/repo> <remote> <branch> <title> <body-file>

  open      <owner/repo> <remote> <branch> <title> <body-file>
            Open the pull request of the checked-out branch into <branch> of
            <owner/repo> — the agent's way in, on a board whose .jaira/pr says
            agent. The checkout must be the pushed tip of its branch, and
            'jaira roles pr-gate <remote>/<branch>' runs first: unless every
            ticket the branch's commits name is in done — or in signoff,
            where .jaira/pr-when says so — nothing is opened and the script
            exits 3 with the reason.
  threads   <owner/repo> <number>
            Every review thread that is not resolved yet, one JSON object per
            line: id, path, line, outdated, and its comments (author, body,
            url). A thread resolved in an earlier round is not listed, so a
            second round never answers the same point twice.
  comments  <owner/repo> <number>
            Comments that have no thread to resolve — the conversation and
            the bodies of submitted reviews — one JSON object per line:
            kind, author, body, url. Answer them with 'comment'; they cannot
            be resolved.
  reply     <owner/repo> <number> <thread-id> <body-file>
            Answer one thread with the text of <body-file> ('Fixed in <sha>:
            ...' or the reason it stays), then resolve it. Two sessions
            answering the same thread at once leave one answer: the later
            one removes its own copy.
  comment   <owner/repo> <number> <body-file>
            Answer on the conversation, for a comment that has no thread.
            Where the conversation already ends in these words under this
            login, nothing is posted, so a second round never says it twice.
  request   <owner/repo> <number> <reviewer>
            Ask <reviewer> for a review again. 'copilot' asks GitHub Copilot.
  settled   <owner/repo> <number>
            Whether the review round is over: exit 0 when no thread is open,
            no review is still requested and no reviewer's latest review
            requests changes; exit 3 listing each of them. On a board whose .jaira/pr-when
            says signoff, done waits for exit 0.

Exit codes: 0 done; 2 usage; 3 gh missing or not signed in to github.com
(also when the sign-in runs out half way), (open) the checkout is not the
pushed tip or the board does not let an agent open this pull request, or
(settled) the review round is not over;
4 the reviewer could not be requested (the pull request stays as it is — say
so and carry on); 1 anything else that failed, with gh's own message.

GitHub only. On GitLab the role's text names the glab commands instead.
USAGE
}

die() { echo "pr-round.sh: $1" >&2; exit "${2:-1}"; }

[ "${1:-}" = -h ] || [ "${1:-}" = --help ] && { usage; exit 0; }
[ $# -ge 3 ] || { usage >&2; exit 2; }

cmd=$1 repo=$2
case "$repo" in */*) ;; *) die "'$repo' is not owner/repo" 2 ;; esac
if [ "$cmd" = open ]; then
  # open has no number yet: it is the call that makes one.
  [ $# -eq 6 ] || { usage >&2; exit 2; }
  shift 2
else
  number=$3
  case "$number" in ''|*[!0-9]*) die "'$number' is not a pull request number" 2 ;; esac
  shift 3
fi
owner=${repo%%/*} name=${repo#*/}

command -v gh >/dev/null 2>&1 || die "gh is not installed" 3
# Checked once, before anything is written: a token that expires half way
# through would leave some threads answered and others not, and the report
# has to say which. Each write below still fails loudly on its own.
gh auth status --hostname github.com >/dev/null 2>&1 ||
  die "gh is not signed in to github.com — run: gh auth login" 3
# gh ends with 4 when it needs a sign-in, and 4 is this script's 'reviewer
# not requested'. A token that runs out after the check above is a sign-in
# problem, so every gh call below reports it as one: exit 3.
gh() {
  local rc=0
  command gh "$@" || rc=$?
  [ "$rc" -ne 4 ] || die "gh lost its sign-in to github.com half way — run: gh auth login, then run this again" 3
  return "$rc"
}

# Comments of one thread beyond the first 50 are not needed to answer it.
threads_query='query($owner:String!,$name:String!,$number:Int!,$endCursor:String){
  repository(owner:$owner,name:$name){pullRequest(number:$number){
    reviewThreads(first:100,after:$endCursor){
      pageInfo{hasNextPage endCursor}
      nodes{id isResolved isOutdated path line
        comments(first:50){nodes{author{login} body url}}}}}}}'

# said is a text as this script compares it: CR dropped, blanks at the end
# of every line and newlines at the end of the text dropped, as base64 — so
# a forge that trims what it stores still matches, and no reviewer's text
# reaches a shell or a jq program. saidjq is the same in jq, for the bodies
# gh returns.
said() { printf '%s' "$(tr -d '\r' < "$1" | sed 's/[[:space:]]*$//')" | base64 | tr -d '\n'; }
saidjq='def said: gsub("\r"; "") | gsub("[ \t]+\n"; "\n") | sub("[ \t\n]+$"; "") | @base64;'

# keeper reads "<who> <id> <said>" lines, oldest first, and prints the
# first id of the unbroken run of <who> = $1 saying $2 that reaches the
# comment $3 — or the end of the list, when $3 is empty. Any other comment,
# another author's or other words, ends a run: an answer repeated after a
# reviewer objected is a new answer, not a copy. It reads to the end, never
# breaking out: under pipefail a gh cut off by a closed pipe would end the
# script.
keeper() {
  local who=$1 want=$2 stop=$3 run="" out="" found="" w id b
  while read -r w id b; do
    [ -n "$found" ] && continue
    if [ "$w" = "$who" ] && [ "$b" = "$want" ]; then
      [ -n "$run" ] || run=$id
    else
      run=""
    fi
    if [ -n "$stop" ] && [ "$id" = "$stop" ]; then out=$run; found=1; fi
  done
  if [ -z "$stop" ]; then out=$run; fi
  echo "$out"
}

# dedupe keeps one answer where two sessions answered the same thread at
# once. The check before posting cannot see an answer still in flight, so
# after posting the thread is read again: the run of this login's comments
# with this text that ends in the one just posted keeps its first, and if
# that is not the one just posted, that one is removed. Every order of two
# sessions ends with exactly one — whichever reads after both posted sees
# the other first, whichever reads before sees only its own. Only the id
# this run posted is ever removed: a person's comment never is.
dedupe() {
  local mine_id=$1 first
  [ -n "$mine_id" ] || return 0
  # The newest 100, oldest first: they hold the answer just posted and any
  # copy beside it however long the thread has grown.
  local rc=0
  first=$(gh api graphql -f thread="$thread" -f query='
    query($thread:ID!){viewer{login}
      node(id:$thread){... on PullRequestReviewThread{comments(last:100){nodes{id author{login} body}}}}}' \
    --jq "$saidjq"' .data as $d | $d.node.comments.nodes[]
          | "\(if (.author.login // "") == $d.viewer.login then "me" else "other" end) \(.id) \(.body | said)"' |
    keeper me "$(said "$file")" "$mine_id") || rc=$?
  # 3 is gh() reporting a lost sign-in, already said: it ends the run, as
  # every command promises. Any other failure only leaves the check undone.
  [ "$rc" -ne 3 ] || exit 3
  if [ "$rc" -ne 0 ]; then
    echo "pr-round.sh: could not read the thread back to check for a second answer — look at it once" >&2
    return 0
  fi
  if [ -n "$first" ] && [ "$first" != "$mine_id" ]; then
    gh api graphql -f id="$mine_id" -f query='
      mutation($id:ID!){deletePullRequestReviewComment(input:{id:$id}){clientMutationId}}' >/dev/null
    echo "answered: by another session at the same time, removed this copy"
  fi
}

body_of() {
  [ -f "$1" ] || die "no file '$1' to read the answer from" 2
  # Blanks alone are no answer either: GitHub refuses them, or posts nothing.
  [ -n "$(tr -d '[:space:]' < "$1")" ] || die "'$1' is empty — an answer needs words" 2
}

case "$cmd" in
  open)
    remote=$1 branch=$2 title=$3 file=$4
    body_of "$file"
    command -v jaira >/dev/null 2>&1 || die "jaira is not on PATH — the gate cannot run, so nothing is opened" 3
    # The gate reads the commits of the checkout; gh opens the branch as it
    # was pushed. Another session may have pushed more since, so the two must
    # be the same commit, or the gate judges commits the pull request lacks
    # or misses ones it carries.
    cur=$(git symbolic-ref --quiet --short HEAD) ||
      die "HEAD is detached — check out the branch to open, so the gate and the pull request read the same commits" 3
    # Read from the configuration, not from a remote-tracking ref: a
    # single-branch clone or a push-only remote has none, and the branch is
    # pushed all the same. 'git push' sends it under its own name, or under
    # the branch it tracks where push.default says upstream (or tracking,
    # its older name).
    pushremote=$(git for-each-ref --format='%(push:remotename)' "refs/heads/$cur")
    [ -n "$pushremote" ] ||
      die "'$cur' has no remote to push to — push it first: git push -u <remote> HEAD" 3
    pushed=refs/heads/$cur
    case "$(git config push.default || true)" in
      upstream|tracking)
        pushed=$(git config "branch.$cur.merge") || die "'$cur' tracks no branch to push to — push it first: git push -u $pushremote HEAD" 3 ;;
    esac
    # Read where git push sends it: a remote's pushurl, where it has one,
    # is not the URL it fetches from.
    pushurl=$(git remote get-url --push "$pushremote") ||
      die "$pushremote has no URL to push to — nothing is opened" 3
    tip=$(git ls-remote "$pushurl" "$pushed" | cut -f1) ||
      die "could not read '$pushed' from $pushremote — nothing is opened" 3
    here=$(git rev-parse HEAD)
    [ -n "$tip" ] || die "'$cur' is not pushed to $pushremote yet — push it first: git push -u $pushremote HEAD" 3
    [ "$tip" = "$here" ] ||
      die "the pushed ${pushed#refs/heads/} on $pushremote is at ${tip:0:12}, this checkout at ${here:0:12} — pull or push until they agree, then open it again" 3
    # --head names the branch the tip check above read, on the account it
    # was pushed to: left to itself gh guesses the head from the checkout,
    # and in a fork that guess need not be the commits the gate judged. The
    # URL is the same entry get-url --push read — the first pushurl, else
    # the url — as configured, before insteadOf rewrites it. Only
    # <scheme>://host/<owner>/<name> and host:<owner>/<name> name an
    # account; anything else opens nothing rather than a guess,
    # and is refused before the gate runs.
    url=$(git config --get-all "remote.$pushremote.pushurl" | head -n1 || true)
    [ -n "$url" ] || url=$(git config "remote.$pushremote.url") ||
      die "$pushremote has no URL — nothing is opened" 3
    path=${url%/}; path=${path%.git}
    case "$path" in
      *://*) path=${path#*://}; case "$path" in */*) path=${path#*/} ;; *) path= ;; esac ;;
      *:*) path=${path#*:} ;;
      *) path= ;;
    esac
    path=${path#/}
    case "$path" in
      ?*/?*/*|/*|*//*) headowner= ;;
      ?*/?*) headowner=${path%/*} ;;
      *) headowner= ;;
    esac
    [ -n "$headowner" ] ||
      die "cannot tell which account $pushremote ($url) belongs to — nothing is opened" 3
    # The gate is the binary's, not this script's: whether the board says
    # agent, and whether every ticket the branch names is far enough along.
    if ! gate=$(jaira roles pr-gate "$remote/$branch" 2>&1); then
      echo "$gate" >&2
      die "not opened — hand the create line back to a person" 3
    fi
    echo "$gate"
    head=${pushed#refs/heads/}
    # Into its own repository the branch name alone is the head.
    [ "$(printf '%s' "$headowner" | tr '[:upper:]' '[:lower:]')" = "$(printf '%s' "$owner" | tr '[:upper:]' '[:lower:]')" ] ||
      head=$headowner:$head
    gh pr create --repo "$repo" --base "$branch" --head "$head" --title "$title" --body-file "$file"
    ;;
  threads)
    # -f for every string: -F turns a value of digits into a number, and a
    # repository named 2048 would then fail the String! it is declared as.
    gh api graphql --paginate -f owner="$owner" -f name="$name" -F number="$number" \
      -f query="$threads_query" \
      --jq '.data.repository.pullRequest.reviewThreads.nodes[] | select(.isResolved | not)
            | {id, path, line, outdated: .isOutdated,
               comments: [.comments.nodes[] | {author: (.author.login // "ghost"), body, url}]}
            | tojson'
    ;;
  comments)
    gh api --paginate "repos/$repo/issues/$number/comments" \
      --jq '.[] | {kind: "conversation", author: (.user.login // "ghost"), body, url: .html_url} | tojson'
    gh api --paginate "repos/$repo/pulls/$number/reviews" \
      --jq '.[] | select(.body != "") | {kind: "review", author: (.user.login // "ghost"), body, url: .html_url} | tojson'
    ;;
  reply)
    [ $# -eq 2 ] || { usage >&2; exit 2; }
    thread=$1 file=$2
    body_of "$file"
    # Answer and resolve are two calls, and the second can fail after the
    # first went through. Run again, the thread is still listed open; its last
    # comment is then this very answer, so only the resolve is repeated and
    # the reviewer never gets the same answer twice. The text decides, not the
    # author alone: an agent writes under the maintainer's own login, and the
    # maintainer's own remark in a thread is a question still owed an answer.
    resp=$(gh api graphql -f thread="$thread" -f query='
      query($thread:ID!){viewer{login}
        node(id:$thread){... on PullRequestReviewThread{isResolved comments(last:1){nodes{author{login} body}}}}}' \
      --jq "$saidjq"' if .data.node == null then "missing"
            elif .data.node.isResolved then "resolved"
            elif (.data.node.comments.nodes[0].author.login // "") == .data.viewer.login
              then "mine \(.data.node.comments.nodes[0].body // "" | said)"
            else "open" end')
    state=${resp%% *}
    if [ "$state" = mine ]; then
      [ "${resp#* }" = "$(said "$file")" ] && state=answered || state=open
    fi
    case "$state" in
      missing) die "no review thread '$thread' on $repo#$number" 1 ;;
      resolved) echo "resolved: already"; exit 0 ;;
      answered) echo "answered: already, resolving" ;;
      *)
        posted=$(gh api graphql -f thread="$thread" -F body=@"$file" -f query='
          mutation($thread:ID!,$body:String!){
            addPullRequestReviewThreadReply(input:{pullRequestReviewThreadId:$thread,body:$body}){comment{id url}}}' \
          --jq '.data.addPullRequestReviewThreadReply.comment | "\(.id) \(.url)"')
        echo "${posted#* }"
        dedupe "${posted%% *}"
        ;;
    esac
    gh api graphql -f thread="$thread" -f query='
      mutation($thread:ID!){resolveReviewThread(input:{threadId:$thread}){thread{isResolved}}}' \
      --jq '"resolved: \(.data.resolveReviewThread.thread.isResolved)"'
    ;;
  comment)
    [ $# -eq 1 ] || { usage >&2; exit 2; }
    file=$1
    body_of "$file"
    # The conversation has nothing to resolve, so what keeps a second round
    # from saying the same thing twice is the text: where the conversation
    # ends in these words under this login, nothing is posted. Anything said
    # after them — a remark, an objection, the maintainer's own words — means
    # the answer is owed again. After posting, the conversation is read back
    # the way dedupe reads a thread, for two sessions posting at once.
    # viewer, not /user: an app or Actions token is refused on /user, and
    # reply already asks GraphQL the same way.
    viewer=$(gh api graphql -f query='query{viewer{login}}' --jq .data.viewer.login)
    mine=$(said "$file")
    convo() {
      gh api --paginate "repos/$repo/issues/$number/comments" \
        --jq "$saidjq"' .[] | "\(.user.login // "-") \(.id) \(.body | said)"'
    }
    first=$(convo | keeper "$viewer" "$mine" "")
    [ -n "$first" ] && { echo "commented: already"; exit 0; }
    posted=$(gh api --method POST "repos/$repo/issues/$number/comments" -F body=@"$file" --jq '"\(.id) \(.html_url)"')
    echo "${posted#* }"
    rc=0
    first=$(convo | keeper "$viewer" "$mine" "${posted%% *}") || rc=$?
    [ "$rc" -ne 3 ] || exit 3
    if [ "$rc" -ne 0 ]; then
      echo "pr-round.sh: could not read the conversation back to check for a second answer — look at it once" >&2
      first=""
    fi
    if [ -n "$first" ] && [ "$first" != "${posted%% *}" ]; then
      gh api --method DELETE "repos/$repo/issues/comments/${posted%% *}" >/dev/null
      echo "commented: by another session at the same time, removed this copy"
    fi
    ;;
  settled)
    [ $# -eq 0 ] || { usage >&2; exit 2; }
    # Both lists in full before deciding: a round is over only when neither
    # has anything, and the report names everything still in the way.
    open=$(gh api graphql --paginate -f owner="$owner" -f name="$name" -F number="$number" \
      -f query="$threads_query" \
      --jq '.data.repository.pullRequest.reviewThreads.nodes[] | select(.isResolved | not)
            | "open thread: \(.path):\(.line // "-") \(.comments.nodes[0].url // "")"')
    pending=$(gh api "repos/$repo/pulls/$number/requested_reviewers" \
      --jq '(.users // [])[] , (.teams // [])[] | "review pending: \(.login // .slug)"')
    # A review that requested changes may say so only in its body, with no
    # thread to stay open: each reviewer's latest verdict is read as well.
    changes=$(gh api graphql -f owner="$owner" -f name="$name" -F number="$number" -f query='
      query($owner:String!,$name:String!,$number:Int!){repository(owner:$owner,name:$name){
        pullRequest(number:$number){latestOpinionatedReviews(first:100){nodes{state author{login}}}}}}' \
      --jq '.data.repository.pullRequest.latestOpinionatedReviews.nodes[]
            | select(.state == "CHANGES_REQUESTED") | "changes requested: \(.author.login // "ghost")"')
    if [ -z "$open$pending$changes" ]; then
      echo "settled: no thread open, no review requested"
      exit 0
    fi
    printf '%s\n' "$open" "$pending" "$changes" | sed '/^$/d' >&2
    exit 3
    ;;
  request)
    [ $# -eq 1 ] || { usage >&2; exit 2; }
    # A login is written with or without its @; the API takes it without.
    reviewer=${1#@}
    # gh pr edit --add-reviewer takes no special value for Copilot; the REST
    # endpoint takes its bot login.
    [ "$reviewer" = copilot ] && reviewer='copilot-pull-request-reviewer[bot]'
    case "$reviewer" in
      */*) field="team_reviewers[]=${reviewer#*/}" ;;
      *) field="reviewers[]=$reviewer" ;;
    esac
    out=$(gh api --method POST "repos/$repo/pulls/$number/requested_reviewers" -f "$field" 2>&1) || {
      rc=$?
      [ "$rc" -eq 3 ] && { echo "$out" >&2; exit 3; }
      echo "pr-round.sh: could not request a review from $1 — $out" >&2
      exit 4
    }
    echo "requested: $1"
    ;;
  *)
    usage >&2
    exit 2
    ;;
esac
