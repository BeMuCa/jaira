---
id: 01M5MNWPVYGSYTYCR3PVNJBTVC
title: "Ein offener Plan-Punkt sperrt die Ziellane nicht mehr, wenn die Definition of Done erfuellt ist"
status: signoff
ready: true
creator: BeMuCa
goal: "Ob ein Ticket in die Ziellane darf, entscheidet allein die Definition of Done; ein offener Plan-Schritt sperrt nicht mehr, sondern wird beim gelungenen Move als Hinweis genannt - der Plan ist die Methode, die DoD das Ergebnis."
context: |-
  Beobachtet 22.-23.10.2026 (Feedback Punkt 8): GH0GKD mit DoD 9/9 wurde an done abgelehnt, weil Plan-Schritt 30 offen war.
  Nachgeprueft am 23.10.: core/gate/gate.go:411-424 verweigert jede Ziellane mit CodePlanIncomplete, solange ein Plan-Punkt offen ist. Der Kommentar dort sagt, Plan und DoD seien nicht unabhaengig.
  CLAUDE.md dieses Repos sagt das Gegenteil: 'Ticking the plan does not tick the definition of done' - der Plan ist eine zweite, separate Liste.
  Nicht entschieden: Gate weg (nur Warnung), oder der Worker hakt beim Abschluss den letzten Plan-Schritt selbst ab. Brainstorm entscheidet; die einfachere Variante ist das Gate zu einer Warnung zu machen.
definition-of-done: Ein Move in eine Ziellane mit vollstaendiger DoD und offenem Plan-Punkt wird nicht mehr abgelehnt; die Ausgabe nennt den offenen Plan-Punkt als Hinweis; ein Test deckt den Fall
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-10-23T15:25:53Z
updated-at: 2026-10-26T20:27:30Z
assignee: BeMuCa
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-650718
claimed-at: 2026-10-26T19:43:53Z
outcome-what: "Ein Move in eine Ziellane wird nicht mehr wegen eines offenen Plan-Schritts abgelehnt; nur noch das DoD-Gate entscheidet, bei unvollstaendiger DoD erscheinen nur DoD-Punkte. Offene Plan-Schritte kommen als Hinweis zurueck: gate.OpenPlanHints, move.Result.Hints, 'Note:'-Zeilen vor der Move-Zeile in der CLI, 'warnings' im JSON, in der TUI-Meldung (applyMove, forceMove, signoff accept). CodePlanIncomplete entfernt; Doku (SKILL.md, README.md, docs/AGENTS.md, 50-done.md) und Release-Notiz angepasst."
outcome-why: "Plan und DoD sind zwei getrennte Listen; denselben Haken setzt derselbe Worker, der Plan-Haken ist kein unabhaengiger Beleg. Die Sperre hielt GH0GKD mit DoD 9/9 vor done fest, weil ein erledigter Plan-Schritt nicht abgehakt war."
outcome-resolves: "DoD 1: TestMoveIntoDoneWithOpenPlanStepCarriesAHint + TestUnfinishedPlanDoesNotBlockTerminalLane; DoD 3: TestUnfinishedPlanAddsNothingToDoDRefusal; DoD 4: internal/cli/flow.go, internal/tui/model.go, internal/tui/signoff.go; DoD 5: die vier Doku-Dateien; DoD 2: Release-Notiz core/release/notes/NJBTVC.md, die volle Race-Suite laeuft vor dem PR."
review-summary: |-
  none
  review: The plan gate is gone. core/gate/gate.go no longer refuses a terminal lane for an open plan step (the block and the CodePlanIncomplete constant are deleted); only the definition-of-done check decides, so with open criteria the refusal lists only those. New gate.OpenPlanHints names each open plan step with its 'jaira dod <id> <n> --done --plan' command; move.Move puts that into Result.Hints whenever the move landed in a terminal lane (forced or not; a refused move gets none). The CLI prints each as a 'Note:' line before the 'X → done' line, so the move line stays last, and as "warnings" in --json; the TUI adds the lines to the message after a move (m), a forced move (f) and the accept key in signoff. Dry-run is unchanged. SKILL.md, README.md, docs/AGENTS.md and builtin 50-done.md no longer describe a plan gate; core/release/notes/NJBTVC.md has two lines. Tests: plan_test.go turned round (signoff->done passes with the plan open; with the DoD open exactly one CodeNeedsSignal), move_test.go checks Hints.
