---
id: 01M4P5G392KHJKK03ZW8T8QS1Y
title: "Cockpit-Rauschfilter: Agentenberichte erscheinen im Teamlead-Chat als eine Zeile"
status: critique
ready: true
creator: Alexander Sacharov
assignee: Alexander Sacharov
goal: "Im Teamlead-Chat steht statt eines vollen Agentenberichts oder einer task-notification eine Zeile wie 'MG677R · review -> signoff · approve mit Vorbehalt'; der volle Text geht nicht verloren und ist mit einem Befehl abrufbar"
context: |-
  Erster Teil von NZZDRH (jaira-cockpit). Was falsch ist: jeder Bericht eines Subagenten und jede task-notification landet als voller Block im Chat der Teamlead-Sitzung; am 11.10.2026 ueber 40 Bloecke fuer 5 Tickets. Der Mensch findet darin nicht, was er entscheiden muss.
  Hebel: Mod-Event 'session.receive' (Claude Code 2.1.295, Mods early access ab 2.1.287). Es sieht jede Lieferung, bevor sie in die Konversation kommt; next({...e, text}) schreibt sie um, {consumed: grund} schluckt sie. Typen: plugin-authoring Skill, types/claude-code.d.ts, SessionReceiveInput/SessionReceiveResult. Herkunft per e.origin unterscheiden.
  WICHTIG - Kompromiss, beim Start mit Alex bestaetigen: die gekuerzte Zeile sieht auch das MODELL des Teamleads, nicht nur der Mensch (umgeschrieben wird, was in die Konversation kommt). Vorschlag: die Zeile traegt Ticket, Lane vorher/nachher, Urteil und eine offene Frage, falls es eine gibt - das reicht, weil der Teamlead ohnehin vom Board arbeitet. Der volle Text wird gespeichert und ist per Slash-Befehl '/flight <id>' abrufbar, auch fuer das Modell.
  Woraus die Zeile entsteht: Ticket-Id und Lane aus 'jaira show <id> --json' bzw. aus dem Bericht; nicht raten, wenn keine Id im Text steht - dann die Lieferung unveraendert durchlassen.
  Was NICHT gefiltert werden darf: Nachrichten des Menschen, Kommentare auf Artifacts, Lieferungen ohne Agentenbezug, und alles, was nach Fehler aussieht (fehlgeschlagener Agent, Abbruch) - das muss voll sichtbar bleiben.
  Entwicklung: als Mod unter ~/.claude/dev-mods/... oder per 'claude --plugin-dir'; Installation ueber jaira regelt NZZDRH.
definition-of-done: "Ein Agentenbericht und eine task-notification mit Ticket-Id erscheinen im Chat als eine Zeile mit Ticket, Lane-Wechsel, Urteil und offener Frage; ohne erkennbare Id bleibt die Lieferung unveraendert; Fehler und Abbrueche bleiben voll"
tags:
  - cli
blocked-by: []
parent: 01M4P5FK818EAFZ68XFYNZZDRH
related: []
commits: []
created-at: 2026-10-11T19:02:10Z
updated-at: 2026-10-11T19:33:09Z
updated-by: Alexander Sacharov
claimed-by: DESKTOP-RFTCH11-767990
claimed-at: 2026-10-11T19:27:17Z
outcome-what: "VERDICT_WORD in core/mod/builtin/jaira-cockpit/hooks/digest.ts: 'mit vorbehalt' als Zusatz der approve-Alternative, damit der Vorbehalt in der Zeile bleibt; Test 'a verdict keeps its reservation'"
outcome-why: "optimize fand: die Alternation nahm 'approve' vor 'approve mit vorbehalt'"
outcome-resolves: "DoD 1 (Urteil vollstaendig in der Zeile)"
review-summary: none
---

# Cockpit-Rauschfilter: Agentenberichte erscheinen im Teamlead-Chat als eine Zeile

## Definition of Done

- [x] Ein Agentenbericht und eine task-notification mit Ticket-Id erscheinen im Chat als eine Zeile mit Ticket, Lane-Wechsel, Urteil und offener Frage; ohne erkennbare Id bleibt die Lieferung unveraendert; Fehler und Abbrueche bleiben voll
  proof: core/mod/builtin/jaira-cockpit/hooks/register.ts:86 (session.receive) + core/mod/builtin/jaira-cockpit/hooks/digest.ts:91; tests 'a report naming a ticket becomes one line', 'a report without a ticket passes unchanged', 'a failed agent stays whole'
- [x] Der volle Text ist gespeichert und mit '/flight <id>' abrufbar (zeigt die Berichte zu einem Ticket, neueste zuerst)
  proof: core/mod/builtin/jaira-cockpit/hooks/register.ts:33 keep() into $.store, :47 flight() for /flight and tool mcp__jaira-cockpit__flight; tests '/flight lists reports newest first', 'the flight tool gives the model the same full text'
- [x] 'claude plugin validate' und 'claude plugin test' gruen; Tests decken: Bericht mit Id, Bericht ohne Id, Fehlerbericht, Nachricht eines Menschen
  proof: claude plugin validate core/mod/builtin/jaira-cockpit: passed; claude plugin test core/mod/builtin/jaira-cockpit: 8 pass 0 fail (core/mod/builtin/jaira-cockpit/tests/cockpit.test.ts)

