---
id: 01M2NCH8ZK9524JZ00J4GTQHNH
title: "Der Dispatcher bekommt einen Gespraechsmodus, statt dass eine zweite Rolle daneben entsteht"
status: testing
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
updated-at: 2026-09-16T20:25:54Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-67097
claimed-at: 2026-09-16T20:17:17Z
outcome-what: "internal/cli/resume.go fuehrt jetzt den Modus: 'mode' in der items-Map des --json-Zweigs und eine 'mode:'-Zeile im Klartext-Block, letztere nur bei gesetztem Modus. Dazu TestResumeCarriesMode in internal/cli/mode_test.go, eine Ergaenzung der bestehenden Modus-Zeile in core/release/NOTES.md und der Verweis auf den Test in der proof-Zeile von DoD-Punkt 3."
outcome-why: "resume.go baut sein JSON von Hand und geht nicht durch ticketJSON, also fehlte das Feld dort. 'jaira resume' ist aber genau die Wiederanlauf-Stelle, auf die sich jaira-dispatcher/SKILL.md:65 beruft — ein Dispatcher nach einem Sitzungsabbruch haette den Modus dort nicht gesehen und waere autonom weitergelaufen, was die in schema.go:44-47 aufgeschriebene Begruendung des Feldes aushebelt."
outcome-resolves: "Befund aus critique-Runde 7: resume.go:126-141 trug kein 'mode', der Klartext-Block ebenso wenig."
review-summary: "none"
review-gaps: "folded core/validate/mode_test.go's two near-identical bad-mode tests into one table test (hand-written 'chat' and untrimmed ' conversational ' asserted the same five things twice) and pulled setMode/modeOf out of internal/cli/mode_test.go, where the same 'set mode, show --json, unmarshal, read the key' block stood four times — 72 test lines gone, no behaviour touched, suite green. Left alone: CanonicalMode is already the one shared write-path check (set, TUI editor, validate all call it), so there is no second implementation to fold; fieldValue's FieldMode case looks unreachable but is not — mode merges through merge.go's default mergeScalar branch, so the merge driver renders it on a conflict; the long doc comments on FieldMode/ModeConversational and the repeated 'a lane that changed no code hands back no line' in README, docs/AGENTS.md, both SKILL.md files and NOTES.md are five different readers, not duplication, and cutting them is an editorial call rather than a cleanup; NOTES.md's mode line is very long but the format is one line per change. No dead code and no hot-path cost found."
test-verdict: "pass: go build/vet/test ./... -race -count=1 green (RC=0, 29 Pakete), DoD 1-6 in der Working Tree geprueft, Verhalten auf einem Scratch-Board mit dem gebauten Binary durchgespielt"
---

# Der Dispatcher bekommt einen Gespraechsmodus, statt dass eine zweite Rolle daneben entsteht

## Definition of Done

- [x] Der Eintritt in den Modus haengt an einer nachpruefbaren Bedingung, nicht am Bauchgefuehl: vor der Plan-Lane zaehlt der Dispatcher die noch offenen Entscheidungen des Tickets auf. Keine offene - er laeuft weiter wie heute. Mindestens eine - er haelt an und fragt den Menschen.
  proof: core/role/builtin/jaira-dispatcher/SKILL.md:44 — Schritt 1 zaehlt aus Ticket UND Notes, eine beantwortete Entscheidung gilt als geschlossen
- [x] Es gibt weiterhin genau eine Dispatcher-Rolle. Kein zweiter Skill und keine zweite Kommandozeile daneben; der Modus steht im selben Prompt.
  proof: core/role/builtin/ enthaelt unveraendert sieben Rollen; der Modus steht in jaira-dispatcher/SKILL.md:31 und jaira-role-lane/SKILL.md:39, kein neuer Skill und keine neue Kommandozeile
- [x] Was im Gespraech entschieden wird, steht mit 'jaira note' auf dem Ticket, BEVOR die Arbeit daran beginnt - nicht hinterher. Nachgestellt an einem Ticket, dessen Sitzung mittendrin abgebrochen wird: die Entscheidung ist danach noch da.
  proof: core/role/builtin/jaira-dispatcher/SKILL.md:54 (jaira note vor Arbeitsbeginn) plus TestModeSurvivesRoundTrip, TestShowPrintsModeForPeople und TestResumeCarriesMode in internal/cli/mode_test.go
- [x] In diesem Modus committet der Agent nicht selbst. Er legt die Aenderungen bereit und gibt eine fertige Commit-Zeile zurueck, die den Ticket-Handle im Betreff traegt und die Ticket-Datei mitnimmt. Nachgestellt: nach dem Commit des Menschen leitet jaira die Commit-Liste vollstaendig ab und der Zug in die Endlane wird nicht verweigert.
  proof: core/role/builtin/jaira-role-lane/SKILL.md:75 — statt 'git commit' die fertige Zeile mit Handle im Betreff und der Ticket-Datei im 'git add'; Begruendung SKILL.md:84; die Gegenseite in core/role/builtin/jaira-dispatcher/SKILL.md:84 — Zeile nur bei Code-Aenderung, sonst keine und die Lane ist trotzdem fertig
