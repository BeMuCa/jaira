---
id: 01M2NCH8ZK9524JZ00J4GTQHNH
title: "Der Dispatcher bekommt einen Gespraechsmodus, statt dass eine zweite Rolle daneben entsteht"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Ticket mit noch offenen Gestaltungsentscheidungen wird im Gespraech gefuehrt statt geraten: derselbe Dispatcher-Prompt zaehlt vor der Plan-Lane die offenen Entscheidungen, haelt bei mindestens einer an und schreibt die Antwort des Menschen mit 'jaira note' aufs Ticket, traegt den Modus auf dem Ticket selbst (damit er einen Sitzungsabbruch ueberlebt), legt den Code nach jedem DoD-Punkt vor, und committet nicht selbst, sondern gibt die fertige Commit-Zeile mit Handle zurueck"
context: |-
  Alex am 2026-09-16, aus zwei Laeufen dieser Woche.

  Was schiefgeht: laesst man ein Ticket autonom durch die Lanes laufen, dessen Form noch nicht feststeht - UI-Aenderungen und Datenbank-/Dateiformat-Aenderungen sind seine Beispiele -, dann raet der Worker. Am Ende steht viel Arbeit, die neu gemacht werden muss.

  Der Beleg liegt auf diesem Board. 9ZZSFT lief autonom in einem Durchgang durch. 0YGWXQ brauchte sieben critique-Runden und drei Entscheidungen von Alex. 0YGWXQ ist nicht groesser - es war weniger entschieden.

  Damit ist Alex' erster Vorschlag ('wenn ein Ticket gross aussieht') das falsche Mass. Das Mass ist, WIE VIELE ENTSCHEIDUNGEN IM TICKET NOCH OFFEN SIND: ein DoD-Punkt, den man auf zwei verschiedene Arten erfuellen kann und beide kommen durch den Gate, ist eine offene Entscheidung. Null davon heisst, der Lauf braucht niemanden.

  Warum es KEINE zweite Kommandozeile wird, obwohl Alex zuerst danach gefragt hat: was er beschreibt, existiert schon - es ist der Teamlead, auf ein Ticket verengt. Genau so lief diese Woche: mit Alex ueber die Form von 0YGWXQ reden, seine Antwort mit 'jaira note' auf das Ticket schreiben, dann Worker starten. Zwei Prompts, die fast dasselbe tun, driften auseinander. Dieses Repository leidet schon daran: 3YRPXJ und 7MG5GB beschreiben beide, dass spawn.sh nicht erweiterbar ist und deshalb geforkt wird. Alex hat dem Modus-statt-Rolle am 16.09. zugestimmt.

  Zwei Eigenschaften hat Alex am 16.09. nachgereicht, und die zweite hat eine Falle:

  Erstens, der Mensch sieht den Code, bevor darauf aufgebaut wird. Nicht am Ende, wenn das Zurueckdrehen am teuersten ist.

  Zweitens, in diesem Modus wird nicht automatisch committet. Die Falle: 'kein Commit' ist die Abwesenheit eines Mechanismus, nicht einer. Committet der Mensch von Hand, faellt mit hoher Wahrscheinlichkeit der Handle aus der Betreffzeile - und CLAUDE.md sagt seit 9ZZSFT ausdruecklich, dass genau dieser Handle die Quelle ist, aus der jaira die Commit-Liste ableitet, weil die Historie der Ticket-Datei bewusst duenn geworden ist. Ohne Handle bleibt die Liste leer und der Zug in die Endlane wird verweigert. Der Agent muss die fertige Commit-Zeile samt Handle zurueckgeben, so wie jaira-role-pr es mit der PR-Beschreibung schon macht - das Muster gibt es also.

  Noch nicht entschieden, Sache der Brainstorm-Lane: in welchen Haeppchen der Mensch den Code sieht. Der DoD-Punkt ist die naheliegende Einheit, weil er ohnehin das Inkrement ist. Zu fein und der Modus ist langsamer, als es selbst zu tippen.

  Nicht vor dem naechsten Release. 0YGWXQ und 7KX89C sind ein erklaerbares Paar; diese Rolle ist der Release danach. Und sie entsteht im Repository, nicht in ~/.claude - 7KX89C raeumt gerade genau das auf.
