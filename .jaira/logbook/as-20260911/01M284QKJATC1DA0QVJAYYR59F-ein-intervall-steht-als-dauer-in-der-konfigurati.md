---
id: 01M284QKJATC1DA0QVJAYYR59F
title: "Ein Intervall steht als Dauer in der Konfiguration, und das Board haelt sich daran"
status: done
ready: true
creator: Alexander Sacharov
goal: "Wie oft im Hintergrund gefetcht wird, steht an einer Stelle und gilt fuer CLI und Board gleichermassen - und laesst sich fuer eine Vorfuehrung auf Sekunden stellen"
context: |-
  Zwei Dinge, die beim Vorbereiten der Bildschirmaufnahme aufgefallen sind.

  1. Das TUI holt die Refs nach einem eigenen, fest eingebauten Takt: refFetchEvery = 60 * time.Second in internal/tui/refs.go. Die Einstellung 'fetch-every-minutes' aus settings.json, die der Hintergrund-Fetch der CLI benutzt (PP5SCQ), liest es nicht. Zwei Zahlen fuer dieselbe Frage, und wer die Einstellung aendert, wundert sich, dass das Board weiter im Minutentakt laeuft.

  2. Die Felder heissen 'fetch-every-minutes', 'snapshot-every-hours' und 'landing-grace-days' - drei verschiedene Einheiten im Namen, und keine davon erlaubt Sekunden. Fuer eine Aufnahme, in der zu sehen sein soll, wie ein neu angelegtes Ticket von selbst auf dem Board erscheint, braucht es drei Sekunden, nicht eine Minute.

  Noch nichts davon ist veroeffentlicht: die Felder existieren nur auf diesem Branch, also ist die Form jetzt noch frei zu waehlen und spaeter nicht mehr.
definition-of-done: "die drei Intervalle stehen als Dauer in settings.json ('fetch-every', 'snapshot-every', 'landing-grace', gelesen mit time.ParseDuration, also '3s', '10m' und '72h' gleichermassen moeglich); ein unlesbarer Wert faellt auf den Default zurueck statt das Board zu verweigern; das TUI benutzt dieselbe Einstellung wie die CLI statt einer eigenen Konstante; Tests decken Default, gueltigen Wert und Unsinn ab"
tags:
  - cli
blocked-by: []
commits:
  - 22abcfdf2699b5047de0c9e093150350c3140060
  - 7d3a24d8091a6b14c948931bf7cc0e4be57d28f9
created-at: 2026-09-11T11:48:02Z
updated-at: 2026-09-11T13:31:59Z
claimed-by: DESKTOP-RFTCH11-193636
claimed-at: 2026-09-11T11:48:26Z
updated-by: Alexander Sacharov
assignee: Alexander Sacharov
question: "Zwei Sachen: (1) die drei Intervalle heissen jetzt fetch-every, snapshot-every und landing-grace und sind Dauern - passt das, bevor es veroeffentlicht ist? (2) eine Zuweisung, die im Hintergrund ankommt, steht 30 Sekunden als Zeile auf dem Board statt als Bildschirm zum Wegdruecken - reicht dir das, oder soll sie laenger stehen?"
outcome-what: "drei Intervalle als Dauer in settings.json, das TUI liest dieselbe Einstellung wie die CLI, und eine Ankunft meldet sich als Zeile statt als Bildschirm"
outcome-why: "Das Board hatte eine eigene, fest eingebaute Minute und las die Einstellung nicht - und eine im Hintergrund ankommende Zuweisung nahm den Bildschirm weg, mitten in der Arbeit"
outcome-resolves: "Die DoD Punkt fuer Punkt: die drei Intervalle sind Dauern (time.ParseDuration, also 3s und 72h gleich gesagt), ein unlesbarer Wert faellt auf den Default statt das Board zu verweigern, das TUI benutzt settings.FetchInterval statt einer Konstante, und drei Tests decken Default, gueltige Dauer und Unsinn ab"
---

# Ein Intervall steht als Dauer in der Konfiguration, und das Board haelt sich daran

## Definition of Done

- [x] die drei Intervalle stehen als Dauer in settings.json ('fetch-every', 'snapshot-every', 'landing-grace', gelesen mit time.ParseDuration, also '3s', '10m' und '72h' gleichermassen moeglich); ein unlesbarer Wert faellt auf den Default zurueck statt das Board zu verweigern; das TUI benutzt dieselbe Einstellung wie die CLI statt einer eigenen Konstante; Tests decken Default, gueltigen Wert und Unsinn ab
  proof: core/settings: fetch-every/snapshot-every/landing-grace als Dauer mit time.ParseDuration, unlesbar faellt auf den Default; internal/tui/refs.go refFetchEvery() liest dieselbe Einstellung statt einer Konstante; 3 Tests (Defaults, Dauern, Unsinn); von Hand geprueft: Board zeigt '1 tickets' und wird ohne Tastendruck zu '2 tickets'

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-11 11:52 · Alexander Sacharov** — Beim Pruefen mit einem echten Board in einem tmux-Pane kam ein zweiter Fehler heraus, der nicht im Ticket stand: ein Ticket, das im Hintergrund ankommt, wurde ueber m.notify gemeldet - und das schaltet den Bildschirm auf die Nachrichtenansicht um, die man mit esc wegdruecken muss. Also: jemand arbeitet, und weil ihm jemand anders ein Ticket zuweist, verliert er den Blick auf sein Board und moeglicherweise den Tastendruck, den er gerade getippt hat.

Jetzt eine Zeile auf dem Board selbst (flash, 30 Sekunden, in der Hinweiszeile) statt eines eigenen Bildschirms. Die Trennung, die dabei entstanden ist und die es vorher nicht gab: die Nachrichtenansicht gehoert dem, was der Mensch gerade getan hat - eine Flash-Zeile dem, was passiert ist, waehrend er etwas anderes tat.

Gefunden nur, weil die Probe das Board wirklich geoeffnet und den Bildschirminhalt angesehen hat, statt zu pruefen, ob die Funktion aufgerufen wurde.