- [x] Der Mensch sieht den Code, bevor darauf aufgebaut wird: der Dispatcher legt ihn nach jedem Inkrement vor und wartet, statt am Ende alles auf einmal zu zeigen.
  proof: core/role/builtin/jaira-role-lane/SKILL.md:55 — 'git diff' nach jedem DoD-Punkt vorlegen und warten; leerer Diff heisst keine Pause
- [x] Je eine Zeile in core/release/NOTES.md unter ## Unreleased fuer den Modus und fuer das, was ein Benutzer beim Committen anders tut.
  proof: core/release/NOTES.md:17 (Modus, inkl. der Kopfzeile von 'jaira show --for-lane') und :18 (Committen von Hand), beide unter ## Unreleased

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
- [x] critique 4: dieselbe Bedingung (Commit-Zeile nur bei Code-Aenderung) an den vier verbleibenden Stellen — dispatcher/SKILL.md, docs/AGENTS.md, README.md, core/ticket/schema.go
- [x] critique 5: 'mode' in core/validate/validate.go pruefen (CodeBadMode, Warning), plus NOTES.md-Zeile
- [x] critique 6: validate.go meldet auch den untrimmed Wert (canon != t.Mode), plus Testfall
- [x] critique 7: 'jaira resume' fuehrt mode — JSON-items und Klartext-Block in internal/cli/resume.go, plus Testfall

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
- **2026-09-16 15:51 · Alexander Sacharov** — In-progress nach critique: die drei Befunde repariert, einer davon anders als critique vorgeschlagen hat.

Befund 1 (Trim-Divergenz) ist NICHT durch Streichen des TrimSpace repariert, obwohl critique das als erste Option nannte. Streichen haette 'jaira set <id> "mode= conversational"' zu einem Fehler gemacht - eine Schreibweise, die in der Shell naheliegt und nichts Falsches meint. Stattdessen gibt ValidMode jetzt als CanonicalMode den getrimmten Wert MIT dem Urteil zurueck, und beide Schreibwege speichern genau das. Damit kann kein Aufrufer mehr pruefen und danach etwas anderes speichern - derselbe Grund, aus dem die Funktion ueberhaupt existiert (Notiz 15:38), eine Ebene weitergezogen. Ein 'bool'-Rueckgabewert laedt zu genau diesem Fehler ein; ein (value, ok) nicht.

Befund 2: row("mode", ...) steht in beiden Detail-Panes neben row("tier", ...), weil beide dieselbe Art Aussage sind - WIE das Ticket gefahren wird, nicht was drinsteht. row() ueberspringt Leeres, ein Ticket ohne Modus kostet keine Zeile; TestShowPrintsModeForPeople prueft beide Richtungen.

Befund 3: die --no-worktree-Regel steht jetzt nur noch im bestehenden Absatz, dort als dritter Fall neben den zwei vorhandenen. Der Bullet im Modus-Abschnitt verweist darauf. Auch Schritt 2 der Schleife ('in its own worktree') ist angefasst - er war die dritte, von critique nur nebenbei erwaehnte Stelle, und ohne ihn haette der Verweis auf einen Absatz gezeigt, dem die Schleife widerspricht.
- **2026-09-16 15:57 · Alexander Sacharov** — critique (2. Durchgang): zwei Befunde, beide klar reparierbar, keine Entscheidung für den Menschen.

1. internal/cli/flow.go:635 — die Klartextausgabe von 'jaira show --for-lane' trägt den Tier in der Kopfzeile, den Modus nicht; nur --json hat ihn. Der vorige Durchgang hat den Modus in die Detail-Panes geholt (jaira show, TUI); die Lane-Prompt-Ausgabe, aus der ein Worker tatsächlich seine Anweisung liest, blieb dabei aus. Ein bash-fähiger Agent ohne --json läuft damit autonom weiter, ohne dass es jemand merkt — dieselbe stille Fehlerart, gegen die CanonicalMode angetreten ist. Modus an dieselbe Stelle wie den Tier, nur wenn nicht leer.

2. core/role/builtin/jaira-dispatcher/SKILL.md, Schritt 1 des Zählens — kennt keinen Zustand 'bereits beantwortet'. Reihenfolge ist Notes (Schritt 4), dann Modus (Schritt 5); stirbt der Dispatcher dazwischen, liest der nächste keinen Modus, zählt neu und fragt den Menschen dasselbe noch einmal. Die Antworten stehen zu dem Zeitpunkt schon in den Notes — sie werden nur nicht gelesen. Schritt 1 muss die Notes mitlesen und eine beantwortete Entscheidung als geschlossen zählen.

Geprüft und nicht beanstandet: mode fällt im Merge in den default-Zweig (mergeScalar, neuerer gewinnt) — für einen geschlossenen Skalar richtig, kein Eintrag in listFields/proseFields nötig. Der neue case in fieldValue() ist gegenüber dem default-Zweig redundant, folgt aber dem Muster aller anderen Felder dort; das aufzubrechen wäre teurer als der Gewinn. Die doppelte Prüfung in tickets.go und tui/edit.go teilt sich CanonicalMode — der Rest ist je eigene Fehlerausgabe und lässt sich nicht sinnvoll zusammenlegen.
- **2026-09-16 16:00 · Alexander Sacharov** — In-progress nach critique (2. Durchgang): beide Befunde repariert, einer davon enger als critique ihn beschrieben hat.

