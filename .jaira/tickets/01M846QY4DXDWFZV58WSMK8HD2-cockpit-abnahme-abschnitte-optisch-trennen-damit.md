---
id: 01M846QY4DXDWFZV58WSMK8HD2
title: "Cockpit-Abnahme: Abschnitte optisch trennen, damit der Text schnell lesbar ist"
status: critique
ready: true
creator: Alexander Sacharov
goal: "Die Abnahme-Ansicht (Приём) liest sich auf einen Blick: Kopf (Titel, Goal, Done, Why), History, Release note und Review-Notizen sind sichtbar getrennte Blöcke, und innerhalb des Kopfs springt jedes Feld ins Auge."
context: |-
  Was heute falsch ist: in der Abnahme-Ansicht des jaira-cockpit steht fast alles als weißer Fließtext untereinander. Goal, Done und Why im Kopf haben nur ein fettes Wort davor und verschwimmen zu einem Block; History und Release note gehen nahtlos ineinander über.
  Unten bei den Review-Notizen klappt es schon: jede Notiz hat ihre Farbmarke (🔴🟡⚪🔵) und eine Knopfzeile darunter - das liest der Maintainer gut.
  Wunsch des Maintainers (23.11.2026, mit Screenshot): Text so teilen, dass er schnell gelesen wird; 'Streifen' zwischen den Teilen; weniger Weiß, das sich vermischt; oben soll klar sein, wo Goal ist und wo was.
  Wo: core/mod/builtin/jaira-cockpit/hooks/register.tsx (Abnahme-Ansicht, Story/readStory, Kopf und History), Texte in hooks/texts.ts.
  Leser: hat ADHS, liest langsam - Ziel ist Überfliegen, nicht schöner Satz.
  Offen für den Worker: welche Mittel die Engine bietet (Trennlinien, Hintergrund, Farbe, Einrückung) - das stärkste davon nehmen, das im Terminal sicher rendert.
definition-of-done: "Zwischen Kopf, History, Release note und Review-Notizen steht je eine sichtbare Trennung (Linie oder farbiger Abschnittstitel)"
tags:
  - cockpit
  - acceptance
blocked-by: []
related: []
commits: []
created-at: 2026-11-23T12:40:23Z
updated-at: 2026-11-23T12:56:07Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-91230
claimed-at: 2026-11-23T12:40:52Z
outcome-what: "Acceptance pane: History, Release note and Review notes open with a coloured title and a dim rule; Goal/Done/Why each have a coloured bar and a blank line; history side info (round, age, commits) dim, the DoD/hunt/review line bright; release note lines one per line"
outcome-why: the pane was one white block; the maintainer skims and could not find Goal or where History ended
outcome-resolves: "MK8HD2: DoD 1-4, 6; DoD 5 is the maintainer's look at acceptance"
review-summary: "The Acceptance pane in the cockpit (register.tsx) now draws History, Release note and Review notes as blocks: each opens with a bold title in the theme colour 'claude' and a dim ─ rule cut at the pane edge (section()). Goal, Done and Why each start with a '┃' bar in their own theme colour (suggestion/success/warning) and a blank line above (field()). History lines for round/age and commits/PR are dim, the DoD/bughunt/review line stays bright. Release note lines now stand one per line under their heading instead of joined with ' / '. Only theme keys and dimColor, so light and dark terminals both read it. Tests updated (480 pass), release note and shipped hash added."
review-gaps: |-
  note: edge: fix-now: register.tsx section() puts the History heading's side info (id | lane | reviewed in round N) into the same truncate-end Text as the rule, so on a narrow pane that info is cut off - before this change it wrapped; fix: keep the rule in its own truncated Text after a wrapping title+info Text
  note: limit: RULE is 240 chars, a pane wider than that shows a shorter rule - harmless
test-verdict: "pass: claude plugin test . 480/480 (RC=0), claude plugin validate passed (no warning on changed lines), go test ./shipped ./mod ./release green after manifest regen; DoD 1-4,6 checked in register.tsx section()/field()/history + 'the head says what the ticket is for' (colours by key, dimColor per history line, rule after each title); failure path: no release note file -> no Release note block (test 'keeps no release notes'), no review notes -> no notes header; change writes nothing, so no two-writer case; DoD 5 left for the person"
hunt-verdict: "pass: bughunt done at 11089bd4d, 1 change finding (edge, not a blocker: History heading's side info cut at the pane edge on narrow panes; the light line still names the stale review), 0 pre-existing"
review-verdict: "accept: DoD 1-4 and 6 are met by the diff and 480 cockpit tests pass; one fix-now edge (History side info truncated); DoD 5 is the person's"
review-round: 1
review-verdict-gaps: |-
  note: edge: fix-now: register.tsx section() puts the History heading's side info (id | lane | reviewed in round N) into the same truncate-end Text as the rule, so on a narrow pane that info is cut off - before this change it wrapped; fix: keep the rule in its own truncated Text after a wrapping title+info Text
  note: limit: RULE is 240 chars, a pane wider than that shows a shorter rule - harmless
review-check: "1. In a repo with a ticket in signoff, start claude with the branch cockpit: claude --plugin-dir core/mod/builtin/jaira-cockpit  2. In the cockpit band press [Acceptance] on that ticket  3. See: Goal, Done, Why each start with a coloured ┃ bar (three different colours) and a blank line between them  4. See: 'History', 'Release note', 'Review notes - choose for each' each in colour, followed by a dim ──── line to the pane edge  5. See: round/age and commits lines are dim, the 'DoD x/y | bughunt | review' line is bright  6. Switch the terminal to a light theme and repeat 3-5: all still readable  7. DoD 5 (seen by a person): you find Goal, Done, Why and Release note at a glance - tick it yourself"
round: 2
back-from: review
back-by: lane
back-count: 1
stop-after-fix: "true"
---