definition-of-done: "Der Eintritt in den Modus haengt an einer nachpruefbaren Bedingung, nicht am Bauchgefuehl: vor der Plan-Lane zaehlt der Dispatcher die noch offenen Entscheidungen des Tickets auf. Keine offene - er laeuft weiter wie heute. Mindestens eine - er haelt an und fragt den Menschen."
tags:
  - docs
blocked-by: []
related:
  - 01M2KBPPVH5PAKZ98B0X7KX89C
  - 01M2HRF34AS7QDTACC323YRPXJ
  - 01M2E5R7NKRK3ETAEKG14XHZ6N
commits:
  - 9cb1df92380b3e96ca46030822a91b946e288938
created-at: 2026-09-16T15:14:31Z
updated-at: 2026-09-16T15:48:16Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-53684
claimed-at: 2026-09-16T15:25:11Z
outcome-what: "Das Frontmatter-Feld 'mode' mit dem einen Wert 'conversational' quer durch Schema, CLI, TUI und die beiden Rollen-Prompts. Der Dispatcher haelt jetzt VOR der Plan-Lane an, wenn das Ticket noch offene Entscheidungen hat, schreibt die Antworten mit 'jaira note' aufs Ticket und setzt den Modus; der Worker legt in diesem Modus nach jedem DoD-Punkt 'git diff' vor und gibt statt eines Commits die fertige Commit-Zeile mit Handle zurueck. Dazu drei Tests, zwei NOTES.md-Zeilen und der Feld-Eintrag in README.md und docs/AGENTS.md."
outcome-why: "Ein Ticket, dessen Form noch nicht feststeht, wird autonom geraten - 0YGWXQ hat das mit sieben critique-Runden bezahlt. Der Modus musste auf der Platte landen und nicht in der getippten Zeile, weil nur so ein abgebrochener Lauf nicht stumm wieder autonom weiterfaehrt."
outcome-resolves: "Alle sechs DoD-Punkte sind getickt und belegt: der Eintritt haengt an der gezaehlten Zahl offener Entscheidungen (dispatcher SKILL.md:31), es bleibt bei einer Dispatcher-Rolle (core/role/builtin/ unveraendert sieben), die Entscheidung steht vor der Arbeit auf dem Ticket und ueberlebt den Abbruch (TestModeSurvivesRoundTrip), der Agent committet nicht mehr selbst sondern gibt die Zeile mit Handle zurueck (role-lane SKILL.md:67), der Code wird nach jedem Inkrement vorgelegt (role-lane SKILL.md:50) und beide NOTES.md-Zeilen stehen unter ## Unreleased."
review-summary: |-
  internal/cli/tickets.go:929 + core/ticket/schema.go:150: ValidMode trimmt, der Schreibpfad nicht — 'jaira set <id> "mode= conversational "' wird akzeptiert und als mode: " conversational " gespeichert, 'show --json' gibt es mit den Leerzeichen zurueck, und der Worker vergleicht laut Prompt auf genau ein Wort. Das ist derselbe Fehlerfall, den TestSetRefusesUnknownMode ausschliessen soll, nur durch die Vordertuer. Entweder das strings.TrimSpace in ValidMode streichen, dann faellt der gepolsterte Wert durch dieselbe Pruefung wie 'Conversational', oder in beiden Schreibpfaden den getrimmten Wert speichern.
  internal/tui/view.go:1199 und internal/cli/tickets.go:709: der Modus steht in keiner Ausgabe fuer Menschen. Beide Detail-Panes drucken row("tier", t.ModelTier), aber kein row("mode", t.Mode); view.go:1099 traegt FieldMode in fieldsWithTheirOwnRow ein, dessen Kommentar 'the fields this pane already has a place for' behauptet — die Zeile gibt es nicht, und laneFields erreicht das Feld ohnehin nie, weil keine Lane es produziert. internal/tui/edit.go:28 laesst es dagegen bearbeiten: schreibbar ueberall, lesbar nur in --json. Je eine Zeile row("mode", t.Mode) neben row("tier", ...); row() ueberspringt Leeres, ein Ticket ohne Modus kostet es also nichts.
  core/role/builtin/jaira-dispatcher/SKILL.md:79 gegen :155-163: zwei Regeln fuer --no-worktree an zwei Stellen. Der neue Abschnitt sagt 'im Gespraechsmodus immer --no-worktree', der bestehende Absatz sagt 'nimm es nur, wenn der Mensch danach fragt oder die Arbeit eine Lane lang ist', und Schritt 2 der Schleife sagt 'in its own worktree'. Den Modus in den bestehenden Absatz bei :160 aufnehmen und den Bullet bei :79 auf einen Verweis darauf kuerzen, statt dieselbe Regel zweimal zu fuehren.
