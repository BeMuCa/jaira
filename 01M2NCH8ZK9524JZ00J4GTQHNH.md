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
updated-at: 2026-09-17T18:51:48Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-75151
claimed-at: 2026-09-17T18:50:16Z
outcome-what: "Die drei Befunde aus critique-Runde 10 repariert. core/role/builtin/jaira-role-lane/SKILL.md: der Lesebefehl 'jaira show --for-lane --json' steht jetzt VOR 'jaira claim' — claim schreibt den eigenen Namen aufs Ticket, und im Gespraechsmodus darf einer der beiden Worker gar nichts schreiben; 'jaira claim' ist dafuer ein eigener Punkt der Schritt-Liste geworden. Die Verbotsliste des mitlaufenden Kritikers nennt claim und dod und endet mit dem Satz, dass alles aus ist, was diese Datei sonst zu schreiben auftraegt. Der Unterscheider 'lane-Argument != status' wird ausdruecklich EINMAL gelesen und nie neu, plus die Regel fuer den Worker, dessen Kontext weg ist: nicht raten, den Dispatcher fragen, und bis zur Antwort nichts schreiben — und melden und aufhoeren statt weiterzulaufen. core/role/builtin/jaira-dispatcher/SKILL.md Schritt 2 fuehrt dieselbe Einmal-Lesung und verpflichtet den Dispatcher, die Frage zu beantworten. scripts/spawn.sh: der usage-Text zu --no-worktree nennt den lesenden Worker als Ausnahme. core/release/NOTES.md: die Zeile zur mitlaufenden Kritik behauptete genau das, was Befund 1 widerlegt, und die Zeile zu --no-worktree dasselbe Falsche wie spawn.sh — beide korrigiert."
outcome-why: "Der Unterscheider war eine Momentaufnahme des Boards, kein bleibendes Merkmal: der implementierende Worker beendet seine Lane mit 'jaira move --to critique', ab da ist status == 'critique' == dem Lane-Argument der mitlaufenden Kritik, und die haelt sich fuer die ordentliche critique-Lane — genau die zwei Schreiber auf einem Feld, die der Absatz verhindern soll. Und 'jaira claim' stand vor dem Lesebefehl: die Kritik hatte geschrieben, bevor sie wissen konnte, dass sie es nicht darf."
outcome-resolves: "DoD 9 (genau eine Stelle schreibt review-summary und bewegt das Ticket): jaira-role-lane/SKILL.md:69-100 traegt die Regel jetzt so, dass sie einen Statuswechsel und einen Sitzungsabbruch ueberlebt — Einmal-Lesung plus Rueckfrage statt Rateweg — und der Lesebefehl steht vor jedem Schreibweg, claim eingeschlossen. Plan 26-29 [x]; go build ./... und go test ./core/... ./internal/... -count=1 gruen."
review-summary: |-
  core/role/builtin/jaira-role-lane/SKILL.md:62-67 — der Unterscheider kippt mitten im Lauf. Die mitlaufende Kritik startet, waehrend der status 'in-progress' ist; sobald der implementierende Worker 'jaira move --to critique' macht, ist status=critique und das Lane-Argument 'critique' gleich dem status — dieselbe, noch laufende (oder nach Kompaktierung neu lesende) Kritik haelt sich ab diesem Moment fuer die ordentliche critique-Lane und schreibt review-summary und move. Das sind genau die zwei Schreiber auf einem Feld, die der Absatz verhindern soll. Regel ergaenzen: wer als mitlaufende Kritik gestartet ist, bleibt es auch, wenn der status auf sein Lane-Argument wechselt — dann meldet er und hoert auf, und der Dispatcher startet den Worker der Lane.
  core/role/builtin/jaira-role-lane/SKILL.md:13 und :22 — 'jaira claim' und 'jaira dod' fehlen in der Verbotsliste von Punkt 1 ('no jaira note, no jaira set, no jaira move, no review-summary, no commit line'), obwohl derselbe Prompt beide ausdruecklich anordnet. Schlimmer: 'jaira claim' steht in Zeile 13, VOR dem 'show --for-lane --json', aus dem der Modus und der status ueberhaupt erst gelesen werden — die mitlaufende Kritik hat also schon auf das Ticket geschrieben, bevor sie wissen kann, dass sie nichts schreiben darf. Den Lesebefehl vor 'jaira claim' ziehen (oder claim an die Bedingung haengen) und beide Kommandos in die Verbotsliste aufnehmen.
  core/role/builtin/jaira-dispatcher/scripts/spawn.sh:14 — der usage-Text zu --no-worktree sagt weiterhin 'Two workers then share one directory, so run only one at a time'. jaira-dispatcher/SKILL.md:215-220 nennt die mitlaufende Kritik jetzt als die eine Ausnahme; ein Dispatcher, der 'spawn.sh --help' liest, bekommt die alte Regel. Denselben Halbsatz dort nachziehen: lesende Worker sind die Ausnahme.
  core/role/builtin/jaira-role-lane/SKILL.md:9-18 — the head names 'jaira show <ticket-id> --for-lane <lane> --json' as 'the call that tells you which of them you are'. That JSON has no status key (complete, diff, input, lane, missing, model_tier, produces, prompt, ticket_id), so it cannot tell anyone; the deciding read is the separate 'jaira show <id> --json' in section 1, three screens further down, after the 'jaira claim' bullet at :26. A worker reading this file in order claims the ticket before it reaches the read that decides — the exact write cf7d37c moved behind the read. Name the status read in the head as well, or end the head paragraph with: mode=conversational means go to section 1 and decide there BEFORE you claim.
  core/role/builtin/jaira-dispatcher/SKILL.md:113-115 — the read-only list there was extended in cf7d37c with 'does not claim the ticket' but still omits 'jaira dod', which jaira-role-lane/SKILL.md:74-77 names and which writes the very field the implementing worker is ticking. Close it the way role-lane closes its list: end the enumeration with a sentence ('and nothing else its prompt tells a worker to write') instead of a list that the next new write command falls out of.
