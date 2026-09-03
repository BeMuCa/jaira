---
id: 01M1KZ56JYZWZAYTHBGSKA9CFA
title: Mehrzeilige Ticket-Felder behalten ihre Zeilen in der Detailansicht
status: signoff
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Ein Feld, das mit Zeilenumbruechen geschrieben wurde (z.B. review-check als nummerierte Schritte), rendert im TUI-Detail als Liste - eine Zeile je Schritt - statt als flachgedrueckter Prosa-Absatz"
context: "Berk am 03.09.: review-check soll durchnummerierte Schritte als Liste zeigen, nicht als Prosa. Befund: internal/tui/view.go detailBody/row() jagt jeden Feldwert durch wrap(), und wrap bricht nur an Leerzeichen - vorhandene \\n im Feld werden geplaettet. Die CLI (jaira show) erhaelt Umbrueche dagegen (SGPDYKs question-Feld rendert dort als Liste). Fix-Ort: row() bzw. der declared()-Pfad fuer die Prosa-Felder - je Eingabezeile wrappen (haengender Einzug 13 wie heute) und die Zeilen erhalten; wrapLines (view.go, seit HREQJR) ist das Muster, hat aber keinen Einzug-Parameter. Achtung styleLines-Lektion: mehrzeilige Bloecke nie mit Praefix durch Style.Render. Zweite Haelfte ist Schreibdisziplin: check-Felder ab jetzt mit echten Umbruechen schreiben (ein Schritt je Zeile) - der Renderer-Fix macht das sichtbar."
definition-of-done: Ein review-check mit einem Schritt je Zeile zeigt im TUI-Detail eine Zeile je Schritt (haengender Einzug unter der Label-Spalte); bestehende einzeilige Felder rendern unveraendert; ein Test deckt ein mehrzeiliges Feld ab
tags: []
blocked-by: []
commits: []
created-at: 2026-09-03T15:45:47Z
updated-at: 2026-09-03T16:14:58Z
claimed-by: EE-3NX6GL3-2629914
claimed-at: 2026-09-03T15:49:13Z
updated-by: BeMuCa
outcome-what: "wrapField-Helper: Feldwerte behalten die Zeilen des Autors (je Eingabezeile gewrappt, haengender Einzug 13); Detail-Pane row() und Signoff section() nutzen ihn - ein mehrzeiliger review-check rendert als nummerierte Liste"
outcome-why: "Berk am 03.09.: check soll durchnummerierte Schritte als Liste zeigen, nicht als Prosa; row() jagte jeden Wert durch wrap(), das \\n plaettet"
outcome-resolves: TestAMultilineCheckKeepsItsLines prueft beide Screens (Zeile beginnt unter der Label-Spalte); einzeilige Werte gehen durch den Passthrough-Zweig unveraendert (alte reviewcheck-Tests gruen); go test ./... -race RC=0
executed-by: fable
review-summary: "Ein Helper statt zweier Sonderfaelle; sitzt in view.go neben den anderen Umbruch-Helfern (wrap/wrapLines seit HREQJR), Signoff teilt ihn - die zwei Screens bleiben 'ein Werkzeug' (Kommentar in signoff.go verlangt das explizit). Kein Feld-Sonderfall: JEDES mehrzeilig geschriebene Feld profitiert (question, gaps), einzeilige bleiben byte-gleich."
review-gaps: "Nichts entfernt. Gelassen: Leerzeilen im Feldwert werden zu reinen Einzug-Zeilen (kosmetisch); wrapLines nicht wiederverwendet - es hat keinen Einzug-Parameter und sein Passthrough-Kontrakt (byte-identische Zeilen) ist ein anderer Job."
review-check: |-
  1. Neu bauen, ein Ticket mit mehrzeiligem review-check oeffnen (NJPQWE in signoff hat einen): jeder Schritt auf eigener Zeile, buendig unter der Label-Spalte.
  2. Dasselbe auf dem Signoff-Screen des Tickets.
  3. Ein einzeiliges Feld (z.B. assignee, goal) sieht aus wie vorher.
review-verdict: "accept (Zweitmodell: Sonnet-Review am selbst gelesenen Diff, eigene Laeufe; Implementierung Fable). Beide Funde umgesetzt: (F1) Commit-Buendelung mit NJPQWE-Hunks vor dem Push re-splittet (52732e1 + 8e74d58, Code-Baum byte-identisch zum reviewten Stand, diff=0); (F2) fieldRow als dritter Label-Renderer plaettete weiter -> wrapField + TestFieldRowKeepsTheAuthorsLines. Suite -race nach Cache-Loeschung RC=0 (15 Pakete) nach jedem Schritt. String-Level-Verifikation wie die Repo-Tests selbst; der Blick ins echte Terminal ist deine Signoff-Abnahme."
---

# Mehrzeilige Ticket-Felder behalten ihre Zeilen in der Detailansicht

## Definition of Done

- [x] Ein review-check mit einem Schritt je Zeile zeigt im TUI-Detail eine Zeile je Schritt (haengender Einzug unter der Label-Spalte); bestehende einzeilige Felder rendern unveraendert; ein Test deckt ein mehrzeiliges Feld ab
  proof: internal/tui/view.go wrapField (Zeilen erhalten, haengender Einzug 13) + row() nutzt ihn; internal/tui/signoff.go section() ebenso; TestAMultilineCheckKeepsItsLines deckt Detail UND Signoff; einzeilige Felder: wrap-Passthrough, bestehende reviewcheck-Tests unveraendert gruen; go test ./... -race RC=0

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-03 15:57 · BeMuCa** — Entscheidung: EIN Helper (wrapField) fuer beide Screens statt zweier Sonderfaelle; einzeilige Werte gehen unveraendert durch wrap (Passthrough-Zweig), also keine Aenderung an bestehenden Renderings - die alten reviewcheck-Tests pinnen das. Leerzeilen im Feld werden zu Einzug-Zeilen (kosmetisch, bewusst nicht behandelt). CLI (jaira show) erhielt Umbrueche schon immer - der Fix betrifft nur die TUI-Label-Spalten-Renderer.
- **2026-09-03 16:11 · BeMuCa** — Zweitmodell-Review (Sonnet) fand zwei Dinge: (F1) im KA9CFA-Commit ritten die uncommitteten NJPQWE-Runde-4-Hunks mit -> Re-Split vor dem Push (52732e1 + 8e74d58, Baum byte-identisch, diff=0). (F2) fieldRow (view.go:864) ist ein DRITTER Label-Renderer - Lane-declared-Felder auf Detail- UND Signoff-Screen - und plaettete weiter; gleicher wrapField-Fix plus TestFieldRowKeepsTheAuthorsLines. Merke: die Label-Spalte hat drei Renderer (row, section, fieldRow) - wer einen anfasst, prueft alle drei.
