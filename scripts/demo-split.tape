# Two people hand a ticket over without saying a word to each other.
#
#   CGO_ENABLED=0 go build -o /tmp/jaira-static ./cmd/jaira
#   vhs scripts/demo-split.tape        → docs/img/demo-split.gif
#
# Needs docker, tmux, ttyd and ffmpeg. Recorded with vhs v0.11.0: v0.12.0 runs
# the tape and writes no file at all (charmbracelet/vhs#787).
#
# vhs records exactly one terminal and cannot split by itself. The split
# therefore happens inside the recording, with tmux: typing always goes to the
# active pane, and Alt+o moves to the other side. Alt+o rather than tmux's own
# prefix, because a prefix is two events with a timeout between them and a
# recorder sending them back to back loses it — the letter then lands in the
# shell instead.

Output docs/img/demo-split.gif

Set Shell bash
Set FontSize 14
Set Width 1500
Set Height 700
Set Padding 15
Set Theme "Catppuccin Mocha"
Set TypingSpeed 30ms

Hide
Type "CGO_ENABLED=0 go build -o /tmp/jaira-static ./cmd/jaira && bash scripts/demo-machines-start.sh" Enter
Sleep 18s
Show
Sleep 1s

# ── two machines, mid-work, sharing only a remote ────────────────────
Type "# two machines, two clones, and the only thing they share is a remote"
Enter
Sleep 2s
Type "git branch -a"
Enter
Sleep 3s

Alt+o
Sleep 800ms
Type "# grace is on her own branches. fix/302 does not exist here"
Enter
Sleep 2s
Type "git branch -a"
Enter
Sleep 3s
Alt+o
Sleep 800ms

# ── ada writes down a problem. she does not say who should fix it ────
Type "# ada writes down a problem. she does not say who should do it"
Enter
Sleep 2500ms

Type "jaira create 'session cookie dropped on 302' --goal 'the cookie survives the OAuth round-trip' --context 'reported in chat: Safari logs people out mid-flow'  --dod 'session survives the redirect'"
Sleep 500ms
Enter
Sleep 3s

Type "# nobody owns it, and it is not even on her disk — it is on its ref"
Enter
Sleep 2s
Type "ls .jaira/tickets/ | wc -l"
Enter
Sleep 2500ms

# ── grace looks at the board of her own accord ───────────────────────
Alt+o
Sleep 800ms
Type "# grace looks at the board. nobody asked her to"
Enter
Sleep 2500ms
Type "jaira fetch"
Enter
Sleep 3s

Type "# she decides to take it. pulling is a compare-and-swap: one clone wins"
Enter
Sleep 2500ms
Type "jaira pull $(mine)"
Enter
Sleep 3500ms

# ── ada finds out the same way she would find out anything ───────────
Alt+o
Sleep 800ms
Type "# ada looks again"
Enter
Sleep 2s
Type "jaira fetch"
Enter
Sleep 3s
Type "jaira list"
Enter
Sleep 3500ms

Type "# thanks, grace. not a word was exchanged — the board said it"
Enter
Sleep 4s
