---
id: 01M5ZJ4KA5N2PH7Q4WY4F1V59P
title: "Ein Review-Verdikt 'zurueckschicken' laesst das Ticket nicht nach signoff"
status: critique
ready: true
creator: BeMuCa
goal: "Steht im review-verdict der Review-Lane 'send back' (oder das Gegenstueck des Boards), nimmt jaira move aus review nur noch in-progress als Ziel; signoff wird mit Exit 3 abgelehnt und nennt das Verdikt."
context: |-
  Was schief ist: die Review-Lane schrieb als Verdikt 'zurueckschicken', das Ticket landete trotzdem in signoff (Ticket 25DZW3 auf einem anderen Board, Teamlead-Session 25.-27.10.2026).
  Heute prueft kein Gate, ob Verdikt und Ziellane zusammenpassen: review muss review-verdict liefern, aber move liest den Inhalt nicht.
  Der Mensch sah ein Ticket zur Abnahme, das der Reviewer selbst nicht freigegeben hatte.
  Bekannt: die Gates leben in core/gate/gate.go; welche Werte review-verdict annimmt, steht in der Lane-Datei review der Boards (grep review-verdict core/lane/).
  Nicht gewollt: ein Freitext-Parser. Ein fester kleiner Wertevorrat (z. B. accept, send-back) reicht; alles andere ist kein Verdikt.
  Zusammenhang: L312/L327 der Feedback-Liste vom 27.10.2026.
definition-of-done: "Ein Move aus review nach signoff mit review-verdict 'send back' wird mit Exit 3 abgelehnt; die Meldung nennt das Verdikt und die erlaubte Ziellane"
tags:
  - gates
blocked-by:
  - 01M6M7EA5EG5TWJSTVNQH6CD33
related:
  - 01M7WWZWSDHK37DGGHJFDF0YD6
commits: []
created-at: 2026-10-27T20:51:56Z
updated-at: 2026-11-21T20:04:03Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-238286
claimed-at: 2026-11-21T19:47:43Z
outcome-what: "review-verdict now starts with a lane-declared word (accept / send-back); jaira move reads it and refuses a forward move out of review on send-back or on no known word, exit 3, naming the verdict and the lane to send back to; show --for-lane review carries the words as verdicts"
outcome-why: "a review that wrote 'send back' still let its ticket land in signoff, so a person saw unapproved work as ready"
outcome-resolves: "TestCLISendBackVerdictRefusesSignoff, TestCLISendBackVerdictGoesBack, TestVerdictWordClasses, TestLaneWithoutVerdictWordsReadsNothing"
review-summary: "none"
round: 3
back-from: testing
back-by: lane
back-count: 1
review-gaps: "inlined the one-use test helper mustLane in core/gate/verdict_test.go; searched for an existing verdict/word normaliser (none - lane.normVerdict is the only one) and for a quote-list helper (quoteAll now serves both messages); no dead code: VerdictNone is the empty-verdict return, the output contract refuses that case"
test-verdict: "pass: go test ./... green earlier, -race on gate/lane/cli/tui RC=0 after the last commits; DoD 1-6 verified in the tree; scratch board: send back -> signoff exit 3 naming send-back and in-progress, ticket stays in review; free text 'looks fine' exit 3; for-lane review carries verdicts accept/send-back; two parallel moves (signoff + in-progress) - signoff refused exit 3, in-progress landed"
hunt-verdict: "fail: bughunt done at 5b8824110, 2 change (send-back bypass via blocked - blocker; a lane with only one verdict list never lets a ticket on), 1 pre-existing (a move that skips review is not checked - note for the person)"
---

# Ein Review-Verdikt 'zurueckschicken' laesst das Ticket nicht nach signoff

## Definition of Done

- [x] Ein Move aus review nach signoff mit review-verdict 'send back' wird mit Exit 3 abgelehnt; die Meldung nennt das Verdikt und die erlaubte Ziellane
  proof: core/gate/gate.go:486 verdictAgainstMove (every move landing past review, parking included); TestCLISendBackVerdictRefusesSignoff, TestSendBackVerdictRefusesSignoff, TestSendBackVerdictHasNoWayAround
- [x] Ein Move aus review nach in-progress mit 'send back' geht ohne --force
  proof: TestCLISendBackVerdictGoesBack; TestSendBackVerdictLetsTheTicketGoBack
- [x] Die Werte, die review-verdict annimmt, stehen in der Lane-Datei und in jaira show --for-lane review
  proof: core/lane/builtin/40-review.md:11-12 verdict-accept/verdict-send-back + prompt step 3; internal/cli/flow.go:753 payload verdicts; TestForLaneReviewNamesTheVerdicts
- [x] Release-Notiz in core/release/notes/<ID>.md
  proof: core/release/notes/F1V59P.md
- [x] Verdikt-Erkennung ueber die Klasse (Tabellentest): Gross/Klein, 'send back' mit Leerzeichen, fuehrende Leerzeichen, ohne Doppelpunkt, unbekanntes Wort, leer
  proof: core/gate/verdict_test.go TestVerdictWordClasses
