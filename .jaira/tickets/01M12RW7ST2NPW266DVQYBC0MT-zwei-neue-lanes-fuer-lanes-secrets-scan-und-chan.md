---
id: 01M12RW7ST2NPW266DVQYBC0MT
title: "Zwei neue Lanes fuer lanes/: secrets-scan und changelog-writer"
status: done
ready: true
creator: BeMuCa
goal: "Beide Lanes liegen als Dateien in lanes/, parsen in CI, und ihr Prompt sagt, wann die Lane fertig ist"
context: "Berk am 27.08. nach der Lane-Recherche (.planning/research/lane-ideas.md, Ideen 5 und 12): diese zwei gefallen ihm. secrets-scan: billigste Pruefung, faengt den schlimmsten Einzelfehler (ein committetes Credential), laeuft auf jedem Ticket, Tier cheap, Eingabe diff, Ausgabe secrets-status. changelog-writer: klein, billig, dient direkt dem Kernwert 'man verliert nie, wofuer eine Aufgabe war' - fuer Endnutzer, nicht nur Agenten; Eingabe outcome-what/-why, Ausgabe changelog-entry; Tier cheap. Beide Prompt-Entwuerfe stehen im Recherchebericht. Offen: wo sie im Ablauf sitzen (secrets-scan nach in-progress vor critique; changelog-writer nach review vor signoff?) und ob changelog-entry in core/release/NOTES.md geschrieben werden soll oder nur auf dem Ticket steht."
definition-of-done: "lanes/secrets-scan.md und lanes/changelog-writer.md existieren, go test ./core/lane/ parst sie, jaira lanes market listet sie, jede hat after:, model-tier, input-requires, output-produces und einen Prompt mit Endbedingung, lanes/README.md hat je eine Tabellenzeile"
blocked-by: []
commits: []
created-at: 2026-08-27T23:28:57Z
updated-at: 2026-09-03T10:27:44Z
claimed-by: EE-3NX6GL3-2809077
claimed-at: 2026-08-31T17:58:21Z
updated-by: BeMuCa
assignee: BeMuCa
outcome-what: "Drei Textfunde behoben, kein Go-Code angefasst. lanes/README.md: die Sits-Zelle von secrets-scan sagt statt 'ahead of every other gate' jetzt 'after implementing, once you move the column there', und ein neuer Absatz 'Belongs, not lands' erklaert einmal fuer die ganze Tabelle, dass jaira lanes add die Lane als letzte Zeile an .jaira/lanes/order haengt und after: nur ohne order-File gelesen wird. lanes/secrets-scan.md: .pem ist nicht mehr bedingungslos ein Fund, sondern nur mit echtem PRIVATE-KEY-Block; .p12, .pfx und id_rsa bleiben bedingungslos. lanes/secrets-scan.md: 'one per line' ist 'one per finding, separated by \";\"' geworden, mit dem Hinweis, dass das Feld eine Zeile ist."
outcome-why: "Alle drei waren Stellen, an denen der Lane-Text mehr behauptete als das Werkzeug tut, und genau das kostet den, der die Lane adoptiert, Vertrauen. Fund 1 nachgeprueft: jaira init schreibt ein order-File, jaira lanes add haengt hinten an, der Anker wird nicht gelesen -- die Spalte landet also rechts aussen. Der Absatz deckt bewusst die ganze Tabelle ab, weil critique, optimize und changelog-writer dasselbe Problem gehabt haetten; so bleiben die vorhandenen Zeilen unangetastet. order.go bleibt unveraendert, weil die Platzierungsfrage ein eigenes Ticket ist. Fund 2 war ein echter Selbstwiderspruch: die Nicht-Fund-Liste zehn Zeilen weiter unten nennt oeffentliche Zertifikate ausdruecklich, und CA-Bundles sind .pem -- die Lane haette genau das Rauschen erzeugt, das ihr eigener Prompt verbietet. Fund 3 war eine falsche Aussage ueber das Datenmodell: secrets-findings ist ein Skalar (jaira set schreibt per SetScalar), es gibt keine zweite Zeile, und das eigene Beispiel darunter trennte schon mit Semikolon."
outcome-resolves: "Verifiziert nach der Aenderung: go test ./core/lane -count=1 gruen (shipped_test parst beide Lane-Dateien erneut ohne Warnung), go test ./... -race -count=1 exit 0 und kein FAIL, gofmt -l core internal nennt weiterhin nur die schon vorher vorhandene Abweichung internal/cli/tickets.go. Commit 158d4732e3d36d64af147b163ce8fb035af2ad4f, nicht gepusht, kein Binary neu gebaut, nichts unter .jaira von Hand geaendert, keine Lane auf dieses Board adoptiert. Die drei Review-Befunde sind damit alle im Text behoben; order.go wurde absichtlich nicht angefasst."
review-summary: none
review-gaps: "Nichts entfernt. Gelassen: model-tier cheap fuer beide (Mustersuche bzw. Ein-Zeilen-Text; ein Board kann es per Lane-Datei aendern); die Sortier-Zufaelligkeit der after-Anker-Reihenfolge (secrets-scan landet alphabetisch vor critique - ehrlich in der Notiz dokumentiert, gilt nur ohne order-Datei); changelog-entry bleibt auf dem Ticket statt in einer Datei (Katalog-Lane darf keinen Repo-Pfad kennen, geteilte Datei = Merge-Konflikt)."
review-check: "1. go test ./core/lane -count=1: shipped_test parst beide neuen Dateien. 2. lanes/secrets-scan.md und lanes/changelog-writer.md lesen: Frontmatter vollstaendig (id/name/after/precedence/tier/requires/produces), Prompt sagt wann fertig. 3. In einem Wegwerf-Board: jaira lanes adopt lanes/secrets-scan.md && jaira lanes add secrets-scan - Lane erscheint, jaira lanes warnt nicht. 4. Praezedenzen gegen die bestehenden pruefen: 35 zwischen in-progress(30) und human(40), 52 zwischen review(50) und signoff(55)."
review-verdict: "accept (Opus-End-Review auf 5f61f35: Frontmatter/Praezedenzen/Prompts/Feld-Fallback alle verifiziert, reject galt nur den drei Textstellen; Delta 158d473 vom Koordinator gelesen: README verspricht keine Platzierung mehr - 'Belongs, not lands' erklaert die order-Datei ehrlich fuer die ganze Tabelle -, .pem an den PRIVATE-KEY-Block gebunden, Scalar-Wortlaut konsistent)"
---

