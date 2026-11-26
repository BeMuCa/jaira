---
id: 01M8BQAGXPMKWVWG0GANRJW72Z
title: Interreview geht die Tickets einzeln mit der Person durch und beschreibt den Check als Schrittliste
status: critique
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
updated-at: 2026-11-26T10:52:04Z
updated-by: BeMuCa
claimed-by: worker-rjw72z
claimed-at: 2026-11-26T10:46:49Z
outcome-what: "The chat interview shows one ticket per message and waits; a number from the person sets how many come at once, 'all' the rest, 'next' leaves a ticket open. Check is a list: Machine line, Prepare, numbered steps with what is seen afterwards, Wrong if. One example ticket in that format. 'accept all' covers only shown tickets. README line and release note follow."
outcome-why: "51 tickets in one 27 KB message could not be worked through, and a one-line Check was not understood."
outcome-resolves: "Present, Gather and Answers of jaira-role-interreview SKILL.md rewritten; page path untouched"
review-summary: none
round: 2
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