---

# Der Dispatcher bekommt einen Gespraechsmodus, statt dass eine zweite Rolle daneben entsteht

## Definition of Done

- [x] Der Eintritt in den Modus haengt an einer nachpruefbaren Bedingung, nicht am Bauchgefuehl: vor der Plan-Lane zaehlt der Dispatcher die noch offenen Entscheidungen des Tickets auf. Keine offene - er laeuft weiter wie heute. Mindestens eine - er haelt an und fragt den Menschen.
  proof: core/role/builtin/jaira-dispatcher/SKILL.md:31 — 'Before the plan lane: count what is still open': offene Entscheidungen auflisten, keine = weiterlaufen, >=1 = anhalten und fragen
- [x] Es gibt weiterhin genau eine Dispatcher-Rolle. Kein zweiter Skill und keine zweite Kommandozeile daneben; der Modus steht im selben Prompt.
  proof: core/role/builtin/ enthaelt unveraendert sieben Rollen; der Modus steht in jaira-dispatcher/SKILL.md:31 und jaira-role-lane/SKILL.md:39, kein neuer Skill und keine neue Kommandozeile
- [x] Was im Gespraech entschieden wird, steht mit 'jaira note' auf dem Ticket, BEVOR die Arbeit daran beginnt - nicht hinterher. Nachgestellt an einem Ticket, dessen Sitzung mittendrin abgebrochen wird: die Entscheidung ist danach noch da.
  proof: core/role/builtin/jaira-dispatcher/SKILL.md:50 (jaira note vor Arbeitsbeginn) plus TestModeSurvivesRoundTrip in internal/cli/mode_test.go — der Modus kommt nach dem Schreiben von der Platte zurueck
- [x] In diesem Modus committet der Agent nicht selbst. Er legt die Aenderungen bereit und gibt eine fertige Commit-Zeile zurueck, die den Ticket-Handle im Betreff traegt und die Ticket-Datei mitnimmt. Nachgestellt: nach dem Commit des Menschen leitet jaira die Commit-Liste vollstaendig ab und der Zug in die Endlane wird nicht verweigert.
  proof: core/role/builtin/jaira-role-lane/SKILL.md:67 — statt 'git commit' die fertige Zeile mit Handle im Betreff und der Ticket-Datei im 'git add'; Begruendung SKILL.md:79
- [x] Der Mensch sieht den Code, bevor darauf aufgebaut wird: der Dispatcher legt ihn nach jedem Inkrement vor und wartet, statt am Ende alles auf einmal zu zeigen.
  proof: core/role/builtin/jaira-role-lane/SKILL.md:50 — 'git diff' nach jedem DoD-Punkt vorlegen und warten; leerer Diff heisst keine Pause
