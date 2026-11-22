---
id: 01M81G681S9YWPBZWQ7HAJ41SC
title: "Das Verdict im Cockpit zeigt oben eine Ampel: grün annehmen, gelb ansehen, rot zurückgeben"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Alex sieht im Verdict mit einem Blick, ob er annehmen kann, und warum nicht."
context: "Problem: im Verdict des Cockpits sieht Alex nicht auf einen Blick, ob er annehmen kann. Er muss Urteil, Gaps und Notizen lesen, und er liest langsam.\n- Wunsch Alex 22.11.2026: \"mir fehlt im Verdict eine Ampel, damit ich schneller annehmen kann\".\n- Das Verdict zeigt seit PDZ444 Titel + Ziel und das Review-Urteil; Code: core/mod/builtin/jaira-cockpit/hooks/register.tsx, digest.ts, board.ts, texts.ts.\n- Vorschlag für die Farben, die Daten liegen schon am Ticket:\n  - grün: review-verdict accept, Bughunt-Abschlusszeile auf HEAD, Gates ok, keine Notiz als Blocker sortiert, kein offener 'seen by a person:'-Punkt\n  - gelb: accept, aber eine 'seen by a person:'-Punkt offen oder note:-Zeilen in review-gaps (kleine Ränder/bekannte Grenzen)\n  - rot: review gibt zurück, Bughunt fehlt oder ist alt (nicht auf HEAD), ein Gate verweigert, oder eine Blocker-Notiz\n- Die Ampel ist eine Zeile ganz oben im Verdict, vor dem Titel, mit einem Wort Grund (\"rot: Bughunt nicht auf HEAD\").\n- Kein neues Feld am Ticket: die Farbe wird aus vorhandenen Feldern berechnet."
definition-of-done: Das Verdict zeigt als erste Zeile eine Ampel (grün/gelb/rot) mit einem kurzen Grund
tags:
  - acceptance
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-11-22T11:27:46Z
updated-at: 2026-11-22T11:33:59Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-1232016
claimed-at: 2026-11-22T11:29:15Z
outcome-what: "Verdict-Pane des Cockpits zeigt als erste Zeile eine Ampel (grün/gelb/rot) mit dem ersten Grund und '+N more'"
outcome-why: "Alex soll ohne Urteil, Gaps und Notizen zu lesen sehen, ob er annehmen kann"
outcome-resolves: "AJ41SC: fehlende Ampel im Verdict"
---

# Das Verdict im Cockpit zeigt oben eine Ampel: grün annehmen, gelb ansehen, rot zurückgeben

## Definition of Done

- [x] Das Verdict zeigt als erste Zeile eine Ampel (grün/gelb/rot) mit einem kurzen Grund
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1945; test 'a ticket in the band is one row' expects '● red: 1 blocker in the review (+4 more)' as first pane line
- [x] Die Farbe wird aus vorhandenen Ticketdaten berechnet (review-verdict, Bughunt auf HEAD, Gates, Notizen, seen by a person:), Regeln mit Test pro Farbe
  proof: core/mod/builtin/jaira-cockpit/hooks/board.ts:231 lightOf; core/mod/builtin/jaira-cockpit/tests/board.test.ts:3803 'the verdict light' green/yellow/red
- [x] Release-Note core/release/notes/<ID>.md, shipped.sha256
  proof: core/release/notes/AJ41SC.md; core/shipped/shipped.sha256 via scripts/shipped-history.sh --write; go test ./core/role ./core/shipped ok
- [ ] seen by a person: Alex sieht die Ampel im laufenden Cockpit an einem Ticket in signoff

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-22 11:29 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/25712.sock
- **2026-11-22 11:30 · Alexander Sacharov** — decision (worker, no open question): Ampel-Regeln. rot: Review lehnt ab (REFUSAL), Blocker in review-gaps, Bughunt fail (fehlt/nicht auf HEAD), ein DoD-Punkt ausser 'seen by a person:' offen (= das Gate des Terminal-Lanes). gelb: Review mit Vorbehalt, note:-Zeilen, offener 'seen by a person:'-Punkt, offene Frage, Review aus frueherer Runde, Notizen unbekannt, Bughunt unknown/unbekannt (aeltere jaira), noch kein Review. gruen: sonst. 'Gates' = DoD-Gate, weil der Worker-Statuszeile-gates-Wert nicht am Ticket liegt; Texte englisch wie texts.ts. Erster Grund in Reihenfolge rot vor gelb wird gezeigt, plus '+N' weitere.
