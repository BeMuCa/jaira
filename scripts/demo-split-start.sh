#!/usr/bin/env bash
# Two panes side by side — ADA on the left, GRACE on the right, each in its own
# clone of the same board — and attaches to them. Meant for recording:
#
#   scripts/demo-split-start.sh        # or through scripts/demo-split.tape
#
# Why tmux rather than two windows: a recording captures exactly one terminal.
# The split therefore has to happen INSIDE that one terminal, or the second
# side is not in the picture at all.
#
# Why its own socket (-L), and this is not a nicety: a recording that talks to
# the default tmux server can reach the panes somebody is working in. It did —
# a demo command ran in another project's checkout and filed a ticket there.
# On a socket of its own this server shares nothing with the one on the
# machine: it cannot attach to those sessions, and killing it cannot end them.
set -eu
J=${J:-/tmp/jaira}
SOCKET=${SOCKET:-jaira-demo}
S=${SESSION:-demo-$$}
TM="tmux -L $SOCKET"
R=$(J="$J" "$(dirname "$0")/demo-split-setup.sh")

$TM kill-session -t "$S" 2>/dev/null || true
$TM new-session -d -s "$S" -c "$R/ada" "bash --norc"
A=$($TM list-panes -t "$S" -F '#{pane_id}' | head -1)
B=$($TM split-window -h -t "$A" -c "$R/grace" -P -F '#{pane_id}' "bash --norc")

# Names above the panes: without them a viewer cannot tell who is who, and who
# is who is the whole point of the picture.
$TM set -t "$S" -g pane-border-status top
$TM set -t "$S" -g pane-border-format ' #{pane_title} '
$TM set -t "$S" -g status off
$TM select-pane -t "$A" -T " ADA "
$TM select-pane -t "$B" -T " GRACE "

# 'mine' prints the id of the first ticket on the board. It exists so a tape can
# type "jaira pull $(mine)" without nested quotes: vhs parses the line it types,
# and a command substitution carrying its own quotes breaks that parse.
mine_fn='mine() { jaira list --json | grep -o "\"id\": \"[^\"]*\"" | head -1 | cut -d\" -f4; }'
for p in "$A" "$B"; do
  $TM send-keys -t "$p" "export JAIRA_HOME=$R/home PATH=$(dirname "$J"):\$PATH PS1='\$ '; $mine_fn; clear" Enter
done
# One key, no prefix, for changing sides. The prefix (Ctrl+B, then a letter) is
# two events with a timeout between them, and a recorder that sends them back to
# back loses the prefix and types the letter into the shell — which is exactly
# what happened: "o" turned up in front of the next command. A prefix-less
# binding has no timing in it at all.
$TM bind -n M-o select-pane -t :.+

$TM select-pane -t "$A"
echo "$R" > /tmp/jaira-demo-sandbox
exec $TM attach -t "$S"
