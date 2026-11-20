---
id: 01M7WWZWSDHK37DGGHJFDF0YD6
title: "Die Schleifenbremse steht am Ticket: was ein Blocker ist und dass nach einem Fix kein neuer Bughunt kommt"
status: testing
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Ticket traegt ein Blocker-Kriterium und die Option stop-after-fix; testing, bughunt und review lesen sie, und das Cockpit bietet Explain nur bei einem blocker an."
context: |-
  Heute: ein Teamlead schreibt mehrmals am Tag von Hand "fix -> test -> kein Bughunt -> review, Blocker nur bei falschem Ergebnis oder Absturz". Das hat jedes Mal sofort konvergiert, steht aber nur in Notizen und gilt fuer das naechste Ticket nicht.
  Ohne Schwelle ist jeder Randfall ein Blocker; jeder Fix ist neuer Code, den der naechste Bughunt jagt (das Ziel wandert mit).
  review-gaps hat immer 2-6 note:-Zeilen; das Cockpit zeigt Explain schon bei note: allein (core/mod/builtin/jaira-cockpit/hooks/board.ts:196 explainable/namesGap), deshalb hat jedes Ticket in signoff ein Explain.
  Die Schwelle gehoert in DoD/Plan, gesetzt von pre-process (siehe Ticket zu Randfaellen vor dem Code); der Mensch kann sie aendern.
  Verwandt: H6CD33 (max-rounds einer Schleifen-Lane).
definition-of-done: "Ticket-Felder fuer Blocker-Kriterium und stop-after-fix; jaira show --for-lane gibt sie testing, bughunt und review mit (Test)"
tags:
  - gates
blocked-by:
  - 01M7WWZRD7WNC5NT1PS5TW5FE2
related:
  - 01M6M7EA5EG5TWJSTVNQH6CD33
  - 01M5ZJ4KA5N2PH7Q4WY4F1V59P
  - 01M7WWZRD7WNC5NT1PS5TW5FE2
commits: []
created-at: 2026-11-20T16:35:17Z
updated-at: 2026-11-20T23:18:50Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-91898
claimed-at: 2026-11-20T22:44:37Z
outcome-what: "As before, plus: a send-back made in the board (TUI move and force move) is recorded as a person's, whatever lane it leaves; the cockpit lists review notes only for a verdict of the current round."
outcome-why: "The bug hunt found a person's TUI return from a loop lane counted toward the third-return stop, and optimize's rewritten gaps shown as review notes on a later round."
outcome-resolves: "Both bughunt change findings of ac850761: core/move Request.ByPerson (TestMoveRecordsAPersonAtTheBoard), board.ts notes gated by sameRound (test 'a later round's optimize gaps are no review notes under the row'); go test -race core/move internal/tui RC=0, claude plugin test 215 pass."
review-summary: none
review-gaps: "none - round 2 read only 1789d202 (ByPerson field, one sameRound guard, two tests); nothing duplicated, dead or repeated"
test-verdict: "pass: round 3 over the diff since ac850761 (ancestor of HEAD): go test -race core/move internal/tui core/role core/shipped RC=0, claude plugin test 215 pass; both findings checked with a neighbour: TUI plain and force move both pass ByPerson, CLI moves from loop lanes still record lane (TestMoveRecordsWhoSentItBack), current-round notes still listed, stale and pre-review ones not"
hunt-verdict: "fail: bughunt done at ac850761, 2 change (optimize gaps listed as review notes on a later round; a person's TUI send-back from a loop lane recorded as by lane), 0 pre-existing"
round: 2
---

# Die Schleifenbremse steht am Ticket: was ein Blocker ist und dass nach einem Fix kein neuer Bughunt kommt

## Definition of Done

- [x] Ticket-Felder fuer Blocker-Kriterium und stop-after-fix; jaira show --for-lane gibt sie testing, bughunt und review mit (Test)
  proof: internal/cli/flow.go (blocker_rule/stop_after_fix beside mode) + core/lane/builtin/27-testing.md, 40-review.md, jaira-role-bughunt/SKILL.md read them; TestForLaneCarriesTheLoopBrake, TestBlockerRuleRoundTrip, TestSetRefusesUnknownStopAfterFix
