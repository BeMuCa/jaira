---
id: 01M4RA8MXYYF2N18GVKNPNC854
title: "M hakt mehrere Milestones an, so wie u mehrere Leute"
status: critique
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
commits: []
created-at: 2026-10-12T15:03:58Z
updated-at: 2026-10-12T15:20:11Z
updated-by: BeMuCa
claimed-by: EE-3NX6GL3-2718199
claimed-at: 2026-10-12T15:05:10Z
outcome-what: "M now ticks several milestones: space ticks a row, enter writes milestone:a,b into the board filter, x clears it, esc leaves the filter alone, and each row keeps its colour swatch. u and M run on one ticking picker in internal/tui/ticks.go (openTicks/keyTicks/applyTicks/tickMark), which replaces keyUsers/applyUsers and the old modeMilestones key switch. Help line, key bar ('M milestones'), README '## Keys', 'jaira milestone --help' and core/release/notes/PNC854.md updated."
outcome-why: "M could narrow the board to only one milestone; seeing two rounds of work together meant typing milestone:a,b by hand. The filter already understood commas, so only the picker was missing. One mechanism instead of a copy of users.go, so u and M cannot drift apart."
outcome-resolves: "DoD1: space toggles a tick per row and the swatch stays (TestMilestonePickerTicksSeveral, checked in the binary). DoD2: enter goes through setFilterTerm, which replaces only the milestone: term and keeps the rest (TestMilestonePickerOwnsOnlyItsTerm, TestMilestonePickerKeepsTheUserFilter). DoD3: u and M share keyTicks/openTicks, so reopen, x and esc are the same code (same tests). DoD4: tests in internal/tui; go test ./... -race exit 0. DoD5: help view.go:1671, key bar view.go:1043, README.md:785, core/release/notes/PNC854.md. Behaviour change: enter with nothing ticked now removes milestone: (like u) instead of taking the cursor row."
executed-by: opus
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