Befund 1 (Modus fehlt in der Klartext-Kopfzeile): die Zeile wird jetzt vorgebaut, statt einen zweiten Fprintf-Zweig zu bauen. Grund: die Kopfzeile ist ein Format-String mit fester Klammer '(tier: %s)'. Ein zweiter Zweig fuer 'mit Modus' haette den Rest der Zeile dupliziert, und die naechste Aenderung an der Zeile haette nur einen der beiden getroffen - dieselbe Drift-Falle, aus der CanonicalMode entstanden ist. Jetzt gibt es genau eine Stelle, an der die Klammer steht.

Was ich dabei geprueft und verworfen habe: den Modus nicht in die Kopfzeile, sondern als eigene '**mode**'-Zeile zu den Input-Feldern. Falsch, weil die Input-Schleife streng ueber l.InputRequires laeuft - der Modus stuende dann zwischen Feldern, die aus dem Lane-Contract kommen, und suggeriert, er sei eins. Genau das schliesst die Pre-process-Notiz (15:33) aus. Die Kopfzeile sagt schon heute, WIE die Lane zu fahren ist.

Befund 2 (Zaehlen kennt kein 'schon beantwortet'): die Reparatur benennt in Schritt 1 ausdruecklich den Fall, der sie noetig macht - Dispatcher stirbt zwischen Schritt 4 (Notes) und Schritt 5 (Modus). Ohne den genannten Fall liest ein Modell 'Notes mitlesen' als Hoeflichkeit und zaehlt trotzdem neu. Der Prompt muss den Grund tragen, weil niemand sonst ihn traegt.

Nicht angefasst, absichtlich: die Reihenfolge Notes-dann-Modus bleibt. Sie erst den Modus setzen zu lassen waere die andere denkbare Reparatur, aber dann traegt ein Ticket den Modus, bevor die Antworten darauf stehen - ein Worker, der genau dazwischen startet, pausiert nach jedem DoD-Punkt fuer Entscheidungen, die niemand aufgeschrieben hat. Die Notes zuerst ist die richtige Reihenfolge; das Lesen war die Luecke.
- **2026-09-16 16:05 · Alexander Sacharov** — critique (3. Durchgang): ein Befund, ein Fix.

jaira-role-lane/SKILL.md, Abschnitt 2 des Modus ('Do not commit. Hand back the commit line instead.') widerspricht der Regel 30 Zeilen darueber: 'did you change no code? Then commit nothing.' Der Modus haengt laut schema.go am TICKET und gilt ausdruecklich 'through critique and testing too' — also liest ihn auch ein Worker in critique, testing oder review. Der aendert keine Zeile Code, bekommt aber trotzdem 'Write the command out ready to paste' und haendigt dem Menschen 'git add <files> .jaira/tickets/<ticket>.md; git commit -m ...' aus. Fuegt der Mensch das ein, entsteht ein Commit, der nur die Ticket-Datei traegt — die Buchhaltung, die 9ZZSFT gerade abgeschafft hat. Und hier trifft es haerter als sonst: im Gespraechsmodus committet ein Mensch, der die Zeile als fertig gedacht hinnimmt.

Abschnitt 1 hat die Ausnahme bereits ('Empty output? No pause' — ein Doku-Punkt ohne Diff pausiert nicht). Abschnitt 2 hat sie nicht. Fix: dieselbe Bedingung in Abschnitt 2 — eine Commit-Zeile nur, wenn Code geaendert wurde; ohne Code-Aenderung gibt der Worker keine Zeile zurueck und laesst die Ticket-Datei im Worktree fuer den naechsten Commit, der Code traegt.

Nicht beanstandet, damit es nicht noch einmal aufgemacht wird: die Inline-Pruefung 'if k == ticket.FieldMode' in tickets.go:933 steht neben den schon vorhandenen Sonderfaellen FieldID und FieldStatus in derselben Funktion — das ist das bestehende Muster, keine neue Form. CanonicalMode mit zwei Aufrufern (CLI und TUI) ist keine Abstraktion mit einem Aufrufer. 'mode' neben 'model_tier' statt in input ist im Code begruendet und stimmt. Und 'jaira show --json, notes included' im Dispatcher-Prompt traegt: tickets.go:655 haengt den ganzen Body an, und die Notes stehen unter '## Progress' darin.
- **2026-09-16 16:06 · Alexander Sacharov** — Dispatcher-Halt: critique hat dreimal zurueckgeschickt (Runde 1: drei Befunde, Runde 2: zwei, Runde 3: einer). Die Schleife wird kleiner, aber sie endet nicht, weil jeder Durchgang tiefer liest als der vorige. Regel im Dispatcher-Prompt: nach drei Rueckgaben entscheidet der Mensch, nicht der Dispatcher.

