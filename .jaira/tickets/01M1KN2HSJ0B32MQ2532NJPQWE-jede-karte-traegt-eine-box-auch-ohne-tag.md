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
updated-at: 2026-09-03T14:41:37Z
claimed-by: EE-3NX6GL3-2382606
claimed-at: 2026-09-03T12:49:52Z
updated-by: BeMuCa
outcome-what: "Kartenbreite in der Spalte: renderColumn reicht w-1 statt w-4 an renderCardBlock - rechts bleibt eine Spalte Luft statt vier, Inhalt gewinnt 3 Spalten"
outcome-why: "Berks Signoff-Feedback mit Screenshot: 'mach die abstaende kleiner, rechts von den boxen kann es fast bis zum rand'"
outcome-resolves: "Eine Zeile plus Kommentar; alle Breiten-Invarianten haengen an renderCardBlock/renderCard, die unveraendert budgetieren (Box-Gesamt = Parameter-2, Inhalt = -4) - Suite -race RC=0 bestaetigt; Optik selbst ist die Signoff-Abnahme"
executed-by: fable
review-summary: "Runde 2 nach dem Review-Fund. Kritik am Fix-Diff: eine Zeile Verhalten (renderCard erhaelt inner-2) + erklaerender Kommentar an der Stelle, an der der naechste Leser denselben Fehler machen wuerde - kein Workaround weiter aussen, keine doppelte Breitenrechnung. Die Titel-Repins pinnen einen Praefix statt des exakten Schnitts, damit der Test Feinjustagen der Budgets ueberlebt und trotzdem faellt, wenn der Titel ganz verschwindet."
review-gaps: "Nichts entfernt. Gelassen: cardHeight behaelt Receiver+Parameter (Signatur-Kompatibilitaet aller Aufrufer, und die Hoehe kann wieder dynamisch werden); der '✎ someo'-Pin wurde zum Glyph-Pin gelockert statt die Boardbreite im Test zu schrauben - Begruendung in der Ticket-Note."
review-check: "1. Neu bauen, Board oeffnen: JEDE Karte gerahmt (tag-los matt, getaggt farbig), exakt 5 Zeilen je Karte, kein abgeschnittener unterer Rand. 2. Terminal auf ~12 Zeilen stauchen: Boxen behalten ihren Deckel UND Boden (vorher deckellos). 3. Spalte mit vielen Karten: die '+N more'-Zahl stimmt mit dem ueberein, was fehlt (vorher '+0 more' bei versteckten Karten). 4. Auf 150x40 passt eine Karte mehr als vorher - Hoehen sind jetzt ehrlich."
review-verdict: "accept - mit offengelegtem Ausweich: der Opus-Reviewer hatte den Breitenfehler kausal gemessen und den Fix VORAB im Overlay verifiziert (alle seine Proben gruen: Kartenhoehe 5 auf w=16..44, Spalten 11/11 statt 11/10 auf 150x32, +N-more ehrlich, Deckel unter Hoehe 13); beim finalen Re-Check starb er zweimal am 529 Overloaded. Der Koordinator hat den Einzeiler wortgleich uebernommen (diff selbst gelesen), die volle Suite -race nach Cache-Loeschung gefahren (RC=0, 15 Pakete) und die eine ungepinnte Reviewer-Probe als Dauertest nachgepinnt (TestACardHeavyWithFlagsStaysFiveRows: exakt 5 Zeilen bei w=12/18/24/40, keine Zeile ueberbreit). Nicht maschinell geprueft bleibt der Blick ins echte Terminal - review-check Schritt 1-4 ist Berks Abnahme."
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
