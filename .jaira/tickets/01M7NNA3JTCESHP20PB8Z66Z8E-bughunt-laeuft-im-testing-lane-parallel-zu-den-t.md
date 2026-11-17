---
id: 01M7NNA3JTCESHP20PB8Z66Z8E
title: Bughunt laeuft im testing-Lane parallel zu den Tests und konvergiert
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Ticket kommt nach einem Bughunt in hoechstens einer zusaetzlichen Runde durch: der Bughunt laeuft im testing-Lane parallel zu den Tests, gibt nur Fehler zurueck, die die Aenderung selbst eingefuehrt hat, und liest in spaeteren Runden nur den neuen Diff."
context: |-
  Was heute schiefgeht:
  - jaira-role-bughunt liest die beruehrten Module ganz und findet dabei alte Fehler, die die Aenderung nicht eingefuehrt hat.
  - jaira-worker/SKILL.md (Abschnitt Bughunt, ca. Zeile 220) schickt jede Fundstelle zurueck nach in-progress, danach laeuft critique -> optimize -> testing und ein neuer Bughunt.
  - Der neue Bughunt liest die Module wieder ganz und findet die naechsten alten Fehler. Die Schleife konvergiert nicht; viele Runden pro Ticket.
  - Bughunt laeuft heute nach testing, nacheinander, ausserhalb jedes Lanes: auf der Tafel unsichtbar, kein Gate.
  Entscheidungen von Alex (17.11.):
  1. Bughunt wird Teil des testing-Lanes: der Lane-Prompt startet den Hunter zuerst als frischen Agenten im Hintergrund, laesst die Tests parallel laufen, wartet auf beide. Lane liefert test-verdict UND hunt-verdict. Kein Kern-Umbau (ein Ticket steht in genau einem Lane).
  2. Folgerunde: sind die Dateien dieselben wie beim ersten Hunt, liest der Hunter nur den Diff seit dem Commit aus 'bughunt: done at <commit>'. Neue Dateien liest er ganz. Die Dateiliste steht in der Abschlusszeile.
  3. Jede Fundstelle wird markiert: change (von dieser Aenderung eingefuehrt) oder pre-existing (auf master schon da). Nur change schickt zurueck nach in-progress.
  4. pre-existing: gleiches Modul und Thema -> note fuer Review/Signoff, die Person entscheidet ueber dod --add; anderes Thema -> eigenes Ticket mit --follows <id>.
definition-of-done: Lane testing (.jaira/lanes/testing.md und die ausgelieferte Fassung) startet den Bughunt parallel zu den Tests und liefert test-verdict und hunt-verdict; das Gate verweigert den Weg ohne hunt-verdict
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-11-17T21:06:22Z
updated-at: 2026-11-17T21:27:56Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1175
claimed-at: 2026-11-17T21:06:50Z
foo-check: ""
outcome-what: "Der testing-Lane startet den Bughunt als Hintergrund-Agent vor seinen Tests und schuldet hunt-verdict neben test-verdict; jaira-role-bughunt markiert Fundstellen change/pre-existing und liest im Folge-Hunt nur den Diff seit seinem letzten 'done at'; worker und dispatcher starten keinen eigenen Hunter mehr."
outcome-why: "Der Bughunt nach testing las die Module jede Runde ganz und schickte jeden alten Fehler zurueck nach in-progress, so konvergierte die Schleife nicht; ausserhalb jedes Lanes war er zudem unsichtbar und ohne Gate."
outcome-resolves: "Ein Ticket kommt nach einem Bughunt in hoechstens einer weiteren Runde durch: nur change-Fundstellen schicken zurueck, die Folgerunde liest nur neuen Diff, und das Gate verweigert testing ohne hunt-verdict."
review-summary: none
round: 4
review-gaps: "none - Runde 3: Rueckfall-Absatz steht in worker und dispatcher je einmal, weil es zwei getrennt installierte Skills sind; nichts verwaist"
test-verdict: "pass: Runde 2 auf 5b249332 - go build+vet ok, go test -race -count=1 ./... RC=0; Scratch-Board: nur test-verdict -> move aus testing exit 3 (hunt-verdict leer); verankerter grep trifft von zwei Notizen (decision zitiert die Zeile, echte Abschlusszeile) nur die echte (count 1); zwei parallele jaira set landen beide (Runde 1)"
hunt-verdict: "fail: 1 change (bughunt done at 5b249332, Folge-Hunt nur Diff): worker-Schleife Schritt 3 verweist nicht auf den Rueckfall-Hunt"
---