Offen und unrepariert, der einzige Befund aus Runde 3: core/role/builtin/jaira-role-lane/SKILL.md:68-84. Abschnitt 2 des Modus ('Do not commit. Hand back the commit line instead.') hebt die Regel aus Zeile 31-37 auf ('did you change no code? Then commit nothing'). Der Modus steht auf dem Ticket und gilt damit auch in critique, testing und review, wo kein Code entsteht; der Worker reicht dem Menschen trotzdem eine git-commit-Zeile, und der Mensch committet reine Ticket-Buchhaltung. Abschnitt 1 hat die passende Ausnahme schon ('Empty output? No pause'), Abschnitt 2 nicht.

Stand des Codes: zwei Commits, 9cb1df9 und f2c79a9, go test ./... -race gruen, alle sechs DoD-Punkte getickt und belegt.

Ausserhalb des Tickets: die Rollenkopien in ~/.claude sind aelter als das Repository. Wer den Modus benutzen will, braucht 'jaira roles install --global --force'.
- **2026-09-16 19:45 · Alexander Sacharov** — Entscheidung von Alex am 16.09., vor der Arbeit aufgeschrieben — sie schliesst den einen offenen Befund aus critique-Runde 3 (core/role/builtin/jaira-role-lane/SKILL.md:68-84).

Die Commit-Zeile, die der Agent im Gespraechsmodus zurueckgibt, traegt IMMER Code. Eine Lane, die keinen Code geaendert hat — critique, testing, review —, gibt gar keine Zeile zurueck: dort gilt weiter die Regel aus Zeile 31-37 ('did you change no code? Then commit nothing'), und die Ticket-Datei wartet wie bisher auf den naechsten Commit, der Code traegt. Abschnitt 2 des Modus bekommt damit dieselbe Ausnahme, die Abschnitt 1 mit 'Empty output? No pause' schon hat.

Was ausdruecklich NICHT gebaut wird: ein Verbot fuer den Menschen. Committet eine Person von Hand etwas ohne Code, ist das ihre Sache und nichts, was der Prompt untersagt. Geregelt wird nur, was der Agent vorschlaegt.

Zwei weitere Punkte, die Alex am selben Tag entschieden hat:

Der Vorbehalt aus dem Kontext ('Nicht vor dem naechsten Release') ist erledigt: das Release kommt spaeter, das Ticket wartet nicht darauf.

Die Rollenkopien in ~/.claude werden aktualisiert, nachdem dieser Befund gefixt ist. Geprueft am 16.09.: das Repository ist ueberall neuer, in den globalen Kopien steht keine Zeile, die im Repository fehlt — jaira-role-pr ist durch den A3R6YC-Merge auseinandergelaufen, jaira-dispatcher und jaira-role-lane durch dieses Ticket. Es geht also nichts verloren. Uebrig in ~/.claude/skills/jaira-dispatcher/scripts sind nur drei alte spawn.sh.*.bak.
- **2026-09-16 19:47 · Alexander Sacharov** — In-progress nach critique (3. Durchgang): der eine offene Befund ist repariert, nach Alex' Entscheidung vom 16.09.

Abschnitt 2 des Gespraechsmodus (core/role/builtin/jaira-role-lane/SKILL.md) hat jetzt dieselbe Ausnahme wie Abschnitt 1: keine Code-Aenderung, keine Commit-Zeile. Damit gilt die Regel aus Zeile 31-37 auch im Modus, und critique/testing/review reichen dem Menschen nichts mehr aus, dessen Einfuegen einen reinen Ticket-Commit erzeugt.

Warum die Ausnahme mit Begruendung im Prompt steht und nicht als ein Halbsatz: der Grund traegt sie. Ein Modell, das nur 'nur bei Code-Aenderung' liest, wiegt das gegen den umgebenden Absatz ab, der ausfuehrlich erklaert, warum die Zeile wichtig ist, und gibt sie im Zweifel doch heraus. Der Satz 'a line written out ready to paste reads as already decided' ist das, was im Zweifel den Ausschlag gibt — genau wie bei Befund 2 aus Runde 2, wo der genannte Fall erst aus einer Hoeflichkeit eine Regel gemacht hat.

Nicht gebaut, wie Alex entschieden hat: kein Verbot fuer den Menschen. Der Prompt regelt nur, was der Agent vorschlaegt.

NOTES.md: die vorhandene Modus-Zeile ist erweitert statt einer zweiten. Die Ausnahme ist keine eigene Aenderung, sondern die Bedingung derselben — zwei Zeilen haetten dem Leser zwei Mechanismen vorgesetzt, wo einer steht.
- **2026-09-16 19:49 · Alexander Sacharov** — critique (4. Durchgang): ein Befund, klarer Fix, keine Entscheidung fuer den Menschen.

Der Befund aus Runde 3 ist repariert — Abschnitt 2 in jaira-role-lane/SKILL.md hat jetzt die Ausnahme, und NOTES.md nennt sie mit. Dieselbe Regel steht aber an vier weiteren Stellen unveraendert bedingungslos da: core/role/builtin/jaira-dispatcher/SKILL.md:85, docs/AGENTS.md:256, README.md:169, core/ticket/schema.go:41.