round: 2
review-gaps: |-
  Removed: hintMessage (internal/tui/model.go, a one-line forward to strings.Join) - inlined as strings.Join(res.Hints, "\n") at model.go:2012, model.go:2084, signoff.go:206; stale clause in the comment above core/move/move_test.go:178 ('the plan must not be among the overridden ones' described the Overrode loop critique 1 dropped). Considered and left: gate.OpenPlanHints (core/gate/gate.go:691) vs ticket.PlanComplete (core/ticket/schema.go:579) - PlanComplete returns texts without step numbers, the hint needs the number for 'jaira dod <id> <n> --done --plan', and PlanComplete keeps its caller internal/cli/tickets.go:1615; the DoD refusal loop (core/gate/gate.go:358-385) has the same skip-settled-and-number shape but yields Violations with other wording - a shared helper for two sites would be a new abstraction; folding the hint into settleMessage (internal/tui/model.go:2121) - it reports a settle, and the three sites already compose hint-then-trim the same way; no notify/flash helper (model.go:1856, :1016) covers composing the text. Dead code: CodePlanIncomplete already removed, no other reference in Go, cockpit or docs; remaining 'the plan finished' text is history (core/lane/corrections.go:94, core/move/oldboard_test.go:29) and this repo's own .jaira/lanes/done.md:11, which the release note covers. Cost: OpenPlanHints runs once per terminal move (core/move/move.go:149), nothing to hoist.
  review: Nothing that blocks. Minor: (1) On a board whose terminal lane has logbook-on-entry: true, the same move files the ticket, and the command the hint names then fails - reproduced on a scratch board: "jaira: ticket: not found: XMGWDM", exit 5. The shipped done lane does not file on entry, so only opted-in lanes hit it; release-note line 1 ("tick it afterwards") is wrong there as well. (2) No automated test pins the CLI output (Note: line before the move line - the BDV0HM "last line is the move" rule - and the warnings key) or the three TUI messages; checked by hand only. (3) 'jaira dod --help' (internal/cli/checklist.go:37-38) still says an in-progress item means "a ticket carrying one cannot enter a terminal lane"; read with --plan that is now false. Not one of the four files DoD 5 names. (4) Cosmetic: docs/AGENTS.md:228 not re-wrapped. Outside this diff, reproduced: 'jaira dod <id> --plan --add' on a fresh ticket (Plan holds only the placeholder) writes the steps at the top of the body, not under ## Plan (driver note 19:33). DoD 2: test half open by design (full race suite before the PR); release-note half holds - two lines, each one unwrapped line starting '- ', written as instructions, no version heading.
