---
id: 01M4VGJJQZ9PZF2F19NCTQ5G6R
title: "docs/COMMANDS.md nennt 'jaira json get' und 'jaira show --brief'"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die Befehlsreferenz docs/COMMANDS.md fuehrt die zwei Befehle, die TP44SC eingefuehrt hat, so dass man sie dort findet und nicht nur in der Release-Note"
context: |-
  Was fehlt: TP44SC (PR #80, Issue #79) hat 'jaira json get <pfad> [--default <v>]' und 'jaira show <id> --brief' eingefuehrt. docs/COMMANDS.md nennt keinen der beiden.
  Ausloeser: Review-Luecke 2 von TP44SC am 13.10.; Alex hat am 13.10. ein eigenes Ticket dafuer gewollt.
  Wohin: 'jaira json get' in den Abschnitt '## Agent plumbing' (docs/COMMANDS.md ~162); '--brief' in die Zeile von 'jaira show <id>' (~58), neben '--notes-last'.
  Wortlaut: aus 'jaira json get --help' und 'jaira show --help' des TP44SC-Binaries und aus core/release/notes/TP44SC.md uebernehmen, nichts erfinden.
  Abhaengigkeit: erst nach dem Merge von PR #80 auf master bauen, sonst gibt es die Befehle auf dem Zweig nicht.
  Keine Release-Note: reine Doku, fuer Nutzer des Binaries nicht sichtbar.
definition-of-done: "docs/COMMANDS.md hat unter 'Agent plumbing' eine Zeile fuer 'jaira json get <pfad> [--default <v>]' (liest JSON von stdin, druckt ein Feld, Exit 5 ohne Treffer)"
tags:
  - docs
blocked-by: []
related: []
follows: 01M4VCZS1N1T0BHQVR8ATP44SC
commits:
  - 48733322aebaf8f507e12ba0b9922ccac2bafbaa
created-at: 2026-10-13T20:51:58Z
updated-at: 2026-10-14T17:28:01Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-34835
claimed-at: 2026-10-14T17:16:18Z
outcome-what: "docs/COMMANDS.md fuehrt 'jaira json get <path> [--default]' unter Agent plumbing und '--brief' in der Zeile von 'jaira show <id>'"
outcome-why: Beide Befehle kamen mit TP44SC und standen nur in --help und der Release-Note
outcome-resolves: Review-Luecke 2 von TP44SC
review-summary: "Ein Commit (48733322), nur docs/COMMANDS.md, zwei Tabellenzeilen. Zeile 58: 'jaira show <id>' nennt jetzt '--brief' - eine JSON-Zeile mit id, title, goal, question, review, den neuesten vier Notizen als notes_last und ihrer Anzahl als notes_total. Zeile 170 unter Agent plumbing: neue Zeile 'jaira json get <path>' - liest ein JSON-Dokument von stdin, druckt den Wert am Punktpfad (Zahl fuer Array-Element), String ohne Anfuehrungszeichen, null druckt nichts, sonst kompaktes JSON; fehlender Pfad oder kein JSON endet mit Exit 5, ausser --default <v> nennt den Ersatz. Beide Texte sind fast wortgleich mit 'show --help' und 'json get --help'."
review-gaps: "none. Geprueft: Wortlaut gegen beide --help-Texte; Verhalten per go run: a.b.1 -> x, fehlender Pfad und Nicht-JSON -> exit status 5, null -> leere Zeile, --default q -> q; show --brief liefert die genannten Schluessel. Commit-Liste vollstaendig (master..HEAD = 1 Commit, nennt TQ5G6R). Keine Release-Note noetig (reine Doku). Nebenbei, nicht Teil dieser Lane: der TP44SC-Logbook-Commit soll laut Notiz vom 13.10. noch vor dem Push in diesen Branch."
test-verdict: "pass: frisch gebautes Binary; 'json get result.pane.pane_id' -> p1 rc 0, 'items.1' -> 2, fehlender Pfad rc 5, mit --default d rc 0; 'show TQ5G6R --brief' liefert genau id,title,goal,question,review,notes_last,notes_total wie docs/COMMANDS.md:58 und :170 sagen."
review-verdict: "Der Diff erfuellt alle drei DoD-Punkte; keine Defekte gefunden. Kleinigkeit ohne Handlungsbedarf: die Befehlsspalte zeigt '[--default <v>]' nicht, der Flag steht aber in der Beschreibung."
review-round: 1
review-check: "1. sed -n '58p;170p' docs/COMMANDS.md  2. go run ./cmd/jaira show --help und go run ./cmd/jaira json get --help daneben lesen - gleiche Felder, Punktpfad, Exit 5, --default  3. echo '{\"a\":1}' | go run ./cmd/jaira json get b; echo $? - Meldung und 'exit status 5'  4. echo '{\"a\":1}' | go run ./cmd/jaira json get b --default q - druckt q  5. go run ./cmd/jaira show TQ5G6R --brief - eine Zeile mit id, title, goal, question, review, notes_last, notes_total"
---

# docs/COMMANDS.md nennt 'jaira json get' und 'jaira show --brief'

## Definition of Done

- [x] docs/COMMANDS.md hat unter 'Agent plumbing' eine Zeile fuer 'jaira json get <pfad> [--default <v>]' (liest JSON von stdin, druckt ein Feld, Exit 5 ohne Treffer)
  proof: docs/COMMANDS.md:170
- [x] Die Zeile 'jaira show <id>' nennt '--brief' (eine JSON-Zeile mit title, goal, question, review, neuesten vier Notizen)
  proof: docs/COMMANDS.md:58
- [x] Beide Texte stimmen mit '--help' des Binaries ueberein
  proof: docs/COMMANDS.md:58,170 gegen 'go run ./cmd/jaira show --help' und 'json get --help' gelesen: gleiche Felder, Pfadsyntax, Exit 5, --default

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-13 21:36 · Alexander Sacharov** — parent session: jaira-6a
- **2026-10-13 21:37 · Alexander Sacharov** — Offene Entscheidungen: keine. PR #80 ist auf master (e9810d6e enthaelt den Merge), Wortlaut kommt aus --help und core/release/notes/0.4.3/TP44SC.md.
- **2026-10-13 21:37 · Alexander Sacharov** — Wortlaut nah an --help gehalten statt an der Release-Note: die Note spricht den Updater an ('Read a single field...'), die Tabelle beschreibt den Befehl. json get steht hinter merge-driver, weil es wie die anderen Plumbing-Zeilen fuer Skripte ist. Kein Test liest COMMANDS.md (grep nach COMMANDS.md in *.go: leer).
- **2026-10-13 21:38 · Alexander Sacharov** — Teamlead 13.10.: TP44SC (Vorgaenger, Branch feat/tp44sc schon gemergt) wartet in signoff auf Alex' Annahme am 14.10. Danach faehrt sein Logbook-Commit in DIESEM Branch mit: 'jaira pull TP44SC && jaira logbook TP44SC', Commit 'chore(TP44SC): file the ticket away with its commits stamped', vor dem Push/PR. Nicht vorher - TP44SC ist noch nicht done.
- **2026-10-14 17:16 · Alexander Sacharov** — parent session: jaira-34 (neu gestartet; Ticket stand in review, Review-Lane noch nicht gelaufen)
- **2026-10-14 17:25 · Alexander Sacharov** — explain: Der Review hat keinen Defekt gefunden: der eine Commit ergaenzt docs/COMMANDS.md um 'jaira show --brief' (Zeile 58) und 'jaira json get' (Zeile 170), Wortlaut und Verhalten passen zu den --help-Texten. Einzige Kleinigkeit: die Befehlsspalte zeigt '[--default <v>]' nicht, der Flag steht nur in der Beschreibung.
Warum: Die Referenz soll die Befehle so zeigen, dass man sie dort findet; ein Flag nur im Beschreibungstext faellt beim Ueberfliegen der Spalte weg, bricht aber nichts.
1) '[--default <v>]' in die Befehlsspalte von 'jaira json get' aufnehmen und zurueckgeben [give back]
2) Ticket so annehmen, die Beschreibung nennt den Flag bereits [accept]
Empfehlung: 2 - der Flag steht in der Beschreibung, der Reviewer sieht keinen Handlungsbedarf, ein weiterer Durchlauf kostet mehr als die Kleinigkeit wert ist.
- **2026-10-14 17:27 · Alexander Sacharov** — decision: 2) Ticket so annehmen, die Beschreibung nennt den Flag bereits
