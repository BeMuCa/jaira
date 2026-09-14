---
id: 01M2FQEEQN61ZE9AJ4Y4S1VM40
title: "Eine Karte zeigt bis zu drei Tag-Farben, nicht nur die des ersten Tags"
status: optimize
ready: true
creator: Alexander Sacharov
goal: "Auf einer Karte sind bis zu drei Tag-Farben gleichzeitig zu sehen: die ersten beiden Plaetze tragen die Farben der ersten beiden Tags des Tickets, der dritte Platz bleibt fuer die Sprint-Markierung reserviert und in diesem Ticket leer."
context: |-
  Ein Ticket kann viele Tags tragen, aber die Karte zeigt genau einen davon an.

  internal/tui/model.go:1341, cardColor(): die Kartenfarbe ist die Farbe von t.Tags[0]. Jeder weitere Tag ist auf dem Board unsichtbar - er steht im Ticket und in 'jaira show', faerbt aber nichts.

  Warum das stoert: beim Planen wird nach mehreren Achsen gleichzeitig gefiltert, zum Beispiel Thema plus Sprint. Mit einer Farbe je Karte muss man dafuer jedes Mal den Filter umstellen, statt die Verteilung einfach zu sehen.

  Der Platz dafuer ist schon da, die Karte muss nicht wachsen. internal/tui/view.go:488 cardHeight() gibt drei Zeilen zurueck, und renderCardBlock (internal/tui/view.go:523-545) faerbt eine Randspalte von einer Zelle Breite ueber diese drei Zeilen. Das sind drei eingefaerbte Zellen, die heute alle dieselbe Farbe tragen. Wie die drei Plaetze genau angeordnet werden, ist Sache der Plan-Lane, nicht dieses Tickets - die Bedingung ist nur, dass die Karte nicht hoeher wird.

  Entschieden am 2026-09-14 von Alex: die Begrenzung ist eine reine Anzeigegrenze. 'jaira tag' weist nichts zurueck und kein bestehendes Ticket wird ungueltig. Ein vierter Tag bleibt auf dem Ticket und bekommt nur keine Farbe. Das Board traegt 80-90 alte Tickets, von denen manche die Grenze schon heute ueberschreiten - eine erzwungene Grenze haette sie nachtraeglich falsch gemacht.

  Offen und absichtlich NICHT Teil dieses Tickets: was den dritten Platz fuellt. Heute waere das ein Tag namens sprint-xxxx, aber Alex erwaegt stattdessen ein eigenes Feld, weil unerledigte Arbeit im naechsten Sprint wieder auftauchen muss und ein Tag am Ticket haengt - und Tickets liegen auf ihren refs, muessten also einzeln nachgezogen werden. Diese Entscheidung steht aus. Bis sie faellt, bleibt Platz drei reserviert und leer; die Farbe eines Tickets aendert sich dadurch fuer niemanden.

  Die Farben selbst bleiben wie sie sind: eine Zeile 'name: <ansi256>' je Tag in .jaira/tags, handgeschrieben aenderbar. Ein Tag ohne Zeile dort ist weiterhin gueltig und kostet nichts.
