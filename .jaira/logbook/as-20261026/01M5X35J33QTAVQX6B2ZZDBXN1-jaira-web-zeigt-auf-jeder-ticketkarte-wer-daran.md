---
id: 01M5X35J33QTAVQX6B2ZZDBXN1
title: "jaira web zeigt auf jeder Ticketkarte, wer daran arbeitet"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Wer jaira web oeffnet, sieht auf jeder Ticketkarte den Namen der Person, der das Ticket zugewiesen ist, ohne die Karte aufzuklappen."
context: |-
  Heute zeigt eine Ticketkarte in jaira web nur Handle, Titel, Lane und 'only on its ref'. Wer am Ticket arbeitet, steht nirgends auf der Karte.
  Ausloeser: Alex hat am 26.10. ein Epic mit 124 Tickets im Web angesehen und wollte sehen, wer welches Ticket hat.
  Die Daten sind schon da: core/web/web.go liefert pro Ticket das Feld assignee (json 'assignee'), core/web/ui/src/entities/ticket/model/types.ts hat es auch. Gebraucht wird nur die Anzeige.
  Karten gibt es mindestens in core/web/ui/src/widgets/epic-list und widgets/unassigned-list.
  Wunsch von Alex: Name hoechstens 20 Zeichen, laengere abschneiden (voller Name im title-Tooltip). Ohne assignee bleibt die Stelle leer.
  Der Filter in features/filter-tickets/model/filter.ts kennt Personen schon; nicht umbauen.
definition-of-done: "Jede Ticketkarte in jaira web zeigt den assignee, auf hoechstens 20 Zeichen gekuerzt, voller Name als Tooltip; ohne assignee nichts"
tags:
  - web
blocked-by: []
related: []
commits:
  - c0ae756953f44ce6981202d927cec98d200834c7
  - 600a89ecdec00bd0065f90ec7cb79d38d5226089
  - ab18b38cb32d5d35a583a36c33fdd9a3f228c82f
created-at: 2026-10-26T21:51:50Z
updated-at: 2026-10-26T22:41:05Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-75425
claimed-at: 2026-10-26T22:22:49Z
mode: ""
outcome-what: "Ticketkarten in jaira web zeigen den assignee (max 20 Zeichen, voller Name im Tooltip; entities/ticket/model/assignee.ts). Lane-, Tag- und Person-Filter nehmen mehrere Werte als Chips, ein Klick macht 'nicht X' (features/filter-tickets/model/filter.ts, ui/ChipSelect.svelte); Link ?lane=a&not-lane=b, alte ?lane=x-Links gehen weiter."
outcome-why: "Alex wollte bei einem Epic mit 124 Tickets sehen, wer welches hat, und Listen wie 'alle Lanes ausser done' filtern; beides ging nur mit Aufklappen bzw. einem Wert pro Filter."
outcome-resolves: "DoD 1: TicketCard.svelte zeigt span.assignee; DoD 2: npm test 10/10, dist neu gebaut, Release-Notiz; DoD 3: passes() in filter.ts, fromQuery/toQuery mit not-<key>, filter.test.ts."
review-summary: "Jede Ticketkarte (TicketCard.svelte, genutzt von EpicList und UnassignedList) zeigt jetzt unter dem Titel eine kleine graue Zeile mit dem assignee: hoechstens 20 Codepoints inkl. '…' (shortAssignee in entities/ticket/model/assignee.ts), voller Name im title-Tooltip, ohne assignee keine Zeile. Zusaetzlich (von Alex am 26.10. nachbestellt, DoD 3): Lane-, Tag- und Person-Filter sind Chip-Listen statt Einzelwerten - Auswahl im Dropdown fuegt einen Chip hinzu, Klick schaltet 'X'/'not X', × entfernt (ChipSelect.svelte); innerhalb eines Filters ODER, ausgeschlossene Werte UND-NICHT, zwischen Filtern UND (passes() in filter.ts). Der Link traegt das als ?lane=a&lane=b&not-lane=c; alte ?lane=x-Links lesen sich als ein Chip. Hilfetext von 'jaira web' und zwei Zeilen Release-Notiz ergaenzt, dist neu gebaut."
review-gaps: |-
  note: die Kuerzung auf 20 Zeichen ist nur im Unit-Test (assignee.test.ts) belegt; auf diesem Board ist kein Name laenger, live also nicht gesehen
  note: Person 'not X' wirft ein Ticket raus, sobald X creator, assignee ODER updated_by ist - bewusste Spiegelung der Einschlussregel (decision-Note), kann aber ueberraschen, wenn man nur 'nicht zugewiesen an X' meint
  note: DoD 3 (Mehrfachfilter) ist eine Erweiterung ueber das urspruengliche Ziel hinaus; von Alex per Note 21:54 bestellt, deshalb kein Befund, aber Titel und goal des Tickets nennen sie nicht
  note: svelte-check meldet 1 Fehler in vite.config.ts (node:url-Typen); Datei ist vom Diff unberuehrt, laut Note schon auf master so
