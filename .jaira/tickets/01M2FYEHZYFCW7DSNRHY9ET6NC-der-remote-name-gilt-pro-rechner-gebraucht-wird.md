---
id: 01M2FYEHZYFCW7DSNRHY9ET6NC
title: "Der Remote-Name gilt pro Rechner, gebraucht wird er pro Board"
status: review
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
updated-at: 2026-09-14T13:21:34Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-93721
claimed-at: 2026-09-14T12:36:43Z
outcome-what: "Vier lokale Critique-Findings behoben: gitref.Remotes/BoardRemote nutzen jetzt Repo.value statt eigener exec-Aufrufe, der handle-Scan im cli-Test ist weg, der tote dir==\"\"-Zweig und RemoteName sind entfernt."
outcome-why: "Die Critique-Lane hat den Entwurf angenommen, aber vier lokale Doppelungen und eine rateende Testhilfe beanstandet."
outcome-resolves: "Critique-Runde 1, alle vier Findings."
review-summary: "Der Remote-Name fuer Ticket-Refs wird nicht mehr blind aus ~/.jaira/settings.json genommen, sondern gegen das Repository aufgeloest, das gerade offen ist. Neu in core/settings: Settings.RemoteFor(dir) mit vier Stufen - (1) \"git config --local jaira.remote\" im Clone gilt unbesehen und faellt nie zurueck, (2) der Rechner-Vorgabewert aus settings.json gilt nur, wenn dieses Repo einen Remote dieses Namens wirklich hat, (3) hat das Repo genau EINEN Remote, wird der genommen, (4) sonst bleibt der eingestellte Name stehen und der Befehl bricht laut ab. Dafuer liefert core/gitref zwei neue Funktionen: Remotes(dir) (git remote) und BoardRemote(dir) (git config --local --get jaira.remote), beide ueber Repo.value, also mit GIT_TERMINAL_PROMPT=0. Repo.Usable() meldet ueber den neuen noRemote() statt \"no remote \\\"upstream\\\"\" jetzt drei Dinge: gesuchter Name, tatsaechlich vorhandene Remotes (bzw. \"this repository has no remotes\") und \"git -C <dir> config jaira.remote <name>\" als Abhilfe. Aufgeloest wird genau einmal je Prozess (internal/cli/refs.go attachRefs, internal/tui/refs.go newSyncer); snapshot.go und fetch.go lesen den fertigen Namen von refs.Repo.Remote ab statt ihn erneut aufzuloesen, RemoteName() ist ersatzlos weg. Dazu README-Abschnitt mit der Reihenfolge und eine Zeile unter ## Unreleased in core/release/NOTES.md. Damit laeuft requirementsgenie (nur origin) ohne jede Einstellung, jaira selbst (origin=Fork, upstream=BeMuCa) geht weiterhin nach upstream, und ein Ticket kann nicht still im Fork landen."
review-gaps: "Nichts, was zurueckgehen muesste. Alle fuenf DoD-Kriterien sind im Diff belegt, nicht nur behauptet; die Schichtung (core/settings -> core/gitref) stimmt, es gibt genau eine Stelle, die gitref.Repo baut (core/refsync/refsync.go:90), und RemoteFor kann nie \"\" liefern, also bleibt Landing/refs.Repo.Remote so belastbar wie vorher. Eigene Gegenproben: go vet ./... sauber, go test ./core/settings ./core/gitref ./internal/cli gruen. Zwei kleine Beobachtungen, beide kein Blocker und keine Ruecksendung: (1) internal/cli/boardremote_test.go TestTheRefGoesToTheConfiguredRemoteWhenTheRepositoryHasIt haengt origin UND upstream an dasselbe bare-Repo (cloneWithRemotes gibt beide Male denselben Pfad). Der Test prueft deshalb nur refs.Repo.Remote == \"upstream\" und koennte gar nicht merken, wenn der Ref in origin gelandet waere - sein Kommentar (\"origin here is a fork with nothing in it\") verspricht mehr, als er einloest. Abgedeckt ist die Aussage trotzdem: core/settings/remotefor_test.go TestRemoteForKeepsTheConfiguredRemoteWhereItExists und die Handprobe der Testing-Lane mit zwei getrennten Remotes (origin+fleet). (2) attachRefs kostet jetzt bis zu zwei zusaetzliche git-Prozesse je Kommando (config --get, remote) - einmal je Prozess, im Millisekundenbereich, und Schritt 11 des Plans hat es bewusst ohne Cache gelassen. Vorbestehend und nicht von dieser Aenderung: \"jaira create\" bricht bei unaufloesbarem Remote nicht ab und zeigt die neue Meldung nicht (internal/cli/refs.go:75 fileOnRefOnly); die Meldung kommt bei pull, fetch, release, snapshot."
test-verdict: "pass: go build/vet/test ./... gruen (26 Pakete, 0 Fehler); Regressionsnachweis - Elternstand 1c6be9d stirbt am selben Fixture mit no remote \"upstream\" (Exit 1), der Branch laeuft durch; alle fuenf DoD-Kriterien am gebauten Binary auf Wegwerf-Fixtures nachgestellt, nicht nur aus Tests gelesen"
review-verdict: "Der Diff deckt alle fuenf DoD-Kriterien ab, und er loest genau den Fehler aus dem Kontext, ohne die Falle aufzumachen, vor der der Kontext warnt: es gibt keinen stillen Rueckfall auf origin - ausgewichen wird nur, wenn das Repo genau EINEN Remote hat, und ein per Board gesetzter Name faellt nie zurueck. Keine Defekte gefunden, nichts im outcome, was der Diff nicht traegt. Zwei Kleinigkeiten stehen in review-gaps (ein Test, dessen Kommentar mehr verspricht als er prueft; zwei git-Aufrufe mehr beim Start) - beide rechtfertigen keine Ruecksendung. Empfehlung: annehmen."
review-check: "Alles laeuft in Wegwerf-Verzeichnissen; /home/alex/.jaira und das echte Board werden nicht angefasst. 1. Binary bauen: cd /home/alex/projects/.worktrees/jaira-9ET6NC && go build -o /tmp/jc/jaira ./cmd/jaira (vorher mkdir -p /tmp/jc). 2. Wegwerf-Einstellungen anlegen, die den Fehler ausloesen: mkdir -p /tmp/jc/home && echo \"{\\\"remote\\\": \\\"upstream\\\"}\" > /tmp/jc/home/settings.json && export JAIRA_HOME=/tmp/jc/home JAIRA_USER=ada. 3. Repo mit NUR origin bauen: git init --bare -q /tmp/jc/solo.git && git clone -q /tmp/jc/solo.git /tmp/jc/solo && git -C /tmp/jc/solo config user.name ada && git -C /tmp/jc/solo config user.email ada@example.test. 4. cd /tmp/jc/solo && /tmp/jc/jaira init && /tmp/jc/jaira create \"check me\" --goal g --context c --dod d. 5. /tmp/jc/jaira list - in der ersten Spalte steht ein sechsstelliges Handle; das ist unten <H>. 6. git -C /tmp/jc/solo.git for-each-ref refs/jaira/ - es MUSS eine Zeile refs/jaira/tickets/... erscheinen: der Ref ist in origin gelandet, obwohl settings.json upstream sagt. 7. /tmp/jc/jaira pull <H> && /tmp/jc/jaira claim <H> && /tmp/jc/jaira release <H> - alle drei laufen durch, release meldet \"<H> is back on the board\". Genau hier stand vorher: no remote \"upstream\". 8. Gegenprobe alter Stand (optional): dieselben Schritte mit einem Binary aus Commit 1c6be9d - Schritt 7 bricht mit no remote \"upstream\" ab, Exit 1. 9. Lauten Abbruch pruefen: git -C /tmp/jc/solo config --local jaira.remote ghost && /tmp/jc/jaira pull <H> - Exit 1 und die Meldung nennt alle drei Teile: no remote \"ghost\" - this repository has origin / set the one this board uses: git -C /tmp/jc/solo config jaira.remote <name>. 10. Pro Board einstellbar: git -C /tmp/jc/solo config --local jaira.remote origin && /tmp/jc/jaira pull <H> - laeuft wieder durch. 11. Eigenes Repo unveraendert: cd /home/alex/projects/jaira && go test ./core/settings/... ./core/gitref/... ./internal/cli/... - gruen; hier ist origin der Fork und upstream BeMuCa, aufgeloest wird weiterhin upstream (Test TestRemoteForKeepsTheConfiguredRemoteWhereItExists). 12. Aufraeumen: rm -rf /tmp/jc."
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
- **2026-09-14 13:05 · Alexander Sacharov** — critique (1. Durchgang), lokale Befunde abgearbeitet. Was dabei herauskam und nicht im Code steht:

