---
id: 01M4RT77BG3PBG7PQ33V5A6RJZ
title: "Explain- und Discuss-Chats starten mit sonnet und lesen nur, was sie brauchen"
status: done
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
commits:
  - ecfeb117eea7e599b3711f575abc5cb207a2ccba
  - e292e8dd3a65abe05fea176865c99367e96d4d64
  - 83a97f0d832a55fb21e52f69006d6b976f4c7186
  - afc5f89d579ac8a9f2222d5d37dee41a23f12d91
  - 06cc2f32b70d02fecb866d65798efed6a514c4cd
  - 2ce0427a71321154aab44c0db96d290d72c32111
  - 6ef6facbfe5a62300a3f06ad2d479e2d03d625ab
created-at: 2026-10-12T19:42:48Z
updated-at: 2026-10-12T21:03:31Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-83570
claimed-at: 2026-10-12T20:15:31Z
outcome-what: "Setting chat-model: Board-Datei .jaira/chat-model zuerst, dann ~/.jaira/settings.json, sonst Konto-Default (core/settings ChatModelFor); Settings-Screen S 'Chat model' schaltet den Board-Wert, zeigt ihn mit Quelle oder 'none — <was stattdessen gilt>'; ungueltiger Board-Wert = none; jaira roles list --json chat_model/chat_model_source; spawn.sh startet explain/discuss mit claude --model <wert> und tippt nur title/goal/question/review/letzte 4 Notes als JSON in den Prompt"
outcome-why: "Ein Explain-Chat lud ~56k Kontext auf Opus fuer 4-5 Zeilen; Alex: sonnet ueber eine Projekt-Einstellung mit Rechner-Fallback; Review-Kleinigkeiten (none-Label, ungueltiger Wert) auf Alex' Wunsch behoben"
outcome-resolves: "DoD 1-5: settings.go ChatModelFor/BoardChatModel, tui settings.go stepChatModel/chatModelValue, roles.go chat_model, spawn.sh brief + --model; Tests in settings, cli, tui, role; Messung 31741->4404 Zeichen"
review-summary: "Neue Einstellung chat-model: core/settings liest zuerst die Board-Datei .jaira/chat-model (eine Zeile, # Kommentare), dann \"chat-model\" in ~/.jaira/settings.json, sonst nichts (Konto-Default); jeder Wert mit anderen Zeichen als Buchstaben/Ziffern/._-[] gilt als nicht gesetzt. jaira roles list nennt das Modell und seine Quelle (--json: chat_model, chat_model_source). Settings-Screen S hat die Zeile 'Chat model': enter schaltet die Board-Datei reihum sonnet -> haiku -> opus -> none (none loescht die Datei); die Zeile zeigt den Board-Wert mit Quelle oder 'none — <was stattdessen gilt>'; ein ungueltiger Board-Wert zaehlt als none und enter geht von dort auf sonnet. spawn.sh startet nur explain/discuss mit claude --model <wert> (lokal und in der WSL-Startzeile), alle anderen Lanes mit blankem claude; dazu tippt es title/goal/question/review/letzte 4 Notes als einzeiliges JSON in den Prompt statt 'lies jaira show --json' (gemessen 31741 -> 4404 Zeichen bei 13PQB3). Release-Note core/release/notes/5A6RJZ.md mit 4 Zeilen."
review-gaps: "Keine Defekte. Kleinigkeiten, kein Grund zur Rueckgabe: (1) Der WSL-Start mit --model (wsl.exe ... bash -lic \"claude --model 'sonnet'\") ist nur per Herdr-Stub geprueft, nie im echten Windows-Herdr - Alex wollte das live pruefen. (2) spawn.sh: der Notes-Auszug nimmt alles nach '## Progress' als Notes; steht dort einmal ein weiterer Abschnitt dahinter, landet er in der letzten Note - heute ist Progress der letzte Abschnitt, also nur theoretisch. (3) spawn.sh: eine Kommentarzeile im neuen Block (\"than reading the whole of it: 'jaira show --json' carries ...\") ist deutlich laenger umbrochen als der Rest - reine Optik."
test-verdict: "pass (round 3): go vet + go test -race -count=1 ./... RC=0; DoD 1-5 verified; built binary on a scratch board with '$(id)' in .jaira/chat-model: roles list falls through to haiku from settings.json; TestSettingsEnterStepsTheChatModel drives S through the real model: row reads 'none — account default' / 'none — <machine model> from settings.json', invalid value reads none and enter goes to sonnet"
review-verdict: "Alle 5 DoD-Punkte sind im Diff belegt; Tests in settings, role, cli, tui gruen (selbst gelaufen); Board-vor-Rechner und ungueltiger Board-Wert = none mit gebautem Binary selbst gesehen. Annehmen; offen bleibt nur der Live-Start im Windows-Herdr."
review-check: "1. jaira tui im Repo starten, S druecken  2. Zeile 'Chat model' zeigt 'none — account default' (oder 'none — <modell> from settings.json on this machine')  3. enter: Zeile zeigt 'sonnet — from .jaira/chat-model, shared with the board', und cat .jaira/chat-model hat 'sonnet' als letzte Zeile  4. enter noch dreimal: haiku, opus, dann wieder none - .jaira/chat-model ist dann weg  5. echo 'kaputt$(id)' > .jaira/chat-model, S neu oeffnen: Zeile zeigt 'none — ...'; enter -> sonnet  6. jaira roles list --json | jq '{chat_model,chat_model_source}' zeigt dasselbe  7. Im Cockpit an einem Ticket mit Review [Explain] druecken: der neue Herdr-Tab startet mit 'claude --model sonnet', und im Prompt steht ein JSON mit title/goal/review/notes_last statt 'jaira show --json'  8. Eine normale Lane (z. B. review) per spawn.sh starten: Tab startet mit blankem 'claude'"
---

