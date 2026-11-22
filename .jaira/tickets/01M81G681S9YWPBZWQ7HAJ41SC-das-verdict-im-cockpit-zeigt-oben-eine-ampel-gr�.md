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
updated-at: 2026-11-22T11:59:35Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1232016
claimed-at: 2026-11-22T11:29:15Z
outcome-what: "Verdict-Pane des Cockpits zeigt als erste Zeile eine Ampel (grün/gelb/rot) mit dem ersten Grund und '+N more'"
outcome-why: "Alex soll ohne Urteil, Gaps und Notizen zu lesen sehen, ob er annehmen kann"
outcome-resolves: "AJ41SC: fehlende Ampel im Verdict"
review-summary: "Das Verdict-Pane des Cockpits beginnt jetzt mit einer Zeile '● green|yellow|red: <erster Grund> (+N more)'. board.ts lightOf rechnet die Farbe aus dem, was schon am Ticket steht: rot bei send-back-Urteil, blocker:-Zeilen in review-gaps, einem Bughunt, der nicht mehr HEAD deckt, oder einem offenen DoD-Punkt ausser 'seen by a person:' (das Gate des Terminal-Lanes); gelb bei kein/altem Review, Freitext-Urteil, Vorbehalt, note:-Zeilen, offenem 'seen by a person:'-Punkt, offener Frage, keinem oder unbekanntem Bughunt; sonst gruen. Welche Richtung ein Urteil hat, liest das Cockpit nicht mehr aus eigenen Wortlisten, sondern aus den Woertern der Lane: jaira lanes --json gibt dafuer neu 'verdicts' {accept, send_back} aus (tickets.go), und wayOf portiert lane.SendsBackBy; nur ohne diese Woerter (alte jaira) bleibt die alte REFUSAL-Heuristik. huntsAtom haelt jetzt {state, at}, damit 'nie gehuntet' (gelb) von 'Hunt nicht auf HEAD' (rot) getrennt wird. openOf zog von overview.ts nach board.ts. Tests pro Farbe in board.test.ts, Go-Test fuer verdicts in context_test.go, Release-Note mit zwei Zeilen."
round: 5
back-from: review
back-by: lane
back-count: 1
review-gaps: |-
  note: ein Ticket in human mit offener Frage hat fast immer offene DoD-Punkte und oft keinen aktuellen Hunt - das Licht zeigt dort rot 'N DoD items open' (Rat: zurueckgeben), obwohl der Mensch nur die Frage beantworten soll; Fix: in einem Lane ohne terminalen next (acceptTo null) die Frage als ersten Grund zeigen oder dodOpen/staleHunt dort gelb werten (board.ts lightOf)
  note: snapshotOf sammelt verdicts aller Lanes eines Boards in eine Liste und prueft send_back vor accept ueber alle - heute deklariert nur review Woerter, mit einer zweiten Lane koennte ein Wort der einen das Urteil der anderen kippen; Fix: die Woerter des Lanes nehmen, aus dem das Urteil kommt (board.ts:135)
  note: Hunts werden nur fuer die im Band gezeigten Tickets gelesen (register.tsx shownIds); ein Verdict-Pane fuer ein Ticket ausserhalb zeigt gelb 'bughunt unknown' - sichere Richtung, aber nicht die Wahrheit
  note: nicht dieses Ticket: der Worktree traegt uncommittete Aenderungen an .jaira/lanes/*.md, AGENTS.md, CLAUDE.md (Lane-Texte aus einem jaira update) - nicht in einen AJ41SC-Commit mitnehmen
test-verdict: "pass: round 3 - go test ./... RC=0, cockpit 'claude plugin test .' 301/301 RC=0. Binary on scratch repo: lanes --json review verdicts {accept:[accept], send_back:[send-back]}; brief bughunt back to null for a ticket with no hunt (as on master). Light: null hunt -> yellow 'no bughunt yet', hunt at a commit not covering HEAD -> red 'bughunt not on HEAD' (tests). DoD 1-3 verified, proofs re-checked; DoD 4 'seen by a person' left for Alex. Change writes nothing."
hunt-verdict: "pass: bughunt done at c347e752b, 0 change, 0 pre-existing; round-2 blocker (none state) closed"
review-verdict: "accept: DoD 1-3 erfuellt, lightOf deckt die Ticketdaten ab und liest die Urteilsrichtung wie jairas Gate; 301 Cockpit-Tests und Go-Tests gruen; vier note:-Ränder, kein Blocker"
review-round: 4
review-verdict-gaps: |-
  note: ein Ticket in human mit offener Frage hat fast immer offene DoD-Punkte und oft keinen aktuellen Hunt - das Licht zeigt dort rot 'N DoD items open' (Rat: zurueckgeben), obwohl der Mensch nur die Frage beantworten soll; Fix: in einem Lane ohne terminalen next (acceptTo null) die Frage als ersten Grund zeigen oder dodOpen/staleHunt dort gelb werten (board.ts lightOf)
  note: snapshotOf sammelt verdicts aller Lanes eines Boards in eine Liste und prueft send_back vor accept ueber alle - heute deklariert nur review Woerter, mit einer zweiten Lane koennte ein Wort der einen das Urteil der anderen kippen; Fix: die Woerter des Lanes nehmen, aus dem das Urteil kommt (board.ts:135)
  note: Hunts werden nur fuer die im Band gezeigten Tickets gelesen (register.tsx shownIds); ein Verdict-Pane fuer ein Ticket ausserhalb zeigt gelb 'bughunt unknown' - sichere Richtung, aber nicht die Wahrheit
  note: nicht dieses Ticket: der Worktree traegt uncommittete Aenderungen an .jaira/lanes/*.md, AGENTS.md, CLAUDE.md (Lane-Texte aus einem jaira update) - nicht in einen AJ41SC-Commit mitnehmen
review-check: |-
  1. cd /home/alex/projects/.worktrees/jaira-aj41sc && go install ./cmd/jaira
  2. claude --plugin-dir /home/alex/projects/.worktrees/jaira-aj41sc/core/mod/builtin/jaira-cockpit (im Repo)
  3. Im Cockpit-Band AJ41SC (in signoff) waehlen und das Verdict-Pane oeffnen
  4. Erste Zeile: '● yellow: 4 review notes to weigh (+N more)' in Gelb, vor dem Titel - das ist DoD 4 (seen by a person)
  5. Optional rot pruefen: jaira lanes --json | grep -A3 verdicts zeigt accept/send_back; ein Ticket mit 'send-back: ...'-Urteil im Pane zeigt '● red: the review gives it back'
  6. Ohne UI: cd core/mod/builtin/jaira-cockpit && claude plugin test . -> 301 pass, Block 'the verdict light'
question: ""
---

# Das Verdict im Cockpit zeigt oben eine Ampel: grün annehmen, gelb ansehen, rot zurückgeben

## Definition of Done

- [x] Das Verdict zeigt als erste Zeile eine Ampel (grün/gelb/rot) mit einem kurzen Grund
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1945; test 'a ticket in the band is one row' expects '● red: 1 blocker in the review (+4 more)' as first pane line
- [x] Die Farbe wird aus vorhandenen Ticketdaten berechnet (review-verdict, Bughunt auf HEAD, Gates, Notizen, seen by a person:), Regeln mit Test pro Farbe
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts:260 lightOf; core/mod/builtin/jaira-cockpit/tests/board.test.ts 'the verdict light' (green/yellow/red, declared verdict words, no hunt vs stale hunt, question in human); verdicts from jaira lanes --json (internal/cli/context_test.go)
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
- **2026-11-22 11:51 · Alexander Sacharov** — critique pass 4 (nur Fix c347e752b): huntsAtom haelt jetzt jairas {state, at} statt nur state - noetig, um 'keiner fertig' von 'nicht auf HEAD' zu trennen; Go-Seite zurueck auf master-Stand bis auf lanes --json verdicts. Nichts weiter.
- **2026-11-22 11:53 · Alexander Sacharov** — bughunt: done at c347e752b - 0 findings in 6 modules, follow-up hunt since f7dfca8a4 (diff only, no new files; uncommitted .jaira/lanes, AGENTS.md, CLAUDE.md are lane text, read as diff); concurrency, error paths, input, permissions checked; ran: claude plugin test . (301 pass), go test ./core/gitrepo ./internal/cli ./core/shipped (ok); files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/optimize.md .jaira/lanes/pre-process.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/gitrepo/hunt.go core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/overview.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/AJ41SC.md core/shipped/shipped.sha256 internal/cli/context_test.go internal/cli/loopbrake_test.go internal/cli/tickets.go
- **2026-11-22 11:56 · Alexander Sacharov** — review: accept geschrieben, aber move nach signoff verweigert: der letzte Hunt (c347e752b) las auch die uncommitteten .jaira/lanes/*.md, AGENTS.md, CLAUDE.md, die nicht zu diesem Ticket gehoeren (Lane-Texte aus jaira update). Frage: wie weiter?
1. (empfohlen) Alex committet die Lane-Texte separat (z.B. chore: jaira update lanes) oder verwirft sie mit git restore, dann move --to signoff ohne neuen Hunt - kein Code dabei
2. Alex: jaira move AJ41SC --to signoff --force - schnell, Worktree bleibt dreckig
3. Lane-Texte in einen AJ41SC-Commit - vermischt fremde Aenderung mit diesem Ticket
- **2026-11-22 11:57 · Alexander Sacharov** — decision (parent/teamlead 11:57): Lane-Texte verworfen (git checkout), sie kommen ueber GPSAR2 nach master. Gap 1 jetzt fixen: in einer Lane ohne terminalen next (acceptTo null, z.B. human) mit offener Frage entscheidet die Frage das Licht - gelb, Frage als erster Grund, rote Gruende nur mitgezaehlt. Uebrige Review-Notizen bleiben bekannte Grenzen.
