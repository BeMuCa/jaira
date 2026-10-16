---
id: 01M52C52KX820CPG5ZW8F80QZT
title: "Cockpit: Worker als farbige Claude-Icons auf einer Lane-Spur im Overview"
status: testing
ready: true
creator: Alexander Sacharov
goal: "Das Overview-Pane des cockpits hat einen Tab 'Workers', der das ganze Pane füllt: die Lanes als Stationen über die volle Breite, jeder Agent als Claude-Code-Pixelkrebs (Clawd) aus Blockzeichen in der Farbe seines Tickets, mit eigenem Hut und dem Handle darunter. Wechselt das Ticket die Lane, läuft der Krebs mit Laufzyklus hinüber; working bewegt die Arme, idle sitzt gedimmt mit z z, eine offene Frage steht als (?) über dem Hut."
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
updated-at: 2026-10-16T13:29:30Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-51728
claimed-at: 2026-10-16T12:50:47Z
outcome-what: "Workers-Tab neu: eine Spur pro Agent über die volle Pane-Höhe, links die Lane des Tickets, rechts seine next_lane; darin der Claude-Krebs aus Blockzeichen in Ticketfarbe mit eigenem Hut, so weit rechts wie Plan/DoD abgehakt sind, darunter 'Handle · working|resting|waits for an answer · Titel'; läuft mit Beinzyklus, beim Lanewechsel bis zur rechten Wand, dann rücken die Wände; readWorkers liest das Verzeichnis aus openOverview"
outcome-why: "Alex sah die Stationen-Version laufen und verstand nichts; er will auf einen Blick lesen, wer was tut und wie weit, ohne Legende"
outcome-resolves: F80QZT
review-summary: "none"
round: 3
review-gaps: "trackRows, WorkersScene (overview-pane.tsx) und WALK_STEP (overview.ts) nicht mehr exportiert - niemand importiert sie; stateWords-Umweg gestrichen, texts.ts hält states { working, idle, question } direkt. Gesucht nach Fortschritts-Zählung (superseded/plan_items/dod_items) in hooks/: progressOf ist die einzige. Gelassen: Lauf-Runs-Zusammenfassung im Zeichnen, nur dort gebraucht. Proofs geprüft: trackRows, pickOf, progressOf, stepOf bestehen weiter."
test-verdict: "pass: go build ./... rc=0, go test -race -count=1 ./... rc=0, claude plugin test 149 pass rc=0, plugin validate passed, tsc 5.6 strict rc=0; DoD 1-6 je über den genannten Test geprüft; Workers-Tab per Scratch-Test über 4 Frames gerendert: Spinner ✻→✳→✶→✢, idle 'z z', ?-Blase in warning über K7P2QX in human, T8QS1Y läuft in-progress → critique (»)→ review. Nicht geprüft: das echte TUI per --plugin-dir (Bildschirm nicht headless steuerbar) und ob ui.close beim Schließen durch die Person (Escape) feuert"
review-verdict: "Diff erfüllt alle 6 DoD-Punkte; claude plugin test 149 pass selbst nachgelaufen. Keine Defekte, zwei kleine Unsicherheiten (Subagent-Status idle, board?.dir statt dir) und Escape-Stopp nur live prüfbar. Annehmen nach Live-Blick."
review-round: 2
review-check: "1. claude --plugin-dir /home/alex/projects/.worktrees/jaira-cockpit-worker-icons/core/mod/builtin/jaira-cockpit im Repo jaira starten  2. im cockpit-Band auf Overview drücken, dann auf den Tab Workers  3. sehen: Zeile ◯ brainstorm ─ ◯ pre-process ─ ◯ in-progress … ◯ human … ◯ signoff, darunter Legende '✻ at work · z z resting · ? waits for a person'  4. in einem zweiten Terminal einen Agenten an ein Ticket setzen (z. B. /jaira-role-lane <id> critique): über seiner Lane steht ein ✻ <ID>, das sich dreht (✻✳✶✢), in eigener Farbe  5. jaira move <id> --to <spätere Lane>: binnen ~5 s läuft das ✻ Station für Station hinüber, mit » daneben  6. jaira set <id> question=\"x\": über der Figur erscheint ein gelbes ?  7. Tab Queue drücken, dann Escape: Pane zu, kein Flackern mehr im Terminal (Takt gestoppt). Automatisch: cd core/mod/builtin/jaira-cockpit && claude plugin test . → describe 'the workers tab', 149 pass"
---

# Cockpit: Worker als farbige Claude-Icons auf einer Lane-Spur im Overview

## Definition of Done

