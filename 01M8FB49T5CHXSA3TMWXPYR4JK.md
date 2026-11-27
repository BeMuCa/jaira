---
id: 01M8FB49T5CHXSA3TMWXPYR4JK
title: "Cockpit: Ticket mit offener Frage in human zeigt Frage statt Abnahme-Checkliste"
status: testing
ready: true
creator: Sascha
goal: "Steht ein Ticket in einer Frage-Lane (human) mit offener Frage, zeigt der cockpit die Frage und ein Antwortfeld – keine 'Check yourself'-Liste und kein Accept, das danach abgelehnt wird."
context: |-
  Was falsch ist: ein Worker stellt eine Frage und schiebt das Ticket nach human. Der cockpit zeigt dafür denselben Abnahme-Bildschirm wie in signoff: 'Check yourself - tick what you have seen' mit den 'seen by a person'-Punkten der DoD.
  Die Person hakt einen Punkt ab, drückt Execute und bekommt erst dann: 'nothing to accept in this lane - mark a note fix, or type what is missing' (texts.ts:80/89, step 'noNext').
  Folge: die Person glaubt, sie solle abnehmen, und hakt DoD-Punkte ab, die erst am Ende prüfbar sind (die Arbeit ist noch gar nicht fertig).
  Wo: core/mod/builtin/jaira-cockpit/hooks/board.ts:1231 wählt den step; noNext greift erst beim Execute. register.tsx:2665 zeigt die seen-Liste unabhängig davon.
  Richtig wäre: Lane ohne next_lane + Ticket hat question → Frage oben zeigen, Antwort eintippen = Antwort (note + zurück an den Worker), seen-Liste ausblenden.
  Gesehen: 27.11.2026 in einem Projekt mit Video-Board.
definition-of-done: "Ticket in human mit question: cockpit zeigt die Frage, keine seen-Checkliste (Test in tests/board.test.ts)"
tags:
  - cockpit
blocked-by: []
related: []
commits: []
created-at: 2026-11-27T20:28:41Z
updated-at: 2026-11-27T20:48:12Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-69193
claimed-at: 2026-11-27T20:33:49Z
outcome-what: "Acceptance pane: offene Frage in nicht-endender Person-Lane wird beantwortet statt abgenommen — getippter Text ist die Antwort (Note 'answer to …', question= gelöscht, zurück an den Worker), auch neben einer fix-Note; Frage oben, seen-Liste ausgeblendet; ohne Eingabefeld verweist die Ablehnung auf Discuss"
outcome-why: "Bei Lane ohne next_lane lief Execute in 'noNext' und die seen-Liste verleitete zum Abhaken vor Ende der Arbeit"
outcome-resolves: "Person sieht im human-Lane die Frage und beantwortet sie direkt; kein 'nothing to accept'"
review-summary: none
review-gaps: "round 2: folded answer() into the giveBack branch of execute (one path for both steps, clearQuestion as first-callback) — answer() is gone; nothing else orphaned; left t.answer as the fail title"
test-verdict: "pass: claude plugin test . 591 pass 0 fail (exit 0) at 5ca1382c7; go test ./core/release/... ok; DoD 1,2,3,5 verified in the tree; round 2 covers the answer beside a fix note (test) and the no-field surface (needAnswerUntyped, read: no harness for canType=false); failure path: refused 'set question=' toasts, note still lands; writes go through jaira under the row lock; DoD 4 left for the person"
round: 2
back-from: testing
back-by: lane
back-count: 1
---

# Cockpit: Ticket mit offener Frage in human zeigt Frage statt Abnahme-Checkliste

## Definition of Done

- [x] Ticket in human mit question: cockpit zeigt die Frage, keine seen-Checkliste (Test in tests/board.test.ts)
  proof: tests/board.test.ts: 'a question waiting in human with no lane after it / with a lane after it: the pane shows it and no checklist…'; register.tsx openAcceptance seen=[] while waitsForAnswer
- [x] Eingetippte Antwort landet als note am Ticket und geht an den Worker zurück; kein 'nothing to accept'-Fehler (Test)
  proof: tests/board.test.ts: 'a question waiting in human with/without a lane after it…' and 'an answer typed beside a note marked fix is still the answer…'; register.tsx execute giveBack/answer branch + clearQuestion
- [x] signoff-Abnahme unverändert (bestehende Tests grün)
  proof: claude plugin test . (2.1.332): 590 pass, 0 fail — describe 'what the person checks themselves' unchanged