# Zwei neue Lanes fuer lanes/: secrets-scan und changelog-writer

## Definition of Done

- [x] lanes/secrets-scan.md und lanes/changelog-writer.md existieren, go test ./core/lane/ parst sie, jaira lanes market listet sie, jede hat after:, model-tier, input-requires, output-produces und einen Prompt mit Endbedingung, lanes/README.md hat je eine Tabellenzeile

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-08-31 18:05 · BeMuCa** — Precedence: secrets-scan=35, changelog-writer=52. Beide sind Merge-Rang, nicht Spaltenposition (Invariante 1). Vorhandene Raenge: backlog 0, brainstorm 5, blocked 10, todo 20, pre-process 25, in-progress 30, human 40, critique 45, optimize 48, review 50, signoff 55, done 60. 35 ist die Mitte der einzigen freien Luecke hinter in-progress: strikt ueber 30, damit ein Merge das Scan-Ergebnis nicht auf 'wird noch gebaut' zurueckdreht, und strikt unter human/critique, damit ein Merge nie ein spaeteres Gate behauptet. 52 liegt zwischen review 50 und signoff 55 -- der Changelog-Satz entsteht nach dem Review, also darf ein Merge ihn nicht gegen 'nur reviewed' verlieren. Beide Werte kollidieren mit nichts und lassen nach oben und unten Platz fuer weitere Katalog-Lanes.

after: ist NICHT die Spaltenposition auf einem echten Board. Verifiziert: 'jaira lanes add' haengt die Lane nur an .jaira/lanes/order an (probe-Board im Scratchpad), der Anker wird dabei ignoriert. after: greift nur, wenn es kein order-File gibt -- genau der Fall, den core/lane/shipped_test.go faehrt. Wer adoptiert, bestimmt also die Spalte selbst; die Lane-Datei erklaert nur die Absicht.

Reihenfolge ohne order-File haengt an der Dateinamen-Sortierung: order() in core/lane/lane.go:771 schiebt jede Custom-Lane direkt hinter ihren Anker, in alphabetischer Datei-Reihenfolge. secrets-scan.md kommt alphabetisch nach critique.md, deshalb landet secrets-scan trotz gleichem Anker (in-progress) VOR critique. Verifiziert, ergibt in-progress -> secrets-scan -> critique -> optimize -> human -> review -> changelog-writer -> signoff. Das ist Zufall der Sortierung, nicht garantiert -- wenn jemand critique umbenennt, kippt es.