Behaviour, nicht nur Doku: der Dispatcher startet auch critique-, testing- und review-Worker. Die geben nach dem Fix richtigerweise keine Commit-Zeile mehr zurueck — und sein Prompt sagt ihm, eine komme und er solle sie unveraendert an den Menschen weiterreichen. Ein Dispatcher, der die Zeile erwartet, fragt den Worker danach oder haelt die Lane fuer unfertig. Das ist genau der Widerspruch aus Runde 3, nur eine Datei weiter.

Fix an allen vier Stellen dieselbe Bedingung wie in jaira-role-lane/SKILL.md:72-78: Commit-Zeile nur bei Code-Aenderung, sonst keine und die Ticket-Datei bleibt im Worktree.
- **2026-09-16 19:54 · Alexander Sacharov** — In-progress nach critique (4. Durchgang): die Bedingung steht jetzt an allen fuenf Stellen, nicht nur in jaira-role-lane.

Repariert: core/role/builtin/jaira-dispatcher/SKILL.md (zwei Punkte statt einem — 'Worker mit Code-Aenderung gibt eine Zeile zurueck' und 'Worker ohne Code-Aenderung gibt keine, und die Lane ist trotzdem fertig'), docs/AGENTS.md, README.md, core/ticket/schema.go.

Warum der Dispatcher zwei Punkte bekommt und nicht einen Nebensatz: die Gefahr dort ist nicht, dass er eine falsche Zeile weiterreicht, sondern dass er eine fehlende Zeile fuer einen Fehler haelt — den Worker noch einmal fragt oder die Lane fuer unfertig haelt. Ein Nebensatz an der bestehenden Zeile haette nur gesagt, wann eine Zeile kommt; der zweite Punkt sagt ausdruecklich, dass ihr Fehlen die Regel bei der Arbeit ist. Das ist der Satz, der den Dispatcher davon abhaelt, die Lane offen zu halten.

Kein Test dazu: alle fuenf Stellen sind Prosa (vier Dokumente und ein Go-Kommentar). Es gibt im Repository keinen Test, der Prompt-Text auf Aussagen prueft, und einer, der auf Formulierungen matcht, bricht beim naechsten Umschreiben, ohne je ein Verhalten geschuetzt zu haben.

NOTES.md: nichts ergaenzt. Die Modus-Zeile nennt die Ausnahme schon (seit Runde 3) und verweist auf 'jaira roles install --global --force' fuer beide Prompts; die Dispatcher-Aenderung faellt darunter. Eine zweite Zeile haette dieselbe Regel zweimal beschrieben.
- **2026-09-16 19:58 · Alexander Sacharov** — critique (5. Durchgang): ein Befund, klarer Fix, keine Entscheidung fuer den Menschen.

core/validate/validate.go kennt den neuen Wert 'mode' nicht. Die Pruefung des geschlossenen Wertebereichs (leer oder 'conversational') steht an zwei Stellen: internal/cli/tickets.go:928 im 'set'-Pfad und internal/tui/edit.go:57 im Editor der Detailspalte. Beide sind Schreibwege durch die Vordertuer.

Die Hintertueren sind nicht abgedeckt, und internal/cli/validate.go:22-26 zaehlt sie selbst auf: 'damage from a hand edit, a bad merge, or an agent writing something unexpected'. Ein von Hand ins Frontmatter geschriebenes 'mode: chat' laedt fehlerfrei, erscheint in der 'mode'-Zeile von 'jaira show' und im JSON, und der Worker — der laut core/ticket/schema.go:147-151 gegen genau ein Wort vergleicht — laeuft autonom weiter. Der Mensch glaubt, er werde vor jedem Inkrement gefragt, und der Agent committet wie immer. Das ist woertlich das Versagen, das der Kommentar an schema.go:147-151 als Daseinsgrund der Pruefung nennt.

Warum das kein Randfall ist: das Dateiformat ist hier die API (CLAUDE.md: 'the file format IS the API, and must stay hand-editable'). Handedit ist ein vorgesehener Schreibweg, nicht ein Missbrauch. Dazu kommt der Merge-Driver: --take-theirs schreibt den Gegenwert mit SetScalar durch (internal/cli/mergedriver.go:223), ohne CanonicalMode, und 'mode' faellt im Merge in den default-Zweig mergeScalar (core/merge/merge.go:158), kann also ueberhaupt konfliktieren.

Fix: in core/validate/validate.go ein CodeBadMode neben CodeBadTag, geprueft mit ticket.CanonicalMode(t.Mode), Field ticket.FieldMode. SeverityWarning aus demselben Grund, den der Kommentar bei bad_tag schon aufschreibt — das Ticket selbst ist heil, nur der Wert ist unerreichbar —, und die Meldung nennt die Reparatur ('jaira set <handle> mode=conversational' oder 'mode='). CanonicalMode existiert bereits und ist genau die eine Stelle, die beide Pakete lesen sollen; validate ist der dritte Leser, nicht eine dritte Liste.

Dazu eine Zeile in core/release/NOTES.md: ein neuer validate-Problemcode aendert Ausgabe und, mit --strict, den Exit-Code.