- [ ] Von einer Person im cockpit gesehen
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/PYR4JK.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-27 20:33 · Sascha** — parent session: uds:/run/user/1000/cc-socks/24078.sock
- **2026-11-27 20:34 · Sascha** — decision (worker): keine offene Entscheidung. Antwort = neuer step 'answer' in executionOf: question offen + Lane endet nicht (waitsForAnswer) + Text eingetippt -> question= löschen, dann bestehender giveBack nach backTo mit Note 'Given back from <lane>: answer to "<frage>": <text>' (gibt an alten Worker per herdr). Warum giveBack: der Weg an den Worker existiert schon und schreibt die Note. Ohne Text: needAnswer statt noNext, auch wenn acceptTo null. seen-Liste ausgeblendet solange waitsForAnswer.
- **2026-11-27 20:39 · Sascha** — critique: none. Form passt: answer() reuses giveBack (move, note, handover an alten Worker) statt eigenem Pfad; step in executionOf neben den anderen. Bewusst gelassen: die Antwort geht nach backTo (in-progress), nicht in die Lane, aus der der Worker fragte (z.B. brainstorm) — jaira führt die Herkunftslane nicht; der Worker liest seinen Lane aus status. Discuss/Explain-accept bleibt als zweiter Weg (bewegt nach vorn), der needAnswer-Text nennt nur noch das Feld.
- **2026-11-27 20:41 · Sascha** — waiting: bughunt (background subagent) since 21:41
- **2026-11-27 20:42 · Sascha** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/board.ts:1233 (executionOf) with register.tsx:1549 and register.tsx:2716 - class: a ticket that waits on its question AND has a review note marked fix takes the giveBack path, not answer: the field is labelled 'Your answer', but what is typed goes into giveBackWith as an extra 'dazu:' line under the fix heading, the question is never cleared (no jaira set question=), and the note does not say it answers the question; the worker then reads the question as still open and, back in human, the cockpit asks for the same answer again. Reachable whenever w.gaps is non-empty, i.e. a ticket reviewed in an earlier round that comes back to human with a question; example: one note marked fix, typed 'yes, option b', waiting=true -> {step:'giveBack', fix:1}; ran: node --experimental-strip-types scratch calling executionOf(items[now], 'yes, option b', 'review', 0, true) -> giveBack
- **2026-11-27 20:42 · Sascha** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/texts.ts:93 (needAnswer) with register.tsx:1543 - class: a surface without an Input (execute called with canType=false, register.tsx:2634) and a ticket waiting on its question: executionOf can only return needAnswer there (nothing can be typed), and the refusal now says 'type it in the field' - a field that does not exist; the old text pointed to Discuss / an [accept] option under Explain, which work without a field. needWhy has a canType variant (blockerNeedsTyping), needAnswer has none; example: mobile surface, human-lane ticket with a question, Execute -> dead-end instruction; ran: read (execute passes canType only to needWhy)
- **2026-11-27 20:42 · Sascha** — hunt report: 2 findings, both change, neither a blocker, as bughunt notes on PYR4JK | read whole: the Acceptance pane path in hooks/register.tsx (openAcceptance, deciding, execute, answer, giveBack/givingBack, giveBackWith, acceptWith, stillIn, once, pane render), hooks/board.ts (Waiting, snapshotOf, executionOf, waitsForAnswer, decisionLines), hooks/texts.ts, release note; tests/board.test.ts only as the test file | not checked: a live cockpit run (no Herdr/claude UI here); the stale-w race in answer() (no stillIn, unlike giveBackWith) was dropped - the pane re-reads w from boardAtom each render and the same gap already exists on the plain giveBack path
- **2026-11-27 20:42 · Sascha** — bughunt: done at e0d480a8e - 2 findings (2 change, 0 pre-existing) in 5 modules, first hunt; concurrency, error paths, input, permissions checked; files: core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/PYR4JK.md
- **2026-11-27 20:45 · Sascha** — fix round 2 (bughunt): (1) Antwort neben einer fix-Note: execute behandelt giveBack und answer gleich — getippter Text bei waitsForAnswer ist immer 'answer to "…"', question= vor den Decisions (giveBackWith bekommt first-Callback); answer() entfällt, clearQuestion() übernimmt. (2) needAnswer ohne Feld (canType=false): eigener Text needAnswerUntyped, verweist auf Discuss/Explain. Test: 'an answer typed beside a note marked fix…'. Mobile-Pfad nur gelesen, kein Test-Harness dafür.
- **2026-11-27 20:45 · Sascha** — critique round 2 (one pass, 49 lines): both bughunt findings addressed in 5ca1382c7 — answer beside a fix note clears the question and says 'answer to'; canType=false gets needAnswerUntyped. Nothing new.
- **2026-11-27 20:46 · Sascha** — waiting: bughunt round 2 (background subagent) since 21:46
- **2026-11-27 20:48 · Sascha** — hunt report: 0 findings, nothing written to the board | read as diff since e0d480a8e: core/mod/builtin/jaira-cockpit/hooks/register.tsx (execute, clearQuestion, giveBackWith, givingBack, stillIn, deciding, Discuss button), hooks/texts.ts, tests/board.test.ts; board.ts and core/release/notes/PYR4JK.md unchanged since then | not checked live on the mobile surface (canType=false); the needAnswerUntyped text was checked against the code only - Discuss shows whenever question is set, [accept] under Explain clears the question; claude plugin test: 591 pass