## Options

- [ ] brainstorm
- [ ] planning

## Plan

<Steps, in order — filled in by the pre-process step, or by you.>

## Progress
- **2026-10-11 19:16 · Alexander Sacharov** — Alex am 11.10.2026 entschieden (Variante A): die gekuerzte Zeile sieht auch das Modell des Teamleads; der volle Bericht bleibt gespeichert und ist per '/flight <id>' abrufbar, auch fuer das Modell. Damit ist die offene Entscheidung im Kontext geschlossen.
- **2026-10-11 19:28 · Alexander Sacharov** — Offene Entscheidungen gezaehlt (Worker, 11.10.): die einzige, die nur Alex treffen konnte (Modell sieht die Kurzzeile), ist laut Notiz beantwortet. Selbst entschieden: (1) Ort: core/mod/builtin/jaira-cockpit/ im Repo, neben core/role/builtin, weil NZZDRH den Mod im Binary mitliefern will; NICHT ~/.claude/dev-mods, das liegt ausserhalb des Worktrees. (2) Gefiltert werden nur Lieferungen an die Hauptkonversation (kein agentId) mit origin task-notification, peer, peer-send-message, coordinator und ohne e.event (e.event = GitHub/Artifact-Kommentar -> durchlassen); bridge (Mensch per Remote Control), slack, scheduled-trigger, projects-relay, unclassified bleiben unberuehrt. (3) Id: 6-stelliger Crockford-Handle mit mindestens einer Ziffer, oder ein beliebiger Handle direkt hinter 'ticket'/'jaira <verb>'/'fix(' usw., oder eine 26-stellige ULID; mehrere Ids -> die haeufigste. (4) Lane-Wechsel, Urteil, Frage werden aus dem Text gelesen, nicht per jaira show: der Mod hat keinen verlaesslichen Worktree und jaira show kostet einen Prozess pro Lieferung. (5) Speicher: $.store (sitzungsuebergreifend, 4 MiB Grenze), pro Ticket die letzten 20 Berichte, je hoechstens 16000 Zeichen. (6) Fuer das Modell zusaetzlich ein Werkzeug mcp__jaira-cockpit__flight mit derselben Ausgabe wie /flight, weil das Modell Slash-Befehle eines Mods nicht selbst ausfuehren kann.
- **2026-10-11 19:31 · Alexander Sacharov** — Gelernt: (a) im Test-Kit nimmt $.command.run das volle CommandRunInput (origin + presentation), der Hook bekam args=undefined, als ein Plugin ohne args aufrief -> e.args ?? ''. (b) Unter dem Plugin muss der Test selbst einen session.receive-Hook stellen, der {text: e.text} antwortet, sonst wirft der Boden. (c) Typpruefung: die Typen liegen erst nach dem ersten Laden unter .claude-plugin/types; bis dahin tsc mit einer tsconfig ausserhalb des Mod-Ordners, die die claude-code.d.ts des Skills nennt (allowImportingTsExtensions noetig wegen import './digest.ts'). (d) Keine Zeile in core/release/NOTES.md: der Mod ist nur eine Quelle im Repo, jaira installiert ihn noch nicht - das bringt NZZDRH mit seiner eigenen NOTES-Zeile. (e) Nicht live ausprobiert in einer Teamlead-Sitzung; ob ein Hintergrund-Agent seinen Abschluss wirklich als session.receive origin task-notification liefert, sagen nur die Typen (PromptOrigin 'task-notification' = 'A background task's notification').
- **2026-10-11 19:32 · Alexander Sacharov** — critique (Pass 1, ganzer Diff gelesen): nichts zu aendern. Stehen gelassen, mit Grund: digest.ts getrennt von register.ts, weil die Regeln rein sind und ohne Engine lesbar bleiben; LANES fest verdrahtet statt aus .jaira/lanes gelesen, weil der Mod keinen sicheren Repo-Pfad kennt und der Pfeil 'a -> b' umbenannte Lanes ohnehin liest; /flight ohne Id listet die Tickets, sonst antwortet der Befehl auf ein vergessenes Argument mit 'No reports kept for .'; immediate:true, weil der Teamlead meist mitten im Turn ist. Bekannte Grenze: ein Bericht, der SHA256 oefter nennt als seinen Handle, wird SHA256 zugeordnet (6 Zeichen mit Ziffer) - selten, und die Zeile verweist dann auf /flight SHA256, nichts geht verloren.
- **2026-10-11 19:32 · Alexander Sacharov** — optimize: VERDICT_WORD in hooks/digest.ts hat 'approve mit vorbehalt' als eigene Alternative HINTER 'approve[sd]?...'; die Alternation nimmt die erste, also wird aus 'approve mit Vorbehalt' nur 'approve' - genau der Vorbehalt aus dem Ziel-Beispiel geht verloren. Zusammenlegen ist kein Aufraeumen, sondern aendert Verhalten -> zurueck nach in-progress: 'mit vorbehalt' als Zusatz in die approve-Alternative falten, Test dazu.
