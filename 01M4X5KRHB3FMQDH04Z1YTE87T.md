---
id: 01M4X5KRHB3FMQDH04Z1YTE87T
title: "Zwei Worker treffen sich auf einem Join-Board und einigen sich, wie ihre Teile zusammenpassen"
status: brainstorm
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Brauchen zwei verlinkte Tickets eine gemeinsame Entscheidung (Schnittstelle, Reihenfolge), treffen sich ihre Worker auf einem Join-Board, diskutieren und schreiben das Ergebnis auf beide Tickets"
context: |-
  Was heute fehlt: Zwei Agenten an zusammengehoerenden Tickets haben keinen Ort, an dem sie gemeinsam entscheiden. Jeder entscheidet fuer sich, und die Teile passen am Ende nicht.
  Berks Idee vom 14.10.: 'Man koennte auch einbauen, dass sich die beiden Worker auf einem Join-Board treffen, wo beide Agenten durch Diskussion entscheiden, was der beste Weg ist, sich zu vereinen, oder Sachen diskutieren.'
  Beispiel: Import-Button (ui-Board) und S3-Upload (backend-Board) einigen sich auf Endpunkt, Dateigroesse und Fehlerfaelle.
  Offen fuer brainstorm: ist das Join-Board ein eigenes Board mit Lanes (vorschlagen, einwenden, entscheiden) oder ein gemeinsames Ticket? Wer entscheidet bei Uneinigkeit (Mensch in der human-Lane)?
  Haengt an S0EG0M (mehrere Boards) und am Ticket fuer verlinkte Tickets mit Beruehrungspunkt.
definition-of-done: "Design als Spec, von Berk freigegeben: Form des Join-Boards, wer es startet, wie eine Einigung festgehalten wird, was bei Uneinigkeit passiert"
tags: []
blocked-by: []
related:
  - 01M4X5JKRN4AV3KX56SES0EG0M
commits: []
created-at: 2026-10-14T12:18:51Z
updated-at: 2026-11-23T22:45:49Z
updated-by: BeMuCa
claimed-by: teamlead-lane-graph
claimed-at: 2026-11-23T22:45:37Z
---

# Zwei Worker treffen sich auf einem Join-Board und einigen sich, wie ihre Teile zusammenpassen

## Definition of Done

- [ ] Design als Spec, von Berk freigegeben: Form des Join-Boards, wer es startet, wie eine Einigung festgehalten wird, was bei Uneinigkeit passiert
- [ ] Zwei Worker an verlinkten Tickets koennen auf dem Join-Board Vorschlaege austauschen; die Einigung steht danach auf beiden Tickets; Test
- [ ] Ohne Einigung landet die Frage in der human-Lane statt still entschieden zu werden

## Options

- [x] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-23 22:45 · BeMuCa** — Brainstorm (Teamlead, 2026-11-23). Was der Code zeigt: mehrere Boards je Repo gibt es (S0EG0M, done; .jaira/boards/<name>/ mit eigenen Lanes), Tickets kennen related und blocked-by, Loops mit max-rounds und die human-Lane mit Frage gibt es. Es fehlt nur der Ort und die Regel, wer dort wann schreibt. jaira startet keine Agenten - 'treffen' heisst: beide Worker schreiben nacheinander auf dasselbe Ticket.
Weg A - Join-Board: ein Board 'join' mit drei Lanes vorschlagen -> antworten -> geeinigt; antworten hat rejects-to vorschlagen mit max-rounds 3. Ein Join-Ticket (related zu beiden, beide blocked-by darauf) laeuft durch. Wer zuerst an den Beruehrungspunkt kommt, legt es an. Nach 3 Runden ohne Einigung greift die bestehende Loop-Grenze und es geht mit Frage in human. Bei 'geeinigt' schreibt der Worker die Einigung als Notiz auf beide Tickets. Kosten: ein mitgeliefertes Board mit 3 Lane-Dateien und ein Absatz im Worker-Skill, fast kein neuer Binary-Code; im Lane-Graph (WFHPZT) von selbst sichtbar. Gibt auf: nichts Erzwungenes - dass die Notiz auf beiden Tickets landet, steht im Prompt, kein Gate prueft es.
Weg B - wie A, plus ein Kommando 'jaira agree <join-id>', das die Einigung auf beide Tickets schreibt und die blocked-by loest, und ein Gate, das 'geeinigt' ohne das ablehnt. Kosten: ein Kommando, ein Gate, Tests. Gibt: DoD 2 ('steht danach auf beiden Tickets') ist geprueft statt erhofft.
Weg C - kein Board, nur ein gemeinsames Ticket in den normalen Lanes mit Notizen. Am kleinsten; aber die Runden sind nicht sichtbar, es gibt keine Grenze, und im Graph erscheint nichts - das Gegenteil von dem, was Berk am 14.10. und 23.11. beschrieben hat.
Empfehlung: B, in zwei Schritten gebaut (erst A als Lanes und Skill-Text, dann agree und Gate). DoD 1 verlangt Berks Freigabe der Form - deshalb human. Nicht am Code geprueft: ob blocked-by ueber Boardgrenzen hinweg greift.
