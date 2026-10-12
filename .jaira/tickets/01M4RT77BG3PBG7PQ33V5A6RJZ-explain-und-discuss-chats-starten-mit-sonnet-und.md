---
id: 01M4RT77BG3PBG7PQ33V5A6RJZ
title: "Explain- und Discuss-Chats starten mit sonnet und lesen nur, was sie brauchen"
status: in-progress
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Ein Explain- oder Discuss-Chat kostet nur einen Bruchteil von heute: günstigeres Modell, weniger Kontext"
context: |-
  Ein Explain-Chat im Cockpit lädt ~56k Kontext auf Opus, um 4-5 Zeilen zu schreiben.
  - Gesehen von Alex 2026-10-12 im Herdr-Tab eines Explain-Chats.
  - Ursache Modell: core/role/builtin/jaira-dispatcher/scripts/spawn.sh (~Zeile 130) startet 'claude' ohne --model, also das Konto-Default (Opus 5.5) - auch für die Lanes explain und discuss.
  - Ursache Kontext: ~20-30k sind die Grundlast jeder Claude-Code-Sitzung (nicht vermeidbar), ~8k CLAUDE.md+AGENTS.md, ~10-20k 'jaira show <id> --json' mit allen Notes und DoD-Beweisen. Der Explain-Prompt (spawn.sh ~Zeile 171) lässt den ganzen Ticket-JSON lesen, gebraucht werden nur review.* und die letzten Notes.
  - Entscheidung Alex 2026-10-12: sonnet für explain und discuss.
  - Nicht anfassen: worker/dispatch/role-lane behalten das Default bzw. die Lane-Tiers.
  - Erst nach 13PQB3 anfangen: 13PQB3 ändert denselben Explain-Prompt in spawn.sh.
definition-of-done: "spawn.sh startet für die Lanes explain und discuss 'claude --model sonnet'; andere Lanes unverändert; Test"
tags: []
blocked-by: []
related: []
commits: []
created-at: 2026-10-12T19:42:48Z
updated-at: 2026-10-12T20:08:02Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-65469
claimed-at: 2026-10-12T19:43:51Z
outcome-what: "Setting chat-model: Board-Datei .jaira/chat-model zuerst, dann ~/.jaira/settings.json, sonst Konto-Default (core/settings ChatModelFor); Settings-Screen S 'Chat model' schaltet den Board-Wert und zeigt die Quelle; jaira roles list --json chat_model/chat_model_source; spawn.sh startet explain/discuss mit claude --model <wert> und tippt nur title/goal/question/review/letzte 4 Notes als JSON in den Prompt"
outcome-why: "Ein Explain-Chat lud ~56k Kontext auf Opus fuer 4-5 Zeilen; Alex: sonnet ueber eine Projekt-Einstellung mit Rechner-Fallback"
outcome-resolves: "DoD 1-4: settings.go ChatModelFor/BoardChatModel, tui settings.go stepChatModel, roles.go chat_model, spawn.sh brief + --model; Tests in settings, cli, tui, role; Messung 31741->4404 Zeichen"
review-summary: "core/settings/settings.go:188 the empty case's source is 'account default, nothing configured', so S shows 'account default — account default, nothing configured'; make the source 'nothing configured'. internal/tui/settings.go:116 reads BoardChatModel once per choice inside the loop; read it once before the loop"
review-gaps: "1) Die WSL-Startzeile ist nicht automatisch getestet: der Test deckt nur den Zweig 'cd … && claude --model …' ab. Auf Alex' Rechner (Herdr ist eine Windows-Binary) läuft aber der andere Zweig: wsl.exe --cd '…' -- bash -lic \"claude --model 'sonnet'\". Ob die doppelten Anführungszeichen durch die Windows-Shell und wsl.exe bis zu bash kommen, hat niemand mit echtem Herdr geprüft (testing nutzte einen Stub). Geht das schief, startet der Chat still auf dem Konto-Default oder gar nicht. 2) DoD 1 sagt 'Projekt-Einstellung'; gebaut ist eine Einstellung pro Rechner (~/.jaira/settings.json, wie dispatch). Die Teamlead-Note deckt das (core/settings), aber Alex sollte wissen: sonnet gilt dann für alle Repos auf dem Rechner. Sonst keine Defekte gefunden; go test für settings, role, cli, tui grün."
test-verdict: "pass: go build/vet + go test -race -count=1 ./... RC=0; DoD 1-3 verified in the tree; built binary: roles list shows 'chat model: sonnet' and chat_model in --json; spawn.sh with stub Herdr + real jaira on 13PQB3: explain starts claude --model 'sonnet', review lane plain claude, empty setting plain claude; explain prompt 5650 chars incl. ticket excerpt (was: prompt + agent reading 31741 chars of show --json)"
review-verdict: "Der Diff erfüllt alle drei DoD-Punkte, keine Defekte im Code gefunden. Vorbehalt: der WSL/Windows-Herdr-Startweg - der, den Alex tatsächlich benutzt - ist nur per Stub geprüft; vor dem Annehmen einmal echt ausprobieren (review-check Schritt 3-5). Und die Einstellung gilt pro Rechner, nicht pro Projekt."
review-check: "1. Im Repo: go build -o /tmp/jaira ./cmd/jaira  2. /tmp/jaira roles list → letzte Zeile lautet 'chat model: account default …'  3. /tmp/jaira öffnen, S drücken, 4x j, enter → Zeile 'Chat model  sonnet'; ~/.jaira/settings.json enthält jetzt \"chat-model\": \"sonnet\"  4. Mit dieser Binary auf PATH im Cockpit bei einem Ticket mit Review-Lücke Explain drücken → im neuen Herdr-Tab startet claude; dort /status eingeben → Modell ist Sonnet, nicht Opus  5. Im selben Tab: der erste Prompt endet mit '… als JSON; lies das Ticket nicht noch einmal ganz: {\"id\": …' und der Chat ruft nicht 'jaira show --json' auf  6. Settings S zurück auf 'account default' stellen, Explain erneut → /status zeigt wieder das Konto-Default. Automatisch: go test ./core/role -run TestSpawnChatsStartOnTheChatModelWithTheirBrief"
---

