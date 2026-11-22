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
updated-at: 2026-11-22T20:32:01Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2119162
claimed-at: 2026-11-22T20:17:27Z
outcome-what: "Nach 3 Rückgaben aus einer Lane geht die Frage zuerst an den Teamlead: Worker/Dispatcher/Lane-Worker lassen den Ticket mit gesetzter question in seiner Lane und schicken sie an die parent-Session; der Teamlead entscheidet (decision (teamlead <datum>)) oder legt ihn selbst nach human. Gate-Meldung, Board-Loop-Text, lanes show und --for-lane sagen 'the teamlead decides, or a person'."
outcome-why: "Alex sah im Cockpit Schleifenfragen (CRASZS, DF46B9, QTQW15), die der Teamlead fast immer selbst entschied."
outcome-resolves: "Ticket K6PBYR: Schleifenbremse entscheidet zuerst der Teamlead. Geprüft: Nebenläufigkeit n/a (Text + Meldung), Fehlerpfad: parent unerreichbar -> human wie bisher; ungewöhnliche Eingabe n/a; Rechte: --force bleibt beim Entscheidenden."
review-summary: "none"
round: 3
back-from: testing
back-by: lane
back-count: 1
review-gaps: "none — nur Texte, eine Meldung und Tests; keine neue Funktion, kein Duplikat (Gate-Ausweg --force wiederverwendet statt neuem Flag)"
test-verdict: "pass: go build ok, go test core/gate core/board core/lane core/role internal/cli -count=1 RC=0 (ganze Suite vorher RC=0), shipped-history --check ok; DoD 1-4 im Baum geprüft (Proof-Zeilen 20-in-progress.md:98/108, worker SKILL:407, teamlead SKILL:104, gate.go:450 stimmen); Scratch-Board: lanes show 'then the teamlead decides, or a person', 4. Rückgabe aus critique exit 3 mit neuer Meldung (parent -> question= bleibt, sonst human), jaira set question= lässt Ticket in critique, --for-lane zeigt neuen max-rounds-Satz; Fehlerpfad = die Verweigerung selbst (rc 3, nichts geschrieben); kein neuer Schreibpfad, daher kein Zwei-Schreiber-Test"
hunt-verdict: "fail: bughunt done at fe2008000, 1 change (Teamlead-Entscheidung unsichtbar für back-count und in-progress-Prompt), 0 pre-existing"
---

# Schleifenbremse: nach 3 Rückgaben entscheidet zuerst der Teamlead, erst wenn er nicht kann, die Person

## Definition of Done

- [x] Nach der 3. Rückgabe schickt der Worker die Frage mit den Optionen an die parent-Session statt den Ticket nach human zu legen; der Ticket bleibt in seiner Lane mit offener question (Test oder Rollen-Text + Test)
  proof: core/role/builtin/jaira-worker/SKILL.md:407, jaira-dispatcher/jaira-role-lane SKILL, core/lane/builtin/20-in-progress.md:103; tests TestThirdSendBackGoesToTheTeamleadFirst, TestInProgressPromptHandsTheThirdReturnToTheTeamleadFirst
- [x] Der Teamlead entscheidet und schreibt 'decision (teamlead <datum>)' auf den Ticket, oder eskaliert an die Person; erst die Eskalation legt den Ticket in human (jaira-teamlead SKILL)
  proof: core/role/builtin/jaira-teamlead/SKILL.md:104; core/move/move.go:209 overridesMaxRounds (Teamlead-Entscheidung per --force zählt als Entscheidung); TestForcePastMaxRoundsIsRecordedAsADecision, TestThirdSendBackGoesToTheTeamleadFirst
- [x] Ohne erreichbare parent-Session verhält sich alles wie heute (human-Lane) (Test)
  proof: Gate unverändert bis auf Meldung (core/gate/gate.go:450, ohne parent -> human --question); 20-in-progress.md:113; Worker-SKILL 'as without a parent'; TestMaxRoundsRefusalNamesThePersonsLane, TestThirdSendBackGoesToTheTeamleadFirst
- [x] Board-Text der Loop-Regel und jaira move-Meldung sagen 'the teamlead decides, or a person'; Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/board/announce.go:254, core/gate/gate.go:450, internal/cli/flow.go:807, internal/cli/tickets.go:1537; TestLaneSectionSaysTheTeamleadDecidesFirst; core/release/notes/K6PBYR.md; core/shipped/shipped.sha256 (shipped-history.sh --check ok)

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