Nicht erneut aufgemacht: die Commit-Zeilen-Bedingung aus Runde 3 und 4 steht jetzt an allen fuenf Stellen und stimmt ueberein. Der --no-worktree-Absatz, die Trim-Entscheidung und die Platzierung von 'mode' neben 'model_tier' statt in input-requires bleiben stehen — alle drei sind in frueheren Runden entschieden worden. Der FieldMode-Zweig in flow.go:685 fieldValue ist NICHT tot, obwohl der Kommentar darueber sagt, mode gehe nicht durch input-requires: der Merge-Driver ruft dieselbe Funktion (mergedriver.go:260) beim Auflisten konfliktierter Felder auf.
- **2026-09-16 20:01 · Alexander Sacharov** — In-progress nach critique (5. Durchgang): die Wertebereichspruefung fuer 'mode' steht jetzt auch in core/validate/validate.go.

- core/validate/validate.go: CodeBadMode neben CodeBadTag, geprueft mit ticket.CanonicalMode, SeverityWarning aus demselben Grund wie bad_tag — das Ticket selbst ist heil. Die Meldung nennt beide Reparaturen ('mode=conversational' und 'mode=').
- Die Stelle ist bewusst VOR der blocked-by-Schleife und direkt hinter der Tag-Schleife: beide pruefen einen geschriebenen Wert gegen einen geschlossenen Bereich, und wer den einen Kommentar liest, findet den anderen.
- Der Merge-Driver (internal/cli/mergedriver.go:223, --take-theirs) bleibt ungeprueft — absichtlich. Er ist ein Konfliktloeser, kein Schreibweg mit eigener Meinung; ein durchgereichtes 'mode: chat' faengt jetzt 'jaira validate' ab, genau wie jeden anderen Handedit. Eine Pruefung im Driver waere eine zweite Wahrheit ueber den Wertebereich.
- core/validate/mode_test.go: der unbekannte Wert wird als Warnung mit Feld und Reparatur gemeldet; leer und 'conversational' bleiben stumm.
- Kein Doku-Update noetig: weder docs/ noch README zaehlen validate-Codes auf (grep bad_tag findet dort nichts). Die NOTES.md-Zeile traegt die Aenderung.
- **2026-09-16 20:03 · Alexander Sacharov** — critique (6. Durchgang): ein Befund, klarer Fix, keine Entscheidung fuer den Menschen.

core/validate/validate.go:203 ruft ticket.CanonicalMode(t.Mode) auf und wertet nur das zweite Rueckgabewert aus. CanonicalMode trimmt aber, BEVOR es urteilt — deshalb ist ' conversational ' fuer die Pruefung gueltig, und der ungetrimmte Wert bleibt auf dem Ticket stehen.

Nachgemessen mit einer Probe gegen ParseDoc/Decode: 'mode: " conversational "' in der Frontmatter ergibt t.Mode = " conversational " (Doc.Scalar liest bei einem quoted scalar token.Value, das die Leerzeichen behaelt), validate schweigt, und flow.go gibt den Wert unveraendert als JSON-Feld 'mode' und in der Kopfzeile '(tier: ..., mode: ...)' aus. Der Worker vergleicht gegen genau ein Wort und laeuft autonom — der Mensch glaubt, er werde vor jedem Inkrement gefragt.

Das ist nicht der Befund aus Runde 1 noch einmal: den haben die beiden SCHREIBpfade (jaira set, TUI-Editor) repariert, indem sie den kanonischen Wert speichern. validate ist der LESEpfad fuer handgeschriebene Dateien und die einzige Stelle, die diese Datei je zu sehen bekommt — genau der Fall, fuer den die Pruefung in Runde 5 ueberhaupt dazukam.

Reparatur: in validate.go zusaetzlich melden, wenn canon != t.Mode — gleicher CodeBadMode, gleiche Severity, die Meldung nennt schon beide Reparaturen. Dazu ein Fall in core/validate/mode_test.go neben 'chat'; dort steht heute nur ein Wert ausserhalb der Menge, keiner mit Rand-Leerzeichen.
- **2026-09-16 20:06 · Alexander Sacharov** — In-progress nach critique (6. Durchgang): validate.go meldet jetzt auch den nur-Whitespace-Abweichler.

core/validate/validate.go: die Bedingung ist `!ok || canon != t.Mode`, nicht mehr nur `!ok`. CanonicalMode gibt die getrimmte Form zurueck UND ein ok; der alte Code warf die Form weg. ' conversational ' kam damit durch die Pruefung und untrimmed bei flow.go (JSON 'mode') und in der Kopfzeile an, wo der Worker gegen genau ein Wort vergleicht.

Dieselbe Meldung, derselbe Code, dieselbe Reparatur — bewusst kein zweiter Code. Fuer den Leser ist es ein Fehler ('das steht so nicht im Ticket'), und zwei Codes fuer eine Reparatur waeren eine Unterscheidung, die niemand braucht. Das %q in der Meldung zeigt die Anfuehrungszeichen, also sieht man die Leerzeichen.

NOTES.md: keine neue Zeile, sondern die bestehende validate-Zeile erweitert — es ist dieselbe Pruefung, und eine zweite Zeile daneben liest sich wie ein zweites Feature. Die Zeile nennt jetzt ausdruecklich, dass nur 'jaira set' trimmt.

