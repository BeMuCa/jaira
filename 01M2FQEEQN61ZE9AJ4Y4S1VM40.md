---
id: 01M2FQEEQN61ZE9AJ4Y4S1VM40
title: "Eine Karte zeigt bis zu drei Tag-Farben, nicht nur die des ersten Tags"
status: in-progress
ready: true
creator: Alexander Sacharov
goal: "Auf einer Karte sind bis zu drei Tag-Farben gleichzeitig zu sehen: die drei Plaetze der linken Randspalte tragen die Farben der ersten drei Tags des Tickets, in Ticket-Reihenfolge."
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

  Platz 3 traegt die Farbe des dritten Tags: eine Karte mit drei gefaerbten Tags zeigt drei unterscheidbare Farbfelder in Ticket-Reihenfolge. Eine Karte mit einem Tag faerbt nur Platz 1, eine mit zwei nur Platz 1 und 2; unbelegte Plaetze zeigen weiter die Schattierung der Lane.

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
updated-at: 2026-09-14T18:43:09Z
assignee: "Alexander Sacharov"
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-521068
claimed-at: 2026-09-14T18:34:27Z
outcome-what: "Drei Critique-Befunde der dritten Runde abgearbeitet, ohne cardColors selbst anzufassen. (1) internal/tui/tagbox_test.go: TestCardColorsAreTheFirstThreeTagsInTicketOrder ersatzlos geloescht - nach dem Umschreiben der Vorrunde war er Zeile fuer Zeile in TestTagsPastTheThirdColourNoSlot enthalten, der zusaetzlich den vierten Tag traegt. (2) internal/tui/model.go:1351: der Kommentar an cardColors sagt nicht mehr im Praesens, die Sprint-Markierung stehe am rechten Kartenrand, sondern 'the sprint marker moved off this bar (ticket 0YGWXQ)'. (3) Ticketdatei: die proof-Zeilen von DoD-Punkt 2 und 3 zeigen wieder auf existierende Tests - Punkt 3 auf den umbenannten TestFourTaggedCardRendersWithTheExtraTagUncoloured (internal/tui/tagbox_test.go:593), Punkt 2 ohne den geloeschten Test."
outcome-why: "Ein doppelter Test kostet Laufzeit und laesst beim Aendern offen, welcher der beiden die Zusage traegt. Ein Kommentar im Praesens ueber Code, den es nicht gibt, schickt den naechsten Leser an den rechten Kartenrand suchen und nennt nicht einmal das Ticket, in dem das entstehen soll. Und ein proof, der auf einen umbenannten Test zeigt, ist kein Beleg mehr - die Sign-off-Lane liest genau diese Zeilen."
outcome-resolves: "Kein DoD-Punkt aendert sein Verhalten: die Belege zeigen jetzt auf Tests, die es gibt. Punkt 2 (dritter Tag faerbt Platz 3, weniger Tags lassen die Plaetze darunter leer) haelt ueber TestThreeTaggedCardShowsAllThreeColoursInTicketOrder und TestSlotsBelowTheLastTagStayUncoloured, Punkt 3 (vier Tags bleiben erhalten, nur die Farbe fehlt) ueber TestFourTaggedCardRendersWithTheExtraTagUncoloured. go test ./... Exit 0, go vet ./... und gofmt -l ohne Ausgabe."
review-summary: ".jaira/tickets/01M2FQEEQN61ZE9AJ4Y4S1VM40-eine-karte-zeigt-bis-zu-drei-tag-farben-nicht-nu.md:25 das Frontmatter-Feld definition-of-done sagt weiter 'Platz 3 bleibt in diesem Ticket unbelegt und zeigt die Schattierung der Lane, so wie eine ungefaerbte Karte es heute tut' - das Gegenteil von goal (Zeile 7) und von DoD-Punkt 2 (Zeile 63), die genau dieser Diff umgeschrieben hat. Das Feld ist nicht Zierde: 'jaira show --for-lane' liefert es jeder folgenden Lane als input (diese critique-Runde hat es so bekommen), also legen review und signoff die zurueckgezogene Zusage neben den Code, der das Gegenteil tut. Stattdessen: jaira set S1VM40 definition-of-done=<...> und Absatz 2 im Wortlaut von DoD-Punkt 2 ('Platz 3 traegt die Farbe des dritten Tags ...'); die uebrigen fuenf Absaetze bleiben unveraendert."
review-gaps: "Entfernt: die tote Schranke 'i < len(slots)' in renderCardBlock (internal/tui/view.go:552). Die Schleife laeuft ueber die Zeilen von renderCard, und das sind immer genau drei — cardHeight() gibt 3 zurueck, cardSlots ist 3, und renderCard baut Titel/Meta/Flags fest als drei Zeilen. Die Bedingung konnte also nie falsch sein und tat so, als gaebe es Karten mit mehr Zeilen als Slots. Bewusst stehen gelassen: der reservierte dritte Slot (Alex' Entscheidung, Sprint-Markierung), das Feld cardSlot.coloured (Palettenfarbe 0 ist gueltig, ein blosser int kann 'keine Farbe' nicht ausdruecken), der nil-Check auf m.tags (tag.Registry.Colour laeuft nicht auf nil), das variadische 'pairs ...any' in registryWith und die beiden alten Box-Kommentare — alles drei hat schon die Critique-Lane geprueft. Keine Duplikate: cardColors ist die einzige Stelle, die Tag zu Kartenfarbe macht; cardColor hat keine Aufrufer mehr hinterlassen. Das Muster '5;'+strconv.Itoa(...) steht mehrfach in glow.go und view.go, ist aber aelter als dieser Diff und nicht sein Problem. Kosten: cardColors wird einmal pro Karte aufgerufen, nicht pro Zeile — pro Zeile bleibt nur ein Array-Zugriff. Die zwei Tests TestThirdSlotStaysUncolouredHoweverManyTags und TestThirdRowAlwaysCarriesTheLaneShade ueberschneiden sich thematisch, pruefen aber verschiedene Schichten (Modell und Rendering) und bleiben beide. go test ./... gruen."
test-verdict: "pass — Durchgang 2, unabhaengig nachgeprueft. Der Befund aus Durchgang 1 ist weg: mit einem echten Board auf Platte (ticket.Store, Titel/Assignee/updated-by/executed-by je mit Zeilenumbruch, auch \\r\\n) rendert m.render() ohne Panik, die Karte bleibt dreizeilig und der Umbruch wird zum Leerzeichen. Gegenprobe ohne fremde Tests: derselbe Fall gegen internal/tui/view.go aus f2b0077 (per go test -overlay, Arbeitsbaum unveraendert) paniert weiterhin mit 'index out of range [3] with length 3' — der Nachweis haengt also am Fix, nicht am Testtext. Keine neue Stoerung: die gerenderten Karten ohne Umbruch (fuenf Tickets x selected/alt) sind byteweise identisch zu denen vor dem Fix, oneLine ist fuer Werte ohne \\r\\n die Identitaet. renderCardBlock ist der einzige Aufrufer von renderCard (internal/tui/view.go:545). Alle sechs DoD-Punkte mit eigenen Pruefungen bestaetigt: 1+2 Slot 1 = 5;83, Slot 2 = 5;45, umgekehrte Tag-Reihenfolge kehrt die Farben um, Slot 3 traegt die Lane-Schattierung auch bei drei Tags; 3 auf einem Wegwerf-Board per CLI: 'jaira tag <id> ui backend docs ci' Exit 0, 'jaira show --json' listet alle vier; 4 ungefaerbter zweiter Tag faellt auf die Lane-Schattierung, die drei Textzeilen sind zeichengleich mit einer Karte ohne Tags; 5 cardHeight()==3 und in einer Lane mit 30 Karten und Platz fuer wenige traegt jede gezeigte Karte ihre Flag-Zeile — keine halb gezeichnete Karte; 6 core/release/NOTES.md, zwei einzeilige '- '-Zeilen unter '## Unreleased' ueber '## 0.2.0'. Suiten: go test ./... Exit 0, go test ./... -race Exit 0, go vet ./... ohne Ausgabe, gofmt -l ohne Ausgabe. Kosmetischer Restpunkt, kein Fehler: ein per CLI geschriebenes \\r\\n kommt als ' \\n' aus dem YAML zurueck und faltet damit zu zwei Leerzeichen; die NOTES-Zeile nennt nur den Titel, obwohl der Fix auch Assignee, updated-by und executed-by deckt."
question: ""
---