- [x] Je eine Zeile in core/release/NOTES.md unter ## Unreleased fuer den Modus und fuer das, was ein Benutzer beim Committen anders tut.
  proof: core/release/NOTES.md:16 (Modus) und :17 (Committen von Hand), beide unter ## Unreleased

## Options

- [x] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/ticket/schema.go: FieldMode-Konstante, canonicalOrder, Ticket.Mode, Zuweisung in Load
- [x] internal/cli/flow.go: 'mode' als eigener Schluessel im showForLane-JSON, neben model_tier — nicht ueber input-requires
- [x] internal/cli/tickets.go: 'mode' in ticketJSON, damit 'jaira show --json' es fuehrt
- [x] internal/cli/tickets.go newSetCmd: mode nur leer oder 'conversational' annehmen, sonst ExitUsage
- [x] internal/tui: mode in fieldsWithTheirOwnRow (view.go) und in editableFields (edit.go)
- [x] Tests: Feld ueberlebt Round-trip, show --for-lane --json fuehrt mode, set lehnt unbekannten Wert ab
- [x] jaira-dispatcher/SKILL.md: neuer Halt VOR der Plan-Lane — offene Entscheidungen zaehlen, bei >=1 anhalten, Antwort mit 'jaira note' aufs Ticket, dann 'jaira set <id> mode=conversational'
- [x] jaira-dispatcher/SKILL.md: Modus beim Start vom Ticket lesen (jaira show --json), nicht aus dem Kontext — das ist, was den Sitzungsabbruch ueberlebt; --no-worktree in diesem Modus empfehlen
- [x] jaira-role-lane/SKILL.md: im Modus nach jedem DoD-Punkt 'git diff' vorlegen und warten; leerer Diff heisst keine Pause
- [x] jaira-role-lane/SKILL.md: im Modus nicht committen, stattdessen fertige Commit-Zeile mit Handle im Betreff und der Ticket-Datei im git add zurueckgeben (Muster aus jaira-role-pr/SKILL.md)
- [x] docs/AGENTS.md und README.md: mode im Frontmatter dokumentieren
- [x] core/release/NOTES.md unter ## Unreleased: eine Zeile fuer den Modus, eine fuer das Committen von Hand
- [x] go test ./... -race

## Progress
- **2026-09-16 15:28 · Alexander Sacharov** — Brainstorm, Befund aus dem Code — nicht aus der Notiz.

Die Rolle ist seit 7KX89C eine Datei im Repository: core/role/builtin/jaira-dispatcher/SKILL.md, ausgeliefert per 'jaira roles install'. Der Modus ist damit eine Prompt-Aenderung an EINER Datei, kein Go-Code. Die Kopie in ~/.claude ist heute byte-identisch (diff leer), es gibt also noch keine Drift zu heilen.

Was der Prompt heute kann und was nicht: der Abschnitt 'When to stop and hand back' listet drei Haltepunkte — Worker am Freigabedialog, dieselbe Lane dreimal zurueck, Menschen-Lane erreicht. Alle drei sind REAKTIV: der Dispatcher haelt an, nachdem etwas passiert ist. Ein Halt VOR der Plan-Lane, weil im Ticket noch etwas offen ist, existiert nirgends. Das ist die eigentliche Luecke, und sie ist genau die, die DoD-Punkt 1 verlangt.

Der Haken, den die Ticket-Notiz nicht sieht: 'nicht automatisch committen' ist keine Eigenschaft des Dispatchers. Das Committen steht im Worker-Prompt (jaira-role-lane/SKILL.md: 'move the ticket, then git add -A and commit ... whose message names the ticket id'). Ein Modus, der nur im Dispatcher-Prompt steht, erreicht den Worker nie: spawn.sh tippt fest '/jaira-role-lane $ticket $lane' in den Tab und hat keine Stelle, an der ein weiteres Argument mitfaehrt (dieselbe Zeile, die 3YRPXJ aufmacht). Der Modus muss also einen Weg zum Worker haben — das ist die Entscheidung, die dieses Ticket traegt, und sie steht in der naechsten Notiz.

