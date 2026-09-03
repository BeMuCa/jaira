---
id: 01M1EGJFM7MN1Q08YT1H79GEPW
title: "Tickets tragen Tags: Feld, Farb-Registry und Kommandos"
status: done
ready: true
creator: BeMuCa
goal: "Ein Ticket traegt ein tags-Listenfeld; jaira tags zeigt alle aktiven Tags mit Farbe und Anzahl; jaira tag <id> <name>... fuegt hinzu (neuer Tag bekommt --color oder eine zufaellige freie Palettenfarbe); der Board-Filter versteht tag:<name>; die Agenten-Doku sagt: erst jaira tags lesen, nie einen zweiten Tag fuers selbe Thema erfinden"
context: "Berk am 01.09.: Tickets sollen Tags bekommen (z.B. UI), damit man im Backlog alle UI-Tickets filtert und nach todo zieht. Wichtig gegen Wildwuchs: ein Kommando zeigt die existierenden Tags, und Agenten muessen es VOR dem Taggen lesen - das gehoert in den verwalteten Agent-Block (core/board/announce.go schreibt ihn) und in .claude/skills/jaira/SKILL.md und docs/COMMANDS.md. Design-Vorgaben: tags ist ein normales Listen-Frontmatter-Feld (wie blocked-by; jaira set tags=a,b funktioniert dann von selbst - pruefen ob listFields in internal/cli/tickets.go tags kennen muss); Tag-Namen kebab/lowercase; Farben leben NICHT je Ticket sondern in einer geteilten, handeditierbaren Datei .jaira/tags (eine Zeile 'name: <ansi256>' plus Kommentarkopf - diff-lesbar, merge-arm, Format ist die API); Palette = feste Liste gut unterscheidbarer ANSI-256-Farben, 'random' = zufaellige noch freie. jaira tags listet Name, Farbfeld, Anzahl offener Tickets (Logbook zaehlt nicht). create bekommt --tag (wiederholbar). Der Filter (view.go/model.go, versteht heute assignee:/lane:/ticket:) lernt tag:. tags erscheint als Basis-Zeile im Detail- und Signoff-View (Entscheidung 8-Erweiterung: immer sichtbar wenn gesetzt)."
definition-of-done: "tags-Feld wird gespeichert und in Detail+Signoff angezeigt; jaira tags listet aktive Tags mit Farbe+Anzahl; jaira tag <id> <name> [--color] legt neue Tags mit freier Zufallsfarbe an und meldet Wiederverwendung; Filter tag:ui trifft; Agent-Block, SKILL.md und COMMANDS.md dokumentieren die Erst-jaira-tags-lesen-Regel; go test ./... -race gruen"
blocked-by: []
commits: []
created-at: 2026-09-01T12:54:42Z
updated-at: 2026-09-03T10:29:03Z
claimed-by: EE-3NX6GL3-3597872
claimed-at: 2026-09-01T12:55:37Z
updated-by: BeMuCa
assignee: BeMuCa
outcome-what: "Runde 3: Zaehl-Dedup je Ticket, /tags-Anker, Newline-Idempotenz, Suggest zieht sich bei Buchstaben-Verlust zurueck, Lock-Kommentar praezisiert"
outcome-why: "Rest-Nits aus dem Re-Review; Commit vom Koordinator, da der Agent nach fertigen Edits am Rate-Limit starb"
outcome-resolves: "Voller -race-Lauf EXIT=0 auf dem Commit-Baum; gofmt nur bekanntes tickets.go"
review-summary: none
review-gaps: "Runde 2: nichts entfernt (hasTag geloescht zugunsten des einen tag.Matches - Duplikat weg). Alle 12 Funde geschlossen, Concurrency-Fix gegengeprobt (ohne Lock 7/8 Tags verloren, mit Lock gruen). Bewusst: CodeBadTag warnt NICHT bei reinen Gross/Klein-Abweichungen (UI zaehlt und filtert als ui - Warnung waere Rauschen ueber Funktionierendes); L14-Kosmetik unangetastet."
review-check: "Neu bauen, dann der 20-Sekunden-Concurrency-Check des Reviewers: Wegwerf-Board, ein Ticket, 'jaira tag $ID alpha & jaira tag $ID beta & wait', cat .jaira/tags -> BEIDE Zeilen da. Danach: jaira tag <id> ui --json -> tags_new/tags_reused-Arrays; jaira validate auf einem Ticket mit set-gesetztem 'front/end'-Tag -> Warnung mit fertigem Umbenennungs-Befehl; TUI-Filter tag:cur trifft security NICHT mehr."
review-verdict: "accept (Opus, 2 adversariale Runden mit eigenen Proben: Concurrency-Verlust, Gruppen-Fehleinsortierung, check-attr, Zaehlprobe - alle an der Wurzel gefixt und revert-gepinnt; die 5 finalen Nits vom Koordinator im Diff verifiziert und mit gruener Suite committet)"
---