- gitref.Remotes/BoardRemote nutzen jetzt (&Repo{Dir: dir}).value(...). Verhalten unveraendert, aber value() setzt GIT_TERMINAL_PROMPT=0 - das fehlte in der handgebauten exec.Command-Variante und haette bei einem credential-Prompt haengen koennen.
- RemoteName ist weg, nicht nur unexportiert: die vier Aufrufe in core/settings/settings_test.go pruefen jetzt ueber RemoteFor(t.TempDir()). Ein TempDir ist kein Repo, also liefern BoardRemote und Remotes nichts und RemoteFor faellt auf denselben Vorgabewert durch, den RemoteName geliefert hat. Dadurch braucht kein Test mehr eine zweite Einstiegstuer.
- Der 'if dir == ""'-Zweig in RemoteFor ist weg. Kein Aufrufer uebergibt ""; s.Root ist nie leer.
- Abweichung von der Kritik bei boardremote_test.go: ticket.At(dir).Create(...) allein reicht hier NICHT. 'jaira create' legt das Ticket auf seinem git-ref ab und nicht auf der Platte ('On its ref, not on your disk'), und genau dieser ref ist das, was die drei Tests pruefen. Store-Create baut keinen ref, und 'release' scheitert dann mit 'no ref for this ticket'. Deshalb bleibt das Anlegen ueber die CLI; nur das Raten des Handles aus stdout ist weg - der Handle kommt aus (&gitref.Repo{Dir: clone}).List(). Und weil im dritten Test der Board-Remote fehlt, kann create keinen ref schreiben und legt das Ticket auf die Platte; dafuer faellt der Helfer auf ticket.At(clone).List() zurueck.
- **2026-09-14 13:08 · Alexander Sacharov** — critique (2. Durchgang): keine Befunde. Alle vier Findings aus Runde 1 sind zu: gitref.Remotes/BoardRemote gehen ueber (&Repo{Dir: dir}).value(...) und erben damit GIT_TERMINAL_PROMPT=0 und die getrennte stderr-Behandlung; der if dir==""-Zweig in RemoteFor ist weg; RemoteName ist ganz entfallen und in RemoteFor hineingezogen, die vier Aufrufe in settings_test.go gehen jetzt ueber RemoteFor(t.TempDir()); handleOf ist weg.

