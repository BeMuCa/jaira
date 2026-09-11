---
id: 01M28MHSDBABYVD8785A74VM40
title: "Das Logbuch wird abgelegt, wenn ein Mensch es sagt, nicht wenn ein Ticket fertig wird"
status: in-progress
ready: true
creator: Alexander Sacharov
goal: "Fertige Tickets sammeln sich in done, und wer seine Stunden eintraegt, legt sie mit einem Befehl als Tagesordner ab - das Board sagt Bescheid, wenn sich viel angesammelt hat, entscheidet aber nichts"
context: |-
  Issue #6 im Original, geschrieben aus dem taeglichen Gebrauch: ein Move nach done nahm 49 fremde fertige Tickets mit ins Logbuch. Der Commit hiess nach einem Handle und enthielt 52 geaenderte Dateien, davon 49 Umbenennungen fremder Arbeit. Rueckgaengig machen ging nur ganz oder gar nicht.

  Im Code steht es so: core/lane/settle.go ruft bei logbook-on-entry s.FileLane(lane), und core/ticket/trim.go nimmt JEDEN Ticket mit t.Status == lane - nicht den, der gerade eingetreten ist.

  Der Kern der Beschwerde ist nicht der Sweep selbst, sondern WANN er passiert: 'dieses Ticket ist fertig' ist eine Aussage ueber die Arbeit, 'diese Tickets sind abgelegt' eine ueber die Buchhaltung. Die zweite faellt Tage spaeter und betrifft einen Satz, den jemand bewusst zusammenstellt - meist beim Eintragen der Stunden. Eine terminale Lane ist kein Beleg dafuer, dass jemand ablegen will.

  Entschieden von Alexander (11.09.2026): der Sweep bleibt, aber er wird von Hand ausgeloest. Tickets sammeln sich in done. Wer den Schnitt macht, ruft das Logbuch auf, und es legt den Tagesordner mit allem an, was fertig ist. Das TUI darf darauf hinweisen, dass sich viel angesammelt hat - entscheiden darf es nichts.

  Damit wird auch WXQ9PT abgeloest, das in signoff auf Abnahme wartet: dessen DoD verlangt 'done ist danach leer'.
definition-of-done: "ein Move nach done legt nichts mehr ab: das Ticket bleibt in done stehen, und die Meldung sagt, dass es fertig ist, nicht dass es abgelegt wurde; 'jaira logbook --all' legt alles aus der terminalen Lane in den heutigen Ordner und nennt, was es abgelegt hat; 'jaira logbook <id>' legt weiterhin genau eines ab und 'jaira logbook' ohne Argument listet nur; das Board zeigt eine Zeile, sobald sich in der terminalen Lane mehr als eine einstellbare Zahl angesammelt hat, und nennt den Befehl - es legt nie selbst ab; Tests decken ab, dass ein Move nichts mitnimmt, dass --all den ganzen Satz nimmt, und dass die Zeile ab der Schwelle erscheint"
tags:
  - cli
blocked-by: []
commits: []
created-at: 2026-09-11T16:24:28Z
updated-at: 2026-09-11T16:24:49Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-630371
claimed-at: 2026-09-11T16:24:41Z
---

# Das Logbuch wird abgelegt, wenn ein Mensch es sagt, nicht wenn ein Ticket fertig wird

## Definition of Done

- [ ] ein Move nach done legt nichts mehr ab: das Ticket bleibt in done stehen, und die Meldung sagt, dass es fertig ist, nicht dass es abgelegt wurde; 'jaira logbook --all' legt alles aus der terminalen Lane in den heutigen Ordner und nennt, was es abgelegt hat; 'jaira logbook <id>' legt weiterhin genau eines ab und 'jaira logbook' ohne Argument listet nur; das Board zeigt eine Zeile, sobald sich in der terminalen Lane mehr als eine einstellbare Zahl angesammelt hat, und nennt den Befehl - es legt nie selbst ab; Tests decken ab, dass ein Move nichts mitnimmt, dass --all den ganzen Satz nimmt, und dass die Zeile ab der Schwelle erscheint

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

