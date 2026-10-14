---
id: 01M4VGJJQZ9PZF2F19NCTQ5G6R
title: "docs/COMMANDS.md nennt 'jaira json get' und 'jaira show --brief'"
status: review
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
commits: []
created-at: 2026-10-13T20:51:58Z
updated-at: 2026-10-13T21:38:53Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-26486
claimed-at: 2026-10-13T21:36:55Z
outcome-what: "docs/COMMANDS.md fuehrt 'jaira json get <path> [--default]' unter Agent plumbing und '--brief' in der Zeile von 'jaira show <id>'"
outcome-why: Beide Befehle kamen mit TP44SC und standen nur in --help und der Release-Note
outcome-resolves: Review-Luecke 2 von TP44SC
review-summary: "Zwei Tabellenzeilen, beide in vorhandener Form: '--brief' neben '--notes-last' in der Zeile von 'jaira show <id>', wie es die DoD verlangt; 'json get' als eigene Zeile unter Agent plumbing wie tasks/checkpoint/hook print. Laenge vergleichbar mit der tags-Zeile. Kein einfacherer Zuschnitt, nichts Spekulatives, keine Release-Note noetig (reine Doku, Kontext sagt das)."
review-gaps: "Nichts zu entfernen: zwei Tabellenzeilen, keine Dopplung in docs/ oder README.md (grep nach 'json get' und '--brief'), kein Code."
test-verdict: "pass: frisch gebautes Binary; 'json get result.pane.pane_id' -> p1 rc 0, 'items.1' -> 2, fehlender Pfad rc 5, mit --default d rc 0; 'show TQ5G6R --brief' liefert genau id,title,goal,question,review,notes_last,notes_total wie docs/COMMANDS.md:58 und :170 sagen."
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
