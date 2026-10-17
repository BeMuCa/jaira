---
id: 01M55C9V60QFM9V58MD8PXYN99
title: "Der Launcher zaehlt auch die Tickets, die nur auf ihrem Ref liegen"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die Zeile eines Projekts im Launcher (Startbild 'Projects:') zeigt dieselbe Zahl offener und gesamter Tickets wie 'jaira list' in diesem Repo - Tickets auf refs/jaira/tickets/<ID> eingeschlossen"
context: |-
  Was falsch ist: der Launcher zaehlt nur Ticket-Dateien auf der Platte (.jaira/tickets/). Tickets, die nur auf ihrem Ref liegen, fehlen.
  Belege 17.10.: jaira zeigt '29 open / 29', 'jaira list --json' liefert 75 Tickets, davon 46 nur auf dem Ref. demo-video zeigt '0 open / 0' bei 13 Refs unter refs/jaira/tickets und 0 Dateien.
  Seit Tickets bis zum Logbook auf ihrem Ref leben, ist 'nur auf dem Ref' der Normalfall - der Launcher zeigt also fast nichts.
  Code: internal/tui/home.go, func describe (~Zeile 160): ticket.At(root) + s.List() liest nur die Platte; e.Total = len(tickets), e.Open zaehlt nicht-terminale Lanes. 'jaira list' fuehrt Platte und Refs schon zusammen (Feld on-ref-only im JSON) - denselben Weg nehmen, nicht neu bauen.
  Achtung Startzeit: der Launcher liest viele Repos auf einmal; Refs lesen darf das Startbild nicht spuerbar bremsen (Constraint 'Instant startup' in CLAUDE.md). Ein Ticket, das auf Platte und Ref liegt, zaehlt einmal.
  Mehrere Boards (S0EG0M): describe prueft die Lane gegen die Lanes des Standard-Boards. Ein Ticket mit board: wird gegen die Lanes seines Boards gezaehlt, sobald S0EG0M auf master ist.
definition-of-done: "describe zaehlt Tickets auf der Platte und auf refs/jaira/tickets/ zusammen, jedes Ticket genau einmal; Zahlen gleich wie 'jaira list --json' im selben Repo; Test"
tags:
  - tui
blocked-by: []
related: []
commits:
  - a51bf3026a48d70037745ab646694789790f0ceb
  - 3385ffcac3e0c91cb256dce48f2887d28c737dcb
  - 35ad5dcec4499746fec0e6da2ba2b38f6c55c29b
  - 65bd1ede07900f56c52e8bf578933412aa1f0827
  - ded927e367ad5bbff29682c02ee7b50f739477cb
