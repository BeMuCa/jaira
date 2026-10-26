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
updated-at: 2026-10-26T19:55:30Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-56301
claimed-at: 2026-10-26T19:53:57Z
outcome-what: "round 2: review prompt shows how two findings are written as two lines through jaira set (never joined with ';') and says 'Write four things' for its four outputs; board copy synced. Round 1: prompt splits findings into blocker:/note: lines, blocker back to in-progress, note never a ticket; acceptance page shows note: lines as candidates; worker/dispatcher/teamlead never file a note as a ticket; release note 503ADG."
outcome-why: "signoff (Alex) sent back review notes 1 and 3: a joined line shows up as one candidate on the acceptance page, and the count was wrong"
outcome-resolves: "review no longer reads 'none' as failure, findings carry a severity one per line, and only a person turns a review note into a ticket"
review-summary: "The review lane prompt (core/lane/builtin/40-review.md and this board's identical .jaira/lanes/review.md) now asks for one finding per line in review-gaps, each prefixed blocker: (breaks a DoD item or adds a defect) or note: (everything else); 'Default to finding problems' is gone, replaced by judging against the DoD with none as a normal answer; a blocker sends the ticket back to in-progress with a jaira note, otherwise it goes on to the person; a note may never become a ticket by an agent. The acceptance role puts note: lines into tickets[].candidates, the page shows them in an open details box on the ticket card (ru/de/en strings, escaped via fmt), and build.py plan lists them with the ticket comment so a ticket is made only for a candidate the person names. Worker, dispatcher and teamlead skills each gained one sentence saying a note is never a ticket they create. Release note 503ADG.md has two lines (lanes use review --force; roles install)."
review-gaps: |-
  note: the prompt's only multi-finding example is one line; nothing shows how to write two lines through jaira set (e.g. $'note: a\nnote: b'), so an agent may join findings with ';' and the acceptance page then shows them as one candidate
  note: DoD 4 names core/shipped/shipped.sha256, which does not exist in the repo; the worker replaced it with the jaira lanes use review --force line in the release note (see the 17:34 decision note) - the person decides whether that meets the item
  note: 'Write three things, in this order' in 40-review.md lists four items - older than this ticket
test-verdict: "pass: go build RC=0, go test -count=1 ./... RC=0 (31 packages ok); DoD 1-5 checked in the tree (built-in and board review.md identical, old sentence gone, release note passes core/release); fresh jaira init gives the blocker/note prompt, jaira lanes use review --force restores it over an edited copy; acceptance build.py page+plan on a scratch data.json lists candidates and the naming comment, page JS passes node --check"
review-verdict: "The diff meets DoD 1-3 and 5 as written; DoD 4 is met except the shipped.sha256 part, which names a file that does not exist - accepted as the documented decision, not a blocker. go test ./... exits 0, built-in and board review.md are byte-identical, build.py plan prints candidates with and without a comment (checked with a fixture). The HTML candidate box was read, not rendered in a browser."
review-round: 1
review-check: "1. cd /home/alex/projects/.worktrees/jaira-503adg  2. grep -n 'Default to finding' core/lane/builtin/40-review.md .jaira/lanes/review.md - no output  3. diff core/lane/builtin/40-review.md .jaira/lanes/review.md - no output  4. sed -n 30,52p core/lane/builtin/40-review.md - you see the blocker:/note: list and 'Never create a ticket from it'  5. go test ./... - every line ok, exit 0  6. Next acceptance page with a reviewed ticket that has note: lines: its ticket card shows an open box 'Review notes: candidates for a ticket' with those lines"
round: 2
---

# Review trennt Blocker von Hinweisen und legt keine Folge-Tickets an

## Definition of Done

- [x] core/lane/builtin/40-review.md teilt jeden Fund in blocker und note und enthaelt 'Default to finding problems' nicht mehr
  proof: core/lane/builtin/40-review.md:30
- [x] Ein blocker schickt nach in-progress zurueck; eine note bleibt eine Zeile in review-gaps und der Prompt verbietet, daraus ein Ticket anzulegen
  proof: core/lane/builtin/40-review.md:54 (note: never a ticket), :83 (blocker -> in-progress)
- [x] Keine Rolle unter core/role/builtin legt aus review-gaps ein Ticket an; die Abnahme-Seite zeigt notes als Kandidaten, ueber die die Person entscheidet
  proof: core/role/builtin/jaira-role-acceptance/template.html:220 (candidates on the ticket card), SKILL.md:236 (ticket only when named), build.py:155; worker SKILL.md:196, dispatcher SKILL.md:478, teamlead SKILL.md:184
- [x] .jaira/lanes/review.md dieses Boards traegt dieselbe Aenderung; shipped.sha256 kennt die neue Fassung; go test ./... gruen
  proof: diff core/lane/builtin/40-review.md .jaira/lanes/review.md empty; shipped.sha256 does not exist - Alex 2026-10-26: release-note line 'jaira lanes use review --force' meets it; go test ./... exit 0
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
- **2026-10-26 17:37 · Alexander Sacharov** — critique pass 1: none. Let stand: the same no-ticket sentence in worker, dispatcher and teamlead (DoD 3 asks it of every role that reads review output); candidates reuse details.dod rather than a new CSS class; 'Write three things' over four items in 40-review.md is older than this ticket and out of scope.
- **2026-10-26 19:53 · Alexander Sacharov** — decision (Alex, 2026-10-26, signoff, ueber Teamlead): zurueck nach in-progress mit den Review-notes 1 und 3 - (1) im Review-Prompt ein Beispiel, wie zwei Zeilen in review-gaps geschrieben werden (nicht mit ';' verbinden); (3) 'Write three things' -> 'four'. DoD 4 gilt durch die Zeile in der Release-Note als erfuellt (shipped.sha256 existiert nicht), note 2 also erledigt.
- **2026-10-26 19:53 · Alexander Sacharov** — Teamlead jaira-7a: alter Worker-Tab geschlossen (Cache 72 min kalt), neuer Worker fuer die Runde aus signoff.
- **2026-10-26 19:55 · Alexander Sacharov** — round 2 (from signoff): 40-review.md:22 'Write four things'; :43-46 two findings are two lines, never joined with ';', with a $'blocker: ...\nnote: ...' example for jaira set. Board copy synced. Release note unchanged - its first line already covers the review prompt change.