# Cockpit-Abnahme: Abschnitte optisch trennen, damit der Text schnell lesbar ist

## Definition of Done

- [x] Zwischen Kopf, History, Release note und Review-Notizen steht je eine sichtbare Trennung (Linie oder farbiger Abschnittstitel)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx section() + tests/board.test.ts 'the head says what the ticket is for'
- [x] Im Kopf hat jedes Feld (Goal, Done, Why) eine eigene Farbe oder Marke am Zeilenanfang und eine Leerzeile Abstand, sodass man ein Feld findet, ohne zu lesen
  proof: register.tsx field(): BAR + colour per field, marginTop 1; test checks suggestion/success/warning
- [x] Langer Fließtext ist nicht mehr gleichfarbig weiß: Beschriftungen und Nebeninfos (Runde, Zeit, Commits) sind gedimmt, der Inhalt hell
  proof: register.tsx history {dim}; test: round/commits dimColor, DoD line not
- [x] Hell- und Dunkelmodus des Terminals bleiben lesbar
  proof: only theme keys (claude/suggestion/success/warning) and dimColor, no hex colours: register.tsx section()/field()
- [ ] seen by a person: der Maintainer findet in der Ansicht Goal, Done, Why und Release note mit einem Blick
- [x] Release-Note in core/release/notes/<ID>.md
  proof: core/release/notes/MK8HD2.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] section header helper: coloured title + dim rule
- [x] head: mark + colour per field, blank line between
- [x] history: side info dimmed, content bright
- [x] tests + release note

## Progress
- **2026-11-23 12:40 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6958.sock
- **2026-11-23 12:42 · Alexander Sacharov** — decision: engine has no per-side Box borders (BoxProps: borderStyle only, all four sides). Separator = section title in colour + dim '─' run cut at pane width (Text wrap=truncate-end). Head fields: coloured '▌' mark + coloured bold label per field (Goal suggestion, Done success, Why warning), marginTop 1. Theme keys only, so light/dark follow the theme. No decision left for a person.
- **2026-11-23 12:47 · Alexander Sacharov** — the release note lines are now one per line under their title, no longer joined with ' / '. Section titles lost their trailing colon (texts.ts). shipped.sha256 regenerated (gate shipped-manifest).
- **2026-11-23 12:47 · Alexander Sacharov** — critique: none. Considered e.props.bodyColumns (used by the Overview pane, register.tsx:2539) for an exact rule width: left it - title width with wide glyphs would need measuring, truncate-end cuts the 240-char rule at the edge with nothing to compute; a pane wider than 240 columns gets a shorter rule, harmless.
- **2026-11-23 12:50 · Alexander Sacharov** — waiting: bughunt (background agent) since 13:50
- **2026-11-23 12:50 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:2337-2341,2359 - class: side information that used to wrap is now cut at the pane edge, because section() renders the title plus its 'after' text in one Text wrap="truncate-end" (needed for RULE), so whatever follows the heading on that row is lost on a narrow pane; example: 'History | MK8HD2 | signoff | review of round 2, changed since' is 59 columns, so on a narrower Acceptance pane the stale-review warning ('changed since') is cut. Before the change that row used wrap="wrap". Not a blocker: the light line at the top still names oldReview as a reason, although after the first reason it can only appear as '+N more'. Fix idea: keep the dim rule in its own truncated Text and let 'after' wrap. Confirmed by reading; the pane width depends on the terminal and was not measured
- **2026-11-23 12:50 · Alexander Sacharov** — hunt report: 1 finding, 1 change (edge, not a blocker), noted on MK8HD2 | read whole: the Acceptance pane render and the story loader openAcceptance/readStory in register.tsx, plus releaseLines in board.ts; read as a diff: texts.ts, board.test.ts, release note, shipped.sha256 (hashes match HEAD files); the uncommitted .jaira/lanes, AGENTS.md and CLAUDE.md edits are jaira update 0.6.4 output (dry run: 0 updated), not this change; ran claude plugin test (480 pass, 0 fail) | not checked: real pane widths in a live cockpit, so how often the History row gets cut is unmeasured; register.tsx outside the acceptance pane (2600 lines) was not read whole, since the change touches only rendering
- **2026-11-23 12:50 · Alexander Sacharov** — bughunt: done at 11089bd4d - 1 finding (1 change, 0 pre-existing) in 3 modules, first hunt; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/MK8HD2.md core/shipped/shipped.sha256
- **2026-11-23 12:51 · Alexander Sacharov** — testing: note: edge: section() truncates its 'after' text with the rule; example: 'review of round 2, changed since' lost on a narrow pane; fix: rule in its own truncated Text, 'after' wraps - left for review/person, not a blocker under the default rule
- **2026-11-23 12:51 · Alexander Sacharov** — waiting: review (fresh agent via run-lane.sh) since 13:51
- **2026-11-23 12:53 · Alexander Sacharov** — review: done at 11089bd4d
- **2026-11-23 12:53 · Alexander Sacharov** — review: fix-now: section() in register.tsx truncates the History heading's side info (id | lane | reviewed in) together with the rule; put the rule in its own truncate-end Text so the title+info wraps as before
- **2026-11-23 12:54 · Alexander Sacharov** — set stop-after-fix=true: fix-now is a few render lines, no lock, no parser
- **2026-11-23 12:56 · Alexander Sacharov** — fix-now closed: History's id | lane | reviewed-in moved off the rule row onto its own dim wrapping line (register.tsx section() takes only a title now); test checks the line. 480 pass.
