---
id: 01M4ZAV5XPAN4B89W2V1TPCKZB
title: "Web-Seite: Filter und Tickets aus Logbook und Archiv anzeigen"
status: signoff
ready: true
creator: Alexander Sacharov
goal: "Auf der Seite von 'jaira web' sieht man jedes Ticket, das es je gab — offen, nur auf dem Ref, im Logbuch oder im Archiv — und filtert es nach Text, Lane, Tag, Person und Datum, mit den Filtern in der URL."
context: |-
  Was heute fehlt:
  - 'jaira web' zeigt nur offene Tickets: die Liste 'Tickets in no epic' und die Epics (core/web/web.go, Routen /api/epics und /api/unassigned).
  - Es gibt keinen Filter auf der Seite. Wer ein Ticket sucht, scrollt.
  - Tickets, die 'jaira logbook' abgelegt hat (.jaira/logbook/<wer>-<datum>/), und Tickets in .jaira/archive/ sind auf der Seite gar nicht zu sehen.
  Auslöser: Alex am 2026-10-15 — will alle Tickets aus logbook und archive auf der Seite ansehen können.
  Bekannt:
  - Die UI liegt in core/web/ui/src (Svelte, widgets/ und features/); dist ist eingebettet.
  - 'jaira list' kennt --query, --tag, --milestone, --lane, --assignee — dieselben Filter wären für die Seite naheliegend.
  Offen für brainstorm: welche Filter zuerst (Text, Tag, Lane, Assignee?), und ob logbook/archive eine eigene Ansicht oder ein Schalter 'auch abgelegte zeigen' sind.
definition-of-done: "Auf der Seite gibt es ein Textfeld, das Tickets nach Titel, Ziel und ID filtert"
tags:
  - web
