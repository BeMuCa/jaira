---
id: 01M2FYEHZYFCW7DSNRHY9ET6NC
title: "Der Remote-Name gilt pro Rechner, gebraucht wird er pro Board"
status: critique
ready: true
creator: Alexander Sacharov
goal: "Auf einem Board, dessen Repository den eingestellten Remote nicht hat, funktionieren die ref-Befehle wieder - ohne dass ein Ticket dadurch im falschen Repository landet."
context: |-
  'jaira release' ist auf jedem Board ausser diesem tot. Es bricht ab mit: no remote "upstream".

  Damit fehlt genau die Haelfte des Mechanismus, die ein Mensch braucht: ein Ticket nehmen und, wenn er es doch nicht macht, wieder hergeben. Alex ist am 2026-09-14 auf dem requirementsgenie-Board bei FAQ3MW darauf gelaufen und hat den assignee von Hand aus der Datei geloescht.

  Die Ursache ist eine Einstellung, die auf der falschen Ebene liegt:
  - ~/.jaira/settings.json enthaelt {"remote": "upstream"}. Diese Datei gilt pro Rechner und damit fuer JEDES Board, das auf ihm geoeffnet wird.
  - core/settings/settings.go:130 RemoteName() gibt den eingestellten Namen unbesehen zurueck. Auf den Vorgabewert origin faellt es nur zurueck, wenn der String leer ist.
  - core/gitref/gitref.go:112-117 Repo.remote() macht dasselbe noch einmal.
  - core/gitref/gitref.go:100-105 fuehrt dann 'git remote get-url upstream' aus und liefert ErrNoRepo: no remote %q.

  Welches Repository welche Remotes hat:
  - /home/alex/projects/jaira: origin = sashasoft90/jaira (Fork), upstream = BeMuCa/jaira. Hier stimmt die Einstellung.
  - /home/alex/projects/requirementsgenie: nur origin = git.esprit-engineering.de/.../requirementsgenie. Hier gibt es kein upstream, und jeder ref-Befehl faellt um.

  Die Einstellung kam am 2026-09-13 dazu, als das jaira-Board auf upstream umzog. Sie war fuer dieses eine Repository richtig und hat alle anderen mitgerissen.

  Wichtig fuer die Plan-Lane, sonst wird die naheliegende Loesung die falsche: 'wenn der eingestellte Remote fehlt, nimm origin' ist NICHT sicher. In genau diesem Repository ist origin der Fork. Ein Ticket-Ref, das still nach origin geht statt nach upstream, landet im Fork - und das ist der Fehler, gegen den die Einstellung ueberhaupt eingefuehrt wurde. Ein stiller Rueckfall tauscht einen lauten Abbruch gegen einen leisen Datenverlust.

  Nicht Teil dieses Tickets: wohin die Einstellung genau wandert (Board-Datei, Eintrag je Pfad in ~/.jaira/settings.json, oder etwas Drittes). Das entscheidet die Plan-Lane; die beiden Bedingungen unten gelten fuer jeden dieser Wege.
definition-of-done: |-
  In einem Repository, dessen einziger Remote origin heisst, laeuft 'jaira release <id>' durch, waehrend ~/.jaira/settings.json weiterhin remote: upstream sagt. Nachgestellt an einem Fixture-Repository mit genau einem Remote.

  Im jaira-Repository selbst (origin = Fork, upstream = BeMuCa) geht der Ticket-Ref weiterhin nach upstream. Ein Test haelt das fest, damit die Loesung nicht darin bestehen kann, ueberall still auf origin auszuweichen.

  Der Remote laesst sich pro Board festlegen, nicht nur pro Rechner - auf welchem Weg auch immer die Plan-Lane das loest.

  Bricht eine ref-Operation doch am Remote ab, nennt die Meldung drei Dinge: den eingestellten Namen, die Remotes, die dieses Repository tatsaechlich hat, und den Befehl, der es geradezieht. Nicht nur 'no remote "upstream"'.

  Eine Zeile in core/release/NOTES.md unter ## Unreleased.
tags:
  - cli
