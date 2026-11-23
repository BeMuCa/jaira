---
id: 01M846QY4DXDWFZV58WSMK8HD2
title: "Cockpit-Abnahme: Abschnitte optisch trennen, damit der Text schnell lesbar ist"
status: done
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
commits:
  - 11089bd4d9ee3d5e339f6e78f2f6f70dbfb58a39
  - 408407f6cdfada7885b5e34e023213c87d335fb7
  - cb1cf932ea87fb7eaf174e5434ed1ea283a3e37e
  - 18b976ab8f8e7f0ef45db1c5a726f6d2af156949
created-at: 2026-11-23T12:40:23Z
updated-at: 2026-11-23T13:33:31Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-91230
claimed-at: 2026-11-23T12:40:52Z
outcome-what: "Acceptance pane: History, Release note and Review notes open with a coloured title and a dim rule; Goal/Done/Why each have a coloured bar and a blank line; history side info (round, age, commits) dim, the DoD/hunt/review line bright; release note lines one per line"
outcome-why: the pane was one white block; the maintainer skims and could not find Goal or where History ended
outcome-resolves: "MK8HD2: DoD 1-4, 6; DoD 5 is the maintainer's look at acceptance"
review-summary: "Delta review from 11089bd4d to 408407f6c. section() in register.tsx now takes only a title: the coloured title and the dim rule share one truncated row, and History's id | lane | reviewed-in moved onto its own dim line (key history-where) that wraps instead of being cut. Closes the earlier fix-now; the rest of the change (coloured bar per Goal/Done/Why, dim side info in History, Release note lines one per row under a ruled title, theme keys only) was read whole in round 1 and is unchanged. claude plugin test: 480 pass."
review-gaps: |-
  note: limit: the rule is 240 '─' cut at the pane edge; a pane wider than 240 columns shows a shorter rule (noted by the implementer)
  note: limit: DoD 5 is the person's to check in a live cockpit; tests check theme keys and dimColor, not real light/dark rendering
test-verdict: "pass: round 2, changed files register.tsx + board.test.ts: claude plugin test . 480/480 (RC=0), go test ./shipped green after manifest regen; History's id | lane line now its own dim wrapping Text (test 'the head says what the ticket is for' checks it); writes nothing; DoD 5 left for the person"
hunt-verdict: "pass: stop-after-fix, no hunt in round 2"
review-verdict: "accept: fix-now answered - History's id/lane/reviewed-in wrap on their own line; no new defects in the delta; 480 tests pass"
review-round: 2
review-verdict-gaps: |-
  note: limit: the rule is 240 '─' cut at the pane edge; a pane wider than 240 columns shows a shorter rule (noted by the implementer)
  note: limit: DoD 5 is the person's to check in a live cockpit; tests check theme keys and dimColor, not real light/dark rendering
review-check: "1. In this worktree start: claude --plugin-dir core/mod/builtin/jaira-cockpit  2. In the cockpit band press Acceptance and open a ticket waiting in signoff (e.g. MK8HD2)  3. Look at the head: Goal, Done and Why each start with a coloured ┃ bar in its own colour, a blank line between them  4. Below: History, Release note and Review notes - choose for each each open with a coloured title and a dim ─── line to the pane edge  5. Under History the line \"<ID> | <lane>\" and the round/commits lines are dim, the \"DoD … | bughunt … | review: …\" line is bright  6. Narrow the terminal: the ID/lane line wraps, nothing after the History title is lost  7. Switch the terminal to a light theme: everything stays readable  8. Seen by a person (DoD 5): you find Goal, Done, Why and Release note at one glance"
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
- [x] seen by a person: der Maintainer findet in der Ansicht Goal, Done, Why und Release note mit einem Blick
  proof: Alex looked at the Acceptance pane and pressed Accept in the cockpit, 23.11.2026
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
- **2026-11-23 12:56 · Alexander Sacharov** — waiting: review round 2 (fresh agent via run-lane.sh) since 13:56
- **2026-11-23 12:58 · Alexander Sacharov** — review: done at 408407f6c
- **2026-11-23 12:59 · Alexander Sacharov** — report: Acceptance pane in blocks - coloured section title + dim rule, coloured bar per Goal/Done/Why, dim side info; one fix-now (History side info wraps) closed | signoff, context 166k of 1000k, cache warm 60 min | DoD 5 open: the maintainer looks at the pane (light and dark); branch not pushed; ticket file modified in worktree, rides with the next commit or logbook
- **2026-11-23 13:01 · Alexander Sacharov** — hunt report: 0 findings, nothing written to MK8HD2 | read whole: none (no new file); read as a diff since 11089bd4d: core/mod/builtin/jaira-cockpit/hooks/register.tsx (section() without its third argument, all 3 callers checked; new history-where line), core/mod/builtin/jaira-cockpit/tests/board.test.ts, core/shipped/shipped.sha256 (new entry = sha256 of HEAD register.tsx, old hash kept) | nothing left unchecked; ran: claude plugin test . (480 pass, 0 fail)
- **2026-11-23 13:01 · Alexander Sacharov** — bughunt: done at 408407f6c - 0 findings (0 change, 0 pre-existing) in 3 modules, follow-up hunt of 11089bd4d..408407f6c; concurrency, error paths, input, permissions checked; files: .jaira/lanes/critique.md .jaira/lanes/in-progress.md .jaira/lanes/review.md .jaira/lanes/testing.md AGENTS.md CLAUDE.md core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/release/notes/MK8HD2.md core/shipped/shipped.sha256
- **2026-11-23 13:29 · Alexander Sacharov** — Accepted from signoff: angenommen:
- note: limit: the rule is 240 '─' cut at the pane edge; a pane wider than 240 columns shows a shorter rule (noted by the implementer)
- note: limit: DoD 5 is the person's to check in a live cockpit; tests check theme keys and dimColor, not real light/dark rendering
- **2026-11-23 13:29 · Alexander Sacharov** — cockpit: accepted by the person, signoff -> done