Zur bewussten Abweichung bei Finding 2 (ticket.At(dir).Create statt CLI): die Begruendung traegt. internal/cli/refs.go:74 fileOnRefOnly loescht die lokale Datei, sobald der Flush den ref gesendet hat - auf einem Board mit brauchbarem Remote lebt ein frisches Ticket also NUR auf refs/jaira/tickets/<id>. Store-Create wuerde eine Datei ohne ref anlegen und damit genau den Mechanismus umgehen, den diese drei Tests pruefen. Die Muster-Tests (claimrelease_test.go) laufen in einem t.TempDir() ohne git, dort ist refsync ohnehin inert - deshalb passt ihr Muster hier nicht. Das Raten aus stdout ist das eigentliche Problem gewesen und ist weg; der Handle kommt aus gitref.Repo.List().

Der Store-Rueckfall in ticketIn ist kein toter Zweig: im dritten Test scheitert Usable() am fehlenden Board-Remote, refsync ist inert, die Datei bleibt liegen und es gibt keinen ref.

Nichts Neues eingeschleppt: go vet ./core/... ./internal/... laeuft sauber, Remotes/BoardRemote liefern bei fehlendem git, fehlendem Repo und leerer Ausgabe dasselbe wie vorher (nil bzw. ""), und value() trimmt bereits, was BoardRemote vorher von Hand tat.
- **2026-09-14 13:12 · Alexander Sacharov** — optimize: keine Aenderung noetig. Vier Durchgaenge (Duplikation, toter Code, Fluff, Kosten) im Detail in review-gaps. Zwei Dinge bewusst NICHT angefasst, damit die naechste Runde sie nicht neu aufmacht: der tote 'remote == ""'-Zweig in Settings.Landing ist vorbestehend (auch RemoteName gab nie "" zurueck) und gehoert nicht zu diesem Ticket; das doppelte strings.TrimSpace in RemoteFor ist folgenlos, aber das Entfernen wuerde RemoteFor an eine undokumentierte Zusicherung von Repo.value binden.
- **2026-09-14 13:16 · Alexander Sacharov** — Testing-Lane: Verdikt BESTANDEN.

