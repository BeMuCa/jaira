---
id: 01M1KZ56JYZWZAYTHBGSKA9CFA
title: Mehrzeilige Ticket-Felder behalten ihre Zeilen in der Detailansicht
status: backlog
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
updated-at: 2026-09-03T15:45:47Z
---

# Mehrzeilige Ticket-Felder behalten ihre Zeilen in der Detailansicht

## Definition of Done

- [ ] Ein review-check mit einem Schritt je Zeile zeigt im TUI-Detail eine Zeile je Schritt (haengender Einzug unter der Label-Spalte); bestehende einzeilige Felder rendern unveraendert; ein Test deckt ein mehrzeiliges Feld ab

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress

