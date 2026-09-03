---
id: 01M1KN2HSJ0B32MQ2532NJPQWE
title: "Jede Karte traegt eine Box, auch ohne Tag"
status: critique
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Auf dem Board ist jede Karte umrandet: mit Registry-Farbe des ersten Tags, ohne Tag in neutraler Rahmenfarbe - keine ungeboxten Karten mehr"
context: "Berk am 03.09. mit Screenshot der Backlog-Spalte: tag-lose Karten (F7K9MF, 8WDYZW, KMNM6Z) stehen ohne Box da - 'auch die ohne tags sollen von einer box umgeben sein!'. Revidiert bewusst die AXTFG3-Ausnahme 'nur Karten mit Registry-Farbe bekommen die Box' (damals als Hoehen-Schoner gedacht); die 2 Zeilen je Karte hat Berk explizit akzeptiert ('das mit den 2 zeilen ist ok, ich will den rahmen'). Umsetzung: internal/tui/view.go - renderCardBlock boxt immer (Fallback-Rahmenfarbe neutral, z.B. colFaint wie der Spaltenrahmen), cardHeight wird fuer alle Karten 5, cardsInBudget rechnet damit; tagbox-/fitwindow-Tests pinnen das alte Verhalten und muessen beide Kartenarten neu decken."
definition-of-done: "Jede Karte ist umrandet (Tag-Farbe wenn vorhanden, sonst neutral); das Kartenbudget rechnet mit der einheitlichen Hoehe; Tests decken geboxte Karten mit und ohne Tag; go test ./... -race gruen"
tags: []
blocked-by: []
commits: []
created-at: 2026-09-03T12:49:35Z
updated-at: 2026-09-03T13:25:08Z
claimed-by: EE-3NX6GL3-2382606
claimed-at: 2026-09-03T12:49:52Z
updated-by: BeMuCa
outcome-what: "Breiten-Fix aus dem Review: renderCard budgetiert auf die Inhaltsflaeche der Box (inner-2, lipgloss Width = Gesamtbreite) - Karten sind wieder exakt 5 Zeilen, Spalten zaehlen ehrlich (11/11 statt 11/10 auf 150x32, eine Karte MEHR auf 150x40); zwei Titel-Pins auf stabilen Praefix umgestellt"
outcome-why: "Das Opus-Review mass den Wrap kausal: Inhaltszeilen waren 2 Spalten breiter als die Box-Innenflaeche - deckellose Boxen unter Hoehe 13, versteckte Karten hinter '+0 more'; betraf auch die alten AXTFG3-Farbboxen"
outcome-resolves: "Reviewer-Fix 1:1 uebernommen (von ihm im Overlay vorab verifiziert: alle seine Proben gruen); go test ./... -race RC=0 nach Cache-Loeschung, Binary neu gebaut"
executed-by: fable
review-summary: "Kritik am Diff: (1) einfachste Form - kein neuer Renderer, renderCardBlock verliert nur den Nicht-Box-Zweig, cardHeight wird konstant; (2) Muster gehalten - Fallback-Farbe ist colFaint, dieselbe wie der Spaltenrahmen (view.go:361), keine neue Farbe erfunden; (3) nichts Spekulatives; (4) richtiger Ort - der Overflow-Fix sitzt in renderCard, wo die Zeile gebaut wird, nicht als Nach-Klammerung. Preis, benannt: alle Karten sind innen 2 Spalten schmaler und je 2 Zeilen hoeher - weniger Karten pro Spalte sichtbar; das war Berks explizite Abwaegung ('das ist ok, ich will den rahmen')."
review-gaps: "Nichts entfernt. Gelassen: cardHeight behaelt Receiver+Parameter (Signatur-Kompatibilitaet aller Aufrufer, und die Hoehe kann wieder dynamisch werden); der '✎ someo'-Pin wurde zum Glyph-Pin gelockert statt die Boardbreite im Test zu schrauben - Begruendung in der Ticket-Note."
review-check: "1. Neu bauen, Board oeffnen: JEDE Karte hat einen Rahmen - tag-lose in der matten Spaltenrahmen-Farbe, getaggte (mit Registry-Farbe) farbig. Dein Screenshot-Fall (Backlog, F7K9MF etc.) ist der Testfall. 2. Karte mit vielen Flags (blocked + unworked + fremder ✎) ansehen: genau 3 Inhaltszeilen in der Box, kein Umbruch, unterer Rand intakt. 3. Spalten scrollen: '+N more' rechnet mit der neuen Hoehe."
---

