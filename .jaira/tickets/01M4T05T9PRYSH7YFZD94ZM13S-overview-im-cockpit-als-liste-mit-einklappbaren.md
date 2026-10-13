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
updated-at: 2026-10-13T14:58:09Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-145384
claimed-at: 2026-10-13T14:56:12Z
outcome-what: "Overview-Pane: Tabs Queue/Log, Abschnitte untereinander mit Klapp-Knopf, Titel umbrechen; Gruppen- und Personenfarben, Lane gedimmt; Runde 3: Personenfilter im Log (All + Namensknoepfe, Mehrfachauswahl, Sitzungs-Atom)"
outcome-why: "Drei Spalten zu 33% schnitten jeden Titel ab; danach wollte Alex Gruppen und Leute auf einen Blick auseinanderhalten und das Log nach Leuten filtern"
outcome-resolves: "Alle neun DoD-Punkte: Tabs, Umbruch, Klappen, Zustand, Release-Note, Personenfarben, Gruppenfarben, Lesbarkeit, Personenfilter"
review-summary: "Das Overview-Pane (core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx) hat oben zwei Tab-Knoepfe Queue und Log statt drei Spalten zu 33%. Queue zeigt In work, Queue und Put back untereinander ohne Breitenvorgabe; jede Ueberschrift ist ein Knopf '▾/▸ Name (n)', der den Abschnitt klappt. Log zeigt die Geschichte der letzten 7 Tage. Zeilen haben wrap=wrap statt truncate-end, Titel brechen also um. Tab und Klapp-Zustand liegen in zwei Sitzungs-Atomen (register.tsx:61-62, overviewTab Vorgabe 'queue', overviewFolded leer = alles offen) und ueberleben Schliessen/Oeffnen. Runde 2: die drei Gruppen faerben Ueberschrift und Handle mit Theme-Farben (claude/suggestion/warning); jeder Personenname bekommt per FNV-1a-Hash eine von 6 ANSI-Farben, in beiden Tabs gleich; Lane und Trenner gedimmt; im Log ist 'put back' in der Put-back-Farbe. texts.ts ersetzt createdBy/parkedBy durch created/putBack/by. Release-Note mit zwei Zeilen in core/release/notes/4ZM13S.md; 5 neue Tests, 137 pass."
round: 5
review-gaps: "Keine Defekte. Klein: (1) Der Klapp-Knopf (overview-pane.tsx:48) gibt kein plain={plain} mit wie die Band-Knoepfe (register.tsx:1079), im Terminal zeichnet die Engine die Klammern daher ungefaerbt um die farbige Ueberschrift - rein optisch, beim Handcheck ansehen. (2) Das Feld definition_of_done nennt noch 'Geschichte' als vierten Abschnitt; durch Alex' Tab-Entscheidung vom 13.10. ueberholt, die DoD-Punkte selbst stimmen. (3) Die Tests pruefen Props (wrap, color), nicht die echte Darstellung im Terminal - Umbruch und Farben nur per Hand pruefbar."
test-verdict: "pass: go build ./... RC=0, go test ./... RC=0, claude plugin test core/mod/builtin/jaira-cockpit RC=0 (137 pass, 2 neu: Gruppenfarben an Ueberschrift+Handle, Personenfarbe gleich in beiden Tabs und nach Wiederoeffnen); DoD 1-8 an den proof-Zeilen im Baum geprueft; Pane im Engine-Testharness gemountet und gelesen. Nicht gesehen: wie das Terminal die ANSI-Namen (cyanBright usw.) tatsaechlich zeichnet - das zeigt erst der Handcheck"
review-verdict: "Der Diff erfuellt alle acht DoD-Punkte, jeder Beleg stimmt mit der Datei ueberein; 137 Tests gruen. Keine Defekte gefunden. Unsicher nur bei der echten Terminal-Darstellung (Umbruch, Farbkontrast der 6 Personenfarben auf hellem Theme) - das prueft kein Test."
review-round: 4
review-check: "1. claude --plugin-dir /home/alex/projects/.worktrees/jaira-4zm13s/core/mod/builtin/jaira-cockpit  2. Im Band 'Overview' druecken  3. Erwartet: oben 'Queue' (hervorgehoben) und 'Log' (gedimmt); darunter untereinander '▾ In work (n)', '▾ Queue (n)', '▾ Put back (n)', jede Ueberschrift in eigener Farbe, die Handles darunter in derselben Farbe, Lane grau  4. Personennamen unter den Zeilen sind farbig; derselbe Name hat ueberall dieselbe Farbe  5. Terminal schmal ziehen: ein langer Titel bricht um statt mit … zu enden  6. Auf '▾ In work' druecken: Zeilen verschwinden, Knopf zeigt '▸ In work (n)'  7. 'Log' druecken: Tage mit Uhrzeit, 'put back' farbig, Namen in derselben Farbe wie im Queue-Tab  8. 'Close', dann wieder 'Overview': Log ist noch aktiv; auf 'Queue': In work noch eingeklappt  9. Automatisch: claude plugin test core/mod/builtin/jaira-cockpit -> 137 pass"
---

# Overview im Cockpit als Liste mit einklappbaren Abschnitten statt drei Spalten

## Definition of Done

