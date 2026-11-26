---
id: 01M86H5QJWEKFHXJXA710TREV1
title: "reap.sh:418 reads any failed 'herdr pane get' as a closed tab, not only…"
status: in-progress
ready: true
creator: Sascha
context: |-
  Was falsch ist: core/role/builtin/jaira-dispatcher/scripts/reap.sh:418 - schlaegt 'herdr pane get <pane>' fehl, bleibt tab leer, und reap.sh behandelt den Tab als schon geschlossen.
  Folge: hat herdr nur einen Moment nicht geantwortet (Timeout, Neustart), entfernt reap.sh Worktree und Branch, obwohl claude im Tab noch darin arbeitet.
  Der Check herdr_up (reap.sh:133, 'tab list') laeuft einmal am Anfang der Runde; ein Ausfall danach faellt nicht auf.
  8NWE4P (24.11.) hat nur den Fall 'hc get antwortet bei tab list nicht' abgefangen (Exit 3 -> reap behaelt alles).
  Fehlerformate: Herdr App (hc) schreibt {"error":{"kind":"not_found"}} auf stderr, Exit 3; altes herdr.exe kennt den Code pane_not_found (aus strings herdr.exe).
  Gleiches Muster, nicht dieses Ticket: close-tab.sh:159 (dort parkt es nur aus, entfernt nichts).
tags:
  - cockpit
blocked-by: []
related: []
follows: 01M866HXTHQS29K2ZMNQ8NWE4P
commits: []
created-at: 2026-11-24T10:21:10Z
updated-at: 2026-11-26T07:59:46Z
assignee: Sascha
updated-by: Sascha
claimed-by: DESKTOP-RFTCH11-35248
claimed-at: 2026-11-26T07:54:32Z
goal: "reap.sh liest ein gescheitertes 'herdr pane get' nur dann als geschlossenen Tab, wenn herdr den Pane nicht kennt (not_found); jede andere Fehlschlag (Timeout, herdr weg) behaelt Tab, Worktree und Record fuer die naechste Runde."
---

# reap.sh:418 reads any failed 'herdr pane get' as a closed tab, not only…

## Definition of Done

- [x] reap.sh:418 schliesst nur bei not_found/pane_not_found auf einen geschlossenen Tab; jeder andere Fehlschlag von 'pane get' gibt 'kept <ticket>: ...' aus und laesst Worktree, Branch und Record stehen; Test in core/role/reap_test.go mit einem herdr-Stub, dessen pane get ohne not_found scheitert
  proof: reap.sh hget() + pane/tab get in the close loop; TestReapReadsOnlyNotFoundAsAClosedTab (timeout, silent, tab get timeout: worktree, branch, records kept)
- [x] Ein Pane, den herdr nicht kennt (not_found bzw. pane_not_found), wird wie bisher als geschlossen behandelt; Test
  proof: TestReapReadsOnlyNotFoundAsAClosedTab cases 'pane unknown to hc', 'pane unknown to the old herdr', 'tab unknown to hc'
- [x] Release-Note in core/release/notes/0TREV1.md
  proof: core/release/notes/0TREV1.md
- [ ] Gleicher Fehler in close-tab.sh:159 (Geschwister): ein stummes herdr hebt den Park nicht mehr auf; Test

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-11-26 07:54 · Sascha** — parent session: uds:/run/user/1000/cc-socks/62036.sock
- **2026-11-26 07:56 · Sascha** — decision: keine offenen Entscheidungen. Unbekannter Fehlschlag von pane get -> 'kept <t>: herdr gave no answer for pane <p>' und die ganze Runde fuer das Ticket bleibt stehen (wie 'closing failed'). Erkennung per grep auf not_found im kombinierten stdout+stderr: deckt hc (kind not_found) und altes herdr (pane_not_found) ab, ohne den Stream des alten herdr zu kennen. close-tab.sh:159 bleibt draussen (entfernt nichts).
