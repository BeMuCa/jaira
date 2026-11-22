---
id: 01M81G681S9YWPBZWQ7HAJ41SC
title: "Das Verdict im Cockpit zeigt oben eine Ampel: grün annehmen, gelb ansehen, rot zurückgeben"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Alex sieht im Verdict mit einem Blick, ob er annehmen kann, und warum nicht."
context: "Problem: im Verdict des Cockpits sieht Alex nicht auf einen Blick, ob er annehmen kann. Er muss Urteil, Gaps und Notizen lesen, und er liest langsam.\n- Wunsch Alex 22.11.2026: \"mir fehlt im Verdict eine Ampel, damit ich schneller annehmen kann\".\n- Das Verdict zeigt seit PDZ444 Titel + Ziel und das Review-Urteil; Code: core/mod/builtin/jaira-cockpit/hooks/register.tsx, digest.ts, board.ts, texts.ts.\n- Vorschlag für die Farben, die Daten liegen schon am Ticket:\n  - grün: review-verdict accept, Bughunt-Abschlusszeile auf HEAD, Gates ok, keine Notiz als Blocker sortiert, kein offener 'seen by a person:'-Punkt\n  - gelb: accept, aber eine 'seen by a person:'-Punkt offen oder note:-Zeilen in review-gaps (kleine Ränder/bekannte Grenzen)\n  - rot: review gibt zurück, Bughunt fehlt oder ist alt (nicht auf HEAD), ein Gate verweigert, oder eine Blocker-Notiz\n- Die Ampel ist eine Zeile ganz oben im Verdict, vor dem Titel, mit einem Wort Grund (\"rot: Bughunt nicht auf HEAD\").\n- Kein neues Feld am Ticket: die Farbe wird aus vorhandenen Feldern berechnet."
definition-of-done: Das Verdict zeigt als erste Zeile eine Ampel (grün/gelb/rot) mit einem kurzen Grund
tags:
  - acceptance
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-11-22T11:27:46Z
updated-at: 2026-11-22T11:51:18Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1232016
claimed-at: 2026-11-22T11:29:15Z
outcome-what: "Verdict-Pane des Cockpits zeigt als erste Zeile eine Ampel (grün/gelb/rot) mit dem ersten Grund und '+N more'"
outcome-why: "Alex soll ohne Urteil, Gaps und Notizen zu lesen sehen, ob er annehmen kann"
outcome-resolves: "AJ41SC: fehlende Ampel im Verdict"
review-summary: none
round: 4
back-from: testing
back-by: lane
back-count: 2
review-gaps: "round 2: nichts entfernt. Geprueft: wayOf vs core/lane/verdict.go (TS kann Go nicht rufen, bewusst portiert); BriefHunt vs HuntCheck (wiederverwendet, keine Kopie); freeVerdict-Text neu, kein vorhandener Text passt. Kein toter Code."
test-verdict: "fail: suite green; the hunt's blocker gives a wrong (green) light"
hunt-verdict: "fail: 1 change blocker - BriefHunt says none whenever the ticket's commits are not visible from the checkout (master, uncommitted work): a ticket with code and no hunt lights green"
---

# Das Verdict im Cockpit zeigt oben eine Ampel: grün annehmen, gelb ansehen, rot zurückgeben

## Definition of Done

