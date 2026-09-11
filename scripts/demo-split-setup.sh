#!/usr/bin/env bash
# Builds the sandbox for the split-screen walkthrough and prints its path.
# Called during the hidden phase of a recording, so the picture carries only
# the commands that are the point.
set -eu
J=${J:-/tmp/jaira}
R=$(mktemp -d)
git init -q --bare "$R/board.git"
for n in ada grace; do
  git clone -q "$R/board.git" "$R/$n" 2>/dev/null
  git -C "$R/$n" config user.name "$n"
  git -C "$R/$n" config user.email "$n@example.test"
  JAIRA_HOME="$R/home" "$J" -C "$R/$n" init >/dev/null
done

# A recording cannot wait ten minutes for the board to notice a new ticket, so
# the sandbox asks for three seconds. This is a setting, not a demo-only hack:
# the same field is what a person would slow down or speed up for themselves.
mkdir -p "$R/home"
printf '{ "fetch-every": "%s" }\n' "${FETCH_EVERY:-3s}" > "$R/home/settings.json"

echo "$R"
