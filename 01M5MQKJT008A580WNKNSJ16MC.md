---
id: 01M5MQKJT008A580WNKNSJ16MC
title: "Worker-Rolle: kein Version-Bump/Changelog, nur eigene Dateien formatieren"
status: done
ready: true
creator: Alexander Sacharov
goal: "Ein Worker ändert weder Version noch Changelog des Repos und formatiert nur Dateien, die er selbst geändert hat."
context: |-
  Gefunden am 23.10.2026 in einer langen Teamlead-Session für requirementsgenie (Notizen: /home/alex/.claude/acceptance/requirementsgenie/teamlead-staerken.md).
  Was falsch ist: Worker bumpten package.json und den Changelog selbst — beim Zusammenführen mehrerer Tickets auf einem Abnahme-Stand gab es bei jedem Ticket Konflikte; der Teamlead vergibt die Version beim Release. Außerdem formatiert 'task code:format' (ruff format) das ganze Repo (238 Dateien), Worker mussten das jedes Mal zurücknehmen.
  Beides wurde per Hinweis an jeden einzelnen Worker gelöst — gehört in die Rolle.
  Wo: Rollen jaira-worker / jaira-role-lane (core/role/builtin). Formatieren: nur 'git diff --name-only' der eigenen Änderung an den Formatter geben. Repos, die den Bump im Ticket wollen, brauchen eine Einstellung (z. B. settings.json 'release': 'teamlead'|'worker'), Standard so wie heute dokumentiert.
definition-of-done: "Rolle sagt: Version/Changelog nicht anfassen, wenn das Repo es dem Teamlead überlässt (Einstellung)"
tags:
  - cli
blocked-by: []
related: []
commits:
  - 16353616e2f8a46b7c7bbbbf559107dbfe8eceed
  - 6a87d6b936b37227fa24c3d1ce482f896f17d75b
  - 44dc85d0ad452f24a7e4fec692481bec9cae4c14
created-at: 2026-10-23T15:55:52Z
updated-at: 2026-10-23T16:29:07Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-54412
claimed-at: 2026-10-23T15:56:31Z
mode: conversational
outcome-what: "Board-Datei .jaira/release (teamlead|worker, Standard worker) über settings.ReleaseFor, in 'jaira roles list' als release/release_source; jaira-worker und jaira-role-lane lassen bei teamlead Version, Changelog und Tags in Ruhe und formatieren immer nur eigene Dateien"
outcome-why: Worker bumpten Version und Changelog pro Ticket und formatierten das ganze Repo — jeder Merge mehrerer Tickets auf einen Abnahmestand gab Konflikte; bisher per Hinweis an jeden Worker gelöst
outcome-resolves: "DoD 1 (Rolle fasst Version/Changelog nicht an, wenn .jaira/release teamlead sagt) und DoD 2 (nur eigene geänderte Dateien formatieren)"
review-summary: "Neue Board-Datei .jaira/release (eine Zeile, # Kommentare erlaubt): steht dort 'teamlead', liefert settings.ReleaseFor teamlead, sonst worker (Standard, auch bei Tippfehler). 'jaira roles list' zeigt das als release/release_source in --json und als Zeile 'release: …' im Text, --help nennt die Datei. jaira-worker und jaira-role-lane lesen den Wert und fassen bei teamlead weder Versionsfeld noch Changelog noch Tag an – das sticht auch eine Repo-Anweisung, die pro Änderung bumpt. Beide Rollen sagen zusätzlich ohne Einstellung: nur eigene Pfade an den Formatter (git diff --name-only --diff-filter=d HEAD + git ls-files --others --exclude-standard), nie task code:format/ruff format ./npm run format. Nebenbei: Zeilen-Scan von BoardChatModel in boardLine() ausgelagert und mitbenutzt. Release-Note core/release/notes/SJ16MC.md mit zwei Zeilen."
round: 2
review-gaps: "Nichts blockiert. Drei Kleinigkeiten: 1) Groß-/Kleinschreibung: 'Teamlead' in .jaira/release fällt still auf worker zurück (getestet) – sichtbar nur über release_source='the default'; strings.EqualFold wäre billig. 2) Auf einem noch nicht geteilten Board ist .jaira/ gitignored (git check-ignore bestätigt) – die Release-Note sagt 'commit it', das geht erst nach 'jaira share'; die Einstellung wirkt lokal trotzdem. 3) 'git ls-files --others --exclude-standard' listet alle unversionierten Dateien, auch die einer anderen Session im selben Worktree – der Rollentext sagt 'your new files', der Befehl filtert das nicht."
test-verdict: "pass: go build/vet + go test -race ./... green (RC=0); DoD 1-2 in the tree (proof lines checked); scratch board: no file → worker (the default), '# who\\nteamlead' → teamlead from .jaira/release, typo 'teamlaed' → worker"
review-verdict: "Der Diff erfüllt beide DoD-Punkte; go test ./core/settings ./core/role ./internal/cli grün; Verhalten auf einem Scratch-Board nachgeprüft (Standard worker, teamlead mit Kommentar und Leerzeichen erkannt, validate meckert nicht). Keine Defekte, nur die drei kleinen Punkte aus review-gaps – keiner davon ein Grund zum Zurückschicken."
review-round: 2
review-check: "1. go build -o /tmp/jaira ./cmd/jaira  2. mkdir /tmp/rb && cd /tmp/rb && git init && /tmp/jaira init  3. /tmp/jaira roles list | tail -1 → 'release: worker — … (the default)'  4. echo teamlead > .jaira/release  5. /tmp/jaira roles list --json | /tmp/jaira json get release → 'teamlead'  6. /tmp/jaira roles list | tail -1 → '(from .jaira/release, shared with the board)'  7. Rollentext lesen: core/role/builtin/jaira-worker/SKILL.md ab Zeile 179 und core/role/builtin/jaira-role-lane/SKILL.md ab Zeile 97 – Abschnitt Version/Changelog und 'format only the files you changed'"
---