definition-of-done: |-
  Eine Karte mit zwei Tags zeigt zwei unterscheidbare Farbfelder: Platz 1 traegt die Farbe des ersten Tags des Tickets, Platz 2 die des zweiten, in genau der Reihenfolge, in der sie im Ticket stehen.

  Platz 3 bleibt in diesem Ticket unbelegt und zeigt die Schattierung der Lane, so wie eine ungefaerbte Karte es heute tut.

  Ein Ticket mit vier Tags behaelt alle vier: 'jaira tag' nimmt den vierten an und gibt keinen Fehler, 'jaira show' listet ihn, und nur die Farbe fehlt ihm.

  Ein Tag ohne Zeile in .jaira/tags laesst seinen Platz in der Lane-Schattierung, und die drei Textzeilen der Karte stehen an derselben Stelle wie bei einer Karte ohne jeden Tag - nachgestellt an einer Karte mit einem gefaerbten und einem ungefaerbten Tag.

  cardHeight() gibt weiterhin 3 zurueck: die Karte wird durch diese Aenderung keine Zeile hoeher, nachgestellt an einer Lane mit mehr Karten als Platz.

  Eine Zeile in core/release/NOTES.md unter ## Unreleased.
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-09-14T10:29:46Z
updated-at: 2026-09-14T16:06:49Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-273100
claimed-at: 2026-09-14T15:48:04Z
outcome-what: "Die Randspalte der Karte hat jetzt drei Farb-Slots statt einer Kartenfarbe. cardColor in internal/tui/model.go ist durch cardColors ersetzt: gibt [3]cardSlot zurueck, Slot 1 = erster Tag, Slot 2 = zweiter Tag, Slot 3 immer ungefaerbt und fuer die Sprint-Markierung reserviert. renderCardBlock (internal/tui/view.go) waehlt die barParams je Zeile aus dem Slot dieser Zeile; ein ungefaerbter Slot faellt wie bisher auf die Lane-Schattierung bzw. die Selektionsfuellung zurueck. selectionFill und Glow bleiben an Slot 1. cardHeight gibt weiterhin 3 zurueck. Neue Tests in internal/tui/tagbox_test.go fuer beide Farben in Ticket-Reihenfolge, die immer ungefaerbte Zeile 3, den ungefaerbten zweiten Tag ohne Textversatz und die Vier-Tag-Karte; registryWith nimmt jetzt mehrere name/colour-Paare. Eine Zeile unter neuem ## Unreleased in core/release/NOTES.md."
outcome-why: "Ein Ticket kann viele Tags tragen, aber nur der erste faerbte die Karte; jeder weitere war auf dem Board unsichtbar. Beim Planen nach zwei Achsen (Thema plus Sprint) musste dafuer jedes Mal der Filter umgestellt werden. Die drei eingefaerbten Zellen der Randspalte waren schon da und trugen alle dieselbe Farbe - sie auf Slots aufzuteilen zeigt zwei Tags gleichzeitig, ohne dass die Karte eine Zeile waechst."
outcome-resolves: "Alle sechs DoD-Punkte sind mit Test bzw. Datei belegt: zwei unterscheidbare Farbfelder in Ticket-Reihenfolge (TestTwoTaggedCardShowsBothColoursInTicketOrder), Platz 3 in Lane-Schattierung (TestThirdRowAlwaysCarriesTheLaneShade), vier Tags bleiben erhalten ohne neue Validierung (TestFourTaggedCardRendersWithTheExtraTagsUncoloured, per CLI nachgestellt: jaira tag nimmt den vierten an, jaira show listet ihn), ungefaerbter Tag ohne Textversatz (TestUncolouredSecondTagFallsBackWithoutMovingTheText), cardHeight bleibt 3 (internal/tui/view.go:488 plus TestCardHeightIsTheThreeContentRows und TestACardHeavyWithFlagsStaysThreeRows), NOTES.md-Zeile unter ## Unreleased. go test ./... ist gruen."
review-summary: none
---

# Eine Karte zeigt bis zu drei Tag-Farben, nicht nur die des ersten Tags

## Definition of Done

- [x] Eine Karte mit zwei Tags zeigt zwei unterscheidbare Farbfelder: Platz 1 traegt die Farbe des ersten Tags des Tickets, Platz 2 die des zweiten, in genau der Reihenfolge, in der sie im Ticket stehen.
  proof: TestTwoTaggedCardShowsBothColoursInTicketOrder (internal/tui/tagbox_test.go)
- [x] Platz 3 bleibt in diesem Ticket unbelegt und zeigt die Schattierung der Lane, so wie eine ungefaerbte Karte es heute tut.
  proof: TestThirdRowAlwaysCarriesTheLaneShade (internal/tui/tagbox_test.go)
- [x] Ein Ticket mit vier Tags behaelt alle vier: 'jaira tag' nimmt den vierten an und gibt keinen Fehler, 'jaira show' listet ihn, und nur die Farbe fehlt ihm.
  proof: TestFourTaggedCardRendersWithTheExtraTagsUncoloured; CLI nachgestellt: jaira tag <id> ui backend docs ci exit 0, jaira show --json listet alle vier
- [x] Ein Tag ohne Zeile in .jaira/tags laesst seinen Platz in der Lane-Schattierung, und die drei Textzeilen der Karte stehen an derselben Stelle wie bei einer Karte ohne jeden Tag - nachgestellt an einer Karte mit einem gefaerbten und einem ungefaerbten Tag.
  proof: TestUncolouredSecondTagFallsBackWithoutMovingTheText (internal/tui/tagbox_test.go)
