---
id: 01M4RA8MXYYF2N18GVKNPNC854
title: "M hakt mehrere Milestones an, so wie u mehrere Leute"
status: done
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Mit M waehlt man auf dem Board beliebig viele Milestones aus, und das Board zeigt nur noch deren Tickets - genauso wie u es fuer Leute tut."
context: |-
  Heute waehlt M auf dem Board genau EINEN Milestone und schreibt milestone:<name> in den Filter. Wer zwei Runden Arbeit zusammen sehen will, muss den Filter von Hand tippen.
  Der User-Picker u kann das schon fuer Leute: Leertaste hakt an, Enter schreibt user:a,b in den Filter, x leert (internal/tui/users.go, keyUsers/applyUsers). M soll sich genauso bedienen.
  Der Filter selbst kann es schon: ein Komma heisst bei jedem Schluessel 'eins davon', milestone:a,b funktioniert also heute schon, wenn man es tippt (internal/tui/model.go, Kommentar 'A comma lists alternatives', auf master c4bbdce). Es fehlt nur der Picker.
  Entstanden am 2026-10-12 im Brainstorm zu D0SAHM (Epic-Layer ueber Milestones). Berk will es als eigenes Ticket VOR dem Layer: es ist fuer sich nuetzlich, und der Layer nutzt es spaeter fuer 'Board mit diesem Milestone oeffnen'.
definition-of-done: "Im Milestone-Picker (M) hakt die Leertaste einen Milestone an oder ab; mehrere koennen gleichzeitig angehakt sein, und jeder zeigt seine Farbe"
tags:
  - tui
blocked-by: []
related: []
follows: 01M1KFMETY9MKH4V38TVD0SAHM
commits:
  - aa6b76e2e3e4a61e2d6247433a3d80f9d59410ba
