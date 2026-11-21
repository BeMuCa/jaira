#!/usr/bin/env bash
# Look at every commit a push would publish for a string that must not leave
# this machine, before it does: the denylist's patterns and this machine's own
# name. Prints where a hit is — commit, file, line — and never the text that
# matched or the pattern, so the report can go anywhere the branch goes.
# Usage: pre-push-scan.sh <base>
set -euo pipefail

usage() {
  cat <<'USAGE'
pre-push-scan.sh <base>

Checks every commit of <base>..HEAD — each one, not only the sum, because a
push publishes the whole history — in five places: the lines it adds, the
names of the files it changes, its message, its author email and its
committer email. A merge is checked for what it adds itself, a conflict
resolution, and not for what it brings in from its parents.

Two lists are matched, case-insensitively:

  the denylist   one extended regular expression per line, '#' comments and
                 blank lines ignored, read from $JAIRA_MENTION_DENYLIST or
                 else ~/.config/jaira/mention-denylist. It lives outside the
                 repository on purpose: a list inside it would be the mention.
                 A hit stops the push. No file, or one with no pattern in it:
                 the scan says the strings were not checked and goes on.
  this machine   its name as 'uname -n' gives it, up to the first dot
                 ($JAIRA_SCAN_HOSTNAME overrides it, empty switches it off).
                 A hit is a warning: it names the place and does not stop.

A hit is printed as '<commit> <file>:<line>: ...', '<commit> commit
message: ...' or '<commit> author email: ...'. Where a file's own name
matches, the file is named by its number in
'git diff-tree --cc -r --name-only <commit>' instead, because its name is
the mention.

Exit codes: 0 nothing stops the push (warnings may be printed); 3 the
denylist matched — do not push; 2 usage, <base> is not a commit here, or the
denylist cannot be read or holds an invalid pattern — fix that, a scan that
cannot read its list never passes; 1 git failed.
USAGE
}

die() { echo "pre-push-scan.sh: $1" >&2; exit "${2:-1}"; }

case "${1:-}" in -h|--help) usage; exit 0 ;; esac
[ $# -eq 1 ] && [ -n "$1" ] || { usage >&2; exit 2; }
base=$1
git rev-parse --verify --quiet "$base^{commit}" >/dev/null ||
  die "$base is not a commit in this repository — fetch it first" 2

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

list=${JAIRA_MENTION_DENYLIST:-$HOME/.config/jaira/mention-denylist}
deny=$tmp/deny
: >"$deny"
denylist_state="not found at $list, its strings were not checked"
if [ -e "$list" ]; then
  [ -r "$list" ] && [ -f "$list" ] || die "the denylist $list cannot be read" 2
  # A pattern file written on Windows ends its lines in CR, which would make
  # every pattern demand a CR the content no longer has.
  tr -d '\r' <"$list" | { grep -v -e '^[[:space:]]*#' -e '^[[:space:]]*$' || true; } >"$deny"
  rc=0
  grep -Eqi -f "$deny" </dev/null 2>/dev/null || rc=$?
  [ "$rc" -le 1 ] || die "the denylist $list holds a pattern grep cannot read" 2
  if [ -s "$deny" ]; then
    denylist_state="$(wc -l <"$deny" | tr -d ' ') pattern(s) from $list"
  else
    denylist_state="$list holds no pattern, its strings were not checked"
  fi
fi

host=${JAIRA_SCAN_HOSTNAME-$(uname -n 2>/dev/null || true)}
host=${host%%.*}
case "$host" in localhost) host="" ;; esac
hostfile=$tmp/host
printf '%s' "$host" >"$hostfile"
host_state="this machine's name checked"
[ -n "$host" ] || host_state="this machine's name not checked (none known)"

# lines <patfile> <grep mode> <file>: the numbers of the lines of <file> that
# match, one per line. Exit 2 from grep is an error, never "no match". Call it
# with its output redirected to a file, never inside $(...): die in a
# subshell ends only the subshell, and the error would read as "no match".
lines() {
  [ -s "$1" ] || return 0
  local rc=0
  grep -n -a -i "$2" -f "$1" -- "$3" >"$tmp/grep" || rc=$?
  [ "$rc" -le 1 ] || die "grep failed on $3" 1
  cut -d: -f1 <"$tmp/grep"
}