test-verdict: "pass: go vet RC=0, go test -race -count=1 ./... RC=0 (31 packages ok), npm ci+build+test 10/10 RC=0, committed dist matches the build (git diff --exit-code RC=0); DoD 1-3 checked in the tree; scratch board in headless Chromium: 'Maximiliane Müller-Lüdenscheidt' shows as 'Maximiliane Müller-…' with the full name as title, a ticket without assignee shows no line; tag web / not web / not-person / old ?tag=cli each leave exactly the expected card"
review-verdict: "Diff erfuellt DoD 1-3; keine Defekte gefunden. Selbst geprueft: npm test 10/10 gruen, frischer vite build ist byte-gleich mit core/web/ui/dist, go build und go test ./core/web/... ./internal/cli/... gruen, Release-Notiz im richtigen Format. Keine Blocker, nur Notes."
review-round: 1
review-check: |-
  1. cd core/web/ui && npm test  -> 10 Tests, 0 fail
  2. Im Repo-Root: jaira web  -> Browser oeffnet sich
  3. Auf jeder Karte steht unter dem Titel klein der Name (z.B. Alexander Sacharov); Maus drauf -> Tooltip mit vollem Namen; Karten ohne assignee haben keine solche Zeile
  4. Filterleiste: bei Lane 'done' waehlen  -> Chip 'done', nur erledigte Tickets
  5. Auf den Chip klicken  -> er heisst 'not done', gestrichelter Rand, alle Tickets ausser done; Adresszeile enthaelt not-lane=done
  6. Bei Lane zusaetzlich 'review' waehlen  -> zwei Chips; × am Chip entfernt ihn wieder aus Liste und Adresse
  7. Alten Link testen: Adresse auf ?lane=backlog setzen und laden  -> ein Chip 'backlog', nur Backlog-Tickets
---

# jaira web zeigt auf jeder Ticketkarte, wer daran arbeitet

## Definition of Done

- [x] Jede Ticketkarte in jaira web zeigt den assignee, auf hoechstens 20 Zeichen gekuerzt, voller Name als Tooltip; ohne assignee nichts
  proof: core/web/ui/src/entities/ticket/ui/TicketCard.svelte:27 (span.assignee, title=voller Name, nur bei assignee) + model/assignee.ts shortAssignee, Test assignee.test.ts; live: 250 von 267 Karten mit Namen, Rest ohne
- [x] npm test und svelte-check gruen, dist neu gebaut (web-dist-Check gruen), Release-Notiz core/release/notes/<ID>.md
  proof: npm test 5/5 exit 0; svelte-check nur der bekannte vite.config.ts node:url-Fehler (master); npm run build -> core/web/ui/dist; core/release/notes/ZDBXN1.md