# Explain- und Discuss-Chats starten mit sonnet und lesen nur, was sie brauchen

## Definition of Done

- [x] Das Modell für explain- und discuss-Chats kommt aus einer Projekt-Einstellung (in den jaira-Settings, auch auf dem Settings-Screen S sichtbar), z. B. chat-model; spawn.sh startet dann 'claude --model <wert>'. Ohne Einstellung bleibt es beim Konto-Default wie heute, nichts ist fest verdrahtet; andere Lanes unverändert; Test
  proof: core/settings/settings.go:179 ChatModelFor (Board, dann settings.json:121); internal/cli/roles.go:59 chat_model; internal/tui/settings.go:61 Chat model row; core/role/builtin/jaira-dispatcher/scripts/spawn.sh:149 --model nur bei explain/discuss; tests TestChatModelName, TestRolesListReportsTheDispatchSetting, TestSettingsEnterStepsTheChatModel, TestSpawnChatsStartOnTheChatModelWithTheirBrief (review lane bleibt plain claude)
- [x] Der Explain-Prompt liest nur review.* und die letzten Notes (z. B. jaira show --json | Auszug), nicht den ganzen Ticket-JSON; der Startkontext eines Explain-Chats ist messbar kleiner (vorher/nachher in einer Note)
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:139 Auszug title/goal/question/review/letzte 4 Notes in den Prompt; test TestSpawnChatsStartOnTheChatModelWithTheirBrief (core/role/role_test.go:495); Messung vorher/nachher in Note (31741 -> 4404 Zeichen bei 13PQB3)
- [x] Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/5A6RJZ.md
- [x] chat-model wird zuerst aus der Board-Einstellung im Repo (.jaira/, mit dem Board geteilt) gelesen, erst wenn dort leer aus ~/.jaira/settings.json; Settings-Screen S zeigt, woher der Wert kommt; Test
  proof: core/settings/settings.go:148 BoardChatModel (.jaira/chat-model), :179 ChatModelFor Board vor settings.json mit Quelle; internal/tui/settings.go:134 stepChatModel schreibt den Board-Wert, :117 chatModelValue zeigt Quelle; internal/cli/roles.go:59 bestEffortRoot; spawn.sh:135 roles list -C $wt; tests TestChatModelForPrefersTheBoard (core/settings/settings_test.go:247), TestSettingsEnterStepsTheChatModel, TestRolesListReportsTheDispatchSetting (board wins)
