---
id: 01M5659BR7RYZ6RVC84HGK1QGT
title: "Die Release-Notes zeigen eine Version pro Seite, die neueste zuerst, und links/rechts blaettert zwischen den Versionen"
status: in-progress
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
updated-at: 2026-10-18T00:40:02Z
updated-by: berkc
claimed-by: berk-ca-703016
claimed-at: 2026-10-18T00:32:08Z
---

# Die Release-Notes zeigen eine Version pro Seite, die neueste zuerst, und links/rechts blaettert zwischen den Versionen

## Definition of Done

- [x] Popup und 'w' zeigen genau eine Version pro Seite, die neueste zuerst, mit Kopfzeile Version und Position (z.B. 0.5.4 · 1/12)
  proof: internal/tui/releasenotes.go:90 header + lines() one entry; TestNotesViewStartsOnTheNewestAndPagesOneVersionAtATime, TestHomeWOpensReleaseNotesNewestFirst
- [x] Pfeil links/rechts (und h/l) blaettert zur aelteren/neueren Version, an den Enden passiert nichts
  proof: internal/tui/releasenotes.go:42-49; TestNotesViewStartsOnTheNewestAndPagesOneVersionAtATime (ends: right at newest, left at oldest)
- [x] Nach einem Upgrade ueber mehrere Versionen sieht man, wie viele davon ungelesen sind
  proof: internal/tui/releasenotes.go:164 title '· N unread'; TestPopupAfterAnUpgradeAcrossVersionsCountsTheUnreadAndPagesThem
- [x] Test pinnt Startseite und Blaettern; Release-Note
  proof: internal/tui/releasenotes_test.go:39,201; internal/tui/home_test.go:416; core/release/notes/GK1QGT.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-18 00:32 · berkc** — Decided (open point in context): the popup after an upgrade pages only the unread versions (release.Since(base)), starts on the newest, and its title says how many: 'What is new since 0.5.0 · 3 unread'; the page header is '0.5.4 · 1/3'. 'w' pages every version, '0.5.3 · 1/12'. Why: the popup's job is the unread, 'w' already keeps the whole history one key away; marking unread pages inside a 12-page popup would add a second concept for the same count. Direction per DoD 2: left/h = older (next page number), right/l = newer; ends do nothing. j/k scrolling stays for a version longer than the screen and resets on a page change. Brainstorm and planning options unticked, so the ticket goes todo -> in-progress directly.
- **2026-10-18 00:40 · berkc** — Found out: release.Since(base) includes the Unreleased entry whenever note files sit directly in core/release/notes/, so a local build (tmux, -X main.version=0.5.3, stamp 0.5.0) says '4 unread' with an 'Unreleased' first page. A cut release has no loose note files, so real builds count only versions. Left as is: Since is T8A8KM's selection, not this ticket's.
