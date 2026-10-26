---
id: 01M5WMAAP0J9BN35Y49B503ADG
title: Review trennt Blocker von Hinweisen und legt keine Folge-Tickets an
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Review schickt nur zurueck, was die DoD bricht oder einen Defekt einfuehrt; alles andere ist eine Zeile in review-gaps, und kein Agent macht daraus ein neues Ticket - das entscheidet der Mensch bei der Abnahme."
context: |-
  Was falsch ist: jeder Review findet etwas, und aus fast jedem Fund wird ein neues Ticket. Die Boards wachsen schneller, als sie abgearbeitet werden.
  Zahlen (26.10.2026, ein anderes Board des Maintainers mit derselben review-Lane): 329 Tickets, 137 davon (42 %) mit follows auf ein anderes; 197 in 14 Tagen angelegt, 42 davon nennen den Review im Kontext.
  Ursache 1: core/lane/builtin/40-review.md:65 und .jaira/lanes/review.md:65 sagen 'Default to finding problems. A review that approves everything is worth nothing.' Das Modell liest 'none' als Versagen.
  Ursache 2: review-gaps ist ein Feld ohne Schwere. Ein DoD-Bruch und ein 'waere schoen' stehen gleich da; der Agent danach kann sie nicht unterscheiden und legt fuer jeden ein Ticket an.
  Ursache 3: jedes Folge-Ticket laeuft wieder durch review und erzeugt wieder welche.
  Entschieden von Alex (26.10.): Varianten 1+2 in einem Ticket. 1 = Prompt: jeder Fund ist 'blocker' (bricht DoD oder fuehrt Defekt ein -> zurueck nach in-progress, kein neues Ticket) oder 'note' (eine Zeile in review-gaps, sonst nichts); der Satz 'Default to finding problems' wird ersetzt durch 'gegen die DoD urteilen; none ist eine normale, erwartete Antwort'. 2 = ein Agent legt aus review-gaps nie ein Ticket an; ob aus einer note eins wird, entscheidet die Person in signoff/auf der Abnahme-Seite.
  Nicht gewaehlt: Tiefenlimit fuer follows in jaira create (Variante 3) - erst wenn 1+2 nicht reichen.
  Wo nachsehen: core/lane/builtin/40-review.md (eingebaut), .jaira/lanes/review.md (dieses Board; X9GJ51-Stil: von Hand geaendert?), core/shipped/shipped.sha256 (neue Version der eingebauten Lane listen), Rollen unter core/role/builtin/ die review-gaps lesen oder Tickets anlegen (grep -rn 'review-gaps\|jaira create' core/role/builtin), jaira-role-acceptance (dort muss die note als Kandidat fuer ein Ticket erscheinen).
definition-of-done: "core/lane/builtin/40-review.md teilt jeden Fund in blocker und note und enthaelt 'Default to finding problems' nicht mehr"
tags:
  - gates
blocked-by: []
related: []
commits: []
created-at: 2026-10-26T17:32:18Z
updated-at: 2026-10-26T17:37:23Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-81120
claimed-at: 2026-10-26T17:32:43Z
outcome-what: "review lane prompt (built-in and this board) splits every finding into blocker: and note: lines; blocker sends back to in-progress, note stays a line and the prompt forbids a ticket from it; 'Default to finding problems' replaced by judging against the DoD with none as a normal answer. Acceptance page shows note: lines as candidates and creates a ticket only for one the person names in the comment; worker, dispatcher and teamlead say a note is never a ticket they create. Release note 503ADG."
outcome-why: "every review found something and almost every finding became a follow-up ticket, so boards grew faster than they were worked off"
outcome-resolves: "review no longer reads 'none' as failure, findings carry a severity, and only a person turns a review note into a ticket"
---

# Review trennt Blocker von Hinweisen und legt keine Folge-Tickets an

## Definition of Done

- [x] core/lane/builtin/40-review.md teilt jeden Fund in blocker und note und enthaelt 'Default to finding problems' nicht mehr
  proof: core/lane/builtin/40-review.md:30
- [x] Ein blocker schickt nach in-progress zurueck; eine note bleibt eine Zeile in review-gaps und der Prompt verbietet, daraus ein Ticket anzulegen
  proof: core/lane/builtin/40-review.md:50 (note: never a ticket), :79 (blocker -> in-progress)
- [x] Keine Rolle unter core/role/builtin legt aus review-gaps ein Ticket an; die Abnahme-Seite zeigt notes als Kandidaten, ueber die die Person entscheidet
  proof: core/role/builtin/jaira-role-acceptance/template.html:220 (candidates on the ticket card), SKILL.md:236 (ticket only when named), build.py:155; worker SKILL.md:196, dispatcher SKILL.md:478, teamlead SKILL.md:184
- [x] .jaira/lanes/review.md dieses Boards traegt dieselbe Aenderung; shipped.sha256 kennt die neue Fassung; go test ./... gruen
  proof: diff core/lane/builtin/40-review.md .jaira/lanes/review.md empty; shipped.sha256 does not exist - other boards: jaira lanes use review --force (see note); go test ./... exit 0
- [x] core/release/notes/<ID>.md sagt in einer Zeile, was sich fuer Nutzer der review-Lane aendert
  proof: core/release/notes/503ADG.md:1; go test ./core/release ok

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-26 17:32 · Alexander Sacharov** — parent session: jaira-7a
- **2026-10-26 17:34 · Alexander Sacharov** — decision (worker, before plan): 1) review-gaps: one line per finding, each starting 'blocker:' or 'note:'; 'none' when nothing. A blocker also goes into a jaira note and the ticket moves to in-progress instead of human/signoff. 2) core/shipped/shipped.sha256 does not exist in the repo - DoD 4 names a file nobody has. Lane files are not hashed; core/lane/corrections.go only removes single frontmatter fields by design and must not rewrite a prompt. Other boards get the new prompt with 'jaira lanes use review --force'; the release note says so. 3) Acceptance page: data.json tickets[].candidates = the note: lines of review-gaps, shown on the ticket card under a heading 'candidates for a ticket - say in the comment which one becomes one'. Read-back creates a ticket only for a candidate the person names in the comment; an unnamed one stays a line. No new db collection, no new buttons.
- **2026-10-26 17:37 · Alexander Sacharov** — found out: no role under core/role/builtin created tickets from review-gaps in so many words - the tickets came from agents acting on the field unprompted, so DoD 3 is met by saying it out loud in worker/dispatcher/teamlead/acceptance rather than by removing a line. Acceptance read-back keeps candidates out of decisions.json: they are not cases of a failed step; a ticket made from one is recorded by its follows and a quoted comment.