blocked-by: []
related: []
commits: []
created-at: 2026-10-15T08:28:46Z
updated-at: 2026-10-19T07:38:24Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-74288
claimed-at: 2026-10-19T06:35:35Z
outcome-what: "jaira web: Filterleiste, Schalter 'All tickets ever' (Board, Ref, Logbuch, Archiv mit Ort), Ablehnung im Detail, Filter im URL-Query; dazu Lane-Farben an Badge und Lane-Filter aus der Board-Reihenfolge (GET /api/lanes)"
outcome-why: "Alex will jedes Ticket, das es je gab, auf der Seite sehen und nach Personen, Datum und Lane filtern; bisher zeigte die Seite nur offene Tickets ohne Filter"
outcome-resolves: "GET /api/unassigned?all=1 und neue Kartenfelder (web_source.go cards/rejection), Filtermodell features/filter-tickets mit URL-Abgleich, TestWebListsFiledTicketsAndOpensThem; DoD 1-8 mit Proofs"
review-summary: "jaira web bekommt eine Filterleiste ueber beiden Spalten (Text ueber Titel/Ziel/ID, Lane, Tag, Person = creator|assignee|updated_by, Datumsbereich auf erstellt/geaendert/abgelegt), alles UND-verknuepft und im Browser gefiltert; der Zustand steht per history.replaceState im URL-Query. Der Schalter 'All tickets ever' ruft GET /api/unassigned?all=1: der Server (internal/cli/web_source.go cards) liefert dann auch fertige Board-Tickets, Ref-Tickets und alle Tickets aus .jaira/logbook und .jaira/archive ohne Epic, jedes mit place, abgelegte nach Ablagedatum sortiert (Logbuch-Ordnerdatum bzw. rejected-at). Abgelegte Karten sind nicht ziehbar; die Detailansicht zeigt den Ort und bei abgelehnten Archiv-Tickets rejected-kind/why/by/at (direkt aus dem Frontmatter gelesen, damit schema.go fuer A0C081 frei bleibt). Neu GET /api/lanes: Lanes in Board-Reihenfolge mit kind; LaneBadge faerbt Badge und Lane-Filter in oklch (Arbeits-Lanes nach Position, human amber, terminal gruen, parking rot). link.placeOf wurde zu link.PlaceOf exportiert. Release-Note und web --help ergaenzt."
round: 3
review-gaps: "Keine Luecke gegen die DoD 1-9. Klein, nicht blockierend: 1) ReviewPage.svelte load(): schaltet man 'All tickets ever' schnell um, kann eine spaete Antwort von ?all=1 eine neuere von ?all=0 ueberschreiben (kein Abbruch/Versionszaehler) - die Liste zeigt dann abgelegte Tickets bei ausgeschaltetem Schalter bis zum naechsten Minuten-Reload. 2) EpicList <details open={expanded}>: Filter setzen und wieder leeren klappt auch Epics zu, die man von Hand aufgeklappt hatte. 3) Epics() liest jetzt bei jedem Minuten-Reload alle abgelegten Header (cards) - in den Notizen gemessen 0,08 s, akzeptiert. 4) Lane-Farben haengen an der Reihenfolge der Lanes; in Notiz und Release-Note offen benannt. Kein JS-Test fuer filter.ts im Repo (nur Wegwerf-Skript laut Notiz). Geprueft: go vet, go test core/web internal/cli core/link gruen; npm run build laesst core/web/ui/dist unveraendert."
test-verdict: "pass: go vet ./... RC=0, go test -count=1 -race ./... RC=0, dist matches npm ci && npm run build; DoD 1-9 verified in headless Chrome/Playwright on this repo and a scratch board with a rejected archive ticket"
review-verdict: "Der Diff erfuellt alle neun DoD-Punkte; Outcome und Diff stimmen ueberein. Keine Defekte, die eine Rueckgabe rechtfertigen - nur die zwei kleinen UI-Unschoenheiten aus review-gaps. Unsicher bleibt nur die Optik im Browser, die ich nicht selbst angesehen habe (testing hat Headless-Chrome/Playwright-Laeufe notiert)."
review-round: 3
review-check: "1. cd /home/alex/projects/.worktrees/jaira-tpckzb && go run ./cmd/jaira web  2. die ausgegebene Adresse im Browser oeffnen: ueber beiden Spalten steht eine Leiste Search / Lane / Tag / Person / Date / 'All tickets ever'  3. 'web' ins Suchfeld tippen: beide Spalten schrumpfen, die Adresszeile zeigt ?q=web  4. Haken 'All tickets ever' setzen: links erscheinen viele Tickets mit Zusatzzeile 'in the logbook' bzw. 'in the archive', Adresse enthaelt all=1  5. Date auf 'filed away', von und bis 2026-10-17: es bleiben 8NHD8Q und S0EG0M  6. 8NHD8Q anklicken: Dialog zeigt farbiges Lane-Badge 'done' (gruen) und 'in the logbook'  7. eine logbook-Karte ziehen wollen: sie laesst sich nicht greifen  8. Lane-Auswahl oeffnen: jede Lane hat einen farbigen Punkt, dieselbe Farbe wie die Badges; Seite neu laden (F5): Filter und Farben sind unveraendert  9. Adresse kopieren, in neuem Tab oeffnen: dieselbe gefilterte Liste. Abgelehntes Archiv-Ticket: go test ./internal/cli -run TestWebListsFiledTicketsAndOpensThem (prueft kind und why)"
---

# Web-Seite: Filter und Tickets aus Logbook und Archiv anzeigen

## Definition of Done

- [x] Auf der Seite gibt es ein Textfeld, das Tickets nach Titel, Ziel und ID filtert
  proof: core/web/ui/src/features/filter-tickets/model/filter.ts:71 (title/goal/handle/id, any case), FilterBar.svelte Search field
- [x] Filter nach Lane, Tag und Person (creator, assignee, zuletzt geaendert von) - mehrere Filter wirken zusammen (UND)
  proof: core/web/ui/src/features/filter-tickets/model/filter.ts:70-84 matches() ANDs lane, tag and person (creator|assignee|updated_by); FilterBar.svelte selects