- [x] Mit stop-after-fix laeuft nach einem Fix-Durchgang nur der Testlauf, kein neuer Bughunt
  proof: core/lane/builtin/27-testing.md: stop_after_fix and round > 1 -> no hunt, hunt-verdict 'pass: stop-after-fix'
- [x] Das Cockpit zeigt Explain nur, wenn review-gaps eine blocker:-Zeile hat oder das Urteil einen Vorbehalt nennt; reine note:-Zeilen stehen als Liste ohne Knopf (Test in board.test.ts)
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts explainable/notesOf; tests/board.test.ts 'Explain is offered for a blocker or a reservation; notes alone are listed without a button'
- [x] jaira zeichnet bei jedem Rueckweg die Herkunfts-Lane auf und ob eine Person oder eine Schleifen-Lane zurueckgeschickt hat; jaira show --json gibt beides aus, und die Schleifenzaehlung zaehlt nur Rueckgaben derselben Schleifen-Lane (Test: signoff->in-progress zaehlt nicht)
  proof: core/move/move.go recordBack; TestMoveRecordsWhoSentItBack (signoff->in-progress is person, resets the count), TestMoveCountsOneReturnForTwoAtOnce; show --json 'back'; 20-in-progress.md reads back.count
- [x] Die Explain-Vorlage (core/mod/builtin/jaira-cockpit/hooks/texts.ts:76-80) ordnet jede Note zuerst ein (Defekt / Rand in diesem Ticket / vorbestehend / nur Hinweis), sagt ausdruecklich, wenn nichts davon die Abnahme verlangt, und gibt keine feste Reihenfolge 1) [give back] 2) [accept] mehr vor - heute zwingt die Vorlage 'Warum: why it matters' und give back als erste Option, das liest sich immer wie 'muss geaendert werden'
  proof: core/mod/builtin/jaira-cockpit/hooks/texts.ts explainBySubagent + core/role/builtin/jaira-dispatcher/scripts/spawn.sh explain prompt: classes first, no Warum:, order by recommendation, neutral case 'can be accepted as it is'; test 'without Herdr, Explain asks this session for a subagent...'
- [x] Steht ein Ticket in signoff/human und ist sein review-verdict aus einer frueheren Runde (review_round < round), schweigt das Cockpit nicht: das Band sagt 'Review von Runde N, seitdem geaendert' und zeigt die alten Funde mit Explain; heute blendet sameRound() in board.ts:189 Explain dann aus (Test)
  proof: board.ts reviewedIn + register.tsx t.reviewedIn line; tests/board.test.ts 'a verdict of an earlier round keeps its Explain and says the ticket changed since' (BLK789: round 8, review_round 7)
- [x] Die Explain-Ansicht im Cockpit ist farbig und auf einen Blick lesbar: Statuszeile mit gruenem Haken/rotem Kreuz je Pruefung (DoD, Tests, Blocker, bughunt auf HEAD); jede Note mit farbiger Marke nach Einordnung - rot 'Blocker', gelb 'kleiner Rand', grau/gedimmt 'bekannte Grenze', blau 'nicht aus diesem Ticket'; die Empfehlung fett mit eigener Farbe; die Notiz-Texte normal, nur die Marken farbig; funktioniert im hellen und dunklen Terminal-Thema (Theme-Farben statt fester Hex-Werte, wie das Band heute)
  proof: hooks/board.ts explanationOf (status/note/recommendation parts) + register.tsx pane (success/error ticks, error/warning/dim/suggestion marks, bold claude recommendation); tests 'a classed explanation shows its checks...' and 'a blocker mark is red'
- [x] class: the loop brake reads recorded facts, never prose — every return is recorded by jaira move, whatever path moved it (CLI, cockpit, --force); parking in blocked and back records nothing; a ticket without back-* fields reads as no return (table test)
  proof: core/move/move_test.go TestMoveRecordsWhoSentItBack (parking, old ticket without fields, --force path)
