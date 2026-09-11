# Two machines: ADA on this host, GRACE in a container with its own hostname
# and its own filesystem. The only path they share is the bare repository —
# which is to say, the remote.
#
#   CGO_ENABLED=0 go build -o /tmp/jaira-static ./cmd/jaira
#   vhs scripts/demo-machines.tape     → docs/img/demo-machines.gif
#
# Needs docker as well as vhs, ttyd, ffmpeg and tmux. Recorded with vhs v0.11.0:
# v0.12.0 runs the tape and writes no file at all (charmbracelet/vhs#787).

Output docs/img/demo-machines.gif

Set Shell bash
Set FontSize 14
Set Width 1500
Set Height 700
Set Padding 15
Set Theme "Catppuccin Mocha"
Set TypingSpeed 30ms

Hide
Type "CGO_ENABLED=0 go build -o /tmp/jaira-static ./cmd/jaira && bash scripts/demo-machines-start.sh" Enter
Sleep 15s
Show
Sleep 1s

# ── who is who: two different hosts ──────────────────────────────────
Type "hostname"
Enter
Sleep 2s
Alt+o
Type "hostname"
Enter
Sleep 2s

# ── and nothing in common but the remote ─────────────────────────────
Type "git branch -a | wc -l"
Enter
Sleep 2s
Alt+o

Type "jaira create 'session cookie dropped on 302' --goal 'the cookie survives the OAuth round-trip' --context 'reported in chat: Safari logs people out mid-flow' --dod 'session survives the redirect' --assignee grace"
Sleep 500ms
Enter
Sleep 3s

Alt+o
Type "jaira fetch"
Enter
Sleep 3s
Type "jaira pull $(mine)"
Enter
Sleep 4s