# Worker-Rolle: kein Version-Bump/Changelog, nur eigene Dateien formatieren

## Definition of Done

- [x] Rolle sagt: Version/Changelog nicht anfassen, wenn das Repo es dem Teamlead überlässt (Einstellung)
  proof: core/role/builtin/jaira-worker/SKILL.md:179, core/role/builtin/jaira-role-lane/SKILL.md:97; setting core/settings/settings.go:262 (ReleaseFor, .jaira/release); TestReleaseDefaultsToTheWorker, TestRolesListReportsTheBoardsRelease, TestCodeChangingRolesLeaveTheReleaseAndFormatOnlyTheirFiles
- [x] Rolle sagt: nur eigene geänderte Dateien formatieren
  proof: core/role/builtin/jaira-worker/SKILL.md:189, core/role/builtin/jaira-role-lane/SKILL.md:106; TestCodeChangingRolesLeaveTheReleaseAndFormatOnlyTheirFiles

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/settings: Board-Datei .jaira/release lesen (BoardRelease), Standard worker, Tippfehler = worker
- [x] jaira roles list: release + release_source in --json und Textausgabe, --help nennt die Datei
- [x] Rollen jaira-worker + jaira-role-lane: Abschnitt Version/Changelog (liest release) und Formatieren nur eigener Dateien
- [x] Tests (settings, roles list, Rollentext) + core/release/notes/SJ16MC.md

## Progress
- **2026-10-23 15:56 · Alexander Sacharov** — parent session: requirementsgenie-92
- **2026-10-23 15:57 · Alexander Sacharov** — question (open decision): Wo lebt die Einstellung 'Version/Changelog gehört dem Teamlead'? 1) (Recommended) Board-Datei .jaira/release (eine Zeile teamlead|worker, committed, wie .jaira/chat-model) — gilt für alle im Repo, Rolle liest sie über 'jaira roles list --json | jaira json get release'; Standard worker = heutiges Verhalten. 2) ~/.jaira/settings.json 'release' — pro Maschine, wie im Ticket vorgeschlagen; gilt dann aber für jedes Repo auf der Maschine (dasselbe Problem wie 'remote'). 3) Keine Einstellung in jaira: Rolle sagt nur 'wenn CLAUDE.md/AGENTS.md des Repos Version/Changelog dem Release überlässt, nicht anfassen' — kein Go-Code. Formatieren (DoD 2) ist nicht offen: nur 'git diff --name-only' der eigenen Änderung an den Formatter, ohne Einstellung.
- **2026-10-23 15:59 · Alexander Sacharov** — decision (parent requirementsgenie-92, 23.10.): Option 1 — Board-Datei .jaira/release, eine Zeile teamlead|worker, committed; Standard worker (heutiges Verhalten). Rolle liest sie über 'jaira roles list --json'. requirementsgenie setzt danach teamlead.
- **2026-10-23 16:03 · Alexander Sacharov** — Formatieren ohne Einstellung: immer nur eigene Pfade (git diff --name-only --diff-filter=d HEAD + ls-files --others). Kein TUI-Schalter für .jaira/release — Datei von Hand, wie lanes/order; erst bauen, wenn jemand ihn vermisst. 'teamlead' sticht ausdrücklich eine Repo-Anweisung, die pro Änderung bumpt (sonst wäre die Einstellung in requirementsgenie wirkungslos). Eine Datei mit 'worker' meldet source 'the default', weil sie dasselbe bewirkt.
- **2026-10-23 16:04 · Alexander Sacharov** — decision (parent requirementsgenie-92): Diff abgenommen inkl. teamlead > Repo-Anweisung; Commit 16353616 vom Parent. Ab jetzt ohne Zwischenfreigabe bis signoff, Commits selbst (Ticket-ID, ohne Co-Authored-By).
- **2026-10-23 16:04 · Alexander Sacharov** — critique: ReleaseFor dupliziert den Zeilen-Scan von BoardChatModel (gleiche Datei) — gemeinsamer Helfer boardLine(path). Sonst nichts: release_source folgt chat_model_source/python_source; doppelter Rollentext folgt den schon doppelten Commit-Regeln; Dispatcher-Pfad deckt jaira-role-lane ab.
- **2026-10-23 16:05 · Alexander Sacharov** — critique Runde 2: nichts mehr — boardLine-Fund in 6a87d6b9 behoben.
- **2026-10-23 16:09 · Alexander Sacharov** — review (fresh): ok, keine Rücksendung. Kleinigkeiten für signoff: 'Teamlead' groß geschrieben → still worker; .jaira/release auf ungeteiltem Board gitignored; ls-files --others nimmt auch fremde unversionierte Dateien mit.