# Explain- und Discuss-Chats starten mit sonnet und lesen nur, was sie brauchen

## Definition of Done

- [x] Das Modell für explain- und discuss-Chats kommt aus einer Projekt-Einstellung (in den jaira-Settings, auch auf dem Settings-Screen S sichtbar), z. B. chat-model; spawn.sh startet dann 'claude --model <wert>'. Ohne Einstellung bleibt es beim Konto-Default wie heute, nichts ist fest verdrahtet; andere Lanes unverändert; Test
  proof: core/settings/settings.go:179 ChatModelFor (Board, dann settings.json:121); internal/cli/roles.go:59 chat_model; internal/tui/settings.go:62 Chat model row; core/role/builtin/jaira-dispatcher/scripts/spawn.sh:149 --model nur bei explain/discuss; tests TestChatModelName, TestRolesListReportsTheDispatchSetting, TestSettingsEnterStepsTheChatModel, TestSpawnChatsStartOnTheChatModelWithTheirBrief (review lane bleibt plain claude)
- [x] Der Explain-Prompt liest nur review.* und die letzten Notes (z. B. jaira show --json | Auszug), nicht den ganzen Ticket-JSON; der Startkontext eines Explain-Chats ist messbar kleiner (vorher/nachher in einer Note)
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:139 Auszug title/goal/question/review/letzte 4 Notes in den Prompt; test TestSpawnChatsStartOnTheChatModelWithTheirBrief (core/role/role_test.go:495); Messung vorher/nachher in Note (31741 -> 4404 Zeichen bei 13PQB3)
- [x] Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/5A6RJZ.md
- [x] chat-model wird zuerst aus der Board-Einstellung im Repo (.jaira/, mit dem Board geteilt) gelesen, erst wenn dort leer aus ~/.jaira/settings.json; Settings-Screen S zeigt, woher der Wert kommt; Test
  proof: core/settings/settings.go:148 BoardChatModel (.jaira/chat-model), :179 ChatModelFor Board vor settings.json mit Quelle; internal/tui/settings.go:109 stepChatModel schreibt den Board-Wert, Zeile zeigt Quelle; internal/cli/roles.go:59 bestEffortRoot; spawn.sh:135 roles list -C $wt; tests TestChatModelForPrefersTheBoard (core/settings/settings_test.go:247), TestSettingsEnterStepsTheChatModel, TestRolesListReportsTheDispatchSetting (board wins)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

- [x] core/settings: chat-model + ChatModelName() (validiert), Test
- [x] jaira roles list: chat_model in --json und Textzeile
- [x] TUI Settings-Screen S: Zeile 'Chat model', enter schaltet reihum, Test
- [x] spawn.sh: --model für explain/discuss aus roles list; Ticketauszug (title, goal, question, review, letzte Notes) in den Prompt; Test mit jaira-Stub; vorher/nachher messen
- [x] core/release/notes/5A6RJZ.md

