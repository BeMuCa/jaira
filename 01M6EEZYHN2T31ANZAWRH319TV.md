---
id: 01M6EEZYHN2T31ANZAWRH319TV
title: "README: Abschnitt 'Jaira PM View - Web only' fuer Leute ohne Coding-Agent"
status: testing
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
updated-at: 2026-11-23T01:46:46Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: worker-h319tv
claimed-at: 2026-11-23T01:38:34Z
outcome-what: "README has a new section '## Jaira PM View - Web only' (after '## Start'): binary-only install lines for sh and PowerShell with JAIRA_NO_ROLES=1, git clone + jaira init + jaira web, when the tickets show up on a fresh clone, a screenshot docs/img/web.png, one sentence each for what works and what does not. 'Web UI' dropped from the 'Not built' row; VS Code extension stays. Release note core/release/notes/H319TV.md."
outcome-why: "jaira web is meant for people without a coding agent, but nothing told them how to get there, and the README still listed a web UI as deliberately not built."
outcome-resolves: "All four DoD items; install line read from the script, not executed. Checked: concurrency n/a (docs); error path: web without init exits 1 -> init step included; unusual input: fresh clone without refs -> wait sentence; permissions: drops push, so push access named."
review-summary: "none"
round: 2
back-from: critique
back-by: lane
back-count: 1
review-gaps: "removed the '--port/--no-open' sentence from the new section: it repeated the Commands table (README '## Commands') and docs/COMMANDS.md:57, and a PM does not need it (jaira web prints the address when no browser opens). Left the clone/init steps that overlap '## Start' (teammate join): the section must stand alone. No code, so passes 2-4 n/a."
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