blocked-by: []
related: []
commits: []
created-at: 2026-09-14T12:32:09Z
updated-at: 2026-09-14T13:00:32Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-93721
claimed-at: 2026-09-14T12:36:43Z
outcome-what: "Der Remote fuer die Ticket-Refs wird jetzt pro Board aufgeloest statt pro Rechner: gitref.Remotes/BoardRemote lesen den Clone, settings.RemoteFor entscheidet in der Reihenfolge git config jaira.remote > settings.json (nur wenn das Repo den Remote hat) > einziger Remote > lauter Abbruch, und Repo.Usable nennt im Fehlerfall eingestellten Namen, vorhandene Remotes und den korrigierenden Befehl."
outcome-why: "Ein einziges \"remote\": \"upstream\" in ~/.jaira/settings.json galt fuer jedes Board auf dem Rechner und hat auf jedem Repository ohne upstream (requirementsgenie) saemtliche ref-Befehle lahmgelegt - u. a. jaira release, also genau die Haelfte, die ein Mensch zum Zurueckgeben eines Tickets braucht."
outcome-resolves: "jaira release und die uebrigen ref-Befehle laufen auf einem Board mit nur origin durch, waehrend settings.json weiterhin upstream sagt; im jaira-Repo selbst geht der Ref unveraendert nach upstream."
review-summary: |-
  core/gitref/gitref.go:681-731 - Remotes und BoardRemote bauen exec.LookPath, exec.Command und bytes.Buffer neu, obwohl run() in derselben Datei (Zeile 150) und value() (Zeile 174) genau das schon tun, inklusive GIT_TERMINAL_PROMPT=0 und getrenntem stderr; stattdessen (&Repo{Dir: dir}).value("remote") bzw. .value("config", "--local", "--get", "jaira.remote") benutzen - die Funktionen bleiben dabei Paketfunktionen auf einem Verzeichnis.
  internal/cli/boardremote_test.go:51 - handleOf durchsucht die Ausgabe nach einem beliebigen sechsstelligen Grossbuchstaben-Wort; jeder solche Titel- oder Lane-Teil in der Ausgabe trifft genauso. Die uebrigen cli-Tests machen es andersherum: Ticket ueber ticket.At(dir).Create(...) anlegen und ticket.Handle(tk.ID) benutzen (internal/cli/checklist_test.go:14, internal/cli/claimrelease_test.go:34). Diesem Muster folgen und handleOf loeschen.
  core/settings/settings.go:175-177 - der Zweig if dir == "" liefert genau dasselbe wie der Durchlauf darunter (BoardRemote("") ist leer, Remotes("") ist leer, also wird want == RemoteName() zurueckgegeben). Zweig streichen, oder wenn er nur die zwei git-Aufrufe sparen soll, das im Kommentar sagen - als Fall, den es geben kann, liest er sich falsch, s.Root ist nie leer.
  core/settings/settings.go:147 - RemoteName() hat nach dieser Aenderung ausser RemoteFor und den Tests keinen Aufrufer mehr (geprueft mit grep ueber alle Nicht-Test-Dateien); der Kommentar begruendet es mit dem Fall ohne Verzeichnis, den es nicht gibt. Entweder in RemoteFor hineinziehen oder unexportieren.
---

# Der Remote-Name gilt pro Rechner, gebraucht wird er pro Board

## Definition of Done

- [x] In einem Repository, dessen einziger Remote origin heisst, laeuft 'jaira release <id>' durch, waehrend ~/.jaira/settings.json weiterhin remote: upstream sagt. Nachgestellt an einem Fixture-Repository mit genau einem Remote.
  proof: internal/cli/boardremote_test.go TestRefCommandsWorkWhenTheMachineSettingNamesAnAbsentRemote, TestTheRefGoesToTheConfiguredRemoteWhenTheRepositoryHasIt, TestABoardRemoteThatIsGoneStopsLoudly; core/settings/remotefor_test.go TestBoardRemoteWinsAndNeverFallsBack; core/gitref/gitref_test.go TestUsableExplainsAMissingRemote; core/release/NOTES.md:16