- [x] Das Verdict zeigt als erste Zeile eine Ampel (grün/gelb/rot) mit einem kurzen Grund
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1945; test 'a ticket in the band is one row' expects '● red: 1 blocker in the review (+4 more)' as first pane line
- [x] Die Farbe wird aus vorhandenen Ticketdaten berechnet (review-verdict, Bughunt auf HEAD, Gates, Notizen, seen by a person:), Regeln mit Test pro Farbe
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts:260 lightOf; core/mod/builtin/jaira-cockpit/tests/board.test.ts 'the verdict light' (green/yellow/red, declared verdict words, no hunt vs stale hunt); verdicts from jaira lanes --json (internal/cli/context_test.go)
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/AJ41SC.md (2 lines); core/shipped/shipped.sha256 via scripts/shipped-history.sh --write; go test ./... RC=0
- [ ] seen by a person: Alex sieht die Ampel im laufenden Cockpit an einem Ticket in signoff

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 11:29 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 11:30 · Alexander Sacharov** — decision (worker, no open question): Ampel-Regeln. rot: Review lehnt ab (REFUSAL), Blocker in review-gaps, Bughunt fail (fehlt/nicht auf HEAD), ein DoD-Punkt ausser 'seen by a person:' offen (= das Gate des Terminal-Lanes). gelb: Review mit Vorbehalt, note:-Zeilen, offener 'seen by a person:'-Punkt, offene Frage, Review aus frueherer Runde, Notizen unbekannt, Bughunt unknown/unbekannt (aeltere jaira), noch kein Review. gruen: sonst. 'Gates' = DoD-Gate, weil der Worker-Statuszeile-gates-Wert nicht am Ticket liegt; Texte englisch wie texts.ts. Erster Grund in Reihenfolge rot vor gelb wird gezeigt, plus '+N' weitere.
- **2026-11-22 11:34 · Alexander Sacharov** — critique: (1) openOf aus overview.ts nach board.ts ziehen und in dodOpenOf benutzen statt den Filter zu kopieren; (2) gelber Fallback ohne dod_items bekommt eine eigene kind 'dodIncomplete' statt dodOpen n=0.
- **2026-11-22 11:36 · Alexander Sacharov** — critique pass 2: beide Befunde beantwortet (openOf in board.ts, kind dodIncomplete); nichts weiter.
- **2026-11-22 11:40 · Alexander Sacharov** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/board.ts:240 - class: the light decides 'the review gives it back' by its own word list (REFUSAL, board.ts:196: request|reject|refuse|fail|zurück|ablehn|nicht|not) instead of the review lane's declared verdict words (verdict-send-back: [send-back], verdict-accept: [accept]; jaira lanes --json carries them as send_back, internal/cli/flow.go:755; the gate reads them in core/lane/verdict.go:40). The lane's own send-back word is not in the list, and plain negations are; example: verdict 'send-back: DoD 2 - an empty export still writes a header' with gaps 'none' lights GREEN 'accept' (the same verdict with a 'note:' gap lights yellow, never red); verdict 'Not a single defect found; accept.' or 'Nicht zu beanstanden - accept' lights RED 'the review gives it back'. Reachable whenever a send-back verdict sits in a person's lane: a move made in the board (Interactive skips verdictAgainstMove) or with --force. REFUSAL itself is older (75637606), but lightOf is the new caller that turns it into a colour; ran: node --experimental-strip-types scratch script calling snapshotOf+lightOf(w,'ok') on these verdicts (outputs green / red as stated)
- **2026-11-22 11:40 · Alexander Sacharov** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/board.ts:243 - class: a ticket that owes no bug hunt because it changed no code is lit red 'no bughunt on HEAD'. jaira show --brief --json reports bughunt:null for every ticket without a finished hunt (internal/cli/tickets.go:830 calls LastHunt, which never checks whether code changed), huntOf maps null to fail with at='' (board.ts:856), register.tsx:398 skips its ancestor check for at='', and lightOf pushes red noHunt. The gate does not ask for a hunt there (core/gitrepo/hunt.go HuntCheck: CodeChanged false when the ticket's commits touch only board data or there are none, so huntAgainstMove passes). Example: a decision or research ticket moved into human with --question, or into signoff, with no code commits: Verdict light red 'no bughunt on HEAD (+1 more)' although nothing is owed and the move would pass. The band's dim cross for the same state is older (d16585f9); the red 'give it back' meaning is new with the light; ran: fresh repo, jaira init, jaira create 'decide x' -> jaira show XV0ZCB --brief --json prints "bughunt":null (exit 0)
- **2026-11-22 11:40 · Alexander Sacharov** — bughunt: done at c2eab0fab - 2 findings (2 change, 0 pre-existing) in 5 modules, first hunt; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/optimize.md .jaira/lanes/pre-process.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/overview.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/AJ41SC.md core/shipped/shipped.sha256
- **2026-11-22 11:42 · Alexander Sacharov** — decision (round 2, nach bughunt): (1) jaira lanes --json gibt 'verdicts' {accept, send_back} einer Lane aus, die sie deklariert (wie der --for-lane-Payload); das Cockpit liest den Verdict damit wie lane.SendsBackBy: send-back -> rot, accept -> kein Grund, anderes Freitext auf einem Board mit Woertern -> gelb 'verdict outside the lane's words'; nur ohne deklarierte Woerter (alte jaira) bleibt die REFUSAL-Heuristik. (2) show --brief: bughunt-State 'none', wenn kein Hunt fertig ist und HuntCheck sagt, die Commits aenderten keinen Code; das Cockpit zaehlt 'none' wie skipped (neutral, Strich im Band). Kleinere Variante als Wortlisten im Cockpit zu pflegen (memory: edges grow -> reuse jaira code).
- **2026-11-22 11:46 · Alexander Sacharov** — round 2 fix: beide bughunt-Blocker ueber jaira selbst geloest (lanes --json verdicts, brief bughunt none via HuntCheck). Band zeigt none wie skipped als '–'. Explain-Prompt kennt none.
- **2026-11-22 11:46 · Alexander Sacharov** — critique pass 3 (nur der Fix): wayOf in board.ts portiert lane.SendsBackBy (8 Zeilen) - gelassen, weil das Cockpit nur jaira-JSON liest und ein berechnetes Ticketfeld in list --json mehr Format waere; BriefHunt ruft LastHunt im nil-Fall zweimal (via HuntCheck), nur Notizen-Scan, gelassen.
- **2026-11-22 11:48 · Alexander Sacharov** — bughunt: change: blocker: core/gitrepo/hunt.go:63-71 (BriefHunt) with core/gitrepo/derive.go:42,66 and core/mod/builtin/jaira-cockpit/hooks/board.ts:268,284 - class: 'no code change visible from this checkout' is reported as 'nothing to hunt' (state none). HuntCheck finds commits with git log on HEAD only and returns CodeChanged=false when it finds none or git fails, so any checkout that does not hold the ticket's code commits yields none: the main checkout on master while the code sits on an unmerged feat branch (the cockpit falls back to board.repo whenever worktreeFor finds no tree for the id, register.tsx:390), and a worktree whose code is still uncommitted (conversational mode). lightOf treats none like skipped: no noHunt (red) and no huntUnknown (yellow), so a ticket with code and no hunt can show a green light and a '–' hunt mark; before f7dfca8a4 the same case was bughunt null -> fail -> red. example: ticket with code committed on feat/<id>, no hunt note, cockpit reading from master -> green. ran: scratch repo, binary built from HEAD: uncommitted code.go -> {"state":"none"}; committed on feat branch -> null; same ticket from master -> {"state":"none"}
- **2026-11-22 11:48 · Alexander Sacharov** — bughunt: done at f7dfca8a4 - 1 finding (1 change, 0 pre-existing) in 6 modules, follow-up hunt since c2eab0fab (diff for files read before, new files whole); concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/optimize.md .jaira/lanes/pre-process.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/gitrepo/hunt.go core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/overview.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/AJ41SC.md core/shipped/shipped.sha256 internal/cli/context_test.go internal/cli/loopbrake_test.go internal/cli/tickets.go
- **2026-11-22 11:49 · Alexander Sacharov** — decision (round 3, nach bughunt 2): 'none' im Brief wieder raus - HuntCheck sieht Commits nur von HEAD aus, also rät es auf master/uncommitted falsch. Kleinere Form: Cockpit unterscheidet 'kein Hunt fertig' (brief null) von 'Hunt nicht auf HEAD' (fail mit at). Null -> gelb 'no bughunt yet' (ein Ticket ohne Code hat keinen, eines mit Code kommt ohne Hunt nur per --force in signoff); fail mit at -> rot. Weicht von 'rot: Bughunt fehlt' im Kontext ab; Alex kann widersprechen.