Was NICHT erfunden werden muss: der Kanal zum Menschen gibt es schon. Das Frontmatter-Feld 'question' traegt heute eine Frage an den Menschen (auf 7KX89C im Einsatz), die human-Lane ist der Ort, an dem eine Person antwortet, und 4XHZ6N baut daneben den gesammelten human-note-Kanal. Der Gespraechsmodus soll keinen vierten daneben stellen.
- **2026-09-16 15:28 · Alexander Sacharov** — Brainstorm, Entscheidung 1: wie der Modus den Worker erreicht. Drei Wege, einer empfohlen.

A) Der Modus steht auf dem TICKET. Ein Frontmatter-Feld (Arbeitstitel 'mode: conversational'), gesetzt vom Dispatcher, wenn seine Zaehlung mindestens eine offene Entscheidung findet. 'jaira show --for-lane' reicht es dem Worker mit, so wie es Prompt, Eingabe und Modellstufe schon mitreicht.
  Kostet: ein Feld im Schema, also Go-Code plus NOTES.md-Zeile — der teuerste der drei.
  Gibt: als einziger das, was DoD-Punkt 3 woertlich verlangt ('die Sitzung mittendrin abgebrochen — die Entscheidung ist danach noch da'). Ein frischer Dispatcher nach 'jaira resume' faehrt im selben Modus weiter, weil der Modus auf der Platte steht und nicht in einem Kontext.

B) Der Modus steht in der GETIPPTEN ZEILE. spawn.sh bekommt ein weiteres Argument, im Tab landet '/jaira-role-lane <id> <lane> --conversational'.
  Kostet: haengt an 3YRPXJ, das genau diese fest verdrahtete Zeile aufmacht — und stirbt mit der Sitzung. Ein neu gestarteter Worker laeuft still wieder autonom, und niemand sieht es, bis er committet hat.
  Gibt: kein Go-Code, nur Shell.

C) KEIN Worker fuer in-progress. Der Dispatcher macht die Inkremente im Gespraech selbst.
  Kostet: bricht die Regel, auf der die Rolle steht ('You do not implement'), und verbrennt genau den Kontext, den ein langes Gespraech braucht — der Dispatcher existiert, damit Warten und Logs NICHT im Teamlead-Kontext liegen; hier zieht er sie sich selbst rein.
  Gibt: nichts zu bauen. Ehrlich gesagt ist das, was diese Woche von Hand lief.

Empfehlung: A. Nur A ueberlebt den Sitzungsabbruch, und ohne das ist DoD-Punkt 3 nicht erfuellbar. B ist A ohne Gedaechtnis; C ist gar kein Mechanismus, sondern die Gewohnheit, die dieses Ticket ersetzen soll.

Mitzunehmen, wenn A gebaut wird: die Commit-Zeile, die der Worker statt des Commits zurueckgibt, muss den Handle im Betreff tragen und die Ticket-Datei mit 'git add' nennen. Das Muster steht fertig in jaira-role-pr/SKILL.md — die Rolle gibt die create-Zeile zurueck, wenn ein Agent sie gerufen hat, statt sie auszufuehren. Abschreiben, nicht neu erfinden.
- **2026-09-16 15:29 · Alexander Sacharov** — Brainstorm, Entscheidung 2: in welchen Haeppchen der Mensch den Code sieht. Das ist die Frage, die der Kontext ausdruecklich hierher gegeben hat.

i) EIN DoD-PUNKT = EIN HAEPPCHEN. Der Worker arbeitet einen Punkt, legt den Diff vor, wartet.
  Kostet: ein Ticket mit sechs DoD-Punkten heisst sechs Pausen.
  Gibt: die Einheit existiert schon und ist schon das Inkrement. Sie hat auch schon ihren Beleg — 'jaira dod <n> --done --proof <datei:zeile>' — die Pause hat also von selbst ein Artefakt, das der Mensch anschaut, und nicht nur 'schau mal den Diff an'.

