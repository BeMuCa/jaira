---
id: 01M1E5VPBM0P9VESJ6T7D9SC5A
title: "Ein gefuelltes Lane-Feld nennt die Lane, aus der sein Wert stammt"
status: done
ready: true
creator: BeMuCa
goal: "Im Detail- und Signoff-View traegt jedes gefuellte Feld, das eine Lane deklariert, ein dezentes Herkunfts-Label - z.B. summary (critique) Text - so wie leere Felder schon '— owed by <lane>' sagen"
context: "Berk am 01.09., am Signoff von 8HV58V auf seinem req-Board: die owed-Zeilen nennen die schuldende Lane, aber sobald ein Feld gefuellt ist, verschwindet die Herkunft - man sieht nicht mehr, dass summary aus der critique-Lane stammt. Gewuenschtes Format sinngemaess '(critique-lane)Text' vor dem Wert. Traeger der Herkunft ist heute nur die Deklaration (erste Lane in Board-Reihenfolge, die das Feld in output-produces nennt - dieselbe Zuordnung wie gate.OwedBy); WER zuletzt schrieb, weiss erst die History aus Schema-Cut 2. Deklarieren zwei Lanes dasselbe Feld (critique und review beide review-summary), ist das Label also der Schema-Besitzer, nicht der letzte Schreiber - Grenze dokumentieren. Basis-Felder (goal, context) bekommen KEIN Label: sie gehoeren der Erstellung (Entscheidung 8), auch wenn brainstorm goal produziert. Render-Stellen: detailBody declared()/laneFields() und renderSignOff declared() in internal/tui."
definition-of-done: "Gefuellte Outcome-/Review-/Lane-Felder zeigen im Detail- und Signoff-View ein dezentes (lane)-Praefix; leere weiter '— owed by'; goal/context ohne Label; Breiten-Tests bleiben gruen; go test ./... -race gruen"
blocked-by: []
commits: []
created-at: 2026-09-01T09:47:29Z
updated-at: 2026-09-03T10:28:58Z
claimed-by: EE-3NX6GL3-3294065
claimed-at: 2026-09-01T09:48:34Z
updated-by: BeMuCa
assignee: BeMuCa
outcome-what: "Runde 2: DeclaredBy liefert ALLE Deklarierer (board order), sourced joint '(critique/review) '; Leerzeile rettet OwedBys godoc; Geschwister-Assertions verschaerft"
outcome-why: "Review-Funde: godoc-Umzuordnung und ein Label, das auf Boards mit critique+optimize auf 3 von 4 Review-Feldern die falsche Lane behauptete"
outcome-resolves: "go test ./internal/tui ./core/gate -count=1 EXIT=0; gate-Test pinnt den Doppel-Deklarierer-Fall"
review-summary: none
review-gaps: "Nichts entfernt. Gelassen: Label als Klartext im gewrappten Wert (kein Substring-Styling ueber Umbrueche); goal-Ausnahme statt einer Basis-Feld-Liste (goal ist das einzige Basis-Feld, das eine Lane deklariert - context/dod stehen in keiner produces)."
review-check: "Terminal ~100 Spalten, ein Ticket in review mit gefuelltem review-summary oeffnen: die Zeile liest 'summary (critique|review) <text>' - je nachdem, welche Lane das Feld auf DIESEM Board zuerst deklariert; goal-Zeile ohne Klammer. Signoff-Screen: gleiche Labels in den Sektionen. Leere Felder unveraendert '— owed by <lane>'."
review-verdict: "accept nach 3 Runden (Opus-Reviews): Label joint alle Deklarierer '(critique/review)' - wahr in jeder Board-Konstellation, im Renderer testgepinnt; goal bleibt begruendet bar (Reviewer stimmt zu: bedingtes Label wuerde behaupten, was das File nicht traegt); OwedBys godoc druckt wieder (Commit auf die go-doc-Ausgabe gegated)"
---

# Ein gefuelltes Lane-Feld nennt die Lane, aus der sein Wert stammt

## Definition of Done

- [ ] Gefuellte Outcome-/Review-/Lane-Felder zeigen im Detail- und Signoff-View ein dezentes (lane)-Praefix; leere weiter '— owed by'; goal/context ohne Label; Breiten-Tests bleiben gruen; go test ./... -race gruen

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-01 09:57 · BeMuCa** — Grenze bewusst dokumentiert (auch im Code-Kommentar): das Label ist der Schema-Besitzer (erste deklarierende Lane in Board-Reihenfolge, dieselbe Zuordnung wie OwedBy), NICHT der letzte Schreiber - wer zuletzt schrieb, weiss erst die History aus Schema-Cut 2. Deklarieren critique und review beide review-summary, steht (critique) auch vor einem Wert, den review ueberschrieb. goal traegt kein Label (Basis-Feld, Entscheidung 8). Ein Bestandstest pinnte die unetikettierte Zeile und wurde auf die etikettierte verschaerft (eigener Commit).
- **2026-09-01 10:04 · BeMuCa** — review reject, uebernommen: (L3) DeclaredBy klebte ohne Leerzeile an OwedBys Doku - godoc ordnete sie um; (L4) auf Boards mit critique+optimize adoptiert luegt das Label (erste Deklariererin statt Schreiberin) auf 3 von 4 Review-Feldern - Fix: ALLE Deklarierer joinen '(critique/review) ...', wahr in jeder Konstellation; (L6) zwei Geschwister-Assertions auf die etikettierte Form verschaerfen. Abgelehnt mit Begruendung: (L5) goal bleibt bar - auch on-route kann der Ersteller es geschrieben haben, ein bedingtes Label waere wieder eine Behauptung ohne Beleg; echte Schreiber-Herkunft ist Schema-Cut 2.
