---
id: 01M85APCAKF9JG91BGMAMMY2CS
title: "jaira merkt sich, woher eine uebernommene Lane kommt, und zeigt, wenn der Katalog eine neuere hat"
status: backlog
ready: true
creator: BeMuCa
assignee: Berk
goal: "Fuer jede Lane aus dem Marktplatz weiss jaira Quelle und Pruefsumme, zeigt 'neuer im Katalog' oder 'von dir geaendert' an und ueberschreibt eine geaenderte Fassung nie still."
context: |-
  Heute ist eine uebernommene Lane danach eine Datei ohne Herkunft in ~/.jaira/lanes. Niemand weiss, aus welchem Stand sie kam.
  Es gibt kein Update: 'adopt --force' ueberschreibt, auch eine von Hand angepasste Fassung, ohne Hinweis.
  Denselben Verlust beschreibt VHS67E fuer eingebettete Lanes (angepasste Lane verliert still gegen die eingebettete Fassung). Dieses Ticket darf das Muster nicht wiederholen.
  Loesung:
  - Beim Uebernehmen schreibt jaira eine Zeile in ~/.jaira/lanes/.market: id, Quelle, Commit des Index, sha256. Eine Textdatei, von Hand lesbar. Die Lane-Datei selbst bleibt unveraendert.
  - 'jaira lanes market' markiert je Eintrag: uebernommen, neuer im Katalog (sha256 im Index anders als gemerkt), von dir geaendert (Datei anders als gemerkt).
  - 'jaira lanes market adopt <id> --update' holt die neue Fassung. Ist die lokale Datei geaendert, lehnt es ab und nennt beide Wege: --force oder die eigene Fassung unter anderer id sichern.
  - Kein automatisches Update, kein Hintergrund-Abruf.
  Entfernen: die Datei in ~/.jaira/lanes loeschen; eine verwaiste Zeile in .market wird beim naechsten Lauf entfernt. 'jaira lanes remove' nimmt eine Lane nur vom Board, nicht aus dem Katalog - das bleibt so.
  Vorgabe: docs/MARKET.md.
definition-of-done: adopt schreibt die Herkunftszeile; Test liest sie zurueck
tags:
  - cli
blocked-by:
  - 01M8BQFBMS778FE9WK6YDAW7WP
related:
  - 01M2XPBCFVPGPJ62NNN9VHS67E
commits: []
created-at: 2026-11-23T23:08:41Z
updated-at: 2026-11-26T10:51:15Z
updated-by: Berk
---

# jaira merkt sich, woher eine uebernommene Lane kommt, und zeigt, wenn der Katalog eine neuere hat

## Definition of Done

- [ ] adopt schreibt die Herkunftszeile; Test liest sie zurueck
- [ ] jaira lanes market zeigt die drei Zustaende, auch in --json
- [ ] adopt --update auf eine lokal geaenderte Lane endet mit Exit 3 und aendert nichts (Test)
- [ ] Eine geloeschte Lane hinterlaesst nach dem naechsten Lauf keine Zeile in .market
- [ ] Release-Note core/release/notes/<TICKET>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

