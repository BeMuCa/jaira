---
id: 01M7X20YACKGBHD6P1QKMFYM1A
title: "In signoff hakt die Person die Notes an, die behoben werden sollen, und Give back nimmt genau diese mit"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Die Explain-/Ticket-Ansicht im Cockpit zeigt jede Note aus review-gaps mit einem Haken; die Person waehlt aus, was behoben werden soll, und Give back schreibt genau diese Notes als Auftrag an den Worker; nichts angehakt plus Accept nimmt an."
context: |-
  Heute: Give back im Cockpit bietet nur die Optionen aus der explain:-Note (1) [give back] 2) [accept]); welche der 3-5 Notes behoben werden sollen, muss die Person in Worten sagen oder der Teamlead herausfinden.
  Alex 20.11.2026: "muss man die Moeglichkeit geben, Haken zu setzen, was wir korrigieren wollen?" - nach einer Signoff-Uebersicht als Tabelle mit Marken je Note (Blocker / kleiner Rand / bekannte Grenze / nicht aus diesem Ticket).
  Alex gibt heute fast jedes Ticket zurueck, um offene Raender zu schliessen; die Auswahl ist also der Normalfall, kein Sonderfall.
  Voraussetzung: DF0YD6 (Notes werden eingeordnet und farbig gezeigt).
  Verwandt: Q63EEV (Tickets im Cockpit einzeln abarbeiten: Pane mit Accept/Give back).
definition-of-done: "Jede Note aus review-gaps steht mit Haken in der Ansicht; Blocker und kleine Raender sind vorausgewaehlt, bekannte Grenzen und fremde Notes nicht (Test)"
tags: []
blocked-by:
  - DF0YD6
related:
  - 01M7WWZWSDHK37DGGHJFDF0YD6
  - 01M5M7JTVR8G7DZ1MGHJQ63EEV
commits: []
created-at: 2026-11-20T18:03:14Z
updated-at: 2026-11-21T19:26:41Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2283240
claimed-at: 2026-11-21T19:01:13Z
outcome-what: "The cockpit's Give back pane lists every review note as a row the person switches between fix now, follow-up and accept; Give back sends only the fix-now notes, Accept in the pane records the rest, follow-ups become tickets with follows set, and the next acceptance starts each note as it was decided."
outcome-why: Alex gives nearly every ticket back to close some of its review notes and had to say in words which; the choice per note and its record on the ticket were missing.
outcome-resolves: "DoD 1-6: tests in core/mod/builtin/jaira-cockpit/tests/board.test.ts 'give back with the review notes'; in-progress prompt names the beheben: list."
review-summary: "none"
round: 4
back-from: testing
back-by: lane
back-count: 2
review-gaps: "none"
test-verdict: "pass (round 2): claude plugin test 252/252 RC=0, go test core/lane core/shipped internal/cli RC=0, shipped-history --check RC=0; the four hunt findings each have a test (kept follow-up creates no ticket, '-race...' goes as --title=, refused Accept writes/creates nothing and accepted one notes after the move, dazu: under beheben: and read past); DoD 1-6 proofs still name existing tests"
hunt-verdict: "fail (round 2): bughunt change blocker - jaira create needs the title positional, --title= alone exits 2, every follow-up fails"
---

# In signoff hakt die Person die Notes an, die behoben werden sollen, und Give back nimmt genau diese mit

## Definition of Done