created-at: 2026-10-12T15:03:58Z
updated-at: 2026-10-13T19:18:46Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-1610342
claimed-at: 2026-10-13T11:41:29Z
outcome-what: "M now ticks several milestones: space ticks a row, enter writes milestone:a,b into the board filter, x clears it, esc leaves the filter alone, and each row keeps its colour swatch. u and M run on one ticking picker in internal/tui/ticks.go (openTicks/keyTicks/applyTicks/tickMark), which replaces keyUsers/applyUsers and the old modeMilestones key switch. Help line, key bar ('M milestones'), README '## Keys', 'jaira milestone --help' and core/release/notes/PNC854.md updated."
outcome-why: "M could narrow the board to only one milestone; seeing two rounds of work together meant typing milestone:a,b by hand. The filter already understood commas, so only the picker was missing. One mechanism instead of a copy of users.go, so u and M cannot drift apart."
outcome-resolves: "DoD1: space toggles a tick per row and the swatch stays (TestMilestonePickerTicksSeveral, checked in the binary). DoD2: enter goes through setFilterTerm, which replaces only the milestone: term and keeps the rest (TestMilestonePickerOwnsOnlyItsTerm, TestMilestonePickerKeepsTheUserFilter). DoD3: u and M share keyTicks/openTicks, so reopen, x and esc are the same code (same tests). DoD4: tests in internal/tui; go test ./... -race exit 0. DoD5: help view.go:1671, key bar view.go:1043, README.md:785, core/release/notes/PNC854.md. Behaviour change: enter with nothing ticked now removes milestone: (like u) instead of taking the cursor row."
executed-by: opus
review-summary: "M on the board now works like u. The ticking logic that lived only in internal/tui/users.go (cursor, ticks, what space/enter/x/esc do, how the filter is written) moved into one new file, internal/tui/ticks.go (openTicks, keyTicks, applyTicks, tickMark), and both pickers run on it: u writes user:, M writes milestone:. In the M picker space ticks the row under the cursor; every row shows a tick (filled or empty circle) in front of its colour swatch; enter writes the ticked names as milestone:a,b into the ordinary board filter and leaves every other filter term as it was; x removes milestone: and closes; esc closes without changing anything; reopening M pre-ticks what milestone: already holds. The old single-select key switch and msIdx are gone. Behaviour change, settled by Berk: enter with nothing ticked now removes milestone: instead of filtering to the cursor row. Docs: help line for M, 'M milestones' in the key bar, README '## Keys' line, one sentence in 'jaira milestone --help', core/release/notes/PNC854.md. Tests: two new (TestMilestonePickerTicksSeveral, TestMilestonePickerOwnsOnlyItsTerm); three existing ones now tick with space before enter or read the shared picker state."
review-gaps: "none that block. Checked and fine: all five DoD items are in the diff and work in the binary; u behaves as before (same keys, same filter; its tests only had field names renamed); ticks are kept by name, so when another session creates a milestone while M is open, the right milestone stays ticked - only the cursor row shifts, as it did before this change. Small, not from this change: (1) at 40 columns the picker hint line is cut off ('space tick · enter show only t', 'esc close' missing) - the u picker is cut the same way, because wrapHints wraps at the terminal width, not the popup width; pre-existing. (2) A typed milestone: value that is no milestone on the board (typo, filed milestone) is not shown in the picker and enter drops it - same as u with user:. (3) The worktree has untracked .jaira/milestones/demo-*.md from a ref sync; they are not in ff585cb and must not be swept into a later commit by git add -A."
test-verdict: "pass: go build/vet ./... and go test ./... -race -count=1 RC=0 (no FAIL/race/panic); DoD 1-5 checked in the tree and the named tests run -v (TestMilestonePickerTicksSeveral, TestMilestonePickerOwnsOnlyItsTerm, TestMilestonePickerKeepsTheUserFilter, TestPickerNarrowsTheBoardAndReleasesIt, TestNarrowBoardShowsAllKeysAndFits, user picker tests) all PASS; binary from ff585cb on a scratch board: M space j space enter -> filter milestone:round-one,round-two, Alpha+Beta shown, Gamma hidden; reopen shows both ticks with own colours; untick+esc keeps filter; M x shows all three; u then M gives 'user:tester milestone:round-two'; ? shows the new M line; at 40x24 the key bar wraps and keeps 'M milestones'"
question: "M arbeitet jetzt wie u: enter schreibt nur, was angehakt ist. 'M enter' OHNE Leertaste leert deshalb jetzt milestone: (wie u bei user:), statt wie bisher auf die Zeile unter dem Cursor zu filtern. So lassen (genau wie u, so gebaut und vom Ticket verlangt)? Oder soll enter ohne Haken die Cursor-Zeile nehmen? Kosten der zweiten Variante: wer alles abhakt und enter drueckt, bekommt dann einen Filter auf die Cursor-Zeile statt eines leeren, und M verhaelt sich wieder anders als u."
review-verdict: "pass: the diff meets all five DoD items; no defect found. go build, go vet and go test ./... -race -count=1 green at ff585cb. Binary on a throwaway board: two ticks + enter -> 'filter: milestone:round-one,round-two', 2 tickets, only those two; reopen shows both ticks; untick + esc keeps the filter; u then M x -> 'user:sam' kept; enter with nothing ticked removes milestone: (Berk's settled choice). One mechanism for u and M instead of a copy, so they cannot drift apart."
review-check: |-
  1. Binary bauen (NICHT nach ~/.local/bin): cd /home/berk/git/jAIra-PNC854 && go build -o /tmp/jaira-pnc854 ./cmd/jaira  -> keine Ausgabe = ok.
  2. Testboard ohne Remote: mkdir -p /tmp/pnc854-board && cd /tmp/pnc854-board && git init -q && /tmp/jaira-pnc854 init
  3. Pruefen, dass kein Remote da ist: git remote -v  -> keine Ausgabe.
  4. Vier Tickets: for t in Alpha Beta Gamma Delta; do /tmp/jaira-pnc854 create "$t" --goal g --context c --dod d; done
  5. /tmp/jaira-pnc854 list  -> vier Zeilen; die Kurz-IDs von Alpha, Beta, Gamma merken.
  6. Drei Milestones: /tmp/jaira-pnc854 milestone create round-one <ID-Alpha>  ;  /tmp/jaira-pnc854 milestone create round-two <ID-Beta>  ;  /tmp/jaira-pnc854 milestone create round-three <ID-Gamma>  (Delta bleibt ohne).
  7. Eine Person fuer u: /tmp/jaira-pnc854 set <ID-Alpha> assignee=sam  -> 'Updated ...'.
  8. /tmp/jaira-pnc854 starten. Kommt die Projektliste: Cursor auf pnc854-board, enter. Unten in der Tastenleiste steht 'M milestones'.
  9. M druecken -> Fenster 'Milestones' mit round-one, round-three, round-two. Jede Zeile: leerer Kreis, farbiges Quadrat, Name, 1. Darunter: 'space tick · enter show only the ticked · x show everything again · esc close'.
  10. Leertaste -> round-one hat jetzt einen vollen Kreis, das farbige Quadrat bleibt.
  11. j, j (Cursor auf round-two), Leertaste -> round-two auch voller Kreis; round-three bleibt leer.
  12. enter -> Fenster zu. Oben steht 'filter: milestone:round-one,round-two', rechts '2 tickets'. Nur Alpha und Beta sind zu sehen, Gamma und Delta nicht.
  13. M -> round-one und round-two haben schon den vollen Kreis.
  14. Leertaste (round-one verliert den Haken), dann esc -> Filter oben unveraendert 'milestone:round-one,round-two'.
  15. u -> Personen-Fenster wie bisher: oben du mit 4, darunter sam mit 1. j auf sam, Leertaste, enter -> Filter 'milestone:round-one,round-two user:sam', nur Alpha zu sehen.
  16. M, dann x -> Filter nur noch 'user:sam'. milestone: ist weg, user: ist geblieben.
  17. M, Leertaste, enter -> 'user:sam milestone:round-one'. Dann M, Leertaste (Haken weg), enter -> wieder nur 'user:sam' (enter ohne Haken leert milestone:, wie u - so gewollt).
  18. ? -> Hilfe-Zeile 'M  narrow the board to some milestones: space ticks them, enter applies (x there shows everything again)'.
  19. q beenden. Aufraeumen: rm -rf /tmp/pnc854-board /tmp/jaira-pnc854 (steht das Board danach noch in der Projektliste: dort Cursor drauf, x).
