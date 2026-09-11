#!/usr/bin/env bash
# Two machines side by side: ADA and GRACE, each in a container with its own
# hostname, its own filesystem and its own home. The only path they share is the
# bare repository — which is to say, the remote.
#
# Both sides are containers rather than one being this host, for two reasons:
# the claim "two machines" is then equally true of both, and a recording meant
# to be shown to other people does not carry the name of somebody's laptop.
#
#   scripts/demo-machines-start.sh     # or through the tapes that call it
#
# Honest about what this is: a container shares the host's kernel, so it is not
# a second computer in the strict sense. Everything that matters to the claim is
# separate though — hostname, filesystem, clone, home — and the pane titles say
# "another machine" only because that is demonstrable on camera with hostname,
# not because it reads well.
set -eu
J=${J:-/tmp/jaira-static}      # static build: the container has no glibc
IMAGE=${IMAGE:-alpine/git:latest}
SOCKET=${SOCKET:-jaira-demo}   # a server of its own: a recording must never be
S=${SESSION:-machines-$$}      # able to reach the panes somebody is working in
FETCH_EVERY=${FETCH_EVERY:-3s} # a recording cannot wait ten minutes
TM="tmux -L $SOCKET"

command -v docker >/dev/null || { echo "docker is needed for this one"; exit 1; }
[ -x "$J" ] || { echo "build it first: CGO_ENABLED=0 go build -o $J ./cmd/jaira"; exit 1; }

R=$(mktemp -d)
git init -q --bare "$R/board.git"

# A master branch with something in it, so that "git branch -a" prints a real
# list on both sides. An empty repository lists nothing, and on a recording
# "nothing" is indistinguishable from "the command did not work" — the point
# being made is that ADA's branch is MISSING from a list, which needs a list.
seed=$(mktemp -d)
git clone -q "$R/board.git" "$seed" 2>/dev/null
printf 'the exporter\n' > "$seed/README.md"
git -C "$seed" add README.md
git -C "$seed" -c user.name=ada -c user.email=ada@example.test commit -qm "the project"
git -C "$seed" push -q origin HEAD:master
rm -rf "$seed"

# Each machine's home, prepared from here so both panes start clean.
for who in ada grace; do
  mkdir -p "$R/${who}home/.jaira"
  printf '{ "fetch-every": "%s" }\n' "$FETCH_EVERY" > "$R/${who}home/.jaira/settings.json"
done

# box prints the docker command that is one machine.
box() {
  local who=$1
  echo "docker run --rm -it --entrypoint sh --user $(id -u):$(id -g) --hostname ${who}-laptop \
    -v $J:/usr/local/bin/jaira:ro -v $R/board.git:$R/board.git -v $R/${who}home:/home/$who \
    -e HOME=/home/$who -e TERM=xterm-256color $IMAGE"
}

$TM kill-session -t "$S" 2>/dev/null || true
$TM new-session -d -s "$S" "$(box ada)"
A=$($TM list-panes -t "$S" -F '#{pane_id}' | head -1)
B=$($TM split-window -h -t "$A" -P -F '#{pane_id}' "$(box grace)")

# The labels are written into the border format by pane index, not taken from
# the pane title. A program running in the pane sets that title with an escape
# sequence and tmux always accepts it — the board did, and the right-hand label
# turned into "jaira" mid-recording. A format the pane cannot reach stays put.
$TM set -t "$S" -g pane-border-status top
$TM set -t "$S" -g pane-border-format \
  ' #{?#{==:#{pane_index},0}, ADA · this machine , GRACE · another machine } '
$TM set -t "$S" -g status off

# One key, no prefix, for changing sides. tmux's own prefix is two events with a
# timeout between them, and a recorder sending them back to back loses it — the
# letter then lands in the shell instead of switching panes.
$TM bind -n M-o select-pane -t :.+

# Narration in a tape is typed as a plain shell comment ("# ada writes down a
# problem"), not through a helper: a helper is echoed as the command AND printed
# as its output, so every line of explanation appears twice. A comment is
# executed by nobody and stays on screen exactly once.

# 'mine' prints the id of the first ticket on the board, so a tape can type
# "jaira pull $(mine)" with no nested quotes: vhs parses what it types, and a
# substitution carrying its own quotes breaks that parse.
mine_fn='mine() { jaira list --json | grep -o "\"id\": \"[^\"]*\"" | head -1 | cut -d\" -f4; }'

# Each machine starts mid-work, on branches of its own. Two lists that differ
# say more than one list with a gap in it: neither person can see what the other
# is doing, and master is all they have in common — which is exactly the state a
# ticket has to cross.
branches_ada="spike/safari-cookies fix/302"
branches_grace="chore/bump-deps feat/exporter-retry"

for who in ada:"$A" grace:"$B"; do
  name=${who%%:*}; pane=${who#*:}
  eval "brs=\$branches_$name"
  work="git clone -q $R/board.git /home/$name/work 2>/dev/null; cd /home/$name/work"
  work="$work; git config user.name $name; git config user.email $name@example.test"
  for b in $brs; do
    work="$work; git checkout -q -b $b; echo '$name was here' >> notes-$name.md; git add -A; git commit -qm 'work on $b'"
  done
  $TM send-keys -t "$pane" "export PS1='$name\$ '; $work; jaira init >/dev/null 2>&1; $mine_fn; clear" Enter
done

$TM select-pane -t "$A"
echo "$R" > /tmp/jaira-demo-sandbox
exec $TM attach -t "$S"
