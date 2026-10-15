---
id: 01M2HX6SYCFB27V26R5APWAF22
title: "jaira update bringt das ganze Setup nach, nicht nur den Textblock"
status: pre-process
ready: true
creator: Alexander Sacharov
goal: "Nach einem Upgrade des Binaries genuegt ein Befehl, um ein Board nachzuziehen: Agent-Block, Rollen und Lanes - und er sagt vorher, was er anfassen wuerde."
definition-of-done: "'jaira update' bringt mit den entsprechenden Schaltern Agent-Block, Rollen und Lanes in einem Aufruf auf den Stand des installierten Binaries. Nachgestellt an einem Board, dessen Block, Rollen und Lanes alle drei veraltet sind."
tags:
  - cli
blocked-by: []
related:
  - 01M21PJTWPK4DZ89PRY6T8A8KM
commits: []
created-at: 2026-09-15T06:48:56Z
updated-at: 2026-10-15T12:46:57Z
assignee: "BeMuCa"
updated-by: BeMuCa
context: |-
  Nach 'jaira self update' haengt jedes Board zurueck, und es haengt an drei Stellen zurueck: der erzeugte Agent-Block in CLAUDE.md und AGENTS.md, die Rollen-Prompts, und die Lanes.

  Heute braucht das drei Befehle - 'jaira update', 'jaira roles install --global --force', und je Lane ein 'jaira lanes use <id> --force'. Der Sammelschalter fuer Lanes ist Ticket 0S1A7M; dieses hier ist die Frage, ob ein Befehl alle drei Dinge tut.

  'jaira update' traegt den passenden Namen und verspricht in seiner eigenen Hilfe genau das: 'brings a board's own setup up to date with whichever binary you already have installed'. Es regeneriert aber nur den Textblock. Rollen und Lanes sind ebenso 'a board's own setup'.

  Die Entscheidung, die dieses Ticket treffen muss und die groesser ist als die Umsetzung: was tut der nackte Aufruf ohne Schalter? Wer 'jaira update' heute gewohnheitsmaessig laeuft, erwartet eine Textblock-Regeneration. Faengt der Befehl an, Rollen und Lanes anzufassen, aendert sich das Verhalten unter den Fuessen aller, die ihn schon benutzen. Das ist eine Frage der Vertraeglichkeit, keine der Bequemlichkeit.

  Die zweite Haelfte des Problems, von Alex am 2026-09-15 benannt: auch mit einem Sammelbefehl muss er ihn in JEDEM Projekt einzeln aufrufen. 'jaira projects' kennt die Boards dieses Rechners bereits - acht Eintraege, darunter jaira, requirementsgenie und die Worktrees.

  Aber dieselbe Liste zeigt, warum ein rechnerweiter Lauf heute gefaehrlich waere: drei der acht Eintraege sind Wegwerf-Boards, die testing-Worker im Laufe des Tages unter /tmp/claude-*/scratchpad/ angelegt haben - sandbox, board74, board. Sie stehen dauerhaft in der Liste, obwohl ihre Verzeichnisse Muell sind. Ein Befehl, der ueber alle Boards laeuft, laeuft zuerst in diese hinein.

  Dazu kommt seit Commit 42d8253: die Lane-Dateien dieses Boards liegen unter git. Ein rechnerweiter Lauf macht damit Arbeitsbaeume in mehreren Repositories schmutzig, nicht nur Dateien im Heimverzeichnis.
claimed-by: EE-3NX6GL3-3760272
claimed-at: 2026-10-15T12:36:39Z
---

# jaira update bringt das ganze Setup nach, nicht nur den Textblock

## Definition of Done