- [x] Workers-Tab füllt das Pane (bodyColumns × bodyRows): eine Spur pro Figur, untereinander, durch eine dünne Linie getrennt; links die beschriftete Wand der aktuellen Lane, rechts die der next_lane; darin ein Clawd-Sprite aus Blockzeichen mit Hut
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx trackRows; core/mod/builtin/jaira-cockpit/tests/board.test.ts 'the workers tab > one track per figure fills the pane, walled by its lane and the next, with who does what in plain words'
- [x] Jede Figur trägt den 6-stelligen Handle; zwei Tickets mit Figuren haben nie dieselbe Farbe und nie denselben Hut, solange es nicht mehr Tickets als Farben bzw. Hüte gibt
  proof: core/mod/builtin/jaira-cockpit/hooks/overview.ts pickOf (COLORS, HATS); core/mod/builtin/jaira-cockpit/tests/board.test.ts 'the workers tab > two tickets never wear the same colour or hat, and each keeps both on every read'
- [x] Farbe und Hut eines Tickets bleiben zwischen zwei Aktualisierungen gleich
  proof: pickOf keeps picks from atoms workersColors/workersHats (core/mod/builtin/jaira-cockpit/hooks/register.tsx readWorkers); same test after the 5 s re-read
- [x] Tests in tests/ decken Platzierung, Hut- und Farbvergabe, Fortschritt→Spalte, Laufzyklus und das Rücken der Wände ab; Release-Notiz in core/release/notes/<ID>.md
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts describe 'the workers tab' (9 tests; claude plugin test 152 pass rc=0); core/release/notes/F80QZT.md
- [x] Die Spalte des Krebses folgt dem Fortschritt in der Lane (Plan-Häkchen/Plan, sonst DoD, sonst linke Wand) und wird Spalte für Spalte angelaufen, die Beinzeile wechselt dabei zwischen zwei Frames; beim Lanewechsel läuft er zur rechten Wand, die Wände rücken (alte rechte = neue linke), er startet links neu (Timer per $.clock.every, gestoppt beim Schließen und Tabwechsel)
  proof: core/mod/builtin/jaira-cockpit/hooks/overview.ts progressOf, stepOf; core/mod/builtin/jaira-cockpit/tests/board.test.ts 'the workers tab > the crab stands as far across as its ticket has come…', '…walks there column by column, its legs stepping', '…the walls move along, and it starts again at the left', '…the frames stop when another tab is shown and when the pane closes'