- [x] Ein ungueltiger, von Hand in .jaira/chat-model geschriebener Wert gilt als none (Fallback auf Rechner-Wert/Konto-Default), nicht als sonnet; der letzte Schritt in S heisst 'none'; Kommentar und Release-Note passend; Test
  proof: internal/tui/settings.go:117 chatModelValue ('none — <was gilt>'), :134 stepChatModel Kommentar (ungueltig = none -> sonnet); core/settings/settings.go:148 BoardChatModel (ungueltig -> ''); core/release/notes/5A6RJZ.md Zeile 2; test TestSettingsEnterStepsTheChatModel (internal/tui/settings_test.go:124: none-Label je Schritt, '$(id)' im Board faellt auf settings.json durch, enter -> sonnet)

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
- **2026-10-12 20:09 · Alexander Sacharov** — critique Runde 4: beide Befunde aus Runde 3 behoben. Nichts weiter.
- **2026-10-12 20:13 · Alexander Sacharov** — testing Runde 2: pass - Board vor Rechner auf Scratch-Board mit gebautem Binary gesehen, validate sauber, Datei nicht ignoriert.
- **2026-10-12 20:15 · Alexander Sacharov** — review Runde 2 (frischer Kontext, DoD 4 dabei): alle vier DoD-Punkte im Diff belegt, Tests grün, keine Code-Defekte. Offen nur der Live-Start im Windows-Herdr (Alex prüft) und zwei Formulierungs-Kleinigkeiten (stepChatModel-Kommentar, 'none' in der Release-Note) - kein Grund zurückzugeben.
- **2026-10-12 20:18 · Alexander Sacharov** — Alex 2026-10-12 (signoff): die zwei Review-Kleinigkeiten nicht stehen lassen, sondern fixen -> DoD 5. Zurueck nach in-progress, derselbe Worker.
- **2026-10-12 20:18 · Alexander Sacharov** — Runde 3 (DoD 5) Befund vorab: Lesen behandelt einen ungültigen Board-Wert schon als none (BoardChatModel -> modelName -> '', ChatModelFor fällt auf Rechner/Default durch). Falsch waren der stepChatModel-Kommentar und die Anzeige: S zeigte bei leerem Board den wirksamen Wert ('account default — ...'), nie 'none'. Plan: Zeile zeigt den Board-Wert, bei none 'none — <was stattdessen gilt>'; Enter von none (auch ungültig) -> sonnet, das ist der nächste Schritt nach none.
- **2026-10-12 20:21 · Alexander Sacharov** — critique Runde 5 (DoD 5): nichts. chatBoard neben chatModel ist nötig, weil die Zeile Board-Wert und wirksamen Wert getrennt zeigt; der frische Re-Read in stepChatModel folgt editKey.
- **2026-10-12 20:24 · Alexander Sacharov** — testing Runde 3: pass - ungültiger Board-Wert fällt mit gebautem Binary auf den Rechner-Wert durch; S-Zeile per Test über das echte Model geprüft.
- **2026-10-12 20:27 · Alexander Sacharov** — review Runde 3 (frischer Kontext): DoD 1-5 im Diff belegt, Tests gruen, Board-vor-Rechner und ungueltiger Wert mit gebautem Binary geprueft. Keine Defekte; offen nur der Live-Start im Windows-Herdr.
- **2026-10-12 20:29 · Alexander Sacharov** — Live-Check WSL-Start (Auftrag Teamlead 22:2x): nicht gelaufen. Vorbereitet: Branch-Binary gebaut, Scratch-Board mit .jaira/chat-model=sonnet, 'roles list' sagte 'sonnet — from .jaira/chat-model'. Der Aufruf von spawn.sh (dieser Branch, explain) wurde vom Auto-Mode-Classifier als 'Create Unsafe Agents' abgelehnt; nicht umgangen. Scratch-Board wieder gelöscht, kein Tab geöffnet. Bleibt bei Alex: S -> sonnet, Explain, /status im Tab.