- [x] Tickets aus .jaira/logbook/ sind auf der Seite sichtbar und öffnen sich in der Detailansicht
  proof: internal/cli/web_source.go:120 cards() reads the logbook; TestWebListsFiledTicketsAndOpensThem; live: ?all=1 on this repo lists 129 logbook tickets, 8NHD8Q opens with place logbook
- [x] Tickets aus .jaira/archive/ sind auf der Seite sichtbar und öffnen sich in der Detailansicht
  proof: internal/cli/web_source.go:120 cards() reads the archive; TestWebListsFiledTicketsAndOpensThem; live: 5 archive tickets listed
- [x] core/release/notes/<ID>.md beschreibt die neuen Filter und Ansichten
  proof: core/release/notes/TPCKZB.md
- [x] Filter nach Datum: Zeitraum von-bis auf erstellt, zuletzt geaendert und abgelegt (Logbuch-Datum)
  proof: core/web/ui/src/features/filter-tickets/model/filter.ts:78 date range on created/updated/filed; filed_at from ticket.LogbookFolderDay (web_source.go cards) or rejected-at
- [x] Ein abgelehntes Ticket im Archiv zeigt rejected-kind und rejected-why
  proof: internal/cli/web_source.go:165 rejection(); TicketDetails.svelte:23; TestWebListsFiledTicketsAndOpensThem checks kind and why
- [x] Die gesetzten Filter stehen in der URL; ein Link zeigt dieselbe gefilterte Liste
  proof: core/web/ui/src/features/filter-tickets/model/filter.ts fromQuery/toQuery; ReviewPage.svelte:26 reads, :40 writes via history.replaceState
- [x] Jede Lane hat im Web eine eigene Farbe — dieselbe am Lane-Badge jedes Tickets und im Lane-Filter; die Farbe einer Lane ist stabil (gleiche Lane, gleiche Farbe bei jedem Laden)
  proof: core/web/ui/src/entities/lane/model/colour.svelte.ts laneColour (board order via GET /api/lanes, internal/cli/web_source.go Lanes, TestWebListsLanesInBoardOrderWithKind); LaneBadge in TicketCard/TicketDetails, FilterBar lane select; checked light+dark in headless Chrome

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/web: Ticket bekommt place, tags, creator, assignee, updated_by, created_at, updated_at, filed_at, filed; GET /api/unassigned?all=1; TicketDetail bekommt place + rejected
- [x] internal/cli/web_source.go: Karten aus Board+Ref+Logbuch+Archiv fuellen, rejected-* ueber Doc.Fields lesen, Tests
- [x] UI: features/filter-tickets (Text, Lane, Tag, Person, Datum, Schalter), Zustand im URL-Query, wirkt auf beide Spalten; Ort-Badge; abgelegte nicht ziehbar; Detail zeigt Ort und Ablehnung
- [x] UI bauen (dist committed), web --help und core/release/notes/TPCKZB.md

## Progress
- **2026-10-19 06:35 · Alexander Sacharov** — parent session: jaira-a7
- **2026-10-19 06:35 · Alexander Sacharov** — decision (Alex, 19.10.2026 via teamlead jaira-a7): 'we must see ALL tickets that ever existed and filter by people, date or lane.' Abgelegte Tickets (logbook, archive, refs) kommen in DIESELBE Liste mit einem Schalter 'alle Tickets aller Zeiten'; jedes zeigt, wo es liegt — keine eigene Ansicht. Damit sind die brainstorm-Fragen beantwortet. Offen fuer den Plan: genaue Personenfelder, Server-API vs. Filtern im Browser. core/ticket/store.go Archive/Restore nicht anfassen — A0C081 aendert sie; rejected-kind/why/by/at (flache Keys aus A0C081) anzeigen, wenn vorhanden, nicht auf den Merge warten.
- **2026-10-19 06:37 · Alexander Sacharov** — brainstorm — was der Code zeigt:
- core/web/web.go hat nur /api/epics und /api/unassigned; internal/cli/web_source.go baut sie aus boardIndex + epic.Unassigned — also nur Tickets unter .jaira/tickets/ und ref-only (Store.List haengt Source.Extra an).
- Die Detailansicht kann abgelegte Tickets schon: webSource.Ticket nutzt Store.LoadAnywhere (Board, dann logbook, dann archive). Es fehlt nur ein Weg dorthin.
- Filed-Tickets findet Store.FiledAwayIDs() (archive + alle logbook-Ordner, auch der Legacy-Name); Logbuch-Datum steckt im Ordnernamen, ticket.LogbookFolderDay liest es.
- rejected-* gibt es auf master noch nicht (A0C081 ist uncommitted in ../jaira-a0c081, fuegt Ticket.Rejected in core/ticket/schema.go hinzu). Lesbar ohne schema.go anzufassen: Doc.Fields().Scalar("rejected-kind").
- Groesse heute: 134 abgelegte + ~29 offene Tickets.