---

# M hakt mehrere Milestones an, so wie u mehrere Leute

## Definition of Done

- [x] Im Milestone-Picker (M) hakt die Leertaste einen Milestone an oder ab; mehrere koennen gleichzeitig angehakt sein, und jeder zeigt seine Farbe
  proof: internal/tui/ticks.go:36 keyTicks space; view.go:1619 tick before swatch; TestMilestonePickerTicksSeveral (two ticks, swatches 38;5;45 and 38;5;202); binary: M space j space shows both ● with colours 100/208
- [x] Enter schreibt die angehakten Milestones als milestone:a,b in den Board-Filter, ersetzt dabei eine vorhandene milestone:-Bedingung und laesst alle anderen Filter-Begriffe stehen; das Board zeigt danach nur Tickets dieser Milestones
  proof: internal/tui/ticks.go:70 applyTicks -> setFilterTerm(milestone); TestMilestonePickerTicksSeveral (filter milestone:round-one,round-two, board shows only those two); TestMilestonePickerOwnsOnlyItsTerm (user:berk kept, old milestone: replaced); TestMilestonePickerKeepsTheUserFilter
- [x] Wiederoeffnen, x und esc verhalten sich wie beim User-Picker u (internal/tui/users.go)
  proof: u and M share one handler: model.go:1243 keyTicks(s, "M", ...) and ticks.go:22 openTicks; TestMilestonePickerOwnsOnlyItsTerm (reopen shows ticks, esc changes nothing, x keeps user:berk); TestPickerNarrowsTheBoardAndReleasesIt; userpicker_test.go green
