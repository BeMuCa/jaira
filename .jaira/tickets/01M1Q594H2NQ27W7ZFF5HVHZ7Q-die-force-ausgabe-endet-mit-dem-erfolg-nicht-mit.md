---
id: 01M1Q594H2NQ27W7ZFF5HVHZ7Q
title: "Die force-Ausgabe endet mit dem Erfolg, nicht mit der Verweigerungsliste"
status: critique
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Bei jaira move --force stehen die Override-Bullets VOR der Erfolgszeile: die letzte Zeile eines gelungenen Moves ist immer der Move selbst (bzw. nextStep/filed-Zeilen), nie ein Bullet, das wie ein Fehler aussieht"
context: "Berk am 04.09. (Feedback L121), Vorgeschichte Invariante 11: 'move --force | tail -1' zeigte die Bullet-Zeile des Erfolgsberichts und las sich wie eine Verweigerung - der Fehlalarm kostete einen Tag (BDV0HM, 31.08.). Fix: Reihenfolge der Ausgabe drehen (Overrode-Block zuerst, dann 'X -> lane', dann filed/trimmed, dann nextStep) - Inhalt unveraendert, nur Ordnung; --json unberuehrt."
definition-of-done: Nach einem forcierten Move ist die letzte Ausgabezeile nie ein Override-Bullet; ein Test pinnt die Reihenfolge; bestehende Ausgabe-Tests angepasst
tags: []
blocked-by: []
commits: []
created-at: 2026-09-04T21:30:31Z
updated-at: 2026-09-05T02:24:14Z
claimed-by: EE-3NX6GL3-395982
claimed-at: 2026-09-05T02:24:08Z
updated-by: BeMuCa
outcome-what: "CLI-Ausgabe eines forcierten Moves umgestellt: Override-Bullets zuerst, Erfolgszeile 'X -> lane' danach - die letzte Zeile eines gelungenen Moves ist nie mehr ein Bullet; Kommentar an der Stelle nennt den BDV0HM-Fehlalarm"
outcome-why: "Feedback L121 / Invariante-11-Falle: 'move --force | tail -1' zeigte ein Verweigerungs-Bullet des Erfolgsberichts und kostete am 31.08. einen Tag Umweg"
outcome-resolves: TestForcedMoveEndsOnItsSuccessLine erzwingt einen echten Override und pinnt die Schlusszeile; --json unberuehrt (kehrt vorher zurueck); voller Lauf -race RC=0
executed-by: fable
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