- [x] Eine Lane-Datei ohne Verdikt-Listen (aelteres review.md) aendert nichts am heutigen Verhalten (Test)
  proof: core/gate/verdict_test.go TestLaneWithoutVerdictWordsReadsNothing

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-20 16:49 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-21 19:47 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock
- **2026-11-21 19:49 · Alexander Sacharov** — decision (worker, settled from code+context, no open question): (1) verdict words live in the lane file as two lists, verdict-accept: [accept] and verdict-send-back: [send-back]; review-verdict starts with the word, then ':' and the reason ('send-back: DoD 2 unmet'); compare is case-insensitive, space==dash so 'send back' counts. (2) the gate acts only on a lane that declares the lists, so a board whose review.md predates this keeps today's behaviour until 'jaira update --lanes'. (3) on such a lane a forward, non-parking move with a send-back word is refused (exit 3, new code verdict_sends_back) naming the verdict and rejects-to (review gets rejects-to: in-progress); a forward move with a verdict that is neither word is refused too (verdict_unknown) - 'alles andere ist kein Verdikt'. Backward moves and parking stay free; a move in the TUI (Interactive) is a person deciding and is exempt, --force stays the CLI override. (4) show --for-lane carries the lists as 'verdicts'. Not done: validating at 'jaira set' - the gate at move is the one rail.
- **2026-11-21 19:56 · Alexander Sacharov** — in-progress: dropped rejects-to: in-progress from review.md — TestRejectsToDefaultsEmpty keeps default lanes free of back edges on purpose; the refusal names the work lane instead (last earlier lane producing outcome-what, gate.go sendBackTo). 'Accept.' with a full stop reads as accept (mode_test writes that). Cockpit digest/board.ts not touched: a send-back verdict can no longer reach signoff, so the Explain heuristic there never sees one. Not exempted: ByPerson (typed at a terminal) — --force is the person's override, only a move in the TUI is exempt.
- **2026-11-21 19:57 · Alexander Sacharov** — critique: gate.go sendBackTo duplicates quoteAll for RejectsTo; reuse quoteAll. Rest let stand: lane-declared words are what DoD 3 and the context's 'Gegenstueck des Boards' ask for; the unknown-verdict refusal is the context's 'alles andere ist kein Verdikt'; listOrScalar is the existing rejectsTo reader generalised, not a new pattern.
- **2026-11-21 19:57 · Alexander Sacharov** — critique round 2: the one finding (quoteAll reuse) is addressed in db31aea3a; nothing else.
- **2026-11-21 20:01 · Alexander Sacharov** — bughunt: change: blocker: core/gate/gate.go:486 - class: the verdict gate only judges a move whose source lane is the reviewing lane, and the parking lane (blocked) records no lane to return to, so any detour through it drops the verdict; example: review-verdict="send-back: DoD 1 unmet", 'jaira move X --to signoff' exit 3, then 'jaira move X --to blocked --reason w' exit 0 and 'jaira move X --to signoff' exit 0 - ticket sits in signoff with a send-back verdict, by an agent, no --force. The older human-exit rail (gate.go:405) has no !parking exemption and does refuse signoff->blocked, so only the new rule leaks; ran: branch binary on a scratch board (jaira init, review lane from this branch)
- **2026-11-21 20:01 · Alexander Sacharov** — bughunt: change: core/lane/lane.go:366 (parse) with core/lane/verdict.go ReadVerdict and core/gate/gate.go verdictAgainstMove - class: a lane file that declares only one of the two word lists is accepted, and with verdict-accept empty every verdict that is not a send-back word reads as unknown, so no forward move out of that lane is ever possible without --force; the refusal says 'begin it with (none)' and 'show --for-lane' hands out "accept":null; example: review.md with only 'verdict-send-back: [send-back]', review-verdict="accept: fine" / "lgtm" -> 'jaira move X --to signoff' exit 3 verdict_unknown every time; not a blocker: needs a hand-edited lane file, but the release note invites one; ran: branch binary on a scratch board
- **2026-11-21 20:01 · Alexander Sacharov** — bughunt: pre-existing: core/gate/gate.go:470 - class: a forward move that skips a lane checks only the lane being left, never the lanes jumped over, so a lane's output contract (and now its verdict) is bypassed by skipping it; example: ticket in in-progress, 'jaira move X --to signoff --what w --why y --resolves r' exit 0 - review never ran, no review-verdict at all; blame: 610e33a41a (2026-08-24, before the ticket); ran: branch binary on a scratch board
- **2026-11-21 20:01 · Alexander Sacharov** — bughunt: done at 5b8824110 - 3 findings (2 change, 1 blocker; 1 pre-existing) in 6 modules, first hunt; concurrency, error paths, input, permissions checked; files: .claude/skills/jaira/SKILL.md .jaira/lanes/review.md core/gate/gate.go core/gate/verdict_test.go core/lane/builtin/40-review.md core/lane/lane.go core/lane/verdict.go core/lane/verdict_test.go core/release/notes/F1V59P.md core/role/jaira/SKILL.md core/shipped/shipped.sha256 internal/cli/flow.go internal/cli/verdict_test.go scripts/demoworld.sh
- **2026-11-21 20:02 · Alexander Sacharov** — testing: back to in-progress for the two bughunt change findings. Fix of the class: the send-back rule must hold for every move that lands past the verdict lane, not only a move out of it — blocked is drawn last, so blocked→signoff is not even 'forward'. The unknown-word rule stays on the move out of the verdict lane. A lane declaring only one of the two lists is refused at load.
- **2026-11-21 20:03 · Alexander Sacharov** — fix round 2: verdictAgainstMove now walks every verdict lane the target lies past (send-back refused from review, blocked or a skip ahead; unknown word refused out of review or out of parking); lane parse refuses one list alone, so quoteAll's '(none)' branch went. For the person (pre-existing, bughunt finding 3, gate.go ~470 since 610e33a41a): a move that skips review entirely (in-progress -> signoff) is not checked for review outputs at all - not this ticket's scope; decide note vs. own ticket.
