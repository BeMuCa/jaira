---
id: 01M4X5KGP2QVWD9MG9VP8NHD8Q
title: "Der Teamlead legt Tickets auf das Board, das zum Thema passt, und der Worker faehrt sie durch dessen Lanes"
status: backlog
ready: true
creator: BeMuCa
assignee: Alexander Sacharov
goal: "Aus einer Bitte wie 'in der Datenbank fuer user eine Spalte hinzufuegen' legt der Teamlead ein Ticket an, waehlt das passende Board (hier backend), und der jaira-worker bekommt Ticket und Board und treibt es durch dessen Lanes"
context: |-
  Was heute fehlt: Der Teamlead (Skill jaira-teamlead) legt Tickets auf dem einen Board des Repos an. Ein Board pro Thema gibt es noch nicht (S0EG0M baut es).
  Berks Idee vom 14.10.: 'Wir koennten dann z.B. unserem Teamlead sagen: hey, ich will das und jenes machen, und er erstellt Tickets und matcht sie auch den Boards, die zu dem Thema passen.' Beispiel: 'in unserer Datenbank sollten wir fuer user noch eine Spalte hinzufuegen' -> backend-Board; der jaira-worker kriegt das Ticket und das Board, durch das er es treiben soll.
  Offen fuer brainstorm: woran der Teamlead das Board erkennt (Beschreibung pro Board in seiner Board-Datei? Tags?), und was passiert, wenn kein Board passt (nachfragen statt raten).
  Haengt an S0EG0M (mehrere Boards pro Repo).
definition-of-done: "Jedes Board traegt eine kurze Beschreibung seines Themas, die der Teamlead liest; Test"
tags: []
blocked-by: []
related:
  - 01M4X5JKRN4AV3KX56SES0EG0M
commits: []
created-at: 2026-10-14T12:18:43Z
updated-at: 2026-10-17T12:36:37Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-65566
claimed-at: 2026-10-17T12:36:37Z
---

# Der Teamlead legt Tickets auf das Board, das zum Thema passt, und der Worker faehrt sie durch dessen Lanes

## Definition of Done

- [ ] Jedes Board traegt eine kurze Beschreibung seines Themas, die der Teamlead liest; Test
- [ ] Der Teamlead legt ein Ticket mit Board an und nennt seine Wahl mit Grund; passt kein Board, fragt er den Menschen statt zu raten
- [ ] jaira-worker und Dispatcher lesen das Board vom Ticket und arbeiten dessen Lanes ab, nicht die des Standard-Boards; Test
- [ ] Eine Zeile in core/release/notes/8NHD8Q.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-17 12:36 · Alexander Sacharov** — Alexander Sacharov took this ticket over from BeMuCa
- **2026-10-17 12:36 · Alexander Sacharov** — decision 17.10. (Alexander, Teamlead-Session): Alexander uebernimmt das Ticket von Berk. Der Branch feat/8nhd8q zweigt von feat/s0eg0m ab, weil er --board und das Feld board: aus S0EG0M braucht; S0EG0M ist beim Abzweigen in review. Wird S0EG0M nach master gemergt, master hier einmergen. Ziel dahinter: in ~/projects/demo-video sagt eine Designerin ohne Programmierkenntnis dem Teamlead 'mach ein Video aus dieser Demo', und der legt das Ticket auf das Board video statt auf das SWE-Board.