# Eine Karte zeigt bis zu drei Tag-Farben, nicht nur die des ersten Tags

## Definition of Done

- [x] Eine Karte mit zwei Tags zeigt zwei unterscheidbare Farbfelder: Platz 1 traegt die Farbe des ersten Tags des Tickets, Platz 2 die des zweiten, in genau der Reihenfolge, in der sie im Ticket stehen.
  proof: TestTwoTaggedCardShowsBothColoursInTicketOrder (internal/tui/tagbox_test.go)
- [x] Platz 3 traegt die Farbe des dritten Tags: eine Karte mit drei gefaerbten Tags zeigt drei unterscheidbare Farbfelder in Ticket-Reihenfolge. Eine Karte mit einem Tag faerbt nur Platz 1, eine mit zwei nur Platz 1 und 2; unbelegte Plaetze zeigen weiter die Schattierung der Lane.
  proof: TestThreeTaggedCardShowsAllThreeColoursInTicketOrder, TestSlotsBelowTheLastTagStayUncoloured (internal/tui/tagbox_test.go)
- [x] Ein Ticket mit vier Tags behaelt alle vier: 'jaira tag' nimmt den vierten an und gibt keinen Fehler, 'jaira show' listet ihn, und nur die Farbe fehlt ihm.
  proof: TestFourTaggedCardRendersWithTheExtraTagUncoloured (internal/tui/tagbox_test.go:593); CLI nachgestellt: jaira tag <id> ui backend docs ci exit 0, jaira show --json listet alle vier