created-at: 2026-10-17T16:49:42Z
updated-at: 2026-10-17T19:01:25Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-54481
claimed-at: 2026-10-17T16:55:23Z
outcome-what: "Der Launcher liest jedes Repo ueber denselben refsync-Source wie 'jaira list' (Tickets nur auf dem Ref zaehlen mit, jedes einmal), prueft ein Ticket mit board: gegen die Lanes seines Boards, und zaehlt im Hintergrund: das erste Bild zeigt Name, Reihenfolge, Fehler und '… open', die Zahlen kommen per tea.Cmd nach."
outcome-why: "Seit Tickets bis zum Logbook auf ihrem Ref leben, zeigte der Launcher fast nichts (demo-video 0/0 bei 13 Refs); Refs lesen kostet ~80 ms, die das erste Bild nicht tragen darf (Constraint 'Instant startup', Entscheidung Teamlead 17.10.)."
outcome-resolves: "Launcher-Zeile = 'jaira list'-Zahlen; erstes Bild gleich schnell (Median ~55 ms), Zaehler stehen nach ~0,25 s."
review-summary: "Der Launcher liest jedes Repo jetzt ueber den refsync-Source (newSyncer als Store.Source, nur lokale Refs, kein Netz), also zaehlt s.List() auch Tickets, die nur auf refs/jaira/tickets/ liegen, dedupliziert gegen Dateien und Logbook. Ein Ticket mit board: wird gegen die Lanes seines Boards (lane.LanesOf) auf terminal/human geprueft, today's board weiter gegen den Katalog. Das Zaehlen laeuft parallel (ein Goroutine pro Repo, describeAll) und im Hintergrund: NewHome/refresh liefern sofort Name, Fehler, Agents und '… open', countBoards() als tea.Cmd (Init, r, a, Entfernen) liefert homeCountedMsg, das per Root einsortiert wird. Drei neue Tests in internal/tui/home_refs_test.go; release note core/release/notes/PXYN99.md."
round: 3
review-gaps: "1) Release-Note Zeile 1 sagt 'each row shows the same open and total as jaira list' - stimmt nur ohne benannte Boards. jaira list filtert per onBoard (internal/cli/root.go:378) auf today's board; der Launcher summiert alle Boards (bewusst, Notiz 16:56 'Launcher-Zeile = ganzes Repo'). Auf einem Repo mit board:-Tickets weichen die Zahlen ab, auch vom DoD-Wortlaut 'Zahlen gleich wie jaira list --json'. Live-Belege (demo-video 13, jaira 73) gelten, weil dort keine board:-Tickets liegen. Vorschlag: Note-Zeile auf 'same as jaira list summed over jaira boards' aendern. 2) Kleines Rennen: zwei countBoards-Laeufe (z.B. r zweimal schnell) - kommt der aeltere spaeter an, ueberschreibt counted() die frischeren Zahlen per Root. Selten, harmlos bis zum naechsten r. 3) today's board zaehlt weiter gegen den Katalog lane.Load(\"\") statt die eigenen .jaira/lanes des Repos - vorher schon so, nicht von diesem Diff."
test-verdict: "pass: go test ./... -race RC=0 (vor Commit 65bd1ede, Code seitdem unveraendert); DoD 1-5 im Baum geprueft, die drei Proof-Tests gruen; Binary im pty auf echten 10 Repos: erstes Bild '… open' ueberall, 'unreadable' sofort bei requirementsgenie, zweites Bild demo-video 4 open / 13, jaira 69 open / 73"
review-verdict: "Diff erfuellt den Kern: Ref-Tickets zaehlen mit, jedes einmal, ohne Netz, erstes Bild nicht langsamer; go vet sauber, go test -race ./internal/tui gruen. Abnahme empfohlen mit einer Korrektur an der Release-Note (Gap 1): 'gleich wie jaira list' gilt nur fuer Repos ohne benannte Boards. Gaps 2/3 sind klein bzw. nicht neu."
review-round: 3
review-check: "1. cd /home/alex/projects/.worktrees/jaira-pxyn99 && go build -o /tmp/jaira-pxyn99 ./cmd/jaira  2. /tmp/jaira-pxyn99 (ohne Argumente, Launcher 'Projects:')  3. Erstes Bild: jede Zeile zeigt kurz '… open', nach ~0,3 s Zahlen  4. Zeile demo-video zeigt 13 offene/gesamte Zahlen statt '0 open / 0'  5. In zweitem Terminal: cd ~/projects/demo-video && jaira list --json | grep '\"count\"' - gleiche Zahl als Total  6. Im Launcher r druecken: wieder kurz '… open', dann dieselben Zahlen  7. Benannte Boards: go test ./internal/tui -run 'TestHome' -v - TestHomeCountsATicketAgainstItsOwnBoard PASS"
---

# Der Launcher zaehlt auch die Tickets, die nur auf ihrem Ref liegen

## Definition of Done

- [x] describe zaehlt Tickets auf der Platte und auf refs/jaira/tickets/ zusammen, jedes Ticket genau einmal; Zahlen gleich wie 'jaira list --json' im selben Repo; Test
  proof: internal/tui/home.go:223; TestHomeCountsTicketsThatLiveOnlyOnTheirRef (internal/tui/home_refs_test.go:58); live: demo-video 13, jaira 73, genie-world 7 = jaira list --json
- [x] Ein Repo ohne Ticket-Dateien, aber mit Refs (wie demo-video) zeigt seine offenen Tickets statt '0 open / 0'; Test
  proof: TestHomeCountsTicketsThatLiveOnlyOnTheirRef (home_refs_test.go:58, Refs ohne Dateien -> 1 open / 2); live demo-video 4/13 statt 0/0