- [x] Ein Test in internal/tui deckt mehrere Haken, Enter und den entstehenden Filter ab; go test ./... -race ist gruen
  proof: internal/tui/milestonebar_test.go:205 TestMilestonePickerTicksSeveral, :256 TestMilestonePickerOwnsOnlyItsTerm; go test ./... -race exit 0 (28 ok, 0 FAIL)
- [x] Hilfe (?), Tastenleiste, README '## Keys' und core/release/notes/PNC854.md nennen die Mehrfachauswahl
  proof: help view.go:1671; key bar view.go:1043 'M milestones' (TestNarrowBoardShowsAllKeysAndFits); README.md:785; core/release/notes/PNC854.md; also jaira milestone --help text internal/cli/milestones.go:59

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Failing tests in internal/tui/milestonebar_test.go: two ticks + enter give milestone:round-one,round-two and the board shows only their tickets; reopen shows the ticks; enter replaces an existing milestone: term and keeps the rest; x clears; esc changes nothing; each row shows its colour swatch and its tick
- [x] Change the two tests that pin single-select 'M enter' (TestPickerNarrowsTheBoardAndReleasesIt, TestMilestonePickerKeepsTheUserFilter) to tick with space before enter
- [x] Pull the ticking mechanism out of internal/tui/users.go into one shared picker (cursor, ticks, filter key, one key handler, one apply) and put u on it; userpicker_test.go stays green
- [x] Put M on the shared picker: open seeds ticks from the milestone: term, keys go through the shared handler, renderMilestones shows tick + swatch; remove the old modeMilestones switch and msIdx
- [x] Help (?) line for M, 'M milestones' in the key bar (and its hints_test list), README '## Keys' line, core/release/notes/PNC854.md
- [x] go vet + go test ./... -race green; build to scratchpad and drive M, space, j, space, enter, M, x, esc in the binary under tmux

