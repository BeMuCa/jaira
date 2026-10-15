---
id: 01M21PJTWPK4DZ89PRY6T8A8KM
title: "Eine Pille unter der Version zeigt, dass ein neues Release da ist"
status: pre-process
ready: true
creator: BeMuCa
assignee: BeMuCa
goal: "Steht im Cache ein neueres Release als die laufende Binary, zeigt der Projektkopf unter der Versionszeile eine kleine Pille mit der Zielversion"
context: |-
  Berk am 08.09.: 'einmal woechentlich soll gechecked werden ob eine neue Version existiert und wenn ja soll im Projektfenster unter der Version ein kleines Pille sein'.
  Der Check existiert schon und muss NICHT gebaut werden. core/selfupdate/cache.go: MaxAge ist 24h, PollCache liest den Cache, stempelt ihn bei Ablauf frisch und startet einen abgekoppelten Kindprozess 'jaira self upgrade --check --json'. Der holt per HTTPS https://api.github.com/repos/BeMuCa/jaira/releases/latest. Kein git, kein Blockieren der UI.
  Berk hat am 08.09. entschieden: das 24h-Intervall bleibt wie es ist, woechentlich waere seltener als heute. Nur die Darstellung ist neu.
  Gewuenschter Text der Pille, von Berk gewaehlt: ein Aufwaerts-Pfeil und die Zielversion, also '^ 0.1.2'. Nicht 'new release', nicht 'update ready'.
  Ort: internal/tui/updatecheck.go, versionLine. Die Funktion liefert heute eine Zeile fuer die Fusszeile mit 'jaira 0.1.1 - up to date' bzw. '... - 0.1.2 available - run: jaira self upgrade'.
  Reihenfolge beachten: P1AE82 verlegt die Versionszeile erst nach links oben in den Projektkopf und laesst die Fusszeile weg. Dieses Ticket setzt die Pille darunter und kommt danach.
  Bestehende Regel nicht brechen: ein dev-Build und ein Cache, der nie beschrieben wurde, duerfen nichts behaupten. Ohne bekanntes Release keine Pille.
definition-of-done: "Ist im Cache ein neueres Release als release.Current, steht im Projektkopf unter der Version eine Pille '^ <version>'; bei gleicher Version und bei unbekanntem Cache steht keine Pille; ein Test deckt alle drei Faelle; go test ./... -race gruen"
tags:
  - release
blocked-by: []
follows: 01M1KCKH0MCQ5P1BKTSHP1AE82
commits: []
created-at: 2026-09-08T23:45:19Z
updated-at: 2026-10-15T12:20:10Z
updated-by: BeMuCa
related:
  - 01M2HX6SYCFB27V26R5APWAF22
---

# Eine Pille unter der Version zeigt, dass ein neues Release da ist

## Definition of Done

