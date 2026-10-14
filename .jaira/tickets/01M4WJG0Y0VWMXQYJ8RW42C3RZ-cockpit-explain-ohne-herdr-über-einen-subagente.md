---
id: 01M4WJG0Y0VWMXQYJ8RW42C3RZ
title: "Cockpit: Explain ohne Herdr über einen Subagenten statt Fehler"
status: critique
ready: true
creator: Alexander Sacharov
goal: "Explain im Cockpit funktioniert auch ohne Herdr: die Sitzung startet einen Subagenten, der den Review-Befund erklärt und als explain-Note ans Ticket schreibt. Discuss ohne Herdr läuft im aktuellen Chat. Die Fehlermeldung nennt die tatsächliche Ursache."
context: |-
  Was falsch ist: Auf Windows ohne Herdr zeigt [Explain] im Cockpit nur 'Explain 7MX3EJ failed - A chat of its own needs Herdr (HERDR_ENV=1).' Es passiert sonst nichts.
  Wo: core/mod/builtin/jaira-cockpit/hooks/register.tsx, Funktion openChat (Zeile ~551). Ohne HERDR_ENV=1 oder ohne spawn.sh bricht sie mit t.needsHerdr ab. Einen Weg über einen Subagenten gibt es nicht.
  Erwartet hat Alex (14.10.2026): ohne Herdr soll ein Subagent laufen.
  Was es schon gibt: $.prompt.submit schickt Text an die laufende Sitzung (benutzt in decide, register.tsx ~627). Der Explain-Prompt steht in core/role/builtin/jaira-dispatcher/scripts/spawn.sh (lane explain, Zeile ~247). Er schreibt 'jaira note <id> "explain: ..."', und das Cockpit liest diese Note schon unter [Read explanation].
  Discuss: ein Gespräch mit einer Person, das kann ein Subagent nicht führen. Ohne Herdr gehört es in den aktuellen Chat (per $.prompt.submit).
  Zweiter Fehler: derselbe Text t.needsHerdr erscheint bei drei Ursachen (kein HERDR_ENV, kein spawn.sh, keine Board). Man sieht nicht, welche.
definition-of-done: "Explain ohne HERDR_ENV=1 schickt per $.prompt.submit einen Auftrag an die Sitzung, einen Subagenten mit dem Explain-Prompt aus spawn.sh zu starten; kein Fehler-Toast"
tags:
  - tui