- [x] Im jaira-Repository selbst (origin = Fork, upstream = BeMuCa) geht der Ticket-Ref weiterhin nach upstream. Ein Test haelt das fest, damit die Loesung nicht darin bestehen kann, ueberall still auf origin auszuweichen.
  proof: core/settings/remotefor_test.go: TestRemoteForKeepsTheConfiguredRemoteWhereItExists und TestRemoteForDoesNotGuessBetweenSeveralRemotes; internal/cli/boardremote_test.go: TestTheRefGoesToTheConfiguredRemoteWhenTheRepositoryHasIt
- [x] Der Remote laesst sich pro Board festlegen, nicht nur pro Rechner - auf welchem Weg auch immer die Plan-Lane das loest.
  proof: git config jaira.remote <name>, gelesen von gitref.BoardRemote; core/settings/remotefor_test.go: TestBoardRemoteWinsAndNeverFallsBack, TestBoardRemoteIsSharedByWorktrees
- [x] Bricht eine ref-Operation doch am Remote ab, nennt die Meldung drei Dinge: den eingestellten Namen, die Remotes die dieses Repository tatsaechlich hat, und den Befehl der es geradezieht. Nicht nur: no remote "upstream".
  proof: core/gitref/gitref_test.go:274 TestUsableExplainsAMissingRemote prueft alle drei Teile; :293 TestUsableSaysWhenThereAreNoRemotesAtAll fuer den Fall ohne Remotes
- [x] Eine Zeile in core/release/NOTES.md unter ## Unreleased.
  proof: core/release/NOTES.md: eine Zeile unter ## Unreleased

Im jaira-Repository selbst (origin = Fork, upstream = BeMuCa) geht der Ticket-Ref weiterhin nach upstream. Ein Test haelt das fest, damit die Loesung nicht darin bestehen kann, ueberall still auf origin auszuweichen.

Der Remote laesst sich pro Board festlegen, nicht nur pro Rechner - auf welchem Weg auch immer die Plan-Lane das loest.

Bricht eine ref-Operation doch am Remote ab, nennt die Meldung drei Dinge: den eingestellten Namen, die Remotes, die dieses Repository tatsaechlich hat, und den Befehl, der es geradezieht. Nicht nur 'no remote "upstream"'.

Eine Zeile in core/release/NOTES.md unter ## Unreleased.

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/gitref: Remotes(dir) ergaenzen - liest 'git remote' und liefert die Namen, die dieses Repository wirklich hat
- [x] Ablageort festlegen und in core/settings dokumentieren: pro Board steht der Remote in der git-config des Clones (git config jaira.remote <name>)
- [x] settings.RemoteFor(dir) schreiben: Reihenfolge git config jaira.remote > settings.json remote > einziger Remote des Repos > origin
- [x] Regel gegen stillen Rueckfall festschreiben: ein per Board gesetzter Name faellt nie zurueck; der Rechner-Vorgabewert weicht nur aus, wenn das Repo genau EINEN Remote hat
- [x] Fehlender Test 1: Fixture-Repo mit nur origin, JAIRA_HOME-settings sagt remote=upstream -> 'jaira release <id>' laeuft durch
- [x] Fehlender Test 2: Fixture-Repo mit origin=Fork und upstream, settings sagt upstream -> aufgeloest wird upstream, nicht origin
- [x] Fehlender Test 3: Fixture-Repo mit zwei Remotes, von denen keiner wie der eingestellte heisst -> Abbruch, kein Ausweichen
- [x] gitref.Repo.Usable(): Fehlertext um die drei Angaben erweitern - eingestellter Name, vorhandene Remotes, 'git config jaira.remote <name>' als Abhilfe
- [x] Test 4: der Fehlertext enthaelt alle drei Angaben
- [x] Die vier Aufrufstellen auf RemoteFor(dir) umstellen: internal/cli/refs.go:27, internal/cli/snapshot.go:112, internal/cli/fetch.go:103, internal/tui/refs.go:35
- [x] Startkosten pruefen: RemoteFor kostet bis zu zwei git-Aufrufe je Kommando - einmal je Prozess merken oder erst beim ersten Ref-Zugriff aufloesen
- [x] README.md:338 (settings.json-Block) um den Board-Remote und die Reihenfolge ergaenzen
- [x] core/release/NOTES.md: Abschnitt '## Unreleased' oben anlegen und eine Zeile schreiben
- [x] go test ./... und Gegenprobe im jaira-Repo selbst: der Ticket-Ref geht weiterhin nach upstream

