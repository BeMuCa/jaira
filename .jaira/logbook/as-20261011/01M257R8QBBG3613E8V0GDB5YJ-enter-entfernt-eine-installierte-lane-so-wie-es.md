---
id: 01M257R8QBBG3613E8V0GDB5YJ
title: "Enter entfernt eine installierte Lane, so wie es sie hinzufuegt"
status: done
ready: true
creator: BeMuCa
assignee: Alexander Sacharov
goal: "Enter ist im Lane-Screen der eine Schalter: auf einer noch nicht installierten Lane fuegt es sie hinzu, auf einer installierten nimmt es sie wieder weg"
context: |-
  Berk am 09.09.: 'man sollte eine lane die hinzugefuegt wurde mit enter auch wieder entfernen'.
  WICHTIG - Entfernen gibt es schon, auf x. internal/tui/lanes.go:284-287 ruft startRemove, das einen Ja/Nein-Dialog aufmacht (lanes.go:347); die Fusszeile nennt 'x remove'. Es fehlt also keine Faehigkeit, es fehlt die Taste, die Berk erwartet.
  Der Konflikt, der vor dem Bauen zu entscheiden ist: enter ist heute schon belegt (lanes.go:288-293). Auf der Plus-Spalte oeffnet es den Katalog, auf einer noch nicht installierten Lane fuegt es sie hinzu. Beides ist 'enter fuegt hinzu'.
  Berks Wunsch macht daraus einen Umschalter: enter auf einer installierten Lane = entfernen. Das ist in sich stimmig, aber enter wird damit auf einer Spalte destruktiv, und enter ist die Taste, die man blind drueckt.
  Was dagegen absichert: der Ja/Nein-Dialog existiert schon und zeichnet 'no' zuerst, damit blindes Enter-Haemmern nichts loescht (lanes.go:810-820). Diese Vorsicht muss erhalten bleiben - genau sie ist der Grund, warum enter hier ueberhaupt vertretbar ist.
  x soll bleiben. Zwei Wege zur selben Handlung sind hier kein Widerspruch, sondern der Schutz fuer den, der x schon kennt.
definition-of-done: "Enter auf einer installierten Lane startet dieselbe Entfernung wie x, samt Ja/Nein-Dialog mit 'no' als Vorauswahl; enter auf der Plus-Spalte und auf einer nicht installierten Lane fuegt weiter hinzu; x funktioniert unveraendert; die Fusszeile nennt beide; Tests decken beide Enter-Faelle; go test ./... -race gruen"
tags:
  - tui
blocked-by: []
commits:
  - 6bd21c73399d119c3edd146759e47737c2828ef0
  - 9dafa85f5202e1b6b71bcf850144df8ddf6f6d21
  - b4fbf11768a471c4932185daf2a13c2f46f84be6
created-at: 2026-09-10T08:43:06Z
updated-at: 2026-10-11T20:38:10Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-925641
claimed-at: 2026-10-11T20:04:07Z
outcome-what: "Review: enter-Umschalter im Lane-Screen gegen die DoD geprueft"
outcome-why: "Ein zweites Modell bestaetigt den Diff vor der Abnahme"
outcome-resolves: "DoD 1 erfuellt; Tests gruen mit sauberem HOME"
review-summary: "Im Lane-Screen ruft enter auf einer installierten Lane jetzt startRemove() auf (internal/tui/lanes.go, case \"enter\", dritter Zweig) - denselben Weg wie x: es oeffnet confirmRemove{id} ohne path, also Board-Entfernung ueber lane.Remove, Vorauswahl no. Enter auf '+' (Katalog) und auf einer nicht installierten Lane (addAvailable) bleibt unveraendert; der Datei-Loesch-Zweig fuer verfuegbare Lanes ist nur ueber x erreichbar, weil enter dort vorher in addAvailable endet. Fusszeile in beiden Varianten um 'enter add/remove' ergaenzt, 'x remove' bleibt. Tests: neuer Helper laneScreenWithOneAvailable ersetzt drei Kopien desselben Setups; TestLaneScreenEnterTogglesAddAndRemove prueft hinzufuegen -> enter fragt (no) -> enter,enter entfernt nichts -> enter,right,enter entfernt; TestLaneScreenFooterNamesEnterAndX prueft die Fusszeile. NOTES.md-Zeile unter Unreleased, einzeilig."
review-gaps: "Keine Luecke gegen die DoD. Klein: (1) Der Fusszeilen-Test prueft nur die Variante ohne geteilte Lanes; die zweite help-Liste (len(ls.shared)>0) ist ungetestet, aber identisch geaendert. (2) Der umgebrochene Kommentar bei 'if ls.confirm != nil' hat jetzt eine ueberlange Zeile - kosmetisch. (3) go test ./... -race scheitert auf dieser Maschine mit dem echten HOME (internal/cli Dod*/Lanes*, internal/tui Archive/Delete: 'unrecognized lane backlog') - identisch auf master 15f8057, also nicht von diesem Ticket; mit leerem HOME ist alles gruen. Test-Isolation gegen ~/.jaira fehlt - eigenes Thema."
test-verdict: "pass: go build ./... RC=0; go test ./... -race -count=1 RC=0 (isoliertes HOME); 27 TestLaneScreen* gruen, darunter EnterTogglesAddAndRemove (enter fuegt hinzu, enter auf installierter Lane oeffnet Dialog mit yes=false, enter-enter loescht nichts, enter-right-enter entfernt), PlusColumn (enter auf + oeffnet Katalog), X*-Tests unveraendert gruen, FooterNamesEnterAndX. Interaktiver TUI-Lauf nicht gemacht - key()-Sequenzen in den Tests sind derselbe Pfad"
review-verdict: "Der Diff erfuellt die DoD: enter schaltet hinzufuegen/entfernen, Entfernen nur ueber den Ja/Nein-Dialog mit no als Vorauswahl, x unveraendert, Fusszeile nennt beide, beide Enter-Faelle getestet, go test ./... -race gruen (mit sauberem HOME). Keine Defekte gefunden."
review-check: "1. go run ./cmd/jaira in einem Repo mit .jaira/  2. S druecken, zum Lanes-Screen gehen  3. Fusszeile zeigt 'enter add/remove' und 'x remove'  4. Cursor auf eine installierte Lane ohne Tickets (z. B. eine leere), enter: Ja/Nein-Frage erscheint, 'no' ist markiert  5. enter: nichts wird entfernt, Lane bleibt  6. enter, dann rechts (yes), enter: Lane verschwindet aus der Reihe  7. Cursor auf die nun nicht installierte (gedimmte) Lane, enter: sie wird ohne Frage wieder hinzugefuegt  8. Cursor auf '+', enter: der Katalog oeffnet sich  9. x auf einer installierten Lane: dieselbe Frage wie in Schritt 4. Automatisch: go test ./internal/tui -run 'TestLaneScreenEnterToggles|TestLaneScreenFooterNamesEnterAndX' -race"
---

