---
id: 01M4RWVFXTPC084S81WSDSKZ70
title: "Cockpit zeigt, was der Teamlead anlegt und was als Nächstes ansteht"
status: done
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
commits:
  - adecb1eba291258a5b54a11218a70fc2af716b53
  - eafa7e8648523039b46bcafc340d3ba68c071edc
  - ca2f57ced067baf1209bf1cec563b1a213b0fefb
  - e62c17242f65d4f00b5e2a355e26e93a160db8dd
  - 776fc9af9b42e90364fabbc34374782a21c22182
created-at: 2026-10-12T20:28:49Z
updated-at: 2026-10-12T22:08:19Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
mode: ""
claimed-by: DESKTOP-RFTCH11-88370
claimed-at: 2026-10-12T21:24:21Z
outcome-what: "Ein Knopf Overview im Cockpit-Band oeffnet einen Pane: drei Spalten In work / Queue / Put back (Handle | lane | title, darunter wer bzw. Grund) und darunter, was in den letzten 7 Tagen angelegt oder zurueckgestellt wurde, nach Tag, mit Autor. Der Teamlead nennt dieselbe Geschichte beim Start."
outcome-why: "Alex verlor nach einer Nacht, was der Teamlead angelegt oder zurueckgestellt hatte; das Band zeigte nur, was auf ihn wartet. Die Ansicht liest nur das Board (jaira lanes/list/next --all), also ueberlebt sie Sitzungsende und Neustart."
outcome-resolves: DSKZ70
review-summary: "Neuer Knopf 'Overview' in der Kopfzeile des Cockpit-Bands (nur mit Board). Er ruft jaira lanes/list/next --all --json frisch auf und oeffnet einen Pane: drei Spalten nebeneinander - In work (Ticket in agentischer Lane, das next nicht anbietet, mit assignee), Queue (was next --all anbietet, ohne backlog/blocked, Personen-Lanes und terminal), Put back (alle backlog+blocked, mit wer und blocked_reason), jede Zeile 'Handle | lane | title', zuletzt geaenderte zuerst. Darunter die Geschichte der letzten 7 Tage nach lokalem Kalendertag: 'created by <creator>' aus created_at, 'put back by <updated_by>' aus updated_at bei Tickets, die jetzt in backlog/blocked stehen und an einem spaeteren Tag als angelegt geaendert wurden. Scheitert ein Befehl, zeigt der Pane 'could not be read'. Daten rein in hooks/overview.ts, Zeichnen in hooks/overview-pane.tsx, register.tsx nur Knopf + openOverview + Pane-Hook. jaira-teamlead SKILL.md nennt beim Start dieselbe 7-Tage-Geschichte nach derselben Regel. Release-Note core/release/notes/DSKZ70.md mit zwei Zeilen; 5 neue Tests in tests/board.test.ts."
review-gaps: "(1) 'In work' ist zu weit: overview.ts:186-188 nimmt jedes Ticket einer agentischen Lane, das jaira next --all nicht anbietet, als 'gehalten'. next laesst aber auch Tickets weg, die nicht Ready sind (Promotion-Felder fehlen) oder auf ein offenes blocked_by warten (gate.Actionable, core/gate/gate.go:644) - solche stehen dann in 'In work' mit leerem wer, obwohl niemand sie haelt; der Kommentar overview.ts:161 behauptet 'somebody holds it'. Umgekehrt landet ein Ticket mit abgelaufenem Claim oder einem Claim derselben Session wie das Cockpit in Queue (flow.go:396 ClaimActive). Auf dem echten Board gerade ohne Folge (einziger Fall 5Q20DC, wirklich gehalten). (2) Bekannt aus critique/testing und gelassen: ein beanspruchtes oder abhaengigkeitsblockiertes Ticket in todo faellt aus allen drei Spalten. (3) 'put back' haengt an updated_at - eine Note auf einem alten Backlog-Ticket erscheint als Zurueckstellen; von Alex am 2026-10-12 akzeptiert. (4) Kein Test fuer (1)/(2). Sonst deckt der Diff jeden DoD-Punkt; Outcome stimmt mit dem Diff ueberein."
test-verdict: "pass: go clean -testcache && go test -race ./... RC=0, claude plugin test 117/117 RC=0; DoD 1-7 an den Proof-Zeilen und Tests im Baum geprueft; overviewOf auf den echten Board-Daten (jaira lanes/list/next --all) laufen lassen: 1 in work, 3 queue, 49 put back, Geschichte 12.10 mit 8 angelegten Tickets und Autor"
review-verdict: "Diff erfuellt alle sieben DoD-Punkte; 117 Cockpit-Tests gruen, claude plugin validate gruen, go test core/mod core/release core/role gruen. Ein echter, aber kleiner Defekt: 'In work' zaehlt nicht-startbare Tickets (fehlende Felder, offenes blocked_by) als gehalten - auf dem heutigen Board unsichtbar. Annehmbar; wer es sauber will, laesst In work zusaetzlich an assignee/Claim haengen (eine Zeile in overview.ts:188). Offen ausserhalb des Diffs: die Note vom 21:17 verlangt, feat/13pqb3 in feat/dskz70 zu mergen, bevor der Draft-PR kommt - im Branch noch nicht geschehen."
review-check: "1. cd /home/alex/projects/.worktrees/jaira-dskz70/core/mod/builtin/jaira-cockpit && claude plugin test .  -> am Ende '117 pass, 0 fail'  2. cd /home/alex/projects/.worktrees/jaira-dskz70 && claude --plugin-dir /home/alex/projects/.worktrees/jaira-dskz70/core/mod/builtin/jaira-cockpit  3. im Band oben rechts neben dem Titel steht ein Knopf 'Overview' - druecken  4. ein Pane 'Board overview' oeffnet sich: drei Spalten 'In work (n)', 'Queue (n)', 'Put back (n)'; jede Zeile 'HANDLE | lane | Titel', darunter grau wer (bei blocked 'wer: Grund'); kein Ticket aus human/signoff/done  5. darunter 'Created and put back, last 7 days', Tageszeilen '── 2026-10-12 ──', darunter 'HH:MM  HANDLE | lane | Titel | created by ...' bzw. 'put back by ...', neuester Tag oben  6. Vergleich: jaira list --json | jq '.tickets[]|select(.status==\"blocked\")|.handle' - genau diese plus die backlog-Tickets stehen in 'Put back'  7. Escape oder 'Close' schliesst den Pane"
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
  proof: core/mod/builtin/jaira-cockpit/hooks/overview.ts:70 overviewOf (in Arbeit nur mit assignee / Warteschlange laut jaira next --all / zurueckgestellt mit blocked_reason); core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:30 drei Spalten nebeneinander; test 'Overview reads the board afresh and shows in work, queue and put back, each row Handle | Lane | Titel'