# Bughunt laeuft im testing-Lane parallel zu den Tests und konvergiert

## Definition of Done

- [x] Lane testing (.jaira/lanes/testing.md und die ausgelieferte Fassung) startet den Bughunt parallel zu den Tests und liefert test-verdict und hunt-verdict; das Gate verweigert den Weg ohne hunt-verdict
  proof: core/lane/builtin/27-testing.md:11,24,86 (= .jaira/lanes/testing.md); TestTestingOwesHuntVerdict in core/gate/huntverdict_test.go
- [x] jaira-role-bughunt markiert jede Fundstelle change oder pre-existing, geprueft gegen die Basis; nur change fuehrt zurueck nach in-progress
  proof: core/role/builtin/jaira-role-bughunt/SKILL.md:104-130 (Change or pre-existing, git blame gegen Ticket-Commits); core/lane/builtin/27-testing.md:76-79 (nur change fuehrt zurueck)
- [x] Folge-Hunt liest nur den Diff seit dem letzten 'done at', solange die Dateiliste gleich ist; neue Dateien ganz; die Abschlusszeile nennt die Dateiliste
  proof: core/role/builtin/jaira-role-bughunt/SKILL.md:43-66 (verankerter grep auf die letzte done-at-Zeile; Folge-Hunt nur git diff <commit>, neue Dateien ganz), :151-155 (files: in der Abschlusszeile)
- [x] pre-existing: gleiches Thema als note fuer die Person, anderes Thema als Ticket mit --follows
  proof: core/lane/builtin/27-testing.md:78-84 (gleiches Thema note fuer die Person, anderes Thema jaira create --follows)
- [x] jaira-worker und jaira-dispatcher rufen den Bughunt nicht mehr separat auf; Tests, shipped.sha256 und Release-Notiz
  proof: core/role/builtin/jaira-worker/SKILL.md:145-148,200-223, core/role/builtin/jaira-dispatcher/SKILL.md:268-288 (Hunt im testing-Lane; eigener Hunter nur, wo testing kein hunt-verdict schuldet); go test RC=0; core/shipped/shipped.sha256; core/release/notes/Z66Z8E.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-17 21:06 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/51013.sock