- [ ] Ist im Cache ein neueres Release als release.Current, steht im Projektkopf unter der Version eine Pille '^ <version>'; bei gleicher Version und bei unbekanntem Cache steht keine Pille; ein Test deckt alle drei Faelle; go test ./... -race gruen
- [ ] Ein Agent, der mit jaira arbeitet, erfaehrt von einem neueren Release und sagt es dem Nutzer: jaira nennt es in einer Ausgabe, die ein Agent sieht (z.B. jaira resume/next), und der Agent-Block sagt 'meldet jaira ein neueres Release, sag es dem Nutzer und nenne jaira self upgrade'; Test
- [ ] Nach einem Upgrade oder 'jaira update' zeigt der TUI einmal ein Popup mit den Release-Notes seit der Version, die zuletzt auf diesem Board lief; eine Taste schliesst es; bis zum naechsten Versionswechsel erscheint es nicht wieder; Test
- [ ] Eine Release-Note-Zeile kann als wichtig markiert werden (Vorschlag: '- ! <text>', bleibt fuer alte Binaries eine normale Zeile); sie steht im Popup und in der Ausgabe von 'jaira update' hervorgehoben; der Abschnitt zu core/release/notes/ in CLAUDE.md nennt die Markierung; Test
- [ ] Auf der Launcher-Seite oeffnet ein Hotkey die Release-Notes (neueste Version zuerst); die Hilfezeile nennt die Taste; Test

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-13 20:51 · BeMuCa** — Berk 2026-10-13: der Agent soll einen auf ein neues Release ansprechen (DoD 2). Zur Frage 'feste Uhrzeit?': nein - core/selfupdate/cache.go PollCache prueft, sobald jaira laeuft und die letzte Pruefung aelter als 24 h ist (MaxAge, cache.go:20); wer erst spaeter online kommt, bekommt die Pruefung dann. Das Upgrade selbst bleibt ein Befehl (jaira self upgrade), es laeuft nicht von allein. Gehoert zu PWAF22: nach dem Upgrade muessen Block und Skills nachgezogen werden.
- **2026-10-14 07:47 · BeMuCa** — Berk 2026-10-14: 'das will ich haben' - nach Upgrade/Update ein Popup mit den Release-Notes, wichtige Aenderungen staerker markiert, und ein Hotkey im Launcher, der die Release-Notes zeigt (DoD 3-5). Teil davon gibt es schon: die CLI mahnt 'jaira update' an (internal/cli/update.go:53), 'jaira update' druckt die Notizen seit der letzten Version (update.go:94-95). Anlass: EBTX6X aendert das Milestone-Format, und das Team muss es erfahren. Formatwechsel fuer alte Builds und Agenten: eigenes Ticket, folgt EBTX6X.
- **2026-10-15 12:19 · BeMuCa** — Popup-Merker (DoD 3): neue Datei 'notes-seen' im StateDir des Boards, neben dem Stempel 'version' (StateDir = JAIRA_HOME bzw. ~/.jaira + state/<arbeitsbaum>, core/ticket/store.go:701). Eine Zeile: die Version, deren Notes zuletzt geschlossen wurden.
Warum dort: pro Nutzer und pro Arbeitsbaum wie der Stempel, nie im Repo. 'Ich habe die Notes gelesen' darf nicht per git reisen - sonst quittiert ein Teammitglied fuer alle. Eine Zeile, von Hand editierbar, kein Merge-Konflikt moeglich.
Warum nicht der Stempel allein: 'jaira update' stempelt (update.go:112), und PWAF22 laesst 'self upgrade' update automatisch laufen. Danach ist Stamped==Current und Since(Stamped) leer - der TUI wuesste nicht mehr, ab welcher Version.
Regel: base = notes-seen; fehlt die Datei, base = Stempel. Popup nur wenn release.Newer(Current, base) und base nicht 'dev'. Geschrieben wird notes-seen beim Schliessen des Popups und wenn es nichts zu zeigen gibt. 'jaira update' schreibt es NICHT.
Bekannte Luecke: fehlt notes-seen UND hat update schon gestempelt (nur beim allerersten Wechsel, bevor der TUI auf diesem Baum je lief), kommt kein Popup. Die Notes hat update dann schon gedruckt.
Pro Baum heisst: jeder Worktree mit altem Stempel zeigt es einmal; ein neuer Worktree (kein Stempel) nie. Nur in New (model.go:313), nicht in switchBoard: ein per 1-9 gewechseltes Board zeigt es beim naechsten Oeffnen.
- **2026-10-15 12:19 · BeMuCa** — Pille (DoD 1): versionLine ruft selfupdate.PollCache weiter genau einmal, bei der Konstruktion (home.go:134, model.go:318). Kein Netz, kein Blockieren: PollCache liest nur update-check.json und startet bei >24h einen abgekoppelten Kindprozess (core/selfupdate/cache.go:181).
'Neuer' per release.Newer (bisher unexportiertes newer, release.go:86), nicht per '!='. Ein Cache mit aelterer Version als die laufende Binary darf keine Pille zeigen - heute sagt versionLine in dem Fall faelschlich 'available'.
Der Schwanz '. X available - run: jaira self upgrade' faellt weg: Pille und Zeile wuerden sonst dasselbe zweimal sagen. Den Befehl nennt der Agent (DoD 2). Wer den Hinweis fuer Menschen behalten will, setzt ihn gedimmt neben die Pille - Geschmacksfrage fuer critique/Berk, kein Blocker.
renderBoard zaehlt die Kopfzeilen schon (view.go:256 headLines): eine zweite Zeile im head kostet die Spalten automatisch eine Reihe. P1AE82 ist gelandet (fac61b4e), die Version steht links oben; TestBoardVersionSitsInTheTopLeftCorner haelt die Zeile darunter frei.
Der Launcher bekommt die Pille auch: er zeigt dieselbe Versionszeile, und ohne den weggefallenen Schwanz verloere er sonst die Info.
- **2026-10-15 12:19 · BeMuCa** — Agent-Ausgabe (DoD 2): eine Zeile auf stderr bei 'jaira next' (alle Varianten, vor dem --per-lane-Zweig, flow.go:407) und 'jaira resume'. Das sind die Befehle, mit denen eine Sitzung startet.
stderr, weil stdout das JSON ist, das ein Agent parst - dieselbe Regel wie nudgeIfStale (update.go:28). Der Stop-Hook (hook.go:44) leitet stderr nach /dev/null und bleibt unberuehrt.
Nicht auf jedem Befehl: genau das war der Grund fuer die bisherige CLI-Stille (Kommentar update.go:31-36). Dieser Kommentar und die Doku von versionLine/PollCache werden falsch und muessen mitgeaendert werden.
next/resume rufen dafuer PollCache auf. Erst dadurch stimmt Berks Satz vom 13.10. ('prueft, sobald jaira laeuft'): heute prueft NUR der TUI, wer nur per Agent arbeitet, erfuhr nie von einem Release. Kosten: hoechstens ein abgekoppelter Kindprozess je 24h; JAIRA_NO_UPDATE_CHECK=1 schaltet ab; Testbinaries spawnen nie (cache.go:144).
Der Agent-Block ist englisch, der Satz aus der DoD steht dort als 'when jaira reports a newer release ..., tell the user and name jaira self upgrade', direkt nach dem jaira-resume-Punkt (announce.go:110). Die Bloecke in CLAUDE.md/AGENTS.md dieses Repos werden mit der frisch gebauten Binary neu erzeugt (Vorbild 81eccd1f).
- **2026-10-15 12:20 · BeMuCa** — Launcher-Taste (DoD 5): 'w' (what's new). Auf dem Launcher belegt: a r s x d q j k enter esc ctrl+c (home.go:258-305).
'n' waere frei auf dem Launcher, heisst auf dem Board aber 'neues Ticket' (model.go:1645); 'v' ist auf dem Board die Pipeline (model.go:1587). Beide haetten je nach Bildschirm zwei Bedeutungen. 'w' ist auf Launcher und Board frei (nur die Pipeline nutzt es als WASD-hoch, eigener Bildschirm).
Inhalt: release.Since('') - alle Releases, leere weggelassen, neueste zuerst (composeNotes sortiert absteigend, NOTES.md ist absteigend). Derselbe notesView zeigt die Notes im Popup (DoD 3) und hier.
- **2026-10-15 12:20 · BeMuCa** — Marker '- ! <text>' (DoD 4): changeLines (release.go:100) bleibt unveraendert und liefert daraus die Change '! <text>'. Jede Binary liest die Zeile also als normale Aenderung; das Format bleibt 'eine Zeile, die mit - beginnt'.
Neu nur release.Important(change) (text, ok): Praefix '! ' MIT Leerzeichen; '!text' zaehlt nicht. Kein neues Feld an release.Entry und kein neues JSON-Feld: 'jaira update --json' traegt '! <text>' roh in 'changes', ein Agent sieht das Zeichen.
Hervorhebung: im TUI farbig (colWarn, fett, Bullet '!'); in 'jaira update' ohne ANSI als '  ! <text>' statt '  - <text>', weil die Ausgabe auch Agenten und Pipes lesen.
Hinweis: alte Binaries lesen neue Notes nie (go:embed). 'Bleibt fuer alte Binaries normal' ist damit ein Formatversprechen fuer Menschen und Werkzeuge, kein Laufzeitfall. Der Test belegt es an changeLines.
Doku: Abschnitt 'Client-facing changes go in core/release/notes/' steht in CLAUDE.md UND AGENTS.md, beide ausserhalb des jaira-Blocks; dazu der Kommentar an changeLines. Workflows (.github/workflows/release-notes.yaml) verschieben nur Dateien und lesen das Format nicht.
- **2026-10-15 12:20 · BeMuCa** — Kollisionen mit PWAF22 (dort noch kein Code, gleiche Basis 0aefa5c1):
(1) internal/cli/update.go - hier nur der Kommentar update.go:31-36 und die Notes-Schleife update.go:132-137, die zu writeNotes(w, notes) wird. PWAF22 baut RunE um; ein kleiner Konflikt dort ist wahrscheinlich.
(2) core/release/release.go - newer wird Newer (exportiert). PWAF22 braucht fuer 'Skills aelter als Binary' evtl. denselben Versionsvergleich: dann Newer nutzen statt einen zweiten zu bauen.
(3) core/board/announce.go agentNote - eine Zeile nach dem jaira-resume-Punkt. Aendert PWAF22 den Block auch, nach dem Merge die Bloecke in CLAUDE.md/AGENTS.md neu erzeugen.
(4) Release-Notes: je eigene Datei, kein Konflikt.
Prozess: DoD 1 verlangt 'go test ./... -race gruen', die Implementierer duerfen auf Berks Anweisung kein go test laufen lassen - den Beleg muss die testing-Lane oder CI liefern.
