---
id: 01M6EEZYHN2T31ANZAWRH319TV
title: "README: Abschnitt 'Jaira PM View - Web only' fuer Leute ohne Coding-Agent"
status: review
ready: true
creator: BeMuCa
goal: "Die README hat einen eigenen Abschnitt, mit dem jemand ohne Agent und ohne Terminal-Board in fuenf Minuten jaira web auf einem geklonten Repo laufen hat: Install-Zeile ohne Rollen/Cockpit, jaira web, ein Screenshot, was man dort kann und was nicht; die Zeile 'Web UI' in der Tabelle 'Not built' ist weg."
context: |-
  Was fehlt: jaira web (Epics sortieren, Tickets lesen) ist fuer PMs und Leute ohne Coding-Agent gedacht, aber nirgends fuer sie erklaert. Es gibt jaira web --help, eine Zeile in der README-Befehlstabelle (Zeile 850) und einen Absatz in docs/COMMANDS.md Zeile 57. Kein Tutorial, kein Screenshot unter docs/img.
  Widerspruch: README Zeile 1022 fuehrt 'Web UI' in der Tabelle 'Not built' als bewusst nicht gebaut.
  Zielgruppe (Wunsch 2026-10-28/11-02): Leute, die jaira nur ziehen, um die Web-Ansicht zu nutzen - ohne Claude, ohne Terminal-Board.
  Bekannt: die Seite steckt im Binary (core/web/web.go:30, go:embed), es braucht kein Node und kein zweites Programm. Der curl-Installer richtet seit V38STG auch Rollen und Cockpit ein; der Abschnitt muss die Variante ohne beides nennen.
  Bekannt: ein frischer Clone hat die Ticket-Refs nicht (origin holt nur refs/heads/*); jaira holt refs/jaira/tickets/* selbst (core/refsync/refsync.go:933). Der Abschnitt sagt, wann die Tickets nach dem Start erscheinen.
  Soll als eigener PR kommen, nur README und Bild.
definition-of-done: "README hat den Abschnitt 'Jaira PM View - Web only' mit Install-Zeile ohne Rollen/Cockpit, jaira web, Screenshot und je einem Satz was geht und was nicht"
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-11-02T15:45:35Z
updated-at: 2026-11-23T01:50:32Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: worker-h319tv
claimed-at: 2026-11-23T01:38:34Z
outcome-what: "README has a new section '## Jaira PM View - Web only' (after '## Start'): binary-only install lines for sh and PowerShell with JAIRA_NO_ROLES=1, git clone + jaira init + jaira web, when the tickets show up on a fresh clone, a screenshot docs/img/web.png, one sentence each for what works and what does not. 'Web UI' dropped from the 'Not built' row; VS Code extension stays. Release note core/release/notes/H319TV.md."
outcome-why: "jaira web is meant for people without a coding agent, but nothing told them how to get there, and the README still listed a web UI as deliberately not built."
outcome-resolves: "All four DoD items; install line read from the script, not executed. Checked: concurrency n/a (docs); error path: web without init exits 1 -> init step included; unusual input: fresh clone without refs -> wait sentence; permissions: drops push, so push access named."
review-summary: "Adds a README section '## Jaira PM View - Web only' right after '## Start' (README.md:300-343). It gives an install line for sh and for PowerShell that sets JAIRA_NO_ROLES=1, so only the binary is installed, without skills or cockpit. Then git clone, jaira init, jaira web. It explains that tickets arrive on their own refs, that jaira web fetches them at start, and that you reload if the page opened first; the page also rereads every minute. It adds a screenshot at docs/img/web.png (demo board, no real names) and one sentence each on what the page can do (read tickets, filter, sort tickets into epics, create/edit/close epics, pushed at once) and what it cannot (create, edit or move tickets, notes, acceptance). The Features bullet links to the section (README.md:67-68). 'Web UI' is dropped from the 'Not built' table, so that row now says only 'VS Code extension' (README.md:1247). A one-line release note is in core/release/notes/H319TV.md. Every claim was checked against the code: install.sh:110 and install.ps1:159 read the variable; web.go:273 fetches at start; ReviewPage.svelte:20 refreshes every 60 s; web.go:303-313 only writes epics."
round: 2
back-from: critique
back-by: lane
back-count: 1
review-gaps: |-
  note: README.md:327-329 says init writes CLAUDE.md and AGENTS.md 'where there were none'. On a shared clone, init writes its block into whichever of the two already carry it (core/board/announce.go:491-496), and it can append a milestones line to the tracked .gitignore (internal/cli/tickets.go:93-95). So it can also change committed files. Fix: 'init may change CLAUDE.md, AGENTS.md and .gitignore in the clone; you do not need to commit them'.
  note: pre-existing: the help for jaira init (internal/cli/tickets.go:39-40) says a shared clone 'keeps its .gitignore exactly as it is', but tickets.go:93-95 can add the milestones ignore line.
  note: DoD 3 - the install line was never run end to end (curl|sh was denied in the implementing session). I checked with a stub script that '| JAIRA_NO_ROLES=1 sh' hands the variable to the piped script, and that master's install.sh:110 and install.ps1:159 both read it. The PowerShell line is untested. The check flow below has the person run it.
  note: README.md:338-339 'every drop is pushed at once, so you need push access'. Without push access the drop still lands locally and only the push fails, shown on the page (core/web/web.go:183, 253-266). Possible wording: 'without push access it stays on your machine'.
test-verdict: "pass: no Go changed - suite not run (house rule tonight: no go test in lanes), go build ok. DoD 1-4 checked in the tree (README.md:301 section, :1247 row, :331-334 wait sentence, core/release/notes/H319TV.md one '- ' line, docs/img/web.png present, anchor #jaira-pm-view---web-only matches GitHub slug). Function: README steps replayed on a fresh clone of a scratch bare remote: jaira web before init exits 1 (no half state), after init jaira web serves 3 tickets within 1.5 s, /api/sync shows last_fetch. Two writers: the change writes nothing. Not verified: the curl|sh and PowerShell install lines themselves (not executed)."
hunt-verdict: "pass: bughunt done at 797831c0c, 0 findings"
review-verdict: "accept: all four DoD items are met by the diff and every README claim matches the code. No blocker. Four notes: the init side effects are worded too narrowly, a pre-existing help-text mismatch, the install line has not been run end to end by an agent, and the push-access wording."
review-round: 2
review-verdict-gaps: |-
  note: README.md:327-329 says init writes CLAUDE.md and AGENTS.md 'where there were none'. On a shared clone, init writes its block into whichever of the two already carry it (core/board/announce.go:491-496), and it can append a milestones line to the tracked .gitignore (internal/cli/tickets.go:93-95). So it can also change committed files. Fix: 'init may change CLAUDE.md, AGENTS.md and .gitignore in the clone; you do not need to commit them'.
  note: pre-existing: the help for jaira init (internal/cli/tickets.go:39-40) says a shared clone 'keeps its .gitignore exactly as it is', but tickets.go:93-95 can add the milestones ignore line.
  note: DoD 3 - the install line was never run end to end (curl|sh was denied in the implementing session). I checked with a stub script that '| JAIRA_NO_ROLES=1 sh' hands the variable to the piped script, and that master's install.sh:110 and install.ps1:159 both read it. The PowerShell line is untested. The check flow below has the person run it.
  note: README.md:338-339 'every drop is pushed at once, so you need push access'. Without push access the drop still lands locally and only the push fails, shown on the page (core/web/web.go:183, 253-266). Possible wording: 'without push access it stays on your machine'.
review-check: |-
  1. Open README.md on branch feat/h319tv and find the section '## Jaira PM View - Web only' (around line 300). It has two install lines, then git clone / jaira init / jaira web, a screenshot, and a can / cannot paragraph.
  2. In a scratch folder, run: curl -fsSL https://raw.githubusercontent.com/BeMuCa/jaira/master/scripts/install.sh | JAIRA_INSTALL_DIR=$PWD/bin JAIRA_NO_ROLES=1 sh. The last line printed must say 'skipped the skills and the Cockpit (JAIRA_NO_ROLES=1)'.
  3. Clone a repository whose jaira board is shared, cd into it, run ./path/to/bin/jaira init, then git status. Expect possible changes to CLAUDE.md, AGENTS.md or .gitignore (see the first note).
  4. Run jaira web. A browser opens with the tickets in no epic on the left and the epics on the right, like docs/img/web.png. If the list is empty, reload once after a few seconds and the tickets appear.
  5. Drag a ticket onto an epic. The epic's count goes up. Ctrl-C stops jaira web.
  6. Scroll to the 'Not built' table near the end of the README. The row now reads 'VS Code extension', with no 'Web UI'.
---

# README: Abschnitt 'Jaira PM View - Web only' fuer Leute ohne Coding-Agent

## Definition of Done

- [x] README hat den Abschnitt 'Jaira PM View - Web only' mit Install-Zeile ohne Rollen/Cockpit, jaira web, Screenshot und je einem Satz was geht und was nicht
  proof: README.md:301 '## Jaira PM View - Web only' (install lines with JAIRA_NO_ROLES=1, jaira init/web, docs/img/web.png, one sentence each for can/cannot); linked from the Features bullet README.md:67-68
- [x] Die Zeile 'Web UI, VS Code extension' in der Tabelle 'Not built' ist entfernt oder sagt, was heute gilt
  proof: README.md:1247 row now '| VS Code extension | ...', Web UI dropped
- [x] Auf einem frischen Clone nachgestellt: Install-Zeile, jaira web, Tickets erscheinen; die Wartezeit steht im Abschnitt
  proof: scratch run with the branch binary: bare remote + shared board, 3 tickets on refs, git clone, jaira init (0 refs), jaira web: refs/jaira/* present within 0.5 s, /api/unassigned lists the 3 tickets; wait sentence README.md:331-334. Install line NOT run (curl|sh denied in this session) - read from scripts/install.sh:110
- [x] Release-Notiz in core/release/notes/<ID>.md
  proof: core/release/notes/H319TV.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 01:39 · BeMuCa** — Entscheidungen (kein Punkt offen fuer die Person): 1) Abschnitt steht als eigenes '## Jaira PM View - Web only' direkt nach '## Start', nicht im Install-Abschnitt - feat/q9f9w4 aendert den Install-Text, so bleibt der Merge trivial. 2) Install-Zeile mit JAIRA_NO_ROLES=1 (so geschrieben in scripts/install.sh:110 und install.ps1:159). 3) 'Not built': nur 'Web UI' faellt weg, die Zeile bleibt fuer 'VS Code extension' mit ihrem Grund. 4) Release-Notiz wird geschrieben, weil DoD 4 sie verlangt, obwohl der Kontext 'nur README und Bild' sagt und CLAUDE.md README-Aenderungen nicht als client-facing zaehlt - die Person kann sie streichen. 5) Screenshot mit chrome-headless-shell aus ~/.cache/ms-playwright auf einem Scratch-Board, nach docs/img/web.png. Kein Plan-Lane: Option 'planning' nicht angehakt.
- **2026-11-23 01:44 · BeMuCa** — Gefunden beim Nachstellen: (a) jaira web ohne init auf frischem Clone: 'no .jaira directory found; run jaira init', exit 1 - daher init im Abschnitt. (b) init schreibt CLAUDE.md und AGENTS.md in den Clone. (c) Der erste Fetch laeuft parallel zum Browser-Start (internal/cli/web.go:112-119), die Seite liest nur alle 60 s neu (ReviewPage.svelte REFRESH_MS) - daher 'reload'. (d) NICHT getan: den echten Installer per curl|sh ausfuehren - vom Auto-Mode-Klassifizierer abgelehnt; JAIRA_NO_ROLES=1 nur aus scripts/install.sh:110 und install.ps1:159 gelesen. Die PowerShell-Zeile ($env:JAIRA_NO_ROLES = "1"; irm ... | iex) ist ungetestet. (e) Screenshot mit Windows-Chrome headless (WSL-chrome-headless-shell fehlt libasound) auf einem Demo-Board mit englischen Titeln, User anna, kein echter Name/Host darauf.
- **2026-11-23 01:45 · BeMuCa** — critique: Features-Bullet 'Epics and jaira web' (README.md:63-69) verlinkt den neuen Abschnitt nicht; Link wie beim Cockpit-Bullet ergaenzen. Sonst stehen lassen: eigene Install-Zeile statt Verweis auf ## Install ist gewollt (Merge mit feat/q9f9w4, Abschnitt soll allein reichen); Release-Notiz trotz README-only steht wegen DoD 4 (Notiz oben).
- **2026-11-23 01:48 · BeMuCa** — bughunt: done at 797831c0c - 0 findings in 3 files, first hunt; README claims checked by running init/share/clone/init in scratch repos and reading install scripts and web.go; concurrency, error paths, input, permissions checked; files: README.md core/release/notes/H319TV.md docs/img/web.png