Suite: go build ./... ok, go vet ./... ok, go test ./... alle 26 Pakete ok, 0 Fehler.
Die neun in den DoD-Proofs genannten Tests existieren und laufen gruen.

Gegenprobe mit dem Elternstand (1c6be9d), gleiches Fixture-Repo:
  jaira pull/release -> 'no remote "upstream"', Exit 1.
Mit dem Stand dieses Branches laeuft dasselbe durch. Der Fehler aus dem Kontext ist also wirklich weg, nicht nur wegdefiniert.

Am Binary nachgestellt (Build aus diesem Worktree, JAIRA_HOME-Fixture, Wegwerf-Repos im Scratchpad; das echte Board und ~/.jaira wurden nicht angefasst):
1) Repo mit nur origin, settings.json sagt remote=upstream: create -> Ref landet in origin.git; pull, claim, release laufen durch, Exit 0, assignee danach leer. BESTANDEN.
2) jaira-Repo selbst (origin=Fork, upstream=BeMuCa, kein jaira.remote gesetzt): RemoteFor loest zu "upstream" auf - nur aufgeloest, nichts gepusht. BESTANDEN.
3) Fixture mit origin+fleet, 'git config jaira.remote fleet': der Ticket-Ref liegt danach in fleet, origin hat keinen Ticket-Ref. Pro Board einstellbar. BESTANDEN.
4) Fehlermeldung nennt alle drei Teile, in vier Faellen geprueft (ghost-Board-Remote, settings-Remote fehlt bei zwei Remotes, Repo ganz ohne Remotes):
   no remote "ghost" - this repository has fleet, origin
     set the one this board uses: git -C <pfad> config jaira.remote <name>
   Ohne Remotes: 'this repository has no remotes'. BESTANDEN.
5) core/release/NOTES.md: eine Zeile unter ## Unreleased, eine Zeile, kein Umbruch. BESTANDEN.

Beobachtung, kein Fehler: 'jaira create' bricht bei unaufloesbarem Remote nicht ab und zeigt die neue Meldung nicht - der Ref-Teil wird still uebersprungen (internal/cli/refs.go:75 fileOnRefOnly), das Ticket bleibt als Datei liegen. Das ist bestehendes Verhalten und nicht Teil dieses Tickets. Die Meldung kommt bei pull, fetch, release, snapshot.

Naechster Schritt: nichts zu beheben. Ticket kann in die naechste Lane.
- **2026-09-14 13:21 · Alexander Sacharov** — review-Lane: keine Ruecksendung. Der Diff selbst wurde unabhaengig nachgeprueft (go vet ./... sauber, go test core/settings, core/gitref, internal/cli gruen) und die beiden Kernfaelle von Hand am gebauten Binary nachgestellt: Repo mit nur origin und settings.json=upstream laeuft durch (Ref landet in origin.git, pull/claim/release Exit 0), und mit git config jaira.remote ghost bricht pull mit Exit 1 und der dreiteiligen Meldung ab. Zwei Beobachtungen stehen jetzt in review-gaps - wichtigste: internal/cli/boardremote_test.go TestTheRefGoesToTheConfiguredRemoteWhenTheRepositoryHasIt haengt origin und upstream an dasselbe bare-Repo und kann deshalb gar nicht bemerken, wenn der Ref im Fork landete; er prueft nur den aufgeloesten Namen. Abgedeckt ist die Aussage anderswo, deshalb kein Blocker - aber wer den Test spaeter anfasst, sollte upstream auf ein zweites bare-Repo zeigen lassen. Hinweis fuers Archiv: review-gaps trug vorher den ausfuehrlichen Optimize-Bericht; der steht in der Ticket-Historie (Commit 154b198 ff.) und in der Notiz vom 13:12.