test-verdict: "pass: gofmt/vet/build clean, core/gate+move+lane and cli Move|Flow race tests green, DoD 1/3/4/5 verified in tree, behaviour exercised on scratch board (CLI text Note: before move line, --json warnings, DoD refusal without plan step, TUI NOTE message); DoD 2 left open for the pre-PR full race suite"
review-verdict: "review: The diff meets DoD 1, 3, 4, 5 and the release-note half of DoD 2; no blocking defect. Checked by hand on a scratch board: Note: line before the move line, warnings in --json, DoD-only refusal, TUI move and signoff accept both show the hint; the installed 0.5.4 refuses the same ticket. The gate tests fail with the plan gate put back. Recommend signoff; the gaps are optional follow-ups."
review-round: 2
review-check: |-
  1. cd /home/berk/git/jAIra-NJBTVC && go build -o /tmp/njbtvc/jaira ./cmd/jaira
  2. mkdir -p /tmp/njbtvc/repo && cd /tmp/njbtvc/repo && git init -q && git commit -q --allow-empty -m init
  3. export JAIRA_NO_FETCH=1 JAIRA_HOME=/tmp/njbtvc/jhome; J=/tmp/njbtvc/jaira; $J init
  4. $J create "plan probe" --goal g --context c --dod "the criterion" - note the 6-letter handle it prints; it is ID below
  5. sed -i 's/^<Steps, in order.*/- [x] write the spec\n- [ ] implement/' .jaira/tickets/*.md - two plan steps under ## Plan, step 2 open (by hand: 'jaira dod --plan --add' puts them in the wrong place on a fresh ticket, a separate bug)
  6. $J set ID assignee="$(git config user.name)" outcome-what=w outcome-why=y outcome-resolves=r
  7. $J dod ID 1 --done --proof x
  8. git commit -q --allow-empty -m "feat(ID): probe"
  9. $J move ID --to done - you see: Note: step 2 ("implement") of the plan is still open; mark it with 'jaira dod ID 2 --done --plan' - and as the LAST line: ID → done. (~/.local/bin/jaira refuses this same move with "the plan is not finished".)
  10. $J move ID --to backlog; $J move ID --to done --json - the JSON has "warnings": ["step 2 (\"implement\") of the plan is still open; ..."]
  11. $J move ID --to backlog; $J dod ID 1 --todo; $J move ID --to done; echo $? - refused, exit 3, exactly one bullet "the definition of done is not met: criterion 1 ...", nothing about the plan
  12. $J dod ID 1 --done --proof x; $J move ID --to signoff; run $J, press Enter on the project list, press the Right arrow 7 times to the "Human Review" column, Enter, then a - a NOTE box says "Accepted ID into Done." and below it the step 2 line from 9
---

# Ein offener Plan-Punkt sperrt die Ziellane nicht mehr, wenn die Definition of Done erfuellt ist


## Definition of Done

- [x] Ein Move in eine Ziellane mit vollstaendiger DoD und offenem Plan-Punkt wird nicht mehr abgelehnt; die Ausgabe nennt den offenen Plan-Punkt als Hinweis; ein Test deckt den Fall
  proof: core/gate/plan_test.go TestUnfinishedPlanDoesNotBlockTerminalLane (CheckAdvance Err()==nil); core/move/move_test.go Hints check
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/
  proof: release note core/release/notes/NJBTVC.md exists; full race suite still to run
- [x] Bei unvollstaendiger DoD nennt die Ablehnung nur die offenen DoD-Punkte, keinen Plan-Schritt mehr
  proof: core/gate/plan_test.go TestUnfinishedPlanAddsNothingToDoDRefusal
- [x] Der Hinweis erscheint in der Textausgabe, im JSON und beim Move aus der TUI, das Annehmen in signoff eingeschlossen
  proof: internal/cli/flow.go Note: lines before the move line + warnings key in --json; internal/tui/model.go applyMove/forceMove; internal/tui/signoff.go accept (compiled; output not run by hand)
- [x] Die Texte, die das Plan-Gate beschreiben (.claude/skills/jaira/SKILL.md, README.md, docs/AGENTS.md, core/lane/builtin/50-done.md), sagen nichts Falsches mehr
  proof: SKILL.md, README.md, docs/AGENTS.md, core/lane/builtin/50-done.md

## Options

- [x] brainstorm
- [ ] planning

## Plan

- [x] Gate: Plan-Sperre in gate.go entfernen, Hinweis-Helfer OpenPlanHints
- [x] move.Result.Hints fuellen bei Move in Ziellane
- [x] CLI Text+JSON zeigen den Hinweis (dry-run unveraendert)
- [x] TUI: model.go zwei Pfade + signoff Accept zeigen den Hinweis
- [x] Tests: plan_test.go drehen, move-Test fuer Hinweis
- [x] Docs (SKILL, README, AGENTS, 50-done) + Release-Notiz
- [x] gofmt/vet/build + Paket-Tests, DoD abhaken

## Progress
- **2026-10-26 19:17 · BeMuCa** — BRAINSTORM 1/3 - Befund im Code.
- core/gate/gate.go:417-431: jede Ziellane (target.Terminal) verweigert mit CodePlanIncomplete, eine Meldung pro offenem Plan-Schritt. Die DoD spielt dabei keine Rolle.
- Die DoD-Sperre steht davor, gate.go:358-385, haengt aber an target.RequiresNonModelSignal, nicht an Terminal. Das ist per Default dasselbe (core/lane/lane.go:344), eine Lane kann es abwaehlen.
- Folge: DoD 9/9 + ein offener Plan-Schritt = Ablehnung. So bei GH0GKD (requirementsgenie): offen war Schritt 30 "R6-13 Abnahme-Vorbereitung", die Arbeit war laut Notiz 23.10. 07:56 erledigt, nur nicht abgehakt.
- Der TUI-Accept in signoff (internal/tui/signoff.go:177-183) hat kein Override: die Person kommt nur weiter, indem sie den Plan-Schritt abhakt.
- Der Kommentar gate.go:404-416 begruendet die Sperre mit "erst Design fertig, nicht gebaut". Aber Plan-Haken und DoD-Haken setzt derselbe Worker selbst. Der Plan-Haken ist kein unabhaengiger Beleg, er faengt nur das Vergessen eines zweiten Hakens.
- Die Doku ist gespalten. Agent-Block core/board/announce.go:98-104 (-> CLAUDE.md/AGENTS.md): die DoD ist das, was das Gate der Ziellane liest; der Plan ist eine zweite, getrennte Liste. Dagegen sagen die Plan-Sperre ausdruecklich an: .claude/skills/jaira/SKILL.md:105-107, README.md:247-249, docs/AGENTS.md:227-228, core/lane/builtin/50-done.md:11.
- Einen Kanal fuer "Move gelungen, aber mit Hinweis" gibt es heute nicht: gate.Violation ist nur Ablehnung (gate.go:41-51), move.Result hat Refused/Overrode/Trimmed/SettleErr (core/move/move.go:58-70).
- **2026-10-26 19:18 · BeMuCa** — BRAINSTORM 2/3 - Optionen und Entscheidung.
(a) Plan-Sperre weg, offene Plan-Schritte als Hinweiszeile beim gelungenen Move. Kosten: ein neuer Hinweis-Kanal in move.Result, Ausgabe in CLI (Text + JSON) und an drei TUI-Stellen. Gibt auf: die harte Sperre bei vergessenem Haken - die faengt aber keine echte Luecke, die DoD-Sperre und signoff/review bleiben.
(b) Beim Abschluss den letzten Plan-Schritt automatisch abhaken. Kosten: jaira schreibt einen Haken, den niemand gesetzt hat - die Datei behauptet Arbeit ohne Beleg. Trifft nur "der letzte" Schritt; offen kann jeder sein (Schritt 3 von 5). Alle abhaken waere Faelschen. Im Worker-Skill statt im Code hilft es der Person in der TUI nicht. Verworfen.
(c) Plan-Sperre ersatzlos streichen, kein Hinweis. Am billigsten (Block loeschen, Test drehen). Gibt auf: den sichtbaren Hinweis, den die DoD des Tickets ausdruecklich verlangt. Verworfen.
Entscheidung: (a). Die DoD ist das Ergebnis und allein das Gate; der Plan ist die Methode und wird nur gemeldet. Das passt zu dem, was der Agent-Block heute schon sagt.
Plan-Sperre bei UNVOLLSTAENDIGER DoD: entfaellt auch. Begruendung:
- Die DoD-Sperre lehnt dann sowieso ab (CodeNeedsSignal je offenem Kriterium, gate.go:362-375; ohne DoD-Liste gate.go:376-383).
- Behalten wuerde nur zusaetzliche Bullets und einen zweiten Code liefern, und eine Regel "Plan sperrt nur, wenn die DoD auch sperrt" - ein Gate, das von einem anderen Gate abhaengt.
- Einziger Fall, in dem sie allein sperren wuerde: eine Ziellane mit requires-nonmodel-signal: false. Die hat die Beleg-Pflicht bewusst abgewaehlt; dort per Plan weiter zu sperren widerspricht dieser Wahl.
- Vorschlag: der Hinweis erscheint nur bei einem Move, der landet. Eine Ablehnung sagt schon, was zu tun ist.
- **2026-10-26 19:18 · BeMuCa** — BRAINSTORM 3/3 - Beruehrungspunkte fuer pre-process (nur gefunden, nichts geaendert).
Code:
- core/gate/gate.go:404-432: Block und Kommentar weg oder umbauen; der Hinweis darf kein Violation sein, sonst zeigt die TUI ihn als Ablehnung (Vorpruefung model.go:1982, signoff.go:177).
- core/gate/gate.go:35 CodePlanIncomplete: wird nicht mehr ausgegeben; pre-process entscheidet: loeschen (plan_test.go benutzt ihn) oder behalten. Kein anderer Verbraucher: core/mod (Cockpit) verzweigt nicht darauf, board.ts:213-225; refusalCode flow.go:345-352 kennt ihn nicht.
- core/ticket/schema.go:579 PlanComplete() liefert die offenen Texte, aber keine Nummern; der Hinweis braucht die Nummer fuer 'jaira dod <id> <n> --done --plan'.
- core/move/move.go:58-70 Result + move.go:139-146: Hinweis hier eintragen, dann bekommen alle vier Schreibstellen ihn.
- internal/cli/flow.go:246-252 Text: Hinweis VOR der Zeile "X -> done" ausgeben (Regel BDV0HM dort: die letzte Zeile ist der Move).
- internal/cli/flow.go:232-244 JSON: neuer Schluessel. Vgl. offenes Ticket 01M26PW18D... "Bei --json steht eine Warnung im JSON, nicht auf einem zweiten Stream".
- internal/cli/flow.go:320-341 reportDryRun: Hinweis optional auch hier.
- TUI: internal/tui/model.go:2012-2027 (normaler Move, notify(...,false)), model.go:2072-2090 (forceMove), signoff.go:196-214 (Accept-Meldung).
- internal/cli/tickets.go:1615-1619 plan_complete im JSON: bleibt wahr, keine Aenderung.
Tests:
- core/gate/plan_test.go:36-63 TestUnfinishedPlanBlocksTerminalLane: dreht sich um (kein Refuse, Hinweis nennt "implement"); Kommentar 36-40 wird falsch. plan_test.go:67-104 pruefen nur "kein CodePlanIncomplete" - an den Konstanten-Entscheid anpassen.
- core/ticket/superseded_test.go:80 nutzt PlanComplete: unberuehrt.
Text, der falsch wird:
- .claude/skills/jaira/SKILL.md:105-107; README.md:247-249; docs/AGENTS.md:227-228; core/lane/builtin/50-done.md:11 (description).
- Bestehende Boards tragen die alte done-Beschreibung in .jaira/lanes/done.md:11. Ob das eine Correction braucht (core/lane/corrections.go) oder stehen bleiben darf: pre-process entscheidet.
- NICHT anfassen: core/lane/corrections.go:94 doneDoorway (woertliche Kopie einer ausgelieferten Lane, Geschichte), core/move/oldboard_test.go:29 (Fixture), core/release/NOTES.md:117 (geschlossen).
- core/board/announce.go:98-104 bleibt wahr; hoechstens ein Halbsatz "ein offener Plan-Schritt wird nur gemeldet".
- Release-Notiz core/release/notes/NJBTVC.md (geaenderte Move-Ausgabe, Move wird nicht mehr abgelehnt).
- **2026-10-26 19:24 · BeMuCa** — Driver: DoD um 3 Punkte ergaenzt, die der Brainstorm vorgeschlagen hat (nur-DoD-Ablehnung, Hinweis in Text/JSON/TUI, Doku). DoD-Punkt 2 (go test ./... -race) laeuft einmal vor dem PR, nicht in einer Lane.
- **2026-10-26 19:28 · BeMuCa** — in-progress: CodePlanIncomplete entfernt (nichts emittiert es mehr, kein anderer Verbraucher). Hinweis = gate.OpenPlanHints -> move.Result.Hints (nur bei Ziellane, vor settle berechnet) -> CLI 'Note:' vor der Move-Zeile + JSON 'hints'; TUI-Meldung an 3 Stellen. Dry-run bewusst unveraendert (sagt nur 'would be allowed'). Gate-/Move-Test: die nackte Fixture loest andere Gates aus (assignee, commits, Outputs), daher pruefen die Tests nur, dass der Plan nicht darunter ist; der Move-Test nutzt --force. Nicht von Hand gegen CLI/TUI gelaufen. Boards mit alter done-Beschreibung (.jaira/lanes/done.md) nicht angefasst.
- **2026-10-26 19:33 · BeMuCa** — Driver: der in-progress-Worker hat die Plan-Schritte von Hand unter ## Plan verschoben. Grund: 'jaira dod --plan --add' schrieb sie bei einem Plan mit Platzhalter-Text als erste Checkliste oben in den Body, sie zaehlten dann als DoD 1-7. Nicht auf einer Board-Kopie nachgestellt; moeglicher Folge-Bug.
- **2026-10-26 19:40 · BeMuCa** — critique 1: (1) core/gate/plan_test.go:51/:73 - use the passing-into-done setup from core/gate/commits_test.go:108-124 (status signoff, commits set, Actor a, Interactive) and assert Err()==nil resp. exactly one CodeNeedsSignal instead of filtering refusal texts for plan/implement; checked in a scratch copy, passes on 6caff07d and fails on its parent; then drop the Overrode loop core/move/move_test.go:193-197, keep the Hints check, Force stays like :139/:158. (2) internal/cli/flow.go:239 JSON key hints -> warnings, the key lanes.go:130, tickets.go:1282/1369, market.go:53 already use for non-fatal lines on success; text Note: line stays (brainstorm 3/3). (3) core/release/notes/NJBTVC.md: rewrite as an instruction, replacement text in review-summary; add a second line telling old boards to delete ', the plan finished if there is one' from .jaira/lanes/done.md (precedent 0.5.3/8NHD8Q.md). Left alone: dry-run without hint (stand on in-progress note 19:28, dry-run reports neither trimmed nor filed either, flow.go:322-347); no correction for old done.md (corrections.go only removes one frontmatter field); hintMessage wrapper and gate.OpenPlanHints placement; doc texts.
- **2026-10-26 19:50 · BeMuCa** — critique 2: all four critique-1 findings resolved in e555efb2 - (1) plan_test.go:49-53/:71-76 use the commits_test.go passing setup, assert Err()==nil and exactly one CodeNeedsSignal; Overrode loop gone from move_test.go, Hints check kept; (2) flow.go:239 key is warnings, no other emitter of warnings in flow.go, no doc names a hints key; (3)+(4) release note is the two suggested lines verbatim, both one line starting '- '; the deleted phrase matches the old 50-done.md description exactly, .jaira/boards/<name>/lanes/ is real (core/lane/boards.go:123). Touched tests pass (core/gate TestUnfinishedPlan*, core/move TestMoveIntoDoneWithOpenPlanStepCarriesAHint). Nothing new; rest of the diff not re-read (lane rule).
- **2026-10-26 20:27 · BeMuCa** — Driver, Uebergabe an signoff: offen fuer die Person - (1) DoD 2: volle Race-Suite vor dem PR; (2) internal/cli/checklist.go:37-38 'jaira dod --help' sagt noch, ein [~]-Punkt sperre die Ziellane - fuer Plan-Punkte jetzt falsch, nicht geaendert; (3) Folge-Bug ausserhalb dieses Tickets: 'jaira dod --plan --add' schreibt bei Platzhalter-Plan die Schritte oben in den Body (vom Review nachgestellt); (4) Ziellane mit logbook-on-entry: der Hinweis nennt dann ein dod-Kommando, das scheitert.