# Jede Karte traegt eine Box, auch ohne Tag

## Definition of Done

- [x] Jede Karte ist umrandet (Tag-Farbe wenn vorhanden, sonst neutral); das Kartenbudget rechnet mit der einheitlichen Hoehe; Tests decken geboxte Karten mit und ohne Tag; go test ./... -race gruen
  proof: view.go: renderCardBlock boxt immer (Fallback colFaint), cardHeight konstant 5; Tests: TestUntaggedCardIsBoxedNeutrally, TestTaggedCardWithoutColourIsBoxedNeutrally, TestCardHeightIsUniform..., Budget-Tests neu gerechnet; go test ./... -race RC=0

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-09-03 12:59 · BeMuCa** — Entscheidungen: (1) Neutralrahmen = colFaint (240), dieselbe Farbe wie der Spaltenrahmen - Karten ohne Tag-Farbe treten zurueck, faerben aber nicht falsch. (2) Beifang-Fund beim Umbau: renderCard schrieb meta/flags als '  '+truncate(...,w) - bis zu w+2 breit. Ungeboxed hat clampBlock das verschluckt; IN der Box brach die Zeile um und zerstoerte die 3-Zeilen-Hoehe (Karten mit vielen Flags wurden 4+ Zeilen, cardHeight log). Fix: Budget w-2 fuer die eingerueckten Zeilen - betraf auch die schon existierenden AXTFG3-Farbboxen. (3) updatedby-Pin von '✎ someo' auf das Glyph '✎' gelockert: jede Karte ist jetzt 2 Spalten schmaler innen, wieviel Name ueberlebt haengt von der Boardbreite ab; der Testkommentar sagte selbst 'the glyph, not the whole name'.
- **2026-09-03 13:19 · BeMuCa** — Review-Runde 1 (Opus, kausal gemessen): renderCardBlock gibt renderCard 'inner' als Budget, aber lipgloss Width(inner) ist die GESAMTbreite der Box - Inhaltsflaeche ist inner-2 (gemessen 10->8, 16->14, 28->26; stand sogar schon in den lipgloss-Learnings). Folge: Inhaltszeilen 2 zu breit, Karten rendern 6 statt 5 Zeilen, Spalten melden '+0 more' und verstecken Karten, Terminalhoehe <=12 zeigt deckellose Boxen. Betraf auch die AXTFG3-Farbboxen schon - universelles Boxen machte es sichtbar. Fix (vom Reviewer im Overlay verifiziert, alles gruen): renderCard bekommt max(1, inner-2); Preis: 2 Titel-Truncation-Repins (TestBoardRenders render_test.go:99, TestFilterNarrowsTheBoard :124). Nach Fix: 11/11 Karten auf 150x32 statt 11/10, auf 150x40 passt eine Karte MEHR, weil Hoehen endlich stimmen.
- **2026-09-03 13:25 · BeMuCa** — Fix umgesetzt wie vom Reviewer verifiziert: renderCard bekommt max(1, inner-2) - die Inhaltsflaeche der Box, nicht ihre Gesamtbreite; Kommentar an Ort und Stelle erklaert die lipgloss-Width-Semantik. Die zwei vorhergesagten Titel-Repins gemacht (render_test.go: 'Refactor auth' -> Praefix 'Refactor', uebersteht kuenftige Breiten-Feinjustagen). Voller Lauf -race nach Cache-Loeschung: RC=0.
