---
id: 01M82EEXPS56MEHEHPTVK6PBYR
title: "Schleifenbremse: nach 3 Rückgaben entscheidet zuerst der Teamlead, erst wenn er nicht kann, die Person"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Nach der dritten Rückgabe aus derselben Lane geht die Frage an den Teamlead (parent-Session); er entscheidet und notiert es, oder er gibt sie an die Person weiter. Nur dann landet der Ticket in der human-Lane."
context: |-
  Was heute falsch ist (Alex 22.11.2026): die Board-Regel sagt 'At most 3 times in a row: the next send-back is refused, and a person decides'. Worker legen den Ticket deshalb mit --question in human; er steht dann bei Alex im Cockpit als 'waiting for you' (z.B. CRASZS, DF46B9, QTQW15 am 22.11.).
  In der Praxis hat der Teamlead diese Fragen fast immer selbst entschieden (Regel 'Ränder wachsen -> kleineres Design'), Alex sah sie trotzdem.
  Gewünscht: 3 Rückgaben -> Teamlead entscheidet; kann er nicht (echte Produkt-/Design-Frage der Person) -> human.
  Stellen: Loop-Text in .jaira/lanes/*.md bzw. core/lane (refused send-back, Meldung), jaira-worker/jaira-dispatcher SKILL (wohin die Frage geht), jaira-teamlead SKILL (Teamlead entscheidet, notiert 'decision (teamlead ...)', eskaliert sonst mit AskUserQuestion). Ohne parent-Session (kein Teamlead): wie heute human.
definition-of-done: Nach der 3. Rückgabe schickt der Worker die Frage mit den Optionen an die parent-Session statt den Ticket nach human zu legen; der Ticket bleibt in seiner Lane mit offener question (Test oder Rollen-Text + Test)
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-11-22T20:16:47Z
updated-at: 2026-11-22T20:40:48Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2119162
claimed-at: 2026-11-22T20:17:27Z
outcome-what: "Give-back: Teamlead-Reset per 'jaira set <id> back-count=0' (vorhandener Weg, Test), in-progress-Text ohne die Notiz-Zeitpunkt-Regel, Kommentar overridesMaxRounds repariert"
outcome-why: "Alex sah im Cockpit Schleifenfragen (CRASZS, DF46B9, QTQW15), die der Teamlead fast immer selbst entschied."
outcome-resolves: "Ticket K6PBYR: Schleifenbremse entscheidet zuerst der Teamlead. Geprüft: Nebenläufigkeit n/a (Text + Meldung), Fehlerpfad: parent unerreichbar -> human wie bisher; ungewöhnliche Eingabe n/a; Rechte: --force bleibt beim Entscheidenden."
review-summary: "Nach der dritten Rückgabe aus einer Lane geht die Entscheidung zuerst an den Teamlead: Worker-, Dispatcher- und Lane-Rollen sowie der in-progress-Lane-Text sagen jetzt 'Ticket in der Lane lassen, question setzen, an die parent-Session schicken'; nur ohne parent (oder unerreichbar) geht es wie bisher nach human. Die Teamlead-Rolle bekommt einen Abschnitt 'entscheiden (decision (teamlead <datum>)-Notiz) oder an die Person geben und selbst nach human legen'. Im Code: die max_rounds-Meldung (gate.go), jaira lanes show, --for-lane-Text, Template und der CLAUDE.md-Lane-Abschnitt sagen 'the teamlead decides, or a person'; move.go bucht ein --force, das die max_rounds-Verweigerung überstimmt, als Entscheidung (back-by: person, count 1), damit die Schleife wieder Runden hat und der nächste in-progress-Worker arbeitet statt erneut zu stoppen. Release-Note und shipped.sha256 dabei; Tests: Rollen-/Lane-Text-Contains, Gate-Meldung, TestForcePastMaxRoundsIsRecordedAsADecision."
round: 4
back-from: signoff
back-by: person
back-count: 1
review-gaps: |-
  note: jedes --force auf einer Rückgabe bei count 3 setzt jetzt back-by: person/1, auch wenn es aus einem anderen Grund gegeben wurde (z.B. veralteter hunt, vgl. TestForceOverridesStaleHunt) - ein Worker, der so forciert, hebt die Schleifenbremse unbemerkt auf; die Rollentexte verbieten es, der Code unterscheidet nicht. Fix wenn gewünscht: nur bei vorhandener decision-Notiz oder ByPerson als Entscheidung buchen.
  note: 20-in-progress.md "decision (teamlead) written after the newest send-back" - der Zeitpunkt der Rückgabe steht nirgends im Frontmatter, der Worker muss ihn aus der Notizreihenfolge erraten; eine alte Teamlead-Notiz aus einer früheren Runde kann als aktuell gelesen werden.
  note: Fall count 3 schon in in-progress ohne --force: nach dem Teamlead-Beschluss arbeitet in-progress, aber die nächste Rückgabe derselben Lane wird sofort wieder verweigert (back-count bleibt lane/3) - jede weitere Runde braucht eine zweite Teamlead-Entscheidung plus --force; im Teamlead-Abschnitt erwähnen.
  note: core/move/move.go Kommentar zu overridesMaxRounds endet in einem kaputten Satz ("and the next worker reading back-by: lane, count 4 stopping once more") - umformulieren zu "instead of the next worker reading back-by: lane, count 4 and stopping once more".
  note: DoD 1 und 3 sind nur über strings.Contains auf Rollen-/Lane-Text belegt (bewusst, siehe decision (worker)); das Routing selbst ist nicht ausführbar getestet.
test-verdict: "pass: Runde 2 nach Fix ac250fa6e — go test -race core/move core/gate core/lane core/role RC=0, internal/cli RC=0, shipped-history --check ok; neuer Test TestForcePastMaxRoundsIsRecordedAsADecision: 4. Rückgabe aus testing verweigert (max_rounds), mit --force -> back testing/person/1; DoD 1-4 Proof-Zeilen neu geprüft (20-in-progress.md:103/113, move.go:209); Scratch-Board Runde 1: Meldung exit 3, question= lässt Ticket in critique; kein neuer Schreibpfad außer back-by im bestehenden Mutate-Lock"
hunt-verdict: "pass: bughunt done at ac250fa6e, 0 change, 0 pre-existing; Befund aus Runde 1 (fe2008000) in ac250fa6e behoben"
review-verdict: "accept: alle vier DoD-Punkte sind im Diff erfüllt, go test ./core/... ./internal/... grün, shipped-history --check ok; kein Blocker, fünf kleine Notizen (vor allem: --force aus anderem Grund bei count 3 hebt die Bremse auf)."
review-round: 3
review-verdict-gaps: |-
  note: jedes --force auf einer Rückgabe bei count 3 setzt jetzt back-by: person/1, auch wenn es aus einem anderen Grund gegeben wurde (z.B. veralteter hunt, vgl. TestForceOverridesStaleHunt) - ein Worker, der so forciert, hebt die Schleifenbremse unbemerkt auf; die Rollentexte verbieten es, der Code unterscheidet nicht. Fix wenn gewünscht: nur bei vorhandener decision-Notiz oder ByPerson als Entscheidung buchen.
  note: 20-in-progress.md "decision (teamlead) written after the newest send-back" - der Zeitpunkt der Rückgabe steht nirgends im Frontmatter, der Worker muss ihn aus der Notizreihenfolge erraten; eine alte Teamlead-Notiz aus einer früheren Runde kann als aktuell gelesen werden.
  note: Fall count 3 schon in in-progress ohne --force: nach dem Teamlead-Beschluss arbeitet in-progress, aber die nächste Rückgabe derselben Lane wird sofort wieder verweigert (back-count bleibt lane/3) - jede weitere Runde braucht eine zweite Teamlead-Entscheidung plus --force; im Teamlead-Abschnitt erwähnen.
  note: core/move/move.go Kommentar zu overridesMaxRounds endet in einem kaputten Satz ("and the next worker reading back-by: lane, count 4 stopping once more") - umformulieren zu "instead of the next worker reading back-by: lane, count 4 and stopping once more".
  note: DoD 1 und 3 sind nur über strings.Contains auf Rollen-/Lane-Text belegt (bewusst, siehe decision (worker)); das Routing selbst ist nicht ausführbar getestet.
review-check: "1. go test ./core/move -run TestForcePastMaxRoundsIsRecordedAsADecision -v  -> PASS (vierte Rückgabe wird verweigert, mit --force steht back-by: person, count 1)  2. go test ./core/gate ./core/role ./core/lane ./core/board  -> alles ok  3. go run ./cmd/jaira lanes show testing  -> Zeile 'Max rounds:  3 in a row, then the teamlead decides, or a person'  4. core/role/builtin/jaira-teamlead/SKILL.md öffnen, Abschnitt 'Three rounds of one lane: you decide first'  -> zwei Wege: selbst entscheiden mit decision-Notiz, oder per jaira move --to human an die Person  5. cat core/release/notes/K6PBYR.md  -> drei einzeilige Einträge"
---

# Schleifenbremse: nach 3 Rückgaben entscheidet zuerst der Teamlead, erst wenn er nicht kann, die Person

## Definition of Done

- [x] Nach der 3. Rückgabe schickt der Worker die Frage mit den Optionen an die parent-Session statt den Ticket nach human zu legen; der Ticket bleibt in seiner Lane mit offener question (Test oder Rollen-Text + Test)
  proof: core/role/builtin/jaira-worker/SKILL.md:407, jaira-dispatcher/jaira-role-lane SKILL, core/lane/builtin/20-in-progress.md:103; tests TestThirdSendBackGoesToTheTeamleadFirst, TestInProgressPromptHandsTheThirdReturnToTheTeamleadFirst
- [x] Der Teamlead entscheidet und schreibt 'decision (teamlead <datum>)' auf den Ticket, oder eskaliert an die Person; erst die Eskalation legt den Ticket in human (jaira-teamlead SKILL)
  proof: core/role/builtin/jaira-teamlead/SKILL.md:104; core/move/move.go:209 overridesMaxRounds; TestForcePastMaxRoundsIsRecordedAsADecision, TestMaxRoundsBackCountResetIsTheTeamleadsDecision, TestThirdSendBackGoesToTheTeamleadFirst
- [x] Ohne erreichbare parent-Session verhält sich alles wie heute (human-Lane) (Test)
  proof: Gate unverändert bis auf Meldung (core/gate/gate.go:450, ohne parent -> human --question); 20-in-progress.md:113; Worker-SKILL 'as without a parent'; TestMaxRoundsRefusalNamesThePersonsLane, TestThirdSendBackGoesToTheTeamleadFirst
- [x] Board-Text der Loop-Regel und jaira move-Meldung sagen 'the teamlead decides, or a person'; Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/board/announce.go:254, core/gate/gate.go:450, internal/cli/flow.go:807, internal/cli/tickets.go:1537; TestLaneSectionSaysTheTeamleadDecidesFirst; core/release/notes/K6PBYR.md; core/shipped/shipped.sha256 (shipped-history.sh --check ok)
- [x] Nach 'decision (teamlead ...)' ist auch im Fall count 3 ohne --force die nächste Rückgabe derselben Lane wieder erlaubt — der Teamlead setzt 'jaira set <id> back-count=0' (decision teamlead 2026-11-22)
  proof: internal/cli/maxrounds_test.go TestMaxRoundsBackCountResetIsTheTeamleadsDecision; core/role/builtin/jaira-teamlead/SKILL.md:134; core/lane/builtin/20-in-progress.md (back-count back to 0)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Gate-Meldung core/gate/gate.go + Test
- [x] Board-Text: announce.go, flow.go, tickets.go, template.go
- [x] in-progress-Lane-Text (builtin)
- [x] jaira-worker/jaira-dispatcher/jaira-role-lane SKILL: Frage an parent
- [x] jaira-teamlead SKILL: entscheiden oder eskalieren
- [x] Rollen-Text-Tests core/role
- [x] Release-Note + shipped.sha256 + Gates

## Progress
- **2026-11-22 20:17 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 20:18 · Alexander Sacharov** — decision (worker): keine offene Personen-Entscheidung. (1) Warten: Ticket bleibt in der verweigernden Lane, question per 'jaira set <id> question=...'. (2) Teamlead entscheidet 'trotzdem zurück/kleineres Design' mit 'jaira move --force' — der Ausweg, den die Gate-Meldung schon nennt ('as the one deciding'); vorwärts (Arbeit annehmen) braucht kein --force. Kein neues Flag: kleiner. (3) Gate prüft keine parent-Session; das Routing steht in Worker/Dispatcher/in-progress-Lane-Text, Gate-Meldung nennt beide Wege. (4) Tests: strings.Contains auf SKILL.md/Lane-Text (core/role-Muster) + Gate-Meldung in core/gate/maxrounds_test.go.
- **2026-11-22 20:24 · Alexander Sacharov** — critique: 20-in-progress.md:99 Klammer '(a teamlead, or the worker or dispatcher it started)' falsch für Lane-Worker ohne Session-Namen -> '(--parent, usually a teamlead's)'. Stehen gelassen: --force für Teamlead-Rücksendung (bestehender Ausweg der Gate-Meldung, kein neues Flag); Regel in 4 SKILLs + Lane wiederholt = bestehendes Muster je Rolle.
- **2026-11-22 20:25 · Alexander Sacharov** — critique Runde 2: Befund aus Runde 1 (20-in-progress.md Klammer) in fe2008000 behoben; sonst nichts.
- **2026-11-22 20:28 · Alexander Sacharov** — bughunt: change: blocker: core/role/builtin/jaira-teamlead/SKILL.md:127-130 with core/lane/builtin/20-in-progress.md:88-97 and core/move/move.go:186-193 - class: a teamlead decision after three rounds is invisible to the round counter and to the in-progress prompt, so every fresh in-progress worker after it escalates again instead of working; example: testing sends back a 4th time, max_rounds refuses, the teamlead writes 'decision (teamlead ...)' and passes --force as the skill says; recordBack sees an agent (atATerminal false, testing is not a person's lane) and writes back-by: lane, back-count: 4; the in-progress lane worker the dispatcher starts fresh reads 'by: lane with count 3 or more: do not start another round' (line 92), and the prompt honours a decision: note only for by: person (line 88-90), so it sets the question and stops again. Same when the third send-back already landed in in-progress (count 3, no --force needed): the teamlead's 'tell the worker how to go on' reaches no fresh lane worker. A further testing send-back is refused again at once (RoundsLeft 0 at count >= 3), so each later round needs another teamlead decision. Before this change only a person decided, and at the board that is recorded as by: person; the agent decider this change introduces can never be. read (recordBack byPerson=atATerminal(), agent shells have no TTY per flow.go:938-946)
- **2026-11-22 20:28 · Alexander Sacharov** — bughunt: done at fe2008000 - 1 finding (1 change, 0 pre-existing) in 7 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/board/announce.go core/board/lanesection_test.go core/gate/gate.go core/gate/maxrounds_test.go core/lane/builtin/20-in-progress.md core/lane/rounds_test.go core/lane/template.go core/release/notes/K6PBYR.md core/role/builtin/jaira-dispatcher/SKILL.md core/role/builtin/jaira-role-lane/SKILL.md core/role/builtin/jaira-teamlead/SKILL.md core/role/builtin/jaira-worker/SKILL.md core/role/role_test.go core/shipped/shipped.sha256 internal/cli/flow.go internal/cli/tickets.go
- **2026-11-22 20:29 · Alexander Sacharov** — testing: zurück wegen bughunt change. Fix-Plan: (a) move.go: ein --force, das die max_rounds-Verweigerung überstimmt, wird als Entscheidung gebucht (back-by: person, count 1) — wie die Rücksendung einer Person; (b) 20-in-progress.md: eine 'decision (teamlead ...)'-Notiz neuer als die letzte Rückgabe ist das Wort des Entscheiders wie by: person (deckt Fall count 3 schon in in-progress); (c) teamlead SKILL sagt beides.
- **2026-11-22 20:34 · Alexander Sacharov** — bughunt: done at ac250fa6e - 0 findings in 2 new modules read whole (core/move/move.go, core/move/move_test.go) and 6 files read as diff since fe2008000, follow-up hunt; concurrency, error paths, input, permissions checked; ran: go test ./core/move ./core/gate ./core/lane ./core/role ./core/shipped (ok), sha256sum of the two new shipped entries (match); files: core/board/announce.go core/board/lanesection_test.go core/gate/gate.go core/gate/maxrounds_test.go core/lane/builtin/20-in-progress.md core/lane/rounds_test.go core/lane/template.go core/move/move.go core/move/move_test.go core/release/notes/K6PBYR.md core/role/builtin/jaira-dispatcher/SKILL.md core/role/builtin/jaira-role-lane/SKILL.md core/role/builtin/jaira-teamlead/SKILL.md core/role/builtin/jaira-worker/SKILL.md core/role/role_test.go core/shipped/shipped.sha256 internal/cli/flow.go internal/cli/tickets.go
- **2026-11-22 20:38 · Alexander Sacharov** — given back by teamlead (22.11.), short round: (1) Fall count 3 ohne --force: nach 'decision (teamlead ...)' muss auch die nächste Rückgabe derselben Lane wieder erlaubt sein - die Entscheidung setzt den Zähler zurück wie --force (Test); (2) Kommentar zu overridesMaxRounds in core/move/move.go reparieren. Rest bleibt note.
- **2026-11-22 20:39 · Alexander Sacharov** — give-back Runde: (1) Reset per 'jaira set <id> back-count=0' — das existiert schon (set schreibt back-count), RoundsLeft rechnet max-rounds minus count, recordBack zählt ab 0+1. Kein neuer Code; Teamlead-SKILL schreibt es vor, in-progress-Text verliert die Regel 'decision-Notiz nach der letzten Rückgabe' (Review-Notiz: Zeitpunkt nicht erkennbar) — count unter 3 reicht. CLI-Test dazu. (2) Kommentar overridesMaxRounds.
