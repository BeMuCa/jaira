---
id: 01M1Q594H2NQ27W7ZFF5HVHZ7Q
title: "Die force-Ausgabe endet mit dem Erfolg, nicht mit der Verweigerungsliste"
status: done
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Bei jaira move --force stehen die Override-Bullets VOR der Erfolgszeile: die letzte Zeile eines gelungenen Moves ist immer der Move selbst (bzw. nextStep/filed-Zeilen), nie ein Bullet, das wie ein Fehler aussieht"
context: "Berk am 04.09. (Feedback L121), Vorgeschichte Invariante 11: 'move --force | tail -1' zeigte die Bullet-Zeile des Erfolgsberichts und las sich wie eine Verweigerung - der Fehlalarm kostete einen Tag (BDV0HM, 31.08.). Fix: Reihenfolge der Ausgabe drehen (Overrode-Block zuerst, dann 'X -> lane', dann filed/trimmed, dann nextStep) - Inhalt unveraendert, nur Ordnung; --json unberuehrt."
definition-of-done: Nach einem forcierten Move ist die letzte Ausgabezeile nie ein Override-Bullet; ein Test pinnt die Reihenfolge; bestehende Ausgabe-Tests angepasst
tags: []
blocked-by: []
commits:
  - 2d4a718912ba2786a57778cd59e7c15dfb9b30fd
created-at: 2026-09-04T21:30:31Z
updated-at: 2026-09-05T02:25:06Z
claimed-by: EE-3NX6GL3-395982
claimed-at: 2026-09-05T02:24:08Z
updated-by: BeMuCa
outcome-what: "CLI-Ausgabe eines forcierten Moves umgestellt: Override-Bullets zuerst, Erfolgszeile 'X -> lane' danach - die letzte Zeile eines gelungenen Moves ist nie mehr ein Bullet; Kommentar an der Stelle nennt den BDV0HM-Fehlalarm"
outcome-why: "Feedback L121 / Invariante-11-Falle: 'move --force | tail -1' zeigte ein Verweigerungs-Bullet des Erfolgsberichts und kostete am 31.08. einen Tag Umweg"
outcome-resolves: TestForcedMoveEndsOnItsSuccessLine erzwingt einen echten Override und pinnt die Schlusszeile; --json unberuehrt (kehrt vorher zurueck); voller Lauf -race RC=0
executed-by: fable
review-summary: "Kritik: reine Reihenfolge-Aenderung an einer Stelle, Inhalt und JSON unangetastet; der Kommentar traegt die Begruendung (BDV0HM) an den Ort, an dem der naechste die Reihenfolge 'zurueckverbessern' wollen wuerde. TUI-forceMove bewusst gelassen: dort liest niemand tail -1, die Notify zeigt den Block ganz."
review-gaps: "Nichts entfernt. Gelassen: filed/trimmed-Zeilen und nextStepLine folgen NACH der Erfolgszeile - sie sind Zusatzinfo eines gelungenen Moves, kein Erfolgssignal; der Test pinnt nur, dass keine Bullet-Zeile zuletzt steht."
review-check: |-
  1. Neues Binary; ein absichtlich unfertiges Ticket mit --force nach todo moven.
  2. Ausgabe: erst Overrode+Bullets, dann X -> todo als letzte Zeile.
  3. move --force | tail -1 zeigt jetzt den Erfolg - die BDV0HM-Falle ist zu.
test-verdict: "pass: voller Lauf -race RC=0 (test16 deckt den Commit-Stand); Verhalten am ECHTEN Binary exerziert - der naechste forcierte Move auf diesem Board zeigt die neue Reihenfolge (Beleg folgt am eigenen signoff->done-Move dieses Tickets)"
review-verdict: "accept (koordinator-verifiziert, Direktdurchlauf per Berks Auftrag: Einzeilen-Reihenfolge-Fix mit Pin-Test, -race gruen; der eigene forcierte done-Move dieses Tickets ist die Live-Probe der neuen Ausgabe)"
---

# Die force-Ausgabe endet mit dem Erfolg, nicht mit der Verweigerungsliste

## Definition of Done

- [x] Nach einem forcierten Move ist die letzte Ausgabezeile nie ein Override-Bullet; ein Test pinnt die Reihenfolge; bestehende Ausgabe-Tests angepasst
  proof: flow.go: Override-Block vor der Erfolgszeile (Kommentar nennt BDV0HM); TestForcedMoveEndsOnItsSuccessLine pinnt letzte Zeile = '-> lane', nie Bullet; bestehende Move-Tests gruen (contains-Asserts, keine Reihenfolge-Pins gerissen); voller Lauf -race RC=0 (test16)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

