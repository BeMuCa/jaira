---
id: 01M52C52KX820CPG5ZW8F80QZT
title: "Cockpit: Worker als farbige Claude-Icons auf einer Lane-Spur im Overview"
status: testing
ready: true
creator: Alexander Sacharov
goal: "Das Overview-Pane des cockpits bekommt einen dritten Tab 'Workers': die Lanes als Stationen (◯ lane ─ ◯ lane), jeder Agent als ✻ in der Farbe seines Tickets mit dem 6-stelligen Handle an der Station seiner Lane; wechselt das Ticket die Lane, läuft das ✻ sichtbar hinüber. working dreht sich, idle schläft gedimmt, eine offene Frage steht als ?-Blase darüber."
context: |-
  Was heute fehlt: das cockpit zeigt laufende Worker nur als Zahl im Band ('5 worker', hooks/texts.ts:87) und als Textzeilen im Tool 'workers' (hooks/sessions.ts:119 workerLines).
  Bei 5 parallelen Workern sieht Alex nicht auf einen Blick, wer an welchem Ticket in welcher Lane steht.
  Wunsch von Alex am 2026-10-16: alle Worker als Claude-Icon, aber unterschiedlich gefärbt, Name = 6-stelliger Ticket-Handle.
  Entschieden: Ort ist das Overview-Pane (hooks/overview-pane.tsx), als Spur der Lanes, nicht das Band.
  Schon da: personColor in hooks/overview-pane.tsx:23 färbt Personen per Hash aus 6 ANSI-Farben - dasselbe Muster passt für Handles.
  Schon da: welche Agenten laufen, wissen snapshotOf (hooks/board.ts, Herdr-Tabs '<id>/<lane>' + Subagenten) und workerLines; die Lane des Tickets kommt aus lanesOf.
  Offen für die Plan-Lane: bei 6 Farben und 5 Workern kollidiert ein reiner Hash oft. Vorschlag: Hash als Startplatz, bei Kollision unter den gerade laufenden zur nächsten freien Farbe weiterrücken - dann bleibt die Farbe stabil, solange der Worker lebt.
  Offen: nur Lanes mit Agent zeigen oder alle agentischen Lanes (leere als Platzhalter).
definition-of-done: Overview zeigt eine Spur der agentischen Lanes mit je einer ✻-Glyphe pro laufendem Agenten unter der Lane seines Tickets
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-16T12:49:22Z
updated-at: 2026-10-16T13:08:57Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-51728
claimed-at: 2026-10-16T12:50:47Z
outcome-what: "Overview im cockpit hat einen dritten Tab Workers: Lanes als Stationen, jeder Agent ein ✻ in Ticketfarbe mit Handle, dreht bei working, z z bei idle, ?-Blase bei offener Frage, läuft beim Lanewechsel Station für Station; Timer 250 ms nur solange der Tab offen ist, Board alle 5 s neu gelesen"
outcome-why: "Bei 5 parallelen Workern sah Alex nur eine Zahl im Band und Textzeilen im workers-Tool, nicht wer an welchem Ticket in welcher Lane steht; Alex wollte es verspielt, wie ein kleines Spiel"
outcome-resolves: F80QZT
review-summary: none
round: 2
review-gaps: "overview.ts RUNNING war eine Kopie von WORKING in board.ts:66 - entfernt, WORKING exportiert und genutzt; Ticket-Typ in overview.ts bekommt question statt Schnitttyp im parse. Gelassen: workersOf und snapshotOf lesen Tab-Labels beide per split('/'), zählen aber nach verschiedenen Regeln (Band nur working, Workers-Tab auch idle, laut decision-Notiz); readWorkers liest jaira lanes alle 5 s mit, billig und hält das Overview-Muster (openOverview liest auch alles neu). Kein toter Code gefunden."
---

# Cockpit: Worker als farbige Claude-Icons auf einer Lane-Spur im Overview

## Definition of Done

