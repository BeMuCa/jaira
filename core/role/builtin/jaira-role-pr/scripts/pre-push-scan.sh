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
push publishes the whole history — in every place the push publishes: the
lines it adds, the names of the files it changes, its message, and its
author's and committer's name and email. A merge is checked for what it adds
itself, a conflict resolution, and not for what it brings in from its
parents.

Two lists are matched, case-insensitively:

  the denylist   one extended regular expression per line, '#' comments and
                 blank lines ignored, read from $JAIRA_MENTION_DENYLIST or
                 else ~/.config/jaira/mention-denylist. It lives outside the
                 repository on purpose: a list inside it would be the mention.
                 UTF-8, with or without a BOM, or UTF-16 with its BOM; any
                 other encoding is refused rather than half read. Lines may
                 end in LF, CRLF or CR; blanks at a line's end are ignored.
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
cannot read its list never passes; any other non-zero code — git's own —
git failed.
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

# grep -i folds the case of a non-ASCII letter only in a UTF-8 locale; the
# caller's may be C. Take the first locale here that folds Ü to ü.
fold=C
fold_state="non-ASCII letters compared case-sensitively (no UTF-8 locale here)"
for l in "${LC_ALL:-}" "${LC_CTYPE:-}" "${LANG:-}" C.UTF-8 C.utf8 en_US.UTF-8 en_US.utf8; do
  [ -n "$l" ] || continue
  if LC_ALL=$l grep -qi "$(printf '\303\274')" <<<"$(printf '\303\234')" 2>/dev/null; then
    fold=$l fold_state=""
    break
  fi
done

list=${JAIRA_MENTION_DENYLIST:-$HOME/.config/jaira/mention-denylist}
deny=$tmp/deny
: >"$deny"
denylist_state="not found at $list, its strings were not checked"
if [ -e "$list" ]; then
  [ -r "$list" ] && [ -f "$list" ] || die "the denylist $list cannot be read" 2
  # Read the list as it was written, or not at all: a pattern dropped on the
  # way reads as "no hit" and the push goes out.
  case "$(head -c 2 "$list" | od -An -tx1 | tr -d ' \n')" in
  fffe | feff)
    command -v iconv >/dev/null || die "the denylist $list is UTF-16 and iconv is missing" 2
    iconv -f UTF-16 -t UTF-8 <"$list" >"$tmp/list" 2>/dev/null || die "the denylist $list is not valid UTF-16" 2
    ;;
  *) cat <"$list" >"$tmp/list" ;;
  esac
  [ "$(tr -d '\000' <"$tmp/list" | wc -c)" = "$(wc -c <"$tmp/list")" ] ||
    die "the denylist $list holds NUL bytes — save it as UTF-8" 2
  if LC_ALL=C grep -q '[^ -~[:space:]]' "$tmp/list"; then
    command -v iconv >/dev/null || die "the denylist $list holds non-ASCII text and iconv is missing to check it is UTF-8" 2
    iconv -f UTF-8 -t UTF-8 <"$tmp/list" >/dev/null 2>&1 || die "the denylist $list is not UTF-8 — save it as UTF-8" 2
  fi
  # A BOM would stick to the first pattern. A line ends in LF, CRLF or a lone
  # CR (an old Mac file): every CR becomes a line break, and the empty line
  # CRLF leaves behind is dropped with the blank ones. Trailing blanks are an
  # editor's, not the pattern's: kept, they would demand a blank in the text.
  LC_ALL=C sed "1s/^$(printf '\357\273\277')//" "$tmp/list" | tr '\r' '\n' |
    LC_ALL=C sed 's/[[:blank:]]*$//' |
    { LC_ALL=C grep -a -v -e '^[[:space:]]*#' -e '^[[:space:]]*$' || true; } >"$deny"
  rc=0
  LC_ALL=$fold grep -Eqi -f "$deny" </dev/null 2>/dev/null || rc=$?
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
  LC_ALL=$fold grep -n -a -i "$2" -f "$1" -- "$3" >"$tmp/grep" || rc=$?
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
  git log -1 --format='%ae%n%ce%n%an%n%cn' "$c" >"$tmp/who"

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
      case $n in
      1) report "$kind" "$short author email" ;;
      2) report "$kind" "$short committer email" ;;
      3) report "$kind" "$short author name" ;;
      *) report "$kind" "$short committer name" ;;
      esac
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
[ -z "$fold_state" ] || echo "pre-push-scan: $fold_state"
if [ "$stops" -gt 0 ]; then
  echo "pre-push-scan: $stops denylist hit(s) — do not push; rewrite the commits named above" >&2
  exit 3
fi
[ "$warns" -eq 0 ] || echo "pre-push-scan: $warns warning(s) for this machine's name — the push may go on"
exit 0
