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
updated-at: 2026-10-14T07:57:52Z
assignee: Alexander Sacharov
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-52389
claimed-at: 2026-10-14T07:50:01Z
outcome-what: "Cockpit: Explain ohne Herdr über einen Subagenten, Discuss im aktuellen Chat (beide über once()); eigene Meldungen für fehlendes spawn.sh und Board; natives Windows: USERPROFILE und Git Bash; spawn.sh legt den langen Explain/Discuss-Auftrag in eine Datei und tippt nur eine kurze Zeile"
outcome-why: "Ohne Herdr und auf nativem Windows endete Explain in 'needs Herdr'; mit Herdr kam der lange getippte Auftrag verstümmelt an, und keine explain-Note entstand"
outcome-resolves: "DoD 1-8: Explain/Discuss ohne Herdr, irreführende needsHerdr-Meldung, Windows in Herdr, verstümmelter Prompt im Herdr-Tab, Doppelklick auf Discuss"
review-summary: "Cockpit openChat prüft jetzt zuerst das Board (eigene Meldung noBoard), dann HERDR_ENV. Ohne Herdr startet es nichts mehr, sondern schickt per $.prompt.submit einen Auftrag an die eigene Sitzung: Explain = 'starte einen Subagenten mit diesem Prompt', der mit jaira -C <ticket-dir> show --brief liest und eine explain:-Note im Format schreibt, das [Read explanation] liest (Warum:/Empfehlung:/[give back]/[accept]); Discuss = das Gespräch läuft im aktuellen Chat und endet in einer decision:-Note plus question= leeren. Mit Herdr, aber ohne spawn.sh, kommt jetzt noSpawn ('jaira roles install') statt 'needs Herdr'. Zusätzlich für natives Windows: home() nimmt USERPROFILE, wenn HOME fehlt (Backslashes zu /), und spawn.sh wird überall als [bash, spawn.sh] gestartet - auf Windows mit CLAUDE_CODE_GIT_BASH_PATH oder der bash.exe drei Ebenen über git --exec-path, damit nicht die WSL-bash aus System32 greift. Das gilt auch für Accept (pr/next) und Give back. Die Prompts in texts.ts sind englische Kopien von spawn.sh lane explain/discuss, Querverweise an beiden Stellen."
review-gaps: "1) Echtes natives Windows ist ungeprüft (steht selbst in der Note): der Test kann nur POSIX-Pfade, USERPROFILE ohne Laufwerksbuchstabe; ob $.fs.exists 'C:/Users/...' und Herdr unter Git Bash wirklich laufen, zeigt erst ein Windows-Rechner. 2) Die Explain/Discuss-Prompts existieren jetzt zweimal (spawn.sh deutsch, texts.ts englisch) - bewusst entschieden und mit Querverweis, aber Drift ist nur per Hand zu verhindern; kein Test vergleicht die beiden. 3) Discuss hat keinen once()-Schutz (Explain schon): ohne Herdr schickt ein Doppelklick den Auftrag zweimal in den Chat - vorher öffnete er zwei Tabs, also nicht neu, aber jetzt billiger auszulösen. 4) Ein Ticketpfad mit einfachem Anführungszeichen bricht das '${dir}'-Quoting im Prompt - Randfall. Kein Fehler gegen die DoD gefunden; alle sechs DoD-Punkte sind im Diff belegt."
test-verdict: "pass: go build/vet RC=0, go test -race -count=1 ./... RC=0, claude plugin validate RC=0, claude plugin test 142 pass RC=0; DoD 1-6 je am Baum und am genannten Test geprüft; 'bash spawn.sh' startet das Skript (Usage-Fehler wie erwartet). Nicht live gedrückt: die Band-Buttons brauchen eine interaktive Sitzung, und echtes Windows ist hier nicht verfügbar"
review-verdict: "Der Diff erfüllt alle sechs DoD-Punkte; claude plugin test im Cockpit: 142 pass, 0 fail. Keine Defekte im Linux/WSL-Pfad gefunden. Unsicher bin ich nur beim Windows-Teil: er ist logisch richtig (Git-Bash-Wahl, USERPROFILE), aber nur gegen einen Mock getestet - annehmbar, wenn Alex ihn auf Windows einmal von Hand prüft."
review-round: 1
review-check: "Ohne Herdr (WSL/Linux): 1. In einer Sitzung ohne HERDR_ENV ein Repo mit einem Ticket in human/signoff öffnen, das einen Review-Befund hat. 2. Im Cockpit [Explain] drücken. 3. Erwartet: Toast '<ID>: a subagent of this session explains the finding...', kein 'Explain failed'; im Chat erscheint der Auftrag und die Sitzung startet einen Agenten. 4. Danach erscheint [Read explanation] am Ticket mit Befund/Warum/Optionen/Empfehlung. 5. Bei einem Ticket mit offener question [Discuss] drücken: die Frage wird im selben Chat erklärt; nach deiner Entscheidung steht eine decision:-Note am Ticket und question ist leer. Auf Windows in Herdr: 6. Herdr-Sitzung auf nativem Windows starten, [Explain] drücken. 7. Erwartet: ein neuer Herdr-Tab öffnet sich (nicht 'needs Herdr', nicht WSL). 8. jaira roles install rückgängig/spawn.sh umbenennen und [Explain] erneut: Meldung nennt 'jaira roles install'. Automatisch: cd core/mod/builtin/jaira-cockpit && claude plugin test -> 142 pass."
round: 2
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
  proof: core/mod/builtin/jaira-cockpit/hooks/register.tsx:1118 und :1160 Discuss über once(); Test 'without Herdr, a double click on Explain or Discuss sends one request' (fällt ohne den Fix, 143 pass)

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
