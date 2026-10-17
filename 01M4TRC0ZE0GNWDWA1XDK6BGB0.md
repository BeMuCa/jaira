---
id: 01M4TRC0ZE0GNWDWA1XDK6BGB0
title: "Das Board zeigt den Stand des Refs, wenn die Ticket-Datei dahinter zurueckliegt"
status: backlog
ready: true
creator: BeMuCa
goal: "Ein Mensch sieht auf dem Board die Lane, in der ein Ticket wirklich steht, auch wenn die Datei in seinem Checkout veraltet ist."
context: |-
  Am 2026-10-13 stand D0SAHM auf Berks Board (Haupt-Checkout /home/berk/git/jAIra) in in-progress. Sein Ref refs/jaira/tickets/01M1KFMETY9MKH4V38TVD0SAHM sagte schon signoff.
  Grund: Board und 'jaira list' lesen die Datei im Checkout, und die lag hinter dem Ref. Gearbeitet wurde das Ticket in einem anderen Worktree (/home/berk/git/jAIra-D0SAHM).
  'jaira pull D0SAHM' wusste es: 'D0SAHM was behind its ref and is up to date now - was: in-progress ... now: signoff'. jaira erkennt den Fall also schon.
  Es war das zweite Mal in derselben Session. Berk haette fast ein Ticket fuer haengend gehalten, das auf seine Abnahme wartete.
  JDXV60 beschreibt dasselbe Problem zwischen zwei Worktrees (2026-09-15), vor den Refs. Das hier ist der Rest, der mit Refs bleibt.
  Offen fuer brainstorm: Board liest dann den Ref, oder zeigt nur ein Kennzeichen 'Ref ist weiter' an der Karte.
definition-of-done: "Liegt die Datei eines Tickets hinter seinem Ref, zeigt das Board auf der Karte den Stand des Refs oder ein Kennzeichen dafuer; 'jaira list' und 'jaira show' sagen es ebenso; ein Test baut genau diesen Fall"
tags:
  - concurrency
  - tui
blocked-by: []
related: []
follows: 01M2HW4788Z5G0M05G3PJDXV60
commits: []
created-at: 2026-10-13T13:48:57Z
updated-at: 2026-10-17T20:39:55Z
assignee: berkc
updated-by: berkc
claimed-by: berk-ca-297639
claimed-at: 2026-10-17T20:39:34Z
---

# Das Board zeigt den Stand des Refs, wenn die Ticket-Datei dahinter zurueckliegt

## Definition of Done

- [ ] Liegt die Datei eines Tickets hinter seinem Ref, zeigt das Board auf der Karte den Stand des Refs oder ein Kennzeichen dafuer; 'jaira list' und 'jaira show' sagen es ebenso; ein Test baut genau diesen Fall
- [ ] go test ./... -race gruen; Release-Notiz unter core/release/notes/

## Options

- [x] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-14 14:08 · BeMuCa** — Dritte Variante am 14.10.: D0SAHM stand auf Berks Board in signoff, obwohl es laengst abgenommen und abgelegt war (logbook bemuca-20261013, per PR #74 in master). Grund: Der Haupt-Checkout /home/berk/git/jAIra stand noch auf dem alten master fa222c5d, und dort lag eine veraltete, geaenderte Kopie .jaira/tickets/...D0SAHM...md mit status: signoff. Den Ref refs/jaira/tickets/<ULID> gibt es nach dem Ablegen nicht mehr, also kann 'zeig den Stand des Refs' hier nicht helfen. Das Board muss auch erkennen: Datei ohne Ref, deren Ticket upstream im Logbuch liegt -> als abgelegt bzw. veraltet zeigen statt als offen. Behoben fuer heute durch git stash der veralteten Kopien + git pull master.
- **2026-10-17 20:39 · berkc** — berkc took this ticket over from BeMuCa
- **2026-10-17 20:39 · berkc** — Brainstorm (2026-10-17, entschieden vom Nutzer): Was falsch ist laut Code: Store.List zeigt fuer jede Ticket-Datei die Datei; Syncer.Extra (core/refsync/refsync.go) ueberspringt jeden Ref, fuer den eine Datei existiert ('if local[id] || filed[id] continue'). Liegt die Datei hinter dem Ref, zeigen Board, list und show die alte Lane. 'jaira pull' erkennt den Fall schon mit isStale (refsync.go: ref.updated_at strikt nach file.updated_at) - dieselbe Definition wie CheckFresh. Schreibseite: Mutate/Archive/Logbook rufen schon CheckFresh und verweigern mit StaleError ('run jaira pull'), ein Schreiben auf eine veraltete Datei kann den Ref also heute schon nicht still zurueckdrehen.
Optionen: (A) nur Kennzeichen an der Karte, Datei-Stand bleibt sichtbar - billig, aber Board zeigt weiter die falsche Lane. (B) Board/list/show zeigen den Stand des Refs plus Kennzeichen 'Datei veraltet' - zeigt die echte Lane, Schreiben bleibt an der Datei und wird von CheckFresh verweigert. (C) automatisch pullen beim Lesen - schreibt Dateien ohne Befehl, verletzt 'Lesen schreibt nie'.
Gewaehlt: B (Nutzerentscheidung 2026-10-17): newest wins in List und show, Marker an Karte/Zeile/JSON. Load bleibt bei der Datei, weil Archive/Logbook/move den geladenen Ticket gegen CheckFresh pruefen - ein Ref-Ticket dort wuerde die Pruefung bestehen und die alte Datei ablegen.