ii) JEDE DATEI-AENDERUNG. Zu fein. Genau davor warnt der Kontext: langsamer, als es selbst zu tippen. Damit erledigt.

iii) EINE LANE = EIN HAEPPCHEN, gezeigt am Ende von in-progress.
  Kostet: das ist der heutige Zustand. Der Mensch sieht alles auf einmal, wenn das Zuruecknehmen am teuersten ist — genau das Versagen, das 0YGWXQ mit sieben critique-Runden bezahlt hat.
  Gibt: nichts, was heute fehlt.

Empfehlung: i, mit einer Bremse dagegen, dass daraus doch eine Pause pro Zeile wird — ein DoD-Punkt, der keinen Diff erzeugt (Dokumentation ohne Code, ein Punkt, der nur eine Nachpruefung beschreibt), ist kein Haeppchen und pausiert nicht. Gezeigt wird, was 'git diff' seit der letzten Pause hergibt, nicht eine Zusammenfassung davon.

Damit ist auch klar, WO der Code liegen soll, waehrend der Mensch ihn anschaut: im Verzeichnis, das er offen hat. spawn.sh kann das seit 7KX89C — '--no-worktree' startet den Worker in der ausgecheckten Arbeitskopie statt in ../.worktrees/<repo>-<slug>. Ein Gespraechsmodus, der den Menschen nach jedem Punkt in ein Worktree schickt, das er suchen muss, verliert genau die Bequemlichkeit, fuer die er da ist. Der Preis steht im Hilfetext von spawn.sh und gilt: mit --no-worktree darf nur ein Worker gleichzeitig laufen — was fuer einen Modus, in dem ohnehin ein Mensch mitliest, keine echte Einschraenkung ist.

Nicht entschieden und auch nicht hier zu entscheiden: ob der Mensch im Tab des Workers antwortet oder ueber die human-Lane. Der Dispatcher-Prompt sagt schon, dass ein Mensch, der in einen Worker-Tab tippt, kein Fehler ist — das reicht fuer den ersten Bau.
- **2026-09-16 15:33 · Alexander Sacharov** — Pre-process, Befund am Code: was das 'mode'-Feld (Brainstorm-Entscheidung A) konkret kostet.

Das Feld ist HEUTE schon schreibbar, ohne eine Zeile Go: 'jaira set' (internal/cli/tickets.go:825) nimmt jedes field=value entgegen, es gibt keine Whitelist. 'jaira set GTQHNH mode=conversational' landet im Frontmatter und ueberlebt den Sitzungsabbruch. Der Grund, trotzdem Go-Code zu schreiben, ist nicht das Schreiben, sondern das LESEN: ein Feld, das nirgends im Schema steht, erscheint in keinem JSON, in keinem TUI-Feld und in keinem Merge-Regelsatz. Der Worker wuerde es nie sehen.

Die Lesekette, die dranhaengt (je eine Stelle, alle klein):
- core/ticket/schema.go: FieldMode-Konstante, canonicalOrder (Zeile 124), Ticket.Mode, Zuweisung in Load (~Zeile 552, neben t.Question)
- internal/cli/flow.go:657 fieldValue — braucht KEINEN neuen Case, der default-Zweig liest jeden Doc-Scalar; ein expliziter Case nur der Lesbarkeit halber
- internal/cli/tickets.go:1287 ticketJSON — 'mode' neben 'model_tier', damit 'jaira show --json' und der Dispatcher es sehen
- internal/tui/view.go:1090 fieldsWithTheirOwnRow und internal/tui/edit.go:17 editableFields — sonst faellt das Feld im TUI in den Sammel-Rest bzw. ist von Hand nicht umstellbar
- core/merge/merge.go: NICHTS zu tun. mode ist ein Skalar, kein Prosa- und kein Listenfeld; die Default-Regel (neuere Antwort gewinnt) ist genau richtig fuer einen Modusschalter.

