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
      space="$(labelof "$("$hc" get "$HC_TERMINAL_ID" --json 2>/dev/null || true)" herdr.space)"
      [ -z "$space" ] || args+=(--label "herdr.space=$space")
    fi
    t="$("$hc" "${args[@]}" | jaira json get address)"
    printf '{"result":{"tab":{"tab_id":"%s","label":"%s"},"root_pane":{"pane_id":"%s"}}}\n' \
      "$(js "$t")" "$(js "$label")" "$(js "$t")"
    ;;
  "tab list")
    "$hc" ls --json
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
