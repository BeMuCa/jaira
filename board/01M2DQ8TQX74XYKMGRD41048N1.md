---
id: 01M2DQ8TQX74XYKMGRD41048N1
title: Board hakt beim Pfeiltasten-Navigieren alle 2 Sekunden
status: backlog
ready: false
creator: Alexander Sacharov
goal: "Pfeiltasten-Navigation im Board bleibt flüssig: kein periodisches Stocken mehr, weil Disk- und Git-Arbeit nicht mehr im Update-Loop läuft."
context: |-
  Beim Durchlaufen der Tickets mit den Pfeiltasten stockt das Board alle 2-3 Sekunden spürbar (gemeldet auf 0.1.4).

  Ursache ist sehr wahrscheinlich der periodische Rescan im Bubble-Tea-Update-Loop:
  - internal/tui/model.go:783 tick() feuert alle 2 Sekunden einen tickMsg.
  - internal/tui/model.go:803 ruft darauf m.reload() SYNCHRON in Update() auf. Solange reload läuft, wird keine Taste verarbeitet - genau das Stocken.
  - m.reload() (internal/tui/model.go:345) macht pro Tick: lane.Load, tag.Load, m.store.List() (alle Ticket-Dateien lesen + YAML parsen), session.Load, refMarks, rebuild, refreshLiveBoards.
  - refreshLiveBoards (internal/tui/model.go:387) geht ZUSÄTZLICH über alle anderen aufgezeichneten Boards und liest deren Session-Dateien - boardHasLiveSession. Liegt eines davon unter /mnt/c oder auf einem nicht erreichbaren Pfad, kostet das unter WSL2 richtig Zeit.
  - Derselbe synchrone reload hängt auch an changeMsg (fsnotify) und refFetchedMsg - also nicht nur am Timer.

  Noch nicht gemessen: wie lange ein reload auf diesem Board tatsächlich dauert und welcher Teil dominiert. Das ist der erste Schritt.

  Bekannt: tick() ist Absicht - es ist der Backstop für unzuverlässige fsnotify-Events unter WSL2. Der Timer darf also nicht einfach weg; die Arbeit muss aus dem Update-Loop raus (tea.Cmd in einer Goroutine, Ergebnis als Msg zurück), plus Debounce statt Rescan pro Event.
definition-of-done: |-
  Rescan läuft nicht mehr synchron in Update(): tickMsg/changeMsg/refFetchedMsg stoßen einen tea.Cmd an, der im Hintergrund lädt und das Ergebnis als Msg zurückgibt.
  Überlappende Rescans werden unterdrückt - ein laufender Reload verhindert, dass der nächste Tick einen zweiten startet.
  Cursor, Filter und offenes Detail-Pane überstehen einen Hintergrund-Reload unverändert.
  Messung vor/nach der Änderung ist im Ticket notiert: wie lange ein reload dauert und welcher Teil dominiert.
  Pfeiltasten-Navigation über ein Board mit mehreren Dutzend Tickets stockt nicht mehr spürbar.
tags:
  - tui
  - concurrency
blocked-by: []
related: []
commits: []
created-at: 2026-09-13T15:48:13Z
updated-at: 2026-09-13T15:48:13Z
---

# Board hakt beim Pfeiltasten-Navigieren alle 2 Sekunden

## Definition of Done

- [ ] Rescan läuft nicht mehr synchron in Update(): tickMsg/changeMsg/refFetchedMsg stoßen einen tea.Cmd an, der im Hintergrund lädt und das Ergebnis als Msg zurückgibt.
Überlappende Rescans werden unterdrückt - ein laufender Reload verhindert, dass der nächste Tick einen zweiten startet.
Cursor, Filter und offenes Detail-Pane überstehen einen Hintergrund-Reload unverändert.
Messung vor/nach der Änderung ist im Ticket notiert: wie lange ein reload dauert und welcher Teil dominiert.
Pfeiltasten-Navigation über ein Board mit mehreren Dutzend Tickets stockt nicht mehr spürbar.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