Wege:
A) Ein Endpoint GET /api/tickets?all=1 liefert jede Karte mit den Filterfeldern (lane, tags, creator, assignee, updated_by, created/updated/filed, place); gefiltert wird im Browser, Filter im URL-Query. Kostet: ein Endpoint, eine Filterleiste. Gibt: sofortiges Filtern, kein Request pro Tastendruck; bei tausenden Tickets wird die Antwort gross.
B) Filter serverseitig (?q=&tag=&lane=...), wie 'jaira list'. Kostet: Filterlogik in Go doppelt zu 'jaira list' oder deren Refactor; jeder Tastendruck ein Request mit Board-Scan. Gibt: kleine Antworten.
C) Eigene Seite /all. Von Alex ausgeschlossen (gleiche Liste, Schalter).

Ich nehme A: die Datenmenge ist klein, Filtern im Browser ist sofort, die URL traegt den Zustand ohne Server-Logik.
- **2026-10-19 06:37 · Alexander Sacharov** — decision (worker, aus Code/Default): Die Filterleiste steht ueber beiden Spalten und wirkt auf beide — auf 'Tickets in no epic' und auf die Ticketzeilen jeder Epic-Karte. Der Schalter 'alle Tickets aller Zeiten' fuegt der linken Liste die abgelegten Tickets ohne Epic hinzu (Epic-Karten zaehlen abgelegte schon heute). Jede Karte zeigt, wo sie liegt: Board / Ref / Logbuch / Archiv. Abgelegte Karten lassen sich nicht ziehen. Personenfilter: EIN Feld 'Person', trifft wenn creator ODER assignee ODER updated_by passt — 'filter by people' ohne drei Dropdowns. Datumsfilter: Auswahl des Feldes (erstellt / geaendert / abgelegt) plus von/bis. Abgelegt-Datum: Logbuch = Datum aus dem Ordnernamen; Archiv = rejected-at, wenn gesetzt, sonst keins. Die Snapshot-Branch (core/snapshot) ist ein Backup und bleibt aussen vor.
- **2026-10-19 06:40 · Alexander Sacharov** — Go-Seite steht: web.Ticket traegt place/tags/creator/assignee/updated_by/created_at/updated_at/filed_at/rejected, GET /api/unassigned?all=1 liefert zusaetzlich fertige Board-Tickets und alle abgelegten ohne Epic; TicketDetail traegt place und rejected. rejected-* wird in internal/cli/web_source.go:rejection per Doc.Fields gelesen — schema.go bleibt fuer A0C081 frei; nach dessen Merge kann rejection() auf t.Rejected umgestellt werden. Test: TestWebListsFiledTicketsAndOpensThem. Als naechstes UI.
- **2026-10-19 06:43 · Alexander Sacharov** — UI fertig. Kein JS-Testrunner im Repo (package.json hat nur vite build); filter.ts wurde mit einem Wegwerf-Skript unter node --experimental-strip-types gegen 15 Faelle geprueft (Text, Lane, Tag, Person, Datum je Feld, URL-Roundtrip, ohne filed_at). svelte-check: 0 Warnungen; der eine Fehler in vite.config.ts (node:url-Typen) besteht schon auf master. Live auf diesem Repo: ?all=1 = 222 Tickets in 0,19 s. Die Seite selbst wurde nicht im Browser angesehen — das ist fuer die Abnahme.
- **2026-10-19 06:45 · Alexander Sacharov** — critique (Runde 1): 1) filedPlace in internal/cli/web_source.go ist eine Kopie von placeOf in core/link/link.go:196 — exportieren (link.PlaceOf) statt kopieren. 2) Ortsbezeichnungen stehen doppelt (TicketCard PLACES, TicketDetails if-Kette) — eine Map in entities/ticket/model/place.ts. Stehen gelassen: cards() liest bei jedem /api/epics alle abgelegten Header (gemessen 0,08 s bei 134 Dateien); rejection() liest die Archivdatei ein zweites Mal — beides klein, und rejection() wird nach A0C081 durch t.Rejected ersetzbar; die Doppelung epic.Unassigned / everything-Zweig bleibt, damit die Standardliste genau 'jaira epic ls --unassigned' entspricht.
- **2026-10-19 06:46 · Alexander Sacharov** — critique Runde 1 umgesetzt: link.PlaceOf exportiert (core/link/link.go:196) und in web_source.go benutzt; Ortsbezeichnungen als PLACE_LABELS in entities/ticket/model/place.ts, TicketCard und TicketDetails lesen sie. Proofs 3/4/7 auf neue Zeilen gesetzt.
- **2026-10-19 06:46 · Alexander Sacharov** — critique Runde 2: beide Befunde aus Runde 1 erledigt (filedPlace weg, link.PlaceOf benutzt; PLACE_LABELS in beiden Komponenten). Nichts weiter — gemaess Lane-Regel nur diese Stellen gelesen.
- **2026-10-19 06:48 · Alexander Sacharov** — decision (Alex, 19.10.2026 via teamlead jaira-a7): 'lanes must be told apart by colour' im Web. Vorgaben: kein neues Lane-Frontmatter-Feld ohne Rueckfrage; Quelle selbst waehlen (z.B. feste Palette nach Lane-Reihenfolge des Boards, terminal/human erkennbar); lesbar in hellem und dunklem Theme. Heute: .jaira/lanes/*.md tragen keine Farbe; die TUI wechselt nur einen Lane-Schatten (internal/tui/view.go:656 laneShade).
- **2026-10-19 06:49 · Alexander Sacharov** — testing Runde 1: go vet ./... RC=0, go test -count=1 -race ./... RC=0, npm ci && npm run build laesst dist unveraendert (CI-Check web-dist). Headless-Chrome auf diesem Repo: ?all=1&q=web&person=Alexander+Sacharov zeigt 3 Karten mit Ort (ref, logbook); ?all=1&date=filed&from=2026-10-17&to=2026-10-17 zeigt genau 8NHD8Q und S0EG0M (Ordner as-20261017). Playwright auf Scratch-Board mit abgelehntem Archiv-Ticket: Schalter setzt ?all=1, Karte nicht ziehbar, Dialog zeigt 'in the archive' und 'Rejected: wontfix | Ada Lovelace,10/2/26 | nobody needs it'. Befunde: 1) DoD 9 (Lane-Farben) ist neu und offen. 2) TicketDetails.svelte: zwischen rejected.by und Datum fehlt das Leerzeichen nach dem Komma (Svelte schneidet das Leerzeichen am Ende des {#if}-Blocks ab) — als {', '} schreiben.
- **2026-10-19 06:54 · Alexander Sacharov** — decision (worker, DoD 9): Farbquelle = Lane-Reihenfolge des Boards, kein Frontmatter-Feld. GET /api/lanes liefert Lanes in Board-Reihenfolge mit kind (work/human/terminal/parking; human = requires-human-exit ODER requires-question, wie der Agent-Block in CLAUDE.md es markiert). Die Seite faerbt in oklch: Arbeits-Lanes bekommen der Reihe nach 9 Farbtoene (240,165,290,110,340,215,190,265,315 — wechselnd, damit Nachbarn sich unterscheiden; in-progress bewusst magenta statt blau, weil backlog blau ist), die 10. grau, dann wieder von vorn; human = amber 70, terminal = gruen 145, parking = rot 25. Helligkeit/Chroma aus dem Theme (--lane-l/--lane-c in app/styles.css), Licht und Dunkel per Screenshot geprueft. Eine Lane eines anderen Boards bleibt gedaempft (--muted). Grenze: aendert jemand die Lane-Reihenfolge, wechseln die Farben der Arbeits-Lanes mit. Ausserdem: Epic-Karten klappen ihre Ticketliste auf, solange ein Filter eingrenzt (sonst blieben Treffer zugeklappt); testing-Befund 2 (Komma ohne Leerzeichen) behoben.
- **2026-10-19 06:54 · Alexander Sacharov** — critique Runde 3 (neuer Umfang, Diff f91dfe19 vollstaendig gelesen): keine Befunde. Stehen gelassen: kind-Regel in web_source.go Lanes liest dieselben Felder wie lane.Facts/core/board/announce.go (requires-human-exit oder requires-question = Lane einer Person) — Facts() zu nehmen spart nichts, es ist dieselbe Abbildung; Lane-Filter bietet weiter die Lanes der Tickets an (choices), jetzt in Board-Reihenfolge, damit auch Lanes anderer Boards filterbar bleiben.
- **2026-10-19 06:57 · Alexander Sacharov** — testing Runde 2: vet/race-Suite RC=0, dist unveraendert nach npm ci+build. Playwright: Lane-Farben bei zwei Ladevorgaengen identisch (backlog 240, todo 290, critique 215, testing 265, human/signoff 70, done 145); Lane-Filter ?lane=backlog setzt URL, Rahmen und Option in derselben Farbe wie das Badge; 8NHD8Q aus dem Logbuch oeffnet mit 'done in the logbook'; Scratch-Board: abgelehntes Archiv-Ticket zeigt 'Rejected: wontfix | Ada Lovelace, 10/2/26, 11:00 AM | nobody needs it' (Leerzeichen-Befund behoben), Karte nicht ziehbar.
- **2026-10-19 06:59 · Alexander Sacharov** — review (fresh context): DoD 1-9 vom Diff gedeckt, Tests gruen, dist synchron. Zwei kleine Nachbesserungen moeglich, nicht blockierend: Wettlauf beim schnellen Umschalten von 'All tickets ever' (spaete Antwort ueberschreibt neuere) und details open={expanded} klappt Hand-Aufgeklapptes beim Leeren des Filters zu.
- **2026-10-19 07:34 · Alexander Sacharov** — explain: Der Review findet keine Luecke gegen die DoD 1-9, nur zwei kleine UI-Unschoenheiten: Beim schnellen Umschalten von 'All tickets ever' kann eine spaete Antwort eine neuere ueberschreiben, und beim Leeren des Filters klappen Epics zu, die man von Hand geoeffnet hatte.
Warum: Beides ist nur Kosmetik und kein Datenverlust - die Liste stimmt nach dem naechsten Minuten-Reload wieder, und die Epics lassen sich neu aufklappen. Offen bleibt nur die Optik im Browser, die der Reviewer nicht selbst angesehen hat.
1) Ticket so annehmen und die zwei Unschoenheiten stehen lassen [accept]
2) Zurueckgeben: Wettlauf in ReviewPage.svelte load() beheben und Epic-Aufklappzustand beim Leeren des Filters erhalten [give back]
Empfehlung: 1 - Beide Punkte sind klein und nicht blockierend, der Diff erfuellt alle neun DoD-Punkte.
- **2026-10-19 07:38 · Alexander Sacharov** — decision: 2) Zurueckgeben: Wettlauf in ReviewPage.svelte load() beheben und Epic-Aufklappzustand beim Leeren des Filters erhalten