- [x] Jede Note aus review-gaps steht mit Haken in der Ansicht; Blocker und kleine Raender sind vorausgewaehlt, bekannte Grenzen und fremde Notes nicht (Test)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'every note of the review is a row; blockers and small edges start as fix now...' + 'the Explain note classes a known limit...'
- [x] Give back mit Auswahl schreibt eine Note 'Given back from signoff: beheben: <Notes woertlich>' und der Worker arbeitet genau diese ab; nicht angehakte Notes bleiben als akzeptiert vermerkt (Test)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'Give back writes exactly the notes to fix...'; core/lane/builtin/20-in-progress.md 'A give-back with decisions is the list'
- [x] Eine nicht angehakte Note, die kein Blocker ist, haelt Accept nicht auf; ein nicht angehakter Blocker fragt einmal nach (Test)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'a note that is no blocker, not marked fix now, does not hold Accept...' + 'a blocker not marked fix now asks once; the second press accepts'
- [x] Eingabeklassen: keine Notes, nur ein Blocker, sehr viele Notes (Scrollen), Note ueber mehrere Zeilen, review-gaps aus frueherer Runde, zwei Personen/Sessions entscheiden gleichzeitig
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'no notes: the pane asks...', 'very many notes are all rows', 'a note wrapped over lines is one row...', 'a review of an earlier round lists its own notes...', 'decided meanwhile by someone else...'
- [x] Je Note waehlt die Person eine von drei Entscheidungen: jetzt beheben (geht mit Give back an den Worker) / spaeter im Folge-Ticket (jaira legt ein Ticket mit follows=<dieses> an, Kontext = Note woertlich) / so annehmen (bekannte Grenze); jede Entscheidung steht danach als Note am Ticket, sodass die naechste Abnahme sie nicht neu fragt (Test)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'a press changes a row...', 'Give back writes exactly...' (follow-up via jaira create --follows), 'the next acceptance starts every note decided before as it was decided...'
- [x] Annahmeregel der Heuristik (Klasse aus der Explain-Note nur bei gleicher Anzahl, Fortsetzungszeilen nur bei Praefixen): eine Note falsch eingeordnet ist besser als eine Note verloren oder zwei verschmolzen - keine Zeile aus review-gaps faellt weg (Tabellentest)
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 'an Explain note that classes fewer notes ... is not matched by position', 'gaps without prefixes ... none merged'; core/mod/builtin/jaira-cockpit/hooks/board.ts gapItems

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] gapItems in board.ts: Notes aus review-gaps mit Klasse, Fortsetzung, Vorwahl, frueheren Entscheidungen (Tabellentest)
- [x] decisionsNote/decisionsOf: Note schreiben und zuruecklesen
- [x] Give-back-Pane: Notezeilen mit Umschalten, Freitext, Give back, Accept mit Rueckfrage
- [x] Status neu lesen vor dem Anwenden, Folge-Tickets mit jaira create --follows
- [x] Tests in board.test.ts fuer alle DoD-Punkte, Release-Note