## Progress
- **2026-09-14 12:46 · Alexander Sacharov** — Ablageort: git config jaira.remote im Clone, nicht eine Datei im Board.

Drei Wege standen zur Wahl:
- .jaira/config, mitcommittet: faellt aus. Ein Remote-Name gehoert zum Clone, nicht zum Projekt. Auf diesem Rechner heisst BeMuCa/jaira 'upstream'; wer direkt von BeMuCa klont, hat dasselbe Repo als 'origin'. Ein mitgelieferter Wert waere fuer den naechsten Klon falsch. core/settings/settings.go sagt das im Kopfkommentar selbst.
- boards-Abschnitt in ~/.jaira/settings.json, nach Pfad geschluesselt: geht, bringt aber Pfadschluessel mit sich, die bei git worktrees auseinanderlaufen (jedes worktree ist ein anderer Pfad, gemeint ist derselbe Clone).
- git config des Clones: pro Clone per Konstruktion, von allen worktrees geteilt, nicht committebar, von Hand editierbar, und jaira schreibt dort schon (share.go traegt den Merge-Driver in .git/config ein). Damit ist auch die Abhilfe im Fehlertext ein gewoehnlicher git-Befehl und braucht kein neues jaira-Kommando.
- **2026-09-14 12:46 · Alexander Sacharov** — Wie beide Bedingungen der DoD zusammenpassen, ohne still ins falsche Repository zu schreiben.

Die DoD verlangt, dass 'jaira release' in einem Repo mit nur origin durchlaeuft, OBWOHL ~/.jaira/settings.json weiterhin upstream sagt - also ohne dass dort vorher etwas eingestellt wurde. Gleichzeitig darf die Loesung nicht ueberall still auf origin ausweichen.

Die Trennlinie ist, WER den Namen gesetzt hat:
- Der Name aus der git-config des Boards ist eine Entscheidung fuer genau dieses Repository. Fehlt der Remote, wird abgebrochen. Kein Ausweichen.
- Der Name aus ~/.jaira/settings.json ist nur ein Vorgabewert fuer alle Boards. Hat dieses Repo ihn nicht, gilt er hier nicht.

Und dann die Ausweichregel, die den stillen Datenverlust ausschliesst: ausgewichen wird nur, wenn das Repository GENAU EINEN Remote hat. Dann gibt es keine zweite Moeglichkeit, in die etwas fallen koennte.
- requirementsgenie: nur origin -> laeuft, ohne dass jemand etwas einstellt.
- jaira selbst: upstream ist da -> upstream, unveraendert.
- Fork ohne upstream, aber mit origin und einem dritten Remote: mehrdeutig -> lauter Abbruch statt Fork.

'Nimm origin, wenn der eingestellte Remote fehlt' waere genau der Fehler aus dem Kontext: origin ist hier der Fork.
- **2026-09-14 12:46 · Alexander Sacharov** — Wo die Aufloesung hingehoert und was sie kostet.

Heute faellt der Name an vier Stellen aus settings.RemoteName(), das Repository sieht keine davon: internal/cli/refs.go:27, internal/cli/snapshot.go:112, internal/cli/fetch.go:103, internal/tui/refs.go:35. Deshalb eine Funktion settings.RemoteFor(dir) statt vier Mal derselbe Sonderfall. RemoteName() bleibt fuer den Fall ohne Verzeichnis stehen.

Das Lesen der Remotes gehoert nach core/gitref (kennt das Repo, schaltet schon git auf). core/settings importiert gitref bereits; umgekehrt darf es nicht sein.

Kosten: RemoteFor macht bis zu zwei git-Aufrufe (config get, remote). attachRefs laeuft bei JEDEM Kommando, auch bei 'jaira list' - und 'Instant startup' ist eine Projektbedingung. Falls das messbar wird: entweder je Prozess einmal merken oder den Namen erst beim ersten Ref-Zugriff aufloesen (refsync.New nimmt dann eine Funktion statt eines Strings). Erst messen, dann entscheiden - Schritt 11.