secrets-status, secrets-findings und changelog-entry sind keine modellierten Ticket-Felder, und das ist in Ordnung: gate.fieldFilled (core/gate/gate.go:530) fuellt unbekannte Felder aus Doc().Scalar, und 'jaira set' schreibt beliebige Keys per SetScalar (internal/cli/tickets.go:719). Beide Lanes sind also ohne eine Zeile Go-Code erfuellbar und ihr output-produces wird am Gate wirklich geprueft.

changelog-entry bleibt auf dem Ticket -- kein core/release/NOTES.md. Zwei Gruende: eine Katalog-Lane darf keinen Pfad aus jairas eigenem Repo kennen, sonst ist sie nicht adoptierbar; und eine gemeinsame Datei waere bei parallelen Tickets genau der Merge-Konflikt, den das Format vermeiden soll. Die offene Frage aus dem Ticket-Kontext ist damit im Lane-Design entschieden, nicht im Release-Prozess: Release Notes entstehen spaeter durch Einsammeln des Feldes.

secrets-scan liest nur hinzugefuegte Zeilen. Eine geloeschte Zeile steht schon in der History -- das ist ein Rotationsproblem, kein Review-Problem, und der Prompt sagt das explizit, damit die Lane nicht auf Diff-Rauschen anspringt.

Arbeitsbaum: core/gate/gate.go, core/validate/*, internal/tui/{view,signoff,home_test}.go waren beim Start schon geaendert (fremdes Ticket, OwedBy/leere Felder). Nicht von mir, nicht mitcommittet.
- **2026-08-31 18:44 · BeMuCa** — End-Review reject, drei Text-Funde (kein Go-Code): README:41 'ahead of every other gate' ist falsch (lanes add haengt ans order-Ende - eigenes Ticket dafuer angelegt); secrets-scan.md:32 .pem-Regel widerspricht der not-a-finding-Liste (CA-Bundles sind .pem) - .pem nur bei PRIVATE-KEY-Block, .p12/.pfx/id_rsa bleiben unbedingt; secrets-scan.md:58 'one per line' vs. Scalar-Feld - 'one per finding, separated by ; '.
- **2026-08-31 18:48 · BeMuCa** — Review-Runde 2, drei Textfunde, kein Go-Code angefasst.

1. lanes/README.md: 'ahead of every other gate' war falsch. Nachgeprueft im Wegwerf-Board: 'jaira init' schreibt .jaira/lanes/order, und 'jaira lanes add' haengt die Lane als LETZTE Zeile an -- der Anker after: wird dabei nicht gelesen. after: greift nur, wenn ueberhaupt kein order-File existiert, und das ist auf einem initialisierten Board nie der Fall. Die Tabellenzelle sagt jetzt 'after implementing, once you move the column there', und ein neuer Absatz 'Belongs, not lands' erklaert einmal fuer die ganze Tabelle, dass Sits und after: nur die Absicht sind und das Umsortieren dem Adoptierenden gehoert. Damit sind auch die schon vorhandenen Zeilen (critique, optimize, changelog-writer) mit abgedeckt, ohne dass ich sie umschreibe -- sie haetten dasselbe Problem gehabt. order.go bleibt unangetastet, die Platzierungsfrage ist ein eigenes Ticket.

2. lanes/secrets-scan.md: '.pem ... regardless of its contents' widersprach der Nicht-Fund-Liste zehn Zeilen weiter unten, weil CA-Bundles und oeffentliche Zertifikate ganz normal .pem sind. .p12, .pfx und id_rsa bleiben bedingungslos -- die tragen praktisch nie etwas Oeffentliches. .pem ist jetzt an einen echten '-----BEGIN ... PRIVATE KEY-----'-Block gebunden. Der Punkt war nicht nur formal: eine Lane, die jedes .pem meldet, produziert genau das Rauschen, das der Prompt zwei Absaetze weiter verbietet.

3. lanes/secrets-scan.md: 'one per line' war doppelt falsch. secrets-findings ist ein Skalar -- 'jaira set' schreibt es per Doc().SetScalar, es gibt also gar keine zweite Zeile -- und das eigene Beispiel direkt darunter trennte schon mit '; '. Steht jetzt als 'one per finding, separated by "; "' plus dem Hinweis, dass das Feld eine Zeile ist, damit niemand Absaetze hineinschreibt.

Alle drei waren Behauptungen der Textdateien ueber Tool-Verhalten, das ich vorher nur teilweise geprueft hatte -- Fund 1 hatte ich in der Notiz der ersten Runde selbst korrekt beschrieben und dann trotzdem falsch ins README geschrieben. Lehre fuer den naechsten Lane-Text: was die Notiz als verifiziert festhaelt, muss woertlich in die Datei, nicht in einer optimistischeren Fassung.
