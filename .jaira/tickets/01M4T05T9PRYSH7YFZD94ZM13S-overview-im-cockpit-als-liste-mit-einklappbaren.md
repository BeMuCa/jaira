---
id: 01M4T05T9PRYSH7YFZD94ZM13S
title: Overview im Cockpit als Liste mit einklappbaren Abschnitten statt drei Spalten
status: critique
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
updated-at: 2026-10-13T14:51:24Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-131987
claimed-at: 2026-10-13T14:44:42Z
outcome-what: "Overview-Pane: zwei Tabs (Queue, Log), Abschnitte untereinander mit Klapp-Knopf und Anzahl, Titel umbrechen; Runde 2: In work/Queue/Put back je eigene Farbe an Ueberschrift und Handle, jede Person eine aus dem Namen abgeleitete Farbe in beiden Tabs, Lane gedimmt"
outcome-why: "Drei Spalten zu 33% schnitten jeden Titel ab; nach dem Handcheck wollte Alex Gruppen und Personen auf einen Blick auseinanderhalten"
outcome-resolves: "Alle acht DoD-Punkte: Tabs, Umbruch, Klappen, Zustand, Release-Note, Personenfarben, Gruppenfarben, Lesbarkeits-Durchgang"
review-summary: "core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:35 person(who, key) hat einen key-Parameter, den beide Aufrufer (:56, :79) mit 'who' fuellen - key='who' fest in person() schreiben und den Parameter streichen"
round: 4
review-gaps: "Keine Defekte gefunden. Schwaechen: (1) Umbrechen ist nur am Prop belegt (test 'a title is never cut off' prueft wrap='wrap'), nicht an einer echten schmalen Pane - ob die verschachtelte <Text bold>-Zeile im Terminal wirklich umbricht, zeigt erst der Handcheck. (2) Der Log-Tab hat keine Klapp-Ueberschrift und keine Anzahl - laut Note vom 13.10. 06:46 so gewollt, die Context-Zeile 'vierter Abschnitt' ist ueberholt. (3) Dass die Atome beim Sitzungsende verschwinden, belegt kein Test; es folgt dem bestehenden expandedAtom-Muster. (4) Kleinkram: die Box flexDirection=row um jeden Fold-Button (overview-pane.tsx:31) ist nur Huelle; review-summary stand vorher auf 'none'."
test-verdict: "pass: go test ./... RC=0; claude plugin test core/mod/builtin/jaira-cockpit 135/135 RC=0; plugin validate ok; DoD 1-5 gegen Baum und Proof-Zeilen geprueft; Verhalten im Test-Harness gemountet (Tabs, Klappen, Schliessen/Oeffnen, wrap). Nicht gesehen: das Pane in einem echten Terminal"
review-verdict: "Der Diff erfuellt alle fuenf DoD-Punkte; keine Defekte gefunden. Unsicher bin ich nur beim echten Umbruch im Terminal - der Test prueft das Prop, nicht die Darstellung."
review-round: 2
review-check: "1. Im Worktree: claude --plugin-dir /home/alex/projects/.worktrees/jaira-4zm13s/core/mod/builtin/jaira-cockpit  2. Im Band auf 'Overview' druecken  3. Erwartet: oben zwei Knoepfe 'Queue' (hervorgehoben) und 'Log'; darunter untereinander '▾ In work (n)', '▾ Queue (n)', '▾ Put back (n)'  4. Terminal schmal ziehen: ein langer Titel bricht in die naechste Zeile um statt mit … zu enden  5. Auf '▾ In work' druecken: Zeilen verschwinden, Knopf zeigt '▸ In work (n)'  6. Auf 'Log' druecken: nur noch die Geschichte der letzten 7 Tage  7. 'Close', dann wieder 'Overview': Log-Tab ist noch aktiv; auf 'Queue' - In work ist noch eingeklappt  8. Automatisch: claude plugin test core/mod/builtin/jaira-cockpit -> 135 pass"
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
- [x] Jeder Personenname im Pane (Queue und Log) hat eine Farbe, stabil aus dem Namen abgeleitet: dieselbe Person hat ueberall und bei jedem Oeffnen dieselbe Farbe, verschiedene Personen meist verschiedene
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:23 personColor (FNV-1a auf 6 ANSI-Farben), :35 person() in Queue :56 und Log :79; test 'a person has one colour, derived from the name: the same in both tabs and after reopening'
- [x] In work, Queue und Put back haben je eine eigene Farbe an Ueberschrift und Zeilen, sodass die drei Gruppen auf einen Blick unterscheidbar sind
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:16 GROUP claude/suggestion/warning, :48 Ueberschrift, :54 Handle jeder Zeile; test 'in work, queue and put back each carry their own colour, at the heading and at every row'
- [x] Ein Lesbarkeits-Durchgang ueber das Pane ist gemacht: was klein war, ist umgesetzt, und eine Note nennt, was geaendert und was bewusst gelassen wurde
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:54 Lane gedimmt, :78 put back farbig; Note 'Lesbarkeits-Durchgang' vom 13.10.

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
- **2026-10-13 14:44 · Alexander Sacharov** — Alex am 13.10. nach dem Handcheck: 'super' - Tabs, Klappen und Umbruch passen. Zurueck aus signoff fuer eine zweite Runde Lesbarkeit: (1) Namen der Leute farbig, jede Person eine eigene Farbe. Alex fragte 'vielleicht zufaellig' - entschieden (Teamlead): Farbe stabil aus dem Namen abgeleitet (Hash auf eine kleine Palette), damit dieselbe Person bei jedem Oeffnen und in beiden Tabs dieselbe Farbe hat; echtes Zufall wechselt bei jedem Oeffnen und hilft dann nicht beim Wiedererkennen. (2) In work, Queue und Put back je eine eigene Farbe, damit man sieht, dass es verschiedene Tickets sind. (3) Selbst pruefen, was das Pane sonst lesbarer macht, umsetzen was klein ist, den Rest in einer Note nennen.
- **2026-10-13 14:44 · Alexander Sacharov** — Zweite Runde an einen neuen Worker (Tab w3:tE9 geschlossen): Cache des alten seit ~10 h kalt (letzte Antwort 06:58, jetzt 16:44).
- **2026-10-13 14:44 · Alexander Sacharov** — parent session: jaira-6a
- **2026-10-13 14:48 · Alexander Sacharov** — Runde 2, entschieden (Worker, offene Entscheidungen keine fuer Alex): Gruppenfarben als Theme-Keys - In work 'claude', Queue 'suggestion', Put back 'warning' (folgen dem Theme). Personenfarben als ANSI-Namen (cyan, magenta, green + je Bright) per FNV-1a-Hash des Namens: getrennt von den drei Gruppenfarben, und ANSI folgt dem Terminal-Schema, also hell und dunkel lesbar. Farbe am Ticket-Handle (fett) und an der Ueberschrift (Text-Kind im Fold-Button, label bleibt), Lane gedimmt, damit der Titel heraussticht.
- **2026-10-13 14:50 · Alexander Sacharov** — Lesbarkeits-Durchgang: GEAENDERT - Lane und Trenner gedimmt (overview-pane.tsx:54, :77), damit Handle und Titel heraustreten; im Log ist 'put back' in der Put-back-Farbe (:78), wie im Queue-Tab; die Zeile unter einem Ticket zeigt nur noch, was da ist (kein leerer Text ohne Grund). BEWUSST GELASSEN - relative Zeiten ('vor 2 h') statt Uhrzeit: aendert overview.ts/clockOf, nicht klein; eigene Farben je Lane: wuerde mit den drei Gruppenfarben konkurrieren; Log-Tab einklappbar oder mit Anzahl: Alex hat am 13.10. ausdruecklich keine Klapp-Abschnitte im Log gewollt; Close-Knopf oben statt unten: Escape schliesst ohnehin.
- **2026-10-13 14:50 · Alexander Sacharov** — critique R2-1: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:35 person(who, key) hat einen key-Parameter, den beide Aufrufer (:56, :79) mit 'who' fuellen - key='who' fest in person() schreiben und den Parameter streichen. Stehen gelassen: Fold-Button mit label UND Text-Kind (overview-pane.tsx:48) - Button hat kein color-Prop, und label+farbiges Text-Kind ist das Muster von register.tsx:1080 (Accept/Give back); anders als critique 1 (aktiver Tab), wo variant existierte.