- [x] Ein Tag ohne Zeile in .jaira/tags laesst seinen Platz in der Lane-Schattierung, und die drei Textzeilen der Karte stehen an derselben Stelle wie bei einer Karte ohne jeden Tag - nachgestellt an einer Karte mit einem gefaerbten und einem ungefaerbten Tag.
  proof: TestUncolouredSecondTagFallsBackWithoutMovingTheText (internal/tui/tagbox_test.go)
- [x] cardHeight() gibt weiterhin 3 zurueck: die Karte wird durch diese Aenderung keine Zeile hoeher, nachgestellt an einer Lane mit mehr Karten als Platz.
  proof: internal/tui/view.go cardHeight returns 3; TestCardHeightIsTheThreeContentRows, TestACardHeavyWithFlagsStaysThreeRows, TestColumnDrawsEveryCardItCountsInFull, TestRenderCardAlwaysReturnsAsManyLinesAsThereAreSlots, TestACardWhoseFieldsCarryNewlinesStaysThreeRows
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
- **2026-09-14 16:09 · Alexander Sacharov** — optimize: nur die tote Schranke 'i < len(slots)' entfernt, die die Critique-Lane uebergeben hatte. Der Bezug ist jetzt direkt: die Schleife indiziert slots[i] ueber die drei Zeilen von renderCard. Das haelt, solange cardHeight()==cardSlots gilt — genau das sagt der Kommentar an cardSlots bereits. Sonst nichts gefunden: keine Duplikate, kein verwaister Code, keine Arbeit in der Schleife, die sich heben liesse.
- **2026-09-14 16:15 · Alexander Sacharov** — testing (Durchgang 1): FEHLGESCHLAGEN — eine neue Absturzquelle im Diff.

Suite: go test ./... gruen, go test ./... -race gruen (Exit 0), go vet ./... gruen, gofmt -l meldet nichts. 0 Fehler, 0 Skips in der regulaeren Suite.

Der Befund liegt nicht in der Suite, sondern in einem Fall, den sie nicht abdeckt:

internal/tui/view.go:553 — 'panic: runtime error: index out of range [3] with length 3' in renderCardBlock. Die Schleife laeuft ueber die Zeilen von renderCard und indiziert slots[i] mit genau drei Slots. renderCard liefert normalerweise drei Zeilen, aber der Titel wird ungeprueft eingesetzt (internal/tui/view.go:617 truncate(t.Title, w-1)); enthaelt der Titel einen Zeilenumbruch, sind es vier Zeilen und die Karte laesst das Board abstuerzen.

