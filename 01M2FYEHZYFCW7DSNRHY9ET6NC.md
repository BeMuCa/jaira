---
id: 01M2FYEHZYFCW7DSNRHY9ET6NC
title: "Der Remote-Name gilt pro Rechner, gebraucht wird er pro Board"
status: in-progress
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
updated-at: 2026-09-14T12:52:57Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-93721
claimed-at: 2026-09-14T12:36:43Z
---

# Der Remote-Name gilt pro Rechner, gebraucht wird er pro Board

## Definition of Done

- [ ] In einem Repository, dessen einziger Remote origin heisst, laeuft 'jaira release <id>' durch, waehrend ~/.jaira/settings.json weiterhin remote: upstream sagt. Nachgestellt an einem Fixture-Repository mit genau einem Remote.

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
- [ ] README.md:338 (settings.json-Block) um den Board-Remote und die Reihenfolge ergaenzen
- [ ] core/release/NOTES.md: Abschnitt '## Unreleased' oben anlegen und eine Zeile schreiben
- [ ] go test ./... und Gegenprobe im jaira-Repo selbst: der Ticket-Ref geht weiterhin nach upstream

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
