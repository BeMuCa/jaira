---
id: 01M8BQAGXPMKWVWG0GANRJW72Z
title: Interreview geht die Tickets einzeln mit der Person durch und beschreibt den Check als Schrittliste
status: review
ready: true
creator: Berk
assignee: BeMuCa
goal: "Die Abnahme im Chat zeigt ein Ticket nach dem anderen (oder so viele auf einmal, wie die Person sagt), und unter Check steht der Pruefweg als nummerierte Schritte, die jemand ohne Vorwissen nachmachen kann."
context: |-
  Was falsch ist: /jaira-role-interreview zeigt alle Tickets in einer Nachricht. Am 24.11. waren das 51 Tickets in signoff, eine Nachricht von rund 27 KB - nicht abzuarbeiten.
  Zweites Problem: die Zeile Check ist ein Einzeiler (Vorgabe in SKILL.md, Abschnitt Present: 'Each ticket is exactly four lines'). Berk verstand daraus nicht, was er tun soll.
  Berks Wunsch (26.11.): 1 Ticket pro Runde, alternativ so viele wie die Person sagt; Check als Bullet-Liste vom Ablauf, den man zum Pruefen macht; darf mehrere Zeilen haben.
  Fable soll die Check-Beschreibung gegenlesen, damit sie verstaendlich ist.
  Datei: der Skill jaira-role-interreview (Branch feat/5QR36R-interreview, installiert unter ~/.claude/skills/jaira-role-interreview/SKILL.md).
definition-of-done: "SKILL.md: ohne Angabe zeigt die Rolle genau ein Ticket, wartet auf die Antwort und zeigt dann das naechste; mit einer Zahl der Person (z.B. '3 auf einmal') so viele"
tags:
  - acceptance
blocked-by: []
related: []
follows: 01M5ZJ555685ZXARWNEF5QR36R
commits: []
created-at: 2026-11-26T10:44:50Z
updated-at: 2026-11-26T10:56:15Z
updated-by: BeMuCa
claimed-by: worker-rjw72z
claimed-at: 2026-11-26T10:46:49Z
outcome-what: "The chat interview shows one ticket per message and waits; a number from the person sets how many come at once, 'all' the rest, 'next' leaves a ticket open. Check is a list: Machine line, Prepare, numbered steps with what is seen afterwards, Wrong if. One example ticket in that format. 'accept all' covers only shown tickets. README line and release note follow."
outcome-why: "51 tickets in one 27 KB message could not be worked through, and a one-line Check was not understood."
outcome-resolves: "Present, Gather and Answers of jaira-role-interreview SKILL.md rewritten; page path untouched"
review-summary: "Only the chat path of core/role/builtin/jaira-role-interreview/SKILL.md changed; the page path (--page) is untouched. Present now shows one ticket per message and waits; a number from the person (\"3 at once\") sets how many come per message, 'all' shows the rest, 'next' leaves the shown ticket open. Each ticket gets a header with '(n of total)', a What line, a separate Machine line (what each already-run check tests, in words, command in brackets, pass/fail), then Check as Prepare (one line per thing, done by the role where a build or temp copy suffices) + numbered steps 'action -> what you see' + 'Wrong if:', then Recommend. Seven writing rules and one full example ticket (A3K9QP) sit below the shape. Gather now orders the tickets once by the page path's block grouping and runs show + machine checks per ticket just before it is shown. 'accept all' now covers only tickets already shown. The record report.md gets its header before the first ticket and one line per ticket as its answer is carried out, plus 'open' and 'not shown' lines. Also: frontmatter description, README paragraph rewrapped, release note core/release/notes/RJW72Z.md with 3 lines, and both new SKILL.md hashes appended to core/shipped/shipped.sha256 (current file hash 43cf1088... is listed; go test ./core/shipped passes)."
round: 2
review-gaps: |-
  note: DoD 4 (Fable cross-read) rests only on the worker's own note 'Gegengelesen (Fable, 26.11.)' - the diff cannot show that a second model actually read it
  note: Record says append a line per ticket 'as soon as its answer is carried out'; explain <id> is also an answer but leaves the ticket open, so a model may write 'explained, open' and later a second line for the same ticket - say that explain writes no line until the ticket is decided or the interview ends
  note: nothing proves a model running the role follows the new shape (one ticket, Prepare, Wrong if) - only a real interview shows it; the worker says the same