core/validate/mode_test.go: TestUntrimmedModeIsReportedWithTheRepair neben dem 'chat'-Fall.
- **2026-09-16 20:09 · Alexander Sacharov** — critique (7. Durchgang): ein Befund, klarer Fix, keine Entscheidung fuer den Menschen.

Befund: 'jaira resume' traegt den Modus nicht.

internal/cli/resume.go:126-141 baut seine JSON-Items als eigenes Map-Literal (id, handle, title, status, reason, current_step, goal, notes) statt ueber ticketJSON. Jeder andere Lesepfad ist versorgt — ticketJSON (tickets.go:1324), show --for-lane (flow.go), der Detail-Block in printDetail, die TUI-Zeile. resume ist der einzige, der sein eigenes Literal baut, und genau deshalb durchgefallen. Der Klartext-Zweig darunter (Zeile 153-165) druckt Status, Grund, 'was on' und die Notizen — auch dort keine mode-Zeile.

Warum das mehr ist als eine fehlende Ausgabe: core/ticket/schema.go:44-47 begruendet die ganze Bauform damit, dass der Modus einen Sitzungsabbruch ueberlebt, und core/role/builtin/jaira-dispatcher/SKILL.md:65 sagt woertlich, ein frischer Dispatcher lese ihn 'after jaira resume ... back off disk'. Das stimmt nicht. Genau der Dispatcher, der nach einem Abbruch neu startet und seiner eigenen Anweisung aus Zeile 28 folgt ('reads jaira resume and carries on'), sieht den Modus dort nicht — und laeuft autonom weiter, was der Modus verhindern soll. Gerettet wird es nur dadurch, dass der Abschnitt 30 Zeilen weiter oben zusaetzlich 'jaira show <id> --json' vorschreibt; der Prompt widerspricht sich also selbst.

Fix: "mode": i.t.Mode in die items-Map (i.t ist bereits das per s.Load nachgeladene volle Ticket, Zeile 133) und eine mode-Zeile in den Klartext-Block, wie printDetail sie neben 'tier' hat. Zwei Zeilen. Alternative waere, SKILL.md:65 zu streichen und dort nur show --json zu nennen — das laesst aber die Begruendung des Feldes (Abbruch ueberleben) an einer Stelle haengen, die der Dispatcher-Prompt selbst nicht als Wiederanlauf fuehrt. Empfehlung daher: resume versorgen.
- **2026-09-16 20:12 · Alexander Sacharov** — In-progress nach critique (7. Durchgang): 'jaira resume' fuehrt jetzt den Modus.

- internal/cli/resume.go baut sein JSON von Hand (items-Map, ~Zeile 138) und geht NICHT durch ticketJSON. Jedes neue Frontmatter-Feld muss dort einzeln nachgetragen werden — das ist die Stelle, die bei einem neuen Feld still zurueckbleibt, und sie ist nicht die einzige ihrer Art in internal/cli.
- Der Klartext-Zweig bekommt die Zeile nur bei gesetztem Modus, wie 'was on:' und wie row() in tickets.go: ein Ticket ohne Modus kostet keine Zeile. Der Test prueft beide Richtungen.
- Warum ueberhaupt hier und nicht nur in 'jaira show --json': jaira-dispatcher/SKILL.md:65 nennt 'jaira resume' ausdruecklich als die Stelle, an der ein frischer Dispatcher den Modus von der Platte liest. Die Alternative aus critique — den Satz in SKILL.md streichen — laesst die Begruendung des Feldes (schema.go:44-47: es sitzt auf dem Ticket, um einen Sitzungsabbruch zu ueberleben) unbelegt, also wurde der Code an die Doku angeglichen und nicht umgekehrt.
- Die NOTES.md-Zeile zum Modus wurde ergaenzt statt eine zweite geschrieben: es ist dieselbe Aenderung fuer den Leser (wo sehe ich den Modus), und zwei Zeilen zum selben Feld lesen sich wie zwei Features.
- **2026-09-16 20:16 · Alexander Sacharov** — critique (8. Durchgang): kein Befund. Das ist der Ausgang der Schleife, nicht ein Durchwinken.

Nachgeprueft, jede Stelle einzeln, und alle tragen 'mode':
- Schreibwege: 'jaira set' (internal/cli/tickets.go:935, nach der --append-Faltung, speichert den kanonischen Wert) und der TUI-Editor (internal/tui/edit.go:61). Kein dritter. 'jaira create' setzt model-tier (tickets.go:245), aber keinen Modus — richtig so, der Dispatcher setzt ihn erst vor der Plan-Lane. 'jaira set' kennt keine Feld-Aliase, k wird nur getrimmt, also kann die Pruefung nicht ueber eine andere Schreibweise des Schluessels umgangen werden.
- Lesewege: ticketJSON (tickets.go:1327) und damit 'jaira show --json', 'jaira list --json', 'jaira next --json', 'jaira move --json', 'jaira dod --json' — alle gehen durch dieselbe Funktion. showForLane fuehrt es in JSON und in der Klartext-Kopfzeile (flow.go). printDetail und die TUI-Detailspalte haben die Zeile. 'jaira resume' war der letzte Handbau-Fall und ist seit 28d53bd versorgt, JSON wie Klartext. 'jaira claim --json' ist eine Quittung, kein Ticket — dort gehoert es nicht hin.
- Hintertueren: 'jaira validate' faengt jeden handgeschriebenen Wert ab, auch den nur durch Leerzeichen abweichenden. Der Merge-Driver bleibt bewusst ungeprueft (Runde 5), mode faellt in core/merge zu Recht in den default-Zweig — proseFields (merge.go:50-59) fuehrt question, nicht mode, und das ist fuer einen geschlossenen Skalar richtig.
- flow.go:685 (case FieldMode in fieldValue) ist nicht tot: der Merge-Driver ruft dieselbe Funktion beim Auflisten konfliktierter Felder.
- Die Commit-Zeilen-Bedingung steht an allen fuenf Stellen gleichlautend (jaira-role-lane, jaira-dispatcher, docs/AGENTS.md, README.md, core/ticket/schema.go).

