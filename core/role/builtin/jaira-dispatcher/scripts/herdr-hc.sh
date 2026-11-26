#!/usr/bin/env bash
# The part of the herdr CLI that spawn.sh, run-lane.sh and reap.sh call, carried
# out by Herdr App's hc (herdr-core). Herdr App sets no HERDR_ENV and ships no
# herdr binary with tab and pane commands: it runs terminals, and the app shows
# each one as a tab. One terminal stands for one tab and its one pane, so the
# terminal's address (t7) is both ids below.
#
# Without the herdr.space label the app puts a terminal by its directory, and a
# worker in a worktree beside the repository lands outside the space the person
# is looking at. The label of the terminal this runs in is copied onto the new
# one, which the app then shows as a new tab of that same space.
set -euo pipefail

hc="${HC_BIN:-hc}"

# The value of one label out of hc's JSON. Label keys carry dots, which
# 'jaira json get' reads as a path, so they are matched as text.
labelof() {
  grep -o "\"$2\":\"[^\"]*\"" <<<"$1" | head -1 | sed 's/^"[^"]*":"//; s/"$//' || true
}
# Ids and labels written here come from spawn.sh: a ticket id, a lane name and
# a t-address. Quotes and backslashes are escaped all the same.
js() { printf '%s' "$1" | sed 's/\\/\\\\/g; s/"/\\"/g'; }

# hc with a deadline. hc get waits for the terminal it asks, and one that never
# answered hung 'tab list' and every caller with it: reap.sh under its lock, so
# no later round ran, and the cockpit, which gives up after 15 s and then sees
# no tab at all. JAIRA_HC_TIMEOUT seconds (5) are the budget of this whole call,
# not of each hc in it: 'tab list' asks one terminal after the other, and five
# mute ones at 3 s each took the cockpit's 15 s again. Once the budget is spent
# hc is killed, or not started, and the call fails like any refusal. Plain
# bash, not timeout(1): macOS ships none, and Git Bash may find Windows'
# timeout.exe first (see reap.sh). hc writes into files, not into the caller's
# pipes: what a killed hc left running would hold those open, and the caller
# would wait for it all the same.
secs="${JAIRA_HC_TIMEOUT:-5}"
case "$secs" in
  ''|*[!0-9]*|0) echo "herdr-hc.sh: JAIRA_HC_TIMEOUT must be seconds above 0, not '$secs'" >&2; exit 2 ;;
esac
end=$((SECONDS + secs))
bounded() {
  local o e p w rc=0 left=$((end - SECONDS))
  if [ "$left" -le 0 ]; then
    echo "herdr-hc.sh: hc $1 not asked: the ${secs}s were up" >&2
    return 1
  fi
  o="$(mktemp)"; e="$(mktemp)"
  "$hc" "$@" >"$o" 2>"$e" &
  p=$!
  # The watchdog: a TERM, once hc has answered, ends its sleep and it.
  ( s=""; trap '[ -z "$s" ] || kill "$s" 2>/dev/null; exit 0' TERM
    sleep "$left" >/dev/null 2>&1 & s=$!; wait "$s"
    kill -KILL "$p" 2>/dev/null && echo "herdr-hc.sh: hc $1 gave no answer before the ${secs}s were up" >&2 ) >/dev/null &
  w=$!
  wait "$p" || rc=$?
  kill "$w" 2>/dev/null || true
  wait "$w" 2>/dev/null || true
  cat "$o"; cat "$e" >&2; rm -f "$o" "$e"
  return "$rc"
}