Der Fall ist erreichbar, nicht theoretisch: 'jaira create "Zeile eins<NL>Zeile zwei"' wird angenommen und schreibt title als YAML-Blockskalar (|-); Decode liest den Umbruch zurueck (nachgestellt ueber core/ticket.ParseDoc/Decode). Auch handgeschriebenes 'title: "a\\nb"' genuegt — das Format ist ausdruecklich handeditierbar.

Regression, nicht geerbt: auf 1c6be9d (master) und auf 1276d7b (Implementierungs-Commit dieses Tickets) rendert derselbe Fall ohne Absturz; erst f2b0077 (optimize: Entfernen von 'i < len(slots)') macht daraus einen Panic. Nachgestellt ueber 'go test -overlay' mit einer Testdatei ausserhalb des Repos — im Arbeitsbaum wurde nichts geaendert.

DoD-Belege sonst nachgeprueft und alle haltbar: TestTwoTaggedCardShowsBothColoursInTicketOrder, TestThirdRowAlwaysCarriesTheLaneShade, TestUncolouredSecondTagFallsBackWithoutMovingTheText, TestFourTaggedCardRendersWithTheExtraTagsUncoloured, TestACardHeavyWithFlagsStaysThreeRows und TestColumnDrawsEveryCardItCountsInFull laufen gruen. cardHeight gibt weiterhin 3 zurueck (internal/tui/view.go:488-490). Punkt 3 per CLI auf einem Wegwerf-Board geprueft: 'jaira tag <id> ui backend docs ci' Exit 0, 'jaira show --json' listet alle vier Tags, nur Farben fehlen keinem — .jaira/tags hat sogar fuer alle vier eine Zeile. Punkt 6: core/release/NOTES.md:16, eine Zeile unter neuem '## Unreleased'.

Nicht abgedeckt: es gibt keinen Test, der renderCardBlock mit einer Karte konfrontiert, deren Inhalt mehr oder weniger Zeilen hat als es Slots gibt. Genau dort faellt der Diff um.
- **2026-09-14 16:22 · Alexander Sacharov** — in-progress (Fixrunde nach testing): der Panic aus internal/tui/view.go:553 ist behoben, aber NICHT durch das Zurueckholen der Schranke 'i < len(slots)'.

Warum nicht die Schranke: sie haette nur den Absturz weggenommen und den eigentlichen Schaden stehen lassen. Eine Karte mit vierzeiligem Inhalt ist vier Zeilen hoch, waehrend cardHeight() 3 zurueckgibt und cardsInBudget (internal/tui/view.go:497) mit 3 rechnet. Die vierte Zeile wird dann ueber die naechste Karte gemalt und die Lane rechnet ihr Zeilenbudget falsch - eine stille Layoutstoerung statt eines sichtbaren Absturzes. Der verlorene Balkenfarbe waere dabei noch das kleinste Problem.

Stattdessen an der Quelle: neue Hilfsfunktion oneLine (internal/tui/view.go, bei truncate) faltet Zeilenumbrueche zu Leerzeichen, angewandt in renderCard auf alle drei Zeilen - Titel, meta (traegt Assignee) und die Flag-Zeile (traegt UpdatedBy und ExecutedBy). Damit haelt die Zusage 'renderCard gibt genau cardHeight Zeilen zurueck' fuer jeden Wert, den ein Ticket tragen kann, nicht nur fuer den Titel. \r\n wird mitgefaltet, ein einzelnes \r ebenso - sonst haette das CR die Zeile im Terminal zurueckgesetzt.

Die tote Schranke bleibt entfernt: der Befund der optimize-Lane war richtig, sie war nur an der falschen Stelle geheilt worden.

Gegenprobe gemacht: mit zurueckgedrehtem oneLine faellt TestACardWhoseFieldsCarryNewlinesStaysThreeRows/title mit genau 'index out of range [3] with length 3'; mit oneLine gruen.