- [x] cardHeight() gibt weiterhin 3 zurueck: die Karte wird durch diese Aenderung keine Zeile hoeher, nachgestellt an einer Lane mit mehr Karten als Platz.
  proof: internal/tui/view.go:488 cardHeight returns 3; TestCardHeightIsTheThreeContentRows, TestACardHeavyWithFlagsStaysThreeRows, TestColumnDrawsEveryCardItCountsInFull
- [x] Eine Zeile in core/release/NOTES.md unter ## Unreleased.
  proof: core/release/NOTES.md:16

Platz 3 bleibt in diesem Ticket unbelegt und zeigt die Schattierung der Lane, so wie eine ungefaerbte Karte es heute tut.

Ein Ticket mit vier Tags behaelt alle vier: 'jaira tag' nimmt den vierten an und gibt keinen Fehler, 'jaira show' listet ihn, und nur die Farbe fehlt ihm.

Ein Tag ohne Zeile in .jaira/tags laesst seinen Platz in der Lane-Schattierung, und die drei Textzeilen der Karte stehen an derselben Stelle wie bei einer Karte ohne jeden Tag - nachgestellt an einer Karte mit einem gefaerbten und einem ungefaerbten Tag.

cardHeight() gibt weiterhin 3 zurueck: die Karte wird durch diese Aenderung keine Zeile hoeher, nachgestellt an einer Lane mit mehr Karten als Platz.

Eine Zeile in core/release/NOTES.md unter ## Unreleased.

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] internal/tui/tagbox_test.go durchgehen und aufschreiben, welche Tests heute welche Zusage der Randspalte pruefen (170-232, 370, 395)
- [x] cardColor (internal/tui/model.go:1341) durch cardColors ersetzen: gibt genau drei Slots zurueck, je Slot Farbe plus gefaerbt-ja-nein; Slot 1 = erster Tag, Slot 2 = zweiter Tag, Slot 3 immer ungefaerbt; ein Tag ohne Zeile in .jaira/tags ergibt einen ungefaerbten Slot, weitere Tags werden ignoriert
- [x] renderCardBlock (internal/tui/view.go:523) auf die Slots umstellen: barParams je Zeile aus dem Slot dieser Zeile, ungefaerbter Slot faellt wie heute auf bgParams zurueck; die Schleife bleibt eine Zeile je Inhaltszeile, cardHeight bleibt 3
- [x] selectionFill weiter mit Slot 1 speisen, damit die Fuellung der ausgewaehlten Karte und das Glow-Verhalten unveraendert bleiben
- [x] TestTaggedCardCarriesItsColourAsAFilledCell anpassen: bei einem Tag traegt nur noch Zeile 1 die Tag-Farbe, Zeile 2 und 3 die Lane-Schattierung
- [x] Test: Karte mit zwei gefaerbten Tags traegt in Zeile 1 die Farbe des ersten und in Zeile 2 die des zweiten Tags, in Ticket-Reihenfolge, und die beiden Farben sind verschieden
- [x] Test: Zeile 3 traegt bei jeder Karte die Lane-Schattierung, auch bei drei oder mehr Tags
- [x] Test: erster Tag gefaerbt, zweiter ohne Zeile in .jaira/tags - Zeile 2 faellt auf die Lane-Schattierung zurueck, und die drei Textzeilen stehen Spalte fuer Spalte dort, wo sie bei einer Karte ohne jeden Tag stehen
- [x] Test: Ticket mit vier Tags rendert ohne Fehler, Tag 3 und 4 faerben nichts; dazu per CLI nachstellen, dass jaira tag den vierten annimmt und jaira show ihn listet - ohne Validierung zu ergaenzen
- [x] Kommentare an cardColors und renderCardBlock neu schreiben: sie beschreiben heute eine einzige Kartenfarbe; der reservierte dritte Platz und der Grund dafuer gehoeren dazu
- [x] core/release/NOTES.md: Abschnitt ## Unreleased ueber ## 0.2.0 anlegen und eine Zeile schreiben, die sagt, was der Leser jetzt sieht - die ersten beiden Tags eines Tickets faerben die Randspalte in Ticket-Reihenfolge
- [x] go test ./... laufen lassen, besonders internal/tui; TestCardHeightIsTheThreeContentRows, TestACardHeavyWithFlagsStaysThreeRows und TestColumnDrawsEveryCardItCountsInFull muessen unveraendert gruen bleiben

## Progress
- **2026-09-14 15:52 · Alexander Sacharov** — Anordnung der drei Plaetze (Entscheidung der Plan-Lane): die drei Plaetze sind die drei Zeilen der bestehenden Randspalte, von oben nach unten. Platz 1 = Zeile 1 = erster Tag, Platz 2 = Zeile 2 = zweiter Tag, Platz 3 = Zeile 3 = reserviert und leer.