test-verdict: "pass: go build and go vet (core/role, core/shipped) RC=0; go test -count=1 ./core/role ./core/shipped ./core/release and ./core/lane -run 'Shipped|Refresh' all ok (RC=0); DoD 1-3 read in the tree at the named sections, DoD 4 note present; branch binary installed the roles into a scratch folder - installed SKILL.md is byte-identical, and 'jaira update' on a scratch board prints the three new note lines. NOT run: the full race suite (house rule: main session before the PR), a live interview with the new text, the bug hunt (no executable code changed; no subagents allowed in this run)."
---

# Interreview geht die Tickets einzeln mit der Person durch und beschreibt den Check als Schrittliste

## Definition of Done

- [x] SKILL.md: ohne Angabe zeigt die Rolle genau ein Ticket, wartet auf die Antwort und zeigt dann das naechste; mit einer Zahl der Person (z.B. '3 auf einmal') so viele
  proof: core/role/builtin/jaira-role-interreview/SKILL.md section Present ('Show ONE ticket per message, then wait') + section Answers (a number, or all; next)
- [x] SKILL.md: Check ist eine nummerierte Liste - Vorbereitung, je Schritt eine Handlung und was man danach sieht; Maschinen-Ergebnisse stehen davor als eigene Zeile
  proof: core/role/builtin/jaira-role-interreview/SKILL.md section Present: ticket shape (Machine: line above Check:, Prepare, numbered steps with →) and the seven rules below it
- [x] Ein Beispiel fuer ein Ticket im neuen Format steht in SKILL.md
  proof: core/role/builtin/jaira-role-interreview/SKILL.md section Present: 'An example - a ticket in signoff' (A3K9QP)
- [x] Gegengelesen von Fable: die Check-Vorgabe ist fuer jemanden ohne Vorwissen verstaendlich (Notiz am Ticket)
  proof: note 'Gegengelesen (Fable, 26.11.)' on this ticket

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-26 10:47 · BeMuCa** — decisions (worker, no open decision needed the person): (1) default batch is 1; a number from the person sets it for the rest of the interview, 'all' shows the rest. (2) 'accept all' now covers only tickets already shown and unanswered - otherwise it would accept tickets nobody saw. (3) new word 'next' leaves the shown tickets open and goes on; without it a person is stuck on a ticket they do not want to decide. (4) Gather runs per ticket just before it is shown, not for all up front: with 51 tickets the first message would otherwise wait for 51 machine runs. (5) Order stays the page path's block grouping, so tickets of one theme come one after another. (6) Prepare is done by the role where a build or a copy suffices (same rule the page path's setup has), else it is commands to paste. Page path (--page) untouched: it shows all tickets by design.
- **2026-11-26 10:50 · BeMuCa** — Gegengelesen (Fable, 26.11.): die Check-Vorgabe in SKILL.md, Abschnitt Present, als jemand ohne Vorwissen gelesen und dreimal nachgeschaerft. (1) Maschinen-Zeile nannte nur das Kommando ('go test ./tui -run ...') - sagt einem Leser nichts; jetzt: in Worten was geprueft wird, Kommando in Klammern. (2) Prepare war ein langer Satz - jetzt eine Zeile pro Sache, und die Rolle bereitet selbst vor, wo ein Build oder eine Kopie reicht. (3) Es fehlte, woran man ein 'Nein' erkennt - neue Zeile 'Wrong if:'. Dazu Regel 'No shorthand': kein interner Name ohne Halbsatz, voller Pfad, woertliche Taste. Nicht geprueft: ob ein Modell, das die Rolle ausfuehrt, die Vorgabe auch so umsetzt - das zeigt erst ein echter Lauf.
- **2026-11-26 10:51 · BeMuCa** — critique: core/role/builtin/jaira-role-interreview/SKILL.md section Answers: 'an answer other than explain' misses a fix that ended in a question of the role's own (asked) - say 'no shown ticket still waits on an explanation or on a question of yours'
core/role/builtin/jaira-role-interreview/SKILL.md section Record: the record is written only at the end; one ticket per message over 51 tickets is a long session, and a killed one leaves no record - append each ticket's line when its answer is carried out
README.md:849 the changed line is 130 characters in a 80-column paragraph - rewrap
Let stand (note 'decisions (worker...)'): next, all, accept-all narrowing, per-ticket gather.
- **2026-11-26 10:52 · BeMuCa** — critique pass 2: the three findings of pass 1 are answered (Answers wording, record appended per ticket, README rewrapped); nothing else read again. review-summary=none.
- **2026-11-26 10:54 · BeMuCa** — report: interreview SKILL.md shows one ticket per message (number/all/next change that), Check is Machine line + Prepare + numbered steps + Wrong if with one example, record appended per ticket; README line + release note | review, context unknown | open: model review not run (worker told to stop here); full race suite and a live interview run are not done; needs a person only to try one real interview
