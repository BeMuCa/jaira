---
id: 01M5GFCE8T532CJF8NZMEZ3V6D
title: "Die Marken pull und behind verschwinden: Schreiben geht auch ohne Datei, und eine unveraenderte Datei holt jaira von selbst nach"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Niemand muss mehr jaira pull tippen, um ein Ticket zu lesen oder zu aendern; die Marken '⇢ pull', '[pull it]', '⇣ pull' und '[behind ref: pull it]' gibt es nicht mehr."
context: |-
  Wunsch von berkc am 22.10.2026: 'the pull it seems useless, throw it out. and regarding the fetching, Shouldnt we just updates the one that we do not work on? because this way we do not see the changes until its done right?'
  Was er loesen will: Board und Datei sollen nie auseinanderlaufen, und Aenderungen anderer sollen sofort sichtbar sein - ohne Handarbeit.
  Warum es die Marken heute gibt (geprueft):
  - '⇢ pull' / '[pull it]' (internal/tui/view.go:767, internal/cli/tickets.go:745): das Ticket liegt nur auf seinem Ref, keine Datei auf der Platte. Schreiben (note, move, dod) wird verweigert: 'on its ref and not on your disk; pull it first'.
  - '⇣ pull' / '[behind ref: pull it]' (K6BGB0, view.go:776): die Datei auf der Platte ist aelter als der Ref. Der 2-s-Takt liest nur Platte und lokale Refs, der 2-min-Fetch aktualisiert nur Refs; die Datei im Arbeitsbaum schreibt jaira nie selbst um, weil sie zum Branch gehoert und ungespeicherte Aenderungen haben kann.
  Loesung, die berkc mit 'nur die aktualisieren, an denen wir nicht arbeiten' meint:
  1. Schreiben auf ein Ticket ohne Datei geht direkt auf den Ref (kein pull noetig). Die Datei braucht man nur, wenn das Ticket mit dem Code committet wird - dann holt jaira sie beim ersten Schreiben in einem Branch, der nach dem Ticket heisst.
  2. Ist die Datei hinter dem Ref UND unveraendert gegenueber HEAD (niemand arbeitet gerade daran), schreibt jaira sie selbst auf den Ref-Stand nach.
  3. Nur wenn die Datei lokale Aenderungen hat UND der Ref weiter ist, bleibt es ein Konflikt - dann eine klare Meldung statt einer Dauermarke.
definition-of-done: jaira note/move/dod/claim auf ein Ticket ohne Datei schreiben auf den Ref; Test in internal/cli
tags:
  - tui
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-10-22T00:15:11Z
updated-at: 2026-10-22T00:15:11Z
---

# Die Marken pull und behind verschwinden: Schreiben geht auch ohne Datei, und eine unveraenderte Datei holt jaira von selbst nach

## Definition of Done

- [ ] jaira note/move/dod/claim auf ein Ticket ohne Datei schreiben auf den Ref; Test in internal/cli
- [ ] Eine Datei, die hinter ihrem Ref liegt und keine lokalen Aenderungen hat, wird beim naechsten Lesen nachgezogen; Test
- [ ] Die Marken '⇢ pull', '[pull it]', '⇣ pull', '[behind ref: pull it]' sind entfernt; nur der echte Konflikt (lokale Aenderung + Ref weiter) meldet sich; Test
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