- [x] Das Pane hat oben zwei Tabs: 'Queue' (in Arbeit, Warteschlange, zurueckgestellt - untereinander, volle Breite) und 'Log' (die Geschichte der letzten Tage); Queue ist beim Oeffnen aktiv
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:102 Tabs Queue/Log (aktiver variant=primary), :47 Abschnitte untereinander ohne width; core/mod/builtin/jaira-cockpit/hooks/register.tsx:61 overviewTabAtom Vorgabe 'queue'; test 'Overview reads the board afresh and shows in work, queue and put back, each row Handle | Lane | Titel'
- [x] Ein Ticket-Titel wird nicht mehr abgeschnitten, sondern umbricht bei Bedarf
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:55 wrap="wrap" statt truncate-end; test 'a title is never cut off: rows wrap'
- [x] Jede Abschnitt-Ueberschrift im Queue-Tab ist ein Knopf mit Anzahl, der den Abschnitt ein- und ausklappt; beim ersten Oeffnen sind alle aufgeklappt
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:43 section(): Knopf '▾/▸ Name (n)' klappt (:49); Vorgabe offen (folded leer); test 'a heading folds its section, and tab and folds outlive closing the pane'
- [x] Klapp-Zustand und gewaehlter Tab bleiben erhalten, solange die Sitzung laeuft, auch ueber Schliessen und Wiederoeffnen des Panes
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:61-63 overviewTabAtom/overviewFoldedAtom (Sitzungs-Atome, nicht beim Oeffnen zurueckgesetzt), :1273; test 'a heading folds its section, and tab and folds outlive closing the pane'
- [x] Release-Note als Datei core/release/notes/4ZM13S.md (seit AR97K9 eine Datei je Aenderung)
  proof: core/release/notes/4ZM13S.md:1
- [x] Jeder Personenname im Pane (Queue und Log) hat eine Farbe, stabil aus dem Namen abgeleitet: dieselbe Person hat ueberall und bei jedem Oeffnen dieselbe Farbe, verschiedene Personen meist verschiedene
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:23 personColor (FNV-1a auf 6 ANSI-Farben), :36 person() in Queue :57 und Log :93; test 'a person has one colour, derived from the name: the same in both tabs and after reopening'
- [x] In work, Queue und Put back haben je eine eigene Farbe an Ueberschrift und Zeilen, sodass die drei Gruppen auf einen Blick unterscheidbar sind
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:16 GROUP claude/suggestion/warning, :49 Ueberschrift, :55 Handle jeder Zeile; test 'in work, queue and put back each carry their own colour, at the heading and at every row'
- [x] Ein Lesbarkeits-Durchgang ueber das Pane ist gemacht: was klein war, ist umgesetzt, und eine Note nennt, was geaendert und was bewusst gelassen wurde
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:55 Lane gedimmt, :92 put back farbig; Note 'Lesbarkeits-Durchgang' vom 13.10.
- [x] Im Log-Tab steht unter der Tab-Zeile eine Filterzeile: 'All' und je ein Knopf pro Person aus dem Log, in deren Farbe; ein Namensknopf schaltet die Person an/aus (mehrere gleichzeitig), 'All' hebt die Auswahl auf; das Log zeigt nur Ereignisse der gewaehlten Personen, Tage ohne Treffer fallen weg; die Auswahl haelt fuer die Sitzung
  proof: core/mod/builtin/jaira-cockpit/hooks/overview-pane.tsx:72-82 Filterzeile All + Namensknoepfe in Personenfarbe, :74 Tage ohne Treffer fallen weg; core/mod/builtin/jaira-cockpit/hooks/register.tsx:63 overviewPeopleAtom, :1274-1275 person/everybody; test 'the Log filters to the people chosen, several at once, All clears it, and the choice outlives closing the pane'

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
- **2026-10-13 14:51 · Alexander Sacharov** — critique R2-2: Befund R2-1 behoben (person() ohne key-Parameter, overview-pane.tsx:35); nichts weiter.
- **2026-10-13 14:55 · Alexander Sacharov** — Alex am 13.10., Runde 3 (zurueck aus signoff): im Log-Tab Leute filtern. Unter der Tab-Zeile eine zweite Zeile mit Knoepfen: 'All' und dann je ein Knopf pro Name, der im Log vorkommt; man waehlt, wen man sehen will. Entschieden (Teamlead): Mehrfachauswahl - jeder Namensknopf schaltet seine Person an/aus, wie u auf dem Board mehrere Leute ankreuzt; 'All' hebt die Auswahl auf; nichts gewaehlt = alle. Namensknoepfe in der Personenfarbe aus Runde 2. Auswahl haelt die Sitzung wie Tab und Klapp-Zustand. Ein Tag ohne passendes Ereignis verschwindet. Nur Log; Queue bleibt ungefiltert.
- **2026-10-13 14:56 · Alexander Sacharov** — Runde 3 an denselben Worker (jaira-4zm13s-7d) per SendMessage: Cache warm, letzte Antwort 14:55 UTC, jetzt 14:55 UTC.
- **2026-10-13 14:57 · Alexander Sacharov** — Runde 3: Personenfilter im Log - Atom overviewPeople (string[], leer = alle), Namen alphabetisch aus allen Log-Ereignissen; ein Ereignis ohne Namen faellt bei aktiver Auswahl weg (gehoert zu niemandem Gewaehlten). Gewaehlter Knopf variant=primary, sonst dimColor, wie die Tabs. Dabei alle --proof-Zeilen neu gesetzt: die aus Runde 1 waren nach Runde 2 um einige Zeilen verrutscht.
