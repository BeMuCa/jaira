---
id: 01M8BQAGXPMKWVWG0GANRJW72Z
title: Interreview geht die Tickets einzeln mit der Person durch und beschreibt den Check als Schrittliste
status: backlog
ready: true
creator: Berk
assignee: Berk
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
updated-at: 2026-11-26T10:44:50Z
---

# Interreview geht die Tickets einzeln mit der Person durch und beschreibt den Check als Schrittliste

## Definition of Done

- [ ] SKILL.md: ohne Angabe zeigt die Rolle genau ein Ticket, wartet auf die Antwort und zeigt dann das naechste; mit einer Zahl der Person (z.B. '3 auf einmal') so viele
- [ ] SKILL.md: Check ist eine nummerierte Liste - Vorbereitung, je Schritt eine Handlung und was man danach sieht; Maschinen-Ergebnisse stehen davor als eigene Zeile
- [ ] Ein Beispiel fuer ein Ticket im neuen Format steht in SKILL.md
- [ ] Gegengelesen von Fable: die Check-Vorgabe ist fuer jemanden ohne Vorwissen verstaendlich (Notiz am Ticket)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