- [x] input classes: blocker-rule with quotes, umlauts, multiple lines and CRLF survives set/show round-trip; stop-after-fix rejects anything but true/false/empty (exit 2); two parallel moves on one ticket leave consistent back-* fields (one refused or both counted)
  proof: internal/cli/loopbrake_test.go TestBlockerRuleRoundTrip (quotes, umlauts, multiline, CRLF, lone CR), TestSetRefusesUnknownStopAfterFix; core/move TestMoveCountsOneReturnForTwoAtOnce
- [x] acceptance rule for the explain parser: an explain note line without a known class mark or status shape is shown as plain text, never dropped - an old note must still read as before
  proof: board.ts partOf falls back to a text part; legacy SHAPED/MARKED explanation tests still pass (214 pass, claude plugin test)

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] schema: fields blocker-rule, stop-after-fix, back-from, back-by, back-count in core/ticket/schema.go; jaira set checks stop-after-fix is true/false/empty like mode
- [x] move: on every send-back write back-from, back-by (person when the from-lane is not agentic or needs a human exit/question, else lane) and back-count (consecutive returns from the same lane since a person's or another lane's return); table test in core/move
- [x] show --json prints the five fields; --for-lane payload carries blocker_rule (ticket's, else the default four classes) and stop_after_fix beside mode, also in the plain-text header; test in internal/cli
- [x] lane prompts + twins: testing skips the hunt in a fix round when stop_after_fix, reads blocker_rule; review reads blocker_rule; bughunt skill marks a change finding that meets the rule; in-progress third-return reads back-from/back-by/back-count instead of counting notes; pre-process may set the two fields; shipped.sha256
- [x] cockpit board.ts: Explain only for a blocker: line or a reservation; note: lines listed under the row without a button; a review of an older round shows 'Review von Runde N, seitdem geaendert' with Explain; tests in board.test.ts
- [x] explain template in texts.ts and spawn.sh: status line, one line per note with [class], one Empfehlung, no Warum:, no fixed give-back-first order, neutral case says it can be accepted
- [x] explain pane: parse status line and class marks, colour them with theme colours (error/warning/dim/suggestion), recommendation bold; old notes without marks render as before; tests
- [x] release note core/release/notes/DF0YD6.md; go test ./... and claude plugin test in the cockpit dir, exit codes checked
- [-] (duplicate of step 5, added by mistake)

## Progress
- **2026-11-20 16:48 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-20 17:05 · Alexander Sacharov** — Aus TW5FE2 (bughunt e12f36d5): der Prompt-Zaehler in in-progress kann eine Rueckgabe der Person aus signoff nicht von einer review-Rueckgabe unterscheiden - jaira speichert die Herkunfts-Lane nicht (show --json unterscheidet sich nur in status/next_lane/round; sendsBack in core/move/move.go:133 erhoeht round bei jedem Rueckweg, updated_by ist derselbe git-Nutzer). Folge: nach zwei echten review-Rueckgaben eskaliert in-progress, statt zu tun, was die Person wollte. Die echte Zaehlung gehoert ins Binary - neuer DoD-Punkt hier.
- **2026-11-20 17:45 · Alexander Sacharov** — Alex 20.11.2026: 'aus dem Text der Beschreibung im Cockpit habe ich immer das Gefuehl, dass man es korrigieren muss.' Ursache in texts.ts:78-80: Explain gibt es nur bei einem Fund, die Vorlage verlangt 'Warum: why it matters' und setzt [give back] als Option 1. Alex schickt deshalb (und weil er keine offenen Raender will) jedes Ticket zurueck.
- **2026-11-20 17:57 · Alexander Sacharov** — Beobachtet 20.11.2026 in einem anderen Projekt: review Runde 7 schickt mit Blocker zurueck, Stopp-Regel nach drei Rueckgaben, Teamlead entscheidet 'Befund richtig, Verhalten so gewollt' und schiebt das Ticket ohne neues Review nach signoff. round=8, review_round=7 -> sameRound() false -> kein Explain, obwohl review-verdict 'Zurueck an in-progress: ein Blocker' und eine blocker:-Zeile auf dem Ticket stehen. Die Person sieht nichts. Verwandt: F1V59P (Verdikt 'zurueck' sperrt signoff).
- **2026-11-20 18:00 · Alexander Sacharov** — Alex 20.11.2026 zur Explain-Form (DoD 5): 'so wie du es hier erklaerst, ist es sofort verstaendlich; die Beschreibung im Cockpit klingt wie eine Katastrophe.' Vorbild ist die Form, die der Teamlead im Chat benutzt:
- eine Kopfzeile mit dem Stand: Lane, DoD x/y, Tests, 'kein Blocker' / 'Blocker: ...', bughunt auf HEAD ja/nein
- eine Tabelle, eine Zeile je Note: | # | was (ein Satz, Alltagssprache) | Einordnung |; Einordnung ist eine von: bekannte Grenze (nicht loesbar, dokumentiert), kleiner Rand (formal ok, eine Runde ~10 min), nicht aus diesem Ticket (wo es erledigt wird), Blocker
- dann genau eine Empfehlung mit Grund in einem Satz, und welche Note sie betrifft
- kein 'Warum: why it matters' pro Note, keine feste Reihenfolge give back zuerst; der neutrale Fall (nur bekannte Grenzen / fremde Notes) sagt ausdruecklich 'kann so angenommen werden'
- **2026-11-20 22:44 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/35112.sock
- **2026-11-20 22:45 · Alexander Sacharov** — decision: feat/df0yd6 builds on origin/feat/tw5fe2 (merged in), because TW5FE2 (blocked_by) is done but not on master yet and changes the same lane prompts (testing, review, bughunt). Once TW5FE2's PR lands, its commits drop out of this branch's diff.
- **2026-11-20 22:47 · Alexander Sacharov** — decision (worker, settled from code): 1) two frontmatter fields: 'blocker-rule' (one sentence, free text) and 'stop-after-fix' (true/false). Empty blocker-rule = the four classes TW5FE2 already put into testing/review prompts (wrong result, data loss, crash, security hole); --for-lane hands both to testing, review and the bughunt role. 2) return origin: move writes 'back-from: <lane>', 'back-by: person|lane' (person = the from-lane is owned by a person, e.g. signoff/human) and 'back-count: <n>' = consecutive returns from that same lane, reset by a return from another lane or a person; show --json prints all three; the in-progress prompt's note-based count is replaced by reading them. 3) explain note format: one line per note with its class in brackets ([blocker] / [kleiner Rand] / [bekannte Grenze] / [nicht aus diesem Ticket]), a status line, one Empfehlung; no 'Warum:' and no fixed [give back]-first order; the cockpit colours the class marks with theme colours. Same template in texts.ts and spawn.sh.
- **2026-11-20 22:47 · Alexander Sacharov** — question: stop-after-fix default for a ticket where nobody set it? 1) on (Recommended) - matches the rule the teamlead writes by hand daily: after a fix round only tests run, no hunt; the field turns it off per ticket. 2) off - TW5FE2's diff-only hunt stays the default; pre-process or a person turns stop-after-fix on per ticket.
- **2026-11-20 22:48 · Alexander Sacharov** — decision: stop-after-fix is off by default (teamlead jaira-64, 2026-11-20). The diff-only hunt after a fix stays the default: it found real defects in fixes (TW5FE2 rebase fallback, JG3REW foreign comment delete, NNV39Y not-ignored). stop-after-fix=true is a deliberate per-ticket opt-out set by pre-process or a person, e.g. for text-only tickets. The brake for endless loops is the count (3 returns from the same lane -> teamlead), not skipping the hunt. Mode stays autonomous.
- **2026-11-20 22:55 · Alexander Sacharov** — in-progress: CRLF in any 'jaira set' value made the write fail ('value position is at end of line') - pre-existing for every field; fixed at the cause in Doc.SetScalar (core/ticket/frontmatter.go), CR/CRLF stored as LF. Concurrency test needs one Store per goroutine: a shared Store races in idIndex, but jaira never shares one (two processes = two stores).
- **2026-11-20 22:57 · Alexander Sacharov** — in-progress pause: Go side and lane prompts done and pushed (fbc56cd1, 8a396712); DoD 1,2,4,8,9 ticked. Left: cockpit (DoD 3,5,6,7,10) = plan steps 5-7, then release-note lines for the cockpit and the full test run (go test ./..., claude plugin test in core/mod/builtin/jaira-cockpit). in-progress prompt no longer counts back: notes; the in-progress twin changed with it.
- **2026-11-20 23:08 · Alexander Sacharov** — in-progress done, moved to critique; the first send of the move was refused as 'written first' by this session's own last plan tick (ref 0c810cc8 = disk minus the move) - a race between two sends of one session, the disk file is the superset.
- **2026-11-20 23:09 · Alexander Sacharov** — critique: none. Let stand: core/move/move.go isPersons repeats the person-lane predicate gate.go:212, web_source.go:256 and snake.go:67 inline too - a Lane method would touch three unrelated files, not this ticket's scope; explain parser accepts German class aliases (kleiner Rand, bekannte Grenze, nicht aus diesem Ticket) because the explain chat answers in the person's language and may translate the marks - three map entries, kept.
- **2026-11-20 23:13 · Alexander Sacharov** — testing: test-verdict pass written; bughunt running in background at ac850761 - a restarted worker waits for its 'bughunt: done at ac850761' line or starts a new hunt.
- **2026-11-20 23:13 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/board.ts:128 - notes are listed from review-gaps whenever a verdict exists, with no sameRound check (the old explainable() had one for exactly this reason). From round 2 on, the optimize lane overwrites review-gaps before review runs again, so a ticket waiting in human with a round-1 verdict shows optimize's removal list ('removed helper fooBar ...') under its row as the review's notes. Trigger: round=2, review_round=1, status human, verdict set, gaps written by optimize; ran: node --experimental-strip-types on snapshotOf with that ticket -> notes = the two optimize lines, reviewedIn 1
- **2026-11-20 23:13 · Alexander Sacharov** — bughunt: change: core/move/move.go:171 - recordBack decides person vs lane only from the lane left (isPersons), never from who moved. A person who sends a ticket back from critique/optimize/testing in the TUI (internal/tui/model.go:2056 move.Move without Interactive) is recorded as back-by=lane and raises back-count, so in-progress can hit its third-return stop on a person's own send-back and the 'by: person, whatever the count' rule in 20-in-progress.md never applies to it. Trigger: TUI move testing->in-progress three times; read (move.Request.Interactive exists but recordBack does not get it; TUI path 47b1737 is old, the classification is new)
- **2026-11-20 23:13 · Alexander Sacharov** — bughunt: done at ac850761 - 2 findings (2 change, 0 pre-existing) in 8 modules, first hunt; concurrency, error paths, input, permissions checked; files: .jaira/lanes/in-progress.md .jaira/lanes/pre-process.md .jaira/lanes/review.md .jaira/lanes/testing.md core/lane/builtin/15-pre-process.md core/lane/builtin/20-in-progress.md core/lane/builtin/27-testing.md core/lane/builtin/40-review.md core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/move/move.go core/move/move_test.go core/release/notes/DF0YD6.md core/role/builtin/jaira-dispatcher/scripts/spawn.sh core/role/builtin/jaira-role-bughunt/SKILL.md core/role/role_test.go core/shipped/shipped.sha256 core/ticket/frontmatter.go core/ticket/schema.go internal/cli/flow.go internal/cli/loopbrake_test.go internal/cli/tickets.go
- **2026-11-20 23:16 · Alexander Sacharov** — testing: class: a send-back's sender read off the lane alone, and review notes read off gaps whatever round wrote them; example: TUI return from testing counted toward the third-return stop, a human-lane ticket with round 2 / review_round 1 listed optimize's lines as review notes; fix: move.Request.ByPerson set by both TUI move sites (Interactive counts too), notes only for a verdict of the current round; tests TestMoveRecordsAPersonAtTheBoard and 'a later round's optimize gaps are no review notes under the row'.
- **2026-11-20 23:16 · Alexander Sacharov** — critique round 2: none - read only 1789d202; ByPerson beside Interactive is the narrow shape: Interactive also opens the human-exit gate, which a TUI move must not get by this change.
- **2026-11-20 23:18 · Alexander Sacharov** — testing round 3: test-verdict pass; bughunt running over the diff since ac850761 at 1789d202.