- [x] Drei Zustände am Sprite: working (Arme bewegen sich, Funke), idle (sitzt, gedimmt, z z über dem Hut), Frage ((?) über dem Hut in Warnfarbe); unter dem Krebs eine Zeile 'Handle · working|resting|waits for an answer · Titel gekürzt'; Test deckt die Zuordnung ab
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx trackRows (ARMS, LEGS.sit, '(?)', 'z z', caption); core/mod/builtin/jaira-cockpit/tests/board.test.ts 'the workers tab > working moves its arms in place; resting sits still', '…one track per figure…' (captions, z z, (?))
- [x] readWorkers liest das Verzeichnis, das openOverview bekommt, nicht board?.dir (review-Lücke 1)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx openOverview sets overviewDir, readWorkers reads it

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
- **2026-10-16 13:11 · Alexander Sacharov** — testing: grün, siehe test-verdict. Offen für den Menschen: Live-Blick mit claude --plugin-dir <worktree>/core/mod/builtin/jaira-cockpit, Escape-Schließen stoppt Frames (nicht testbar im Kit).
- **2026-10-16 13:15 · Alexander Sacharov** — review (fresh): 149 pass nachgelaufen, DoD 1-6 vom Diff gedeckt. Offen, klein: readWorkers nutzt board?.dir statt openOverview-dir; Subagent-Status 'idle' nur im Test belegt; Escape-Stopp nur live prüfbar.
- **2026-10-16 13:16 · Alexander Sacharov** — worker zu review-gaps 2: AgentStatus 'idle' existiert in der Plugin-API (claude-code.d.ts: 'pending' | 'running' | 'waiting' | 'idle' | 'completed' | 'failed' | 'killed', idle = zwischen zwei Turns) - geschlossen. Gap 1 (readWorkers liest board?.dir statt dir aus openOverview) bleibt offen; im signoff nicht mehr meins.
- **2026-10-16 13:19 · Alexander Sacharov** — Alex' Entscheidung 2026-10-16 (über jaira-28) nach Blick in signoff: zurück nach in-progress. Statt ✻-Glyphen der Claude-Code-Pixelkrebs (Clawd) aus Blockzeichen in Ticketfarbe, Handle darunter; jeder Worker mit eigenem Hut (1 Zeile, kollisionsfrei wie die Farbe); Workers-Tab füllt das Pane 100% Breite und Höhe, Stationen über die ganze Breite, beim Lanewechsel echter Laufzyklus (Beinzeile wechselt zwischen 2 Frames, 1-2 Spalten pro Frame); working = Arme/Beine bewegen sich am Platz, idle = sitzt still, gedimmt, 'z z' über dem Hut, Frage = ?-Sprechblase über dem Hut in Warnfarbe; vorher klären, wie das Pane seine Größe erfährt; review-Lücke readWorkers liest dir aus openOverview statt board?.dir.
- **2026-10-16 13:21 · Alexander Sacharov** — decision: Größe kommt aus den Pane-Props bei jedem ui.render: bodyColumns (Zellen innen) und scroll.bodyRows (die Zeilen, die das Pane höchstens haben darf; Dock = volle Höhe, inline passt sich dem Baum an, also ist bodyRows dort das Maximum). Größenänderung zeichnet neu mit neuen Props. 100% Höhe geht also: Holst = bodyRows minus Tab-, Legenden- und Close-Zeile.
- **2026-10-16 13:21 · Alexander Sacharov** — decision: Workers-Tab wird ein Zeichen-Holst (Zeilen aus Zellen), nicht Box-Layout, damit Figuren spaltengenau laufen. Position je Ticket in Stationseinheiten (Bruch), 1/8 Station pro 250-ms-Frame, so bleibt sie bei Größenänderung richtig. Hut wie Farbe pro Ticket (nicht pro Agent): Hash-Start, nächster freier; zwei Agenten desselben Tickets sehen gleich aus und tragen denselben Handle. Mehrere Figuren an einer Station stapeln sich nach oben, ohne Platz eine Spalte nach rechts. working = Arme heben/senken + ✻-Funke rechts, idle = sitzt (Beine eingezogen), gedimmt, 'z z' über dem Hut, Frage = (?) über dem Hut in warning.
- **2026-10-16 13:22 · Alexander Sacharov** — Alex' Entscheidung 2026-10-16 (über jaira-28), ersetzt 'Stationen über die Breite': er sah es laufen und verstand nichts. Neu: eine Spur (Zeile) pro Worker, untereinander, durch eine dünne Linie getrennt, zusammen über die volle Pane-Höhe; links die Wand der Lane, in der das Ticket steht, rechts die Wand seiner next_lane, beide beschriftet; x = Fortschritt in der Lane aus dem Ticket (Plan-Häkchen/Plan gesamt, sonst DoD, sonst linke Wand); beim Lanewechsel verschieben sich die Wände (alte rechte wird neue linke), der Krebs fängt links neu an; unter dem Krebs eine Zeile in Klartext: Handle · Zustand · Titel gekürzt. Hüte, Farben, Laufzyklus, Zustände, volle Größe, readWorkers-dir bleiben.
- **2026-10-16 13:22 · Alexander Sacharov** — decision: Zustandswörter auf Englisch (working / resting / waits for an answer), wie alle cockpit-Texte in texts.ts - Projektsprache Englisch; jaira-28 nannte russische Beispiele nur für die Bedeutung. Fortschritt: abgelöste Punkte ([-]) zählen weder als erledigt noch zur Summe. In Lanes nach in-progress ist die DoD meist voll, der Krebs steht dort an der rechten Wand = 'bereit für die nächste Lane'. Lanewechsel: Krebs läuft erst bis zur rechten Wand, dann springen die Wände um und er startet links, dann läuft er zu seinem Fortschritt. Geschwindigkeit 1/40 der Spur pro 250-ms-Frame. Eine Spur pro Figur (Agent oder Ticket mit offener Frage); Position, Farbe und Hut je Ticket. Legende entfällt, die Bildunterschrift sagt es.
- **2026-10-16 13:28 · Alexander Sacharov** — in-progress Runde 3 (Spuren-Layout): Holst aus Zellen in overview-pane.tsx trackRows, Modell in overview.ts (workersOf -> Walker[] mit lane/next/title/progress, progressOf, pickOf für Farbe+Hut, stepOf mit Spot {lane,x}). Spurhöhe = (bodyRows-2 - Regeln)/n, mindestens 7 Zeilen; mehr Figuren als Platz -> Pane scrollt. Fließkomma: beim Hinauslaufen wird x knapp unter 1 auf 1 gesetzt, sonst käme der Wandwechsel nie. Tests 152 pass, validate, tsc grün.
- **2026-10-16 13:28 · Alexander Sacharov** — critique Runde 3 (Spuren): ganzer Diff 0763467b gelesen. Zellen-Holst in trackRows ist die einfachste Form, die spaltengenaues Laufen erlaubt (Box-Layout kann keine Bruchteil-Spalten); pickOf ersetzt colorsOf statt daneben zu stehen; Modell bleibt in overview.ts, Zeichnung in overview-pane.tsx wie bisher. none. Für optimize: trackRows und WALK_STEP sind exportiert, ohne dass jemand sie importiert.