# Tickets tragen Tags: Feld, Farb-Registry und Kommandos

## Definition of Done

- [ ] tags-Feld wird gespeichert und in Detail+Signoff angezeigt; jaira tags listet aktive Tags mit Farbe+Anzahl; jaira tag <id> <name> [--color] legt neue Tags mit freier Zufallsfarbe an und meldet Wiederverwendung; Filter tag:ui trifft; Agent-Block, SKILL.md und COMMANDS.md dokumentieren die Erst-jaira-tags-lesen-Regel; go test ./... -race gruen

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/tag: neues Paket mit Normalize (lowercase kebab), Palette (16 ANSI-256-Farben) und Registry (.jaira/tags lesen/schreiben, unbekannte Zeilen erhalten)
  proof: core/tag/tag.go: Normalize, Palette (16), Registry Load/Set/Assign/Save; core/tag/tag_test.go gruen
- [x] core/ticket/schema.go: FieldTags, canonicalOrder neben blocked-by, Ticket.Tags in Decode
  proof: core/ticket/schema.go: FieldTags, canonicalOrder vor blocked-by, Ticket.Tags; core/ticket/tags_test.go gruen
- [x] core/merge: tags als Listenfeld vereinigen statt eine Seite waehlen
  proof: core/merge/merge.go listFields kennt FieldTags
- [x] internal/cli: listFields lernt tags; ticketJSON traegt tags; create bekommt --tag; list bekommt --tag
  proof: internal/cli/tickets.go: listFields, ticketJSON tags, create --tag, list --tag; TestListFiltersByTag gruen
- [x] internal/cli/tags.go: jaira tags (Name, Swatch, offene Anzahl, --json) und jaira tag <id> <name>... [--color]
  proof: internal/cli/tags.go: newTagsCmd/newTagCmd, in root.go registriert; Smoke-Test gegen Temp-Board OK
- [x] internal/tui: tags-Basiszeile in detailBody und renderSignOff; Filter lernt tag:<name>
  proof: internal/tui view.go/signoff.go/model.go; internal/tui/tags_test.go gruen
- [x] Doku: agentNote in core/board/announce.go, SKILL.md, docs/COMMANDS.md — erst jaira tags lesen, nie ein Synonym erfinden
  proof: core/board/announce.go agentNote, .claude/skills/jaira/SKILL.md, docs/COMMANDS.md
- [x] Tests: Feld-Round-Trip, Registry-Lesen/Schreiben mit unbekannten Zeilen, Farbvergabe bis Palette erschoepft, tags-Listing, tag-Kommando, Filter tag:
  proof: core/tag/tag_test.go, core/ticket/tags_test.go, internal/cli/tags_test.go, internal/tui/tags_test.go
- [x] Verifikation: go test ./... -race -count=1 gruen, gofmt -l core internal nur internal/cli/tickets.go
  proof: go test ./... -race -count=1 -> EXIT=0 (race-final.log, alle 15 Pakete ok); gofmt -l core internal listet nur internal/cli/tickets.go (vorher schon unformatiert)