- [ ] 'jaira update' bringt mit den entsprechenden Schaltern Agent-Block, Rollen und Lanes in einem Aufruf auf den Stand des installierten Binaries. Nachgestellt an einem Board, dessen Block, Rollen und Lanes alle drei veraltet sind.
- [ ] Der Aufruf ohne Schalter aendert sein Verhalten nicht still: was er heute tut, tut er weiterhin, und alles darueber hinaus wird ausdruecklich verlangt. Die Entscheidung steht begruendet im Ticket, nicht nur im Code.
- [ ] Der Befehl sagt vor dem Schreiben, was er anfassen wuerde, und was er auslaesst weil ein Mensch es geaendert hat - wie 'jaira roles install' es schon tut.
- [ ] Eine Zeile je Aenderung in core/release/notes/PWAF22.md, als Anweisung an den Leser formuliert und wahr
- [ ] jaira self upgrade fuehrt danach jaira update und die Rollen-Installation an den Orten aus, an denen Rollen installiert sind (global und/oder im Projekt), und sagt, was es angefasst hat; Test
- [ ] Sind installierte Skills (jaira, jaira-role-*) aelter als das laufende Binary, warnt jaira einmal und nennt den Befehl, der sie nachzieht; Test

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [ ] 1 core/shipped: manifest core/shipped/shipped.sha256 ('<sha256>  <repo path>' per version of every installed file) + shipped.go Known(path,b)/Line(path,b) + scripts/shipped-history.sh (bash, --full-history -m, maps lanes/<id>.md and .claude/skills/jaira history, refuses shallow clones, adds working-tree versions) + .gitattributes eol=lf pins; tests first: well-formed lines, Known per path+bytes, every listed path exists
- [ ] 2 core/role/install.go: writeFile recognises an older shipped version (new Action Updated, written without --force) + dryRun param on install/Agent.Install; fixture core/role/testdata (v0.4.0 jaira-role-acceptance SKILL.md); tests first: TestEveryInstalledFileIsAShippedVersion, TestInstallReplacesAnOlderShippedVersion, TestInstallDryRunWritesNothing; internal/cli/roles.go: printRoleResults/roleResultsJSON extracted, 'updated' lines + count, help text; release-note line
- [ ] 3 core/lane/refresh.go: RefreshShipped(root, dryRun) - unchanged / updated (older shipped version) / skipped (edited), only lanes the board holds and jaira ships; fixture core/lane/testdata (catalogue testing.md before 2fb8f177); tests first incl. TestEveryBuiltinLaneIsAShippedVersion
- [ ] 4 core/role/target.go: Installed(root) []Target (own project dir per agent, then personal; only dirs already holding roles; deduped) + Behind(results); tests first
- [ ] 5 core/board/announce.go: PreviewAgentFiles(root, choice, lanes) by threading dry through announceIn/announceInAgentFile/removeBlock/setClaudeImport (lines 397-576 only, not agentNote); test: preview names what the write does and writes nothing
- [ ] 6a internal/cli: 'jaira update --roles --lanes' (new updatesetup.go; update.go flags/Long/RunE 80-127 only, never 128-139): lanes refreshed before lane.Load, roles where installed, --roles works outside a board (personal only), JSON keys roles/lanes/skipped only when asked, exit 3 via named return+defer; tests first incl. DoD-1 reproduction TestUpdateRolesLanesBringsAStaleBoardUpToDate and TestUpdateWithoutSwitchesWritesOnlyTheBlock; two release-note lines
- [ ] 6b internal/cli: 'jaira update --dry-run' (would-lines, 'nothing was written (--dry-run)', no stamp, same exit code, return before the notes block) + read-only hint lines of the bare call when roles/lanes are behind; tests first: TestUpdateDryRunWritesNothing, TestUpdateDryRunExitCodeMatchesTheRealRun, TestUpdateWithoutSwitchesNamesWhatIsBehind; two release-note lines
- [ ] 7 internal/cli/self.go: after Replace run the NEW binary '<target> -C <dir> update --roles [--json]' through seam var runUpgraded; child exit 3 is not a failure; JSON setup/setup_exit; no lanes; tests first in self_test.go (recorder seam + unix test executing a shell-script payload); release-note line
- [ ] 8 internal/cli/skillsnudge.go: nudgeSkillsIfStale in openStore after nudgeIfStale - personal dirs only, once per version via JAIRA_HOME/skills-check/version, stderr only, silent on dev/update/roles install; tests first in skillsnudge_test.go (own HOME/USERPROFILE/JAIRA_HOME); release-note line
- [ ] 9 ship the jaira skill (driver decided yes, see note) - core/role/jaira/SKILL.md embedded and installed beside the roles (not a role), test keeps .claude/skills/jaira/SKILL.md byte-equal, manifest regenerated, .gitattributes, README.md:649/683 + docs/AGENTS.md:263/271; release-note line
- [ ] 10 docs: README (Upgrading, Bringing a board up to date, Roles) and docs/COMMANDS.md:130 name --roles/--lanes/--dry-run and what self upgrade now runs
- [ ] 11 reproduce DoD 1 by hand: scratch build with -ldflags version, HOME/JAIRA_HOME in the scratchpad, scratch board with stale block/roles/lane; update --dry-run, update --roles --lanes, rerun; skills warning once; transcript in a note, DoD proofs ticked

