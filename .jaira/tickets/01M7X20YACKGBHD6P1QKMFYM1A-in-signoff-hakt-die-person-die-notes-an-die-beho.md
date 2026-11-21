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
updated-at: 2026-11-21T19:15:48Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-2283240
claimed-at: 2026-11-21T19:01:13Z
outcome-what: "The cockpit's Give back pane lists every review note as a row the person switches between fix now, follow-up and accept; Give back sends only the fix-now notes, Accept in the pane records the rest, follow-ups become tickets with follows set, and the next acceptance starts each note as it was decided."
outcome-why: Alex gives nearly every ticket back to close some of its review notes and had to say in words which; the choice per note and its record on the ticket were missing.
outcome-resolves: "DoD 1-6: tests in core/mod/builtin/jaira-cockpit/tests/board.test.ts 'give back with the review notes'; in-progress prompt names the beheben: list."
review-summary: "core/mod/builtin/jaira-cockpit/hooks/board.ts gapItems: its first line re-implements findings() in the same file (split, MARKER, oneLine, NO_GAP) - call findings(gaps) and only merge the continuation lines"
round: 2
back-from: critique
back-by: lane
back-count: 1
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
