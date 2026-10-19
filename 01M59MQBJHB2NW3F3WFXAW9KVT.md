---
id: 01M59MQBJHB2NW3F3WFXAW9KVT
title: "Human Review zeigt zuerst, wie der Mensch die Aenderung selbst testet"
status: backlog
ready: true
creator: berkc
assignee: berkc
goal: "Wer ein Ticket in Human Review oeffnet, sieht oben kurze Schritte zum Ausprobieren (was starten, was druecken, was sehen) und erst danach die Review-Details."
context: |-
  Gemeldet von berkc am 19.10.2026: 'die Human Review Lane ist schlecht, frueher standen da Beschreibungen, wie man testet'.
  Die Daten gibt es: die review-Lane muss review-check liefern (.jaira/lanes/review.md, output-produces). Der Sign-off-Screen zeigt es (internal/tui/signoff.go:83), das Detail auch (view.go:1375) - aber als letzten Punkt nach summary, gaps, verdict.
  Beispiel GK1QGT: check ist ein Block mit 9 Schritten voller Pfade wie /home/berkc/code/X_General/jaira-T8A8KM und go build -ldflags ... - fuer einen Entwickler, nicht zum schnellen Ausprobieren.
  Zu klaeren beim Bearbeiten: was berkc 'frueher' gesehen hat (git log von internal/tui/signoff.go und .jaira/lanes/review.md pruefen).
  Ausserdem: der Code liegt in einem anderen Worktree/Branch. Ohne einen fertigen Befehl zum Bauen dieses Branches kann man nichts testen (siehe X3A8RZ: berkc sah den Fix nicht, weil er nur im Branch ist).
definition-of-done: "Der Sign-off-Screen und das Detail einer Karte in Human Review zeigen review-check als ersten Block, als nummerierte Schritte; Test in internal/tui"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-19T08:33:51Z
updated-at: 2026-10-19T08:33:51Z
---

# Human Review zeigt zuerst, wie der Mensch die Aenderung selbst testet

## Definition of Done

- [ ] Der Sign-off-Screen und das Detail einer Karte in Human Review zeigen review-check als ersten Block, als nummerierte Schritte; Test in internal/tui
- [ ] Der Prompt der review-Lane verlangt fuer check: hoechstens ~6 Schritte, beginnend mit einem kopierbaren Befehl, der den Branch baut und startet, und je Schritt was man sehen soll
- [ ] Release-Note-Zeile in core/release/notes/<ID>.md

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

