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
commits: []
created-at: 2026-09-13T18:57:58Z
updated-at: 2026-09-13T19:56:51Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-6599
claimed-at: 2026-09-13T19:56:51Z
---

# Rollen-Prompts im Binary ausliefern: jaira roles install

## Definition of Done

- [ ] 'jaira roles install --project' legt sieben Ordner .claude/skills/jaira-<id>/SKILL.md an; ein zweiter Lauf aendert nichts; eine von Hand geaenderte Datei bleibt ohne --force unberuehrt und wird gemeldet; 'jaira roles install --global' schreibt nach ~/.claude/skills; 'jaira roles list' nennt die eingebetteten Rollen; eine Zeile in core/release/NOTES.md unter ## Unreleased; go test ./... -race gruen

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [ ] die sieben Prompts einfrieren: ~/.claude/skills/<id>/ nach core/role/builtin/jaira-<id>/ kopieren, im Frontmatter name: auf jaira-<id> setzen und jede Querverweis-Zeile (/role-lane, /role-tester, /role-research) auf den praefixierten Namen umschreiben
- [ ] core/role/role.go: //go:embed all:builtin, Typ Role{ID,Name,Description,Files}, Builtins() und Get(id) lesen name:/description: aus SKILL.md
- [ ] core/role/target.go: Zielordner bestimmen - im Projekt die vorhandenen von .claude/ .codex/ .agents/, keiner da -> .claude/ anlegen; global immer ~/.claude/skills
- [ ] core/role/install.go: Install(dstSkillsDir, force) schreibt je Datei und meldet written | unchanged | modified-skipped | overwritten; Vergleich gegen die eingebetteten Bytes, nicht blosses Stat wie lane.Export
- [ ] Tests in core/role: Erstinstallation, zweiter Lauf komplett unchanged, handgeaenderte Datei ohne --force uebersprungen und gemeldet, mit --force ueberschrieben, Pfad kann dstDir nicht verlassen
- [ ] internal/cli/roles.go: 'jaira roles list' und 'jaira roles install --project|--global [--force]', --json-Form, Exit 3 wenn eine geaenderte Datei uebersprungen wurde; in root.go registrieren
- [ ] CLI-Tests: Textausgabe, --json, Exit-Codes
- [ ] eine Zeile in core/release/NOTES.md unter ## Unreleased
- [ ] go test ./... -race gruen

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
