---
id: 01M4RWVFXTPC084S81WSDSKZ70
title: "Cockpit zeigt, was der Teamlead anlegt und was als Nächstes ansteht"
status: testing
ready: true
creator: Alexander Sacharov
goal: "Ein Knopf im Cockpit-Band oeffnet eine Board-Geschichte ueber mehrere Tage: Spalten wie Lanes (in Arbeit, in der Warteschlange, zurueckgestellt) und darunter, was in den letzten Tagen angelegt und zurueckgestellt wurde - damit man nach einer Nacht nicht verliert, was der Teamlead angelegt hat"
context: |-
  Was heute fehlt: Der Teamlead legt oft Tickets an oder stellt sie zurück (backlog, blocked), und Alex sieht davon im Cockpit nichts.
  Das Band zeigt nur 'Waiting for you', Entscheidungen und die Tickets der Kollegen, also nichts Neues und nichts Kommendes.
  Anlass: Alex mag die Ticket-Tabelle, die der Teamlead im Chat schreibt (Handle | Lane | Titel), und will so eine Ansicht im Cockpit hinter einem Knopf.
  Daten gibt es schon: jaira list --json (Status, Erstellzeit), jaira next --all --json (was startbar ist), blocked = zurückgestellt.
  Ort: core/mod/builtin/jaira-cockpit/hooks/register.tsx (Band, Panes über $.ui.open wie der Message-Pane), Texte in hooks/texts.ts.
  Offen: welcher Zeitraum zählt als 'neu' (seit Sitzungsstart? letzte 24h?) — vor dem Bauen klären.
definition-of-done: "Ein Knopf im Band öffnet einen Pane mit drei Abschnitten: neu angelegt, zurückgestellt (blocked/backlog), als Nächstes startbar"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-12T20:28:49Z
updated-at: 2026-10-12T21:17:45Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
mode: ""
claimed-by: DESKTOP-RFTCH11-41026
claimed-at: 2026-10-12T20:34:22Z
outcome-what: "Ein Knopf Overview im Cockpit-Band oeffnet einen Pane: drei Spalten In work / Queue / Put back (Handle | lane | title, darunter wer bzw. Grund) und darunter, was in den letzten 7 Tagen angelegt oder zurueckgestellt wurde, nach Tag, mit Autor. Der Teamlead nennt dieselbe Geschichte beim Start."
outcome-why: "Alex verlor nach einer Nacht, was der Teamlead angelegt oder zurueckgestellt hatte; das Band zeigte nur, was auf ihn wartet. Die Ansicht liest nur das Board (jaira lanes/list/next --all), also ueberlebt sie Sitzungsende und Neustart."
outcome-resolves: DSKZ70
review-summary: none
review-gaps: "overview.ts nutzt parse/oneLine aus board.ts statt eigener Kopien (beide exportiert); Parameter days von overviewOf entfernt (kein Aufrufer mit anderem Wert, HISTORY_DAYS direkt); Table:any in overview-pane.tsx durch ElementTable aus claude-code ersetzt (tsc 5.6 mit den Engine-Typen gruen). Gelassen: Ticket-Typ in overview.ts neben dem in board.ts - andere Felder, ein gemeinsamer Typ waere eine Vereinigung ohne Gewinn; extra(r) zweimal pro Zeile - billig."
---

# Cockpit zeigt, was der Teamlead anlegt und was als Nächstes ansteht

## Definition of Done

- [x] Ein Knopf im Band öffnet einen Pane mit drei Abschnitten: neu angelegt, zurückgestellt (blocked/backlog), als Nächstes startbar
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:985 Knopf Overview im Band; :1103 Pane; Abschnitte in core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:30 (in Arbeit, Warteschlange, zurueckgestellt) + Geschichte (= neu angelegt); test 'the band has an Overview button only where there is a board'
- [x] Jede Zeile zeigt Handle | Lane | Titel, wie die Tabelle des Teamleads
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:22 Handle | Lane | Titel; test 'Overview reads the board afresh and shows in work, queue and put back, each row Handle | Lane | Titel'
- [x] Tests in core/mod/builtin/jaira-cockpit/tests decken Knopf und Pane ab; Zeile in den Release-Notes
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts:1544 describe 'the overview' (5 Tests: Knopf, kein Board, Spalten, Geschichte, Lesefehler); core/release/notes/DSKZ70.md:1
- [x] Pane-Aufbau: Spalten wie Lanes - 'in Arbeit' (Handle, Lane, Titel, Autor/Worker), 'Warteschlange' (startbar laut jaira next), 'zurueckgestellt' (backlog + blocked mit Grund); Test
  proof: core/mod/builtin/jaira-cockpit/hooks/overview.ts:68 overviewOf (in Arbeit/Warteschlange laut jaira next --all/zurueckgestellt mit blocked_reason); core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:30 drei Spalten nebeneinander; test 'Overview reads the board afresh and shows in work, queue and put back, each row Handle | Lane | Titel'
