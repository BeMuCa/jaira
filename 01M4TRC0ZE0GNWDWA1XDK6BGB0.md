---
id: 01M4TRC0ZE0GNWDWA1XDK6BGB0
title: "Das Board zeigt den Stand des Refs, wenn die Ticket-Datei dahinter zurueckliegt"
status: pre-process
ready: true
creator: BeMuCa
goal: "Board, 'jaira list' und 'jaira show' zeigen den neuesten Stand eines Tickets: liegt die Datei im Checkout hinter ihrem Ref, den des Refs, mit einem Kennzeichen, dass 'jaira pull' faellig ist - und ein Schreiben auf die veraltete Datei bleibt verweigert."
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
updated-at: 2026-10-17T20:40:30Z
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
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [ ] ticket: Feld BehindRef am Ticket (Laufzeit, kein Frontmatter) und optionales Interface AheadSource; Store.List ersetzt Datei-Tickets, deren Ref weiter ist, durch das Ref-Ticket (Path bleibt die Datei, BehindRef=true); Store.Newest(t) fuer show
- [ ] refsync: Syncer.Ahead(files) - ListHashes + extraCache nach Ref-Hash + ReadMany nur fuer neue Hashes, Urteil per isStale (dieselbe Definition wie pull/CheckFresh)
- [ ] Ausgabe: TUI-Karte Kennzeichen 'pull due', jaira list Zeile [behind ref - pull], JSON behind-ref, show-Text Zeile; Load bleibt Datei
- [ ] Tests: refsync (Datei hinter Ref -> List zeigt Ref-Lane+BehindRef; Mutate verweigert mit ErrStale, Ref unveraendert), TUI-Karte zeigt Kennzeichen
- [ ] go vet + go test -race; manuell mit zwei Worktrees eines Temp-Repos; Release-Notiz core/release/notes/K6BGB0.md

## Progress
- **2026-10-14 14:08 · BeMuCa** — Dritte Variante am 14.10.: D0SAHM stand auf Berks Board in signoff, obwohl es laengst abgenommen und abgelegt war (logbook bemuca-20261013, per PR #74 in master). Grund: Der Haupt-Checkout /home/berk/git/jAIra stand noch auf dem alten master fa222c5d, und dort lag eine veraltete, geaenderte Kopie .jaira/tickets/...D0SAHM...md mit status: signoff. Den Ref refs/jaira/tickets/<ULID> gibt es nach dem Ablegen nicht mehr, also kann 'zeig den Stand des Refs' hier nicht helfen. Das Board muss auch erkennen: Datei ohne Ref, deren Ticket upstream im Logbuch liegt -> als abgelegt bzw. veraltet zeigen statt als offen. Behoben fuer heute durch git stash der veralteten Kopien + git pull master.
- **2026-10-17 20:39 · berkc** — berkc took this ticket over from BeMuCa
- **2026-10-17 20:39 · berkc** — Brainstorm (2026-10-17, entschieden vom Nutzer): Was falsch ist laut Code: Store.List zeigt fuer jede Ticket-Datei die Datei; Syncer.Extra (core/refsync/refsync.go) ueberspringt jeden Ref, fuer den eine Datei existiert ('if local[id] || filed[id] continue'). Liegt die Datei hinter dem Ref, zeigen Board, list und show die alte Lane. 'jaira pull' erkennt den Fall schon mit isStale (refsync.go: ref.updated_at strikt nach file.updated_at) - dieselbe Definition wie CheckFresh. Schreibseite: Mutate/Archive/Logbook rufen schon CheckFresh und verweigern mit StaleError ('run jaira pull'), ein Schreiben auf eine veraltete Datei kann den Ref also heute schon nicht still zurueckdrehen.
Optionen: (A) nur Kennzeichen an der Karte, Datei-Stand bleibt sichtbar - billig, aber Board zeigt weiter die falsche Lane. (B) Board/list/show zeigen den Stand des Refs plus Kennzeichen 'Datei veraltet' - zeigt die echte Lane, Schreiben bleibt an der Datei und wird von CheckFresh verweigert. (C) automatisch pullen beim Lesen - schreibt Dateien ohne Befehl, verletzt 'Lesen schreibt nie'.
Gewaehlt: B (Nutzerentscheidung 2026-10-17): newest wins in List und show, Marker an Karte/Zeile/JSON. Load bleibt bei der Datei, weil Archive/Logbook/move den geladenen Ticket gegen CheckFresh pruefen - ein Ref-Ticket dort wuerde die Pruefung bestehen und die alte Datei ablegen.
- **2026-10-17 20:40 · berkc** — Plan-Begruendung: Load wird NICHT umgestellt. Archive, Logbook und der move-Gate laden per Load und pruefen danach CheckFresh(t) - wuerde Load das Ref-Ticket liefern, bestuende es die Pruefung (gleiches updated_at) und Archive wuerde die alte Datei ablegen und deren Stand per RecordFiled auf den Ref schreiben. Deshalb nur List (Board, list, next, web) und show ueber Store.Newest. Mutate liest unter Lock die Datei neu und ruft CheckFresh -> StaleError; das bleibt die Schreib-Sperre. Separates optionales Interface statt Extra-Signatur zu aendern: ein zusaetzliches 'git for-each-ref' pro List, Decodes nur bei neuem Ref-Hash (extraCache geteilt). Nicht abgedeckt: D0SAHM-Variante (Datei ohne Ref, Ticket upstream abgelegt) - dort gibt es keinen Ref zum Vergleichen.
