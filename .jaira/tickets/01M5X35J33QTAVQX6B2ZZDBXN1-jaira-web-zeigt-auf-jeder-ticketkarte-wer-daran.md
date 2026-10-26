---
id: 01M5X35J33QTAVQX6B2ZZDBXN1
title: "jaira web zeigt auf jeder Ticketkarte, wer daran arbeitet"
status: critique
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
commits: []
created-at: 2026-10-26T21:51:50Z
updated-at: 2026-10-26T22:17:19Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-72427
claimed-at: 2026-10-26T21:52:15Z
mode: ""
outcome-what: "Ticketkarten in jaira web zeigen den assignee (max 20 Zeichen, voller Name im Tooltip; entities/ticket/model/assignee.ts). Lane-, Tag- und Person-Filter nehmen mehrere Werte als Chips, ein Klick macht 'nicht X' (features/filter-tickets/model/filter.ts, ui/ChipSelect.svelte); Link ?lane=a&not-lane=b, alte ?lane=x-Links gehen weiter."
outcome-why: "Alex wollte bei einem Epic mit 124 Tickets sehen, wer welches hat, und Listen wie 'alle Lanes ausser done' filtern; beides ging nur mit Aufklappen bzw. einem Wert pro Filter."
outcome-resolves: "DoD 1: TicketCard.svelte zeigt span.assignee; DoD 2: npm test 10/10, dist neu gebaut, Release-Notiz; DoD 3: passes() in filter.ts, fromQuery/toQuery mit not-<key>, filter.test.ts."
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