cmd="${1:-} ${2:-}"
[ $# -ge 2 ] && shift 2 || shift $#
case "$cmd" in
  "tab create")
    cwd=""; label=""
    while [ $# -gt 0 ]; do
      case "$1" in
        --cwd)       cwd="$2"; shift 2 ;;
        --label)     label="$2"; shift 2 ;;
        --workspace) shift 2 ;;
        *)           shift ;;
      esac
    done
    args=(run --bg --json --no-notify --outlive-parent)
    # --outlive-parent: by default hc ends a terminal together with the program
    # of the one that started it, and a dispatcher has to outlive its teamlead.
    # --no-notify: hc would otherwise type notes about the worker into the
    # caller's prompt; the roles report through SendMessage and the ticket.
    [ -z "$cwd" ] || args+=(--cwd "$cwd")
    [ -z "$label" ] || args+=(--label "jaira.tab=$label")
    if [ -n "${HC_TERMINAL_ID:-}" ]; then
      space="$(labelof "$(bounded get "$HC_TERMINAL_ID" --json 2>/dev/null || true)" herdr.space)"
      [ -z "$space" ] || args+=(--label "herdr.space=$space")
    fi
    # The app's agent list shows a terminal by its hc name and nothing else:
    # without one, every tab reads 'Claude Code'. The label <ID>/<role> gives
    # <role>-<id>, in the characters and length hc accepts.
    name=""
    if [ -n "$label" ]; then
      case "$label" in */*) name="${label##*/}-${label%%/*}" ;; *) name="$label" ;; esac
      # LC_ALL=C: tr works on bytes in GNU and BSD alike. Under a UTF-8 locale
      # BSD tr reads ü as one character and GNU as two bytes, so the same label
      # gave two names; in C both give one dash per byte, ranges plain ASCII.
      name="$(printf '%s' "$name" | LC_ALL=C tr 'A-Z' 'a-z' | LC_ALL=C tr -c 'a-z0-9_-' '-')"
      name="${name:0:32}"
    fi
    # A name belongs to one live terminal, and the old tab of the same ticket
    # may still be up: -2 to -9, then no name, rather than no tab at all.
    # hc prints its refusal, JSON or not, on stderr and nothing on stdout.
    out=""; err="$(mktemp)"
    trap 'rm -f "$err"' EXIT
    for n in "" -2 -3 -4 -5 -6 -7 -8 -9 none; do
      try=("${args[@]}")
      if [ -n "$name" ] && [ "$n" != none ]; then
        try+=(--name "${name:0:$((32 - ${#n}))}$n")
      fi
      if out="$("$hc" "${try[@]}" 2>"$err")"; then break; fi
      # Only a taken name is worth another try; any other refusal is the error.
      if [ -z "$name" ] || [ "$n" = none ] || ! grep -q '"name_taken"' "$err"; then
        cat "$err" >&2
        exit 1
      fi
    done
    t="$(jaira json get address <<<"$out")"
    printf '{"result":{"tab":{"tab_id":"%s","label":"%s"},"root_pane":{"pane_id":"%s"}}}\n' \
      "$(js "$t")" "$(js "$label")" "$(js "$t")"
    ;;
  "tab list")
    # As herdr lists tabs: the cockpit finds a ticket's worker by its label,
    # which hc ls leaves out, so each live terminal is asked for its own.
    # A terminal hc get calls not_found went since hc ls and is left out. Any
    # other failure (no answer in time, or not asked once JAIRA_HC_TIMEOUT
    # was spent) may hide a live worker's tab: the list is still
    # printed, but stderr names the terminal and the exit is 3. Every caller
    # already reads a failed list as "tabs unknown" - Give back starts no
    # worker, reap.sh removes nothing. A flag in the JSON instead would have to
    # be checked by each of them, and reap.sh would go on to read the mute
    # terminal's 'pane get' failure as a closed tab.
    ls="$(bounded ls --json)"
    err="$(mktemp)"
    trap 'rm -f "$err"' EXIT
    sep=""; mute=""; printf '{"result":{"tabs":['
    for t in $(grep -o '"address":"[^"]*"' <<<"$ls" | sed 's/^"address":"//; s/"$//'); do
      if ! j="$(bounded get "$t" --json 2>"$err")"; then
        grep -q '"not_found"' "$err" || mute="$mute $t"
        continue
      fi
      printf '%s{"tab_id":"%s","label":"%s","agent_status":"%s"}' "$sep" "$(js "$t")" \
        "$(js "$(labelof "$j" jaira.tab)")" "$(js "$(jaira json get status.state --default - <<<"$j")")"
      sep=","
    done
    printf ']}}\n'
    if [ -n "$mute" ]; then
      echo "herdr-hc.sh: tab list incomplete: hc get gave no answer for$mute" >&2
      exit 3
    fi
    ;;
  "tab get")
    j="$("$hc" get "$1" --json)"
    printf '{"result":{"tab":{"tab_id":"%s","label":"%s"}}}\n' \
      "$(js "$1")" "$(js "$(labelof "$j" jaira.tab)")"
    ;;
  "tab close")
    "$hc" kill "$1" >/dev/null 2>&1 || true
    for _ in 1 2 3 4 5; do
      "$hc" rm "$1" >/dev/null 2>&1 && exit 0
      sleep 1
    done
    "$hc" rm "$1" >/dev/null
    ;;
  "pane get")
    j="$("$hc" get "$1" --json)"
    app="$(jaira json get status.app --default - <<<"$j")"
    [ "$app" = claude-code ] && app=claude
    printf '{"result":{"pane":{"pane_id":"%s","tab_id":"%s","agent":"%s","agent_status":"%s"}}}\n' \
      "$(js "$1")" "$(js "$1")" "$(js "$app")" \
      "$(js "$(jaira json get status.state --default - <<<"$j")")"
    ;;
  "pane run"|"pane send-text")
    # hc send pastes and presses Enter in one step, and refuses to type into a
    # program blocked at a dialog.
    "$hc" send "$1" "$2" >/dev/null
    ;;
  "pane send-keys")
    p="$1"; shift
    # The Enter spawn.sh presses after send-text: hc send pressed it already,
    # and a second one is not what spawn.sh meant.
    [ "$*" = enter ] || "$hc" keys "$p" "$@" >/dev/null
    ;;
  *)
    echo "herdr-hc.sh: '$cmd' is not one of the herdr calls it carries out" >&2
    exit 2
    ;;
esac
