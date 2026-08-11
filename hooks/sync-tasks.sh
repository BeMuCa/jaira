#!/usr/bin/env bash
# Mirror Claude Code's task list into the jaira backlog.
#
# Install as a PostToolUse hook matched on the task tools. It only writes ticket
# files — it never invokes a tool or a model — so it cannot trigger another round
# of syncing. An unchanged task list writes nothing at all.
#
# Failures are swallowed deliberately: a board that cannot be updated is a
# nuisance, but a hook that breaks the user's session is a much worse outcome.
set -uo pipefail

command -v jaira >/dev/null 2>&1 || exit 0
[ -d .jaira ] || exit 0

payload=$(cat 2>/dev/null) || exit 0

# The hook receives the tool call envelope; jaira wants the task list itself.
# Pull it out if jq is available, otherwise hand the whole thing over and let
# jaira's tolerant parser find it.
if command -v jq >/dev/null 2>&1; then
  tasks=$(printf '%s' "$payload" | jq -c '{tasks: (.tasks // .tool_response.tasks // .tool_input.tasks // [])}' 2>/dev/null) || exit 0
  [ "$tasks" = '{"tasks":[]}' ] && exit 0
else
  tasks=$payload
fi

printf '%s' "$tasks" | jaira sync-tasks --json >/dev/null 2>&1 || exit 0
exit 0