core/release/NOTES.md hat zur Zeit keinen Abschnitt '## Unreleased'; der oberste ist ## 0.2.0. Der Abschnitt muss neu angelegt werden.
- **2026-09-14 12:50 · Alexander Sacharov** — Ablageort ist git config --local jaira.remote; gelesen von gitref.BoardRemote(dir). Warum nicht memoisieren: RemoteFor wird genau einmal je Prozess aufgerufen, in attachRefs. snapshot.go und fetch.go lesen den aufgeloesten Namen jetzt von refs.Repo.Remote ab, statt ihn ein zweites und drittes Mal aufzuloesen. Damit kosten die bis zu zwei git-Aufrufe nur einmal, und es braucht keinen Cache, der in Tests veraltet.
- **2026-09-14 12:54 · Alexander Sacharov** — Gegenprobe gelaufen: settings.RemoteFor loest sowohl /home/alex/projects/jaira als auch das worktree .worktrees/jaira-9ET6NC auf 'upstream' auf - unveraendert. Verworfen: 'fehlt der eingestellte Remote, nimm origin'. In diesem Repo ist origin der Fork; das waere stiller Datenverlust. Stattdessen weicht nur der Rechner-Vorgabewert aus, und auch nur wenn das Repo genau EINEN Remote hat. Gegenprobe, dass der Test 1 den Fehler wirklich nachstellt: RemoteFor testweise auf RemoteName zurueckgedreht -> TestRefCommandsWorkWhenTheMachineSettingNamesAnAbsentRemote faellt mit genau 'no remote "upstream"' um.
- **2026-09-14 12:59 · Alexander Sacharov** — critique (1. Durchgang): Der Entwurf selbst bleibt stehen. git config jaira.remote als Ablageort ist begruendet (pro Clone, von worktrees geteilt, nicht committet, jaira schreibt in share.go schon in .git/config), die Reihenfolge in RemoteFor deckt beide Haelften der DoD ab, und die Schichtung stimmt: core/settings importiert core/gitref, nicht umgekehrt. Dass fetch.go und snapshot.go den aufgeloesten Namen von refs.Repo.Remote ablesen statt ein zweites Mal aufzuloesen, ist die richtige Stelle - nicht angefasst.

Vier Befunde, alle lokal, keiner ruehrt an der Form:
1. core/gitref/gitref.go:681-731 - Remotes/BoardRemote dupliziert die exec-Mechanik von run()/value() in derselben Datei. Ueber (&Repo{Dir: dir}).value(...) gehen, dann gilt auch dort GIT_TERMINAL_PROMPT=0 und die getrennte stderr-Behandlung.
2. internal/cli/boardremote_test.go:51 - handleOf raet das Handle aus der stdout-Ausgabe (jedes sechsstellige Grossbuchstaben-Wort trifft). Die anderen cli-Tests legen das Ticket ueber ticket.At(dir).Create(...) an und nehmen ticket.Handle(tk.ID); diesem Muster folgen.
3. core/settings/settings.go:175-177 - der dir==""-Zweig ist mit dem Durchlauf darunter gleichbedeutend.
4. core/settings/settings.go:147 - RemoteName() hat keinen Aufrufer mehr ausser RemoteFor.

Nicht beanstandet, damit es niemand nochmal aufmacht: dass ein Repo mit genau EINEM Remote diesen still nimmt, auch wenn settings.json etwas anderes sagt - das verlangt die DoD ausdruecklich (der Fall requirementsgenie soll ohne Einstellung laufen).
- **2026-09-14 13:00 · Alexander Sacharov** — DoD-Reparatur durch den Dispatcher: beim Anlegen sind vier der fuenf Kriterien verloren gegangen, weil "jaira create --dod" nur EIN Item nimmt - die Absaetze 2-5 landeten als Fliesstext im Body unter "## Definition of Done", ohne Checkbox. Damit hat das Gate der Terminal-Lane nur 1 von 5 Kriterien geprueft. Die vier fehlenden sind mit "jaira dod --add" nachgetragen und gegen die vorhandene Implementierung geprueft: alle vier waren bereits gebaut und sind mit Testnamen als Proof abgehakt. Nichts an der Implementierung widerspricht ihnen.