## Progress
- **2026-10-13 20:51 · BeMuCa** — Berk 2026-10-13: nach einem Upgrade sollen die jaira-Skills immer neu installiert werden und jaira update automatisch laufen (DoD 5 und 6). Anlass: am 13.10. war ~/.claude/skills/jaira-role-pr noch vom 06.10. und kannte das Ablegen vor dem PR nicht - Alex musste im PR #74 nachfordern. Weder git pull noch jaira update ziehen die Rollen-Skills heute nach (jaira roles install --help). Gehoert mit T8A8KM zusammen.
- **2026-10-15 12:34 · BeMuCa** — Entscheidung DoD 2 - der nackte Aufruf: 'jaira update' ohne Schalter schreibt weiter NUR den Agent-Block und stempelt die Version, wie heute.
- Neue Schalter: --roles (Skills dort, wo die Rollen schon liegen: Projekt-Ordner dieses Boards + persoenlicher Ordner), --lanes (Lane-Dateien dieses Boards, die eine alte ausgelieferte Version sind), --dry-run.
- Kein --all: 'jaira update --roles --lanes' ist schon ein Aufruf.
- Kein --force an update: veraltete Dateien brauchen es nicht mehr. Eine selbst geaenderte Datei ersetzt man gezielt mit 'jaira roles install --force' oder 'jaira lanes use <id> --force'.
- Ohne Schalter kommt hoechstens eine Hinweiszeile dazu ('N Dateien sind aelter - jaira update --roles'). Sie schreibt nichts. Das JSON des nackten Aufrufs bleibt gleich.
- Grund: wer update gewohnheitsmaessig laeuft, bekommt nichts Neues geschrieben. Lane-Dateien liegen auf manchen Boards unter git und gehoeren dem Team.
- **2026-10-15 12:34 · BeMuCa** — Entscheidung DoD 3 - 'vorher sagen': 'jaira update --dry-run' zeigt, was passieren wuerde, und schreibt nichts.
- Gleiche Zeilen wie der echte Lauf, nur 'would update/write' statt 'updated/written'. 'skipped ... (edited here; ...)' wie bei 'jaira roles install'. Letzte Zeile: 'nothing was written (--dry-run)'.
- Gleicher Exit-Code wie der echte Lauf (3, wenn etwas uebersprungen wuerde). So macht es schon 'jaira move --dry-run' (internal/cli/flow.go:317) und 'sync-tasks --dry-run'.
- Verworfen: immer erst einen Plan drucken und dann schreiben - das 'vorher' waere nur die Reihenfolge der Zeilen. Verworfen: Rueckfrage j/n - haelt Agenten an.
- Bekannte Grenze: mit --lanes --dry-run zeigt die Block-Vorschau den Block fuer die Lanes von jetzt, nicht fuer die erneuerten.
- **2026-10-15 12:34 · BeMuCa** — Entscheidung 'von Hand geaendert' (die 11 falsch als edited gemeldeten Dateien vom 14.10.): eine Datei gilt nur dann als von Hand geaendert, wenn ihre Bytes KEINE Version sind, die jaira je ausgeliefert hat.
- Heute: core/role/install.go:101 writeFile vergleicht nur mit den Bytes DIESES Binaries. Eine Datei einer aelteren Version gilt als 'edited' und wird uebersprungen.
- Neu: core/shipped/shipped.sha256, eine Zeile '<sha256>  <Pfad im Repo>' je Version jeder installierten Datei (Rollen, cockpit-Mod, jaira-Skill, eingebaute Lanes). Per go:embed im Binary, ca. 430 Zeilen, ca. 45 KB.
- Erzeugt aus der git-Historie mit scripts/shipped-history.sh: --full-history -m (sonst fehlen 22 Versionen aus Merge-Aufloesungen), verweigert flache Klone.
- Historie: Rollen seit b9eae909 (v0.2.0) in core/role/builtin, vorher nie im Binary. Lanes seit v0.1.0 in core/lane/builtin. critique/optimize/testing lagen vorher in lanes/ (Katalog, verschoben in 2fb8f177) - deren Versionen zaehlen mit.
- Aktuell gehalten: ein Test je Paket verlangt, dass jede Datei, die dieses Binary installiert, gelistet ist, und druckt die fehlenden Zeilen. Das Skript nimmt auch die Dateien im Arbeitsbaum mit.
- Verworfen: alte Inhalte einbetten (Megabytes; lane/corrections.go macht das fuer EINE Datei). Verworfen: Install-Protokoll im Zielordner - hilft nicht fuer Dateien, die alte Binaries schon geschrieben haben.
- sha256 statt git-Blob-ID: mit sha256sum nachpruefbar, gleiches Format wie checksums.txt der Releases, unabhaengig vom git-Objektformat.
- **2026-10-15 12:35 · BeMuCa** — Entscheidung Lanes: 'jaira update --lanes' nutzt dieselbe Liste ausgelieferter Versionen wie die Rollen, nicht lanesEquivalent.
- Gleich der eingebauten Lane: unchanged. Gleich einer frueher ausgelieferten Version: wird ersetzt ('updated'). Sonst: 'skipped', nichts angefasst, Exit 3.
- Nur Lanes, die das Board hat UND die jaira einbaut. Eigene Lanes und Lanes, die das Board nicht hat, bleiben unberuehrt (hinzufuegen ist 'jaira lanes add').
- Warum nicht lanesEquivalent (core/lane/lane.go:893): es vergleicht Verhalten, nicht Bytes. Ein umformulierter Prompt waere 'gleich' und wuerde ueberschrieben. Und es braeuchte die alten Inhalte.
- Auf diesem Board: optimize und testing = alte Katalog-Version aus lanes/ -> werden ersetzt. critique ist von Hand geaendert (X9GJ51, 'notes' in input-requires) -> skipped.
- Der Kern liegt in core/lane (RefreshShipped). 0S1A7M ('lanes use --all') kann ihn spaeter mit einem force-Parameter nach aussen geben. 0S1A7M wird hier NICHT umgesetzt.
- **2026-10-15 12:35 · BeMuCa** — Entscheidung DoD 5 - self upgrade: nach dem Ersetzen startet 'jaira self upgrade' das NEUE Binary am Zielpfad als Kindprozess: '<ziel> -C <dir> update --roles' (plus --json, wenn das Upgrade --json hat).
- Warum Kindprozess: der laufende Prozess ist danach immer noch das alte Binary (internal/cli/self.go:170). Was er aus seinem eigenen embed schreibt, ist die alte Version.
- Im Board: Agent-Block + Skills (Projekt + persoenlich). Ausserhalb eines Boards: 'update --roles' zieht nur die persoenlichen Skills nach und sagt das.
- Lanes NICHT: sie liegen auf manchen Boards unter git und gehoeren dem Team. Ein Upgrade auf einem Rechner soll den Ablauf aller nicht still aendern. Die Hinweiszeile von update nennt veraltete Lanes.
- Exit 3 des Kinds (eine Datei von Hand geaendert) ist kein Fehler des Upgrades: Exit 0, die Zeilen des Kinds nennen, was uebersprungen wurde. JSON: 'setup' (Ausgabe des Kinds) und 'setup_exit'.
- Test: Variable runUpgraded (wie selfExecutable, self.go:21) + Fake-Release ueber JAIRA_RELEASE_API/JAIRA_RELEASE_DOWNLOADS. Dazu ein unix-Test, der ein Shell-Skript als neues Binary wirklich ausfuehrt.
- Nicht drin: die Verweigerung fuer Dev-Builds (self.go:80-84, 'a dev build has no published release'). Offene Frage an Berk.
- **2026-10-15 12:35 · BeMuCa** — Entscheidung DoD 6 - Warnung bei aelteren Skills: einmal je Binary-Version, auf stderr.
- 'aelter' heisst: eine Datei, die 'jaira update --roles' schreiben oder ersetzen wuerde - sie fehlt, oder sie ist eine frueher ausgelieferte Version. Von Hand geaenderte Dateien sind NICHT aelter.
- Geprueft werden nur die persoenlichen Ordner (~/.claude/skills, ~/.copilot/skills), in denen Rollen liegen. Projekt-Skills reisen mit dem Repo; die erreicht der bestehende Hinweis 'run jaira update' (nudgeIfStale) und die Hinweiszeile von update.
- Wo: openStore (internal/cli/root.go:252), direkt nach nudgeIfStale. Nicht im TUI (board.go:22 und home.go:49 gehen nicht durch openStore). Nie auf stdout, auch nicht bei --json.
- 'einmal': Stempel $JAIRA_HOME/skills-check/version (release.Stamp/Stamped). Steht dort die laufende Version, kostet es eine kleine Datei lesen - die Startzeit bleibt.
- Der Text nennt 'jaira update --roles'. Still bei dev-Builds und waehrend 'update' und 'roles install' (die sind die Abhilfe).
- Kein Ueberlapp mit T8A8KM: das zeigt Release-Notizen (TUI-Popup, Pille, Hervorhebung in update). Dies hier sagt nur, dass Skills nachzuziehen sind.
- **2026-10-15 12:35 · BeMuCa** — Entscheidung jaira-Skill (Schritt 9, braucht Berks Ja): das Binary liefert ihn mit.
- Heute: .claude/skills/jaira/SKILL.md liegt nur im Repo. 'roles install' schreibt ihn nicht (README.md:683). Das Binary kennt ihn nicht - ohne das kann DoD 6 ihn nicht pruefen und keinen Befehl nennen.
- go:embed kann nicht nach ../ greifen. Daher eine Kopie unter core/role/jaira/SKILL.md, die eingebettet wird. Ein Test verlangt, dass .claude/skills/jaira/SKILL.md byte-gleich ist (sonst: eine bearbeiten, ueber die andere kopieren).
- Installiert NEBEN den Rollen wie das cockpit-Mod, KEINE Rolle: 'roles list' zeigt ihn nicht, InstalledIn zaehlt ihn nicht. Sonst wuerde 'update --roles' in diesem Repo alle Rollen nach .claude/skills schreiben.
- Die 33 alten Versionen aus .claude/skills/jaira kommen unter dem neuen Pfad in die Liste - Berks von Hand kopierte ~/.claude/skills/jaira wird damit erkannt.
- Folge, die Berk bestaetigen sollte: 'jaira roles install --project' schreibt in jedem Repo auch .claude/skills/jaira/SKILL.md.
- **2026-10-15 12:35 · BeMuCa** — Entscheidung DoD 4 - Release-Notizen, je Schritt eine Zeile in core/release/notes/PWAF22.md, im selben Commit wie der Code:
- roles install: Datei gleich einer frueheren ausgelieferten Version wird ohne --force ersetzt und als 'updated' gemeldet.
- update --roles (wo installiert; ausserhalb eines Boards nur persoenlich).
- update --lanes (bearbeitete Lane: genannt, nicht angefasst, Exit 3).
- update --dry-run.
- nackter update: schreibt weiter nur den Block, nennt aber, was zurueckhaengt.
- self upgrade fuehrt danach 'jaira update --roles' des neuen Binaries aus; Lanes bleiben Handarbeit.
- einmalige stderr-Zeile je Version, wenn persoenliche Skills aelter sind.
- (nur mit Schritt 9) roles install schreibt den jaira-Skill mit.
Wortlaut steht im Design (Scratchpad design.md, je Schritt). Kein '- !'-Marker - der gehoert zu T8A8KM.
- **2026-10-15 12:46 · BeMuCa** — Entscheidung jaira-Skill (Treiber, 15.10.): Schritt 9 wird gebaut, ohne auf Berk zu warten.
- Grund 1: DoD 6 nennt den jaira-Skill ausdruecklich. Ohne eingebettete Kopie kann das Binary ihn weder pruefen noch nachziehen.
- Grund 2: Berk am 14.10.: nach einem Upgrade sollen die jaira-Skills immer neu installiert werden. Der jaira-Skill ist der, der sich am oeftesten aendert.
- Folge, die Berk beim signoff zuruecknehmen kann: 'jaira roles install --project' schreibt in jedem Repo auch .claude/skills/jaira/SKILL.md. Der Schritt ist ein eigener Commit und laesst sich allein zuruecknehmen.
- Offen fuer Berk (nicht in diesem Ticket): die Dev-Build-Verweigerung von 'self upgrade' (self.go:80-84), und Rollen, die mit --into installiert wurden (update --roles und die Warnung sehen sie nicht).
