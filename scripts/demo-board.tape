# The board picks a ticket up by itself: ADA files one on this machine, and the
# open board on the other machine notices it without anybody pressing a key.
#
# The right pane is a container with its own hostname, filesystem and home. The
# only path it shares with the host is the bare repository — the remote.
#
#   go install github.com/charmbracelet/vhs@latest
#   sudo apt install ttyd ffmpeg tmux
#   vhs scripts/demo-board.tape        → docs/img/demo-board.gif
#
# The sandbox sets "fetch-every": "3s" in its settings, because a recording
# cannot wait ten minutes for the default. It is the same field a person would
# change for themselves — nothing here is a demo-only switch.

Output docs/img/demo-board.gif

Set Shell bash
Set FontSize 14
Set Width 1500
Set Height 750
Set Padding 15
Set Theme "Catppuccin Mocha"
Set TypingSpeed 30ms

Hide
Type "CGO_ENABLED=0 go build -o /tmp/jaira-static ./cmd/jaira && bash scripts/demo-machines-start.sh" Enter
Sleep 8s

# One ticket exists already, so the board has something to show before the
# interesting part: a card that is on a ref and in nobody's branch here.
Type "jaira create 'session cookie dropped on 302' --goal 'the cookie survives the OAuth round-trip' --context 'reported in chat: Safari logs people out mid-flow' --dod 'session survives the redirect' --assignee grace" Enter
Sleep 3s
Alt+o
Type "jaira fetch --quiet" Enter
Sleep 3s
# Opened the way a person opens it: 'jaira' puts you on the launcher, and enter
# goes into the board. 'jaira board' would skip the step everybody actually
# takes.
Type "clear && jaira" Enter
Sleep 5s
Enter
Sleep 4s
Alt+o
Type "clear" Enter
Show
Sleep 2s

# ── GRACE's board is open on the right. ADA files a second ticket on the left.
Type "jaira create 'retry the exporter on a 503' --goal 'a flaky upload finishes' --context 'nightly export fails once a week and nobody notices' --dod 'one retry, logged' --assignee grace"
Sleep 500ms
Enter

# Nobody presses anything on the right. The board fetches on its own.
Sleep 12s
