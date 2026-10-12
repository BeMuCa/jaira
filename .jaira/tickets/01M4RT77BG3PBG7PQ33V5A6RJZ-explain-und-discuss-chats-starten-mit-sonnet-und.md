---
id: 01M4RT77BG3PBG7PQ33V5A6RJZ
title: "Explain- und Discuss-Chats starten mit sonnet und lesen nur, was sie brauchen"
status: critique
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
updated-at: 2026-10-12T19:51:38Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-65469
claimed-at: 2026-10-12T19:43:51Z
outcome-what: "Neues Setting chat-model (~/.jaira/settings.json, Settings-Screen S 'Chat model', jaira roles list --json chat_model); spawn.sh startet explain/discuss mit claude --model <wert> und tippt nur title/goal/question/review/letzte 4 Notes als JSON in den Prompt"
outcome-why: "Ein Explain-Chat lud ~56k Kontext auf Opus fuer 4-5 Zeilen; Alex: sonnet, aber ueber die Einstellung"
outcome-resolves: "DoD 1-3: settings.go ChatModelName, tui settings.go stepChatModel, roles.go chat_model, spawn.sh brief + --model; Tests in settings, cli, tui, role; Messung 31741->4404 Zeichen"
---

# Explain- und Discuss-Chats starten mit sonnet und lesen nur, was sie brauchen

## Definition of Done

- [x] Das Modell für explain- und discuss-Chats kommt aus einer Projekt-Einstellung (in den jaira-Settings, auch auf dem Settings-Screen S sichtbar), z. B. chat-model; spawn.sh startet dann 'claude --model <wert>'. Ohne Einstellung bleibt es beim Konto-Default wie heute, nichts ist fest verdrahtet; andere Lanes unverändert; Test
  proof: core/settings/settings.go:117 chat-model, :128 ChatModelName; internal/cli/roles.go:67 chat_model; internal/tui/settings.go:62 Chat model row; core/role/builtin/jaira-dispatcher/scripts/spawn.sh:154 --model nur bei explain/discuss; tests TestChatModelName, TestRolesListReportsTheChatModel, TestSettingsEnterStepsTheChatModel, TestSpawnChatsStartOnTheChatModelWithTheirBrief (review lane bleibt plain claude)
- [x] Der Explain-Prompt liest nur review.* und die letzten Notes (z. B. jaira show --json | Auszug), nicht den ganzen Ticket-JSON; der Startkontext eines Explain-Chats ist messbar kleiner (vorher/nachher in einer Note)
  proof: core/role/builtin/jaira-dispatcher/scripts/spawn.sh:144 Auszug title/goal/question/review/letzte 4 Notes in den Prompt; test TestSpawnChatsStartOnTheChatModelWithTheirBrief (core/role/role_test.go:495); Messung vorher/nachher in Note (31741 -> 4404 Zeichen bei 13PQB3)
- [x] Zeile in core/release/notes/<ID>.md
  proof: core/release/notes/5A6RJZ.md

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