blocked-by: []
related: []
commits: []
created-at: 2026-10-14T06:44:46Z
updated-at: 2026-10-14T08:07:46Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-52389
claimed-at: 2026-10-14T07:50:01Z
outcome-what: "Cockpit: Explain ohne Herdr über einen Subagenten, Discuss im aktuellen Chat (beide über once()); eigene Meldungen für fehlendes spawn.sh und Board; natives Windows: USERPROFILE und Git Bash; spawn.sh legt den langen Explain/Discuss-Auftrag in eine Datei und tippt nur eine kurze Zeile"
outcome-why: "Ohne Herdr und auf nativem Windows endete Explain in 'needs Herdr'; mit Herdr kam der lange getippte Auftrag verstümmelt an, und keine explain-Note entstand"
outcome-resolves: "DoD 1-8: Explain/Discuss ohne Herdr, irreführende needsHerdr-Meldung, Windows in Herdr, verstümmelter Prompt im Herdr-Tab, Doppelklick auf Discuss"
review-summary: "Ohne Herdr fragt das Cockpit jetzt die eigene Sitzung per $.prompt.submit: Explain bittet um einen Subagenten, der den Befund erklärt und als explain:-Note ans Ticket schreibt; Discuss führt das Gespräch im selben Chat und hinterlässt eine decision:-Note. Beide Aufträge stehen in texts.ts (Querverweis zu spawn.sh). needsHerdr ist ersetzt durch noBoard und noSpawn. Auf nativem Windows liest home() USERPROFILE, und spawn.sh startet überall als 'bash spawn.sh' - auf Windows mit Git Bash aus CLAUDE_CODE_GIT_BASH_PATH oder neben git --exec-path, nie mit der WSL-bash aus System32. spawn.sh schreibt den langen Explain/Discuss-Auftrag in eine Datei in einem mktemp-Ordner, startet claude mit --add-dir auf diesen Ordner und tippt nur 'Lies die Datei ...' - damit kommt der Auftrag nicht mehr verstümmelt an. Discuss läuft jetzt wie Explain durch once()."
review-gaps: "1) DoD 8 gilt nur für gleichzeitige Klicks: once() gibt die Sperre im finally frei, sobald openChat zurückkehrt. Ohne Herdr ist das nach git worktree list + say + void submit - einige Millisekunden. Ein echter Doppelklick (zweiter Klick 100-300 ms später) schickt also zwei Aufträge. Der Test 'a double click ... sends one request' drückt per Promise.all gleichzeitig und sieht das nicht. Mit Herdr hält spawn.sh die Sperre Sekunden lang, dort stimmt es. 2) Windows-Teil (USERPROFILE, Git-Bash-Wahl, cygpath-Pfad, --add-dir in doppelten Anführungszeichen) nur gegen Mock getestet. 3) Der getippte Satz 'Lies die Datei <pfad> ...' setzt den Pfad nicht in Anführungszeichen; ein %TEMP% mit Leerzeichen (C:\\Users\\Vor Nach\\...) ist für das Modell mehrdeutig. 4) mktemp-Ordner werden nie gelöscht (bewusst, Note 07:58). Nicht vom Ticket: worktree hat ein untracked .jaira/milestones/ - nicht mitcommitten."
test-verdict: "pass: go build/vet RC=0, go test -race -count=1 ./... RC=0, plugin validate RC=0, plugin test 143 pass RC=0; DoD 1-8 am Baum geprüft; spawn.sh explain mit Herdr-Stub am echten Ticket 42C3RZ: Datei 5551 Zeichen mit 'jaira note 42C3RZ \"explain:' und ganzem brief, getippt nur 'Lies die Datei ... ganz und tu, was darin steht.', Start 'claude --add-dir <ordner>'; claude -p --add-dir liest die Datei ohne Rückfrage. Nicht geprüft: echtes Windows, echter Herdr-Tab"
review-verdict: "Annehmbar mit einem Vorbehalt. DoD 1-7 sind im Diff erfüllt; go test ./core/role ok, claude plugin test 143 pass/0 fail (selbst gelaufen 14.10.). DoD 8 ist schwächer als angehakt: nur gleichzeitige Klicks werden abgefangen, ein normaler Doppelklick ohne Herdr schickt zwei Aufträge (Gap 1). Wenn das wichtig ist: zurück nach in-progress und die Sperre im Weg ohne Herdr bis BUSY_MS halten. Windows-Teil logisch richtig, aber nur gemockt - von Hand prüfen."
review-round: 2
review-check: "Automatisch: 1. cd core/mod/builtin/jaira-cockpit && claude plugin test -> '143 pass, 0 fail'. 2. go test ./core/role -> ok. Ohne Herdr (WSL/Linux, HERDR_ENV nicht gesetzt): 3. Sitzung im Repo starten, Ticket mit Review-Befund in human/signoff. 4. Im Cockpit [Explain] drücken -> Toast '<ID>: a subagent of this session explains the finding...', im Chat erscheint der Auftrag, die Sitzung startet einen Agenten. 5. Danach steht [Read explanation] am Ticket mit Warum:/Optionen/Empfehlung:. 6. Bei Ticket mit offener question [Discuss] -> Frage wird im selben Chat erklärt; nach deiner Antwort: jaira show <id> --json zeigt decision:-Note und question leer. 7. Gap 1 prüfen: [Discuss] zweimal im Abstand von ca. 0,3 s drücken -> erscheint der Auftrag zweimal im Chat, ist Gap 1 bestätigt. Natives Windows in Herdr: 8. jaira roles install, dann [Explain] -> neuer Herdr-Tab öffnet sich (nicht WSL, kein 'needs Herdr'); im Tab steht nur 'Lies die Datei C:\\...\\<ID>-explain.md ...' und Claude liest sie ohne Rückfrage. 9. spawn.sh unter %USERPROFILE%\\.claude\\skills\\jaira-dispatcher\\scripts umbenennen, [Explain] -> Meldung nennt 'jaira roles install'."
round: 3
---

# Cockpit: Explain ohne Herdr über einen Subagenten statt Fehler

## Definition of Done

- [x] Explain ohne HERDR_ENV=1 schickt per $.prompt.submit einen Auftrag an die Sitzung, einen Subagenten mit dem Explain-Prompt aus spawn.sh zu starten; kein Fehler-Toast
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:568-571 openChat ohne HERDR_ENV -> $.prompt.submit(t.explainBySubagent) (texts.ts, Prompt wie spawn.sh lane explain); Test 'without Herdr, Explain asks this session for a subagent that writes the explain note'
- [x] Discuss ohne HERDR_ENV=1 führt das Gespräch per $.prompt.submit im aktuellen Chat; die Entscheidung landet als Note am Ticket
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:568-571 -> $.prompt.submit(t.discussInThisChat): decision-Note + question= leeren; Test 'without Herdr, Discuss talks the question through in this chat'
- [x] Kein spawn.sh bei HERDR_ENV=1 und keine Board haben je eine eigene Meldung in texts.ts
  proof: core/mod/builtin/jaira-cockpit/hooks/texts.ts noSpawn und noBoard statt needsHerdr; register.tsx:566 (noBoard) und :574 (noSpawn); Test 'with Herdr but no spawn.sh, Explain says spawn.sh is missing'
