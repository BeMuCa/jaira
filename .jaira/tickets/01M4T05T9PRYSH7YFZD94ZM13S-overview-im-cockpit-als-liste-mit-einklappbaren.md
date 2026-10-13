---
id: 01M4T05T9PRYSH7YFZD94ZM13S
title: Overview im Cockpit als Liste mit einklappbaren Abschnitten statt drei Spalten
status: optimize
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Das Overview-Pane zeigt seine Abschnitte untereinander statt nebeneinander, jeder Titel steht ganz da, und jeder Abschnitt laesst sich an seiner Ueberschrift ein- und ausklappen"
context: |-
  Was heute falsch ist: das Overview-Pane (seit 0.4.0) zeichnet 'in Arbeit', 'Warteschlange' und 'zurueckgestellt' als drei Spalten zu je 33%.
  Folge: jede Zeile wird abgeschnitten (wrap=truncate-end), Alex kann die Ticket-Titel nicht ganz lesen.
  Ausloeser: Alex am 13.10. nach dem Release von 0.4.0: 'zu wenig Platz, ich will den Namen ganz lesen'.
  Wunsch von Alex: eine Ueberschrift je Abschnitt, die man einklappen kann, darunter die Liste von oben nach unten.
  Entschieden (Teamlead, 13.10.): die Liste ERSETZT die Spalten, kein zweiter Modus daneben - ein Umschalter waere eine Einstellung mehr fuer dasselbe Pane.
  Ort: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx, drawOverview - die Funktion column() mit width=33% und die Box flexDirection=row darum.
  Die Geschichte der letzten Tage darunter ist schon eine Liste; sie wird der vierte Abschnitt mit derselben Ueberschrift-Geste.
  Vorbild fuer Klapp-Zustand: register.tsx expandedAtom + toggleMore (more/fewer-Knopf im Band) - ein Atom, ein Button, kein neues Muster.
  Daten bleiben wie sie sind: overview.ts (overviewOf) aendert sich nicht.
  Texte in hooks/texts.ts (overviewWords).
definition-of-done: "Die Abschnitte in Arbeit, Warteschlange, zurueckgestellt und Geschichte stehen untereinander, volle Breite"
tags:
  - tui
blocked-by: []
related: []
follows: 01M4RWVFXTPC084S81WSDSKZ70
commits: []
created-at: 2026-10-13T06:46:08Z
updated-at: 2026-10-13T06:55:23Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-6022
claimed-at: 2026-10-13T06:47:43Z
outcome-what: "Overview-Pane: zwei Tabs (Queue, Log); im Queue-Tab stehen In work, Queue, Put back untereinander in voller Breite, jede Ueberschrift ein Klapp-Knopf mit Anzahl; Titel umbrechen statt abgeschnitten zu werden; Tab und Klapp-Zustand in zwei Sitzungs-Atomen"
outcome-why: "Drei Spalten zu 33% schnitten jeden Titel ab; Alex konnte die Ticket-Namen nicht ganz lesen"
outcome-resolves: "Alle fuenf DoD-Punkte: Tabs mit Queue als Vorgabe, wrap statt truncate, Klapp-Knoepfe offen per Vorgabe, Zustand ueber Schliessen/Oeffnen, Release-Note"
review-summary: "none"
round: 2
review-gaps: "overview-pane.tsx:14-16 OverviewTab/State/Actions nicht mehr exportiert (niemand importiert sie). Gelassen: Klapp-Atom neben expandedAtom (register.tsx:57) - verschiedene Aufgaben, toggleMore laedt zusaetzlich das Board neu; texts.ts tabQueue='Queue' neben queue='Queue' - Tab und Abschnitt sind zwei Dinge, die nur heute gleich heissen; types/index.d.ts wiederholt den Tab-Typ absichtlich (self-contained contract)."
---

# Overview im Cockpit als Liste mit einklappbaren Abschnitten statt drei Spalten

## Definition of Done

- [x] Das Pane hat oben zwei Tabs: 'Queue' (in Arbeit, Warteschlange, zurueckgestellt - untereinander, volle Breite) und 'Log' (die Geschichte der letzten Tage); Queue ist beim Oeffnen aktiv
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:68 Tabs Queue/Log (aktiver variant=primary), :46 Abschnitte untereinander ohne width; core/mod/builtin/jaira-cockpit/hooks/register.tsx:61 overviewTabAtom Vorgabe 'queue'; test 'Overview reads the board afresh and shows in work, queue and put back, each row Handle | Lane | Titel'
- [x] Ein Ticket-Titel wird nicht mehr abgeschnitten, sondern umbricht bei Bedarf
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:38 wrap="wrap" statt truncate-end; test 'a title is never cut off: rows wrap'
- [x] Jede Abschnitt-Ueberschrift im Queue-Tab ist ein Knopf mit Anzahl, der den Abschnitt ein- und ausklappt; beim ersten Oeffnen sind alle aufgeklappt
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:26 section(): Knopf '▾/▸ Name (n)' klappt; Vorgabe offen (folded leer); test 'a heading folds its section, and tab and folds outlive closing the pane'
- [x] Klapp-Zustand und gewaehlter Tab bleiben erhalten, solange die Sitzung laeuft, auch ueber Schliessen und Wiederoeffnen des Panes
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:61-62 overviewTabAtom/overviewFoldedAtom (Sitzungs-Atome, nicht beim Oeffnen zurueckgesetzt), :1269; test 'a heading folds its section, and tab and folds outlive closing the pane'
- [x] Release-Note als Datei core/release/notes/4ZM13S.md (seit AR97K9 eine Datei je Aenderung)
  proof: core/release/notes/4ZM13S.md:1

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-13 06:46 · Alexander Sacharov** — Alex am 13.10., nachgereicht: oben Tabs statt Tabelle-und-darunter-Log. Tab 'Queue' traegt die einklappbaren Abschnitte, Tab 'Log' die Geschichte. Die Geschichte ist damit KEIN vierter Abschnitt im Queue-Tab (Context-Zeile dazu ueberholt). Tab-Umschalter wie expandedAtom: ein Atom, Buttons.
- **2026-10-13 06:47 · Alexander Sacharov** — parent session: jaira-78
- **2026-10-13 06:52 · Alexander Sacharov** — Gefunden: ein Atom-update aus einem Pane-Knopf zeichnet das gemountete Pane neu (Test drueckt und liest auf demselben Handle) - kein refresh noetig. Zweimal $.ui.mount(OVERVIEW) ohne unmount wirft 'mounted already'. Abschnitt-Knopf traegt Pfeil+Anzahl im label (▾ offen, ▸ zu); aktiver Tab fett, der andere gedimmt, als Button-Kind. Tests: claude plugin test core/mod/builtin/jaira-cockpit
- **2026-10-13 06:53 · Alexander Sacharov** — critique 1: tabButton (overview-pane.tsx:23) markiert den aktiven Tab ueber ein Text-Kind; die Engine hat dafuer Button-Props variant='primary' und dimColor - die nehmen, kein Kind. Stehen gelassen: types/index.d.ts wiederholt 'queue'|'log' absichtlich (Kopf der Datei: self-contained contract, mirrors hooks/*); inline update-Callbacks in register.tsx:1269 statt eigener toggle-Funktion, weil anders als toggleMore kein refreshBoard folgt.
- **2026-10-13 06:54 · Alexander Sacharov** — critique 2: Befund 1 behoben (292d41a9, overview-pane.tsx:23 variant/dimColor); nichts weiter.