## Progress
- **2026-10-12 19:43 · Alexander Sacharov** — Alex 2026-10-12: sonnet ja, aber nur über die Projekt-Einstellungen, nicht fest in spawn.sh. Alex setzt in diesem Projekt sonnet. DoD 1 umformuliert.
- **2026-10-12 19:43 · Alexander Sacharov** — Teamlead 2026-10-12: Branch feat/5a6rjz von feat/13pqb3 (353ac57d) abgezweigt: 13PQB3 hat den Explain-Prompt in spawn.sh (Optionen mit [give back]/[accept]) umgebaut; darauf aufsetzen. 13PQB3 ist noch im review - kommt es zurück, auf dessen Ende rebasen. Die Projekt-Einstellung gehört zu den übrigen jaira-Settings (core/settings), sichtbar auf dem Settings-Screen S.
- **2026-10-12 19:43 · Alexander Sacharov** — parent session: jaira-f5
- **2026-10-12 19:45 · Alexander Sacharov** — decisions (worker, settled without asking): (1) Einstellung = chat-model in ~/.jaira/settings.json (core/settings), wie die Teamlead-Note sagt; maschinenweit wie dispatch, nicht pro Repo. (2) Settings-Screen S: Zeile 'Chat model', enter schaltet reihum account default -> sonnet -> haiku -> opus -> default (wie Agents per ticket); jeder andere Wert (z. B. claude-sonnet-5-5) per Hand in settings.json, wird angezeigt. (3) spawn.sh liest den Wert über 'jaira roles list --json' (neues Feld chat_model, neben dispatch); Werte ausser [A-Za-z0-9._[]-] gelten als leer, weil der Wert in eine Shell-Zeile geht. (4) Kontext: spawn.sh holt selbst 'jaira show --json' und tippt nur title, goal, question, review.* und die letzten 4 Notes als kompaktes JSON in den Prompt - der Chat braucht keinen Tool-Call zum Lesen. Ohne jaira/python3 fällt der Prompt auf den alten Lesehinweis zurück. Gilt für explain und discuss (Ziel nennt beide).
- **2026-10-12 19:51 · Alexander Sacharov** — Messung Startkontext (Zeichen, was der Chat vom Ticket liest): vorher 'jaira show --json' = 31741 bei 13PQB3 (28 Notes, 14 DoD mit Beweisen), 7114 bei 5A6RJZ; nachher der Auszug in spawn.sh = 4404 bzw. 2030. Also ~7x weniger auf einem langen Ticket, ~6-7k Tokens gespart, plus kein Tool-Call zum Lesen. Grundlast der Claude-Code-Sitzung (~20-30k) und CLAUDE.md bleiben, wie die Context-Notiz sagt. Dead end: --notes-last wirkt auf --json nicht (beide 31741), daher die Extraktion in spawn.sh statt eines jaira-Flags.
- **2026-10-12 19:52 · Alexander Sacharov** — critique Runde 1: ein Befund - spawn.sh: der neue chat-model/brief-Block steht zwischen dem WSL-Kommentar und dem case, den der erklärt; Block über den WSL-Kommentar ziehen. Bewusst stehen gelassen: Prüfung des Modellnamens doppelt (Go ChatModelName für roles list/TUI, Bash-Regex dort, wo der Wert in die Kommandozeile geht) - die Bash-Prüfung schützt die Stelle, an der es zählt, auch gegen ein jaira, das den Wert anders liefert; Lesen über 'jaira roles list --json' statt eines eigenen Kommandos, weil dispatch schon dort steht und der Teamlead es so liest.
- **2026-10-12 19:53 · Alexander Sacharov** — critique Runde 2: der eine Befund aus Runde 1 ist behoben (e292e8dd, Block über dem WSL-Kommentar). Nichts weiter.
- **2026-10-12 19:57 · Alexander Sacharov** — testing Runde 1: pass - Suite mit -race grün, Funktion mit gebautem Binary und Herdr-Stub geprüft (Startzeile je Lane und Setting, Prompt 13PQB3 5650 Zeichen).
- **2026-10-12 19:59 · Alexander Sacharov** — review (frischer Kontext): DoD erfüllt, keine Code-Defekte. Vorbehalte für den Menschen: WSL-Startzeile mit --model nur per Stub geprüft - einmal echt im Cockpit testen; chat-model gilt pro Rechner, nicht pro Projekt.
- **2026-10-12 20:02 · Alexander Sacharov** — Alex 2026-10-12: Projekt-Einstellung mit Fallback auf den Rechner (DoD 4). Review-Gap 1 (WSL-Startzeile nur per Stub) prüft Alex live, wenn der Ticket wieder in signoff ist.
- **2026-10-12 20:07 · Alexander Sacharov** — Runde 2 (DoD 4) decisions: Board-Wert in eigener Datei .jaira/chat-model (eine Zeile, # Kommentare, wie lanes/order und tags - kein neues key:value-Format für einen Wert). S schaltet den Board-Wert (der gewinnt), nicht settings.json; 'none' löscht die Datei, dann gilt der Rechner-Fallback - den setzt man von Hand. roles list liest mit bestEffortRoot(), spawn.sh ruft es mit -C $wt, ausserhalb eines Boards nur der Rechner. ChatModelName() ist weg, ChatModelFor(root) liefert Wert + Quelle (roles list --json chat_model_source, S zeigt sie).
- **2026-10-12 20:07 · Alexander Sacharov** — critique Runde 3 (DoD 4): zwei Befunde - (1) leere Quelle heisst 'account default, nothing configured', S zeigt dann 'account default — account default, ...' -> Quelle 'nothing configured'. (2) stepChatModel liest .jaira/chat-model in der Schleife je Wahl -> einmal davor. Stehen gelassen: Guard root=='' in stepChatModel (verhindert Schreiben relativ zum cwd), settings-Paket hält jetzt auch eine Board-Datei (wie RemoteFor schon git config jaira.remote liest).