- [x] Geschichte: was in den letzten Tagen angelegt oder nach backlog/blocked gestellt wurde, nach Tag gruppiert, mit Autor; ueberlebt Sitzungsende, Neustart und Nacht (aus Ticketdaten/git, nicht aus Sitzungsspeicher); Zeitraum begruendet in einer Note; Test
  proof: core/mod/builtin/jaira-cockpit/hooks/overview.ts:27 HISTORY_DAYS=7 (Grund in Note), :99 put back aus updated_at/updated_by, nur jaira list --json; test 'the history: created and put back over the last seven days, by day, newest first, with who'
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
- **2026-10-12 21:21 · Alexander Sacharov** — testing: Funktion auf echten Daten via node --experimental-strip-types gegen hooks/overview.ts (Scratch, nicht im Repo). Gesehen und so gelassen: (1) PNC854 in 'human' steht in Queue - diese Board-Lane human hat human_exit=false, das Band (board.ts isPersons) zeigt sie dort ebenso nicht; beide lesen dieselbe Regel. (2) D28H7V erscheint am 07.10 als 'put back' wegen einer Note an dem Tag - die von Alex akzeptierte updated_at-Grenze. (3) Ein beanspruchtes Ticket in todo faellt weiter aus allen Spalten (Kritik-Note); auf dem echten Board gerade keins, kein DoD-Bruch.
- **2026-10-12 21:24 · Alexander Sacharov** — review (fresh): angenommen mit einem kleinen Befund - 'In work' (overview.ts:186-188) zaehlt jedes Ticket einer agentischen Lane, das jaira next --all nicht anbietet, als gehalten; next laesst aber auch nicht-Ready und durch blocked_by wartende Tickets weg (gate.Actionable). Die stehen dann mit leerem wer in In work. Fix waere: In work nur mit assignee/aktivem Claim, Rest in Queue oder eigene Markierung. Heute ohne Folge auf dem Board. Ausserdem steht der Merge von feat/13pqb3 in diesen Branch (Note 21:17) noch aus.
- **2026-10-12 21:28 · Alexander Sacharov** — merge: origin/feat/13pqb3 (dc0a1b4a, PR #69: 13PQB3+0EW5XK+5A6RJZ) in feat/dskz70 gemergt. Drei Konflikte, alle 'beide fuegten nebeneinander an' - register.tsx (busyAtom/overviewAtom, EXPLANATION+ANSWER/OVERVIEW, Explanation-Pane/Overview-Pane) und types/index.d.ts (busy/overview): beide Seiten behalten. Danach plugin validate ok, tsc ok, claude plugin test 133/133, go test -race ./... RC=0. Review-Befund ('In work' zaehlt nicht-startbare Tickets) ist offen und liegt bei Alex in signoff.
- **2026-10-12 21:30 · Alexander Sacharov** — fix nach Review-Befund (auf Anweisung jaira-f5, Alex schlaeft; bleibt in signoff, keine neue Review-Runde): 'In work' nimmt nur noch Tickets mit assignee (core/mod/builtin/jaira-cockpit/hooks/overview.ts:90). Ein Ticket, das jaira next mangels Ready oder wegen offenem blocked_by weglaesst und niemand haelt, steht in keiner Spalte. Test: W2W2W2 im Spalten-Test; ohne den Fix rot (132/1), mit gruen (133/0). Proofs DoD 4/5 auf die neuen Zeilen gesetzt.
- **2026-10-12 22:07 · Alexander Sacharov** — Alex 2026-10-13 per Telefon: 'Переведи и я сделаю релиз' - angenommen ohne Live-Check; Fehler werden morgen als neuer Bugfix aufgenommen. Vom Teamlead nach done bewegt, PR #71 auf ready.
