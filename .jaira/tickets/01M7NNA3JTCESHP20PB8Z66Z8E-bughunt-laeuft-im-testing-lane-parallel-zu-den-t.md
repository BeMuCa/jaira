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
updated-at: 2026-11-17T21:18:37Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1175
claimed-at: 2026-11-17T21:06:50Z
foo-check: ""
outcome-what: "Der testing-Lane startet den Bughunt als Hintergrund-Agent vor seinen Tests und schuldet hunt-verdict neben test-verdict; jaira-role-bughunt markiert Fundstellen change/pre-existing und liest im Folge-Hunt nur den Diff seit seinem letzten 'done at'; worker und dispatcher starten keinen eigenen Hunter mehr."
outcome-why: "Der Bughunt nach testing las die Module jede Runde ganz und schickte jeden alten Fehler zurueck nach in-progress, so konvergierte die Schleife nicht; ausserhalb jedes Lanes war er zudem unsichtbar und ohne Gate."
outcome-resolves: "Ein Ticket kommt nach einem Bughunt in hoechstens einer weiteren Runde durch: nur change-Fundstellen schicken zurueck, die Folgerunde liest nur neuen Diff, und das Gate verweigert testing ohne hunt-verdict."
review-summary: |-
  core/lane/builtin/27-testing.md:72: 'jaira show <handle> | grep bughunt:' liefert die change-Notizen ALLER Runden, auch schon behobene - nur die Notizen zwischen der vorletzten und der letzten 'done at'-Zeile (die HEAD nennt) lesen
  core/role/builtin/jaira-role-bughunt/SKILL.md:114: 'git log -1 --format=%s' prueft nur den Betreff, die Dateiliste oben findet Commits per --grep im ganzen Text - beide gleich: 'git log -1 --format=%B <commit> | grep -i <id>'
round: 2
---

# Bughunt laeuft im testing-Lane parallel zu den Tests und konvergiert

## Definition of Done

- [x] Lane testing (.jaira/lanes/testing.md und die ausgelieferte Fassung) startet den Bughunt parallel zu den Tests und liefert test-verdict und hunt-verdict; das Gate verweigert den Weg ohne hunt-verdict
  proof: core/lane/builtin/27-testing.md:11,24,85 (= .jaira/lanes/testing.md); TestTestingOwesHuntVerdict in core/gate/huntverdict_test.go
- [x] jaira-role-bughunt markiert jede Fundstelle change oder pre-existing, geprueft gegen die Basis; nur change fuehrt zurueck nach in-progress
  proof: core/role/builtin/jaira-role-bughunt/SKILL.md:100-126 (Change or pre-existing, git blame gegen Ticket-Commits); core/lane/builtin/27-testing.md:75-78 (nur change fuehrt zurueck)
- [x] Folge-Hunt liest nur den Diff seit dem letzten 'done at', solange die Dateiliste gleich ist; neue Dateien ganz; die Abschlusszeile nennt die Dateiliste
  proof: core/role/builtin/jaira-role-bughunt/SKILL.md:43-62 (Folge-Hunt nur git diff <commit>, neue Dateien ganz), :147-151 (files: in der Abschlusszeile)
- [x] pre-existing: gleiches Thema als note fuer die Person, anderes Thema als Ticket mit --follows
  proof: core/lane/builtin/27-testing.md:77-83 (gleiches Thema note fuer die Person, anderes Thema jaira create --follows)
- [x] jaira-worker und jaira-dispatcher rufen den Bughunt nicht mehr separat auf; Tests, shipped.sha256 und Release-Notiz
  proof: core/role/builtin/jaira-worker/SKILL.md:196, core/role/builtin/jaira-dispatcher/SKILL.md:268 (kein eigener Hunter mehr); go test ./... RC=0; core/shipped/shipped.sha256 (5 Zeilen); core/release/notes/Z66Z8E.md

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