Verworfene Alternative: drei Zellen nebeneinander in einer Zeile. Das kostet zwei Spalten Titelbreite auf jeder Karte. Der Kommentar an renderCardBlock (internal/tui/view.go:523) haelt fest, dass der abgeschnittene Titel die meistbeklagte Flaeche des Boards war und der Rahmen genau deswegen gefallen ist - Breite wieder wegzunehmen wuerde das zurueckdrehen. Die Zeilenvariante kostet nichts: die drei Zellen werden heute schon gemalt, nur alle in derselben Farbe.
- **2026-09-14 15:52 · Alexander Sacharov** — Folge fuer Karten mit genau einem Tag - bitte vor der Umsetzung nicht uebersehen: heute traegt die Randspalte einer einfach getaggten Karte ueber alle drei Zeilen die Tag-Farbe. Nach der Aenderung nur noch Zeile 1; Zeile 2 und 3 zeigen die Lane-Schattierung. Das folgt zwingend aus der DoD (Platz 2 gehoert dem zweiten Tag, Platz 3 bleibt leer), ist aber eine sichtbare Aenderung fuer die grosse Mehrheit der Karten.

Konkret bricht das den bestehenden Test TestTaggedCardCarriesItsColourAsAFilledCell (internal/tui/tagbox_test.go:207): er verlangt, dass JEDE Zeile der Karte die Farbe 83 traegt. Dieser Test wird angepasst, nicht geloescht - der Nachweis bleibt, nur auf Zeile 1 verengt.

Falls Alex das nicht will (durchgehender Balken bei einem Tag, Aufteilung erst ab zwei Tags), ist das eine Entscheidung fuer die human-Lane, kein Freihand-Umbau in der Implementierung.
- **2026-09-14 15:52 · Alexander Sacharov** — Zwei kleinere Punkte, die der Plan mitnimmt:

1. selectionFill/glow (internal/tui/glow.go:88): die Fuellung der ausgewaehlten Karte wird heute aus derselben einen Tag-Farbe gemischt. Mit mehreren Farben bleibt es bei Slot 1 (erster gefaerbter Tag). Kein Mischen ueber zwei Tags - das waere eine eigene Gestaltungsfrage und steht nicht in der DoD.

2. core/release/NOTES.md hat aktuell KEINEN Abschnitt ## Unreleased; die oberste Ueberschrift ist ## 0.2.0. Der Abschnitt muss angelegt werden, ueber 0.2.0, nicht in 0.2.0 hineingeschrieben.