## Progress
- **2026-10-12 15:05 · BeMuCa** — DoD 5 umformuliert: das Ticket nannte core/release/NOTES.md unter '## Unreleased', aber CLAUDE.md auf master c4bbdce sagt NOTES.md ist geschlossene Historie (bis 0.3.6); neue Notizen gehoeren nach core/release/notes/<TICKET>.md. Planning-Option angehakt, weil der Auftrag einen Plan im pre-process verlangt.
- **2026-10-12 15:08 · BeMuCa** — Plan-Begruendung: (A) users.go fuer Milestones spiegeln = ~40 Zeilen doppelte Tastenlogik neben keyUsers; (B) den Ankreuz-Mechanismus aus users.go herausziehen (Cursor, Haken, Filter-Schluessel, ein Tasten-Handler) und u und M beide darauf setzen. Gewaehlt B: der Auftrag verlangt ausdruecklich keinen zweiten Mechanismus, und optimize wuerde die Kopie streichen. Kosten: userIdx/userPicked/msIdx werden ein gemeinsamer Zustand, zwei Zeilen in userpicker_test.go aendern sich. Enter ohne Haken entfernt milestone: aus dem Filter, genau wie u bei user: - damit aendern sich die zwei Tests, die 'M enter' ohne Leertaste als Ein-Milestone-Geste festnageln (TestPickerNarrowsTheBoardAndReleasesIt, TestMilestonePickerKeepsTheUserFilter). Verworfen: 'enter ohne Haken nimmt die Cursor-Zeile' - dann wuerde Abhaken-und-enter ploetzlich auf die Cursor-Zeile filtern statt zu leeren. Milestone-Namen sind kebab-case (tag.Normalize), brauchen also nie Anfuehrungszeichen im Filter. README '## Keys' und die Tastenleiste nennen M heute gar nicht; beide bekommen es neu.
- **2026-10-12 15:12 · BeMuCa** — Umsetzung: internal/tui/ticks.go traegt jetzt den einen Ankreuz-Mechanismus (tickPicker: Filter-Schluessel, Cursor, Haken; openTicks/keyTicks/applyTicks/tickMark). Die Zeilenliste wird bei JEDER Taste frisch uebergeben (userNames/milestoneNames), nicht beim Oeffnen gespeichert: M zeichnet live aus m.milestones, und ein reload (Watcher, andere Session legt Milestone an) bei offenem Picker haette sonst Cursor-Zeile und angehakte Zeile auseinanderlaufen lassen. u behaelt seinen Schnappschuss m.users vom Oeffnen, wie vorher.
- **2026-10-12 15:21 · BeMuCa** — critique (pass 1, whole diff ff585cb): keine Befunde. Geprueft und stehen gelassen: (1) Namen pro Taste statt beim Oeffnen gespeichert - steht auf der Umsetzungs-Notiz (reload bei offenem M-Picker); (2) kein anderer Ankreuz-Picker in internal/tui (links/home/legend sind Ein-Cursor), ticks.go ist also kein zweites Muster; (3) nichts Spekulatives: own und tickMark haben je zwei Aufrufer, das Quoting in applyTicks ist der Bedarf von u; (4) milestoneNames neben milestoneColors, userNames in users.go. Verhaltensaenderung 'M enter ohne Haken leert milestone:' ist vom Ticket so verlangt (genauso wie u) - geht als Abnahme-Frage an Berk, nicht als Befund.
- **2026-10-12 15:21 · BeMuCa** — optimize: nichts entfernt. staticcheck -checks U1000 ./internal/tui meldet laneStart (model.go:267) und currentLane (model.go:640) als ungenutzt - beides schon auf master c4bbdce, nicht von dieser Aenderung, stehen gelassen. .jaira/milestones/demo-*.md im Worktree ist untracked und kam von jairas Ref-Sync waehrend eines Board-Kommandos, nicht von dieser Aenderung; nicht committet.
- **2026-10-12 15:26 · BeMuCa** — testing: bestanden, nichts gefunden. Geprueft im Binary (scratchpad/jaira-pnc854, Testboard mit HOME-Override, damit Berks Boards unberuehrt bleiben). Offen fuer Berk nur die Abnahme-Frage zum Verhalten von 'M enter' ohne Haken.
- **2026-10-12 16:03 · BeMuCa** — 2026-10-12 Berk zur Frage in human: ja - enter ohne Haken leert milestone: wie bei u, so lassen. Die Korrektur des Hilfetexts von 'jaira milestone --help' ist ok. Berks erster Handtest schlug fehl; im echten Terminal (tmux, Binary aus ff585cb, Board jAIra mit 3 demo-Milestones) laufen M/space/j/space/enter -> filter milestone:demo-board-dateien,demo-naechste-version (4 Tickets), erneutes M hakt beide vor, x -> 71 Tickets. Vermutete Ursache des Fehlschlags: 'jaira' im PATH ist das alte Binary vom 05.10. (~/.local/bin), und das Test-Binary startet in der Projektliste mit dem Cursor auf requirementsgenie.
- **2026-10-13 11:40 · BeMuCa** — 2026-10-13 Berk beantwortet die human-Frage endgueltig: 'M soll wie u arbeiten und so bleiben' - enter ohne Haken leert milestone:, genau wie u bei user:. Keine Aenderung noetig; weiter nach review.
- **2026-10-13 19:18 · BeMuCa** — Abnahme 2026-10-13: Berk hat im Chat gesagt, er habe PNC854, D0SAHM, 9WE1QK und J1ZVT9 durchgetestet, und das Ablegen ins Logbuch verlangt (Alex' Review auf PR #74: 'schliss mal tickets ab'). Claude hat den Move nach done ausgefuehrt - mit --force, weil eine Agent-Session signoff sonst nicht verlassen darf.
