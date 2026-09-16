---
id: 01M2KBPPVH5PAKZ98B0X7KX89C
title: "Der Dispatcher-Skill faehrt im Repository mit, samt spawn.sh und seinem --no-worktree"
status: testing
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Jeder, der dieses Repository klont, hat den Dispatcher-Skill und sein spawn.sh mitsamt --no-worktree - ohne dass ihm jemand Dateien aus seinem Heimverzeichnis schickt."
context: |-
  Der Dispatcher-Skill liegt heute NUR in ~/.claude/skills/jaira-dispatcher/ auf Alex' Maschine. Wer das Repository klont, bekommt ihn nicht. Das Repository traegt aber schon Skills: .claude/skills/jaira/SKILL.md ist da und wird mitgeklont - der Ort ist also erprobt, nur der Dispatcher steht nicht drin.

  Ausgeloest am 2026-09-15: Alex wollte, dass ein Worker nicht immer in einem eigenen Worktree startet, sondern auch in der Hauptordner-Arbeitskopie. Der Schalter wurde gebaut - scripts/spawn.sh nimmt jetzt --no-worktree (und liest JAIRA_NO_WORKTREE=1), SKILL.md beschreibt ihn. Beides steht in ~/.claude und ist damit auf genau einer Maschine vorhanden. Alex' Satz dazu: 'сделай тикет под это чтобы все могли пользоваться'.

  Was --no-worktree tut: wt = repo-root statt ../.worktrees/<repo>-<slug>; kein Worktree, kein eigener Branch. Gedacht fuer eine einzelne Lane auf dem schon ausgecheckten Branch. Der Preis steht im Hilfetext: zwei Worker teilen sich dann ein Verzeichnis, also nur einer gleichzeitig.

  Zu klaeren und Sache der Plan-Lane, nicht hier entschieden:
  - ob der Skill ins Repository kopiert wird oder ob ~/.claude/skills/jaira-dispatcher ein Symlink dorthin wird, damit es nicht zwei driftende Kopien gibt
  - ob die anderen jaira-role-*-Skills denselben Weg gehen; sie haben dasselbe Problem, sind aber nicht das, was Alex verlangt hat
  - spawn.sh setzt HERDR_* voraus. Auf einer Maschine ohne Herdr muss der Skill das sagen und nicht scheitern

  Die aktuellen Dateien liegen zum Uebernehmen bereit unter ~/.claude/skills/jaira-dispatcher/ (SKILL.md und scripts/spawn.sh, Stand 2026-09-15); scripts/*.bak sind alte Sicherungen und gehoeren NICHT mit ins Repository.
definition-of-done: "Der Skill jaira-dispatcher liegt unter .claude/skills/ im Repository, neben dem schon vorhandenen jaira-Skill; scripts/spawn.sh faehrt mit und kennt --no-worktree; ein frisch geklontes Repository kann einen Worker starten, ohne dass etwas aus ~/.claude kopiert wird"
tags:
  - docs
blocked-by: []
related: []
commits: []
created-at: 2026-09-15T20:21:31Z
updated-at: 2026-09-16T07:12:47Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-30509
claimed-at: 2026-09-16T06:49:05Z
outcome-what: "Zwei Aufraeumungen in core/role/builtin/jaira-dispatcher/scripts/spawn.sh, ohne Verhaltensaenderung: der dispatch-Zweig ruft send-text nicht mehr zweimal auf, sondern setzt $prompt und der Aufruf steht einmal danach; der Kommentar am --no-worktree-Zweig ist von 6 auf 4 Zeilen gekuerzt, weil er usage() wortgleich wiederholte."
outcome-why: "optimize: eine Idee soll an einer Stelle stehen. Der doppelte send-text-Aufruf war die einzige echte Duplikation dieser Aenderung, und sie zu falten macht auch RPXJ (testing -> /jaira-role-tester) zu einem elif statt zu einem dritten kopierten Aufruf."
outcome-resolves: "bash -n gruen, 'go test ./...' gruen (inkl. TestEmbeddedScriptsParse, das spawn.sh --help ohne Herdr startet); DoD 1 unveraendert erfuellt, kein Kommando und keine Ausgabe geaendert, also keine NOTES.md-Zeile."
review-summary: none
review-gaps: "entfernt: der doppelte send-text-Aufruf im dispatch-Zweig (beide Zweige setzen jetzt $prompt, der Aufruf steht einmal danach - und RPXJ braucht dort nur noch ein elif); gekuerzt: der 6-Zeilen-Kommentar am --no-worktree-Zweig, der usage() 100 Zeilen weiter oben wortgleich wiederholte, auf die eine Begruendung, die dort nicht steht. Gesucht und nicht gefunden: eine zweite spawn.sh oder ein zweites usage() im Repository (find/grep, es gibt genau eines). Gelassen und warum: der '--'-Fall der Flag-Schleife (eine Zeile Absicherung, loeschen waere selbst eine Verhaltensaenderung), der 'scripts == 0'-Guard in TestEmbeddedScriptsParse (deckt sich mit TestDispatcherShipsItsScript, ist aber der Guard dieses Tests und liest keinen anderen), und die vorbestehende 'dir := t.TempDir(); Install(dir, false)'-Wiederholung in core/role/role_test.go. NICHT generalisiert: die Lane-zu-Kommando-Tabelle - siehe Notiz, das ist RPXJ."
test-verdict: "pass: go build + 'go test -count=1 -race ./...' green (RC=0), GOOS=windows vet+build green; DoD 1 verified in the tree — 'jaira roles install --into <leer>' schreibt SKILL.md und scripts/spawn.sh (0755) aus dem Binary, kein Zugriff auf ~/.claude; Verhalten mit Herdr-Stub durchgespielt: --no-worktree und JAIRA_NO_WORKTREE=1 starten in $root ohne Worktree, lane 'dispatch' sendet /jaira-dispatcher, jede andere Lane /jaira-role-lane, Default-Pfad legt weiterhin .worktrees/repo-SLUG auf feat/SLUG an"
---

# Der Dispatcher-Skill faehrt im Repository mit, samt spawn.sh und seinem --no-worktree

## Definition of Done

- [x] Der Rollen-Ordner core/role/builtin/jaira-dispatcher/ traegt SKILL.md und scripts/spawn.sh im Stand von ~/.claude (2026-09-15); spawn.sh kennt --no-worktree; 'jaira roles install' in ein leeres Verzeichnis liefert beide Dateien, spawn.sh ausfuehrbar, ohne dass etwas aus ~/.claude kopiert wird
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:7-40,55-68; core/role/builtin/jaira-teamlead/SKILL.md:44-49; core/role/role_test.go TestEmbeddedScriptsParse + TestInstallWritesEveryFile

## Options

- [ ] brainstorm
- [x] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] Hunk --no-worktree aus ~/.claude/skills/jaira-dispatcher/scripts/spawn.sh nach core/role/builtin/jaira-dispatcher/scripts/spawn.sh uebernehmen: usage(), Flag-Schleife, wt-Zweig
- [x] Hunks 'dispatch'-Lane und proj-Variable mitnehmen; die vier .env-Zeilen (VITE_PORT_HOST, BACKEND_PORT_HOST, IMAGE_NS, COMPOSE_PROFILES) auslassen
- [x] Den --no-worktree-Absatz in core/role/builtin/jaira-dispatcher/SKILL.md nachziehen, wortgleich zum Stand in ~/.claude
- [x] Pruefen, dass ohne Herdr weder Skript noch Skill ins Leere laufen: HERDR_ENV-Check meldet sauber, SKILL.md nennt die Ersatzwege
- [x] Test in core/role: 'bash -n' ueber jede eingebettete .sh-Datei, und 'spawn.sh --help' endet mit 0 - heute faellt ein kaputtes Skript erst beim Worker auf
- [x] Zeile unter '## Unreleased' in core/release/NOTES.md: --no-worktree da, und 'jaira roles install --force' noetig, weil eine veraenderte Kopie uebersprungen wird
- [x] Abnahme: 'jaira roles install' in ein leeres Verzeichnis, beide Dateien da, spawn.sh ausfuehrbar, --no-worktree drin; 'go test ./core/role/...' gruen

## Progress
- **2026-09-16 06:52 · Alexander Sacharov** — Die Praemisse des Tickets ist ueberholt: der Dispatcher liegt laengst im Repository, nur an einem anderen Ort als die DoD sagte. core/role/builtin/jaira-dispatcher/ traegt SKILL.md und scripts/spawn.sh, sie sind per go:embed im Binary und 'jaira roles install' schreibt sie nach ~/.claude/skills. .claude/skills/ waere ein zweiter Ort und damit genau die driftende Zweitkopie, gegen die das Ticket existiert. DoD entsprechend umgeschrieben; Symlink-Frage damit erledigt, das Installationskommando ist der Weg.
- **2026-09-16 06:52 · Alexander Sacharov** — Der Diff ~/.claude gegen core/role/builtin ist groesser als --no-worktree. Vier Hunks: (1) --no-worktree samt usage() und Flag-Schleife - das ist der Auftrag; (2) Lane-Name 'dispatch' startet /jaira-dispatcher statt /jaira-role-lane - generisch, kommt mit; (3) proj-Variable statt Inline-Ausdruck, reines Refactoring, kommt mit; (4) vier .env-Zeilen VITE_PORT_HOST, BACKEND_PORT_HOST, IMAGE_NS, COMPOSE_PROFILES - die nennen 'task app:dev' eines anderen Projekts und gehoeren nicht in eine generische Builtin-Rolle. Bleiben draussen, Folgeticket.
- **2026-09-16 06:52 · Alexander Sacharov** — jaira roles install ueberschreibt eine veraenderte Datei nicht, es meldet sie als 'skipped' (core/role/install.go). Alex' ~/.claude-Kopie ist veraendert, also bleibt sie nach diesem Ticket die Abweichung, bis er einmal 'jaira roles install --force' laeuft. Das gehoert in die NOTES.md-Zeile, sonst merkt es niemand.
- **2026-09-16 06:52 · Alexander Sacharov** — Ohne Herdr scheitert spawn.sh nicht stumm: es prueft HERDR_ENV=1 und sagt 'not inside a Herdr pane', exit 1. SKILL.md nennt zwei Ersatzwege (Peer-Sessions, Subagents). Die dritte offene Frage aus dem Kontext ist damit beantwortet, es bleibt ein Pruefschritt und keine Aenderung.
- **2026-09-16 06:57 · Alexander Sacharov** — Planschritt 2 nannte eine "proj-Variable" und vier .env-Zeilen (VITE_PORT_HOST, BACKEND_PORT_HOST, IMAGE_NS, COMPOSE_PROFILES) aus ~/.claude. Beides existiert dort heute nicht mehr: die ~/.claude-Kopie von spawn.sh ist am 16.09. auf den .jaira/worktree-setup-Hook umgebaut worden, der die .env-Logik ganz aus dem Skript nimmt. Dieser Hook ist Ticket 7MG5GB und bleibt hier draussen. Uebernommen wurden daher nur zwei Hunks: --no-worktree und die dispatch-Lane.
- **2026-09-16 06:57 · Alexander Sacharov** — spawn.sh und ~/.claude/skills/jaira-dispatcher/scripts/spawn.sh sind nach diesem Ticket bewusst NICHT byte-gleich - die Builtin-Version traegt noch die inline .env-Logik, die ~/.claude-Version schon den worktree-setup-Hook aus 7MG5GB. SKILL.md dagegen ist byte-gleich (diff leer). Wer die beiden vergleicht und einen Rueckstand vermutet: es ist der andere Ticket-Weg, kein vergessener Hunk.
- **2026-09-16 06:57 · Alexander Sacharov** — Die Flag-Schleife steht vor dem HERDR_ENV-Check, nicht danach. Damit beantwortet spawn.sh --help die Flags auch auf einer Maschine ohne Herdr mit exit 0, statt mit "not inside a Herdr pane" abzubrechen. TestEmbeddedScriptsParse setzt HERDR_ENV= genau deswegen explizit - wer den Check nach oben schiebt, faellt dort auf.
- **2026-09-16 07:00 · Alexander Sacharov** — critique: der --no-worktree-Hunk ist sauber und dokumentiert, der zweite Hunk nicht. Die 'dispatch'-Lane (spawn.sh:137) wurde bewusst mitgenommen - aber niemand erfaehrt davon: teamlead/SKILL.md:45 startet den Dispatcher mit spawn.sh und nennt den Lane-Wert nicht, dispatcher/SKILL.md:84 beschreibt die Signatur ohne ihn, NOTES.md:17 nennt nur --no-worktree. Damit ist der Zweig heute unerreichbar, obwohl er ausgeliefert wird. Nicht rausgeworfen, sondern dokumentiert - das ist die kleinere Aenderung. Viertens: --no-worktree liest den Slug nie (wt=$root), verlangt ihn aber weiter; ein Satz in usage() reicht. Bewusst NICHT aufgemacht: ob der Skill zusaetzlich unter .claude/skills/ liegen soll - die Zweitkopie-Entscheidung steht in der Pre-Process-Notiz und bleibt stehen.
- **2026-09-16 07:01 · Alexander Sacharov** — Die vier Critique-Punkte sind an vier Stellen beantwortet, weil der dispatch-Zweig vier Leser hat: spawn.sh usage() (wer --help liest), dispatcher/SKILL.md (der Dispatcher selbst - dort steht ausdruecklich 'du uebergibst es nicht, du bist was es startet'), teamlead/SKILL.md (der einzige, der 'dispatch' je tippt - dort mit der Warnung, dass jeder andere Name eine Lane statt eines Dispatchers startet) und NOTES.md. Den Zweig zu loeschen waere kleiner gewesen, haette aber den einzigen Weg entfernt, einen Dispatcher in einen eigenen Tab zu bekommen - teamlead/SKILL.md verlangt genau das.
- **2026-09-16 07:03 · Alexander Sacharov** — critique (2. Durchgang): nichts mehr zu sagen. Alle vier Punkte des ersten Durchgangs sind beantwortet - dispatch steht jetzt in spawn.sh usage(), dispatcher/SKILL.md, teamlead/SKILL.md (mit dem Aufrufbeispiel) und als eigene NOTES.md-Zeile; der ungenutzte Slug steht in usage() und SKILL.md. Bewusst NICHT neu aufgemacht: (a) dass derselbe Satz an vier Stellen steht - genau das hat der erste Durchgang verlangt, die vier Leser sind verschieden; (b) dass JAIRA_NO_WORKTREE nur den Wert 1 akzeptiert und ein 'true' stumm ignoriert - usage() sagt =1, das ist die dokumentierte Schnittstelle; (c) die .claude/skills-Zweitkopie, die in der Pre-Process-Notiz entschieden wurde.
- **2026-09-16 07:06 · Alexander Sacharov** — optimize: der 'dispatch'-Zweig in spawn.sh (jetzt Zeile 139) ist die erste Haelfte von Ticket RPXJ ('spawn.sh kann nur eine Rolle starten, obwohl der Prompt zwei verlangt'): dort soll testing nach /jaira-role-tester abgebogen werden, hier biegt dispatch nach /jaira-dispatcher ab. Zwei Namen, eine Idee. Bewusst NICHT hier generalisiert - eine Lane-zu-Kommando-Tabelle zu bauen waere Verhaltensaenderung und RPJXs Auftrag, nicht Aufraeumen. Vorbereitet ist es: der Zweig setzt jetzt nur noch $prompt, der send-text-Aufruf steht einmal danach, also kommt RPXJ mit einem weiteren elif durch.
- **2026-09-16 07:06 · Alexander Sacharov** — optimize: zwei Dinge geprueft und absichtlich gelassen. (1) Der '--)'-Fall in der Flag-Schleife ist Absicherung fuer einen Slug mit Bindestrich, den es auf diesem Board nicht gibt - eine Zeile, und sie zu loeschen waere selbst eine Verhaltensaenderung. (2) 'dir := t.TempDir(); Install(dir, false)' steht jetzt in sechs Tests in core/role/role_test.go - das ist vorbestehende Wiederholung aus fuenf aelteren Tests, nicht von dieser Aenderung eingefuehrt, also hier nicht angefasst.
- **2026-09-16 07:12 · Alexander Sacharov** — testing: Gates gruen ohne Cache — go build ./..., 'go test -count=1 -race ./...' RC=0, dazu 'GOOS=windows GOARCH=amd64 go vet ./...' und der Windows-Build, weil README:817-828 die vor dem Push verlangt. TestEmbeddedScriptsParse/TestInstallWritesEveryFile/TestDispatcherShipsItsScript einzeln -v gelaufen, alle PASS.