NOTES.md: eine zweite Zeile unter ## Unreleased. Gegen den letzten Release ist das sichtbar - auf master stuerzt derselbe Titel zwar nicht ab, malt aber eine vierzeilige Karte.
- **2026-09-14 16:23 · Alexander Sacharov** — Commits dieses Tickets fuer den move: 1276d7b (Implementierung), f2b0077 (optimize), 4e96b33 (Fixrunde nach testing).
- **2026-09-14 16:29 · Alexander Sacharov** — testing (Durchgang 2): BESTANDEN.

Suiten: go test ./... Exit 0, go test ./... -race Exit 0, go vet ./... ohne Ausgabe, gofmt -l ohne Ausgabe. 0 Fehler, 0 Skips.

Der Absturz aus Durchgang 1 ist unabhaengig nachgestellt und weg. Nicht ueber die neuen Tests geprueft, sondern mit eigenen, per 'go test -overlay' eingehaengten Testdateien (im Arbeitsbaum wurde nichts geaendert):
- Echtes Board auf Platte: Tickets ueber ticket.Store angelegt mit Umbruch in title, assignee, updated-by und executed-by, davon einer mit \r\n, danach neu von Platte gelesen. Der Umbruch ueberlebt den Round-Trip (nachgeprueft), m.render() laeuft durch, der Titel steht als eine Zeile auf der Karte und kein \r landet im Output.
- renderCardBlock direkt: sieben Karten mit Umbruechen, darunter eine aus lauter Leerzeilen — jede genau drei Zeilen.
- Gegenprobe: dieselbe Testdatei gegen internal/tui/view.go aus f2b0077 paniert weiterhin mit 'index out of range [3] with length 3'. Der Fix traegt den Nachweis, nicht die Testformulierung.

Keine neue Stoerung gefunden. Fuenf Karten ohne Umbruch (je selected/alt in allen vier Kombinationen) wurden mit und ohne Fix gerendert und verglichen: byteweise identisch. oneLine gibt Werte ohne \r\n unveraendert zurueck. renderCard hat ausser renderCardBlock (internal/tui/view.go:545) keinen Aufrufer, es gibt also keinen zweiten Kartenpfad.

DoD einzeln und mit eigenen Pruefungen nachgestellt, alle sechs halten:
1+2 Slot 1 = 5;83 (erster Tag), Slot 2 = 5;45 (zweiter), verschieden; umgekehrte Tag-Reihenfolge kehrt die Farben um; Slot 3 traegt die Lane-Schattierung einer ungetaggten Karte, auch bei drei Tags.
3 Wegwerf-Board unter dem Scratchpad, frisch mit 'jaira init': 'jaira tag <id> ui backend docs ci' Exit 0, 'jaira show --json' listet alle vier.
4 Erster Tag gefaerbt, zweiter ohne Registry-Zeile: Slot 2 faellt auf die Lane-Schattierung, und die drei Textzeilen sind Zeichen fuer Zeichen die einer Karte ohne jeden Tag.
5 cardHeight()==3; Lane mit 30 Karten in einem 20 Zeilen hohen Terminal: es werden weniger als 30 gezeigt, und jede gezeigte traegt ihre Flag-Zeile '○ spec' — keine halb gezeichnete Karte.
6 core/release/NOTES.md: zwei Zeilen, beide einzeilig und mit '- ' beginnend, unter '## Unreleased' ueber '## 0.2.0'.

Zwei kosmetische Restpunkte, ausdruecklich kein Fehlschlag:
- Ein per CLI geschriebenes \r\n kommt aus dem YAML als ' \n' zurueck; oneLine macht daraus zwei Leerzeichen statt einem. Sichtbar nur als doppeltes Leerzeichen im Titel.
- Die neue NOTES-Zeile spricht nur vom Titel, der Fix deckt aber auch assignee, updated-by und executed-by ab.
- **2026-09-14 18:22 · Alexander Sacharov** — Entscheidung von Alex (2026-09-14, Antwort der human-Lane): Platz 3 ist NICHT mehr reserviert. Die Sprint-Markierung ist als eigenes Ticket 0YGWXQ herausgeloest und wandert an den RECHTEN Kartenrand. Damit gehoert der dritte linke Platz dem dritten Tag.

