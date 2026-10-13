---
id: 01M4VCZS1N1T0BHQVR8ATP44SC
title: "Rollen-Skripte starten Tabs auch auf nativem Windows (Git Bash): kein python3, WSL nicht an .exe erkannt"
status: done
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "spawn.sh und run-lane.sh starten auf nativem Windows (Git Bash + Herdr fuer Windows) einen Worker- oder Dispatcher-Tab, ohne dass ein python3 auf dem PATH liegt, und eine WSL-Session nur dort, wo wirklich WSL laeuft (GitHub-Issue BeMuCa/jaira#79)"
context: |-
  Was heute falsch ist (jaira 0.4.1, Issue #79 von sashasoft90): auf nativem Windows startet spawn.sh keinen Tab. Zwei unabhaengige Fehler.
  FEHLER 1 python3 fest verdrahtet: core/role/builtin/jaira-dispatcher/scripts/spawn.sh (~125, 138, 141, 171: Pane-ID, Chat-Modell, Brief, Pane-Zustand) und run-lane.sh (~75, 86, 135, 147) parsen JSON mit 'python3 -c'. Auf Windows ist python3 der Store-Stub: druckt 'Python', Exit 49 -> Pane-ID fehlt, spawn.sh stirbt, ein leerer Tab bleibt stehen.
  Viele Windows-Setups haben gar kein python3, nur uv ('uv run python'). Ein python3-Shim auf dem PATH ist laut Issue KEINE akzeptable Loesung.
  FEHLER 2 spawn.sh:~161 'case $herdr in /mnt/*|*.exe)' haelt jedes herdr.exe fuer WSL. Nativ ist HERDR_BIN_PATH C:\Users\<u>\.herdr\...\herdr.exe -> Tab startet 'wsl.exe --cd /j/REPO/...', WSL kann dort nicht hin (chdir failed 2), claude startet in / und haengt am Trust-Dialog.
  Entschieden (Teamlead, 13.10.):
  (a) In den Skripten gar kein Python mehr fuer die einfachen Faelle: es sind einzelne JSON-Felder aus herdr-Ausgaben. Ein kleiner Unterbefehl im Binary, z. B. 'jaira json get <pfad>' (liest JSON von stdin, druckt ein Feld), ersetzt die python3-Aufrufe - jaira ist auf jeder Maschine, die die Rollen hat, Python nicht. Kein jq (auch nicht ueberall installiert).
  (b) Wo eine Rolle wirklich ein Python-Skript laufen laesst (jaira-role-acceptance build.py, jaira-role-logbook-summary collect.py/build.py; SKILL.md acceptance ~159/187/237, logbook-summary ~44/116): das Python-Kommando ist einstellbar - Einstellung 'python' in ~/.jaira/settings.json und Env JAIRA_PYTHON (Env gewinnt), Default python3. Die SKILL.md-Texte nennen den Weg, ihn zu lesen, statt python3 zu schreiben. Hinweis aus dem Issue: Windows-Konsolen auf cp1251/cp1252 brechen an Umlauten - die Skripte setzen stdout selbst auf UTF-8 (sys.stdout.reconfigure), nicht der Aufrufer per -X utf8.
  (c) WSL an WSL erkennen, nicht an .exe: WSL_DISTRO_NAME gesetzt oder /proc/version enthaelt microsoft. Unter Git Bash/MSYS (uname -s MINGW*/MSYS*) startet der Tab nativ: cd in den Worktree, dann claude; Pfad per cygpath -w, falls herdr ihn braucht.
  Repro (Issue): Windows 11, Git Bash, Herdr fuer Windows, HERDR_ENV=1, Store-Alias an -> 'jaira roles install --global' -> 'spawn.sh --parent <name> <slug> <ID> dispatch' endet in 'Python' plus leerem Tab.
  Hier auf WSL nicht vollstaendig nachstellbar: Fall WSL muss weiter gehen (Alex arbeitet so, herdr.exe unter /mnt/c); den nativen Fall per Stub/Test absichern und Alex fuer den Handcheck auf Windows benennen.
  Settings-Vorbild: core/settings/settings.go Dispatch (~102-109). Rollen kommen per go:embed aus core/role/builtin/.
definition-of-done: "spawn.sh und run-lane.sh rufen kein python3 mehr auf; JSON-Felder holt ein jaira-Unterbefehl, getestet"
tags:
  - cli
blocked-by: []
related: []
commits:
  - d9a5f8646dcb33f972d2be476584788c52938633
  - 3b385ef850e8d909f6d34abc9198578112f5bc03
  - 15f125972112c8faa3054cda4fc96e07d63e3fc4
  - b20ff5a637e9f829d9e72bf50bfb4cc7fc626ca8
  - 469b9b57ed1cf44391e18b9144360932976dd010
  - 4f1f517e7ee61d7185f73d5d2e78fe83e045b305
created-at: 2026-10-13T19:49:16Z
updated-at: 2026-10-13T20:49:58Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-6892
claimed-at: 2026-10-13T20:23:02Z
outcome-what: "spawn.sh/run-lane.sh lesen JSON mit dem neuen 'jaira json get' statt python3; der Chat-Brief kommt aus 'jaira show --brief'; der WSL-Start haengt an WSL_DISTRO_NAME//proc/version statt an .exe, unter Git Bash startet claude nativ im per cygpath -w uebersetzten Worktree; Python-Rollen nehmen JAIRA_PYTHON/'python'-Einstellung (roles list --json: python) und schreiben UTF-8"
outcome-why: Auf nativem Windows ist python3 fehlend oder der Store-Stub und herdr.exe wurde fuer WSL gehalten - spawn.sh liess leere Tabs stehen
outcome-resolves: "GitHub-Issue BeMuCa/jaira#79"
review-summary: "Sechs Commits. spawn.sh und run-lane.sh lesen jedes JSON-Feld mit dem neuen 'jaira json get <pfad> [--default v]' (internal/cli/jsonget.go: Pfad mit Punkten, String roh, null leer, sonst kompaktes JSON ohne \\u-Escapes, fehlt -> Exit 5) statt python3 -c. Der Chat-Brief fuer explain/discuss kommt aus 'jaira show <id> --brief' (tickets.go ticketBrief, trennt Notizen nur an '- **Datum'). spawn.sh prueft direkt nach dem HERDR_ENV-Check, ob das jaira auf PATH 'json get' kann, und bricht sonst vor Worktree und Tab ab. start_on() entscheidet einmal: MINGW/MSYS/CYGWIN -> windows (Tab per --cwd im cygpath -w-Pfad, Kommando nur 'claude [--model x]', unabhaengig von cmd/PowerShell/bash im Tab); herdr.exe oder /mnt/* plus WSL_DISTRO_NAME oder microsoft in /proc/version -> wsl.exe --cd wie bisher; sonst 'cd ... && claude'. Python-Rollen (acceptance, logbook-summary) holen das Kommando aus 'jaira roles list --json' (Feld python: JAIRA_PYTHON > settings 'python' > python3) und setzen stdout/stderr selbst auf UTF-8. Release-Note core/release/notes/TP44SC.md, fuenf Zeilen, die erste sagt 'erst Binary, dann roles install'."
round: 4
review-gaps: "(1) Der Kern von Issue #79 ist weiter nicht auf Windows gelaufen: Git-Bash-Start, --cwd mit C:\\-Pfad in Herdr fuer Windows und 'jaira roles list --json | jaira json get python' unter uv sind nur per Stub (TestSpawnStartsClaudeWhereItsShellRuns) bzw. auf Linux geprueft. Der Handcheck aus der Note vom 13.10. 20:20 steht aus. (2) docs/COMMANDS.md nennt weder 'jaira json get' (passt unter 'Agent plumbing') noch 'show --brief' neben '--notes-last' in Zeile 58; die Release-Note nennt beide, die Befehlsreferenz nicht. (3) Nebensache: --brief zaehlt Text zwischen '## Progress' und der ersten datierten Notiz nicht mehr mit, das alte Python tat es; auf echten Tickets steht dort nichts. Keine Defekte im Diff gefunden: Tests gruen (internal/cli, core/role, core/settings, go vet sauber), kein python3 mehr in den Skripten, json get von Hand geprueft (Pfad, Exit 5, --default, null, Umlaute), --brief auf diesem Ticket = 9 Notizen, die mehrzeilige 20:20-Notiz bleibt ganz."
test-verdict: "pass: go test -race -count=1 ./... RC=0 nach Luecken-Fix; DoD 1-8 im Baum geprueft; altes jaira -> spawn.sh RC=1 mit Meldung, kein Tab (live); WSL-Tab live nach letztem Stand (w3:pGB, dir=Worktree); --brief am echten Ticket. Nicht pruefbar hier: natives Windows - Schritte in Note"
review-verdict: "Der Diff erfuellt die DoD 1-8 fuer WSL und Linux, ohne Defekte, die ich finden konnte. Was Issue #79 wirklich belegt, der Start auf nativem Windows, ist ungeprueft. Freigabe empfohlen, sobald der Windows-Handcheck gelaufen ist; vorher sollte 'Fixes #79' das Issue nicht schliessen. Lueckenpunkt 2 (COMMANDS.md) ist eine Zeile und kann vor dem PR mit."
review-round: 4
review-check: "1. cd ~/projects/.worktrees/jaira-tp44sc && go test ./internal/cli ./core/role ./core/settings  -> drei Zeilen 'ok'. 2. go build -o /tmp/j ./cmd/jaira && echo '{\"result\":{\"pane\":{\"pane_id\":\"w1:p2\"}}}' | /tmp/j json get result.pane.pane_id  -> w1:p2. 3. echo '{}' | /tmp/j json get x; echo $?  -> Meldung 'no key \"x\"', dann 5. 4. echo '{}' | /tmp/j json get x --default -  -> '-'. 5. /tmp/j show TP44SC --brief  -> eine JSON-Zeile mit id, title, goal, question, review, notes_last (4 Eintraege), notes_total 9 oder mehr. 6. JAIRA_PYTHON='uv run python' /tmp/j roles list | tail -1  -> 'python: uv run python ... (from JAIRA_PYTHON)'. 7. Windows-Handcheck, Schritte 0-6 aus der Note vom 13.10. 20:20: neuer Tab steht in C:\\<pfad>, claude laeuft, kein 'Python' in der Ausgabe; mit altem jaira -> 'install the new jaira first', Exit 1, kein Tab. Nur Schritt 7 belegt Issue #79."
---

# Rollen-Skripte starten Tabs auch auf nativem Windows (Git Bash): kein python3, WSL nicht an .exe erkannt

## Definition of Done

- [x] spawn.sh und run-lane.sh rufen kein python3 mehr auf; JSON-Felder holt ein jaira-Unterbefehl, getestet
  proof: internal/cli/jsonget.go newJSONGetCmd + show --brief (tickets.go ticketBrief); TestJSONGetPrintsOneField, TestTicketBriefCarriesTheNewestNotes, core/role TestSpawn* gegen das gebaute Binary (jairabin_test.go)
- [x] Rollen, die ein Python-Skript laufen lassen, nehmen das Kommando aus JAIRA_PYTHON bzw. der Einstellung 'python' (Default python3); die Skripte schreiben stdout als UTF-8
  proof: core/settings/settings.go PythonCommand; roles list python key (TestRolesListReportsThePythonCommand); SKILL.md lines ask jaira; build.py/collect.py reconfigure stdout (TestRolesNeverRunPython3ByName)
- [x] spawn.sh nimmt den WSL-Weg nur bei WSL_DISTRO_NAME oder microsoft in /proc/version; unter MINGW/MSYS startet es claude nativ im Worktree; Test mit Stub fuer beide Faelle
  proof: spawn.sh start_on()/on; core/role/role_test.go TestSpawnStartsClaudeWhereItsShellRuns (wsl + git bash stubs)
- [x] Auf dieser WSL-Maschine startet spawn.sh weiter einen Worker-Tab wie bisher (live geprueft)
  proof: live 13.10. zweimal, zuletzt nach Luecken-Fix: Tab w3:pGA, wsl=Ubuntu-24.04, dir=~/projects/.worktrees/jaira-tp44sc, claude working; Tab geschlossen
- [x] Release-Note als Datei core/release/notes/<Handle>.md; der Pull Request nennt 'Fixes #79'
  proof: core/release/notes/TP44SC.md; 'Fixes #79' steht im Commit-Body und gehoert in die PR-Beschreibung
- [x] Der native Tab-Start haengt nicht von der Shell im Herdr-Tab ab (cwd ueber herdr oder expliziter bash-Start); Note nennt den Weg
  proof: spawn.sh start_on()=windows -> --cwd (cygpath -w) + nacktes claude; TestSpawnStartsClaudeWhereItsShellRuns (git bash: Tab laeuft 'claude', --cwd C:\work\repo); Weg in Note 13.10.
- [x] spawn.sh prueft vor 'tab create', ob das jaira auf dem PATH 'json get' kann, und bricht sonst ohne leeren Tab mit Meldung ab; die Release-Note sagt 'erst neues Binary, dann roles install'
  proof: spawn.sh Check vor tab create; TestSpawnStopsBeforeATabWhenJairaIsTooOld; live: altes jaira -> RC=1 mit Meldung, kein Tab; core/release/notes/TP44SC.md Zeile 1
- [x] --brief trennt Notizen nur an '- **Datum'; Folgezeilen mit '#' oder '- ' bleiben in ihrer Notiz; getestet
  proof: internal/cli/tickets.go briefNote/ticketBrief; TestTicketBriefCarriesTheNewestNotes (#79- und '- '-Folgezeile)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-13 19:49 · Alexander Sacharov** — parent session: jaira-6a
- **2026-10-13 19:50 · Alexander Sacharov** — decision (worker, keine offene Frage fuer die Person): (1) Unterbefehl 'jaira json get <pfad> [--default v]': liest JSON von stdin, Pfad mit Punkten (result.pane.pane_id), String roh, Objekt/Array kompakt als JSON ohne \u-Escapes; fehlt der Pfad und kein --default -> Exit 1. (2) Der Brief fuer explain/discuss war eine Python-Umformung, kein Feld: er kommt jetzt aus 'jaira show <id> --brief' (Go, gleiche Felder title/goal/question/review/notes_last(4)/notes_total). (3) Python-Kommando: 'jaira roles list --json' traegt 'python' (JAIRA_PYTHON > settings 'python' > python3); die SKILL.md lesen es so und rufen $py unquoted auf, damit 'uv run python' geht. (4) Script-Tests bauen das echte Binary einmal per go build, statt python3 vorauszusetzen.
- **2026-10-13 19:58 · Alexander Sacharov** — Live-Check DoD4: Lane-Name 'livecheck-noop' benutzt, damit der gestartete claude nichts Echtes tut; Tab sofort geschlossen, Ticket unveraendert. Achtung Deployment: die neuen Skripte brauchen ein jaira mit 'json get' auf dem PATH - roles install aus demselben Binary liefert das automatisch, ein altes jaira + neue Skripte nicht. Native Windows (Git Bash) nur per Stub getestet: Handcheck durch Alex/sashasoft90 auf Windows noetig. Fuer den PR: 'Fixes #79' in die Beschreibung.
- **2026-10-13 19:58 · Alexander Sacharov** — critique runde 1: spawn.sh:~114 start_on() wird zweimal aufgerufen (cwd und start); einmal in on=$(start_on) halten
internal/cli/tickets.go ticketBrief parst Notes mit eigener Regex neben trimNotes/progressNotes (tickets.go, resume.go) - eine Note muss dasselbe heissen wie dort: Zeile mit "- " unter ## Progress, Folgezeilen gehoeren dazu
internal/cli/jsonget.go Long-Text: auch Eingabe ohne JSON endet in Exit 5 / --default; so sagen
- **2026-10-13 20:08 · Alexander Sacharov** — review (fresh): DoD auf WSL/Linux erfuellt, Tests gruen. Offen vor 'Fixes #79': Handcheck auf nativem Windows - welche Shell der Herdr-Tab startet (cmd/PS5.1 koennen "cd '...' && claude" nicht). Altes jaira auf PATH + neue Skripte = leerer Tab wie vorher; Release-Note koennte 'erst jaira aktualisieren' sagen. Details in review-gaps.
- **2026-10-13 20:11 · Alexander Sacharov** — Alex am 13.10. (aus signoff zurueck): review-gaps 1-3 beheben, danach Handcheck auf nativem Windows. (1) Der native Start darf nicht von der Shell im Herdr-Tab abhaengen (cmd, PowerShell 5.1 kennen kein &&): Arbeitsverzeichnis ueber herdr selbst setzen, falls 'tab create'/'pane run' ein cwd kennt (herdr --skill / --help pruefen), sonst explizit ueber Git Bash starten (bash.exe -lc). Welcher Weg und warum: in eine Note. (2) Release-Note sagt: erst das neue Binary installieren, dann 'jaira roles install' - und spawn.sh prueft frueh, VOR 'tab create', ob 'jaira json' da ist, und bricht mit klarer Meldung ab statt einen leeren Tab zu hinterlassen. (3) --brief: eine Notiz endet nur an der naechsten '- **Datum'-Zeile wie im alten Python; Folgezeilen mit '#' oder '- ' gehoeren zur Notiz; Test dafuer. An denselben Worker, Cache warm.
- **2026-10-13 20:12 · Alexander Sacharov** — Weg fuer DoD6 (Luecke 1): 'herdr tab create --help' kennt --cwd <PATH>; auf Windows (Git Bash/MSYS/Cygwin) bekommt der Tab das Verzeichnis nur ueber --cwd (cygpath -w) und als Kommando nur 'claude [--model x]' ohne cd, ohne &&, ohne Quotes - das versteht cmd, PowerShell 5.1/7 und bash gleich. Nicht bash.exe -lc: der Pfad zu bash.exe ist auf jeder Maschine anders und aus der Tab-Shell heraus nicht sicher auffindbar. Linux/macOS behalten 'cd ... && claude', WSL behaelt wsl.exe --cd. Risiko fuer den Handcheck: ob Herdr fuer Windows --cwd mit C:\-Pfad wirklich setzt. Luecke 3 (--brief) folgt Alex' Vorgabe und kehrt damit Critique-Befund 1/2 (gleich wie trimNotes zaehlen) um: getrennt wird nur an '- **Datum'.
- **2026-10-13 20:15 · Alexander Sacharov** — critique runde 3 (nach signoff-Ruecklauf): spawn.sh: der Check auf 'jaira json get' sitzt nach 'git worktree add' - ein altes jaira laesst dann einen leeren Worktree samt Branch stehen; direkt hinter den HERDR_ENV-Check ziehen, vor alles, was etwas anlegt
- **2026-10-13 20:20 · Alexander Sacharov** — Handcheck Windows (fuer review-check Schritt 6), konkret fuer Alex:
0) in WSL bauen: cd ~/projects/.worktrees/jaira-tp44sc && GOOS=windows GOARCH=amd64 go build -o /mnt/c/Users/Alex/jaira-tp44sc/jaira.exe ./cmd/jaira
1) in Herdr fuer Windows einen Git-Bash-Tab oeffnen (HERDR_ENV=1), dann: export PATH="/c/Users/Alex/jaira-tp44sc:$PATH"; jaira --version; echo "{\"a\":\"ä\"}" | jaira json get a   -> ä
2) jaira roles install --into /c/Users/Alex/jaira-tp44sc/skills   (nicht --global: das ueberschreibt die installierten Windows-Rollen)
3) cd in einen nativen Windows-Klon mit .jaira (z. B. /c/<pfad>/jaira), ID=<ein Ticket dort>
4) bash /c/Users/Alex/jaira-tp44sc/skills/jaira-dispatcher/scripts/spawn.sh --no-worktree livecheck $ID livecheck-noop "$PWD"
   erwartet: druckt eine Pane-ID, kein "Python"; neuer Tab steht in C:\<pfad>\jaira (nicht in / oder C:\Users\Alex), claude laeuft dort; bei Trust-Dialog meldet spawn.sh "approval dialog" - dann im Tab pruefen, dass er nach dem richtigen Ordner fragt
5) Tab schliessen: "$HERDR_BIN_PATH" tab close <tab>  (oder per Maus)
6) Gegenprobe altes jaira: PATH ohne Schritt-1-Ordner, Schritt 4 wiederholen -> Meldung "install the new jaira first ...", Exit 1, kein Tab
Branch feat/tp44sc ist lokal, nicht gepusht - Schritt 0 laeuft im WSL-Worktree.
- **2026-10-13 20:23 · Alexander Sacharov** — review (fresh, Runde 3): keine Defekte im Diff, Tests gruen. Offen: Windows-Handcheck (Note 20:20) vor 'Fixes #79'; docs/COMMANDS.md nennt 'json get' und 'show --brief' nicht. Details in review-gaps.