- [x] Overview hat einen dritten Tab 'Workers' mit den Lanes als Stationen (◯ lane ─ ◯ lane) und je einem ✻ mit Handle pro Agent an der Station der Lane seines Tickets
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx workers(); core/mod/builtin/jaira-cockpit/hooks/overview.ts workersOf; test 'the workers tab > every agent stands at the station of its ticket's lane'
- [x] Jede Glyphe trägt den 6-stelligen Handle; zwei gleichzeitig laufende Agenten haben nie dieselbe Farbe, solange es nicht mehr Agenten als Farben gibt
  proof: core/mod/builtin/jaira-cockpit/hooks/overview.ts colorsOf; test 'the workers tab > two tickets at work never share a colour, and each keeps its own on every read'
- [x] Die Farbe eines Tickets bleibt zwischen zwei Aktualisierungen gleich
  proof: colorsOf keeps colours from atom workersColors (core/mod/builtin/jaira-cockpit/hooks/register.tsx readWorkers); same test, colours unchanged after the 5 s re-read
- [x] Test in tests/ deckt Platzierung und Farbvergabe ab; Release-Notiz in core/release/notes/<ID>.md
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts describe 'the workers tab' (6 tests, claude plugin test: 149 pass rc=0); core/release/notes/F80QZT.md
- [x] Wechselt ein Ticket die Lane, läuft sein ✻ über mehrere Frames Station für Station hinüber (Timer per $.clock.every, gestoppt beim Schließen und Tabwechsel)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx showTab/frame/stopFrames + overview.ts stepOf; tests 'a ticket that changed lanes walks there one station at a time', 'the frames stop when another tab is shown and when the pane closes'
- [x] Drei Zustände lesbar: working (✻✳✶✢ dreht), idle (gedimmt, z z), Frage (?-Blase in Warnfarbe); Test deckt die Zuordnung ab
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx figure(); tests 'a working figure turns each frame; a resting one does not', 'every agent stands at the station…' (question bubble), 'a subagent of this session stands at its ticket too'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-16 12:50 · Alexander Sacharov** — parent session: jaira-28
- **2026-10-16 12:51 · Alexander Sacharov** — Offene Entscheidungen vor dem Plan: keine, die Alex treffen muss; alle aus Ticket/Code entschieden, siehe decision-Notizen.
- **2026-10-16 12:51 · Alexander Sacharov** — decision: Farbe je Ticket (Goal), nicht je Agent; mehrere Agenten eines Tickets teilen sie. DoD 2 gelesen als: zwei Tickets mit laufenden Agenten nie gleichfarbig, solange Tickets <= Farben.
- **2026-10-16 12:52 · Alexander Sacharov** — decision: Farbvergabe = Hash des Handles als Startplatz, bei Kollision die nächste freie Farbe; schon vergebene Farben laufender Tickets werden in einem cockpit-Atom gemerkt und behalten, auch wenn ein neuer Hash-Treffer dazukommt. Mehr Tickets als Farben: reiner Hash.
- **2026-10-16 12:52 · Alexander Sacharov** — decision: Spur zeigt alle agentischen Lanes in Board-Reihenfolge, leere als Platzhalter (DoD 1 sagt 'Spur der agentischen Lanes'). Steht im Queue-Tab ganz oben. Gezeichnet als gepaddete Textspalten, nicht als Box width=: der bestehende Test verbietet Boxen mit fester Breite.
- **2026-10-16 12:52 · Alexander Sacharov** — decision: 'laufend' = dieselbe Regel wie die Zahl im Band (snapshotOf): Herdr-Tab mit agent_status working, oder Subagent pending/running/waiting, der einen Handle nennt. Ein Worker, der idle auf eine Antwort wartet, erscheint nicht. Spur wird wie das übrige Overview bei jedem Druck auf Overview neu gelesen, nicht live.
- **2026-10-16 12:55 · Alexander Sacharov** — Designänderung von Alex, 2026-10-16 (über jaira-28): verspielt, wie ein kleines Spiel. Eigener dritter Tab 'Workers' statt Abschnitt im Queue-Tab; Lanes als Stationen (◯ lane ─ ◯ lane); jeder Worker ein ✻ in Ticketfarbe mit Handle an seiner Station; Lanewechsel = sichtbares Laufen über mehrere Frames; drei Zustände working (✻✳✶✢ dreht), idle (gedimmt, z z), Frage (? Blase, Warnfarbe). Ersetzt die decision-Notizen zu 'Queue-Tab ganz oben', 'gepaddete Textspalten' und 'idle erscheint nicht'.
- **2026-10-16 12:55 · Alexander Sacharov** — decision: Timer geht. Plugin-API hat $.clock.every(ms, fn) mit cancel(), und update(atom) zeichnet die Leser neu (plugin-authoring types, clock.every; register.tsx:949 nutzt es schon für das Band). Also echte Frames: ein Takt-Atom, 250 ms, läuft nur solange Overview offen und Tab Workers aktiv; ui.close für cockpit-overview und Tabwechsel stoppen ihn. Board wird in dem Tab alle 5 s neu gelesen, damit ein Lanewechsel ohne Neuöffnen sichtbar wird.
- **2026-10-16 12:55 · Alexander Sacharov** — decision: Stationen = alle agentischen Lanes plus die Lanes einer Person (human, signoff), in Board-Reihenfolge; backlog/todo/blocked/done nicht. Frage-Zustand: Ticket hat question gesetzt (auch ohne laufenden Agenten, z. B. in human), oder Subagent status 'waiting'. Idle: Herdr-Tab <id>/... offen, agent_status nicht working. Laufen: eine Station pro Takt-Schritt von der alten zur neuen Lane, Position je Ticket im Atom gemerkt.
- **2026-10-16 13:05 · Alexander Sacharov** — Gefunden: ein ui.close-Hook des Plugins hört das eigene $.ui.close nicht (im Test belegt, Timer lief weiter). Darum stoppt der Close-Knopf die Frames selbst; der Hook bleibt für das Schließen durch die Person (Escape/Close-Marke). Frage-Figur dreht sich nicht, nur working dreht. Agent eines Tickets in backlog/todo wird nicht gezeichnet (keine Station). Geprüft: claude plugin test 149 pass, plugin validate passed, tsc 5.6 strict ohne Fehler.
- **2026-10-16 13:06 · Alexander Sacharov** — critique (Runde 1): core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:21 PEOPLE ist dieselbe Liste wie TICKET_COLORS in overview.ts - eine Palette behalten, personColor nimmt sie
core/mod/builtin/jaira-cockpit/hooks/register.tsx readWorkers wiederholt das Lesen von Herdr-Tabs und $.agent.list aus refreshBoard (register.tsx ~252-260) - einen Helfer agentsOf($, dir) für beide
core/mod/builtin/jaira-cockpit/hooks/overview.ts Walker.kind liest niemand außer dem React-Key, den n schon eindeutig macht - Feld streichen, auch in types/index.d.ts
core/mod/builtin/jaira-cockpit/hooks/overview.ts colorsOf(palette=TICKET_COLORS): der Parameter hat keinen Aufrufer, der ihn setzt - streichen
- **2026-10-16 13:07 · Alexander Sacharov** — in-progress Runde 2: alle 4 critique-Funde umgesetzt - eine Palette COLORS (overview.ts) für Personen und Tickets, agentsOf($, dir) in register.tsx für refreshBoard und readWorkers, Walker.kind und colorsOf-palette gestrichen. Tests 149 pass, validate, tsc grün. Proofs geprüft: zeigen weiter auf existierende Namen.
- **2026-10-16 13:07 · Alexander Sacharov** — critique Runde 2: die 4 Funde aus Runde 1 sind in a37f79b5 umgesetzt (Palette, agentsOf, kind, palette-Param); nur dieser Diff gelesen, nichts Neues. none.
