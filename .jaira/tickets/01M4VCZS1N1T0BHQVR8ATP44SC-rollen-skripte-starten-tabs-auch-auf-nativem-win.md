---
id: 01M4VCZS1N1T0BHQVR8ATP44SC
title: "Rollen-Skripte starten Tabs auch auf nativem Windows (Git Bash): kein python3, WSL nicht an .exe erkannt"
status: critique
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
commits: []
created-at: 2026-10-13T19:49:16Z
updated-at: 2026-10-13T20:16:07Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-427304
claimed-at: 2026-10-13T19:49:40Z
outcome-what: "spawn.sh/run-lane.sh lesen JSON mit dem neuen 'jaira json get' statt python3; der Chat-Brief kommt aus 'jaira show --brief'; der WSL-Start haengt an WSL_DISTRO_NAME//proc/version statt an .exe, unter Git Bash startet claude nativ im per cygpath -w uebersetzten Worktree; Python-Rollen nehmen JAIRA_PYTHON/'python'-Einstellung (roles list --json: python) und schreiben UTF-8"
outcome-why: Auf nativem Windows ist python3 fehlend oder der Store-Stub und herdr.exe wurde fuer WSL gehalten - spawn.sh liess leere Tabs stehen
outcome-resolves: "GitHub-Issue BeMuCa/jaira#79"
review-summary: "spawn.sh: der Check auf 'jaira json get' sitzt nach 'git worktree add' - ein altes jaira laesst dann einen leeren Worktree samt Branch stehen; direkt hinter den HERDR_ENV-Check ziehen, vor alles, was etwas anlegt"
round: 4
review-gaps: "(1) Nativer Windows-Pfad nur per Stub geprueft: der Tab bekommt \"cd 'C:\\...' && claude\" - das geht in Git Bash und PowerShell 7, nicht in cmd.exe und nicht in Windows PowerShell 5.1 (kein &&); welche Shell Herdr fuer Windows im Tab startet, ist nirgends belegt. Handcheck auf Windows steht aus (Notiz 19:58 sagt das selbst). (2) Ein altes jaira vorne auf dem PATH plus neue Skripte: 'jaira json' ist dort unbekannt, spawn.sh stirbt nach 'tab create' und laesst wieder einen leeren Tab stehen - die Release-Note sagt nur 'roles install again', nicht 'mit dem neuen Binary'. (3) Klein: --brief erbt von progressEntries, dass eine Notiz-Folgezeile, die mit '#' beginnt (z.B. '#79 ...'), den Progress-Abschnitt beendet und alle spaeteren Notizen verschluckt; und jede '- '-Folgezeile zaehlt als eigene Notiz, wo das alte Python nur '- **Datum' trennte. (4) SKILL.md: $(...) unquoted zerlegt einen Python-Pfad mit Leerzeichen (C:\\Program Files\\...) - als 'Kommando, kein Pfad' dokumentiert, also bewusst. Keine Defekte im WSL-/Linux-Pfad gefunden."
test-verdict: "pass: go test -race -count=1 ./... RC=0; DoD 1-5 im Baum geprueft; run-lane.sh mit herdr/jaira-Stubs durchgespielt (Status gelesen, Warten, Zusammenfassung mit Umlaut, Tab per tab_id geschlossen, RC=0); JAIRA_PYTHON='uv run --no-project python' laeuft ueber die SKILL-Zeile; spawn.sh live auf WSL (DoD 4). Nicht pruefbar hier: echtes natives Windows (Git Bash) - nur Stub"
review-verdict: "Erfuellt die DoD fuer WSL und Linux: kein python3 mehr in spawn.sh/run-lane.sh, json get und --brief getestet, WSL-Erkennung richtig und per Stub fuer beide Faelle abgesichert, Tests gruen (internal/cli, core/role, core/settings). Unsicher bleibt genau das, was Issue #79 betrifft: der native Git-Bash-Start ist nicht auf Windows gelaufen. Freigabe empfohlen unter der Bedingung, dass jemand den Handcheck auf Windows macht, bevor 'Fixes #79' das Issue schliesst."
review-round: 2
review-check: "1. cd ~/projects/.worktrees/jaira-tp44sc && go test ./internal/cli ./core/role ./core/settings  -> alle drei 'ok'. 2. go build -o /tmp/j ./cmd/jaira && echo '{\"result\":{\"pane\":{\"pane_id\":\"w1:p2\"}}}' | /tmp/j json get result.pane.pane_id  -> druckt w1:p2. 3. echo '{}' | /tmp/j json get x; echo $?  -> Fehlermeldung, Exit 5; mit --default - druckt '-'. 4. /tmp/j show TP44SC --brief  -> eine JSON-Zeile mit id, title, goal, question, review, notes_last (max 4), notes_total. 5. JAIRA_PYTHON='uv run python' /tmp/j roles list | tail -1  -> 'python: uv run python ... (from JAIRA_PYTHON)'. 6. Auf nativem Windows (Git Bash, Herdr fuer Windows, kein python3, HERDR_ENV=1): neues jaira installieren, 'jaira roles install --global', dann 'spawn.sh --no-worktree <slug> TP44SC livecheck-noop' -> ein Tab oeffnet im Worktree (C:\\...), claude startet ohne Trust-Dialog in /, kein 'Python' in der Ausgabe. Schritt 6 ist der einzige, der Issue #79 wirklich belegt."
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
