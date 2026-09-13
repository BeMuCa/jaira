---
id: 01M2E248SM9X1JRZBNTHC9V7ZV
title: "Rollen-Prompts im Binary ausliefern: jaira roles install"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "jaira liefert die Rollen-Prompts selbst aus: 'jaira roles install' schreibt sie nach .claude/skills/ (--project) oder ~/.claude/skills/ (--global), ohne lokale Aenderungen zu ueberschreiben"
context: |-
  Heute liegen die Rollen-Prompts nur in ~/.claude/skills auf genau einem Rechner. Ein Teamkollege klont das Repo, startet jaira und hat die Lanes, aber niemanden, der sie faehrt.
  Die Rollen sind: teamlead (redet mit dem Menschen, priorisiert), dispatcher (faehrt ein Ticket Lane fuer Lane, delegiert), role-lane (genau eine Lane), role-tester (Tests, ohne zu reparieren), role-pr (PR oeffnen, nie mergen), role-brainstorm, role-research.
  Der Mechanismus existiert schon zweimal im Code und wird nicht neu erfunden: core/lane/lane.go:34 bettet builtin/*.md per go:embed ein; core/board/announce.go:295 schreibt einen verwalteten Block in CLAUDE.md/AGENTS.md und schuetzt Handarbeit mit dem jaira:local-Marker.
  Drei Entscheidungen sind schon gefallen:
  - Flacher Ordnername mit Praefix, nicht verschachtelt: .claude/skills/jaira-teamlead/SKILL.md. Verschachtelte Ordner sind in Claude Code kein Namensraum, nur Plugins haben einen - eine Rolle unter skills/jaira/teamlead/ wuerde stillschweigend nicht geladen.
  - Zielordner wie in announce.go: in die Skill-Ordner schreiben, die es gibt (.claude/, .codex/, .agents/), und wenn es keinen gibt, .claude/ anlegen.
  - Eine Datei, die vom eingebetteten Original abweicht, wird ohne --force nicht ueberschrieben. Wie bei 'lanes adopt'. Eine selbst angepasste Rolle ueberlebt 'jaira update'.
  Die sieben Prompt-Dateien sind bereits geschrieben und liegen unter ~/.claude/skills/ - sie sind der Inhalt von core/role/builtin/, nicht neu zu erfinden.
  Nicht Teil dieses Tickets: ein roles-Marktplatz analog zu lanes/ auf GitHub. Erst wenn jemand eine eigene Rolle teilen will.
definition-of-done: "'jaira roles install --project' legt sieben Ordner .claude/skills/jaira-<id>/SKILL.md an; ein zweiter Lauf aendert nichts; eine von Hand geaenderte Datei bleibt ohne --force unberuehrt und wird gemeldet; 'jaira roles install --global' schreibt nach ~/.claude/skills; 'jaira roles list' nennt die eingebetteten Rollen; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen"
tags:
  - cli
blocked-by: []
related: []
commits:
  - pending
created-at: 2026-09-13T18:57:58Z
updated-at: 2026-09-13T20:49:29Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-34867
claimed-at: 2026-09-13T20:41:45Z
outcome-what: "Moved core/role/builtin/jaira-teamlead/scripts/ to core/role/builtin/jaira-dispatcher/scripts/ and renamed TestTeamleadShipsItsScript to TestDispatcherShipsItsScript, with the stat path and two source comments pulled along."
outcome-why: "jaira-dispatcher/SKILL.md:95 is the only prompt that names scripts/spawn.sh, and that path is relative to its own skill folder, so after 'roles install' the reference pointed at nothing while jaira-teamlead carried a script its own prompt never mentions."
outcome-resolves: "The shipped dispatcher role now finds scripts/spawn.sh where its prompt says it is; go test ./... -race green."
review-summary: |-
  core/role/builtin/jaira-dispatcher/scripts/spawn.sh:20 bricht in jedem Repo ohne .env hart ab. Zeile 5 setzt 'set -euo pipefail', Zeile 20 macht 'cp "$root/.env" "$wt/.env"' unbedingt - jaira selbst hat kein .env (ls im Repo-Root: kein .env, kein compose-File), also beendet sich das Skript mit Exit 1, bevor ueberhaupt eine Pane entsteht. Die Runden 1-3 haben nur geprueft, WO das Skript liegt, nie WAS darin steht; ab diesem Ticket wird der Inhalt mit dem Binary verteilt und ist auf jedem Rechner eines Teamkollegen die Wahrheit. Der projektspezifische Block (Zeilen 14-16 Port-Offset und 20-30 .env-Kopie mit COMPOSE_PROJECT_NAME/HTTP_PORT/DB_PORT_HOST/VITE_PORT_HOST/BACKEND_PORT_HOST) gehoert hinter ein 'if [ -f "$root/.env" ]', damit der generische Teil - Worktree anlegen, Pane splitten, claude starten, Prompt tippen - in einem Repo ohne Container-Stack durchlaeuft. jaira-dispatcher/SKILL.md:94 sagt selbst 'Auf einem Projekt mit einem Container-Stack', das Skript setzt einen aber voraus.
  internal/cli/roles_test.go:217 deklariert 'Unprefixed []string `json:"unprefixed"`' - ein Rest der in Runde 1 entfernten Twins(). reportRoleInstall in internal/cli/roles.go:121-129 gibt kein Feld 'unprefixed' mehr aus, und kein Test liest payload.Unprefixed. Die Zeile ersatzlos streichen.
---

# Rollen-Prompts im Binary ausliefern: jaira roles install

## Definition of Done

- [x] 'jaira roles install --project' legt sieben Ordner .claude/skills/jaira-<id>/SKILL.md an; ein zweiter Lauf aendert nichts; eine von Hand geaenderte Datei bleibt ohne --force unberuehrt und wird gemeldet; 'jaira roles install --global' schreibt nach ~/.claude/skills; 'jaira roles list' nennt die eingebetteten Rollen; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh; TestDispatcherShipsItsScript (core/role/role_test.go:56); go test ./... -race green

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] die sieben Prompts einfrieren: ~/.claude/skills/<id>/ nach core/role/builtin/jaira-<id>/ kopieren, im Frontmatter name: auf jaira-<id> setzen und jede Querverweis-Zeile (/role-lane, /role-tester, /role-research) auf den praefixierten Namen umschreiben
- [x] core/role/role.go: //go:embed all:builtin, Typ Role{ID,Name,Description,Files}, Builtins() und Get(id) lesen name:/description: aus SKILL.md
- [x] core/role/target.go: Zielordner bestimmen - im Projekt die vorhandenen von .claude/ .codex/ .agents/, keiner da -> .claude/ anlegen; global immer ~/.claude/skills
- [x] core/role/install.go: Install(dstSkillsDir, force) schreibt je Datei und meldet written | unchanged | modified-skipped | overwritten; Vergleich gegen die eingebetteten Bytes, nicht blosses Stat wie lane.Export
- [x] Tests in core/role: Erstinstallation, zweiter Lauf komplett unchanged, handgeaenderte Datei ohne --force uebersprungen und gemeldet, mit --force ueberschrieben, Pfad kann dstDir nicht verlassen
- [x] internal/cli/roles.go: 'jaira roles list' und 'jaira roles install --project|--global [--force]', --json-Form, Exit 3 wenn eine geaenderte Datei uebersprungen wurde; in root.go registrieren
- [x] CLI-Tests: Textausgabe, --json, Exit-Codes
- [x] eine Zeile in core/release/NOTES.md unter ## Unreleased
- [x] go test ./... -race gruen

## Progress
- **2026-09-13 19:01 · Alexander Sacharov** — Recherche 2026-09-13, code.claude.com/docs/en/skills (Primaerquelle, HIGH):
- SKILL.md muss genau eine Ebene tief liegen: ~/.claude/skills/<name>/SKILL.md bzw. .claude/skills/<name>/SKILL.md. Verschachtelte Ordner sind kein Namensraum. skills/jaira/teamlead/SKILL.md wuerde nicht geladen - die Praefix-Entscheidung im Kontext ist damit bestaetigt.
- Neu und wichtiger als erwartet: der Kommandoname kommt vom ORDNERNAMEN, nicht vom frontmatter 'name'. Bei einem Projekt- oder Personal-Skill ist 'name' nur ein Anzeigelabel. Nur Plugin-Skills bekommen ihren Namen aus dem Feld.
- Folge fuer dieses Ticket: der Ordnername ist die API. Ordner jaira-teamlead -> Kommando /jaira-teamlead. Das frontmatter 'name' muss denselben String tragen, sonst zeigt die Liste etwas anderes an, als der Mensch tippen muss.
- Unterordner INNERHALB einer Rolle sind erlaubt (reference.md, scripts/). Falls eine Rolle spaeter Hilfsdateien braucht, ist das der Weg, nicht ein zweiter Skill.
- Offener Punkt, nicht recherchiert: auf diesem Rechner liegen die sieben Rollen bereits unter ~/.claude/skills/teamlead, dispatcher, role-* ohne Praefix. Nach 'roles install --global' stehen sie doppelt da. Der Installer muss das erkennen und melden, statt stillschweigend ein zweites Paar anzulegen.
- **2026-09-13 19:43 · Alexander Sacharov** — Board-Umzug am 2026-09-13: die Ticket-Refs liegen jetzt auf upstream (BeMuCa/jaira), nicht mehr nur im Fork. 24 Refs und der Snapshot-Zweig jaira/board sind hinueber gepusht, die Kopien im Fork bleiben vorerst als Backup liegen. ~/.jaira/settings.json traegt {"remote":"upstream"} - ohne diese Datei schreibt jaira wieder in den Fork, die Einstellung ist pro Rechner und wird nicht mitgeliefert.
- **2026-09-13 19:55 · Alexander Sacharov** — Der jaira--Praefix ist nicht nur ein Ordnername. In Claude Code muss das Feld name: im Frontmatter dem Ordner entsprechen, sonst laedt der Skill nicht. Also wird beim Einfrieren auch name: teamlead -> name: jaira-teamlead umgeschrieben. Damit aendert sich der Aufruf: /jaira-role-lane statt /role-lane. Die Prompts rufen sich gegenseitig auf - dispatcher/SKILL.md:38 nennt /role-lane und /role-tester, role-brainstorm/SKILL.md:30 nennt /role-research. Diese Zeilen muessen mitgezogen werden, sonst ruft der Dispatcher eines Teamkollegen einen Namen auf, den es auf seinem Rechner nicht gibt. Das ist Schritt 1 des Plans und der Grund, warum er nicht blosses Kopieren ist.
- **2026-09-13 19:55 · Alexander Sacharov** — Eine Rolle ist nicht immer genau eine Datei. teamlead hat zusaetzlich scripts/spawn.sh, die uebrigen sechs nur SKILL.md. Die Definition of Done nennt nur SKILL.md, aber nur SKILL.md auszuliefern wuerde einen kaputten Verweis mitliefern. Deshalb Role.Files als Liste und //go:embed all:builtin ueber den ganzen Baum - die Kosten sind null, die DoD bleibt erfuellt. Nebenbefund, nicht Teil dieses Tickets: dispatcher/SKILL.md:95 verweist auf scripts/spawn.sh, das Skript liegt aber unter teamlead/scripts/. Der Verweis geht schon heute ins Leere.
- **2026-09-13 19:55 · Alexander Sacharov** — lane.Export/copyLane (core/lane/share.go:34) kann hier nicht wiederverwendet werden, obwohl der Mechanismus gleich aussieht. Es prueft nur os.Stat: Datei da -> Abbruch. Die DoD verlangt drei Faelle statt zwei: zweiter Lauf aendert nichts (Bytes gleich -> unchanged, kein Fehler), handgeaenderte Datei bleibt unberuehrt und wird gemeldet (Bytes verschieden -> skipped), --force ueberschreibt. Das ist ein Byte-Vergleich gegen das eingebettete Original, kein Stat. Deshalb eigenes core/role/install.go und kein Umbau von lane.Export - lane haengt an 'lanes use' und soll sich nicht mitaendern.
- **2026-09-13 19:58 · Alexander Sacharov** — Beim Einfrieren umgeschrieben, ueber blosses Kopieren hinaus: name: in allen sieben Frontmattern auf jaira-<id>, und jeder Querverweis auf einen Rollennamen mit /-Praefix. Betroffen waren mehr Stellen als die Notiz vom 19:55 nannte: zusaetzlich zu dispatcher:38 und role-brainstorm:30 auch fuenf description:-Zeilen ('Invoked as /role-lane ...') und teamlead/scripts/spawn.sh:47, das den Prompt per herdr in eine Pane tippt. Das Regex laesst /-Namen in Pfaden in Ruhe ((?<![\w/.-])), sonst waere ~/.claude/skills/role-lane mitgewandert. Die Quellen unter ~/.claude/skills bleiben unveraendert - sie sind ab jetzt nicht mehr die Wahrheit, core/role/builtin ist es.
- **2026-09-13 20:06 · Alexander Sacharov** — Drei Dinge, die der Code nicht sagt:
- go:embed all:builtin statt builtin/* - ohne all: laesst go:embed Unterverzeichnisse aus, teamlead/scripts/spawn.sh waere still verschwunden und der Prompt haette auf eine Datei gezeigt, die nicht mitkommt.
- Das Ausfuehrbar-Bit ueberlebt embed.FS nicht: alles kommt als 0644 wieder heraus. install.go entscheidet deshalb an der Endung (.sh -> 0755). Ein spawn.sh ohne x-Bit ist eine kaputte Rolle, und nichts haette es gemeldet.
- Der offene Punkt aus der Recherche ist erledigt, aber anders als 'stillschweigend ein zweites Paar anlegen': das Praefix macht jaira-teamlead und teamlead zu zwei verschiedenen Kommandos, es kollidiert also nichts. Gefaehrlich ist nur, dass jemand weiter den alten Namen tippt. role.Twins() meldet die unpraefixierte Kopie und loescht sie nie - ein Verzeichnis, das dieses Werkzeug nicht geschrieben hat, entfernt es auch nicht.
Nicht gemacht, bewusst: 'jaira update' ruft roles install nicht auf. Eine neue Rolle nach einem Upgrade muss man selbst holen. Das gehoert in ein eigenes Ticket, sobald jemand es vermisst.
- **2026-09-13 20:09 · Alexander Sacharov** — critique: sechs Befunde, alle in review-summary mit Datei und Gegenvorschlag. Der gewichtigste ist das Fan-out in core/role/target.go: dieselben sieben Rollen landen in jedem existierenden Agent-Ordner, ein Harness liest mehrere davon, also ist jaira-teamlead danach doppelt registriert. Die DoD nennt genau ein Ziel - das ist kein Geschmacksurteil, sondern eine Abweichung von dem, was das Ticket verlangt. Die uebrigen fuenf sind Ballast: installDirs/skillsDirOf rekonstruiert Zielpfade, die der Aufrufer schon hat; Role.Name ist ein zweiter Name fuer Role.ID; within() prueft einen zur Compile-Zeit unmoeglichen Zustand; Twins() ist eine Einmal-Migration fuer genau einen Rechner und traegt dafuer ein dauerhaftes JSON-Feld. Nicht beanstandet und bewusst stehen gelassen: der Byte-Vergleich statt os.Stat wie lane.Export - drei Zustaende brauchen drei Antworten, die Begruendung im Commit traegt. Nebenbei, kein eigener Befund: der Doc-Kommentar ueber SkippedAny in core/role/install.go:152 beginnt mit 'Skipped reports', also mit dem falschen Namen.
- **2026-09-13 20:16 · Alexander Sacharov** — critique abgearbeitet, alle sechs Befunde. Drei Dinge, die der Diff nicht selbst sagt:
- Statt Fan-out jetzt genau ein Ziel plus --into <dir>. Das Flag ist nicht Zierde: ein Projekt, dessen Agent unter .codex/skills liest, hatte vorher den Fan-out als einzigen Weg. Ohne --into waere der Befund nicht behebbar gewesen, ohne diese Nutzer auszusperren.
- Role.Name ist weg, aber die Invariante dahinter nicht: Claude Code laedt einen Skill nicht, wenn frontmatter name: vom Ordner abweicht. Da nichts im Paket das Feld mehr liest, wuerde ein Tippfehler in einer eingebetteten SKILL.md still durchgehen. core/role/role_test.go:39 prueft deshalb jetzt direkt die Bytes auf 'name: jaira-<id>'. Wer diesen Test fuer redundant haelt und loescht, macht genau das Loch auf.
- Twins() gestrichen wie verlangt, aber die Sache selbst bleibt wahr: wer die Prompts frueher von Hand unter den blanken Namen liegen hatte, hat sie weiter, und sie antworten auf ein anderes Kommando. Das steht jetzt als zweite Zeile in core/release/NOTES.md, weil ein Mensch es dort einmal liest, statt dass Code es bei jedem Lauf neu herausfindet.
Nicht angefasst, bewusst: der Byte-Vergleich statt os.Stat. Die critique hat ihn ausdruecklich stehen lassen.
- **2026-09-13 20:16 · Alexander Sacharov** — Commit dieser Runde: 7e889d4.
- **2026-09-13 20:19 · Alexander Sacharov** — critique, zweite Runde: drei Befunde, alle in review-summary mit Datei und Gegenvorschlag. Keiner davon ist ein Wiederaufwaermen der ersten Runde - das Fan-out, Role.Name, within(), Twins() und installDirs sind erledigt und werden nicht erneut angefasst.
- Der gewichtigste ist der handgeschriebene Frontmatter-Scanner in core/role/role.go:136. Der Kommentar dort begruendet ihn damit, dass eine YAML-Abhaengigkeit fuer einen String nichts einbringt - das stimmt, trifft aber nicht zu: ticket.ParseDoc liegt schon im Modul und core/lane/lane.go:243 liest die Lane-Beschreibung genau so. Es geht nicht um eine Abhaengigkeit, sondern um zwei Parser fuer dasselbe Dateiformat.
- Get() ist keine Geschmacksfrage: ausser den eigenen Tests ruft es niemand.
- firstSentence doppelt und mit abweichender Regel ist der kleinste der drei, aber der einzige mit sichtbarer Folge - jede Beschreibung mit einem Pfad darin wird in der Liste falsch gekuerzt.
Bewusst nicht beanstandet: --into (die erste Runde hat das Flag selbst vorgeschlagen), der Byte-Vergleich statt os.Stat, und der Test in role_test.go:39, der name: gegen den Ordnernamen prueft - der haelt eine Invariante, die sonst niemand mehr haelt.
- **2026-09-13 20:24 · Alexander Sacharov** — Zweite critique abgearbeitet, alle drei Befunde. Was der Diff nicht selbst sagt:
- Der Handscanner ist durch ticket.ParseDoc ersetzt, und das ist kein reiner Umbau: ParseDoc liest echtes YAML, also kommt eine description mit escaptem Anfuehrungszeichen jetzt richtig heraus statt mit Backslash. Kein eingebetteter Prompt hat heute so eine Zeile - der Unterschied faellt erst auf, wenn jemand eine schreibt.
- firstSentence lebt nur noch einmal, als board.FirstSentence, und schneidet an '. '. Die Rollenliste zeigt damit 'Talk to the human about a board: what to work next, why, and what needs a decision.' ganz, statt nach dem ersten Punkt in einem Pfad abzubrechen.
- Get() ist weg; TestTeamleadShipsItsScript sucht die Rolle jetzt aus Builtins() heraus. Damit gibt es keine exportierte Funktion mehr, die nur ihr eigener Test benutzt - wer sie wieder einfuehrt, braucht einen Aufrufer ausserhalb der Tests.
- **2026-09-13 20:24 · Alexander Sacharov** — Commit dieser Runde: 8e6398b.
- **2026-09-13 20:27 · Alexander Sacharov** — critique, dritte Runde: ein Befund, in review-summary mit Datei und Gegenvorschlag. Das Skript spawn.sh liegt unter jaira-teamlead, genannt wird es nur von jaira-dispatcher/SKILL.md:95 - und der Pfad dort ist relativ zum eigenen Skill-Ordner, also zeigt er nach der Installation auf nichts. Die Notiz vom 19:55 hatte das als Nebenbefund abgetan ('Der Verweis geht schon heute ins Leere'), aber ab diesem Ticket liefert jaira die Datei selbst aus: was hier eingebettet wird, ist die Wahrheit auf jedem Rechner eines Teamkollegen, nicht mehr eine gewachsene Kopie in ~/.claude/skills. Ein kaputter Verweis wird damit von einem lokalen Schoenheitsfehler zu etwas, das mit dem Binary verteilt wird. Deshalb jetzt ein Befund und nicht wieder ein Nebensatz.

Bewusst nicht beanstandet, nichts davon wird erneut angefasst: der Byte-Vergleich statt os.Stat wie lane.Export (Runde 1 hat ihn ausdruecklich stehen lassen), --into, ticket.ParseDoc in frontmatterDescription, board.FirstSentence in internal/cli (Runde 2 hat beides selbst vorgeschlagen), und role_test.go:39, das name: gegen den Ordnernamen prueft.

Geprueft und in Ordnung befunden: kein exportierter Bezeichner in core/role hat nur noch Tests als Aufrufer - File, Builtins, Install, SkippedAny, ProjectTarget und GlobalTarget werden alle aus internal/cli heraus benutzt. Der Usage-Zweig in internal/cli/roles.go:82-89 ist handgeschrieben statt cobra.MarkFlagsMutuallyExclusive, folgt damit aber dem, was lanes.go:191 und checklist.go:106 schon tun - das ist das vorhandene Muster, kein neues.
- **2026-09-13 20:27 · Alexander Sacharov** — Commit dieser Runde: dd5e7bc.
- **2026-09-13 20:28 · Alexander Sacharov** — Dispatcher haelt an: die Critique-Lane hat das Ticket dreimal zurueckgeschickt (Regel der dispatcher-Rolle). Die Regel unterstellt eine falsche Definition of Done - hier trifft das nicht zu: die Befunde werden kleiner, 6 -> 3 -> 1, und jeder war ein anderer. Runde 1 (7e889d4): Fan-out in drei Agent-Ordner auf .claude/skills reduziert, --into ergaenzt. Runde 2 (8e6398b): Handgeschriebener Frontmatter-Scanner durch ticket.ParseDoc ersetzt. Runde 3, offen: core/role/builtin/jaira-teamlead/scripts/spawn.sh liegt in der falschen Rolle, nur jaira-dispatcher/SKILL.md:95 nennt das Skript - Verzeichnis nach core/role/builtin/jaira-dispatcher/ verschieben und core/role/role_test.go:55-57 sowie :128 mitziehen. Das ist eine kleine, klar beschriebene Aenderung. Ein Mensch entscheidet, ob noch eine Runde laeuft oder ob das Ticket so weitergeht.
- **2026-09-13 20:43 · Alexander Sacharov** — Dritte critique abgearbeitet, der eine Befund. Was der Diff nicht selbst sagt:
- Verschoben, nicht kopiert: das Skript gibt es genau einmal, und genannt wird es genau einmal (jaira-dispatcher/SKILL.md:95). Solange es keinen zweiten Nutzer gibt, wird auch keine Mechanik fuer Mehrfachnutzung gebaut - kein Sammelordner, kein Aufloesen von Pfaden ueber Rollengrenzen hinweg. Wer spaeter ein Skript aus zwei Rollen heraus braucht, entscheidet dann.
- Der Testname hat mitgewandert: TestTeamleadShipsItsScript -> TestDispatcherShipsItsScript. Der alte Kommentar ('teamlead references a script') war schon vor diesem Ticket falsch und haette den Befund verdeckt, wenn ihn jemand als Beleg gelesen haette.
- Keine Zeile in core/release/NOTES.md: das roles-Feature steht komplett unter ## Unreleased, es hat also nie ein Binary gegeben, das spawn.sh unter jaira-teamlead ausgeliefert haette. Niemandem ist etwas zu erzaehlen, was er nie gesehen hat.
- **2026-09-13 20:43 · Alexander Sacharov** — Commit dieser Runde: 8561b2c.
- **2026-09-13 20:46 · Alexander Sacharov** — critique, vierte Runde: zwei Befunde, beide mit Datei und Gegenvorschlag in review-summary.

- Der gewichtigere ist neu und kein Wiederaufwaermen von Runde 3: die hat entschieden, in WELCHER Rolle spawn.sh liegt, nicht was darin steht. Nachgesehen habe ich es erst jetzt. 'set -euo pipefail' (Zeile 5) plus unbedingtes 'cp $root/.env' (Zeile 20) heisst: in jedem Repo ohne .env endet das Skript mit Exit 1, bevor eine Pane entsteht - auch in jaira selbst, das weder .env noch compose-File hat. Der Block gehoert hinter ein '[ -f $root/.env ]', der generische Teil (Worktree, Pane, claude, Prompt) laeuft dann ueberall.
- Der zweite ist Aufraeumen: das JSON-Feld 'unprefixed' im Test-Struct internal/cli/roles_test.go:217 ist ein Rest der in Runde 1 gestrichenen Twins(). Nichts gibt es mehr aus, nichts liest es.

Bewusst nicht beanstandet, nichts davon wird erneut angefasst: der Byte-Vergleich statt os.Stat (Runde 1), --into (Runde 1), ticket.ParseDoc in frontmatterDescription und board.FirstSentence (Runde 2), der Ort von spawn.sh unter jaira-dispatcher (Runde 3), role_test.go:39 mit der name:-Invariante, und der handgeschriebene Usage-Zweig in internal/cli/roles.go:82-89, der dem Muster aus lanes.go:191 folgt.

Nicht als Befund gefuehrt, weil es ein eigenes Ticket waere: das Skript setzt herdr und python3 voraus, beides steht nirgends als Anforderung. Solange der Prompt es nur als Hilfe fuer einen Container-Stack nennt, traegt das - wer kein herdr hat, bekommt Zeile 10 als klare Meldung, kein stilles Scheitern.
- **2026-09-13 20:47 · Alexander Sacharov** — Commit dieser Runde: 82d7ab7.
- **2026-09-13 20:49 · Alexander Sacharov** — Vierte critique abgearbeitet, beide Befunde. Was der Diff nicht selbst sagt:
- Der .env-Block in spawn.sh ist nicht nur eingerueckt, sondern auch der Port-Offset (off=...) ist mit hineingewandert. Er wird ausserhalb des Blocks von nichts mehr gelesen, und ein cksum auf jedem Lauf in einem Repo ohne Container-Stack zu berechnen waere Arbeit fuer eine Zahl, die niemand benutzt.
- Bewusst NICHT generisch gemacht: die Variablennamen (COMPOSE_PROJECT_NAME, VITE_PORT_HOST, ...) und der Worktree-Name rg-$slug stammen aus genau einem Projekt. Sie bleiben stehen, weil ein Repo ohne .env den Block jetzt gar nicht mehr betritt - ein konfigurierbares Port-Schema waere Mechanik fuer einen Nutzer, den es nicht gibt. Wer ein zweites Projekt mit Stack anschliesst, entscheidet dann.
- Keine Zeile in core/release/NOTES.md: das roles-Feature steht komplett unter ## Unreleased, spawn.sh war nie in einem Binary. Gleiche Begruendung wie in Runde 3.