Wichtigster Befund, weil er eine naheliegende Loesung ausschliesst: der Modus darf NICHT ueber 'input-requires' einer Lane laufen. showForLane (internal/cli/flow.go:552) fuellt 'input' streng aus l.InputRequires, und core/lane prueft beim Laden jedes input-requires gegen ticket.SuppliedFields plus die output-produces frueherer Lanes. 'mode' dort einzutragen hiesse: jede Lane-Datei anfassen UND SuppliedFields aufweichen, dessen Kommentar (schema.go:140ff) ausdruecklich erklaert, warum diese Liste eng bleibt. Richtig ist die Ebene darueber: 'mode' als eigener Schluessel im JSON-Rumpf, neben 'model_tier' und 'prompt' — dort steht schon heute, WIE eine Lane zu fahren ist, nicht WOMIT.
- **2026-09-16 15:33 · Alexander Sacharov** — Pre-process, drei Stellen, an denen der Plan eine Wahl trifft. Je Empfehlung plus Grund.

1. WER LOESCHT DEN MODUS WIEDER? Der Dispatcher setzt ihn vor der Plan-Lane. Wenn ihn niemand raeumt, laeuft das Ticket bis ins Logbuch im Gespraechsmodus und jeder spaetere Worker pausiert nach jedem DoD-Punkt fuer einen Menschen, der laengst weitergezogen ist.
   Empfehlung: NIEMAND loescht ihn automatisch, und das ist Absicht. Der Modus ist eine Aussage ueber das Ticket ('hier waren Entscheidungen offen'), nicht ueber eine Lane. Ein Ticket, das ihn einmal gebraucht hat, braucht ihn in critique und testing genauso. Geraeumt wird von Hand mit 'jaira set <id> mode=' — im Prompt der Menschen-Lane als eine Zeile erwaehnt. Automatik hier waere ein zweiter Mechanismus, der raet.

2. FREIER TEXT ODER GEPRUEFTER WERT? 'jaira set' nimmt heute jeden String. 'mode=conversation', 'mode=konversation', 'mode=chat' waeren alle stumm erfolgreich und alle wirkungslos — der Worker vergleicht gegen genau ein Wort.
   Empfehlung: in 'jaira set' pruefen, erlaubt sind nur leer und 'conversational'. Vier Zeilen Code. Ein Tippfehler, der stumm nichts tut, ist genau die Fehlerklasse, die dieses Ticket ueberhaupt erst gibt: der Mensch glaubt, er laeuft im Gespraechsmodus, und der Worker committet.

3. WORAN MERKT DER WORKER, DASS EIN DoD-PUNKT KEINEN DIFF ERZEUGT? Die Brainstorm-Notiz will die Bremse: ein Punkt ohne Code pausiert nicht. Der Worker kann das nicht vorher wissen, nur hinterher.
   Empfehlung: kein Vorher-Wissen einbauen. Der Worker arbeitet den Punkt, laesst 'git diff' laufen; ist die Ausgabe leer, geht er ohne Pause weiter. Das ist eine Prompt-Zeile und kein Mechanismus, und es kann nicht falsch raten.

Nicht in diesem Ticket, absichtlich: spawn.sh bleibt unveraendert. Der Modus faehrt ueber die Platte, nicht ueber die getippte Zeile — genau deshalb ist 3YRPXJ (spawn.sh nicht erweiterbar) hier KEINE Vorbedingung. Und '--no-worktree', das die Brainstorm-Notiz empfiehlt, ist bereits gebaut; es ist eine Empfehlung im Dispatcher-Prompt, kein Code.
- **2026-09-16 15:38 · Alexander Sacharov** — In-progress, Befund: das 'mode'-Feld braucht zwei Schreibwege, nicht einen.

Der Plan nennt nur 'jaira set' als Pruefstelle. Das TUI schreibt aber an 'jaira set' vorbei: internal/tui/edit.go commitEdit ruft t.Doc().SetScalar direkt. Haette ich die vier Zeilen nur in internal/cli/tickets.go gelegt, waere 'mode' im TUI-Feldeditor frei beschreibbar geblieben und ein dort getippter Tippfehler haette genau den stummen Fehler erzeugt, den Pre-process-Entscheidung 2 verhindern will.