## Progress
- **2026-11-20 18:03 · Alexander Sacharov** — Interim-Regeln (Alex 20.11.2026, gelten ab sofort bis TW5FE2 ausgeliefert ist):
1. pre-process: die DoD nennt die Fehlerklasse statt des Symptoms, die Eingabeklassen (leer, sehr gross, Sonderzeichen/Kodierung/CRLF, Pfade mit Leerzeichen, Windows, zwei Prozesse gleichzeitig, Fehler auf halbem Weg) und bei einer Heuristik eine Annahmeregel - als DoD-Punkte (jaira dod <id> --add).
2. in-progress: vor dem Move selbst die vier Kategorien pruefen (Nebenlaeufigkeit, Fehlerpfade, ungewoehnliche Eingabe, Rechte). Tests ueber die Klasse (Tabellentest), Geschwister-Stellen per grep suchen und mitziehen, Fix an der Ursache statt am Beispiel. Nicht nur Haken setzen: wer umsetzt, denkt die Raender selbst.
3. bughunt/testing: ein Fund nennt Klasse und Beispiel; der Rueckweg verlangt den Fix der Klasse, nicht des einen Beispiels.
4. Schleifenbremse: Blocker sind nur falsches Ergebnis, Datenverlust, Absturz, Sicherheitsluecke. Nach einem Fix-Durchgang laeuft der Testlauf und ein Bughunt nur ueber den Diff seit dem letzten Hunt, dann review. Nach 3 Rueckgaben aus derselben Lane: Frage an den Teamlead statt weiter schleifen.
5. Kleine Raender aus review-gaps in derselben Runde schliessen statt als note: stehen lassen - Alex schickt sie sonst ohnehin zurueck. Vorbestehendes ausserhalb des Umfangs als 'note: pre-existing: ...'.
- **2026-11-20 18:03 · Alexander Sacharov** — Alex 20.11.2026, Praezisierung: 'wenn Dinge schon da sind oder wir entscheiden, sie im naechsten Ticket zu beheben usw., dann muss man das in der Abnahme auswaehlen und markieren.' Also nicht nur Haken 'beheben ja/nein', sondern je Note: jetzt / Folge-Ticket / annehmen - und die Wahl wird festgehalten.
- **2026-11-21 18:58 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/4205.sock
- **2026-11-21 19:01 · Alexander Sacharov** — parent session: uds:/run/user/1000/cc-socks/2274544.sock (neuer Worker; vorheriger Claim von Prozess 2263045 war tot)
- **2026-11-21 19:03 · Alexander Sacharov** — decision (Worker, offene Punkte selbst entschieden, keiner braucht Alex):
1. DoD 1 (Haken) und DoD 5 (drei Wahlen) widersprechen sich nicht: DoD 5 ist Alex' spaetere Praezisierung. Jede Note bekommt eine Zeile, ein Druck schaltet um: jetzt beheben -> Folge-Ticket -> annehmen. 'Angehakt' = jetzt beheben.
2. Ort: der Give-back-Pane (Knopf Give back im Band). Er zeigt die Notes aus review-gaps (bei frueherer Runde die verdict_gaps, wie das Band), darunter das Freitextfeld wie bisher (optional, wenn Notes da sind) und Accept. Ohne Notes bleibt der Pane wie heute. Der Explain-Pane bleibt unveraendert.
3. Klasse je Note: 'blocker:' -> Blocker, 'note: pre-existing:' -> nicht aus diesem Ticket, sonst kleiner Rand. Hat die Explain-Note genau so viele [Klasse]-Zeilen wie review-gaps Notes, gilt ihre Klasse nach Reihenfolge (so kommt 'bekannte Grenze' herein). Vorwahl: Blocker/kleiner Rand -> jetzt, bekannte Grenze -> annehmen, fremd -> Folge-Ticket.
4. Eine Zeile ohne 'blocker:'/'note:' setzt die vorige Note fort, wenn review-gaps solche Praefixe hat (Note ueber mehrere Zeilen); ohne Praefixe ist jede Zeile eine Note (alte Reviews).
5. Festgehalten wird in der einen Note, die Give back ohnehin schreibt: 'Given back from signoff: beheben:' mit Abschnitten 'Folge-Ticket <ID>:' und 'angenommen:', je Note eine Zeile '- <woertlich>'. Accept schreibt dieselbe Form als 'Accepted from signoff: ...'. Die naechste Abnahme liest beide und waehlt eine schon entschiedene Note (gleicher Text) so vor, wie sie entschieden wurde, mit dem Vermerk 'entschieden'.
6. Gleichzeitig: vor dem Anwenden liest der Pane den Ticketstatus neu; steht das Ticket nicht mehr in der Lane, wird nichts geschrieben und gesagt, wer schneller war. Folge-Tickets erst nach erfolgreichem Move bzw. vor Accept.
7. Accept im Pane mit Notes auf 'jetzt' wird abgelehnt (Hinweis: Give back oder umstellen); ein Blocker nicht auf 'jetzt' fragt einmal (zweiter Druck nimmt an). Der Accept-Knopf im Band bleibt wie er ist.
- **2026-11-21 19:14 · Alexander Sacharov** — in-progress: vier Kategorien geprueft. Nebenlaeufigkeit: Status wird vor dem Anwenden neu gelesen (Test 'decided meanwhile'); zwei Druecke in einer Session sperrt once()/close(). Fehlerpfade: jaira create scheitert -> '(not created: ...)' in der Note, Note scheitert bei Accept -> Meldung mit Text; refused dry-run schreibt nichts. Eingabe: CRLF, Leerzeilen, none, Markdown-Betonung, Fortsetzungszeilen, 40 Notes getestet; Notes stehen als '- <text>' in der Note, jaira erkennt Notes nur am Datum, also bricht '- **blocker:**' nichts. Rechte: keine neuen. Bewusst nicht: der Accept-Knopf im Band fragt nicht nach (nur der im Pane); Freitext gilt nur fuer Give back, Accept im Pane hat kein Feld.
- **2026-11-21 19:15 · Alexander Sacharov** — critique 1: gapItems in hooks/board.ts dupliziert findings() (gleiche Datei) - findings(gaps) nehmen, nur Fortsetzungszeilen zusammenfuehren. Stehen gelassen: Explain-Pane-[give back] ohne Notes-Auswahl (decision-Note Punkt 2), deutsche Ueberschriften beheben:/angenommen: (DoD 2 woertlich).
- **2026-11-21 19:15 · Alexander Sacharov** — critique 2: Fund aus critique 1 behoben (f8ec0a06, gapItems ruft findings()). Nichts Neues.
- **2026-11-21 19:19 · Alexander Sacharov** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/register.tsx:974 (createFollowUps) with board.ts gapItems/decisionsOf - class: a review note decided 'follow-up' before gets a second follow-up ticket every time it is decided again; createFollowUps creates a ticket for every item with choice 'later' and never looks at 'earlier', and decisionsOf strips the handle (FOLLOW) so the ticket already created cannot be reused. Example: round 1 Give back marks 'note: pre-existing: load.go accepts empty files' follow-up -> ABC123; the next review lists the same note word for word, the pane preselects it as [follow-up] (decided before), and Give back or Accept runs jaira create again -> a duplicate ticket on the shared board. Same path: an Accept retried after its move failed (see the Accept finding). ran: node --experimental-strip-types script feeding decisionLines output back through decisionsOf+gapItems -> round2 item0 choice=later earlier=later
- **2026-11-21 19:19 · Alexander Sacharov** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/register.tsx:977 (createFollowUps) with board.ts:750 followTitle - class: a finding whose text after 'note:'/'pre-existing:' starts with '-' (a flag name: --json, -race, --force) is passed as jaira create's positional title and parsed as a flag, so the follow-up the person ordered is never created; the failure only lands as '(not created: unknown flag ...)' inside the ticket note, no fail() pane, the toast says given back/accepted. Example: 'note: --dry-run skips the gate' -> title '--dry-run skips the gate'. ran: jaira create "--dry-run skips the gate" --context ... --json in a scratch repo (exit 2, 'unknown flag'); jaira create "-race is not set in CI" (exit 2, 'unknown shorthand flag'); --context with a leading dash is accepted (exit 0), only the title breaks. A '--' before the title or --title would avoid it
- **2026-11-21 19:20 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/register.tsx:549 (accept, before hook from acceptWith:990) - class: Accept's side effects run before the move they describe; follow-up tickets are created and the 'Accepted from <lane>:' note written after the dry run but before 'jaira move --force', so a move that then fails (timeout of the 30s run, a concurrent move between dry run and move, a ref push refused) leaves the ticket in signoff with a note saying it was accepted and follow-ups already filed; pressing Accept again repeats both (duplicate follow-ups through the finding above). Example: the real move exits non-zero after a passing dry run -> fail() pane, ticket unchanged, note + tickets stay. ran: read (accept/acceptWith call path); not reproduced against a failing move
- **2026-11-21 19:20 · Alexander Sacharov** — bughunt: change: core/mod/builtin/jaira-cockpit/hooks/board.ts:689 (decisionLines) with core/lane/builtin/20-in-progress.md:65 and .jaira/lanes/in-progress.md:65 - class: what the person typed in the Give back pane is written as a 'dazu: ...' line after the last heading's list, and the in-progress prompt names only beheben:/angenommen:/Folge-Ticket: ('fix exactly what stands under beheben:', 'leave both alone'), so the typed request reads as part of the accepted or follow-up list; and a give-back with nothing marked fix now but text typed starts 'Given back from signoff: Folge-Ticket:', which the prompt's pattern 'Given back from <lane>: beheben:' does not match at all. Example: one note fix now, one follow-up, typed 'also check the docs' -> note ends 'Folge-Ticket: / - ABC123: ... / dazu: also check the docs'. ran: node script printing decisionLines output (dazu: line directly under the Folge-Ticket list)
- **2026-11-21 19:20 · Alexander Sacharov** — bughunt: done at cdc8f51a - 4 findings (4 change, 2 blocker, 0 pre-existing) in 9 modules, first hunt; concurrency, error paths, input, permissions checked; files: .jaira/lanes/in-progress.md core/lane/builtin/20-in-progress.md core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/sessions.ts core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/MFYM1A.md core/shipped/shipped.sha256
- **2026-11-21 19:23 · Alexander Sacharov** — Fix nach Bughunt 1 (Klasse statt Beispiel): (1) eine frueher als Folge-Ticket entschiedene Note traegt ihr Ticket (decisionsOf liest den Handle, GapItem.follow) - createFollowUps legt nur fuer Notes ohne Ticket an; (2) jaira create bekommt jeden Wert mit '=' (--title=, --context=, --follows=, --board=), keine Note kann mehr Flag werden; ein nicht angelegtes Folge-Ticket meldet zusaetzlich einen Fehler-Pane; (3) Accept schreibt Entscheidungen und legt Folge-Tickets erst nach erfolgreichem move --force an, wie Give back; (4) 'dazu:' steht immer unter 'beheben:' (auch ohne Note zum Beheben), decisionsOf liest darueber hinweg, Lane-Prompt nennt die dazu:-Zeile. Dazu: Submit im Give-back-Pane laeuft unter once(), gegen doppeltes Enter (Hunter konnte es nicht pruefen).
- **2026-11-21 19:23 · Alexander Sacharov** — critique 3 (nur der Fix 81052fbd): die vier Funde an der Ursache behoben, keine neue Form noetig. Nichts.
- **2026-11-21 19:24 · Alexander Sacharov** — bughunt: change: blocker: core/mod/builtin/jaira-cockpit/hooks/register.tsx:981 - class: createFollowUps passes the title only as --title=, but jaira create still demands a positional title (internal/cli/tickets.go: Args minArgs(1); --title is a hidden override, not a replacement), so every follow-up creation exits 2 on both paths, Give back and Accept; example: any note marked Folge-Ticket gets '(not created: ... create needs at least 1 argument(s), received 0)' and no ticket exists; the tests in tests/board.test.ts mock jaira and assert the broken argv, so they pass; a positional after '--' keeps a leading '-' safe (jaira create -- '-race fails' --context=x exits 0); ran: jaira create --title=foo --context=bar --follows=none --json in a scratch repo (exit 2, usage error); blame: 81052fbd (names MFYM1A)
- **2026-11-21 19:25 · Alexander Sacharov** — bughunt: done at 81052fbd - 1 finding (1 change, 1 blocker, 0 pre-existing) in 8 modules, follow-up hunt since cdc8f51a; concurrency, error paths, input, permissions checked; files: .jaira/lanes/in-progress.md core/lane/builtin/20-in-progress.md core/mod/builtin/jaira-cockpit/hooks/board.ts core/mod/builtin/jaira-cockpit/hooks/register.tsx core/mod/builtin/jaira-cockpit/hooks/sessions.ts core/mod/builtin/jaira-cockpit/hooks/texts.ts core/mod/builtin/jaira-cockpit/tests/board.test.ts core/mod/builtin/jaira-cockpit/types/index.d.ts core/release/notes/MFYM1A.md core/shipped/shipped.sha256
- **2026-11-21 19:26 · Alexander Sacharov** — Fix nach Bughunt 2: jaira create nimmt den Titel nur positional (--title ist versteckt und ueberschreibt nur). Jetzt: Werte mit '=', dann '--', dann der Titel. Gegen echtes jaira in einem Scratch-Repo geprueft: 'jaira create --context=... --follows=<id> --json -- -race fails' -> exit 0, handle/follows/title/context richtig. Die Mock-Tests hatten den falschen argv bestaetigt - deshalb diesmal am echten Binary.