- [x] Filter Lane, Tag und Person nehmen mehrere Werte (ODER innerhalb eines Filters, UND zwischen Filtern), und jeder Wert laesst sich ausschliessen ('nicht X'); beides steht im Link (fromQuery/toQuery) und ist in filter.ts getestet
  proof: core/web/ui/src/features/filter-tickets/model/filter.ts:87 passes() (ODER/UND-NICHT), :45-47 fromQuery getAll lane+not-lane, :62 toQuery; Tests filter.test.ts (5 Faelle inkl. alter ?lane=x-Link); UI ui/ChipSelect.svelte; live: all 267 = done 146 + not done 121

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-26 21:52 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/6787.sock
- **2026-10-26 21:52 · Alexander Sacharov** — decision: Name steht in einer eigenen Zeile unter dem Titel (grid-column 2/-1, klein, muted) wie 'only on its ref' — keine neue Spalte, damit die schmale Karte nicht breiter wird. Kuerzung: insgesamt hoechstens 20 Zeichen inklusive '…' (19 Zeichen + …), gezaehlt in Codepoints, nicht UTF-16-Einheiten; Helper in entities/ticket/model/assignee.ts mit Test. title-Tooltip nur wenn gekuerzt? Nein: immer der volle Name, einfacher und schadet nicht.
- **2026-10-26 21:54 · Alexander Sacharov** — Erweitert am 26.10. auf Wunsch von Alex: Filter werden Mehrfachauswahl mit Ausschluss, damit man z.B. alle Lanes ausser done oder alle Personen ausser einer sieht. Heute ist jeder Filter in core/web/ui/src/features/filter-tickets/model/filter.ts ein einzelner String (lane, tag, person). Bestehende Links mit ?lane=x muessen weiter funktionieren.
- **2026-10-26 21:54 · Alexander Sacharov** — svelte-check meldet 1 Fehler in vite.config.ts (node:url-Typen) — besteht auf master, nicht angefasst. Live-Check: Headless Chromium (uv run --with playwright) gegen 'jaira web --no-open' mit ?all=1: 267 Karten, 250 mit .assignee, title = voller Name; Karten ohne assignee zeigen keine Zeile. Kuerzung auf 20 nur im Unit-Test geprueft, auf diesem Board ist kein Name laenger.
- **2026-10-26 21:55 · Alexander Sacharov** — decision (selbst entschieden, DoD 3): Link-Format = wiederholte Parameter, Ausschluss unter eigenem Schluessel: ?lane=review&lane=testing&not-lane=done. Kein '!'/'-'-Praefix im Wert, weil ein Tag oder Name damit anfangen kann und URLSearchParams '!' zu %21 kodiert. Alter Link ?lane=x liest sich unveraendert als ein eingeschlossener Wert. Semantik: eingeschlossene Werte ODER, dann jeder ausgeschlossene UND-NICHT; Person 'nicht X' = X ist weder creator noch assignee noch updated_by (Spiegel der Einschluss-Regel).
- **2026-10-26 21:55 · Alexander Sacharov** — Frage (offen, an Alex): Wie waehlt man in der Filterleiste mehrere Werte und schliesst einen aus?
1. (Empfohlen) Chips: das Dropdown bleibt, eine Auswahl fuegt einen Chip hinzu; Klick auf den Chip schaltet zwischen 'X' und 'nicht X' um, x entfernt ihn. Kompakt, Zustand sofort lesbar, ein Klick mehr fuers Ausschliessen.
2. Aufklapp-Liste mit Checkboxen je Wert, dreistufig (egal / nur / nicht). Alles auf einen Blick, aber ein Popup pro Filter und mehr Code.
3. Zwei Dropdowns je Filter ('ist' und 'ist nicht'), jeweils mit Chips. Am einfachsten zu verstehen, verdoppelt aber die Leiste.
- **2026-10-26 22:10 · Alexander Sacharov** — Antwort Alex (26.10., via Parent): Variante 1, Chips. Auswahl im Dropdown fuegt einen Chip hinzu, Klick auf den Chip schaltet X / nicht X, x entfernt ihn.
- **2026-10-26 22:13 · Alexander Sacharov** — DoD 3 live geprueft (Headless Chromium gegen 'jaira web --no-open'): ?lane=backlog alter Link -> 76 Karten, ein Chip; 'done' waehlen 146, Chip umschalten 'not done' 121 (=267-146), Link not-lane=done&all=1; x entfernt den Chip aus dem Link. Nebenbei gesehen, nicht angefasst: steht ein Filter, versteckt ReviewPage leere Epics und zeigt dann 'No epics yet.' — irrefuehrend, aber schon vorher so. Hilfetext von 'jaira web' (internal/cli/web.go) um die Mehrfachauswahl ergaenzt, zweite Zeile in der Release-Notiz.
- **2026-10-26 22:14 · Alexander Sacharov** — Pause (conversational): DoD 3 fertig, nicht committed; Diff + Commit-Zeile an Parent geschickt, warte auf Alex' OK. Danach: Commit pruefen, dann move critique. .jaira/milestones/ (untracked) stammt von einem jaira-Aufruf, gehoert nicht in den Commit.
- **2026-10-26 22:17 · Alexander Sacharov** — Alex (via Parent): kein conversational mode auf diesem Ticket — mode geleert, DoD 3 selbst committet, weiter ohne Halt bis zur frischen Review.
- **2026-10-26 22:17 · Alexander Sacharov** — critique Runde 1: nichts. Geprueft und stehen gelassen: Select-Stil doppelt in ChipSelect.svelte und FilterBar.svelte — das Repo wiederholt Kontrollstile pro Komponente (CreateEpicForm.svelte:49, EditEpicForm.svelte:69), keine globale Regel in app/styles.css; Chip[] statt {is,not} haelt die Reihenfolge der Chips; Person-'nicht' spiegelt die Einschluss-Regel (decision-Note 21:52+).
- **2026-10-26 22:23 · Alexander Sacharov** — review (frischer Kontext): keine Blocker. Selbst nachgeprueft: npm test 10/10, vite build byte-gleich mit dist, go test core/web + internal/cli gruen. 4 Notes in review-gaps, keine davon verlangt Arbeit vor der Abnahme.
- **2026-10-26 22:23 · Alexander Sacharov** — Zu review-gaps Note 1: die Kuerzung wurde in testing doch live gesehen — Scratch-Board, 'Maximiliane Müller-Lüdenscheidt' -> 'Maximiliane Müller-…', title voll (test-verdict). Die Review las die aeltere Note.
- **2026-10-26 22:39 · Alexander Sacharov** — Alex (via Parent): angenommen wie es ist; Person-Filter bleibt 'jede Beteiligung' (review-gaps Note 2 bleibt, keine Aenderung).
- **2026-10-26 22:39 · Alexander Sacharov** — explain: Review ohne Blocker. Karten zeigen den assignee, Lane-/Tag-/Person-Filter sind Chip-Listen mit 'X'/'not X'. Offen ist nur: Person 'not X' wirft ein Ticket raus, sobald X creator, assignee ODER updated_by ist.
Warum: Wer 'nicht zugewiesen an X' meint, bekommt auch Tickets weg, die X nur angelegt oder zuletzt angefasst hat. Das ist bewusst so (Spiegel der Einschlussregel), kann aber ueberraschen.
1) Person 'not X' nur auf assignee beschraenken, Einschluss und Ausschluss sind dann nicht mehr symmetrisch [give back]
2) Ticket so annehmen, Person-Filter bleibt 'jede Beteiligung' [accept]
Empfehlung: 2 - das ist die dokumentierte decision-Note, die Review fand keinen Defekt, und Alex hat genau das schon angenommen.
- **2026-10-26 22:40 · Alexander Sacharov** — decision: 2) Ticket so annehmen, Person-Filter bleibt 'jede Beteiligung'
