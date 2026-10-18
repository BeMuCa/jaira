---
id: 01M5659BR7RYZ6RVC84HGK1QGT
title: "Die Release-Notes zeigen eine Version pro Seite, die neueste zuerst, und links/rechts blaettert zwischen den Versionen"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Wer die Release-Notes liest (Popup nach einem Upgrade und 'w' im Launcher), sieht zuerst nur die neueste Version kurz und knapp und blaettert mit links/rechts zu aelteren und zurueck."
context: |-
  Berk am 18.10.: 'Ich wuerde es pro Version aufteilen. Standardmaessig die neueste sehen, kurz, kompakt, praezise. Dann mit links oder rechts zwischen den Versionen navigieren.'
  Heute (PR #91, Tickets T8A8KM/PWAF22): das Popup nach einem Upgrade zeigt alle Notes seit der zuletzt gelesenen Version in einer langen Liste (internal/tui/releasenotes.go, openReleaseNotes/newNotesView). 'w' im Launcher zeigt alle Versionen untereinander, neueste zuerst.
  Die Daten sind schon pro Version da: release.Notes() liefert []Entry{Version, Changes}, neueste zuerst (core/release/release.go).
  Abgelehnt: eine feste Zahl Tabs (zuerst als 'letzte 5 Updates als Tabs' besprochen) - Blaettern pro Version deckt beliebig viele ab.
  Offen: was das Popup nach einem Upgrade ueber mehrere Versionen zeigt - Vorschlag: es startet auf der neuesten und sagt '1/3', die ungelesenen sind die ersten Seiten; Zeilen mit '!' bleiben hervorgehoben.
  Haengt an PR #91: baut auf dessen releasenotes.go auf, erst nach dessen Merge anfangen.
definition-of-done: "Popup und 'w' zeigen genau eine Version pro Seite, die neueste zuerst, mit Kopfzeile Version und Position (z.B. 0.5.4 · 1/12)"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-18T00:06:20Z
updated-at: 2026-10-18T00:32:00Z
updated-by: berkc
---

# Die Release-Notes zeigen eine Version pro Seite, die neueste zuerst, und links/rechts blaettert zwischen den Versionen

## Definition of Done

- [ ] Popup und 'w' zeigen genau eine Version pro Seite, die neueste zuerst, mit Kopfzeile Version und Position (z.B. 0.5.4 · 1/12)
- [ ] Pfeil links/rechts (und h/l) blaettert zur aelteren/neueren Version, an den Enden passiert nichts
- [ ] Nach einem Upgrade ueber mehrere Versionen sieht man, wie viele davon ungelesen sind
- [ ] Test pinnt Startseite und Blaettern; Release-Note

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