# Enter entfernt eine installierte Lane, so wie es sie hinzufuegt

## Definition of Done

- [x] Enter auf einer installierten Lane startet dieselbe Entfernung wie x, samt Ja/Nein-Dialog mit 'no' als Vorauswahl; enter auf der Plus-Spalte und auf einer nicht installierten Lane fuegt weiter hinzu; x funktioniert unveraendert; die Fusszeile nennt beide; Tests decken beide Enter-Faelle; go test ./... -race gruen
  proof: internal/tui/lanes.go:289-302 (enter -> startRemove auf installierter Lane), Footer lanes.go:860-862; Tests TestLaneScreenEnterTogglesAddAndRemove, TestLaneScreenFooterNamesEnterAndX, TestLaneScreenPlusColumnOpensCatalogueListingOnlyMissingLanes; go test ./... -race exit 0 (isoliertes HOME)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-11 19:44 · Alexander Sacharov** — Code aus dem A/B-Versuch vom 11.10.2026 uebernommen (Arm A, Lane-Pipeline; blinder Richter knapp fuer diese Variante; Review: approve). Offen laut Review: Doc-Kommentare in internal/tui/lanes.go um startRemove nennen teils nur x. Weiter im Worker-Modus; Alex hat das Ticket von BeMuCa uebernommen, um es zum Release zu bringen.
- **2026-10-11 20:04 · Alexander Sacharov** — parent session: jaira-61
- **2026-10-11 20:04 · Alexander Sacharov** — Keine offenen Entscheidungen: Ticket-Kontext legt alles fest (enter=Umschalter, x bleibt, Dialog mit no zuerst). Code liegt in 9dafa85; offen nur der Review-Punkt zu Doc-Kommentaren um startRemove.
- **2026-10-11 20:11 · Alexander Sacharov** — go test ./... -race ist auf dieser Maschine nur mit isoliertem absolutem HOME gruen (Scratchpad-HOME, exit 0). Mit echtem HOME scheitern ~30 Tests in internal/cli und TestArchiveFromTheBoard (nil-Panic), Meldung 'unrecognized lane backlog' - Ursache ist der maschinenlokale ~/.jaira/default-board.md, nicht diese Aenderung. Tests isolieren HOME nicht; eigenes Thema, hier nicht angefasst.
- **2026-10-11 20:11 · Alexander Sacharov** — critique: nichts. enter ruft denselben startRemove wie x (keine zweite Entfernungslogik); selectedAvailable wird vorher abgefangen, also erreicht enter den Datei-Loesch-Zweig nie. laneScreenWithOneAvailable ersetzt drei identische Setups und traegt den neuen Test - kein Scope-Kriechen.
- **2026-10-11 20:28 · Alexander Sacharov** — review: Diff gegen DoD geprueft, keine Defekte. go test ./... -race scheitert hier mit echtem HOME auch auf master 15f8057 ('unrecognized lane backlog' aus ~/.jaira) - mit leerem HOME komplett gruen; Test-Isolation gegen ~/.jaira ist ein eigenes Thema, nicht dieses Ticket.