- [x] Startbild mit 10 Repos baut sich nicht spuerbar langsamer auf als vorher; gemessene Zeit vor/nach als proof
  proof: countBoards internal/tui/home.go:181, Platzhalter home.go:460, TestHomeShowsAPlaceholderUntilItHasCounted (home_refs_test.go:150). 10 echte Repos, Aufwaermlauf + 15 Laeufe, alt/neu abwechselnd, 2 Runden: erstes Bild alt Median 55/57 ms, neu 55/54 ms; Zaehler stehen neu nach Median 261/231 ms
- [x] Eine Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/PXYN99.md (2 Zeilen)
- [x] Ein Ticket mit board: zaehlt gegen die Lanes seines Boards (terminal/human dort), nicht gegen die des Standard-Boards; Test
  proof: internal/tui/home.go:250; TestHomeCountsATicketAgainstItsOwnBoard (home_refs_test.go:105)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-17 16:55 · Alexander Sacharov** — parent session: jaira-79
- **2026-10-17 16:56 · Alexander Sacharov** — decision: offene Entscheidungen keine. (1) Refs ueber denselben Weg wie 'jaira list': refsync.New(...) als Store.Source, dann liefert s.List() die Ref-Tickets mit (ReadOnly), dedupliziert gegen Dateien und Logbook. (2) Mehrere Boards: Launcher-Zeile = ganzes Repo, also alle Boards zusammen (Summe von 'jaira boards'); ein Ticket mit board: wird gegen lane.LoadBoard(root, board) geprueft, Standard-Board weiter gegen den Katalog wie bisher (Load(root) koennte das Board anlegen - der Launcher schreibt nicht). In einem Repo mit nur einem Board ist das gleich 'jaira list --json'. (3) Startzeit: Refs kosten ~4 git-Aufrufe pro Repo; describe laeuft deshalb parallel pro Repo.
- **2026-10-17 16:57 · Alexander Sacharov** — Messung vorher (NewHome, echte 10 Repos aus ~/.jaira/projects.json, 5 Laeufe, WSL): 120/123/186/165/125 ms, Median 124 ms; demo-video 0/0, jaira 28/28. Gemessen mit einem Wegwerf-Test (internal/tui/zz_timing_test.go, nicht committet) mit JAIRA_HOME=~/.jaira.
- **2026-10-17 17:01 · Alexander Sacharov** — critique: internal/tui/home.go:188 baut refsync.New+RemoteFor von Hand; internal/tui/refs.go:41 newSyncer tut genau das im selben Paket - newSyncer(s, "") nehmen
internal/tui/home.go:204 baut einen eigenen Cache Board->Lanes; core/lane/boards.go:162 lane.LanesOf tut das schon - fuer benannte Boards LanesOf(root) nehmen, Standard-Board und nil weiter auf den Katalog
- **2026-10-17 17:02 · Alexander Sacharov** — critique Runde 2: beide Befunde behoben (newSyncer, lane.LanesOf); nichts weiter. Ungebremste Goroutinen in describeAll bewusst stehen gelassen: eine pro registriertem Repo, ~10.
- **2026-10-17 17:08 · Alexander Sacharov** — testing: DoD 3 war falsch belegt. Die erste Vorher-Messung (124 ms) enthielt den kalten ersten Lauf. Fair gemessen (NewHome, echte 10 Repos, 1 Aufwaermlauf + 15 Laeufe, alt/neu abwechselnd, 2 Runden): alt min 38 / Median 53 bzw. 43 ms, neu min 102-107 / Median 124 bzw. 132 ms. Ursache requirementsgenie: 159 Refs, describe ~99 ms, davon git cat-file --batch ~50 ms und ParseDoc+Decode der 159 Tickets ~85 ms (Disk-List 35 ms). Parallelisierung hilft nur bis zum langsamsten Repo. Ohne Entscheidung geht es nicht weiter - siehe Frage.
- **2026-10-17 17:08 · Alexander Sacharov** — Frage (DoD 3, Startzeit): Refs lesen kostet beim Start ~80 ms (10 Repos, warm: Median 48 -> 128 ms; requirementsgenie mit 159 Refs allein ~100 ms). Wie weiter?
1. (Empfohlen) Asynchron: Startbild sofort mit den Zahlen von der Platte (so schnell wie vorher), die Ref-Zahlen kommen ~0,1 s spaeter per tea.Cmd und ersetzen sie. Kosten: Zahlen springen beim Start einmal; ca. 40 Zeilen.
2. Hinnehmen: Startbild ~130 ms statt ~50 ms, kein weiterer Code; DoD 3 wird mit diesen Zahlen belegt.
3. Zaehl-Cache auf der Platte (State-Dir, Schluessel = Ref-Hashes): warm schnell, kalt weiter langsam; neue State-Datei, mehr Code.
- **2026-10-17 17:09 · Alexander Sacharov** — decision (Teamlead, 17.10., Grund: Constraint 'Instant startup'): Variante 1, asynchron. Bis die Ref-Zahlen da sind, zeigt die Zeile keinen Platten-Stand, sondern einen Platzhalter ('… open'), damit nie eine falsche Zahl wie 0/0 aufblitzt und springt. Name, Reihenfolge und Fehlerzeile ('unreadable') sofort wie bisher. DoD 3 wird mit dem Zeitpunkt 'erstes Bild' belegt, als zweite Zahl, wann die Zaehler stehen.
- **2026-10-17 17:13 · Alexander Sacharov** — critique Runde 3 (async): passt zum vorhandenen Muster fetchRefs in internal/tui/refs.go (Arbeit als tea.Cmd, Ergebnis als Msg). Bewusst stehen gelassen: (a) ein aelterer Zaehllauf, der nach einem neueren ankommt, ueberschreibt dessen Zahlen - beide liegen Sekundenbruchteile auseinander, ein Generationszaehler waere Code fuer keinen sichtbaren Fall; (b) 'r' und Hinzufuegen zeigen kurz wieder '… open' - so verlangt es die Entscheidung des Teamleads (nie eine Zahl, die springt).
- **2026-10-17 17:17 · Alexander Sacharov** — review: Abnahme empfohlen. Offen fuer die Person: Release-Note Zeile 1 verspricht 'same as jaira list', der Launcher summiert aber alle Boards, jaira list zeigt nur today's board (onBoard, internal/cli/root.go:378). Eine Wortaenderung in core/release/notes/PXYN99.md.
- **2026-10-17 17:18 · Alexander Sacharov** — Review-Gaps nach Teamlead (17.10.): Gap 1 behoben - Release-Note sagt jetzt 'as jaira list, summed over jaira boards' (Launcher summiert alle Boards, jaira list zeigt nur today's board). Gap 2 (aelterer Zaehllauf kann nach doppeltem r einen frischeren ueberschreiben, bis zum naechsten r) und Gap 3 (today's board zaehlt gegen den Katalog lane.Load(""), so schon vor diesem Diff) bleiben bewusst stehen.
- **2026-10-17 18:59 · Alexander Sacharov** — explain: Der Launcher zaehlt jetzt auch Tickets, die nur auf ihrem Ref liegen, im Hintergrund, und zeigt bis dahin '… open'. Die Release-Note ist korrigiert ('summed over jaira boards'); offen bleiben zwei kleine, bewusst stehengelassene Punkte.
Warum: Zahlen im Launcher stimmen jetzt mit jaira list ueberein, ohne den ersten Bildaufbau zu bremsen. Die Reste (aelterer Zaehllauf ueberschreibt nach doppeltem r einen neueren; today's board zaehlt gegen den Katalog) sind selten und nicht neu.
1) Ticket so annehmen, die zwei Restpunkte bleiben bewusst offen [accept]
2) Zurueckgeben: Generationszaehler gegen das Ueberschreiben durch einen aelteren Zaehllauf einbauen [give back]
Empfehlung: 1 - der Fehler zeigt sich nur bei schnellem doppeltem r und heilt beim naechsten r, ein Zaehler waere Code fuer fast keinen sichtbaren Fall.
- **2026-10-17 19:00 · Alexander Sacharov** — decision: 1) Ticket so annehmen, die zwei Restpunkte bleiben bewusst offen
