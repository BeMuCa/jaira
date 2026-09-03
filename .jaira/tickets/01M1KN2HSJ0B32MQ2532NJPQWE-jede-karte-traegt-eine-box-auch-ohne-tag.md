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
updated-at: 2026-09-03T18:49:35Z
claimed-by: EE-3NX6GL3-2382606
claimed-at: 2026-09-03T12:49:52Z
updated-by: BeMuCa
outcome-what: "Gestapelte Karten teilen sich eine Border-Reihe: ab der zweiten Karte im Fenster faellt deren Top-Border weg (renderColumn), cardsInBudget rechnet gestapelte Karten mit 4 statt 5 Zeilen - der optische Abstand zwischen Boxen ist halbiert und pro weiterer Karte wird eine Zeile frei"
outcome-why: "Berks 4. Screenshot: Luecken zwischen den Boxen halbieren; Messung zeigte keine Leerzeile, sondern zwei aneinanderstossende Border-Reihen, deren Glyphen nur die halbe Zelle fuellen"
outcome-resolves: "TestStackedCardsShareOneBorderRow verbietet Bottom-ueber-Top-Border boardweit; der Box-Balance-Waechter zaehlt jetzt exakt (1 Stack-Top + 1 Bottom je Karte + Spaltenrahmen, positionsgenau gegen cardsInBudget); Budget-Tests neu gerechnet (5+4+4...); go test ./... -race RC=0"
executed-by: fable
review-summary: "Runde 4 (Berks 3. Screenshot): Selektions-Balken entfernt - Selektion ist die gefaerbte Titel-Schrift (Tag-Farbe der Karte, sonst Akzent; Bold bleibt als zweiter Kanal, Zustand haengt nie an Farbe allein); Innen-Einzug der Box von 2 auf 1 Spalte, Budgets folgen (w-1). Der Board-Umschalter behaelt seinen Balken - andere Liste, keine Box. Kein einfacherer Schnitt: beide Aenderungen leben in renderCard, wo die Zeilen entstehen."
review-gaps: "Nichts entfernt. Gelassen: cardHeight behaelt Receiver+Parameter (Signatur-Kompatibilitaet aller Aufrufer, und die Hoehe kann wieder dynamisch werden); der '✎ someo'-Pin wurde zum Glyph-Pin gelockert statt die Boardbreite im Test zu schrauben - Begruendung in der Ticket-Note."
review-check: |-
  1. Neu bauen, Board oeffnen: jede Karte gerahmt (tag-los matt, getaggt farbig), exakt 5 Zeilen.
  2. Abstand: Box endet EINE Spalte vor dem rechten Spaltenrand; innen EIN Leerzeichen Einzug (dein 3. Screenshot).
  3. Selektierte Karte: kein blauer Balken mehr - der Titel selbst ist gefaerbt (Tag-Farbe, sonst blau) und fett.
  4. Terminal auf ~12 Zeilen stauchen: Deckel und Boden der Boxen bleiben.
  5. Spalte mit Ueberlauf: die +N-more-Zahl stimmt mit dem Fehlenden ueberein.
review-verdict: "accept - Runden 3+4 koordinator-verifiziert (Opus-Reviewer seit zwei 529s nicht erneut bemueht; Deltas: eine Konstante in renderColumn, dann Balken->Titel-Farbe und Einzug 2->1 in renderCard). Die Box-Arithmetik aus Runde 2 (kausal gemessen, TestACardHeavyWithFlagsStaysFiveRows) budgetiert relativ und blieb unangetastet; Suite -race nach Cache-Loeschung RC=0 (15 Pakete) nach JEDER Runde. Zustand haengt weiter nie an Farbe allein (Bold bleibt auf dem selektierten Titel). Optik - Einzug, Abstand rechts, Selektion ohne Balken - ist deine Signoff-Abnahme; review-check zaehlt die Schritte."
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
- **2026-09-03 14:36 · BeMuCa** — Berks Signoff-Feedback (03.09., Screenshot): Abstaende kleiner - rechts von den Boxen kann es fast bis zum Rand. Heute: Box flush links, 4 Spalten Luft rechts (renderColumn gibt w-4, Box-Gesamt = w-6, Spalteninnen = w-2). Ziel: 1 Spalte Luft rechts => renderCardBlock bekommt w-1 (Box w-3); Inhalt gewinnt 3 Spalten fuer Titel/Flags.
- **2026-09-03 14:41 · BeMuCa** — Umgesetzt: renderColumn gibt den Karten w-1 statt w-4 - Box endet eine Spalte vor dem rechten Spaltenrand (vorher vier), flush links wie gehabt; Titel/Handle/Flags gewinnen 3 Spalten. Kommentar im Code nennt den Grund (Berks Screenshot). Volle Suite -race nach Cache-Loeschung RC=0, 15 Pakete; kein Test pinnt den exakten Abstand - die Optik ist Berks Abnahme im Signoff.
- **2026-09-03 15:38 · BeMuCa** — Berks 3. Signoff-Feedback (Screenshot): (1) Innenabstand links in der Box zu gross; (2) der blaue Selektions-Balken soll weg - Selektion reicht als gefaerbte Schrift, blau bzw. in der Tag-Farbe der Karte; (3) seine Frage 'Text vergroessert?' - nein, Terminalschrift ist fix; Karten sind seit w-1 drei Spalten breiter, daher mehr sichtbarer Titel.
- **2026-09-03 16:11 · BeMuCa** — Nachtrag: die Runde-4-Hunks (Balken raus, Einzug 1) waren nie eigenstaendig committet und ritten im KA9CFA-Commit mit - vom Zweitmodell-Review gefunden (bekannte Buendelungs-Klasse). Re-Split vor dem Push: eigener Commit 52732e1, Code-Baum danach byte-identisch zum reviewten Stand (diff=0 verifiziert).
- **2026-09-03 18:37 · BeMuCa** — Berks 4. Feedback: Luecken ZWISCHEN den Boxen halbieren. Messung (TestBoardRenders-Log): es gibt KEINE Leerzeile - der Eindruck entsteht durch zwei aneinanderstossende Border-Reihen (Glyphen zeichnen halbe Zellhoehe). Halbieren = gestapelte Karten teilen sich EINE Border-Reihe: ab der zweiten Karte im Fenster faellt die eigene Top-Border weg, die Bottom-Border der vorigen ist der Trenner (traegt deren Tag-Farbe). Hoehen: erste Karte 5, jede weitere 4 - cardsInBudget rechnet positionsabhaengig.