review-gaps: |-
  Drei Befunde, der erste ist ein echter Defekt.

  1. 'git diff' ist blind fuer neue Dateien, und genau daran haengt die Pause. jaira-role-lane/SKILL.md sagt nach jedem DoD-Punkt 'git diff' und dann woertlich: 'Empty output? No pause'. Eine untracked Datei steht in 'git diff' nicht drin. Ein DoD-Punkt, der aus einer neuen Datei besteht, erzeugt also leere Ausgabe, der Worker haelt nicht an und baut weiter — die eine Sache, die der Modus verhindern soll. Der Beleg ist dieses Ticket selbst: seine Implementierung hat core/validate/mode_test.go und internal/cli/mode_test.go neu angelegt, beide waeren im Modus stumm durchgelaufen. Fix ist eine Zeile im Prompt: 'git add -A -N .' vor dem 'git diff', oder 'git status --short' danebenstellen und bei untracked Dateien ebenfalls anhalten.

  2. Der Block, den 'jaira update' in CLAUDE.md/AGENTS.md eines Boards schreibt, kennt den Modus nicht. core/board/announce.go:64 beschreibt die Nutzlast von 'show <id> --for-lane --json' namentlich ('the lane's prompt, the bounded input, the model tier, and the outputs the lane expects back') und nennt die Modellstufe — den Modus nennt es nicht. Dokumentiert ist er in docs/AGENTS.md und README.md, also im Repository von jaira, nicht auf dem Board eines Benutzers. Ein Agent, der die ausgelieferten Rollen nicht benutzt — und CLAUDE.md fordert ausdruecklich, dass die CLI von jedem bash-faehigen Agenten benutzbar ist —, sieht einen 'mode'-Schluessel, den ihm nichts erklaert.

  3. DoD-Punkt 1 ist ausschliesslich als Prosa erfuellt. 'Der Eintritt haengt an einer nachpruefbaren Bedingung' ist ein Prompt-Absatz in jaira-dispatcher/SKILL.md; nichts im Code zaehlt etwas, nichts prueft, ob der Dispatcher den Halt ausgelassen hat, und kein Test deckt ihn ab. Das ist die in der Brainstorm-Lane getroffene Entscheidung (der Modus ist eine Prompt-Aenderung, der Go-Code ist nur der Traeger) und insofern kein Widerspruch — aber es heisst, dass der zentrale DoD-Punkt dieses Tickets durch keinen Mechanismus gehalten wird. Kleiner Nachbrenner in derselben Ecke: der zweite Schreibpfad, internal/tui/edit.go:60, hat keinen Test; die Ablehnung eines unbekannten Modus ist nur fuer 'jaira set' abgedeckt.
test-verdict: "pass: go build/vet/test ./... -race -count=1 green (RC=0, 29 Pakete), DoD 1-6 in der Working Tree geprueft, Verhalten auf einem Scratch-Board mit dem gebauten Binary durchgespielt"
question: "Der Gespraechsmodus ist gebaut und getestet, aber noch nie an einem echten Ticket gelaufen — der Beleg liegt bisher nur in Tests und einem Scratch-Board. Willst du ihn einmal selbst fahren ('jaira set <id> mode=conversational' auf einem Ticket mit offener Form, dann Dispatcher starten), bevor das hier weitergeht, oder reicht dir der Testbericht und es geht direkt in review? Zweitens: die Rollen-Prompts liegen im Repository, deine Kopien in ~/.claude sind noch die alten — 'jaira roles install --global --force' muesste laufen, damit du den Modus ueberhaupt siehst."
review-verdict: |-
  Das Traegerwerk ist sauber und deckt sich mit dem Bericht des Implementierers: ein Feld, ein geschlossener Wertebereich, eine Funktion (CanonicalMode) hinter beiden Schreibpfaden, fuenf Lesestellen inklusive 'jaira resume' — womit DoD-Punkt 3 (ueberlebt den Sitzungsabbruch) wirklich getragen ist und nicht nur behauptet —, dazu validate als Netz fuer alles, was an den Schreibpfaden vorbeikommt. Build, vet und 'go test ./... -count=1' sind gruen; ich habe sie selbst laufen lassen. Alle sechs DoD-Punkte sind formal erfuellt.

  Eine Einschraenkung, die ich nicht zur Freigabe aufrunde: der Defekt aus review-gaps 1 trifft nicht die Randbedingung, sondern den Zweck. Ein DoD-Punkt, der aus einer neuen Datei besteht, laeuft im Gespraechsmodus ohne Pause durch, weil 'git diff' untracked Dateien nicht zeigt — der Fall ist haeufig (Tests, neue Pakete) und der Fix ist eine Zeile im Prompt. Ich empfehle, das vor dem Signoff zu aendern; es geht in eine Person-Minute und nicht in eine Runde durch in-progress. Die anderen beiden Befunde sind Dokumentationsschulden und koennen als eigenes Ticket nachlaufen.
review-check: "Alles unten ist auf diesem Branch von Hand nachgelaufen. Dauer etwa 5 Minuten.\n\n1. Ins Worktree wechseln und bauen:\n   cd /home/alex/projects/.worktrees/jaira-GTQHNH && go build -o /tmp/jp ./cmd/jaira\n   Erwartet: keine Ausgabe, Datei /tmp/jp existiert.\n\n2. Tests laufen lassen:\n   go test ./core/ticket/... ./core/validate/... ./internal/cli/... -count=1\n   Erwartet: drei Zeilen 'ok'.\n\n3. Ein Wegwerf-Board anlegen:\n   rm -rf /tmp/sb && mkdir /tmp/sb && cd /tmp/sb && git init -q . && /tmp/jp init\n   Dann ein Ticket: /tmp/jp create Probe --goal g --context c --dod d --assignee me\n   Das Handle merken (sechs Zeichen, steht in der Ausgabe).\n\n4. Einen falschen Modus setzen: /tmp/jp set <handle> mode=chat\n   Erwartet woertlich: 'jaira: mode is \"conversational\" or empty, got \"chat\"', Exit-Code 2. Nichts wird gespeichert.\n\n5. Den richtigen setzen und ihn an den drei Stellen wiederfinden:\n   /tmp/jp set <handle> mode=conversational\n   /tmp/jp show <handle>            -> eine Zeile 'mode       conversational' direkt unter 'tier'\n   /tmp/jp show <handle> --for-lane in-progress | head -1\n                                    -> '# Lane: Implementing   (tier: cheap, mode: conversational)'\n   /tmp/jp show <handle> --for-lane in-progress --json | grep mode\n                                    -> '\"mode\": \"conversational\",' neben \"model_tier\"\n\n6. Den Modus von Hand kaputt machen und validate fragen:\n   sed -i 's/^mode: conversational/mode: chat/' /tmp/sb/.jaira/tickets/*.md && /tmp/jp validate\n   Erwartet: 'warning <handle>  mode: mode \"chat\" is not a mode: ... set it with jaira set <handle> mode=conversational'.\n\n7. Den Defekt aus review-gaps 1 selbst sehen — der eine Punkt, bei dem ich zum Zurueckgeben rate:\n   cd /home/alex/projects/.worktrees/jaira-GTQHNH && touch /tmp/neu.go && cp /tmp/neu.go ./neu.go && git diff\n   Erwartet: LEERE Ausgabe, obwohl eine neue Datei da liegt. Genau diese leere Ausgabe sagt dem Worker laut core/role/builtin/jaira-role-lane/SKILL.md:55 'keine Pause'. Danach aufraeumen: rm ./neu.go\n\nNicht von Hand pruefbar ist der Halt vor der Plan-Lane (DoD 1): er ist ein Prompt-Absatz in core/role/builtin/jaira-dispatcher/SKILL.md:31-60, kein Code. Lesen ist die einzige Pruefung, die es dafuer gibt."
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
  proof: core/role/builtin/jaira-role-lane/SKILL.md:107-128 — statt 'git commit' die fertige Zeile mit Handle im Betreff und der Ticket-Datei im 'git add'; die Gegenseite in core/role/builtin/jaira-dispatcher/SKILL.md:84 — Zeile nur bei Code-Aenderung
- [x] Der Mensch sieht den Code, bevor darauf aufgebaut wird: der Dispatcher legt ihn nach jedem Inkrement vor und wartet, statt am Ende alles auf einmal zu zeigen.
  proof: core/role/builtin/jaira-role-lane/SKILL.md:74-105 — 'git status --short' und 'git diff' nach jedem DoD-Punkt vorlegen und warten; beide leer heisst keine Pause
- [x] Je eine Zeile in core/release/NOTES.md unter ## Unreleased fuer den Modus und fuer das, was ein Benutzer beim Committen anders tut.
  proof: core/release/NOTES.md:17 (Modus, inkl. der Kopfzeile von 'jaira show --for-lane') und :18 (Committen von Hand), beide unter ## Unreleased
- [x] Im Gespraechsmodus laeuft eine Kritik mit, waehrend an dem Ticket gearbeitet wird, und meldet ihre Befunde sofort - nicht erst, nachdem die implementierende Lane fertig ist. Der Modus und die mitlaufende Kritik stehen im selben Prompt und wirken zusammen.
  proof: core/role/builtin/jaira-dispatcher/SKILL.md:100-136 — Abschnitt 'In conversational mode, a critique runs beside the work', im selben Prompt wie der Modus und aus der Modus-Liste heraus verlinkt (SKILL.md:96); Schritt 3 verlangt den Befund in dem Moment, in dem er da ist
- [x] Eine Kritik sieht, was frueher schon geprueft und was schon repariert wurde, und hebt einen erledigten Befund nicht erneut auf. Nachgestellt an einem Ticket, dessen zweite Runde etwas ausdruecklich fuer gut befunden hat: die dritte liest das und prueft es nicht noch einmal.
  proof: .jaira/lanes/critique.md:10 fuehrt 'notes' in input-requires und der Prompt-Absatz darunter ('Read the notes before you read the diff') macht einen in einer Notiz erledigten Befund zu einem geschlossenen; nachgestellt an diesem Ticket: 'jaira show GTQHNH --for-lane critique --json' liefert missing=null und den notes-Schluessel mit den acht critique-Runden darin
- [x] Genau eine Stelle schreibt review-summary und bewegt das Ticket. Laufen mehrere Kritiker gleichzeitig, lesen sie nur; das Zusammenfuehren und der eine 'jaira move' gehoeren dem Dispatcher.
  proof: core/role/builtin/jaira-role-lane/SKILL.md:69-100 — die Nur-Lese-Rolle wird einmal entschieden und nie neu; wer es nicht mehr weiss, fragt den Dispatcher
- [x] Die Pause nach einem DoD-Punkt erkennt auch eine NEU angelegte Datei. Nachgestellt an einem Punkt, der nur aus einer neuen Datei besteht: der Modus haelt an, statt ihn als leere Ausgabe durchlaufen zu lassen.
  proof: core/role/builtin/jaira-role-lane/SKILL.md:50-80 — 'git status --short' steht jetzt neben 'git diff', ein '??' ist eine Pause und die neue Datei wird mit vorgelegt; nachgestellt in diesem Worktree: eine angelegte probe_neu.go liefert 'git diff' 0 Zeilen und 'git status --short' die Zeile '?? probe_neu.go'
- [x] Der Halbsatz 'Testing is not a lane' steht nicht mehr in core/role/builtin/jaira-dispatcher/SKILL.md; scripts/spawn.sh hat weiterhin genau einen Sonderfall ('dispatch'), und core/role/builtin/jaira-role-tester/SKILL.md ist unveraendert.
  proof: grep 'Testing is not a lane' core/ findet nichts mehr; core/role/builtin/jaira-dispatcher/SKILL.md:157 lautet jetzt '/jaira-role-lane <id> <lane> — every lane, testing included'; spawn.sh und jaira-role-tester/SKILL.md stehen unveraendert in 'git status --short'
- [x] Je eine Zeile in core/release/NOTES.md unter ## Unreleased fuer die mitlaufende Kritik und fuer die Pause, die neue Dateien sieht.
  proof: core/release/NOTES.md:16 (mitlaufende Kritik), :17 (notes als Lane-Eingabe) und :18 (die Pause, die neue Dateien sieht), alle drei unter ## Unreleased

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
- [x] .jaira/lanes/critique.md: 'notes' in input-requires plus Prompt-Absatz — eine Kritik liest, was frueher schon geprueft und repariert wurde
- [x] jaira-dispatcher/SKILL.md: neuer Abschnitt 'mitlaufende Kritik' — im Gespraechsmodus ein zweiter, nur lesender Kritiker, Befunde sofort, genau ein Schreiber von review-summary und genau ein 'jaira move'
- [x] jaira-role-lane/SKILL.md: die Pause nach einem DoD-Punkt sieht auch eine neu angelegte Datei
- [x] jaira-dispatcher/SKILL.md: Halbsatz 'Testing is not a lane' streichen; spawn.sh und jaira-role-tester unveraendert
- [x] core/release/NOTES.md: je eine Zeile unter ## Unreleased fuer die mitlaufende Kritik und fuer die Pause, die neue Dateien sieht
- [x] go test ./... -count=1
- [x] critique 9: jaira-role-lane/SKILL.md — der mitlaufende Kritiker erkennt sich am Ticket (lane-Argument != status), nicht an der Startzeile, und schreibt nichts
- [x] critique 9: core/release/NOTES.md — die falsche Behauptung 'no lane shipped with jaira asked for it' korrigieren
- [x] critique 10: jaira-role-lane/SKILL.md — der Unterscheider wird EINMAL beim ersten Lesen entschieden und nicht neu, wenn der status auf das Lane-Argument wechselt; wer es nicht mehr weiss, schreibt nicht
  proof: core/role/builtin/jaira-role-lane/SKILL.md:69-100
- [x] critique 10: jaira-role-lane/SKILL.md — 'show --for-lane --json' VOR 'jaira claim', und claim/dod in die Verbotsliste der mitlaufenden Kritik
  proof: core/role/builtin/jaira-role-lane/SKILL.md:10-29,73-76
- [x] critique 10: spawn.sh:14 usage von --no-worktree — lesende Worker als Ausnahme nachziehen
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:11-16
- [x] critique 10: core/release/NOTES.md — die Zeile zur mitlaufenden Kritik beschreibt den neuen Unterscheider
  proof: core/release/NOTES.md:17,29
- [x] critique 11: jaira-role-lane/SKILL.md Kopf — 'show --for-lane --json' fuehrt kein status und kann den Unterscheider nicht liefern; der entscheidende Lesebefehl 'jaira show <id> --json' steht im Kopf, vor 'jaira claim'
- [ ] critique 11: jaira-dispatcher/SKILL.md — die Verbotsliste der mitlaufenden Kritik nennt 'jaira dod' und schliesst mit einem Satz statt einer Aufzaehlung

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
- **2026-09-16 20:27 · Alexander Sacharov** — Befund vom 16.09. beim Ausliefern der Rollen — er trifft die NOTES.md-Zeilen dieses Tickets direkt.

'jaira roles install --global --force' verteilt NICHT die Rollen aus dem Arbeitsbaum, sondern die, die im installierten Binary einkompiliert sind (/home/alex/.local/bin/jaira, Version 0.2.1). Das Binary ist aelter als dieser Branch. Der Aufruf hat die globalen Kopien deshalb ZURUECKGEDREHT: ~/.claude/skills/jaira-dispatcher/SKILL.md hatte danach null Treffer auf 'conversational' gegen sechs im Repository, und jaira-teamlead verlor den Absatz ueber 'spawn.sh <slug> <id> dispatch', den das Repository fuehrt. Von Hand wieder gleichgezogen, alle sieben Rollen stimmen jetzt mit core/role/builtin ueberein.

Warum das dieses Ticket angeht: beide NOTES.md-Zeilen schicken den Leser genau auf diesen Befehl, damit er den Modus sieht. Vor einem Release, das den neuen Prompt einbettet, tut der Befehl das Gegenteil — er nimmt ihm den Modus weg, schweigend. Entweder die Zeilen sagen dazu, dass erst das Release den Modus bringt, oder 'roles install' lernt, aus dem Arbeitsbaum zu installieren. Das zu entscheiden ist nicht Sache dieses Tickets, aber die NOTES.md-Zeilen so stehen zu lassen, ist falsch.

Zweiter, kleinerer Befund aus demselben Lauf: scripts/spawn.sh kennt nur 'dispatch' als Sonderfall und tippt fuer jede andere Lane '/jaira-role-lane <id> <lane>'. Fuer die testing-Lane ist das laut Dispatcher-Prompt falsch — dort gehoert '/jaira-role-tester <id>' hin. Der Lauf ging gut aus, weil der Lane-Worker denselben Prompt aus .jaira/lanes/testing.md liest und test-verdict liefert, aber der Prompt und das Skript widersprechen sich.
- **2026-09-16 21:05 · Alexander Sacharov** — Antwort von Alex am 16.09. auf die Frage aus der human-Lane.

Kein Probelauf an einem echten Ticket vor review. Begruendung von Alex: wie viel geredet wird, haengt ohnehin am einzelnen Ticket; es kommen noch viele Aenderungen, und ein Review findet sowieso statt. Das Ticket geht damit direkt in die review-Lane.

Dazu eine Praezisierung des Modus, die Alex im selben Satz macht und die auf das Ticket gehoert: der Gespraechsmodus ist kein Modus, in dem staendig gefragt wird. Alex redet viel, WAEHREND entschieden wird — was aber einmal entschieden ist, wird danach autonom ausgefuehrt und nicht noch einmal aufgemacht. Der Halt vor der Plan-Lane zaehlt deshalb OFFENE Entscheidungen; eine, die in einer Notiz schon beantwortet ist, ist keine offene mehr. Genau das steht seit critique-Runde 2 im Zaehlschritt des Dispatcher-Prompts — diese Antwort bestaetigt es als gewollt und nicht als Zufall.

Zwei Punkte aus demselben Austausch, nachgeprueft statt vermutet:

Zur Auslieferung der Rollen: core/role/role.go:33 traegt '//go:embed all:builtin'. Alles unter core/role/builtin liegt also beim Bauen automatisch im Binary, es geht nichts verloren. Das Binary 0.2.1 ist schlicht aelter als dieser Branch — 'roles install' war nicht falsch, nur die Quelle war alt. Mit dem naechsten Release loest sich das von selbst; die NOTES.md-Zeilen sollten trotzdem sagen, dass erst dieses Release den Modus bringt.

Zum Widerspruch um die testing-Lane: entschieden wird auf EINE Regel, und zwar die, die heute gelaufen ist. testing ist eine gewoehnliche Lane — .jaira/lanes/testing.md ist 'agentic: true' mit 'output-produces: [test-verdict]' —, sie wird mit '/jaira-role-lane <id> testing' gefahren, und spawn.sh behaelt seinen einen Sonderfall 'dispatch'. jaira-role-tester bleibt, was seine eigene Beschreibung sagt: eine Einstiegsstelle mit zwei optionalen Argumenten fuer 'pruef die Suite', kein Lane-Worker. Zu streichen ist damit der Halbsatz in core/role/builtin/jaira-dispatcher/SKILL.md:118 ('Testing is not a lane'), nicht ein zweiter Zweig in spawn.sh.
- **2026-09-16 21:10 · Alexander Sacharov** — Review-Lane. Vier Felder gesetzt, nach eigenem Lauf von Build, vet, 'go test ./... -count=1' (gruen) und einem Scratch-Board mit dem gebauten Binary: set lehnt 'chat' mit Exit 2 ab, 'jaira show' druckt die mode-Zeile unter tier, die Klartext-Kopfzeile von '--for-lane' lautet '(tier: cheap, mode: conversational)', validate meldet einen handgeschriebenen Modus als Warnung mit der Reparaturzeile.

Der Befund, der einen Menschen braucht: jaira-role-lane/SKILL.md sagt nach jedem DoD-Punkt 'git diff' und 'Empty output? No pause'. 'git diff' zeigt untracked Dateien nicht. Ein DoD-Punkt, der aus einer NEUEN Datei besteht, laeuft im Gespraechsmodus also ohne Pause durch — der Fall ist haeufig (neue Tests, neue Pakete), und die Implementierung dieses Tickets ist selbst ein Beispiel: core/validate/mode_test.go und internal/cli/mode_test.go sind neu. Ein-Zeilen-Fix im Prompt: 'git add -A -N .' vor dem 'git diff', oder 'git status --short' danebenstellen.

Nicht zurueck nach in-progress geschickt, weil die Aenderung selbst korrekt ist und der Fix eine Prompt-Zeile ist — das ist billiger am Signoff als eine Runde durch die Lanes. Die zwei restlichen Befunde (Modus fehlt im von 'jaira update' erzeugten CLAUDE.md-Block in core/board/announce.go:64; internal/tui/edit.go hat keinen Test fuer die Ablehnung) sind Nachlauf und koennen ein eigenes Ticket werden.
- **2026-09-16 21:12 · Alexander Sacharov** — Erweiterung des Tickets, von Alex am 16.09. angeordnet: die mitlaufende Kritik gehoert in DIESES Ticket und in denselben Prompt, nicht in ein eigenes. Beides soll zusammen wirken - der Gespraechsmodus haelt an und fragt, und waehrend gearbeitet wird, meldet eine Kritik nebenher, ob irgendwo schon ein Problem liegt. Zwei Tickets, die ich dafuer angelegt hatte, sind wieder geloescht.

WARUM das noetig ist, gemessen an diesem Ticket selbst, nicht vermutet:

Die critique-Lane laeuft heute erst, wenn in-progress fertig ist. GTQHNH brauchte so acht Durchgaenge: 3 Befunde, dann 2, dann fuenfmal je 1, dann 0.

Jeder dieser Durchgaenge fand in einer DATEI etwas, die die vorigen nie geoeffnet hatten - Runde 3 in core/role/builtin/jaira-role-lane/SKILL.md, Runde 4 an vier Doku-Stellen, Runde 5 in core/validate/validate.go, Runde 7 in internal/cli/resume.go. Der Diff wurde ihr dabei jedes Mal VOLLSTAENDIG gereicht: internal/cli/flow.go:581 baut ihn aus allen Commits des Tickets, nicht aus dem letzten. Es ist also kein Fenster-Problem. Ein einzelner Leser nimmt einen Ausschnitt und hoert auf, wenn es reicht.

Der teuerste Befund kam in Runde 7: 'jaira resume' trug das neue Feld nicht, womit der Wiederanlauf nach einem Sitzungsabbruch - der Daseinsgrund dieses Tickets - nur auf dem Papier funktioniert haette. Bei einer Abbruchregel von drei Runden waere er nie gefunden worden.

Zweiter Grund, eine fehlende Eingabe: .jaira/lanes/critique.md:10 fuehrt 'input-requires: [goal, definition-of-done, outcome-what, outcome-resolves, diff]'. 'notes' steht dort nicht - von allen Lanes bekommt nur in-progress sie. Der Prompt der Lane verlangt aber ausdruecklich, einen schon appreparierten Befund nicht erneut aufzumachen. Dafuer fehlt ihr die Eingabe: Runde 2 hat freiwillig notiert, was sie geprueft und fuer gut befunden hat, Runde 3 hat das nicht gelesen und dieselbe Gegend noch einmal gelesen.

Der Haken, der mitentschieden werden muss, wenn mehrere Kritiker gleichzeitig laufen: die Lane schreibt heute selbst review-summary UND bewegt selbst das Ticket. Mehrere Worker wuerden sich um beides pruegeln. Dann duerfen die Kritiker nur lesen, und das Zusammenfuehren samt dem einen 'jaira move' gehoert dem Dispatcher.

Zwei weitere Punkte, die in denselben Durchgang gehoeren:

Aus dem review-Lane-Befund vom 16.09.: core/role/builtin/jaira-role-lane/SKILL.md:55 laesst nach jedem DoD-Punkt 'git diff' laufen und sagt 'Empty output? No pause'. 'git diff' zeigt keine untracked Dateien. Ein DoD-Punkt, der aus einer NEUEN Datei besteht, laeuft im Gespraechsmodus also ohne Pause durch - genau das, was der Modus verhindern soll. Dieses Ticket ist selbst das Beispiel: core/validate/mode_test.go und internal/cli/mode_test.go sind neu angelegt. Fix ist eine Zeile: 'git add -A -N .' vor dem 'git diff', oder 'git status --short' danebenstellen.

Aus dem Lauf selbst: core/role/builtin/jaira-dispatcher/SKILL.md:118 sagt 'Testing is not a lane: /jaira-role-tester <id>', waehrend scripts/spawn.sh:139-145 nur 'dispatch' als Sonderfall kennt und testing zu '/jaira-role-lane <ticket> testing' macht. Gelaufen ist es als Lane und es ging gut aus. Alex hat auf die gelaufene Regel entschieden: testing ist eine gewoehnliche Lane, der Halbsatz in Zeile 118 wird gestrichen, spawn.sh behaelt seinen EINEN Sonderfall, und jaira-role-tester bleibt unangetastet als Einstiegsstelle mit zwei optionalen Argumenten.
- **2026-09-17 18:13 · Alexander Sacharov** — Alex am 2026-09-17: die Punkte 7-12 werden jetzt gemacht, das Ticket geht aus signoff zurueck in die Arbeit. Anlass: 'jaira move GTQHNH --to done' wurde abgewiesen, weil 7-11 offen sind - der Dispatcher hatte den Zweig schon in release/0.3.0 gezogen, nur nach der Lane-Marke signoff und ohne die Definition of Done zu lesen. Punkte 1-6 (der Modus selbst) sind fertig und ihre drei NOTES-Zeilen stehen bereits unter '## 0.3.0' in der Release-Ветке. Was jetzt dazukommt, gehoert in denselben Release.
- **2026-09-17 18:19 · Alexander Sacharov** — In-progress, Punkte 7-12: vier Befunde, die das Repository nicht selbst sagt.

'git add -A -N .' vor dem 'git diff' — der Fix, den sowohl die review-Lane als auch Alex' Notiz vom 16.09. vorgeschlagen haben — ist NICHT genommen worden. Grund: CLAUDE.md und jaira-role-lane/SKILL.md:27 verbieten ausdruecklich 'git add -A', weil eine andere Sitzung denselben Worktree halten kann. Der Gespraechsmodus startet den Worker per '--no-worktree' im ausgecheckten Verzeichnis — er ist damit genau der Fall, den das Verbot meint, nicht die Ausnahme davon. 'git status --short' leistet dasselbe und schreibt nichts in den Index. Wer das spaeter 'vereinfachen' will: das ist der Grund.

.jaira/lanes/critique.md ist die Lane-Datei DIESES Boards, keine ausgelieferte. core/lane/builtin/ fuehrt zehn Lanes und critique ist keine davon. Die 'notes'-Eingabe erreicht fremde Boards deshalb nur ueber die NOTES.md-Zeile, die dem Leser sagt, er solle sie selbst eintragen — nicht ueber ein Release. Eine ausgelieferte critique-Lane zu erfinden waere ein eigenes Ticket, nicht dieses.

Dafuer war kein Go-Code noetig: core/ticket/schema.go:193-197 fuehrt 'notes' bereits in SuppliedFields, und internal/cli/notesinput_test.go deckt den Weg ab. Nachgeprueft statt vermutet — 'show GTQHNH --for-lane critique --json' liefert missing=null und den notes-Schluessel.

Die mitlaufende Kritik widersprach dem '--no-worktree'-Absatz, der sagt 'never run a second one anywhere while such a worker is live'. Aufgeloest, indem der Absatz sie als die eine Ausnahme benennt: sie liest nur, und deshalb kann nichts an ihr mit dem schreibenden Worker kollidieren. Haette man den Absatz stehen lassen, widerspraeche der Prompt sich selbst an zwei Stellen und der Dispatcher folgt der, die er zuerst liest.
- **2026-09-17 18:22 · Alexander Sacharov** — critique (9. Durchgang, erste Runde ueber die Punkte 7-12): zwei Befunde, beide mit klarem Fix, keine Entscheidung fuer den Menschen.

1. Der wichtigere. dispatcher/SKILL.md:115-117 traegt die Nur-Lese-Eigenschaft der mitlaufenden Kritik in der GETIPPTEN STARTZEILE ('Say that in the line you start it with, because its lane prompt tells it to do all three'). Das ist Weg B aus der Brainstorm-Entscheidung dieses Tickets, dort verworfen mit der Begruendung, die hier woertlich wieder zutrifft: er stirbt mit der Sitzung. Was ein neu gestarteter oder kompaktierter Worker liest, ist jaira-role-lane/SKILL.md, und die sagt ihm weiterhin note, review-summary und move — plus 'Boundaries: the lane you were given is the deliverable'. Dann ueberschreibt der Kritiker review-summary und zieht das Ticket unter dem implementierenden Worker weg. Fix an derselben Stelle, an der der Modus selbst schon sitzt: ein dritter Punkt in der Liste 'Two things change' (jaira-role-lane/SKILL.md:48), Unterscheidung ueber 'jaira show <id> --json' — Lane-Argument != status heisst mitlaufende Kritik, also nur melden. Die ordentliche critique-Lane laeuft mit status=critique und schreibt unveraendert weiter.

2. NOTES.md:18 behauptet 'no lane shipped with jaira asked for it'. core/lane/builtin/20-in-progress.md:9 fuehrt notes in input-requires und wird ausgeliefert. Halbsatz streichen oder richtigstellen.

Nicht erneut aufgemacht, mit Angabe der Notiz, auf der ich stehe: 'git add -A -N' bleibt verworfen (Notiz 17.09. 18:19, Grund --no-worktree); die ausgelieferte critique-Lane ist bewusst ein eigenes Ticket (dieselbe Notiz); der --no-worktree-Absatz nennt die Kritik schon als die eine Ausnahme (dieselbe Notiz); die ueber fuenf Dateien wiederholte Prosa steht (optimize-Notiz 16.09. 20:21); 'Testing is not a lane' ist von Alex angeordnet (DoD 11).
- **2026-09-17 18:27 · Alexander Sacharov** — In-progress nach critique (9. Durchgang): beide Befunde repariert, einer davon an anderer Stelle in der Liste als critique vorgeschlagen hat.

Befund 1 — der mitlaufende Kritiker wusste seine Nur-Lese-Regel nur aus der Startzeile. Die Reparatur steht jetzt in core/role/builtin/jaira-role-lane/SKILL.md, also dort, wo ein neu gestarteter oder komprimierter Worker sie wieder liest. Der Unterscheider ist der von critique vorgeschlagene: 'jaira show <id> --json', und ein lane-Argument, das nicht der status des Tickets ist, heisst 'du bist die mitlaufende Kritik'. Er traegt in beide Richtungen — die gewoehnliche critique-Lane laeuft mit status=critique und schreibt weiter.

Anders als vorgeschlagen: critique wollte einen DRITTEN Punkt in der Liste 'Two things change'. Er steht jetzt als ERSTER. Grund: Punkt 'Do not commit' sagt woertlich 'You still jaira move the ticket'. Ein Kritiker, der die Liste der Reihe nach liest, hat dann schon gehandelt, bevor der Widerspruch kommt — und das ist genau die Form von Fehler ('spaet ist die teure Zeit'), die dieser ganze Modus verhindern soll. Die Punkte 'Show the code' und 'Do not commit' sind dadurch zu 2. und 3. geworden; die --proof-Zeilen von DoD 4, 5 und 9 zeigen auf die neuen Zeilennummern.

Zusaetzlich mitgenommen (nicht von critique verlangt): jaira-dispatcher/SKILL.md Schritt 2 sagte 'because its lane prompt tells it to do all three' — das stimmte nach der Reparatur nicht mehr und waere als stehengebliebene Begruendung die naechste Runde wert gewesen. Die Startzeile bleibt trotzdem verlangt, sie kostet einen Nebensatz.

Befund 2 — 'no lane shipped with jaira asked for it' in NOTES.md war falsch: core/lane/builtin/20-in-progress.md:9 fuehrt notes in input-requires und wird ausgeliefert. Die Zeile nennt jetzt den wahren Grund, warum die Aenderung ueberhaupt von Hand gemacht werden muss: critique ist gar keine ausgelieferte Lane.

Gepruefte Alternative und verworfen: critique als Builtin-Lane nach core/lane/builtin/ zu legen, damit die Zeile von selbst wahr wird. Das ist eine andere Aenderung an einem anderen Vertrag (welche Lanes jaira mitbringt) und gehoert nicht in dieses Ticket.
- **2026-09-17 18:31 · Alexander Sacharov** — critique (10. Durchgang, ueber die Punkte 7-12): drei Befunde, alle mit klarem Fix, keine Entscheidung fuer den Menschen. Alle drei entstehen erst durch die Reparatur aus Runde 9 (3eb6ff6) und sind in keiner frueheren Notiz behandelt.
- Befund 1 (der schwere): der Unterscheider 'Lane-Argument != status' in jaira-role-lane/SKILL.md:62-67 ist keine Eigenschaft des Workers, sondern eine Momentaufnahme des Boards. Die mitlaufende Kritik startet bei status=in-progress; der implementierende Worker beendet seine Lane mit 'jaira move --to critique'; ab da liest dieselbe Kritik status=critique == ihr Lane-Argument und haelt sich fuer die ordentliche critique-Lane. Genau der Fall, gegen den der Absatz geschrieben wurde. Der Prompt verlangt das Nachlesen ausdruecklich ('not out of the line that started you'), also trifft es auch jeden Neustart und jede Kompaktierung nach dem move.
- Befund 2: die Verbotsliste in Punkt 1 zaehlt note/set/move/review-summary/Commit-Zeile auf, aber nicht 'jaira claim' (SKILL.md:13) und nicht 'jaira dod' (SKILL.md:22) — beides Schreibwege, beide vom selben Prompt angeordnet, 'jaira dod' auf dasselbe Feld, an dem der Worker daneben arbeitet. Und claim steht VOR dem show, aus dem Modus und status kommen: die Kritik hat geschrieben, bevor sie wissen konnte, dass sie es nicht darf. Reihenfolge umdrehen statt nur die Liste verlaengern.
- Befund 3: spawn.sh:14 (usage von --no-worktree) sagt weiter 'run only one at a time', waehrend jaira-dispatcher/SKILL.md:215-220 die mitlaufende Kritik als Ausnahme fuehrt. Dokumentationsschuld, aber im Skript, das der Dispatcher selbst aufruft.
- Ausdruecklich stehen gelassen, mit Begruendung: die NOTES.md-Zeile zur critique-Lane beschreibt eine Hand-Aenderung an '.jaira/lanes/critique.md' und nicht am Binary — sie gehoert trotzdem dorthin, weil sie sagt, was ein Leser TUN muss, und weil die Zeile das selbst offenlegt ('critique is not a shipped lane at all, so this one reaches your board only by hand'). Das ist der Befund aus Runde 9, dort repariert; ich hebe ihn nicht wieder auf. DoD 11 nachgeprueft: 'Testing is not a lane' ist raus, spawn.sh hat genau einen Sonderfall (Zeile 142, 'dispatch'), jaira-role-tester/SKILL.md unveraendert auf diesem Branch.
- **2026-09-17 18:35 · Alexander Sacharov** — In-progress nach critique (10. Durchgang): die drei Befunde repariert. Zwei davon gehen ueber das hinaus, was critique vorgeschlagen hat.

- Befund 1, der eigentliche Defekt. critique wollte 'wer als mitlaufende Kritik gestartet ist, bleibt es'. Das ist nicht ausfuehrbar: nach einer Kompaktierung weiss der Worker nicht mehr, als was er gestartet ist — genau der Fall, gegen den der Absatz geschrieben ist. Die Regel lautet deshalb jetzt zweiteilig: (a) die Lesung 'lane-Argument != status' wird EINMAL gemacht, vor dem ersten Schreiben, und danach nie wieder — damit kippt sie nicht, wenn der implementierende Worker 'move --to critique' macht; (b) fuer den Fall, dass das eigene Gedaechtnis weg ist, gibt es keinen Rateweg, sondern eine Frage an den Dispatcher, und der Dispatcher-Prompt ist verpflichtet zu antworten (SKILL.md Schritt 2). Ein Default 'im Zweifel nicht schreiben' steht daneben, damit die Frage nicht blockiert.
- Der mitlaufende Kritiker meldet und hoert auf. Das war vorher nicht gesagt und ist die zweite Haelfte von (a): je kuerzer er lebt, desto kleiner das Fenster, in dem sich irgendetwas unter ihm bewegen kann.
- Befund 2, enger als vorgeschlagen: 'jaira claim' wandert nicht nur in die Verbotsliste, sondern der ganze Lesebefehl steht jetzt VOR claim, und der Grund steht dabei. Die Verbotsliste nennt claim und dod ausdruecklich und endet mit einem Satz statt einer Aufzaehlung ('alles, was diese Datei einem Worker zu schreiben auftraegt, ist aus'), damit der naechste neue Schreibbefehl nicht wieder durch die Ritze faellt.
- Befund 3 hat eine zweite Stelle, die critique nicht genannt hat: die NOTES.md-Zeile zu --no-worktree sagte dasselbe Falsche wie spawn.sh:14 ('only one such worker at a time'). Beide nennen jetzt den lesenden Worker als Ausnahme.
- Und die NOTES.md-Zeile zur mitlaufenden Kritik behauptete woertlich das, was Befund 1 widerlegt ('stays read-only instead of writing over the dispatcher'). Korrigiert, nicht ergaenzt — eine Release-Zeile, die eine Eigenschaft verspricht, die der Prompt nicht traegt, ist schlimmer als keine.
- **2026-09-17 18:39 · Alexander Sacharov** — critique (11. Durchgang, ueber die Punkte 7-12): zwei Befunde, beide mit klarem Fix, keine Entscheidung fuer den Menschen. Beide entstehen erst durch cf7d37c und sind in keiner frueheren Notiz behandelt.
- Befund 1 (der schwere): jaira-role-lane/SKILL.md:9-18. Der Kopf sagt woertlich, 'jaira show <id> --for-lane <lane> --json' sei 'the call that tells you which of them you are'. Am laufenden Board nachgesehen: dieses JSON hat die Schluessel complete, diff, input, lane, missing, model_tier, produces, prompt, ticket_id — kein status. Es KANN es nicht sagen. Der entscheidende Lesebefehl ist 'jaira show <id> --json' (traegt status) und steht erst in Abschnitt 1, hinter dem 'jaira claim'-Punkt bei :26. Wer die Datei von oben liest, claimt also, bevor er bei der Lesung ankommt, die entscheidet — genau der Schreibvorgang, den Runde 10 hinter die Lesung verschoben hat. Die Reparatur aus Runde 10 ist dadurch am Kopf wirkungslos.
- Befund 2: jaira-dispatcher/SKILL.md:113-115. Dieselbe Liste wurde in cf7d37c um 'does not claim the ticket' ergaenzt, nennt aber weiter kein 'jaira dod' — das jaira-role-lane/SKILL.md:74-77 auffuehrt und das auf dasselbe Feld schreibt, an dem der implementierende Worker daneben arbeitet. Runde 10 hat die Aufzaehlung in role-lane bewusst mit einem Satz statt einer Liste geschlossen; die Dispatcher-Seite blieb eine Liste und faellt beim naechsten neuen Schreibbefehl wieder auseinander.
- Ausdruecklich stehen gelassen: die Prosa-Wiederholung ueber README, AGENTS.md, beide SKILLs und NOTES.md (optimize-Notiz vom 16.09. 20:21 — fuenf verschiedene Leser, kein Befund). Der Merge-Driver bleibt ungeprueft (Runde 5). Die NOTES.md-Zeile zur critique-Lane beschreibt eine Hand-Aenderung und gehoert dorthin (Runde 9/10). Das Restrennen zwischen dem ersten Lesen der mitlaufenden Kritik und einem sehr fruehen 'move --to critique' des Implementierers hebe ich NICHT auf: Runde 10 hat dafuer den Weg 'im Zweifel nicht schreiben und den Dispatcher fragen' eingezogen, und das deckt den Fall ab.
- **2026-09-17 18:50 · Alexander Sacharov** — Dispatcher, vor der vierten in-progress-Runde: im review-summary von critique-Runde 11 stehen FUENF Befunde, aber nur die letzten zwei sind neu. Die ersten drei sind woertlich aus Runde 10 uebernommen und im Baum bereits geschlossen, nachgesehen nach cf7d37c: die Einmal-Lesung samt Begruendung steht in core/role/builtin/jaira-role-lane/SKILL.md:78-84, 'jaira claim' steht in der Verbotsliste in Zeile 73, und die Kopfzeile 11 verbietet jedes Schreiben vor dem Lesen. Offenbar wurde das Feld angehaengt statt ersetzt. NICHT noch einmal anfassen.

Zu tun sind genau diese zwei:
1. core/role/builtin/jaira-role-lane/SKILL.md:9-18 - der Kopf nennt 'jaira show <ticket-id> --for-lane <lane> --json' als den Aufruf, der sagt, welcher der beiden Worker man ist. Dessen JSON hat gar keinen status-Schluessel (complete, diff, input, lane, missing, model_tier, produces, prompt, ticket_id), kann es also nicht sagen; die entscheidende Lesung ist das separate 'jaira show <id> --json' in Abschnitt 1, drei Bildschirme weiter unten und nach dem 'jaira claim'-Punkt in Zeile 26. Wer die Datei der Reihe nach liest, claimt also weiterhin vor der Entscheidung. Den status-Aufruf im Kopf mitnennen, oder den Kopfabsatz mit dem Satz enden lassen: mode=conversational heisst, vor dem claim in Abschnitt 1 zu gehen und dort zu entscheiden.
2. core/role/builtin/jaira-dispatcher/SKILL.md:113-115 - die Nur-Lese-Liste dort hat in cf7d37c 'claimt das Ticket nicht' bekommen, nennt aber 'jaira dod' nicht, das jaira-role-lane/SKILL.md:74-77 sehr wohl nennt und das genau das Feld schreibt, das der implementierende Worker gerade abhakt. Die Aufzaehlung so schliessen, wie role-lane ihre schliesst: mit einem Satz ('und nichts anderes, was sein Prompt einem Worker zu schreiben auftraegt') statt mit einer Liste, aus der das naechste neue Schreibkommando wieder herausfaellt.

Alex hat diese vierte Runde ausdruecklich freigegeben, nachdem der Dispatcher nach drei Ruecklaeufen angehalten hatte.