Endgueltige Regel: ein Tag faerbt nur die oberste Zelle, zwei Tags die oberen zwei, drei Tags alle drei; ab dem vierten bleiben die Tags auf dem Ticket und faerben nichts.

Konkret: cardColors (internal/tui/model.go:1369) laeuft heute 'i < cardSlots-1' und muss 'i < cardSlots' laufen. Die Gesamtfarbe der Karte (selectionFill, internal/tui/view.go:536) kommt weiterhin aus dem ERSTEN Tag - das ist richtig so und wird nicht angefasst.

Mitzuziehen: (1) Tests fuer den Drei-Tag-Fall im Stil der bestehenden zehn; (2) die Zeile unter '## Unreleased' in core/release/NOTES.md verspricht heute 'the third cell stays reserved and uncoloured' - diese Zeile wird umgeschrieben, keine zweite daneben gestellt.
- **2026-09-14 18:22 · Alexander Sacharov** — Die alte human-Frage ('faerbt eine Karte mit genau einem Tag nur noch Zeile 1?') ist beantwortet: ja, so soll es sein. Eine Karte mit einem Tag faerbt nur die oberste Zelle. Der question-Eintrag ist deshalb geleert.

Veraltet im context-Abschnitt oben: der Absatz 'Offen und absichtlich NICHT Teil dieses Tickets: was den dritten Platz fuellt ... bleibt Platz drei reserviert und leer'. Diese Frage ist entschieden - die Sprint-Markierung ist Ticket 0YGWXQ und liegt am rechten Kartenrand, der dritte linke Platz gehoert dem dritten Tag.
- **2026-09-14 18:28 · Alexander Sacharov** — in-progress (Runde nach der human-Antwort): Slot 3 gehoert jetzt dem dritten Tag. Geaendert wurde genau eine Schleifengrenze - cardColors (internal/tui/model.go:1367) laeuft 'i < cardSlots' statt 'i < cardSlots-1'. renderCardBlock brauchte keine Codeaenderung: es las schon immer slots[i] je Zeile, der reservierte Slot war allein ein Modell-Detail. Nur die Kommentare dort und an cardColors sagten 'reserviert' und wurden nachgezogen.

Tests umgeschrieben statt ergaenzt, weil drei von ihnen die alte Zusage FESTSCHRIEBEN und nach der Aenderung rot geworden waeren: TestThirdSlotStaysUncolouredHoweverManyTags (jetzt TestTagsPastTheThirdColourNoSlot - prueft, dass ein vierter Tag nichts faerbt und die ersten drei nicht verschiebt), TestThirdRowAlwaysCarriesTheLaneShade (jetzt TestThreeTaggedCardShowsAllThreeColoursInTicketOrder - Tabelle ueber 0/1/2/3/4 Tags plus umgekehrte Reihenfolge) und TestFourTaggedCardRendersWithTheExtraTagsUncoloured (prueft jetzt nur noch 5;111 als ungefaerbt; 5;200 ist der dritte Tag und MUSS erscheinen - genau das war die alte Zusage, die gekippt ist).

Neu dazu: TestSlotsBelowTheLastTagStayUncoloured. Der Fall 'weniger Tags als Slots' war vorher nur nebenbei gedeckt, weil Slot 3 ohnehin immer leer war; jetzt ist er die einzige Stelle, die verhindert, dass eine Karte mit einem Tag den Balken wieder durchgehend faerbt.