3. DoD-Punkt 3 (vierter Tag bleibt erhalten) verlangt keine Codeaenderung: cardColor ist der einzige Aufrufer der Registry-Farbe auf der Karte (Suche ueber internal/tui), jaira tag und jaira show fassen die Grenze nicht an. Der Punkt wird nachgestellt, nicht gebaut - und ausdruecklich keine Validierung ergaenzt.
- **2026-09-14 16:00 · Alexander Sacharov** — TestColumnDrawsEveryCardItCountsInFull musste doch angepasst werden (der Plan erwartete ihn unveraendert): er zaehlte '48;5;83m' und rechnete mit shown*3, weil die Tag-Farbe frueher alle drei Zeilen fuellte. Diese Zusage gibt es nicht mehr. Die Absicht des Tests - keine Karte wird halb gezeichnet - bleibt, jetzt ueber drei Zaehlungen: Farbe des ersten Tags = Zeile 1, Farbe des zweiten = Zeile 2, und die Flag-Zeile '○ spec' = Zeile 3. Jede halb gezeichnete Karte faellt bei einer der drei durch.
- **2026-09-14 16:00 · Alexander Sacharov** — registryWith in internal/tui/tagbox_test.go nimmt jetzt beliebig viele name/colour-Paare (vorher genau eines). Variadisch als '...any', damit alle bestehenden Aufrufe registryWith(t, "ui", 83) unveraendert weiterkompilieren - sonst haetten ~20 Aufrufstellen angefasst werden muessen, ohne dass sich an ihnen etwas aendert.
- **2026-09-14 16:00 · Alexander Sacharov** — Bewusst NICHT geaendert: selectionFill und das Glow-Verhalten speisen sich weiter allein aus Slot 1. Die ausgewaehlte Karte wird als Ganzes in einer Farbe gefuellt; zwei Farben dort waeren eine zweite, eigene Entscheidung und gehoeren nicht in dieses Ticket.
- **2026-09-14 16:02 · Alexander Sacharov** — Commit dieser Lane: 1276d7b (feat(S1VM40): show the first two tag colours down a card's bar). Fuer den move als --commits 1276d7b mitgeben.
- **2026-09-14 16:06 · Alexander Sacharov** — critique (Durchgang 1): keine Befunde, review-summary=none. Gelesen: der volle Diff gegen master (1276d7b, 2fc0140) in internal/tui/model.go, internal/tui/view.go, internal/tui/tagbox_test.go und core/release/NOTES.md.

Was geprueft wurde und stand hielt:
- Die Slots liegen auf den drei Zeilen der bestehenden Randspalte. cardHeight (internal/tui/view.go:488) gibt weiterhin 3 zurueck, die Schleife in renderCardBlock bleibt eine Zeile je Inhaltszeile - die Karte waechst nicht, und der Titel verliert keine Spalte (die verworfene Nebeneinander-Variante haette zwei gekostet).
- Anzeigegrenze statt Validierung: cardColors (internal/tui/model.go:1361) ignoriert Tags ab dem dritten stillschweigend; in core/tag und in den tag-/show-Kommandos wurde nichts angefasst. Alex' Entscheidung ist damit eingehalten, kein Altticket wird nachtraeglich ungueltig.
- Slot 3 bleibt hart ungefaerbt (Schleifengrenze cardSlots-1), also kann kein dritter Tag in den reservierten Sprint-Platz rutschen - genau das prueft TestThirdSlotStaysUncolouredHoweverManyTags.
- cardColor hatte ausserhalb der Tests keinen zweiten Aufrufer (grep ueber internal/); selectionFill/Glow speisen sich weiter aus Slot 1, wie im Plan entschieden.

Zu den beiden angepassten Tests - die Anpassung erhaelt, was sie geschuetzt haben:
- TestTaggedCardCarriesItsColourAsAFilledCell -> ...InTheFirstRowOnly: die Zusage war 'die Farbe fuellt eine ganze Zelle, nicht ein halbes Rahmenglyph, und jede Zeile traegt eine Balkenzelle'. Beides steht noch da: Zeile 1 wird auf 48;5;83 geprueft, Zeile 2 und 3 auf das Oeffnen mit der Lane-Schattierung. Dass der Text trotz ungefaerbter Zellen an derselben Spalte steht - der eigentliche Sinn der alten Jede-Zeile-Pruefung - deckt jetzt TestUncolouredSecondTagFallsBackWithoutMovingTheText per Vergleich mit einer ungetaggten Karte ab.
- TestColumnDrawsEveryCardItCountsInFull: die Zusage war 'keine Karte wird halb gezeichnet'. Die alte Zaehlung shown*3 funktionierte nur, weil eine Farbe alle drei Zeilen fuellte. Der Ersatz zaehlt Zeile 1 (48;5;83), Zeile 2 (48;5;45) und Zeile 3 ueber die Flag-Zeile '○ spec' - die Testtickets stehen auf ready=false, die Flag-Zeile ist damit auf jeder Karte garantiert und liegt in renderCard (internal/tui/view.go:706-710) tatsaechlich auf der dritten Zeile. Eine auf ein oder zwei Zeilen abgeschnittene Karte faellt bei einer der drei Zaehlungen durch. Deckung erhalten.

Bewusst nicht als Befund erhoben:
- registryWith (internal/tui/tagbox_test.go:17) nimmt 'pairs ...any' mit Laufzeit-Typpruefung statt eines getippten Paar-Typs. Vertretbarer Tausch: jede getippte Variante haette alle ~20 bestehenden Aufrufstellen angefasst, ohne dass sich an ihnen etwas aendert; ein Fehlgriff faellt mit klarem t.Fatalf sofort auf.
- 'i < len(slots)' in renderCardBlock ist bei einem Array fester Laenge und drei Inhaltszeilen nie falsch. Eine Zeile toter Schutz, Sache der optimize-Lane, nicht der critique.
- Veraltete Kommentare zur alten Rahmen-Aera in internal/tui/view.go:443 und internal/tui/model.go:1264 ('Box') sind aelter als dieser Diff und gehoeren nicht in dieses Ticket.

go build ./..., go vet ./internal/tui und go test ./... sind gruen, gofmt meldet nichts.