- **2026-11-17 21:07 · Alexander Sacharov** — decision (Alex, 17.11.): EIN Lane testing, zwei Aufgaben. Der Agent, der testing faehrt, startet zuerst den Bughunt als Hintergrund-Subagent (Agent mit run_in_background, nur die Ticket-id im Prompt) und fuehrt die Tests danach SELBST aus, ohne zu warten. Danach wartet er auf den Hunter und schreibt test-verdict und hunt-verdict. Kein zweiter Lane, kein zweiter Worker fuer die Tests.
- **2026-11-17 21:07 · Alexander Sacharov** — decision: EIN Lane testing, zwei Aufgaben. Der Agent, der testing faehrt, startet zuerst den Bughunt als Hintergrund-Subagent (Agent mit run_in_background, nur die Ticket-id im Prompt) und fuehrt die Tests danach SELBST aus, ohne auf den Hunter zu warten. Danach wartet er auf den Hunter und schreibt test-verdict und hunt-verdict. Kein zweiter Lane, kein zweiter Worker fuer die Tests. (Alex, 17.11., via teamlead)
- **2026-11-17 21:07 · Alexander Sacharov** — decision: Branch feat/z66z8e steht auf feat/kcrs3w (a7c5f3fc), nicht auf master: jaira-role-bughunt und der Bughunt-Abschnitt in worker/dispatcher existieren nur dort, KCRS3W steht in signoff. Der PR dieses Tickets traegt KCRS3W mit, bis KCRS3W gemerged ist; danach rebase auf master.
- **2026-11-17 21:08 · Alexander Sacharov** — decision (worker): Offene Punkte aus Code und Default geklaert: (a) Hunter markiert change/pre-existing per git blame der Zeilen auf dem ausloesenden Pfad - Commit nennt die Ticket-id oder uncommitted = change, sonst pre-existing; im Zweifel change, mit Begruendung. (b) Hunter bleibt blind fuer das Ziel; ob pre-existing gleiches Thema ist (note) oder anderes (Ticket --follows), entscheidet der testing-Agent, der das Ticket kennt. (c) Folge-Hunt findet seine letzte Zeile ueber 'jaira show <id> | grep "bughunt: done at"' - nur diese Zeile, sonst nichts vom Ticket. Dateien aus der alten Liste: nur git diff <commit>; neue Dateien ganz; Commit kein Vorfahr von HEAD (rebase) = voller Hunt. (d) hunt-verdict ist ein freies Frontmatter-Feld wie test-verdict; gate.fieldFilled liest es ueber den default-Zweig, kein Go-Umbau.
- **2026-11-17 21:12 · Alexander Sacharov** — Basis: feat/z66z8e per fast-forward auf upstream/master be549519 (0.5.13) gesetzt; KCRS3W ist dort als 5000cd95 gemerged, die Notiz zur Basis feat/kcrs3w ist damit erledigt - kein KCRS3W-Commit mehr im PR dieses Tickets. Hinweis vom teamlead 17.11.
- **2026-11-17 21:14 · Alexander Sacharov** — pre-existing, nicht von diesem Ticket: 'jaira update --lanes --dry-run' rendert die Agent-Datei-Vorschau aus den Lanes auf der Platte VOR dem Refresh (internal/cli/update.go:131-134, RefreshShipped schreibt bei dry nicht) - die Vorschau meldet AGENTS.md (updated), obwohl der echte Lauf sie unveraendert liesse. Sichtbar geworden, weil testing jetzt hunt-verdict produziert; TestUpdateDryRunWritesNothing prueft deshalb nur noch, dass CLAUDE.md in der Zeile steht. Die Person entscheidet, ob daraus ein Ticket wird.
- **2026-11-17 21:16 · Alexander Sacharov** — critique: (1) 27-testing.md - der testing-Agent soll nur die bughunt-Notizen dieser Runde lesen (nach der vorletzten done-at-Zeile, bis zur Zeile mit HEAD), sonst schickt er schon behobene change-Fundstellen erneut zurueck und die Schleife konvergiert wieder nicht. (2) bughunt SKILL.md - change-Pruefung per %s (Betreff) passt nicht zur Dateiliste per --grep (ganzer Text): auf %B | grep -i <id> angleichen. Gelassen: review-Lane (40-review.md) schickt laut KCRS3W auch uebernommene Fehler als blocker zurueck - das ist nicht Teil dieser DoD, offene Frage fuer die Person bei der Abnahme.
- **2026-11-17 21:18 · Alexander Sacharov** — critique Runde 2: beide Fundstellen aus Runde 1 in 338a5ede behoben (testing liest nur die Notizen dieser Runde; Hunter prueft %B | grep -i <id>). Nichts weiter.
- **2026-11-17 21:21 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-role-bughunt/SKILL.md:47 (and core/lane/builtin/27-testing.md:69-73, .jaira/lanes/testing.md same lines) - grep "bughunt: done at" is not anchored to a closing line: any note that merely quotes the phrase (a decision, a testing note) is taken as the last round. The hunter then treats a first hunt as a follow-up or misreads files:, and testing cuts its "this round" window at the wrong note. ran: jaira show Z66Z8E | grep "bughunt: done at" | tail -n 1 on this ticket before any hunt returns the worker decision note of 21:08, not a closing line; reproduced on a scratch board (a decision note added after a real closing line wins tail -n 1). Anchor it to the note head, e.g. grep -E "— bughunt: done at [0-9a-f]{7,}"
- **2026-11-17 21:21 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-worker/SKILL.md:198-204 and core/role/builtin/jaira-dispatcher/SKILL.md:268-274 - the hunt now runs only inside a testing lane, and worker says no hunt is due where no lane owes hunt-verdict. A board from plain jaira init has no testing lane (lanes: brainstorm todo pre-process in-progress human review signoff done blocked), so on every default board worker and dispatcher no longer hunt at all; before this change they hunted before review. Only the teamlead check (jaira-teamlead/SKILL.md:205-219) still catches it, and only when a teamlead runs. ran: jaira init on a scratch repo, jaira lanes (no testing; testing/critique/optimize listed as shipped-not-on-board)
- **2026-11-17 21:21 · Alexander Sacharov** — bughunt: done at 338a5ede - 2 findings (2 change, 0 pre-existing) in 12 files, first hunt; concurrency, error paths, input, permissions checked; gate verified on scratch board (no hunt-verdict -> exit 3, round-2 overwrite works); files: .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/gate/huntverdict_test.go core/lane/builtin/27-testing.md core/release/notes/Z66Z8E.md core/role/builtin/jaira-dispatcher/SKILL.md core/role/builtin/jaira-role-bughunt/SKILL.md core/role/builtin/jaira-teamlead/SKILL.md core/role/builtin/jaira-worker/SKILL.md core/shipped/shipped.sha256 internal/cli/updatedry_test.go
- **2026-11-17 21:23 · Alexander Sacharov** — in-progress Runde 3: beide bughunt-change-Fundstellen behoben. (1) grep verankert: '— bughunt: done at [0-9a-f]{7,} ' in bughunt SKILL, testing-Lane und teamlead - auf diesem Ticket getroffen nur die echte Abschlusszeile. (2) worker/dispatcher starten den Hunter selbst, wenn der testing-Lane des Boards kein hunt-verdict schuldet (kein testing nach jaira init, oder alte Lane) - wie vor diesem Ticket, vor fresh/Person-Lane.
- **2026-11-17 21:24 · Alexander Sacharov** — critique Runde 3: nur die Reparatur 5b249332 gelesen - verankerter grep und Rueckfall-Hunt in worker/dispatcher beantworten beide bughunt-Fundstellen; keine neue Fundstelle.
- **2026-11-17 21:25 · Alexander Sacharov** — bughunt: change: core/role/builtin/jaira-worker/SKILL.md:145 - the fallback hunt for boards whose testing lane owes no hunt-verdict (lines 205-217) is never reached from the loop: adc11846 dropped step 3's 'Run the bug hunt first - see The bug hunt below' and 5b249332 restored only the section, so a worker following steps 1-4 on a fresh 'jaira init' board goes from its last lane straight into the fresh/person lane with no hunt; the dispatcher keeps its fallback inside its loop section (jaira-dispatcher/SKILL.md:276), the worker does not; read (git blame -L 145,148 shows step 3 rewritten in adc11846)
- **2026-11-17 21:25 · Alexander Sacharov** — bughunt: done at 5b249332 - 1 finding (1 change, 0 pre-existing), follow-up hunt over diff since 338a5ede in 6 changed files; concurrency, error paths, input, permissions checked; anchored grep verified against jaira show (matches the closing line, UTF-8 and LC_ALL=C); shipped.sha256 holds every new builtin hash; go test core/gate core/lane core/role internal/cli rc=0; files: .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/gate/huntverdict_test.go core/lane/builtin/27-testing.md core/release/notes/Z66Z8E.md core/role/builtin/jaira-dispatcher/SKILL.md core/role/builtin/jaira-role-bughunt/SKILL.md core/role/builtin/jaira-teamlead/SKILL.md core/role/builtin/jaira-worker/SKILL.md core/shipped/shipped.sha256 internal/cli/updatedry_test.go
- **2026-11-17 21:27 · Alexander Sacharov** — in-progress Runde 4: worker-Schleife Schritt 3 (SKILL.md:145) verweist wieder auf den Rueckfall-Hunt fuer Boards ohne hunt-verdict im testing-Lane. Zweiter Ruecklauf aus testing durch den Hunt; die Folgerunde las nur 5b249332..HEAD - die Konvergenz-Regel greift.
