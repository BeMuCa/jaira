---
id: 01M5ZJ4QSWY1WCTRR5QR195QXK
title: "Ein Worker uebergibt an einen Nachfolger, bevor ihm der Kontext ausgeht"
status: backlog
ready: false
creator: BeMuCa
goal: "Die Rolle jaira-worker nennt eine Kontext-Grenze; erreicht ein Worker sie, schreibt er eine Uebergabe (Notizen, Branch, offene DoD, naechster Schritt) auf das Ticket und stoppt, und der Dispatcher startet einen frischen Worker aus dieser Uebergabe."
context: |-
  Was schief ist: ein Worker mit 245k Kontext meldete dreimal denselben Status statt zu arbeiten (Teamlead-Session 27.10.2026, Feedback L315/L328). Erst ein frischer Worker mit handgeschriebenem Uebergabebrief kam weiter.
  Heute: core/role/builtin/jaira-worker/SKILL.md kennt keine Grenze und keine Uebergabe; jaira-dispatcher startet keinen Nachfolger.
  Was es schon gibt: jaira note bei jeder Pause (nach Abstuerzen war der Stand aus Ticket und Branch rekonstruierbar, Feedback L333/O25) - die Uebergabe ist eine letzte, vollstaendige Notiz plus Stopp.
  Die Grenze kann jaira nicht messen; der Worker schaetzt sie selbst (Zahl der Lanes, Dateigroesse der Antworten) oder der Dispatcher setzt eine Zahl in den Brief.
  Zusammenhang: ATQFVJ (Subagent-Namen), X2TGXG (Worker findet Teamlead wieder).
definition-of-done: "jaira-worker/SKILL.md nennt die Grenze und den Inhalt der Uebergabe-Notiz (Branch, Stand je DoD, naechster Schritt, Fallen)"
tags:
  - docs
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-10-27T20:52:01Z
updated-at: 2026-10-27T20:52:01Z
---

# Ein Worker uebergibt an einen Nachfolger, bevor ihm der Kontext ausgeht

## Definition of Done

- [ ] jaira-worker/SKILL.md nennt die Grenze und den Inhalt der Uebergabe-Notiz (Branch, Stand je DoD, naechster Schritt, Fallen)
- [ ] jaira-dispatcher startet aus einer Uebergabe-Notiz einen Nachfolger, der ohne Rueckfrage weiterarbeitet
- [ ] Auf einem Scratch-Board nachgestellt: ein Worker stoppt mit Uebergabe, der Nachfolger beendet die Lane

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