stops=0
warns=0
report() { # <kind> <text>
  if [ "$1" = deny ]; then
    echo "$2: matches the denylist"
    stops=$((stops + 1))
  else
    echo "warning: $2: names this machine"
    warns=$((warns + 1))
  fi
}

# One commit. names: its changed files; adds/locs: its added lines and, line
# for line, '<path><TAB><line>' where each sits in the new file.
scan() {
  local c=$1 short kind pat mode n
  short=$(git rev-parse --short "$c")
  git -c core.quotePath=false diff-tree --cc --root -r --name-only --no-commit-id "$c" >"$tmp/names"
  git -c core.quotePath=false show --format= --cc --unified=0 --no-renames \
    --no-prefix --no-color --no-ext-diff --no-textconv --text "$c" >"$tmp/diff"
  tr -d '\000' <"$tmp/diff" | awk -v adds="$tmp/adds" -v locs="$tmp/locs" '
    BEGIN { printf "" > adds; printf "" > locs }
    /^diff --(git|cc|combined) / { h = 1; path = ""; next }
    h && /^\+\+\+ / { path = substr($0, 5); sub(/\t$/, "", path); next }
    /^@@/ {
      h = 0
      match($0, /^@+/); n = RLENGTH - 1
      ln = 0
      if (match($0, / \+[0-9]+/)) ln = substr($0, RSTART + 2, RLENGTH - 2) + 0
      next
    }
    h || /^\\/ { next }
    {
      pre = substr($0, 1, n)
      if (index(pre, "-")) next
      if (index(pre, "+")) {
        s = substr($0, n + 1); sub(/\r$/, "", s)
        print s > adds; print path "\t" ln > locs
      }
      ln++
    }'

  # Into files and read from there, never piped into grep -q: grep stops at
  # its first match, git dies of SIGPIPE, and under pipefail the hit is lost.
  git log -1 --format=%B "$c" >"$tmp/msg"
  git log -1 --format='%ae%n%ce' "$c" >"$tmp/who"

  : >"$tmp/flagged"
  for kind in deny host; do
    if [ "$kind" = deny ]; then pat=$deny mode=-E; else pat=$hostfile mode=-F; fi
    lines "$pat" "$mode" "$tmp/names" >"$tmp/n"
    while IFS= read -r n; do
      printf '%s\t%s\n' "$n" "$(sed -n "${n}p" "$tmp/names")" >>"$tmp/flagged"
      report "$kind" "$short file $n of its changed files, by its name"
    done <"$tmp/n"
  done
  for kind in deny host; do
    if [ "$kind" = deny ]; then pat=$deny mode=-E; else pat=$hostfile mode=-F; fi
    lines "$pat" "$mode" "$tmp/adds" >"$tmp/hits"
    awk -v c="$short" '
      FILENAME == ARGV[1] { i = index($0, "\t"); idx[substr($0, i + 1)] = substr($0, 1, i - 1); next }
      FILENAME == ARGV[2] { want[$0] = 1; next }
      FNR in want {
        ln = $0; sub(/.*\t/, "", ln)
        p = substr($0, 1, length($0) - length(ln) - 1)
        if (p in idx) p = "file " idx[p] " of its changed files"
        print c " " p ":" ln
      }' "$tmp/flagged" "$tmp/hits" "$tmp/locs" >"$tmp/where"
    while IFS= read -r w; do report "$kind" "$w"; done <"$tmp/where"
    lines "$pat" "$mode" "$tmp/msg" >"$tmp/n"
    [ ! -s "$tmp/n" ] || report "$kind" "$short commit message"
    lines "$pat" "$mode" "$tmp/who" >"$tmp/n"
    while IFS= read -r n; do
      if [ "$n" = 1 ]; then report "$kind" "$short author email"; else report "$kind" "$short committer email"; fi
    done <"$tmp/n"
  done
}

git rev-list --reverse "$base..HEAD" >"$tmp/commits"
count=0
while IFS= read -r c; do
  scan "$c"
  count=$((count + 1))
done <"$tmp/commits"

echo "pre-push-scan: $count commit(s) of $base..HEAD; denylist: $denylist_state; $host_state"
if [ "$stops" -gt 0 ]; then
  echo "pre-push-scan: $stops denylist hit(s) — do not push; rewrite the commits named above" >&2
  exit 3
fi
[ "$warns" -eq 0 ] || echo "pre-push-scan: $warns warning(s) for this machine's name — the push may go on"
exit 0
