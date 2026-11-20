#!/usr/bin/env bash
# One round of review comments on a GitHub pull request: list what is open,
# answer and resolve one thread, answer a comment that has no thread, and ask
# for the review again. Never merges, never approves.
# Usage: pr-round.sh <command> <owner/repo> <number> [args]
set -euo pipefail

usage() {
  cat <<'USAGE'
pr-round.sh <command> <owner/repo> <number> [args]

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
            ...' or the reason it stays), then resolve it.
  comment   <owner/repo> <number> <body-file>
            Answer on the conversation, for a comment that has no thread.
  request   <owner/repo> <number> <reviewer>
            Ask <reviewer> for a review again. 'copilot' asks GitHub Copilot.

Exit codes: 0 done; 2 usage; 3 gh missing or not signed in to github.com;
4 the reviewer could not be requested (the pull request stays as it is — say
so and carry on); 1 anything else that failed, with gh's own message.

GitHub only. On GitLab the role's text names the glab commands instead.
USAGE
}

die() { echo "pr-round.sh: $*" >&2; exit "${2:-1}"; }

[ "${1:-}" = -h ] || [ "${1:-}" = --help ] && { usage; exit 0; }
[ $# -ge 3 ] || { usage >&2; exit 2; }

cmd=$1 repo=$2 number=$3
shift 3
case "$repo" in */*) ;; *) die "'$repo' is not owner/repo" 2 ;; esac
case "$number" in ''|*[!0-9]*) die "'$number' is not a pull request number" 2 ;; esac
owner=${repo%%/*} name=${repo#*/}

command -v gh >/dev/null 2>&1 || die "gh is not installed" 3
# Checked once, before anything is written: a token that expires half way
# through would leave some threads answered and others not, and the report
# has to say which. Each write below still fails loudly on its own.
gh auth status --hostname github.com >/dev/null 2>&1 ||
  die "gh is not signed in to github.com — run: gh auth login" 3

# Comments of one thread beyond the first 50 are not needed to answer it.
threads_query='query($owner:String!,$name:String!,$number:Int!,$endCursor:String){
  repository(owner:$owner,name:$name){pullRequest(number:$number){
    reviewThreads(first:100,after:$endCursor){
      pageInfo{hasNextPage endCursor}
      nodes{id isResolved isOutdated path line
        comments(first:50){nodes{author{login} body url}}}}}}}'

body_of() {
  [ -f "$1" ] || die "no file '$1' to read the answer from" 2
  [ -s "$1" ] || die "'$1' is empty — an answer needs words" 2
}

case "$cmd" in
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
      --jq 'if .data.node == null then "missing"
            elif .data.node.isResolved then "resolved"
            elif (.data.node.comments.nodes[0].author.login // "") == .data.viewer.login
              then "mine\n" + (.data.node.comments.nodes[0].body // "")
            else "open" end')
    state=${resp%%$'\n'*}
    if [ "$state" = mine ]; then
      last=$(printf '%s' "${resp#*$'\n'}" | tr -d '\r')
      mine=$(tr -d '\r' < "$file")
      [ "$last" = "$mine" ] && state=answered || state=open
    fi
    case "$state" in
      missing) die "no review thread '$thread' on $repo#$number" 1 ;;
      resolved) echo "resolved: already"; exit 0 ;;
      answered) echo "answered: already, resolving" ;;
      *)
        gh api graphql -f thread="$thread" -F body=@"$file" -f query='
          mutation($thread:ID!,$body:String!){
            addPullRequestReviewThreadReply(input:{pullRequestReviewThreadId:$thread,body:$body}){comment{url}}}' \
          --jq '.data.addPullRequestReviewThreadReply.comment.url'
        ;;
    esac
    gh api graphql -f thread="$thread" -f query='
      mutation($thread:ID!){resolveReviewThread(input:{threadId:$thread}){thread{isResolved}}}' \
      --jq '"resolved: \(.data.resolveReviewThread.thread.isResolved)"'
    ;;
  comment)
    [ $# -eq 1 ] || { usage >&2; exit 2; }
    body_of "$1"
    gh pr comment "$number" --repo "$repo" --body-file "$1"
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
    if ! out=$(gh api --method POST "repos/$repo/pulls/$number/requested_reviewers" -f "$field" 2>&1); then
      echo "pr-round.sh: could not request a review from $1 — $out" >&2
      exit 4
    fi
    echo "requested: $1"
    ;;
  *)
    usage >&2
    exit 2
    ;;
esac