## Progress
- **2026-09-01 13:21 · BeMuCa** — Farb-Registry als eigenes Paket core/tag, nicht in core/ticket: die Datei ist eine Board-Tatsache wie core/lane/order (Plaintext, eine Zeile pro Eintrag, absente Datei ist kein Fehler). core/tag importiert core/ticket nur fuer DirName, core/ticket kennt tags weiter nur als []string - kein Zyklus.
Farbe wird nur beim Schreiben vergeben. 'jaira tags' schreibt nichts: ein Lesekommando, das das Board aendert, streitet sich mit jeder parallelen Session um dieselbe Datei. Handgesetzte Tags (z.B. via 'jaira set tags=a,b') stehen deshalb farblos in der Liste ('-' plus 'no colour yet'); Farbe kommt beim naechsten 'jaira tag' oder per Hand. Das war die Gabel assign-on-list vs. show-as-colorless - gewaehlt: colorless, weil unmagischer.
Neue Zeilen werden alphabetisch eingefuegt, nicht angehaengt. Anhaengen laesst zwei parallele Sessions dieselbe letzte Zeile schreiben - genau die Stelle, an der git garantiert konfligiert. Eine unsortierte, handgruppierte Datei wird NICHT umsortiert; neue Zeilen landen dann nach dem letzten Eintrag. Kommentare, Leerzeilen und unparsebare Zeilen bleiben wortwoertlich stehen, ein Recolour ersetzt genau eine Zeile inkl. ihres Trailing-Kommentars und dessen Einrueckung.
Palette hat 16 Farben und laesst 39/214/203/78/141 bewusst aus, weil das Board die schon fuer accent/warn/err/ok/agentic ausgibt - ein Test pinnt das. Nach 16 Tags wiederholen sich Farben deterministisch ueber fnv32a(name)%16 statt weiter zufaellig: sonst schreibt jeder Lauf eine neue Farbe und die Datei hoert nie auf zu wackeln.
tags wird in Detail- und Signoff-View absichtlich NICHT eingefaerbt. styleLines-Regel: ein umgebrochener Block, der nach einer Label-Spalte gestylt wird, wird auf die breiteste Zeile gepolstert und laeuft aus dem Pane. Die Farben zahlen sich in 'jaira tags' aus.
core/merge listFields kennt FieldTags jetzt - zwei Sessions, die aus zwei Themen taggen, haben beide recht; eine Seite waehlen hiesse eines still verwerfen.
Ticket.Tags liegt im kommentierten Feldblock nach UpdatedBy, nicht neben Commits: ein Doc-Kommentar dort haette die gofmt-Alignment-Gruppe geteilt und sechs unbeteiligte Zeilen umformatiert.
Der Agent-Block ist nicht versioniert - NoteIsCurrent vergleicht den generierten Text mit dem im File. Nach dem naechsten Binary-Update meldet 'jaira validate' auf JEDEM Board AGENTS.md/CLAUDE.md als veraltet, bis 'jaira update' oder irgendein 'jaira lanes'-Kommando den Block neu schreibt. Dieses Repos eigene Kopien wurden absichtlich nicht mitregeneriert: der Dev-Build hat release.Current='dev', und 'jaira update' wuerde damit den Board-Stempel auf 'dev' setzen.
colourable() nutzt os.File + ModeCharDevice statt eines neuen Terminal-Pakets (charmbracelet/x/term waere von indirect auf direct gewandert). Effekt: Tests (bytes.Buffer) und Pipes bekommen plain Text, nur ein echtes TTY sieht ANSI.
- **2026-09-01 13:32 · BeMuCa** — review accept mit 12 Funden, Fix-Runde fuer: (L4) Save nicht-atomar+ohne Lock - Probe verlor einen Tag bei zwei parallelen Sessions; writeAtomic+Store-Lock existieren schon; (L5) Kommentar verspricht 'hand-gruppiert bleibt gruppiert', Code sortiert ein - Probe filed einen backend-Tag unter den frontend-Header; (L6) kein merge=union fuer .jaira/tags; (L7) set-geschriebene Nicht-Kebab-Tags sind unauffindbar - validate-Warnung; (L8) Set schreibt erste, Load liest letzte Duplikatzeile -> --color wirkungslos; (L9) gate fieldFilled kennt FieldTags nicht - eine tags-produzierende Lane kann nie erfuellen; (L10) tag --json verliert das new/reused-Signal - fuer Agenten das Wichtigste; (L11) JSON colour vs Flag --color; (L12) tag:-Filter TUI substring vs CLI exakt - beide exakt + Test; (L13) eine vakuum-Assertion + tags --json ungetestet; (L15) docs/AGENTS.md lernt tags nie + NOTES-0.1.1-Zeile fuer tags und Agent-Block-Drift. Belassen: (L14) kosmetische Feld-Reihenfolge beim Retrofit.
- **2026-09-01 13:44 · BeMuCa** — Review-Nachlese, die zehn Befunde und was daraus wirklich folgte.
Der bewiesene Datenverlust war nicht die fehlende Atomizitaet, sondern die fehlende Serialisierung. Ein atomarer Write macht jeden einzelnen Write ganz, nicht das Paar geordnet: zwei Sessions lesen beide die Datei, haengen beide ihre Zeile an das an, was sie gelesen haben, und der zweite Save schreibt eine Datei, die vor der Zeile des ersten zusammengebaut wurde. Beides ist jetzt drin - ticket.WriteAtomic (exportiert, statt einer zweiten Kopie des tmp+rename-Tanzes) plus s.Lock('tags') um das ganze Load-Set-Save in registerTags. Gegenprobe gemacht: mit entferntem Lock verliert TestRegisterTagsLosesNoTagUnderConcurrentWriters 7 von 8 Tags, mit Lock gruen.
Store.Lock ist ein exportierter Wrapper um s.lock. Kollisionsfrei per Konstruktion: Lock-Dateien heissen <name>.lock und Ticket-Ids sind 26-stellige ULIDs, 'tags' kann also nie ein Ticket-Lock sein.
insertEntry hat vorher gelogen: der Kommentar versprach 'handgruppierte Datei bleibt gruppiert', der Code sortierte trotzdem ein. Jetzt wird die ganze Datei erst gescannt und nur einsortiert, wenn die vorhandenen Eintraege schon sortiert sind - sonst haengt die neue Zeile nach dem letzten Eintrag an. Grund: in eine per Hand gruppierte Datei einsortieren stellt ein backend-Tag unter eine '# frontend'-Ueberschrift, und dann beschreibt der Kommentar ueber einer Zeile sie nicht mehr. Das ist schlimmer als unsortiert.
Set schreibt jetzt die LETZTE passende Zeile um, nicht die erste - Load ist last-wins, die erste umzuschreiben haette eine Zeile geaendert, die niemand liest, und die Farbe unveraendert gelassen. Duplikate sind kein Sonderfall mehr, seit .jaira/.gitattributes 'tags merge=union' bekommt: union laesst genau dann zwei Zeilen fuer einen Tag stehen, wenn beide Seiten ihn umfaerben. Harmlos, weil last-wins liest und Set dieselbe letzte Zeile ueberschreibt - die beiden Fixes haengen also zusammen.
writeGitAttributes prueft die zwei Zeilen jetzt getrennt. Ein Board, das vor der Registry eingerichtet wurde, nennt den Ticket-Driver schon, und ein einziges 'ist irgendwas konfiguriert'-Test haette es fuer immer ohne die union-Zeile gelassen.
Der TUI-Filter tag: ist jetzt exakt statt Substring, ueber tag.Matches - dieselbe Funktion, die 'jaira list --tag' nutzt. 'tag:cur' hat vorher alles mit 'security' geliefert: eine falsche Antwort, keine unscharfe, und der Board-Filter widersprach dem CLI-Flag. Volltextsuche enthaelt Tags weiter unveraendert.
validate warnt nur bei Namen, die Normalize ABLEHNT, nicht bei Gross-/Kleinschreibung: 'UI' normalisiert zu 'ui', wird also gezaehlt und von beiden Filtern getroffen - dafuer zu warnen waere Rauschen ueber etwas, das funktioniert. Der Vorschlag nennt die ganze reparierte Liste, weil 'jaira set' ein Listenfeld ersetzt und nicht anhaengt; dafuer gibt es tag.Suggest (best-effort-Reparatur), das bewusst NUR beraet - dasselbe stille Kuerzen automatisch anzuwenden waere genau die Synonym-Erfindung, die das Vokabular verhindern soll.
JSON-Key heisst 'color', nicht 'colour', passend zum --color-Flag; eine Schreibweise auf der Maschinenoberflaeche. 'jaira tag --json' traegt jetzt tags_new/tags_reused, immer Arrays und nie null, weil SKILL.md Agenten genau auf --json schickt und das neu-vs-wiederverwendet-Signal der ganze Punkt des Kommandos ist.
gate.fieldFilled kannte tags nicht: als Liste haette es eine Lane, die tags produziert, ewig schuldig gelassen. Jetzt neben FieldCommits behandelt (Liste = gefuellt wenn nicht leer).
- **2026-09-01 13:54 · BeMuCa** — review Runde 2 accept mit 4 Rest-Nits + 1 Kommentar: (L14) countTags dedupliziert nicht je Ticket - ui+UI auf einem Ticket zaehlt '2 open', und merge-union produziert genau solche Paare; counted-Map je Ticket. (L15) 'tags merge=union' ist unverankert - matcht auch .jaira/tickets/tags und logbook/*/tags (git check-attr belegt); '/tags'. (L16) Idempotenz-Check scheitert an fehlendem End-Newline und haengt die Zeile doppelt an; gegen newline-normalisierte Kopie pruefen. (L17) Suggest('ueber')='ber' wird als ausfuehrbarer set-Befehl gedruckt - Vorschlag unterdruecken, wenn die Reparatur einen Buchstaben verwirft. (L18) Lock-Kommentar verspricht 'arbitrary name' - auf 'kein 26-Zeichen-ULID' praezisieren. (L19) vorbestehende WriteAtomic-Kandidaten nur benannt, v.a. core/lane/order.go als genanntes Vorbild - eigenes Thema, nicht dieses Ticket.
- **2026-09-01 18:57 · BeMuCa** — Finale: der Implementierer-Agent beendete alle fuenf Nit-Fixes und starb am Opus-Session-Limit direkt vor dem Race-Lauf. Koordinator hat den Diff gelesen (Suggest-Rueckzug bei Buchstaben-Verlust, /tags-Anker mit check-attr-Begruendung, Newline-Normalisierung VOR dem Idempotenz-Test, counted-Map je Ticket, ehrlicher Lock-Kommentar), go test ./... -race -count=1 auf exakt diesem Baum mit EXIT=0 (0 FAIL/RACE) gefahren und committet.