NOTES.md: die bestehende Zeile umgeschrieben, keine zweite daneben - eine Release-Zeile beschreibt, was der Leser sieht, und 'zwei Farben, dritte reserviert' hat nie ein Binary erreicht (Abschnitt ist ## Unreleased).

Nicht angefasst, bewusst: selectionFill/Glow speisen sich weiter aus Slot 1 (Alex' Vorgabe), und in core/tag bleibt es bei einer reinen Anzeigegrenze ohne Validierung.
- **2026-09-14 18:33 · Alexander Sacharov** — critique (Runde nach der human-Antwort): drei Befunde, alle mit klarem Fix, deshalb zurueck nach in-progress.

1. internal/tui/tagbox_test.go: TestCardColorsAreTheFirstThreeTagsInTicketOrder (:92) und TestTagsPastTheThirdColourNoSlot (:107) sind nach dem Umschreiben identisch - gleicher Aufruf, gleiche Schleife ueber 83/45/200, gleiche Schicht (cardColors). Der zweite haengt nur einen vierten Tag an und ist damit Obermenge. Vorher pruefte jeder etwas anderes (Slots 1+2 gefaerbt vs. Slot 3 ungefaerbt), die Verdopplung entsteht erst in diesem Diff. Fix: TestCardColorsAreTheFirstThreeTagsInTicketOrder loeschen; die Ticket-Reihenfolge ist ausserdem in TestThreeTaggedCardShowsAllThreeColoursInTicketOrder samt Umkehrprobe gedeckt. Dann auch die proof-Zeile von DoD-Punkt 2 nachziehen, die ihn nennt.

2. internal/tui/model.go:1352: der Kommentar sagt 'The sprint marker no longer waits for slot 3 - it goes to the card's right edge instead'. Das ist Praesens fuer Code, den es nicht gibt; wer am rechten Rand danach sucht, findet nichts, und das Ticket 0YGWXQ steht nicht dabei. Fix: Satz streichen oder auf 'the sprint marker moved off this bar (ticket 0YGWXQ)' kuerzen.

3. Ticketdatei Zeile 63: die proof-Zeile von DoD-Punkt 3 nennt TestFourTaggedCardRendersWithTheExtraTagsUncoloured, den dieser Diff in ...ExtraTagUncoloured umbenannt hat. Fix: jaira dod S1VM40 3 --done --proof mit dem neuen Namen.

Nicht beanstandet und bewusst stehen gelassen: die Aenderung selbst ist eine Schleifengrenze - kleiner geht die Form nicht; renderCardBlock unangetastet ist richtig, der reservierte Slot war nur ein Modell-Detail; selectionFill/Glow weiter aus Slot 1 ist Alex' Vorgabe; die eine umgeschriebene statt zweite NOTES-Zeile ist richtig, weil '## Unreleased' nie ein Binary erreicht hat; die Ueberschneidung zwischen Modell- und Rendering-Tests hat schon die vorige Runde stehen lassen und wird nicht neu aufgemacht.
- **2026-09-14 18:36 · Alexander Sacharov** — in-progress (Fixrunde nach der dritten critique): alle drei Befunde behoben, keine Codeaenderung an cardColors selbst.

1. TestCardColorsAreTheFirstThreeTagsInTicketOrder ersatzlos geloescht (internal/tui/tagbox_test.go). Nicht ersetzt oder umgeschrieben: TestTagsPastTheThirdColourNoSlot prueft denselben Aufruf mit einem zusaetzlichen vierten Tag und ist damit echte Obermenge; die Ticket-Reihenfolge deckt TestThreeTaggedCardShowsAllThreeColoursInTicketOrder samt Umkehrprobe.
2. Der Kommentar an cardColors (internal/tui/model.go:1351) sagt jetzt 'the sprint marker moved off this bar (ticket 0YGWXQ)' statt im Praesens einen rechten Kartenrand zu versprechen, den kein Code traegt. Gegengeprueft: sonst nennt keine Stelle in internal/tui einen Sprint-Marker oder einen reservierten Slot (grep auf 'right edge|0YGWXQ|reserved').
3. proof von DoD-Punkt 3 auf den umbenannten TestFourTaggedCardRendersWithTheExtraTagUncoloured (internal/tui/tagbox_test.go:593) gezogen. Punkt 2 nannte den geloeschten Test und ist mitgezogen - der Befund erwaehnte das; ohne wuerde Fix 1 einen zweiten toten Beleg hinterlassen.

go test ./... Exit 0, go vet ./... und gofmt -l ohne Ausgabe.
- **2026-09-14 18:38 · Alexander Sacharov** — Dispatcher-Hinweis: die Fixrunde nach der ersten critique hat ihre Aenderungen im Arbeitsbaum liegen lassen, ohne zu committen. Der Dispatcher hat sie unveraendert als 82c111e committet (fix(S1VM40): drop the duplicated slot test and correct the sprint-marker comment), zusammen mit der Ticketdatei. Es wurde nichts am Inhalt geaendert; go test ./internal/tui war vor dem Commit gruen und gofmt meldete nichts.
- **2026-09-14 18:42 · Alexander Sacharov** — critique (Runde 4): ein Befund, Fix klar, deshalb zurueck nach in-progress.

1. Ticketdatei Zeile 25, Frontmatter definition-of-done: Absatz 2 lautet weiter 'Platz 3 bleibt in diesem Ticket unbelegt und zeigt die Schattierung der Lane'. Der Diff hat goal (Zeile 7) und den Checklistenpunkt 2 (Zeile 63) auf die human-Entscheidung gezogen, dieses Feld aber nicht. Warum das zaehlt: 'jaira show --for-lane <lane> --json' gibt definition-of-done als input aus - diese critique-Runde hat den Widerspruch im eigenen Payload vorgelegt bekommen, goal und definition-of-done sagen einander entgegengesetzte Dinge. review und signoff bekommen ihn genauso. Fix: jaira set S1VM40 definition-of-done=<sechs Absaetze>, Absatz 2 im Wortlaut von DoD-Punkt 2, die uebrigen fuenf unveraendert.

Nachgeprueft, damit der Fix nicht halb wirkt: 'jaira set definition-of-done' schreibt nur das Frontmatter, nie den Body (auf einem Wegwerf-Board unter dem Scratchpad nachgestellt: Feld geaendert, die Prosa unter '## Definition of Done' blieb stehen).

Nicht als Befund erhoben:
- Die losen Prosa-Absaetze im Body unter der Checkliste (Zeilen 74-82) sind derselbe alte Satz. Sie stammen aber aus der Ticketerstellung - ein mehrabsaetziges --dod wird zu EINEM Kaestchen plus Prosa (auf dem Wegwerf-Board reproduziert) - sind aelter als dieser Diff und ueber die CLI gar nicht erreichbar. Das ist ein jaira-Fehler fuer ein eigenes Ticket, nicht der dieses Diffs.
- TestTwoTaggedCardShowsBothColoursInTicketOrder (internal/tui/tagbox_test.go:500) ueberschneidet sich mit der neuen Tabelle in TestThreeTaggedCardShowsAllThreeColoursInTicketOrder (:530), ist aber keine Obermenge: der Vertauschungsfall {backend,ui} steht nur dort. Kein zweiter Loeschkandidat wie in Runde 3.
- Die Aenderung selbst bleibt unbeanstandet: eine Schleifengrenze cardSlots-1 -> cardSlots (internal/tui/model.go:1364), renderCardBlock unangetastet, eine umgeschriebene statt einer zweiten NOTES-Zeile. Kleiner geht die Form nicht.
- Die drei Befunde aus Runde 3 sind erledigt und werden nicht neu aufgemacht: der doppelte Slot-Test ist weg, der Kommentar an cardColors (internal/tui/model.go:1350) nennt jetzt 'the sprint marker moved off this bar (ticket 0YGWXQ)', und die proof-Zeilen der Punkte 2 und 3 zeigen auf Tests, die es gibt.
- Die veralteten Testnamen im Feld review-gaps (TestThirdSlotStaysUncolouredHoweverManyTags, TestThirdRowAlwaysCarriesTheLaneShade) sind Ausgabe der optimize-Lane und werden von deren naechstem Durchgang ueberschrieben - anders als die DoD, die der Vertrag des Tickets ist.
- **2026-09-14 18:43 · Alexander Sacharov** — Befund der zweiten critique behoben, und zwar vom Dispatcher selbst, weil es sein eigener Fehler war: beim Wiedereroeffnen wurde DoD-Punkt 2 im Body umformuliert, das Frontmatter-Feld definition-of-done aber nicht. Absatz 2 dieses Feldes traegt jetzt denselben Wortlaut wie DoD-Punkt 2; die uebrigen fuenf Absaetze sind unveraendert. Keine Codeaenderung in dieser Runde - der Diff aus 50f15a2/82c111e steht unberuehrt.
