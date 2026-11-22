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
updated-at: 2026-11-22T20:25:13Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2119162
claimed-at: 2026-11-22T20:17:27Z
outcome-what: "Nach 3 Rückgaben aus einer Lane geht die Frage zuerst an den Teamlead: Worker/Dispatcher/Lane-Worker lassen den Ticket mit gesetzter question in seiner Lane und schicken sie an die parent-Session; der Teamlead entscheidet (decision (teamlead <datum>)) oder legt ihn selbst nach human. Gate-Meldung, Board-Loop-Text, lanes show und --for-lane sagen 'the teamlead decides, or a person'."
outcome-why: "Alex sah im Cockpit Schleifenfragen (CRASZS, DF46B9, QTQW15), die der Teamlead fast immer selbst entschied."
outcome-resolves: "Ticket K6PBYR: Schleifenbremse entscheidet zuerst der Teamlead. Geprüft: Nebenläufigkeit n/a (Text + Meldung), Fehlerpfad: parent unerreichbar -> human wie bisher; ungewöhnliche Eingabe n/a; Rechte: --force bleibt beim Entscheidenden."
review-summary: "core/lane/builtin/20-in-progress.md:99 '(a teamlead, or the worker or dispatcher it started)' mischt zwei Fälle: ein Lane-Worker unter einem Dispatcher kennt dessen Session nicht (jaira-role-lane: 'nothing tells you its name') — schreibe '(--parent, usually a teamlead's)' und verweise den Lane-Worker auf seine Rolle"
round: 2
back-from: critique
back-by: lane
back-count: 1
---

# Schleifenbremse: nach 3 Rückgaben entscheidet zuerst der Teamlead, erst wenn er nicht kann, die Person

## Definition of Done

- [x] Nach der 3. Rückgabe schickt der Worker die Frage mit den Optionen an die parent-Session statt den Ticket nach human zu legen; der Ticket bleibt in seiner Lane mit offener question (Test oder Rollen-Text + Test)
  proof: core/role/builtin/jaira-worker/SKILL.md:407, jaira-dispatcher/jaira-role-lane SKILL, core/lane/builtin/20-in-progress.md:99; tests TestThirdSendBackGoesToTheTeamleadFirst, TestInProgressPromptHandsTheThirdReturnToTheTeamleadFirst
- [x] Der Teamlead entscheidet und schreibt 'decision (teamlead <datum>)' auf den Ticket, oder eskaliert an die Person; erst die Eskalation legt den Ticket in human (jaira-teamlead SKILL)
  proof: core/role/builtin/jaira-teamlead/SKILL.md:104 (Three rounds of one lane: you decide first); TestThirdSendBackGoesToTheTeamleadFirst
- [x] Ohne erreichbare parent-Session verhält sich alles wie heute (human-Lane) (Test)
  proof: Gate unverändert bis auf Meldung (core/gate/gate.go:450, ohne parent -> human --question); 20-in-progress.md:107; Worker-SKILL 'as without a parent'; TestMaxRoundsRefusalNamesThePersonsLane, TestThirdSendBackGoesToTheTeamleadFirst
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
