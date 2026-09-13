---
id: 01M2E248SM9X1JRZBNTHC9V7ZV
title: "Rollen-Prompts im Binary ausliefern: jaira roles install"
status: critique
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
commits: []
created-at: 2026-09-13T18:57:58Z
updated-at: 2026-09-13T20:09:43Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-45975
claimed-at: 2026-09-13T20:07:20Z
outcome-what: "core/role: seven role prompts embedded via go:embed all:builtin, plus target resolution and a three-way installer (written/unchanged/skipped/overwritten). internal/cli/roles.go: 'jaira roles list' and 'jaira roles install --project|--global [--force]', --json on both, exit 3 when an edited file is left alone. Prompts frozen from ~/.claude/skills with name: and every cross-reference rewritten to the jaira- prefix."
outcome-why: "The role prompts lived in one person's ~/.claude/skills. A teammate who cloned got the board and nobody to drive it. Shipping them inside the binary makes them travel with the tool, the way the lanes already do."
outcome-resolves: "jaira roles install --project writes .claude/skills/jaira-<id>/SKILL.md for all seven roles, --global writes ~/.claude/skills, a second run reports only unchanged and exits 0, an edited file is left alone, reported and exits 3, --force replaces it, jaira roles list names them, core/release/NOTES.md carries a line under ## Unreleased, and go test ./... -race is green."
review-summary: |-
  core/role/target.go:20 ProjectTargets liefert bis zu drei Zielverzeichnisse (.claude, .codex, .agents), und internal/cli/roles.go:110 installiert in jedes davon dieselben sieben Rollen. Ein Harness liest Skills aus mehreren dieser Ordner, also steht jaira-teamlead dann doppelt registriert da, und das Repo traegt dreimal dieselben Prompts. Die DoD nennt genau ein Ziel. Auf ein Verzeichnis reduzieren: .claude/skills, und wenn ein anderer Ordner gewuenscht ist, per Flag statt per Fan-out.
  core/role/target.go:8-17 begruendet die Auswahl mit core/board/announce.go:295, macht aber das Gegenteil: announce.go schreibt AGENTS.md UND CLAUDE.md immer, target.go schreibt nur in existierende Ordner. Entweder die Regel uebernehmen oder den Verweis streichen und schreiben, warum ein Verzeichnisbaum anders behandelt wird als eine Markdown-Datei.
  internal/cli/roles.go:171-199 installDirs/skillsDirOf rekonstruiert die Skills-Verzeichnisse aus den geschriebenen Pfaden zurueck, obwohl RunE sie als targets bereits in der Hand hat. targets an reportRoleInstall durchreichen und beide Helfer loeschen.
  core/role/role.go:40 Role.Name und der name:-Zweig in frontmatter() liefern nie etwas anderes als Role.ID - der Paketkommentar sagt selbst, dass das Verzeichnis die API ist, und internal/cli/roles_test.go:82 pinnt Name == ID fuer alle sieben Rollen fest. Feld, Parse-Zweig und das name-Feld in --json entfernen; description bleibt.
  core/role/install.go:57-60 der within()-Check prueft einen Zustand, den sein eigener Kommentar als unmoeglich beschreibt: Rollen-ID und relativer Pfad kommen aus dem eingebetteten FS und stehen zur Compile-Zeit fest. Check, Helfer und core/role/role_test.go:208-213 streichen.
  core/role/install.go:129 Twins() steht in Ziel und DoD nicht und loest eine einmalige Migration fuer genau die Rechner, auf denen die Prompts vorher von Hand lagen - also den einen. Es kostet jeden spaeteren Leser Code in install.go, eine Zeile CLI-Ausgabe und das dauerhafte JSON-Feld unprefixed. Streichen und die Umbenennung stattdessen als Zeile in core/release/NOTES.md erwaehnen.
---

# Rollen-Prompts im Binary ausliefern: jaira roles install

## Definition of Done

- [x] 'jaira roles install --project' legt sieben Ordner .claude/skills/jaira-<id>/SKILL.md an; ein zweiter Lauf aendert nichts; eine von Hand geaenderte Datei bleibt ohne --force unberuehrt und wird gemeldet; 'jaira roles install --global' schreibt nach ~/.claude/skills; 'jaira roles list' nennt die eingebetteten Rollen; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen
  proof: core/role/install.go:52 Install(); TestRolesInstallProjectWritesSevenRoles, TestRolesInstallGlobalWritesUnderHome, TestRolesInstallSecondRunExitsZero, TestRolesInstallLeavesAnEditedFileAloneAndExitsThree, TestRolesListNamesEveryBuiltin (internal/cli/roles_test.go); core/release/NOTES.md:18; go test ./... -race green

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