- [x] Tests in core/mod/builtin/jaira-cockpit/tests decken Explain und Discuss ohne Herdr ab
  proof: core/mod/builtin/jaira-cockpit/tests/board.test.ts: 5 neue Tests (ohne Herdr Explain/Discuss, kein spawn.sh, Windows USERPROFILE+Git Bash, CLAUDE_CODE_GIT_BASH_PATH); claude plugin test: 142 pass 0 fail, exit 0
- [x] Release-Note core/release/notes/<ID>.md
  proof: core/release/notes/42C3RZ.md
- [x] Natives Windows in Herdr: home() fällt auf USERPROFILE zurück, und spawn.sh wird über bash gestartet (Git Bash), damit Explain/Discuss dort einen Herdr-Tab öffnen
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx home() HOME||USERPROFILE mit /; spawnScript() gibt [bash, spawn.sh]; bash() nimmt auf OS=Windows_NT CLAUDE_CODE_GIT_BASH_PATH, sonst <git --exec-path>/../../../bin/bash.exe; Tests 'on native Windows, ...'
- [x] spawn.sh explain/discuss tippt nur noch eine kurze Zeile per send-text; der lange Auftrag samt brief liegt in einer Datei, die der Chat liest (auf Windows mit Windows-Pfad). Test: ein Prompt über 4000 Zeichen kommt ganz an
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:190 (mktemp-Ordner, claude --add-dir) und :281 (send-text nur 'Lies die Datei ...'); Tests TestSpawnChatsStartOnTheChatModelWithTheirBrief (Notiz über 4000 Zeichen landet ganz in der Datei, getippte Zeile <200) und TestSpawnStartsClaudeWhereItsShellRuns/git bash chat; go test ./core/role ok
- [x] Discuss ohne Herdr läuft durch once() wie Explain; ein Doppelklick schickt nur einen Auftrag
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:511 once() hält die Sperre bis BUSY_MS, wenn die Aktion 'hold' antwortet; :575 openChat ohne Herdr antwortet 'hold'; Discuss-Buttons über once(); Test 'without Herdr, a double click on Explain or Discuss sends one request' (zweiter Klick 300 ms später, fällt ohne 'hold'; 143 pass)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-14 06:45 · Alexander Sacharov** — parent session: jaira-34
- **2026-10-14 06:46 · Alexander Sacharov** — Offene Entscheidungen vor der Arbeit: keine, die nur Alex treffen kann. Selbst entschieden:
1) Prompt-Quelle: der Cockpit-Weg ohne Herdr trägt Explain- und Discuss-Auftrag selbst in texts.ts (englisch, wie alle Cockpit-Texte), inhaltlich gleich zu spawn.sh lane explain/discuss. Warum: ohne Herdr ist spawn.sh oft gar nicht installiert (jaira roles install), und DoD 3 nennt 'kein spawn.sh' nur bei HERDR_ENV=1. Kommentar an beiden Stellen verweist auf die andere.
2) Explain ohne Herdr: $.prompt.submit an die laufende Sitzung: starte einen Subagenten mit diesem Prompt; der Subagent liest selbst jaira show <id> --brief und schreibt die explain:-Note. Das Chat-Modell (.jaira/chat-model) wird dort nicht durchgereicht - Agent kennt nur sonnet/opus/haiku-Namen.
3) Discuss ohne Herdr: $.prompt.submit mit dem Discuss-Auftrag an die Sitzung selbst (kein SendMessage, kein Tab).
4) Meldungen: needsHerdr fällt weg; neu noSpawn (HERDR_ENV=1, spawn.sh fehlt) und noBoard.
- **2026-10-14 06:47 · Alexander Sacharov** — Nachtrag 14.10.: Alex sagt, auf Windows läuft die Sitzung IN Herdr und trotzdem kommt needsHerdr. Verdacht 1: home() in register.tsx:147 liest nur HOME; natives Windows hat oft nur USERPROFILE -> spawnScript sucht /.claude/skills/... und findet nichts -> needsHerdr. Verdacht 2: auch gefunden würde spawn.sh per $.process.run direkt gestartet; Windows führt .sh ohne bash nicht aus. Im Cockpit gibt es keinerlei Windows-Behandlung (kein USERPROFILE, kein bash).
- **2026-10-14 06:53 · Alexander Sacharov** — Umsetzung: spawn.sh startet jetzt überall als 'bash spawn.sh ...' (auch Linux/WSL) - ein Weg, ein Test. Auf Windows NICHT bare 'bash': CreateProcess sucht System32 vor PATH, dort ist bash.exe die WSL-bash. Daher CLAUDE_CODE_GIT_BASH_PATH (setzt Claude Code für Windows), sonst aus git --exec-path abgeleitet. Testmotor nimmt nur absolute POSIX-Pfade (C:/... wird vor dem Mock verworfen) und fällt bei fehlendem HOME offenbar auf das echte HOME zurück - daher HOME:'' und USERPROFILE '\\Users\\a' im Test. Echt auf Windows ungeprüft.
- **2026-10-14 06:54 · Alexander Sacharov** — critique: nichts. Geprüft: Duplikat der Prompts zu spawn.sh steht auf der Note 'Offene Entscheidungen' (bewusst, Querverweis in texts.ts und spawn.sh); 'bash spawn.sh' statt Shebang ist auf Linux/macOS gleichwertig (Shebang ist /usr/bin/env bash, also dieselbe PATH-Suche); bash() außerhalb des try in spawnScript, damit ein Fehler dort nicht als 'kein spawn.sh' erscheint.
- **2026-10-14 07:01 · Alexander Sacharov** — review: DoD erfüllt, 142/0 grün. Offen nur: Windows-Teil gegen echten Windows-Rechner ungeprüft; Discuss ohne once()-Schutz (Doppelklick = zwei Aufträge im Chat); Prompt-Duplikat spawn.sh/texts.ts ohne Test.
- **2026-10-14 07:49 · Alexander Sacharov** — Befund 14.10. (Alex: 'Explain ist nicht zurückgekommen'): Der Herdr-Weg öffnete den Tab 42C3RZ/explain, aber der Prompt kam verstümmelt an. Die Sitzung bekam 1395 Zeichen: 'Erklaere mir in einfachen Worten' + 's. Im Cockpit gibt es...' (Rest des brief-JSON). Die ganze Mitte fehlte, also auch der Auftrag 'jaira note ... explain:'. Deshalb schrieb der Chat keine Note und [Read explanation] erschien nie. Ursache: spawn.sh (Zeile ~256-258) hängt das ganze brief-JSON an den Prompt und tippt alles per 'herdr pane send-text' in die TUI. Lange Texte verlieren auf dem Weg herdr.exe -> wsl.exe -> claude die Mitte. Beleg: ~/.claude/projects/-home-alex-projects--worktrees-jaira-42c3rz/7afdc53f-*.jsonl, erste user-Nachricht.
- **2026-10-14 07:51 · Alexander Sacharov** — DoD 7, Entscheidung: der lange Auftrag liegt in einer Datei in einem eigenen mktemp -d-Ordner (nicht im Worktree: dort stünde er als ?? in git status und könnte mitcommittet werden). claude startet im Tab mit --add-dir <dieser Ordner>, sonst fragt Read nach Erlaubnis - geprüft 14.10.: claude -p liest /tmp/x ausserhalb des cwd nicht ('Read verweigert'). send-text tippt nur 'Lies die Datei <pfad> ganz und tu, was darin steht.' Auf nativem Windows Pfad per cygpath -w, im Startbefehl in doppelten Anführungszeichen (cmd und PowerShell 5.1 kennen beide).
- **2026-10-14 07:58 · Alexander Sacharov** — critique Runde 2: nichts. Live geprüft 14.10.: claude --add-dir <mktemp-Ordner> liest die Datei dort ohne Rückfrage (haiku, -p, default-Modus: 'pelican'), ohne --add-dir verweigert. Achtung: --add-dir ist variadisch und frisst ein folgendes Positionsargument - in spawn.sh unschädlich, weil der Auftrag erst danach getippt wird. Bewusst stehen gelassen: die mktemp-Ordner werden nicht aufgeräumt (je Explain/Discuss eine kleine Datei; /tmp leert das System, %TEMP% auf Windows nicht) - Aufräumen bräuchte einen zweiten Schritt nach dem Lesen, den es heute nicht gibt.
- **2026-10-14 08:05 · Alexander Sacharov** — review (fresh): annehmbar mit Vorbehalt. DoD 8 nur für gleichzeitige Klicks wahr - once() gibt die Sperre frei, sobald openChat ohne Herdr nach Millisekunden zurückkehrt; ein echter Doppelklick schickt zwei Aufträge. Test nutzt Promise.all und sieht das nicht. Tests selbst gelaufen: plugin 143/0, go test ./core/role ok. Windows nur gemockt.
- **2026-10-14 08:05 · Alexander Sacharov** — Runde 3 aus signoff zurück (review gap 1 + 3): once() gab die Sperre im finally frei, sobald openChat zurückkehrte - ohne Herdr nach Millisekunden, ein echter Doppelklick (100-300 ms) schickte zwei Aufträge; der Test drückte gleichzeitig und sah das nicht. Fix: der Weg ohne Herdr hält die Sperre bis BUSY_MS. Dazu: Pfad im getippten Satz 'Lies die Datei ...' in Anführungszeichen (%TEMP% mit Leerzeichen).