- [x] Geschichte: was in den letzten Tagen angelegt oder nach backlog/blocked gestellt wurde, nach Tag gruppiert, mit Autor; ueberlebt Sitzungsende, Neustart und Nacht (aus Ticketdaten/git, nicht aus Sitzungsspeicher); Zeitraum begruendet in einer Note; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/overview.ts:27 HISTORY_DAYS=7 (Grund in Note), :97 put back aus updated_at/updated_by, nur jaira list --json; test 'the history: created and put back over the last seven days, by day, newest first, with who'
- [x] jaira-teamlead SKILL.md: beim Start auch die zuletzt angelegten/zurueckgestellten Tickets der letzten Tage lesen und nennen (dieselbe Datenquelle wie der Pane)
  proof: core/role/builtin/jaira-teamlead/SKILL.md:26 (Abschnitt 'Read the board first': created/put back der letzten 7 Tage aus jaira list --json, Regel wie overview.ts:78)
- [x] Zeile in core/release/notes/DSKZ70.md
  proof: core/release/notes/DSKZ70.md:1-2

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-12 20:33 · Alexander Sacharov** — decision: Alex 2026-10-12 (Teamlead-Gespraech): (1) 'neu' ist kein Zeitfenster seit Sitzungsstart oder 24h, sondern eine Geschichte ueber mehrere Tage - 'ich bin morgens schlafen gegangen und weiss nicht mehr, was du angelegt hast'. Quelle: Erstellzeit + Lane-Wechsel aus den Tickets/ihrer git-Geschichte, kein neuer Speicher, wenn es ohne geht. (2) Alle Tickets des Boards, mit Autor. (3) Aufbau als mehrere Spalten wie Lanes: 'in Arbeit' als Liste (Ticket, Lane, wer), 'in der Warteschlange', 'zurueckgestellt' (backlog/blocked); dazu die Geschichte der letzten Tage, damit nichts von vor zwei Tagen verloren geht. (4) Auch der Teamlead selbst soll beim Start pruefen, was in den letzten Tagen in backlog gelandet ist.
- **2026-10-12 20:34 · Alexander Sacharov** — Teamlead 2026-10-12: parallel laufen 13PQB3 und 0EW5XK in derselben register.tsx (Band/Explain-Pane). Den neuen Pane in eine eigene Datei legen, im Band nur den einen Knopf ergaenzen - so bleibt der Merge klein.
- **2026-10-12 20:34 · Alexander Sacharov** — parent session: jaira-f5
- **2026-10-12 20:35 · Alexander Sacharov** — question (worker, 2026-10-12): Ветка feat/dskz70 стоит на 43080ed7 (локальный master устарел) — там нет core/mod/builtin/jaira-cockpit; он есть только в upstream/master (a41a70ab). 'git merge --ff-only upstream/master' на ветке без собственных коммитов отклонил классификатор прав. Варианты: 1) (Recommended) Alex сам выполняет в worktree: ! git -C /home/alex/projects/.worktrees/jaira-dskz70 merge --ff-only upstream/master — ветка остаётся та же, ничего не теряется; 2) разрешить мне сделать этот fast-forward; 3) пересоздать worktree от upstream/master другим способом.
- **2026-10-12 20:55 · Alexander Sacharov** — decision: Alex hat feat/dskz70 selbst auf upstream/master 9f6e9877 vorgespult; Basis-Frage erledigt.
- **2026-10-12 20:57 · Alexander Sacharov** — decision (worker): Zeitraum der Geschichte = 7 Tage, nach Tag gruppiert. Grund: Alex will 'nach einer Nacht' und 'von vor zwei Tagen' sehen; 7 Tage decken auch Freitag->Montag und einen Feiertag; ein voller Tag hier sind ~8 neue Tickets, also hoechstens ~50 Zeilen, der Pane scrollt. Quelle nur jaira list --json, kein neuer Speicher: 'angelegt' = created_at + creator (exakt); 'zurueckgestellt' = Ticket steht jetzt in backlog/blocked und updated_at liegt im Zeitraum, mit updated_by (blocked mit blocked_reason). Grenze: jaira speichert keinen Zeitpunkt des Lane-Wechsels, nur updated_at - eine spaetere Note auf einem Backlog-Ticket verschiebt seinen Tag. Exakt waere nur git log ueber die Ticket-Refs, zu langsam fuer einen Pane, der jede Minute neu laedt.
- **2026-10-12 20:57 · Alexander Sacharov** — decision (worker): Spalten. 'in Arbeit' = Tickets in einer agentischen Lane, die jaira next --all NICHT als startbar fuehrt (also jemand haelt sie), mit assignee. 'Warteschlange' = jaira next --all ohne backlog/blocked/Personen-Lanes (human, signoff stehen schon im Band). 'zurueckgestellt' = alle Tickets in backlog + blocked (blocked mit Grund), neueste zuerst. DoD 1 'neu angelegt' ist die Geschichte darunter (DoD 5), kein vierter Block. Datei: hooks/overview.ts (Daten, reine Funktionen) + Pane-Render in register.tsx nur als ein Knopf und ein kurzer Pane-Handler, der overview.ts aufruft.
- **2026-10-12 21:03 · Alexander Sacharov** — Gelernt: register.tsx darf $ nicht ueber einen Import weiterreichen (claude plugin validate: '$ is followed only into a function declared in this same file'). Daher: hooks/overview.ts (Daten, rein), hooks/overview-pane.tsx (nur Zeichnen, bekommt die Element-Tabelle), register.tsx nur Knopf + openOverview + ein Pane-Hook (~20 Zeilen). Fund beim Test: Ticket ohne created_at bekam Zeit 0 und landete in der Geschichte; unlesbare Zeit zaehlt jetzt als nie.
- **2026-10-12 21:05 · Alexander Sacharov** — decision: Alex 2026-10-12 (via jaira-f5): Schritt 1 so akzeptiert - 7-Tage-Fenster und updated_at fuer 'put back' sind in Ordnung.
- **2026-10-12 21:14 · Alexander Sacharov** — critique (pass 1): keine Befunde zur Form. Stehen gelassen: eigene Dateien overview.ts/overview-pane.tsx statt alles in register.tsx - Grund ist die Teamlead-Note vom 20:34 (13PQB3 in derselben Datei); die Teilung Daten (overview.ts, rein) / Zeichnen / $ in register.tsx folgt board.ts + register.tsx und der Regel der plugin validate. Fuer optimize uebergeben, nicht Kritik: overview.ts dupliziert parse/oneLine aus board.ts; overviewOf(…, days) hat keinen Aufrufer mit anderem Wert; Table in overview-pane.tsx ist 'any'. Fuer testing: ein beanspruchtes Ticket in todo (nicht-agentische Lane, nicht in next) faellt aus allen Spalten.
- **2026-10-12 21:16 · Alexander Sacharov** — optimize: Proofs DoD 1/2/4/5 nach dem Verschieben der Zeilen neu gesetzt. Typpruefung ohne Projekt-tsconfig: scratch tsconfig nach dem Header von claude-code.d.ts (Skill plugin-authoring) mit npx typescript@5.6.3; Gegenprobe mit absichtlichem Fehler meldet TS2339, prueft also wirklich.
- **2026-10-12 21:17 · Alexander Sacharov** — decision: Alex 2026-10-12 23:2x (geht schlafen): alle PRs vorbereiten, so dass sie nacheinander kommen. Dafuer conversational aus - der Worker committet selbst (auch den offenen optimize-Refactor). DSKZ70 kommt NACH dem gemeinsamen PR feat/13pqb3 (13PQB3+0EW5XK+5A6RJZ): feat/13pqb3 in feat/dskz70 mergen, sobald gepusht. PR fuer DSKZ70 als Draft, gestapelt; Abnahme durch Alex morgen.