Deshalb steht die Pruefung als ticket.ValidMode in core/ticket/schema.go und wird von beiden Wegen gerufen. Eine Funktion statt zweier Vergleiche gegen die Konstante, damit CLI und TUI nicht auseinanderdriften koennen.

Zweiter Befund, kleiner: in 'jaira set' muss die Pruefung NACH der --append-Faltung stehen, nicht davor. Davor haette 'jaira set <id> mode=conversational --append' auf einem bereits gesetzten Ticket 'conversational\nconversational' geschrieben - jeder Einzelwert gueltig, das Ergebnis nicht. Sie steht jetzt direkt vor SetScalar.
- **2026-09-16 15:43 · Alexander Sacharov** — In-progress, Befund an der Dokumentation: 'mode' ist kein Lane-Eingabefeld und wird deshalb bei den Modellstufen dokumentiert, nicht bei den Lane-Contracts.

docs/AGENTS.md hat jetzt 'A note on the conversational mode' direkt neben 'A note on model tiers'. Das ist derselbe Satz Aussage: jaira reicht den Wert durch und handelt nicht darauf, der Treiber entscheidet. Beim Lane-Abschnitt haette es den falschen Eindruck gemacht, man koenne 'mode' in input-requires eintragen - genau das, was die Pre-process-Notiz ausschliesst.

Offen gelassen, absichtlich, weil ausserhalb dieser Lane: die Rollen-Dateien in ~/.claude sind jetzt aelter als die im Repository. Wer den Modus benutzen will, braucht 'jaira roles install --global --force'; beide NOTES.md-Zeilen sagen das. Ein bereits von Hand editierter Dispatcher wird dabei als skipped gemeldet und nicht ueberschrieben.
- **2026-09-16 15:48 · Alexander Sacharov** — critique: drei Befunde, alle mit klarer Reparatur, daher zurueck nach in-progress.
1. ValidMode (core/ticket/schema.go:150) trimmt, der Schreibpfad (internal/cli/tickets.go:929, internal/tui/edit.go:60) nicht. Empirisch geprueft in einem Wegwerf-Board: 'jaira set <id> "mode= conversational "' -> exit 0, Datei traegt 'mode: " conversational "', 'show --json' liefert ' conversational '. Der Mensch glaubt, im Modus zu sein, der Worker vergleicht auf genau ein Wort und committet wie immer — genau der Fall, gegen den TestSetRefusesUnknownMode geschrieben wurde. Fix in einer Zeile: TrimSpace aus ValidMode streichen, oder den getrimmten Wert speichern.
2. Der Modus ist schreibbar ueber CLI und TUI, aber in keiner Ausgabe fuer Menschen sichtbar — nur in --json. Das Ziel des Tickets ist, dass der Modus einen Sitzungsabbruch ueberlebt; wer ihn gesetzt hat, muss auch sehen koennen, dass er noch an ist. Der Eintrag in fieldsWithTheirOwnRow (internal/tui/view.go:1099) behauptet sogar eine Zeile, die es nicht gibt.
3. --no-worktree hat jetzt zwei Regeln an zwei Stellen derselben Datei, die sich widersprechen. Nebenbei: der Worktree-Zwang steht weder im Goal noch in der DoD — wenn er bleibt, gehoert er in den bestehenden Absatz und nicht in einen zweiten.
Nicht beanstandet und bewusst stehen gelassen: mode als freies Frontmatter-Feld statt als Lane-Input (die Begruendung in flow.go:618-623 traegt), das Fehlen einer Merge-Regel in core/merge (Skalar, juengster Schreiber gewinnt, ist hier richtig), und dass jaira den Modus selbst nicht auswertet — das ist dasselbe Muster wie model-tier.