Nicht erneut aufgemacht, weil in frueheren Runden entschieden: mode neben model_tier statt in input-requires (Pre-process 15:33, Runde 2), das Trimmen in CanonicalMode statt Ablehnen (Runde 1), kein automatisches Raeumen des Modus (Pre-process), kein Test auf Prompt-Text (Runde 4).

Ausserhalb dieses Tickets und bewusst nicht hier repariert, gehoert in ein eigenes Ticket: internal/cli baut an mehreren Stellen JSON von Hand statt ueber ticketJSON (resume.go war einer davon und hat Runde 7 gekostet). Ein neues Frontmatter-Feld bleibt dort still zurueck, und nur eine critique-Runde findet es.
- **2026-09-16 20:21 · Alexander Sacharov** — optimize: the only duplication worth removing was in the tests. core/validate/mode_test.go asserted the same five properties twice (unknown value, untrimmed value) - now one table test. internal/cli/mode_test.go repeated the set/show --json/unmarshal block four times - now setMode and modeOf. -72 lines, suite green. Checked and rejected as cleanups: fieldValue's FieldMode case is NOT dead (mode takes merge.go's default mergeScalar path, so the merge driver prints it on a conflict - do not delete it); the prose repeated across README, AGENTS.md, both SKILLs and NOTES.md serves five different readers and cutting it would change what a reader is told, which is not this lane's call.
- **2026-09-16 20:25 · Alexander Sacharov** — testing: drei Durchgaenge, alle gruen.

Gates: 'go build ./...', 'go vet ./...' und 'go test ./... -race -count=1' — RC=0, 29 Pakete ok, kein Skip. Kein Taskfile im Repo, also die Go-Werkzeuge direkt.

Die Forderung, Punkt fuer Punkt an der Working Tree geprueft, nicht am Outcome-Text:
1) jaira-dispatcher/SKILL.md:31-66 — Schritt 1 zaehlt die offenen Entscheidungen aus Ticket UND Notes, Schritt 2 laeuft bei null weiter, Schritt 3 haelt bei >=1 an. Bedingung nachpruefbar, nicht Bauchgefuehl.
2) 'ls core/role/builtin/' zeigt unveraendert sieben Rollen; der Modus steht in denselben zwei Prompts, kein neuer Skill, keine neue Kommandozeile.
3) Schritt 4 schreibt die Antwort mit 'jaira note' VOR Arbeitsbeginn; der Sitzungsabbruch ist durch TestModeSurvivesRoundTrip, TestShowPrintsModeForPeople und TestResumeCarriesMode belegt — alle drei einzeln mit -v laufen gesehen, PASS.
4) jaira-role-lane/SKILL.md:72-95 — fertige 'git add'/'git commit'-Zeile mit Handle im Betreff und der Ticket-Datei im add; die Gegenseite (Lane ohne Code-Aenderung gibt keine Zeile) steht an beiden Stellen.
5) jaira-role-lane/SKILL.md:44-60 — 'git diff' nach jedem DoD-Punkt, leerer Diff heisst keine Pause.
6) core/release/NOTES.md:17 und :18 unter ## Unreleased, je eine Zeile, kein Umbruch.

Funktion, auf einem frischen Scratch-Board mit dem selbst gebauten Binary durchgespielt:
- 'jaira set <id> mode=chat' -> 'jaira: mode is "conversational" or empty, got "chat"', RC=2.
- 'jaira set <id> mode=" conversational "' -> akzeptiert, 'show --json' liefert 'conversational' (getrimmt gespeichert).
- 'jaira show --for-lane in-progress --json' fuehrt mode neben model_tier; der Klartext-Kopf lautet '# Lane: Implementing   (tier: cheap, mode: conversational)'.
- 'jaira resume --json' fuehrt "mode": "conversational" im in_flight-Eintrag, 'jaira resume' im Klartext die Zeile 'mode: conversational' — das ist der Befund aus critique-Runde 7, im laufenden Binary bestaetigt.
- Mode von Hand im Ticket-File auf 'chat' gesetzt: 'jaira validate' meldet die Warnung samt Reparaturzeile.

Nichts gefunden, was zurueckgeht. Diese Lane hat keinen Code geaendert und committet deshalb nichts.
